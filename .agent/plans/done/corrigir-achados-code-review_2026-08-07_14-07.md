# Plano: Correção dos Achados do Code Review

## Pedido do usuário
Corrigir todos os 65 achados identificados no code review do FreddyBot antes de levar o código para produção.

## Objetivo
Resolver bugs críticos, vulnerabilidades de segurança, problemas de performance e melhorias de qualidade identificados na auditoria de código de 2026-08-07.

## Contexto atual
- O FreddyBot é um bot Telegram em Go com API REST, dashboard React, cache Redis e banco SQLite/PostgreSQL.
- Foram encontrados 65 problemas: 6 CRITICAL, 23 HIGH, 25 MEDIUM, 11 LOW/INFO.
- O relatório completo está em `docs/CODE_REVIEW_2026-08-07.md`.

## Arquivos analisados
- `cmd/FreddyBot/main.go`
- `internal/container/container.go`
- `internal/cache/redis.go`
- `internal/cache/cache.go`
- `internal/cache/ranking_cache.go`
- `internal/telegram/bot.go`
- `internal/telegram/handlers.go`
- `internal/telegram/commands.go`
- `internal/telegram/events/channelPost/channelPost.go`
- `internal/core/message_service.go`
- `internal/core/sticker_service.go`
- `internal/core/ranking_service.go`
- `internal/core/group_service.go`
- `internal/core/services.go`
- `internal/core/services/scheduler.go`
- `internal/database/db.go`
- `internal/database/queries.go`
- `internal/database/repository.go`
- `internal/api/router.go`
- `internal/api/handlers.go`
- `internal/api/middleware/rate_limit.go`
- `internal/api/auth/middleware.go`
- `internal/middleware/checkAddBotMiddlewareTelego.go`
- `internal/utils/formatter.go`
- `pkg/config/config.go`
- `pkg/logger/logger.go`
- `pkg/errors/errors.go`
- `pkg/parser/parser.go`
- `Dockerfile`
- `docker-compose.yml`
- `.env-example`

## Arquivos que poderão ser modificados
Todos os listados acima.

## Estratégia de implementação

O plano é dividido em **5 fases** que devem ser executadas em ordem, cada uma gerando commits isolados e testáveis. Cada fase agrupa correções por tema e prioridade, para minimizar risco de regressão.

---

## Fase 1 — CRITICAL: Bugs bloqueantes e crashers (itens 1-5, 51)

> **Prioridade:** Bloqueia deploy  
> **Risco:** Alto — toca em startup, shutdown e inicialização global  
> **Estimativa:** 2-3 horas

### Passo 1.1 — Corrigir race condition no shutdown (item 1)
**Arquivo:** `cmd/FreddyBot/main.go`
- Remover goroutine anônima competindo pelo canal `stop`
- Usar `signal.NotifyContext` para integrar sinais com `context.Context`
- Garantir shutdown sequencial: cancelar contexto → aguardar bot/API pararem → chamar `container.Shutdown()`

### Passo 1.2 — Adicionar sincronização ao container (item 2)
**Arquivo:** `internal/container/container.go`
- Proteger `Initialize` com `sync.Once`
- Adicionar `sync.RWMutex` para proteger acesso às variáveis globais
- Adicionar verificação de nil nos getters com mensagem de panic clara (item 32)

### Passo 1.3 — Adicionar nil checks no Shutdown (item 7)
**Arquivo:** `internal/container/container.go`
- Adicionar `if redisClient != nil` antes de `redisClient.Close()`
- Adicionar `if db != nil` antes de `db.Close()`

### Passo 1.4 — Implementar cleanup em cascata no Initialize (itens 6, 8)
**Arquivo:** `internal/container/container.go` e `cmd/FreddyBot/main.go`
- Se um passo de inicialização falhar, fechar recursos já abertos antes de retornar
- Adicionar `defer container.Shutdown()` no `main.go` após `Initialize`

### Passo 1.5 — Validar conexão Redis com Ping (item 3)
**Arquivo:** `internal/cache/redis.go`
- Adicionar `client.Ping(ctx).Err()` imediatamente após criar o cliente Redis
- Se falhar, retornar erro claro

### Passo 1.6 — Corrigir vazamento de memória no rate limiting (item 51)
**Arquivo:** `internal/api/middleware/rate_limit.go`
- Substituir `c.Request.Context()` por `context.Background()` na chamada `Expire`

### Passo 1.7 — Validar configurações obrigatórias (item 26)
**Arquivo:** `pkg/config/config.go`
- Adicionar função `Validate()` que verifica campos obrigatórios
- Falhar rápido no startup com mensagem descritiva se faltar `BOT_TOKEN`, `DATABASE_URL`, etc.

**Commit:** `fix(critical): resolve race conditions, resource leaks, and startup validation`

---

## Fase 2 — SECURITY: Vulnerabilidades de segurança (itens 5, 12, 14, 17, 19, 20, 25, 27, 28, 29, 30, 47, 59)

> **Prioridade:** Segurança  
> **Risco:** Médio — afeta configuração e infraestrutura  
> **Estimativa:** 2-3 horas

### Passo 2.1 — Proteger portas no docker-compose (item 59)
**Arquivo:** `docker-compose.yml`
- Alterar `6379:6379` para `127.0.0.1:6379:6379`
- Alterar `5432:5432` para `127.0.0.1:5432:5432`
- Remover senha padrão `12345` e usar referência a variável de ambiente

### Passo 2.2 — Remover credenciais hardcoded (item 29)
**Arquivo:** `docker-compose.yml`
- Substituir credenciais por `env_file: .env`
- Documentar no `.env-example`

### Passo 2.3 — Hardening do Dockerfile (itens 27, 28)
**Arquivo:** `Dockerfile`
- Fixar versão da imagem base (ex: `golang:1.24-alpine`)
- Usar imagem final mínima (`gcr.io/distroless/static-debian12` ou `alpine`)
- Adicionar `RUN adduser -D -u 1000 appuser` e `USER appuser`
- Adicionar `.dockerignore` com `.env`, `.git`, `*.png`, etc.

### Passo 2.4 — Sanitizar nomes de usuário (item 25)
**Arquivo:** `internal/utils/formatter.go`
- Criar função `escapeHTML(s string)` usando `html.EscapeString()`
- Aplicar em todos os pontos onde nomes de usuário são inseridos em mensagens

### Passo 2.5 — Mascarar credenciais nos logs (itens 12, 17)
**Arquivos:** `internal/database/db.go`, `internal/telegram/bot.go`
- Garantir que DSN nunca seja logado
- Mascarar token do bot em qualquer log (mostrar apenas últimos 4 caracteres)

### Passo 2.6 — Restringir CORS (item 20)
**Arquivo:** `internal/api/router.go`
- Substituir `*` por lista de origens confiáveis via variável de ambiente `CORS_ORIGINS`

### Passo 2.7 — Usar comparação segura de tokens (item 30)
**Arquivos:** `internal/middleware/`, `internal/api/auth/`
- Substituir `==` por `crypto/subtle.ConstantTimeCompare()` em todas as comparações de tokens

### Passo 2.8 — Sanitizar erros da API (item 47)
**Arquivo:** `internal/api/handlers.go`
- Retornar mensagens genéricas para o cliente (`"Internal server error"`)
- Logar detalhes internos com `logger.Error()`

### Passo 2.9 — Auditar queries contra SQL injection (item 14)
**Arquivo:** `internal/database/queries.go`
- Verificar todas as queries usam placeholders
- Documentar resultado da auditoria

### Passo 2.10 — Verificar .gitignore (item 5)
**Arquivo:** `.gitignore`
- Confirmar que `.env` está incluído
- Adicionar `.env.local`, `.env.production` se necessário

**Commit:** `security: harden docker, sanitize inputs, protect credentials, restrict CORS`

---

## Fase 3 — STABILITY: Estabilidade e resiliência (itens 4, 9, 10, 13, 15, 16, 18, 21, 22, 23, 24, 41, 52, 53, 54)

> **Prioridade:** Estabilidade  
> **Risco:** Médio — afeta runtime e concorrência  
> **Estimativa:** 3-4 horas

### Passo 3.1 — Corrigir `sync.Once` do bot middleware (item 52)
**Arquivo:** `internal/middleware/checkAddBotMiddlewareTelego.go`
- Substituir `sync.Once` por carregamento no startup com retry
- Ou implementar lazy init com verificação de nil e retry

### Passo 3.2 — Cache no middleware de autenticação (item 53)
**Arquivo:** `internal/api/auth/middleware.go`
- Adicionar cache Redis para resultado de `GetUserByID`
- TTL de 5 minutos
- Invalidar na alteração de permissões

### Passo 3.3 — Corrigir starvation dos workers (item 54)
**Arquivo:** `internal/telegram/events/channelPost/channelPost.go`
- Substituir `time.Sleep` bloqueante por requeue assíncrono
- Worker libera e job volta à fila com delay

### Passo 3.4 — Configurar pool de conexões DB (item 13)
**Arquivo:** `internal/database/db.go`
- Adicionar após `sql.Open()`:
  ```go
  db.SetMaxOpenConns(25)
  db.SetMaxIdleConns(5)
  db.SetConnMaxLifetime(5 * time.Minute)
  ```
- Tornar valores configuráveis via env vars

### Passo 3.5 — Configurar pool Redis (item 23)
**Arquivo:** `internal/cache/redis.go`
- Adicionar `PoolSize: 10`, `MinIdleConns: 5`, `MaxRetries: 3`

### Passo 3.6 — Adicionar log de shutdown do bot (item 4)
**Arquivo:** `internal/telegram/bot.go`
- Adicionar log quando contexto é cancelado
- Aguardar handlers ativos finalizarem

### Passo 3.7 — Validação de data com fallback (itens 9, 10)
**Arquivos:** `internal/core/message_service.go`, `internal/core/sticker_service.go`
- Se `messageDate.IsZero()`, usar `time.Now()` como fallback
- Logar warning quando fallback é usado

### Passo 3.8 — Invalidação de cache de ranking (item 24)
**Arquivo:** `internal/cache/ranking_cache.go`
- Invalidar cache quando nova mensagem é salva
- Ou reduzir TTL para 2-3 minutos

### Passo 3.9 — Validar parâmetros da API (item 21)
**Arquivo:** `internal/api/handlers.go`
- Validar IDs positivos, datas válidas, strings não vazias

### Passo 3.10 — Logar erros de Redis (item 22)
**Arquivo:** `internal/cache/redis.go`
- Implementar padrão cache-aside: logar erros e continuar sem cache

### Passo 3.11 — Garantir `defer rows.Close()` (item 15)
**Arquivo:** `internal/database/repository.go`
- Verificar todos os pontos que obtêm `rows`
- Garantir `defer rows.Close()` imediato
- Adicionar verificação de `rows.Err()` após loops (item 41)

### Passo 3.12 — Logar perda de mensagens (item 18)
**Arquivo:** `internal/telegram/handlers.go`
- Garantir que erros de `SaveMessage` são logados com contexto (chatID, userID)

### Passo 3.13 — Verificar reconexão do bot (item 16)
**Arquivo:** `internal/telegram/bot.go`
- Verificar se a lib `telego` tem reconexão automática
- Se não, implementar wrapper com retry + backoff

**Commit:** `fix(stability): connection pools, cache invalidation, worker starvation, error handling`

---

## Fase 4 — PERFORMANCE: Otimizações de performance (itens 11, 35, 37, 42, 43, 55, 56, 57)

> **Prioridade:** Performance  
> **Risco:** Baixo-médio — otimizações podem alterar comportamento  
> **Estimativa:** 2-3 horas

### Passo 4.1 — Corrigir scheduler síncrono (item 55)
**Arquivo:** `internal/core/services/scheduler.go`
- Disparar `sendScheduledPost()` em goroutines separadas
- Usar semaphore para limitar concorrência
- Aplicar rate-limit por canal, não global

### Passo 4.2 — Corrigir tipo inconsistente no cache (item 56)
**Arquivo:** `internal/cache/cache.go`
- Padronizar: serializar com `json.Marshal` antes de inserir no `localCache`
- Ou criar métodos tipados (`SetChannel`/`GetChannel`) que não passem pelo `Get` genérico

### Passo 4.3 — Substituir deep copy JSON por Clone() (item 57)
**Arquivo:** `internal/cache/cache.go` + modelo `Channel`
- Implementar `Clone()` na struct `models.Channel`
- Substituir `json.Marshal/Unmarshal` por chamada a `Clone()`

### Passo 4.4 — Usar timezone explícito no ranking (item 11)
**Arquivo:** `internal/core/ranking_service.go`
- Substituir `time.Now()` por `time.Now().UTC()` ou timezone configurável
- Documentar a decisão

### Passo 4.5 — Adicionar limites ao ranking global (item 37)
**Arquivo:** `internal/core/ranking_service.go`
- Adicionar `LIMIT` à query de ranking global
- Ou usar paginação

### Passo 4.6 — Refatorar código duplicado do ranking (item 38)
**Arquivo:** `internal/core/ranking_service.go`
- Criar função `getRanking(startDate time.Time)` usada pelas três variantes

### Passo 4.7 — Paralelizar queries do ranking (item 35)
**Arquivo:** `internal/core/message_service.go`
- Executar queries de top senders e top sticker senders em goroutines paralelas
- Usar `errgroup` para coordenar

### Passo 4.8 — Adicionar índices ao banco (item 40)
**Arquivo:** migration ou `internal/database/`
- Criar índices em `(chat_id, created_at)`, `(user_id)`, `(chat_id, user_id)`

**Commit:** `perf: optimize scheduler, cache, ranking queries, and database indexes`

---

## Fase 5 — QUALITY: Qualidade de código (itens 31, 33, 34, 36, 39, 44, 45, 46, 48, 49, 50, 58, 60, 61, 62-65)

> **Prioridade:** Qualidade / pós-deploy  
> **Risco:** Baixo — melhorias de manutenibilidade  
> **Estimativa:** 3-4 horas

### Passo 5.1 — Extrair metadata helper (item 36)
**Arquivos:** `internal/core/message_service.go`, `internal/core/sticker_service.go`
- Criar `extractMessageMetadata(msg) (userID, userName int64, chatID int64, date time.Time)`
- Usar em ambos os serviços

### Passo 5.2 — Validar tipo de chat nos comandos (item 45)
**Arquivo:** `internal/telegram/commands.go`
- Adicionar verificação de tipo de chat para comandos de grupo
- Retornar mensagem amigável se usado em chat privado

### Passo 5.3 — Validar chatID em SaveGroup (item 39)
**Arquivo:** `internal/core/group_service.go`
- Validar que chatID é negativo (grupo Telegram)

### Passo 5.4 — Centralizar mensagens no messages.yml (item 46)
**Arquivo:** `internal/telegram/commands.go`, `config/messages.yml`
- Migrar strings hardcoded para `messages.yml`

### Passo 5.5 — Parametrizar limite do GetTopSenders (item 34)
**Arquivo:** `internal/core/message_service.go`
- Aceitar `limit int` como parâmetro

### Passo 5.6 — Garantir chatID na chave de cache (item 49)
**Arquivo:** `internal/cache/ranking_cache.go`
- Verificar e corrigir chaves de cache para incluir chatID

### Passo 5.7 — Circuit breaker para Redis (item 48)
**Arquivo:** `internal/cache/redis.go`
- Implementar circuit breaker simples (contador de falhas + cooldown)

### Passo 5.8 — Corrigir nomenclatura Redis (item 61)
**Arquivos:** `internal/cache/redis.go`, `.env-example`
- Renomear para `REDIS_URL` ou ajustar código para aceitar endereço simples

### Passo 5.9 — Usar go:embed para messages.yml (item 60)
**Arquivo:** `pkg/parser/parser.go`
- Substituir `os.ReadFile("config/messages.yml")` por `//go:embed`

### Passo 5.10 — Melhorias menores (itens 50, 58, 62-65)
- Formatação de números (item 50)
- Fix `generateShortID` (item 58)
- Mover canal `stop` para local (item 62)
- Logger JSON em produção (item 63)
- `Unwrap()` nos erros (item 64)
- Documentar `.env-example` (item 65)

**Commit:** `refactor(quality): DRY, input validation, cache keys, messages centralization`

---

## Riscos
- **Fase 1 (CRITICAL):** Altera fluxo de startup/shutdown — risco de regressão no boot. Requer testes manuais completos.
- **Fase 2 (SECURITY):** Alterações no Docker e CORS podem afetar deploy. Testar em ambiente staging antes de produção.
- **Fase 3 (STABILITY):** Alteração nos workers e pool de conexões pode afetar throughput. Monitorar métricas após deploy.
- **Fase 4 (PERFORMANCE):** Refatoração do ranking e cache pode alterar resultados. Comparar outputs antes/depois.
- **Fase 5 (QUALITY):** Baixo risco — são melhorias incrementais.

## Impactos esperados
- Eliminar todos os crashers e panic em produção
- Proteger contra as vulnerabilidades de segurança mais graves
- Melhorar resiliência do bot contra falhas de rede e banco de dados
- Reduzir carga no banco de dados (~50% menos queries no middleware de auth)
- Melhorar throughput do scheduler e da fila de mensagens
- Código mais limpo, testável e manutenível

## Compatibilidade
- ✅ Linux
- ✅ macOS
- ✅ Docker
- ✅ CI/CD (commits atômicos por fase)
- ⚠️ Windows (verificar paths do `messages.yml` com `go:embed`)

## Como testar

### Build
```bash
cd /home/malbs/Opencode/FreddyBot && go build ./...
```

### Testes
```bash
go test ./... -v -race
```

### Lint
```bash
golangci-lint run ./...
```

### Execução
```bash
go run cmd/FreddyBot/main.go
```

### Docker
```bash
docker-compose build && docker-compose up -d
```

## Rollback
- Cada fase gera um commit separado, permitindo `git revert` granular
- Em caso de regressão crítica: `git revert HEAD~N` (onde N = número de commits da fase)
- Manter branch de backup antes de iniciar: `git branch backup/pre-code-review-fix`

## Observações
- O plano é conservador: prioriza correções cirúrgicas sobre refatorações grandes
- Fases 1-3 devem ser feitas antes do deploy para produção
- Fases 4-5 podem ser feitas após o deploy, de forma incremental
- Recomendado executar `go test -race ./...` após cada fase para detectar race conditions
- O relatório completo está documentado em `docs/CODE_REVIEW_2026-08-07.md`
