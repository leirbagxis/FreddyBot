# Plano: final-hardening-producao_2026-08-14_19-50

## Pedido do usuário
Finalizar o hardening de produção do FreddyBot abordando os 7 pontos específicos solicitados pelo usuário:
1. **Concorrência do Rate Limiter em Memória**: Eliminar Data Race no `memoryLimiterEntry` em `internal/api/middleware/rate_limit.go` usando sincronização thread-safe (`sync.Mutex`) no gerenciador em memória e implementar limpeza periódica oportunística de entradas expiradas sem goroutines extras.
2. **Correção do Endpoint `/readyz`**: Atualizar `internal/api/controllers/healthController.go` para que o `/readyz` retorne obrigatoriamente `503 Service Unavailable` se o banco de dados OU o Redis estiverem indisponíveis (`redisStatus = "down"` $\rightarrow$ `isReady = false`).
3. **Validação do Secret do Webhook**: Adicionar validação estrita dos caracteres permitidos pela Telegram Bot API (`A-Z, a-z, 0-9, _, -`) para `TELEGRAM_WEBHOOK_SECRET` em `pkg/config/config.go`.
4. **Compatibilidade de Criptografia MTProto**: Preservar fallback de `MTPROTO_ENCRYPTION_KEY` para `SECRET_KEY` garantindo retrocompatibilidade de sessões e sem expor credenciais em logs.
5. **Auditoria de Lifecycle & Goroutines**: Garantir que workers/goroutines de longa duração utilizem contexto derivado sem vazamentos ou chamadas desnecessárias a `context.Background()`.
6. **Verificação Completa de Produção**: Executar `gofmt`, `go test ./...`, `go test -race ./...`, `go vet ./...`, `go build ./cmd/FreddyBot` e build do frontend (`npm run build`).
7. **Veredito de Produção**: Emitir veredito final sem realizar commits/pushes.

## Objetivo
Garantir que a aplicação FreddyBot seja 100% thread-safe no teste concorrente de race conditions (`-race`), com o endpoint de prontidão `/readyz` refletindo dependências reais (DB + Redis), validações de webhook mais rígidas e zero regressões.

## Contexto atual
- `rate_limit.go` utilizava `sync.Map`, porém mutava `entry.count++` sem trava interna na struct `memoryLimiterEntry`, gerando data race em chamadas simultâneas.
- `healthController.go` marcava `redisStatus = "down"`, mas não marcava `isReady = false` quando o Redis falhava.
- `config.go` validava apenas o tamanho do `TELEGRAM_WEBHOOK_SECRET`, mas não validava os caracteres permitidos pela Telegram Bot API.

## Arquivos analisados
- `internal/api/middleware/rate_limit.go`
- `internal/api/middleware/rate_limit_test.go` *(a ser criado)*
- `internal/api/controllers/healthController.go`
- `internal/api/controllers/healthController_test.go`
- `pkg/config/config.go`
- `pkg/config/config_test.go`
- `internal/telegram/mtproto/encryption/encryption.go`

## Arquivos que poderão ser modificados / criados
- `internal/api/middleware/rate_limit.go`
- `internal/api/middleware/rate_limit_test.go` *(novo)*
- `internal/api/controllers/healthController.go`
- `internal/api/controllers/healthController_test.go`
- `pkg/config/config.go`
- `pkg/config/config_test.go`

## Estratégia de implementação

1. **`internal/api/middleware/rate_limit.go`**:
   - Refatorar a estrutura do rate-limiter em memória para usar uma struct gerenciadora `memoryLimiter` com um `sync.Mutex` encapsulando um mapa `map[string]*memoryLimiterEntry`.
   - Garantir que a leitura e a mutação de `count` ocorram dentro do bloco protegido por trava `m.mu.Lock() / defer m.mu.Unlock()`.
   - Adicionar limpeza oportunística de entradas expiradas dentro da trava a cada intervalo configurado (ex: 1 minuto).
   - Criar `rate_limit_test.go` com chamadas concorrentes paralelas (`t.RunParallel`) e validar com `go test -race ./internal/api/middleware/...`.

2. **`internal/api/controllers/healthController.go` & `healthController_test.go`**:
   - Atualizar a lógica do `/readyz`: se `cache.HealthCheck(ctx) != nil`, definir `redisStatus = "down"` e `isReady = false`.
   - Atualizar/expandir `healthController_test.go` para testar `/healthz` e `/readyz` nos dois cenários (pronto com DB + Redis ok vs. não pronto quando uma dependência falhar).

3. **`pkg/config/config.go` & `config_test.go`**:
   - Adicionar função auxiliar `isValidWebhookSecret` verificando caracteres válidos (`A-Z, a-z, 0-9, _, -`).
   - Adicionar checagem em `Validate()` para rejeitar secrets contendo caracteres inválidos.
   - Atualizar `config_test.go` testando secrets com caracteres inválidos (ex: espaços ou símbolos especiais).

4. **Verificação de Produção**:
   - Executar `gofmt -w .`.
   - Executar `go test ./...`.
   - Executar `go test -race ./...`.
   - Executar `go vet ./...`.
   - Executar `go build ./cmd/FreddyBot`.
   - Executar build do frontend: `cd dashboard && npm ci && npm run build`.

## Passos detalhados
1. Criar o arquivo de plano `.agent/plans/pending/final-hardening-producao_2026-08-14_19-50.md` e solicitar aprovação do usuário.
2. Refatorar `rate_limit.go` com mutex thread-safe e adicionar `rate_limit_test.go`.
3. Ajustar `healthController.go` para falhar o `/readyz` se o Redis estiver down e atualizar `healthController_test.go`.
4. Adicionar validação de caracteres do webhook secret em `config.go` e atualizar `config_test.go`.
5. Executar toda a suíte de verificação (`go test -race`, `go vet`, `go build`, `npm run build`).
6. Reportar resultados com veredito final (sem commit/push).

## Riscos
- Mínimo. Alterações extremamente pontuais e isoladas focadas em segurança de memória e conformidade de API.

## Compatibilidade
- Linux, macOS, Windows, Docker, CI/CD.

## Como testar

### Testes de Concorrência
```bash
go test -race ./internal/api/middleware/...
```

### Testes da Aplicação
```bash
go test -race ./...
go vet ./...
```

### Build Backend
```bash
go build ./cmd/FreddyBot
```

### Build Frontend
```bash
cd dashboard && npm ci && npm run build
```

## Rollback
`git reset --hard 69ce909`
