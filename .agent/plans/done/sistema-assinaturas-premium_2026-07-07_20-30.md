# Plano: sistema-assinaturas-premium

## Pedido do usuário
Implementar um sistema completo de assinatura premium para o LegendasBr BOT, com **Telegram Stars** como único provedor de pagamento.

A assinatura dá direito a usar a Managed Premium Account (conta MTProto admin gerenciada) como executor, substituindo a conta pessoal do usuário.

## Objetivo
Criar um sistema de assinaturas mensais com:
- Modelos de dados (`Subscription`, `SubscriptionItem`)
- Integração com Telegram Stars (Bot API `sendInvoice`)
- Serviço de assinatura (criação, renovação, cancelamento)
- Dashboard do usuário para gerenciar assinatura
- Integração com `ExecutorFactory` para usar Managed Premium Account
- Sistema de entitlements/features (substitui `if premium {}`)

## Contexto atual

### Backend (Go)
- Clean Architecture com models → repositories → services → controllers
- GORM + PostgreSQL, Redis para cache/sessões
- `TelegramExecutor` interface com 3 implementações: `BotAPIExecutor`, `MTProtoExecutor`, `UserExecutor`
- `ExecutorFactory.ForUser()` atualmente decide entre BotAPI e MTProto (conta conectada do usuário)
- `User` model não tem campo de features/entitlements
- `AdminMTProtoAccount` model já existe (admin gerencia contas MTProto)
- `ConnectedAccount` é a conta pessoal do usuário (1 por user)
- `telego` bot client já disponível no container

### Frontend (React/Next.js)
- Dashboard com tabs `Início | Legendas | Botões | Permissões | Conta Telegram`
- Admin dashboard com tabs incluindo `Contas MTProto`
- UI estilo Cloudflare/TON com cards e toggles

## Arquivos analisados
- `internal/database/models/models.go` — User, Channel, etc.
- `internal/database/models/mtproto_models.go` — ConnectedAccount
- `internal/database/models/admin_mtproto_account.go` — AdminMTProtoAccount
- `internal/core/services/admin_account.go` — AdminAccountService (padrão de serviço)
- `internal/database/repositories/admin_account.go` — AdminAccountRepository (padrão de repo)
- `internal/telegram/executor/executor.go` — TelegramExecutor interface
- `internal/telegram/executor/factory.go` — ExecutorFactory
- `internal/telegram/executor/user_executor.go` — UserExecutor
- `internal/telegram/executor/mtproto.go` — MTProtoExecutor
- `internal/container/appContainer.go` — DI Container
- `internal/api/routes/routes.go` — Rotas da API
- `dashboard/src/App.tsx` — Dashboard principal
- `dashboard/src/components/AdminDashboard.tsx` — Admin dashboard
- `dashboard/src/types.ts` — Tipos TypeScript
- `dashboard/src/api.ts` — Funções API

## Arquivos que poderão ser modificados

### Backend (Go)
- `internal/database/models/models.go` — Add Features field to User
- `internal/database/models/subscription_models.go` — NOVO: Subscription, SubscriptionItem
- `internal/database/repositories/subscription_repository.go` — NOVO
- `internal/core/services/subscription_service.go` — NOVO: lógica de assinatura + Stars
- `internal/api/controllers/subscription_controller.go` — NOVO
- `internal/api/routes/routes.go` — Adicionar rotas de subscription
- `internal/container/appContainer.go` — Registrar novos serviços
- `internal/telegram/executor/factory.go` — Adicionar terceiro caminho (Premium Account)
- `internal/telegram/executor/premium_executor.go` — NOVO: executor para premium

### Frontend (React)
- `dashboard/src/App.tsx` — Adicionar tab de assinatura
- `dashboard/src/api.ts` — Adicionar funções de subscription
- `dashboard/src/types.ts` — Adicionar tipos de subscription
- `dashboard/src/components/PremiumTab.tsx` — NOVO: componente premium dashboard

## Estratégia de implementação

### Fase 1 — Modelos e Repositórios
1. Criar `subscription_models.go` com `Subscription` e `SubscriptionItem`
2. Adicionar campo `Features` (JSON) no model `User` com struct `UserFeatures`
3. Criar `subscription_repository.go`

### Fase 2 — Serviço de Assinatura (Telegram Stars)
4. Criar `subscription_service.go` com:
   - `GetSubscription(ctx, userID)` — status atual
   - `CreateSubscription(ctx, userID)` — inicia fluxo Stars (gera invoice)
   - `HandleStarsPreCheckout(ctx, payload)` — valida `pre_checkout_query`
   - `HandleStarsPayment(ctx, payload)` — ativa subscription via `successful_payment`
   - `CancelSubscription(ctx, userID)` — cancela no fim do período
   - `AddExtraChannel(ctx, userID)` / `RemoveExtraChannel(ctx, userID)`
   - `SyncFeatures(ctx, userID)` — subscription → User.Features
   - `RenewSubscription(ctx, sub)` — renovação mensal automática

### Fase 3 — Subscription Controller + Rotas
5. Criar `subscription_controller.go` com:
   - `GET /api/subscription` — status
   - `POST /api/subscription/create` — criar assinatura (retorna invoice link)
   - `POST /api/subscription/pre-checkout` — webhook interno para pre_checkout_query
   - `POST /api/subscription/successful-payment` — webhook interno para successful_payment
   - `POST /api/subscription/cancel` — cancelar
   - `POST /api/subscription/channels/add` — adicionar canal extra
   - `POST /api/subscription/channels/remove` — remover canal extra
6. Registrar rotas em `routes.go`

### Fase 4 — Integração com ExecutorFactory
7. Modificar `ExecutorFactory.ForUser()` para:
   - 1º: BotAPI (sempre disponível)
   - 2º: ConnectedAccount (se tiver conta pessoal)
   - 3º: Premium Account (se tiver `features.ManagedPremiumAccount`)
8. Criar `PremiumExecutor` que usa `AdminAccountService` para obter sessão da conta admin

### Fase 5 — Frontend
9. Criar `PremiumTab.tsx` com:
   - Card de status da assinatura (ativa/inativa)
   - Preço: ⭐80/mês (base) + ⭐35/mês (canal extra)
   - Botão "Assinar com Telegram Stars"
   - Gerenciamento de canais extras
   - Botão cancelar
10. Adicionar tab "Premium" no array `tabs` em `App.tsx`
11. Adicionar funções API e tipos TypeScript

### Fase 6 — Container
12. Registrar no `appContainer.go`:
    - `SubscriptionRepository`
    - `SubscriptionService`
    - `SubscriptionController`
    - `PremiumExecutor`

## Passos detalhados

### Passo 1: Models de Subscription

`internal/database/models/subscription_models.go`:

```go
package models

import "time"

type SubscriptionStatus string
const (
    SubscriptionActive   SubscriptionStatus = "active"
    SubscriptionCanceled SubscriptionStatus = "canceled"
    SubscriptionExpired  SubscriptionStatus = "expired"
)

type Subscription struct {
    ID                string             `gorm:"type:text;primaryKey" json:"id"`
    UserID            int64              `gorm:"uniqueIndex;not null" json:"userId"`
    Status            SubscriptionStatus `gorm:"type:text;default:active" json:"status"`
    CurrentPeriodStart time.Time         `json:"currentPeriodStart"`
    CurrentPeriodEnd  time.Time          `json:"currentPeriodEnd"`
    ExtraChannels     int                `gorm:"default:0" json:"extraChannels"`
    CancelAtPeriodEnd bool               `gorm:"default:false" json:"cancelAtPeriodEnd"`
    TelegramPaymentID string             `gorm:"type:text" json:"telegramPaymentId"` // charge_id do Stars
    CreatedAt         time.Time          `gorm:"autoCreateTime" json:"createdAt"`
    UpdatedAt         time.Time          `gorm:"autoUpdateTime" json:"updatedAt"`

    User *User `gorm:"foreignKey:UserID" json:"-"`
}

type UserFeatures struct {
    ManagedPremiumAccount bool `json:"managedPremiumAccount,omitempty"`
    CustomEmojis          bool `json:"customEmojis,omitempty"`
    ExtraChannels         int  `json:"extraChannels,omitempty"`
}
```

### Passo 2: Features no User

Adicionar campo no struct `User` em `models.go`:
```go
Features string `gorm:"type:text;default:'{}'" json:"features"` // JSON de UserFeatures
```

### Passo 3: Subscription Repository

`internal/database/repositories/subscription_repository.go`:
- `FindByUserID(ctx, userID) (*Subscription, error)`
- `Create(ctx, sub *Subscription) error`
- `Update(ctx, sub *Subscription) error`
- `Delete(ctx, id string) error`
- `FindExpired(ctx) ([]Subscription, error)` — cleanup job
- `FindDueForRenewal(ctx) ([]Subscription, error)` — renewals

### Passo 4: Subscription Service

`internal/core/services/subscription_service.go`:

```go
type SubscriptionService struct {
    repo          *repositories.SubscriptionRepository
    userRepo      *repositories.UserRepository
    bot           *telego.Bot               // para enviar invoices
    executorCache *executor.ExecutorFactory  // invalidar cache
    adminSvc      *AdminAccountService       // para PremiumExecutor
}
```

Métodos principais:
- `CreateInvoice(ctx, userID int64) (string, error)` — cria invoice no Telegram
- `HandlePreCheckout(ctx, queryID string, payload string) error` — aprova/rejeita
- `HandlePayment(ctx, userID int64, chargeID string) error` — ativa subscription
- `Cancel(ctx, userID int64) error`
- `AddExtraChannel(ctx, userID int64) error` — +1 extra channel (cobra diff no próximo período)
- `RemoveExtraChannel(ctx, userID int64) error`
- `GetStatus(ctx, userID int64) (*Subscription, error)`
- `RenewExpiring(ctx) error` — job periódico para renovar via Stars (re-cobrança automática)
- `SyncFeatures(ctx, userID int64) error` — atualiza User.Features baseado na subscription

### Passo 5: Fluxo Telegram Stars

1. Bot envia `sendInvoice` com `currency="XTR"`, `prices=[{label:"LegendasBr Premium",amount:80}]`
2. Usuário vê a fatura no Telegram, confirma pagamento com Stars
3. Bot recebe update `pre_checkout_query` → `answerPreCheckoutQuery(ok=true)`
4. Bot recebe update `successful_payment` → ativa subscription + sincroniza features
5. Todo mês, job verifica subscriptions ativas e tenta renovar

### Passo 6: Premium Executor

`internal/telegram/executor/premium_executor.go`:

```go
// PremiumExecutor usa uma conta MTProto admin no lugar da conta do usuario.
// Diferente do UserExecutor que usa a conta conectada do proprio usuario,
// o PremiumExecutor obtem uma sessao de admin account gerenciada pelo bot.
type PremiumExecutor struct {
    adminSvc *services.AdminAccountService
    botAPI   TelegramExecutor
}

func (e *PremiumExecutor) EditMessage(ctx, chatID, msgID, text, parseMode, keyboard, opts) error {
    // Obtem a primeira AdminMTProtoAccount ativa
    // Usa MTProtoExecutor com a sessao da conta admin
    // Fallback para BotAPI se falhar
}
```

### Passo 7: Controller

`internal/api/controllers/subscription_controller.go`:
- `GetSubscription` — GET /api/subscription
- `CreateCheckout` — POST /api/subscription/create (retorna invoice link)
- `Cancel` — POST /api/subscription/cancel
- `AddExtraChannel` — POST /api/subscription/channels/add
- `RemoveExtraChannel` — POST /api/subscription/channels/remove

Os webhooks do Telegram (`pre_checkout_query`, `successful_payment`) são tratados via handlers de update do bot, não via API REST.

### Passo 8: Rotas

Em `routes.go`, dentro do grupo protegido:
```go
api.GET("/subscription", subscriptionController.GetSubscription)
api.POST("/subscription/create", subscriptionController.CreateInvoice)
api.POST("/subscription/cancel", subscriptionController.Cancel)
api.POST("/subscription/channels/add", subscriptionController.AddExtraChannel)
api.POST("/subscription/channels/remove", subscriptionController.RemoveExtraChannel)
```

### Passo 9: Telegram Update Handlers

No handler de updates do bot, adicionar:
- `pre_checkout_query` → `subscriptionService.HandlePreCheckout()`
- `successful_payment` → `subscriptionService.HandlePayment()`

### Passo 10: Frontend

`dashboard/src/components/PremiumTab.tsx`:
- Header com "Premium" e badge de status
- Se não assinante: card com preços e botão "Assinar com ⭐ Stars"
- Se assinante: card verde com "Assinatura Ativa", canais extras, botão cancelar
- Se cancelado/no fim: status "Cancelada — expira em DD/MM"

Adicionar em `App.tsx`:
- Nova tab `{ id: 'premium', label: 'Premium', icon: <Zap size={22} /> }`
- Renderizar `<PremiumTab />` quando ativa

### Passo 11: Container

Em `appContainer.go`:
```go
subscriptionRepo := repositories.NewSubscriptionRepository(db)
subscriptionService := services.NewSubscriptionService(subscriptionRepo, userRepo, telegoClient, executorFactory, adminAccountService)
subscriptionController := controllers.NewSubscriptionController(subscriptionService)
```

## Riscos
- **Telegram Stars em dev**: Precisa de testes com stars reais. Incluir modo "stub" que simula `successful_payment` sem cobrar.
- **Job de renovação**: Stars não suporta cobrança recorrente automática. A renovação precisa ser manual (o bot re-envia invoice ou o usuário paga novamente). Implementar lembrete + invoice reenvio.
- **AdminMTProtoAccount offline**: Se a conta admin ficar indisponível, o PremiumExecutor precisa fazer fallback para BotAPI.
- **Concorrência**: Webhooks podem vir duplicados. Usar idempotência via TelegramPaymentID.
- **Executor 3 vias**: Adicionar terceiro caminho no factory requer cuidado para não quebrar os dois existentes.

## Impactos esperados
- **+2 tabelas**: `subscriptions`, coluna `features` em `users`
- **+5 rotas protegidas**: `/api/subscription*`
- **+1 tab no dashboard do usuário**: "Premium"
- **+1 handler de update do bot**: `pre_checkout_query` + `successful_payment`
- **ExecutorFactory**: terceiro caminho de executor (Premium)
- **Nenhuma funcionalidade existente é alterada**: Bot API e Connected Account continuam funcionando

## Compatibilidade
- Linux ✅
- macOS ✅
- Windows ✅ (via WSL)
- Docker ✅
- CI/CD ✅

## Como testar

### Build
```bash
go build ./cmd/FreddyBot/
npm run build --prefix dashboard
```

### Testes
```bash
go test ./internal/core/services/... -v
go test ./internal/database/repositories/... -v
```

### Execução
```bash
go run ./cmd/FreddyBot/
```

## Rollback
1. Reverter mudanças no `ExecutorFactory` para usar apenas 2 caminhos
2. Remover rotas de subscription
3. Reverter models (remover coluna `features` de User)
4. Remover tabela `subscriptions`
5. Reverter container registrations

## Observações
- Telegram Stars não tem cobrança recorrente automática — renovação é manual (re-envio de invoice)
- `sendInvoice` precisa de `provider_token=""` para Stars (empty string, pagamento nativo Telegram)
- Fluxo: `sendInvoice` → `pre_checkout_query` → `answerPreCheckoutQuery` → `successful_payment`
- Preços em Stars (inteiros): 80 ⭐/mês base + 35 ⭐/mês por canal extra
- A conta `AdminMTProtoAccount` usada como Premium Account é configurada pelo admin no painel
- Se o admin não tiver nenhuma conta ativa, usuários premium usam BotAPI (fallback)
