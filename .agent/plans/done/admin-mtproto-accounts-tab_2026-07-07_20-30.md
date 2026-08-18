# Plano: admin-mtproto-accounts-tab

## Pedido do usuário
Criar uma nova aba no dashboard admin onde o administrador possa conectar e gerenciar múltiplas contas MTProto para serem utilizadas no fluxo de edição de postagem.

## Objetivo
Adicionar uma aba "Contas MTProto" no painel administrativo (`/admin/dash`) que permita:
1. Listar todas as contas MTProto gerenciadas pelo admin
2. Conectar novas contas Telegram via fluxo de autenticação (phone → code → 2FA)
3. Desconectar/remover contas existentes
4. Visualizar status (conectado/desconectado), username, Telegram ID, último uso
5. Alternar enable/disable de cada conta

## Contexto atual
- O sistema já possui `ConnectedAccount` (usuário conecta sua própria conta, 1 por usuário)
- O `MTProtoAuthService` gerencia o fluxo de autenticação (SendCode, VerifyCode, VerifyPassword)
- O `ExecutorFactory.ForUser()` checa `HasActiveAccount(userID)` para decidir se usa MTProto
- Contas de usuário são criptografadas com AES-256-GCM e salvas em `connected_accounts`
- Redis guarda estado temporário de autenticação (`mtproto_auth:<userID>`) com TTL de 5 min
- Admin dashboard atual tem 7 abas: Overview, Usuários, Canais, Auditoria, Logs, Broadcast, Config

## Arquivos analisados
- `internal/database/models/mtproto_models.go` — modelos ConnectedAccount e ConnectedAccountChannel
- `internal/database/repositories/connected_account.go` — repositório de contas conectadas
- `internal/core/services/connected_account.go` — serviço de contas conectadas
- `internal/api/controllers/accountController.go` — endpoints /api/account/*
- `internal/telegram/mtproto/auth/auth.go` — serviço de autenticação MTProto
- `internal/api/routes/routes.go` — registro de rotas
- `internal/container/appContainer.go` — DI container
- `dashboard/src/components/AdminDashboard.tsx` — dashboard admin
- `dashboard/src/App.tsx` — tabs do admin
- `dashboard/src/api.ts` — chamadas API
- `dashboard/src/types.ts` — tipos TypeScript

## Arquivos que poderão ser modificados

### Backend (Go)
- NOVO: `internal/database/models/admin_mtproto_account.go` — modelo AdminMTProtoAccount
- NOVO: `internal/database/repositories/admin_account.go` — repositório admin accounts
- NOVO: `internal/core/services/admin_account.go` — serviço admin accounts + AccountSaver
- NOVO: `internal/api/controllers/adminController/adminAccountController.go` — endpoints admin
- MODIFICADO: `internal/api/routes/routes.go` — registrar novas rotas admin
- MODIFICADO: `internal/container/appContainer.go` — registrar novos serviços/controllers
- MODIFICADO: `internal/database/database.go` — AutoMigrate do novo modelo

### Frontend (React/TypeScript)
- NOVO: `dashboard/src/components/AdminMTProtoAccountsTab.tsx` — nova aba
- MODIFICADO: `dashboard/src/App.tsx` — adicionar tab aos adminTabs
- MODIFICADO: `dashboard/src/components/AdminDashboard.tsx` — renderizar nova aba
- MODIFICADO: `dashboard/src/api.ts` — funções de API para admin accounts
- MODIFICADO: `dashboard/src/types.ts` — tipos AdminMTProtoAccount

## Estratégia de implementação

### Modelo de dados separado
Criar modelo `AdminMTProtoAccount` separado do `ConnectedAccount` para:
- Não interferir com o fluxo existente de usuários conectarem suas próprias contas
- Admin pode gerenciar múltiplas contas independentemente de qualquer userID
- Cada conta admin tem um label amigável (ex: "Conta do João - @joaobot")
- Status explícito (connected/disconnected/error)

### Reaproveitamento do Auth Service
O `MTProtoAuthService` será reaproveitado, mas com um `AccountSaver` diferente que salva no modelo admin. Para o Redis, usaremos uma chave prefixada `mtproto_auth:admin:<admin_user_id>:<index>`.

### Integração com ExecutorFactory (futuro)
Nesta primeira iteração, as contas serão gerenciadas visualmente. A integração com o `ExecutorFactory` (para usar as contas admin como fallback quando o usuário não tem conta própria) será feita em uma etapa posterior.

## Passos detalhados

### 1. Modelo AdminMTProtoAccount (Go)
Criar `internal/database/models/admin_mtproto_account.go`:
```go
type AdminMTProtoAccount struct {
    ID               string     `gorm:"type:text;primaryKey" json:"id"`
    Label            string     `json:"label"`
    PhoneNumber      string     `json:"phoneNumber"`
    TelegramUserID   int64      `json:"telegramUserId"`
    Username         string     `json:"username"`
    FirstName        string     `json:"firstName"`
    EncryptedSession string     `gorm:"type:text;not null" json:"-"`
    Enabled          bool       `gorm:"default:true" json:"enabled"`
    Status           string     `gorm:"default:disconnected" json:"status"`
    LastUsedAt       *time.Time `json:"lastUsedAt"`
    CreatedAt        time.Time  `json:"createdAt"`
    UpdatedAt        time.Time  `json:"updatedAt"`
}
```

### 2. Repositório AdminAccountRepository
Criar `internal/database/repositories/admin_account.go`:
- List(ctx) → []AdminMTProtoAccount
- GetByID(ctx, id) → *AdminMTProtoAccount
- Create(ctx, account)
- Update(ctx, account)
- Delete(ctx, id)
- UpdateStatus(ctx, id, status)

### 3. Serviço AdminAccountService
Criar `internal/core/services/admin_account.go`:
- `SaveSession()` — implementa AccountSaver para admin accounts
- `ListAccounts(ctx)` — lista todas as contas
- `GetAccount(ctx, id)` — busca conta por ID
- `DeleteAccount(ctx, id)` — remove conta + chaves Redis
- `GetSessionAndID(ctx, accountID)` — descriptografa sessão para uso MTProto

### 4. Controller AdminAccountController
Criar `internal/api/controllers/adminController/adminAccountController.go`:
- `ListAccounts` → GET /api/admin/accounts
- `ConnectAccount` → POST /api/admin/accounts/connect (inicia auth)
- `VerifyCode` → POST /api/admin/accounts/verify
- `SendPassword` → POST /api/admin/accounts/password
- `DeleteAccount` → DELETE /api/admin/accounts/:id
- `ToggleAccount` → POST /api/admin/accounts/:id/toggle (enable/disable)

### 5. Rotas
Adicionar em `routes.go` dentro do grupo admin:
```go
adminRoute.GET("/accounts", adminAccountController.ListAccounts)
adminRoute.POST("/accounts/connect", adminAccountController.ConnectAccount)
adminRoute.POST("/accounts/verify", adminAccountController.VerifyCode)
adminRoute.POST("/accounts/password", adminAccountController.SendPassword)
adminRoute.DELETE("/accounts/:id", adminAccountController.DeleteAccount)
adminRoute.POST("/accounts/:id/toggle", adminAccountController.ToggleAccount)
```

### 6. Container
Registrar AdminAccountRepository, AdminAccountService, AdminAccountController no container.

### 7. AutoMigrate
Adicionar `AdminMTProtoAccount{}` ao AutoMigrate em `database.go`.

### 8. Componente AdminMTProtoAccountsTab (React)
Criar `dashboard/src/components/AdminMTProtoAccountsTab.tsx`:
- **Lista de contas**: cards mostrando username, label, status (conectado/desconectado), Telegram ID, último uso
- **Botão "Conectar Nova Conta"**: abre modal de autenticação com 3 passos:
  1. Inserir número de telefone + label
  2. Inserir código SMS
  3. Inserir senha 2FA (se necessário)
- **Ações por conta**: toggle enable/disable, botão desconectar
- Design consistente com AdminConfigTab (Cloudflare-inspired)

### 9. Integração no Dashboard
- Adicionar tab 'accounts' ao `adminTabs` em `App.tsx`
- Adicionar renderização condicional em `AdminDashboard.tsx`
- Adicionar funções API em `api.ts`
- Adicionar tipos em `types.ts`

## Riscos
- Fluxo de autenticação MTProto depende de `MTPROTO_APP_ID` e `MTPROTO_APP_HASH` estarem configurados
- Redis necessário para estado temporário de autenticação
- Sessões criptografadas com AES-256-GCM — perda da SECRET_KEY torna sessões irrecuperáveis
- Stub mode (sem credenciais MTProto) não permite testar auth real em dev
- Conflito de chaves Redis se admin e user estiverem autenticando simultaneamente (usar prefixo diferente)

## Impactos esperados
- Nenhuma mudança no fluxo existente de usuários conectarem suas contas
- Admin ganha visibilidade centralizada de contas MTProto
- Base para futura integração com ExecutorFactory (fallback para contas admin)
- +1 tabela no banco de dados

## Compatibilidade
- Linux ✓
- macOS ✓
- Windows ✓
- Docker ✓
- CI/CD ✓

## Como testar

### Backend
```bash
cd /home/malbs/Opencode/FreddyBot
go build ./cmd/FreddyBot/
```

### Frontend
```bash
cd /home/malbs/Opencode/FreddyBot/dashboard
npm run build
```

## Rollback
1. Reverter alterações no routes.go e container
2. Remover migração da tabela admin_mtproto_accounts
3. Reverter commits

## Observações
- O componente AuthFlow.tsx existente (`ConnectedAccountCard` + `AuthFlow`) serve como referência de UI para o fluxo de autenticação
- A criptografia de sessão usa `encryption.Encrypt/Decrypt` — mesma do ConnectedAccountService
- Admin pode conectar quantas contas quiser (sem limite de 1 por user)
