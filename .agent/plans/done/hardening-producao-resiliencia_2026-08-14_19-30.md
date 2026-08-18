# Plano: hardening-producao-resiliencia_2026-08-14_19-30

## Pedido do usuário
Implementar as correções remanescentes para prontidão de produção no repositório `leirbagxis/FreddyBot` (branch `feat/mtproto-custom-emoji-separator`), com foco em:
1. Graceful shutdown propagando o contexto raiz e parando serviços limpos na ordem correta;
2. Validação rigorosa de configuração no boot em `pkg/config`;
3. Migrações de banco de dados seguras com checagem e propagação de erros em DDLs/DMLs;
4. Auditoria e correção de contextos desvinculados, erros ignorados, nil pointers e goroutines sem `recover()`;
5. Fallback em memória para rate-limiting em endpoints sensíveis (ex: `/api/login`) caso o Redis esteja indisponível;
6. Endpoints de saúde e prontidão unauthenticated (`GET /healthz` e `GET /readyz`);
7. Suporte a `MTPROTO_ENCRYPTION_KEY` dedicada com fallback seguro e warning;
8. Revisão e validação do `Dockerfile` de produção;
9. Execução e adição de testes unitários automatizados para validação de config, health, webhook e rate limiter.

## Objetivo
Garantir que a aplicação FreddyBot seja 100% resiliente, segura e confiável para ambiente de produção (PostgreSQL + Redis + Telegram Bot API + MTProto), sem regressões de funcionalidades existentes.

## Contexto atual
- A aplicação possui um contexto raiz (`signal.NotifyContext`) em `main.go`, mas alguns serviços e goroutines (long polling do Telegram, goroutines de background, timers) criam contextos desvinculados (`context.Background()`).
- O parser de rate-limit 429 na Bot API (`botapi.go`) lia o número 429 como tempo de retry.
- O GORM em `database.go` executava `db.Exec(...)` ignorando erros.
- A validação de configuração em `config.go` era parcial (não validava comprimento mínimo do `SECRET_KEY`, formato do webhook secret ou credenciais MTProto parciais).
- Não existiam endpoints padrão de infraestrutura (`/healthz` e `/readyz`).

## Arquivos analisados
- `cmd/FreddyBot/main.go`
- `pkg/config/config.go`
- `internal/database/database.go`
- `internal/cache/redis.go`
- `internal/api/api.go`
- `internal/api/routes/routes.go`
- `internal/container/appContainer.go`
- `internal/telegram/client.go`
- `internal/telegram/executor/botapi.go`
- `internal/telegram/events/channelPost/stage_preflight_telego.go`
- `internal/telegram/events/channelPost/dispatch_telego.go`
- `internal/telegram/events/channelPost/channelPost.go`
- `internal/core/services/scheduler.go`
- `internal/api/controllers/accountController.go`
- `internal/api/controllers/userController.go`
- `internal/api/controllers/adminController/getAllUserAdminController.go`
- `internal/telegram/mtproto/encryption/encryption.go`
- `Dockerfile`

## Arquivos que poderão ser modificados / criados
- `cmd/FreddyBot/main.go`
- `pkg/config/config.go`
- `pkg/config/config_test.go` *(novo)*
- `internal/database/database.go`
- `internal/cache/redis.go`
- `internal/api/api.go`
- `internal/api/routes/routes.go`
- `internal/api/controllers/healthController.go` *(novo)*
- `internal/api/controllers/healthController_test.go` *(novo)*
- `internal/api/middleware/rate_limit.go` *(novo/ajustado)*
- `internal/container/appContainer.go`
- `internal/telegram/client.go`
- `internal/telegram/executor/botapi.go`
- `internal/telegram/events/channelPost/stage_preflight_telego.go`
- `internal/telegram/events/channelPost/dispatch_telego.go`
- `internal/telegram/events/channelPost/channelPost.go`
- `internal/core/services/scheduler.go`
- `internal/api/controllers/accountController.go`
- `internal/api/controllers/userController.go`
- `internal/api/controllers/adminController/getAllUserAdminController.go`
- `internal/telegram/mtproto/encryption/encryption.go`

## Estratégia de implementação

1. **Graceful Shutdown & Pipeline Context Handshake**:
   - Refatorar `telegram.StartBot` para aceitar o `context.Context` raiz vindo de `main.go`.
   - Propagar esse contexto para long polling (`UpdatesViaLongPolling(ctx, ...)`), chamadas de boot da Bot API (`GetMe(ctx)`, `SetWebhook(ctx, ...)`) e para o handler do bot.
   - Guardar/expor o `telegohandler.BotHandler` para invocar `bh.Stop()` durante o encerramento.
   - No `main.go`, usar `sync.WaitGroup` para controlar as goroutines principais (API e bot) e encerrá-las na ordem correta:
     1. Recebe sinal (SIGINT/SIGTERM) -> cancela `ctx`.
     2. Para recepção de novos updates/webhooks (`bh.Stop()`).
     3. Para background workers e schedulers.
     4. Encerra servidor HTTP (`srv.Shutdown()`).
     5. Fecha conexão com Redis.
     6. Fecha pool do banco de dados (PostgreSQL/SQLite).

2. **Reforço de Validação de Configuração (`pkg/config`)**:
   - Atualizar `config.Validate()` e os getters para validar:
     - `SECRET_KEY`: não vazia e comprimento >= 32 bytes.
     - `OWNER_ID`: > 0.
     - `TELEGRAM_BOT_TOKEN`: não vazio.
     - `REDIS_HOST`: não vazio.
     - `DATABASE_FILE` (DSN): obrigatório quando em produção ou com driver Postgres.
     - `APP_ENV`: `"dev"` ou `"prod"`.
     - `WEBHOOK_URL` x `TELEGRAM_WEBHOOK_SECRET`: se webhook ativo, secret é obrigatório e precisa ter comprimento mínimo (e.g. >= 8 caracteres).
     - Credenciais MTProto (`MTPROTO_APP_ID` e `MTPROTO_APP_HASH`): se um estiver presente, o outro também deve estar.
   - Adicionar variável opcional `MTPROTO_ENCRYPTION_KEY`.
   - Criar `pkg/config/config_test.go` para cobrir todos esses cenários de teste.

3. **Migrações Confiáveis de Banco (`internal/database/database.go`)**:
   - Garantir que cada instrução `db.Exec(...)` verifique `.Error` e retorne erro descritivo em caso de falha.
   - Adicionar logs claros das etapas de migração.

4. **Correções de Bugs P0/P1, Nil Checks e Recover**:
   - **Rate Limit Parser (`botapi.go`)**: Corrigir `extractRetryAfter` usando busca direcionada após `"retry after"` ou expressão regular, ignorando o código HTTP `429`.
   - **Nil Pointer em Preflight (`stage_preflight_telego.go`)**: Adicionar verificação `if botInfo == nil` após `GetMe`.
   - **Nil Pointer em Transferência de Canal (`userController.go`)**: Adicionar verificação `if botInfo == nil`.
   - **Type Assertions (`accountController.go`)**: Substituir `userID.(int64)` por `auth.GetUserID(ctx)` com validação segura.
   - **Goroutines soltas**: Adicionar `defer` com `recover()` em `dispatchNotice`, `sendScheduledPost`, `time.AfterFunc` no separador e nos workers da `MessageQueue`.
   - **Redis Concorrência**: Tornar a inicialização do client do Redis resiliente a falhas temporárias sem memorizar `nil` para sempre.

5. **Rate Limiter Fallback**:
   - Criar/ajustar middleware de rate limit em memória para a rota de login (`POST /api/login`) usando um bucket concorrente (baseado em `sync.Map` com expiração de janelas). Se o Redis falhar, esse limiter em memória assume sem derrubar a API.

6. **Endpoints de Health & Readiness (`/healthz` e `/readyz`)**:
   - Criar `HealthController` e rotas:
     - `GET /healthz`: retorna `200 OK` com `{"status":"ok"}` sem dependências externas.
     - `GET /readyz`: testa o ping no banco de dados e no Redis; se um falhar, retorna `503 Service Unavailable` com JSON sanitizado (sem expor credenciais).

7. **Chave de Criptografia MTProto (`encryption.go`)**:
   - Atualizar `encryption.go` para buscar `config.MTProtoEncryptionKey`. Se vazia, utilizar `config.SecretKey` emitindo log de aviso.

8. **Dockerfile**:
   - Validar se a imagem usa multi-stage, usuário não-root `appuser`, e dependências ca-certificates em Alpine.

9. **Execução e Verificação de Testes**:
   - `gofmt -w .`
   - `go test ./...`
   - `go vet ./...`
   - `go build ./cmd/FreddyBot`
   - Dashboard: `cd dashboard && npm ci && npm run build`

## Passos detalhados

1. Criar este plano em `.agent/plans/pending/hardening-producao-resiliencia_2026-08-14_19-30.md` e obter aprovação do usuário.
2. Atualizar `pkg/config/config.go` com regras estritas de validação + `MTPROTO_ENCRYPTION_KEY` e criar `pkg/config/config_test.go`.
3. Ajustar `internal/database/database.go` para propagar erros em todas as chamadas `db.Exec`.
4. Corrigir o parser em `internal/telegram/executor/botapi.go` (`extractRetryAfter`).
5. Adicionar verificações de nil pointer em `stage_preflight_telego.go` e `userController.go`.
6. Adicionar `recover()` e propagação de contexto nas goroutines do `scheduler.go`, `getAllUserAdminController.go` e `dispatch_telego.go`.
7. Ajustar `accountController.go` para extração segura de `userID`.
8. Refatorar `internal/telegram/client.go` e `cmd/FreddyBot/main.go` para garantir o ciclo de vida e graceful shutdown sincronizado de todos os componentes.
9. Implementar `HealthController` (`/healthz` e `/readyz`) e registrar rotas em `routes.go`.
10. Atualizar `internal/telegram/mtproto/encryption/encryption.go` para suportar `MTPROTO_ENCRYPTION_KEY`.
11. Implementar o rate limiter fallback em memória para a rota de login.
12. Revisar e validar o `Dockerfile`.
13. Rodar `gofmt`, `go test ./...`, `go vet ./...`, `go build` e `npm run build` para validar 100%.

## Riscos
- Risco de quebrar testes existentes se a validação de config em `init()` barrar o ambiente de teste: contornaremos verificando `isTestMode()` adequadamente.
- Risco de travamento no shutdown se algum worker de background não responder ao `ctx.Done()`: usaremos timeouts curtos com `sync.WaitGroup`.

## Compatibilidade
- Linux
- macOS
- Windows
- Docker
- CI/CD

## Como testar

### Build Backend
```bash
go build ./cmd/FreddyBot
```

### Testes Go
```bash
go test ./...
go vet ./...
```

### Build Frontend
```bash
cd dashboard && npm ci && npm run build
```

## Rollback
`git reset --hard 405a4cf`

## Observações
Toda a refatoração respeita a arquitetura existente sem alterar o framework Web (Gin), ORM (GORM) ou bot engine (Telego).
