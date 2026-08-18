# Plano: mtproto-connected-accounts

## Pedido do usuário
Implementar suporte completo a contas Telegram dos próprios usuários utilizando o protocolo MTProto (biblioteca `gotd/td`), permitindo que usuários utilizem recursos exclusivos de suas próprias contas (emojis Premium, reações Premium, etc.) dentro dos canais previamente autorizados. A funcionalidade deve coexistir com a Bot API existente.

## Objetivo
Criar sistema completo e modular de contas conectadas via MTProto, com:
- Interface `TelegramExecutor` abstraindo BotAPI e MTProto
- Autenticação segura (telefone → código → 2FA)
- Sessões criptografadas no PostgreSQL
- Pool de clientes MTProto
- API REST para gerenciamento
- UI no dashboard React
- Integração automática no pipeline existente

## Contexto atual
- **Stack**: Go + Telego (Bot API) + GORM + PostgreSQL + Redis + Gin
- **Pipeline de channel posts**: 8 estágios (Preflight → Special → MediaGroup → Queue → Transform → Decorate → Send)
- **Dispatch atual**: usa `telego.Bot` diretamente para `EditMessageText`, `EditMessageCaption`, `EditMessageReplyMarkup`, `SendSticker`
- **AppContainer**: DI centralizada com todos os services/repos/cache
- **Dashboard**: React + Vite + TailwindCSS SPA servida pelo Gin
- **Autenticação Web**: JWT via init-data do Telegram
- **Chaves de sessão Redis**: padrão `channel:v2:{channelID}`, etc.

## Arquivos analisados
- `cmd/FreddyBot/main.go` - Entry point, boot flow
- `internal/container/appContainer.go` - DI container
- `internal/telegram/events/channelPost/dispatch_telego.go` - Dispatch atual (Bot API)
- `internal/telegram/events/channelPost/stage_send_telego.go` - Stage de envio
- `internal/telegram/events/channelPost/pipeline_telego.go` - Pipeline abstraction
- `internal/database/models/models.go` - Todos os modelos GORM
- `internal/api/routes/routes.go` - Todas as rotas da API
- `internal/api/controllers/channelController.go` - Padrão de controller
- `internal/api/auth/jwt.go` - Autenticação JWT
- `internal/api/auth/middleware.go` - Middleware de autenticação
- `pkg/config/config.go` - Configurações via env
- `pkg/errors/errors.go` - Sistema de erros
- `internal/cache/cache.go` - Cache L1+L2
- `.env` - Variáveis de ambiente
- `docker-compose.yml` - Serviços Docker
- `go.mod` - Dependências
- `dashboard/src/App.tsx` - Frontend React
- `dashboard/src/api.ts` - API client
- `dashboard/src/types.ts` - Typescript types

## Arquivos que poderão ser modificados

### Novos arquivos
```
internal/telegram/executor/executor.go          # TelegramExecutor interface
internal/telegram/executor/botapi.go             # BotAPIExecutor implementation
internal/telegram/mtproto/auth/auth.go           # MTProto auth flow (phone, code, 2FA)
internal/telegram/mtproto/auth/types.go          # Auth state types
internal/telegram/mtproto/session/session.go     # Session storage interface + impl
internal/telegram/mtproto/encryption/encryption.go # Encrypt/decrypt session data
internal/telegram/mtproto/client/client.go       # MTProto client wrapper
internal/telegram/mtproto/client/pool.go         # Client pool
internal/telegram/mtproto/executor/executor.go   # MTProtoExecutor implementation
internal/telegram/mtproto/worker/worker.go       # Async worker for MTProto actions
internal/telegram/mtproto/errors.go              # MTProto-specific errors
internal/database/models/mtproto_models.go       # ConnectedAccount + ConnectedAccountChannel models
internal/database/repositories/connected_account.go # Repository for connected accounts
internal/core/services/connected_account.go      # Service for connected accounts
internal/core/services/connected_account_test.go # Tests
internal/api/controllers/accountController.go    # Account API controller
internal/api/types/account.go                    # Account request/response types
pkg/config/config_mtproto.go                     # MTProto-specific env vars
dashboard/src/components/ContaTelegramTab.tsx    # Main account tab
dashboard/src/components/ConnectedAccountCard.tsx # Account status card
dashboard/src/components/AuthFlow.tsx            # Auth flow (phone/code/password)
dashboard/src/components/AuthDisclaimer.tsx      # Terms disclaimer
```

### Arquivos existentes modificados
```
cmd/FreddyBot/main.go                             # Init MTProto system
internal/container/appContainer.go                # Add MTProto services
internal/api/routes/routes.go                     # Add account routes
internal/api/auth/middleware.go                   # Maybe no changes needed
internal/telegram/events/channelPost/dispatch_telego.go  # Use TelegramExecutor
internal/telegram/events/channelPost/types.go     # Add executor to context
internal/telegram/events/channelPost/pipeline_telego.go  # Pass executor
internal/database/database.go                     # AutoMigrate new models
pkg/config/config.go                              # Add MTProto env vars
dashboard/src/App.tsx                             # Add ContaTelegram tab
dashboard/src/api.ts                              # Add account API calls
dashboard/src/types.ts                            # Add account types
docker-compose.yml                                # Fix version warning (optional)
.env-example                                      # Add new env vars
```

## Estratégia de implementação

### Abordagem geral

1. **Camada de abstração primeiro** — Criar a interface `TelegramExecutor` e a implementação `BotAPIExecutor` que encapsula as chamadas telego existentes. Isso estabelece o contrato sem quebrar nada.

2. **Infraestrutura MTProto** — Adicionar `gotd/td` como dependência, criar pacotes de auth, session, encryption, client.

3. **Camada de dados** — Modelos GORM, repositório, service para contas conectadas.

4. **API REST** — Endpoints para gerenciamento de conta (connect, verify, password, disconnect).

5. **Frontend** — Componentes React para a UI de contas conectadas.

6. **Integração no pipeline** — Substituir chamadas diretas ao `telego.Bot` pelo `TelegramExecutor` no dispatch.

7. **Implementação MTProto** — `MTProtoExecutor` que usa `gotd/td` para executar as mesmas operações.

### Fases

**FASE 0 — Preparação e dependências**
- Adicionar `go.uber.org/zap` (opcional, gotd usa zap) ou configurar gotd com logger custom
- Adicionar `github.com/gotd/td` ao go.mod
- Adicionar `ENCRYPTION_KEY` ao config
- Remover `version` obsoleto do docker-compose.yml

**FASE 1 — Interface TelegramExecutor**
- Criar `internal/telegram/executor/executor.go` com a interface
- Migrar chamadas de dispatch para usar a interface via BotAPIExecutor
- Nenhuma mudança comportamental — apenas refatoração segura

**FASE 2 — Modelos e repositórios**
- Criar `ConnectedAccount` e `ConnectedAccountChannel` models
- Criar repositório `ConnectedAccountRepository`
- Adicionar AutoMigrate no `database.go`
- Criar `ConnectedAccountService`

**FASE 3 — Criptografia e sessão**
- Criar `encryption` package (AES-256-GCM com ENCRYPTION_KEY)
- Criar `session` package (load/save do PostgreSQL via repository)

**FASE 4 — Autenticação MTProto**
- Criar `auth` package com fluxo: SendCode → SignIn → 2FA
- Estados de autenticação em memória (com TTL via Redis)
- Criar `client` package com wrapper do `gotd/td`

**FASE 5 — API REST**
- Endpoints: `GET /api/account`, `POST /api/account/connect`, `POST /api/account/verify`, `POST /api/account/password`, `DELETE /api/account`
- Registrar rotas

**FASE 6 — Frontend**
- Componente `ContaTelegramTab` com status, conectar, reconectar, desconectar
- Componente `AuthFlow` com passo-a-passo (telefone → código → 2FA)
- Modal de disclaimer com checkboxes de termos
- Integrar no App.tsx como nova tab

**FASE 7 — MTProto Executor**
- Criar `MTProtoExecutor` implementando `TelegramExecutor`
- Suporte inicial: EditMessageText, EditMessageCaption, EditMessageReplyMarkup
- Pool de clientes MTProto
- Worker para execução assíncrona

**FASE 8 — Integração no pipeline**
- Substituir `pCtx.Bot.EditMessageText(...)` por `executor.EditMessage(...)`
- Lógica de seleção automática: se conta conectada → MTProto, senão → BotAPI
- Injetar executor no ProcessingContextTelego

**FASE 9 — Testes e finalização**
- Testes de criptografia
- Testes de repositório
- Testes de fluxo de autenticação (mock)
- Testes de permissões
- Logging de eventos MTProto

## Passos detalhados

### FASE 0: Preparação

1. Adicionar `ENCRYPTION_KEY` ao `pkg/config/config.go` (min 32 bytes, hex)
2. Adicionar env var no `.env-example`
3. Rodar `go get github.com/gotd/td@latest`
4. Rodar `go mod tidy`

### FASE 1: Interface TelegramExecutor

5. Criar `internal/telegram/executor/executor.go`:
```go
type TelegramExecutor interface {
    EditMessage(ctx context.Context, chatID int64, messageID int, text string, parseMode string, keyboard *telego.InlineKeyboardMarkup, opts *EditOptions) error
    EditCaption(ctx context.Context, chatID int64, messageID int, caption string, parseMode string, keyboard *telego.InlineKeyboardMarkup, opts *EditOptions) error
    EditReplyMarkup(ctx context.Context, chatID int64, messageID int, keyboard *telego.InlineKeyboardMarkup) error
    SendSticker(ctx context.Context, chatID int64, stickerID string) error
    DeleteMessage(ctx context.Context, chatID int64, messageID int) error
}

type EditOptions struct {
    DisableLinkPreview bool
    // Future: PremiumEmoji, CustomEmoji, etc.
}
```

6. Criar `internal/telegram/executor/botapi.go` com `BotAPIExecutor` implementando a interface usando `telego.Bot`

7. Atualizar `ProcessingContextTelego` em `types.go` para incluir `Executor TelegramExecutor`

8. Atualizar o pipeline para passar o executor do container

### FASE 2: Modelos e repositórios

9. Criar `internal/database/models/mtproto_models.go`:
```go
type ConnectedAccount struct {
    ID               string    `gorm:"type:text;primaryKey" json:"id"`
    UserID           int64     `gorm:"uniqueIndex;not null" json:"userId"`
    TelegramUserID   int64     `gorm:"not null" json:"telegramUserId"`
    Username         string    `json:"username"`
    EncryptedSession string    `gorm:"type:text;not null" json:"-"`
    CreatedAt        time.Time `gorm:"autoCreateTime" json:"createdAt"`
    UpdatedAt        time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
    LastUsedAt       *time.Time `json:"lastUsedAt"`
    Enabled          bool      `gorm:"default:true" json:"enabled"`
}

type ConnectedAccountChannel struct {
    ID                 string           `gorm:"type:text;primaryKey" json:"id"`
    ConnectedAccountID string           `gorm:"uniqueIndex:idx_acc_channel;not null" json:"connectedAccountId"`
    ChannelID          int64            `gorm:"uniqueIndex:idx_acc_channel;not null" json:"channelId"`
    Enabled            bool             `gorm:"default:true" json:"enabled"`
    ConnectedAccount   *ConnectedAccount `gorm:"foreignKey:ConnectedAccountID;constraint:OnDelete:CASCADE;" json:"-"`
}
```

10. Adicionar AutoMigrate em `internal/database/database.go`

11. Criar `internal/database/repositories/connected_account.go` com CRUD básico

12. Criar `internal/core/services/connected_account.go` com lógica de negócio

### FASE 3: Criptografia e sessão

13. Criar `internal/telegram/mtproto/encryption/encryption.go`:
- `Encrypt(plaintext []byte) ([]byte, error)` usando AES-256-GCM
- `Decrypt(ciphertext []byte) ([]byte, error)` 
- Chave derivada do `ENCRYPTION_KEY` via HKDF ou SHA256

14. Criar `internal/telegram/mtproto/session/session.go`:
- `SessionStorage` interface (load/save/delete)
- `PostgresSessionStorage` implementando com repository + encryption

### FASE 4: Autenticação MTProto

15. Criar `internal/telegram/mtproto/auth/types.go`:
```go
type AuthState struct {
    UserID       int64
    PhoneNumber  string
    PhoneCodeHash string
    AuthKey      []byte  // temporary, in-memory only
    ExpiresAt    time.Time
}
```

16. Criar `internal/telegram/mtproto/auth/auth.go`:
- `SendCode(ctx, phoneNumber) (hash, error)` 
- `SignIn(ctx, code, hash) error`
- `Password(ctx, password) error`
- Usar Redis com TTL curto (5min) para guardar AuthState

17. Criar `internal/telegram/mtproto/client/client.go`:
- Wrapper around `telegram.Client` from gotd/td
- Connect/Disconnect com base em sessão descriptografada
- Rate limit handling (FloodWait)

### FASE 5: API REST

18. Criar `internal/api/types/account.go`:
```go
type AccountStatusResponse struct {
    Status         string  `json:"status"` // "connected", "disconnected"
    TelegramID     *int64  `json:"telegramId,omitempty"`
    Username       *string `json:"username,omitempty"`
    ConnectedAt    *string `json:"connectedAt,omitempty"`
    LastUsedAt     *string `json:"lastUsedAt,omitempty"`
}

type ConnectRequest struct {
    PhoneNumber string `json:"phoneNumber" binding:"required"`
}

type VerifyRequest struct {
    Code string `json:"code" binding:"required"`
}

type PasswordRequest struct {
    Password string `json:"password" binding:"required"`
}
```

19. Criar `internal/api/controllers/accountController.go`:
- `GetAccountStatus` (GET)
- `ConnectAccount` (POST) - inicia auth
- `VerifyCode` (POST) - verifica código
- `SendPassword` (POST) - envia 2FA
- `DisconnectAccount` (DELETE) - remove conta

20. Atualizar `internal/api/routes/routes.go`:
```go
accountRoutes := api.Group("/account")
accountRoutes.Use(auth.AuthMiddlewareJWT(c))
{
    accountRoutes.GET("", accountController.GetAccountStatus)
    accountRoutes.POST("/connect", accountController.ConnectAccount)
    accountRoutes.POST("/verify", accountController.VerifyCode)
    accountRoutes.POST("/password", accountController.SendPassword)
    accountRoutes.DELETE("", accountController.DisconnectAccount)
}
```

### FASE 6: Frontend

21. Atualizar `dashboard/src/types.ts` com tipos de conta

22. Atualizar `dashboard/src/api.ts` com chamadas de account

23. Criar `dashboard/src/components/AuthDisclaimer.tsx`:
- Checkboxes de Termos de Uso e Política de Privacidade
- Botão "Continuar"
- Texto explicativo sobre segurança

24. Criar `dashboard/src/components/AuthFlow.tsx`:
- Step 1: Input de telefone + validação
- Step 2: Input de código SMS
- Step 3: Input de senha 2FA (se necessário)
- Loading states, error handling

25. Criar `dashboard/src/components/ConnectedAccountCard.tsx`:
- Status badge (conectado/desconectado)
- Username, Telegram ID
- Data de conexão, última utilização
- Botões: Reconectar, Desconectar

26. Criar `dashboard/src/components/ContaTelegramTab.tsx`:
- Container principal que gerencia estados
- Integra AuthDisclaimer, AuthFlow, ConnectedAccountCard
- Lógica de polling de status

27. Atualizar `dashboard/src/App.tsx`:
- Adicionar tab "Conta Telegram" na navegação
- Renderizar ContaTelegramTab quando selecionado

### FASE 7: MTProto Executor

28. Criar `internal/telegram/mtproto/executor/executor.go`:
- `MTProtoExecutor` struct com pool de clientes
- `EditMessage` - usa `messages.EditMessage` do MTProto
- `EditCaption` - usa `messages.EditMessage` com `media` field
- `EditReplyMarkup` - usa `messages.EditMessage` só com reply_markup
- Rate limit handling (FloodWait → sleep + retry)
- Atualizar `last_used_at` na conta

### FASE 8: Integração no pipeline

29. Atualizar `internal/telegram/events/channelPost/dispatch_telego.go`:
- Substituir `pCtx.Bot.EditMessageText(...)` por `pCtx.Executor.EditMessage(...)`
- Substituir `pCtx.Bot.EditMessageCaption(...)` por `pCtx.Executor.EditCaption(...)`
- Substituir `pCtx.Bot.EditMessageReplyMarkup(...)` por `pCtx.Executor.EditReplyMarkup(...)`
- Substituir `pCtx.Bot.SendSticker(...)` por `pCtx.Executor.SendSticker(...)`

30. Criar `internal/telegram/executor/factory.go`:
- `NewExecutor(bot, accountService) TelegramExecutor`
- Se user tem conta conectada ativa → `MTProtoExecutor`
- Senão → `BotAPIExecutor`

31. Atualizar `internal/container/appContainer.go`:
- Adicionar `ConnectedAccountService`, `ConnectedAccountRepository`
- Adicionar `MTProtoAuthService`
- Adicionar `ExecutorFactory` (se necessário)

### FASE 9: Testes e finalização

32. Testes de criptografia (roundtrip, integridade, chave inválida)

33. Testes de repositório (CRUD, busca por user_id, cascade delete)

34. Testes de fluxo de autenticação (mock do client gotd)

35. Testes de permissões (conta ativa, canal autorizado, sessão válida)

36. Adicionar logging específico de MTProto no `pkg/logger/logger.go`

37. Atualizar `.env-example` com `ENCRYPTION_KEY`

38. Atualizar `docker-compose.yml` removendo `version` obsoleto

## Riscos

- **Alto**: `gotd/td` é uma biblioteca complexa com MTProto raw — curva de aprendizado alta e possíveis breaking changes
- **Alto**: Sessões MTProto são específicas por datacenter — pooling entre instâncias precisa de cuidado
- **Médio**: FloodWait do Telegram pode ser agressivo em MTProto (mais que Bot API) — precisa de rate limiting robusto
- **Médio**: Integração no pipeline existente pode introduzir regressões se a interface não capturar todos os casos de uso
- **Médio**: Sincronização de sessão entre múltiplas instâncias do servidor (escalabilidade horizontal)
- **Baixo**: Compatibilidade de tipos entre gotd e telego (chat IDs, message IDs, etc.)
- **Baixo**: Dependências adicionais (gotd puxa várias libs como `go.uber.org/zap`, `golang.org/x/net`, etc.)

## Impactos esperados

- **Performance**: Leve overhead inicial na criação da factory de executor (resolução de conta conectada via Redis/DB)
- **Segurança**: Sessões criptografadas no PostgreSQL, dados sensíveis apenas em memória
- **Manutenção**: Código mais limpo com interface `TelegramExecutor` separando concerns
- **Escalabilidade**: Pool de clientes MTProto e workers assíncronos
- **UX**: Usuários podem usar contas próprias para recursos Premium
- **Compatibilidade**: 100% backward compatible — Bot API continua como fallback

## Compatibilidade
- Linux ✓
- macOS ✓ (desenvolvimento)
- Windows ✓ (WSL2)
- Docker ✓
- CI/CD ✓ (testes automatizados)

## Como testar

### Build
```bash
go build ./cmd/FreddyBot/
```

### Testes unitários
```bash
go test ./internal/telegram/mtproto/encryption/...
go test ./internal/database/repositories/... -run ConnectedAccount
go test ./internal/core/services/... -run ConnectedAccount
```

### Testes de integração (docker)
```bash
sg docker -c "docker compose up -d"
go test ./internal/... -tags=integration
```

### Execução
```bash
go run ./cmd/FreddyBot/
```

### Verificação frontend
```bash
cd dashboard && npm run build
```

## Rollback

1. Reverter alterações no `cmd/FreddyBot/main.go` e `container/appContainer.go`
2. Remover as novas rotas de API
3. Reverter `dispatch_telego.go` para usar `pCtx.Bot` diretamente
4. Remover modelos do AutoMigrate
5. Remover dependência `gotd/td` do `go.mod`
6. Reverter frontend

## Observações

- A implementação deve ser **incremental** — cada fase é funcional por si só
- A Fase 1 (interface) já traz benefício imediato de organização de código, mesmo sem MTProto
- O `gotd/td` requer Go 1.20+ (projeto já usa 1.25.7)
- Sessões MTProto expiram — o worker deve detectar e reconectar automaticamente
- Para o MVP, focar em **edição de mensagens** (EditMessageText, EditMessageCaption, EditMessageReplyMarkup) — operações mais frequentes
- Recursos Premium (emoji, reações) virão após o MVP básico funcionando
- A chave `ENCRYPTION_KEY` deve ter no mínimo 32 bytes (64 caracteres hex ou string raw)
- `gotd/td` usa `context.Context` extensivamente — alinhado com o padrão do projeto
