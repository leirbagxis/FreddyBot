# Plano: refatorar-teste-stars-com-reembolso

## Pedido do usuário
Consertar bug na parte admin onde usuário consegue selecionar o mesmo canal várias vezes e comprar. Mudar política de testes: em vez de ativar grátis, usar valores de apenas 1 star com opção de reembolso na dashboard admin. Salvar o ID da compra (telegram_payment_charge_id) no banco, ter controle de reembolsos (tabela refunds), e realizar o reembolso corretamente via API do Telegram.

## Objetivo
1. Alterar `STARS_TEST_MODE` para definir preços como 1 star em vez de ativar grátis
2. Garantir que `TelegramPaymentID` seja sempre salvo na subscription
3. Criar modelo/tabela `Refund` para rastrear reembolsos
4. Implementar `AdminRefundPayment` no backend (chama `bot.RefundStarPayment()`)
5. Adicionar endpoint `POST /api/admin/subscriptions/refund`
6. Adicionar seção de reembolso no `AdminSubscriptionsTab` (frontend)
7. Validar no backend que canais não sejam duplicados / impedir compras duplicadas

## Contexto atual
- `STARS_TEST_MODE=true` + `?test=true` ativa assinatura **sem pagamento** (apenas exibe invoice)
- Preços reais: base=80 stars, extra=35 stars
- `TelegramPaymentID` já existe no model `Subscription` mas pode ficar vazio em modo teste
- Não existe tabela de reembolsos — cancelamento admin apenas expira assinatura
- AdminSubscriptionsTab já existe com listagem e cancelamento
- Telegram Bot API tem método `RefundStarPayment(userID, chargeID)`
- Frontend `PremiumTab` já deduplica canais no select, mas backend não valida

## Arquivos analisados
- `internal/core/services/subscription_service.go`
- `internal/database/models/subscription_models.go`
- `internal/database/repositories/subscription_repository.go`
- `internal/api/controllers/subscription_controller.go`
- `internal/api/controllers/adminController/adminSubscriptionController.go`
- `internal/database/database.go`
- `pkg/config/config.go`
- `dashboard/src/components/PremiumTab.tsx`
- `dashboard/src/components/AdminSubscriptionsTab.tsx`
- `dashboard/src/api.ts`
- `dashboard/src/types.ts`

## Arquivos que poderão ser modificados

### Backend
- `pkg/config/config.go` — alterar comentário do StarsTestMode
- `internal/database/models/subscription_models.go` — adicionar modelo Refund
- `internal/database/database.go` — adicionar Refund ao AutoMigrate
- `internal/database/repositories/subscription_repository.go` — adicionar FindByChargeID, métodos Refund
- `internal/core/services/subscription_service.go` — alterar CreateInvoice (1 star em test mode), adicionar AdminRefundPayment
- `internal/api/controllers/adminController/adminSubscriptionController.go` — adicionar endpoint Refund
- `internal/api/routes/routes.go` — registrar nova rota (se necessário)

### Frontend
- `dashboard/src/api.ts` — adicionar adminRefundPayment
- `dashboard/src/types.ts` — adicionar interface Refund (opcional)
- `dashboard/src/components/AdminSubscriptionsTab.tsx` — adicionar seção de reembolso por subscription

## Estratégia de implementação

### 1. Alterar StarsTestMode (free → 1 star)
- No `CreateInvoice`: quando `StarsTestMode` estiver ativo, sobrescrever `totalStars = 1` para qualquer channelCount
- Remover o bloco `if testMode && config.StarsTestMode { activateSubscription }` — não ativar sem pagamento
- O usuário SEMPRE paga (1 star em test mode, valor real em prod)
- Mudar comentário no config

### 2. Novo modelo Refund
```go
type Refund struct {
    ID                      string    `gorm:"type:text;primaryKey" json:"id"`
    SubscriptionID          string    `gorm:"index;not null" json:"subscriptionId"`
    UserID                  int64     `gorm:"index;not null" json:"userId"`
    TelegramPaymentChargeID string    `gorm:"type:text;not null" json:"telegramPaymentChargeId"`
    AmountStars             int       `gorm:"not null" json:"amountStars"`
    Status                  string    `gorm:"type:text;default:processed" json:"status"`
    RefundedAt              time.Time `json:"refundedAt"`
    RefundedBy              int64     `json:"refundedBy"`
    CreatedAt               time.Time `gorm:"autoCreateTime" json:"createdAt"`
}
```
Adicionar ao AutoMigrate.

### 3. AdminRefundPayment
- Receber `userID` e `telegramPaymentChargeID`
- Verificar se subscription existe e está ativa
- Verificar se já não foi reembolsado (buscar Refund por chargeID)
- Chamar `bot.RefundStarPayment(ctx, &RefundStarPaymentParams{UserID: userID, TelegramPaymentChargeID: chargeID})`
- Se ok: expirar subscription (instant), salvar Refund record
- Retornar erro se chargeID não existir ou refund já processado

### 4. Endpoint REST
`POST /api/admin/subscriptions/refund`
```json
{"userId": 12345, "telegramPaymentChargeId": "charge_abc123"}
```

### 5. Frontend AdminSubscriptionsTab
- Na listagem, exibir `telegramPaymentId` quando disponível
- Botão "Reembolsar" ao lado de cada subscription ativa com charge_id
- Modal de confirmação
- Indicador visual de já reembolsado

### 6. Validação de canais duplicados (backend)
- Em `CreateInvoice`, garantir que `channelCount` é >= 1 e <= número máximo de canais do usuário
- Validar que usuário não possa ter subscription duplicada (já existe)

## Passos detalhados

### Passo 1: Modelo Refund
- Adicionar struct `Refund` em `subscription_models.go`
- Adicionar ao `AutoMigrate` em `database.go`

### Passo 2: Config + Service — alterar StarsTestMode
- `config.go`: atualizar comentário (agora = preço 1 star, não free)
- `subscription_service.go`:
  - Em `CreateInvoice`: se `StarsTestMode && testMode`, forçar `totalStars = 1`
  - Remover o bloco de ativação automática sem pagamento (linhas 166-171)
  - Sempre retornar invoice URL real

### Passo 3: Repositório — métodos de refund
- `subscription_repository.go`:
  - `FindByChargeID(ctx, chargeID string) (*models.Subscription, error)`
  - `FindRefundByChargeID(ctx, chargeID string) (*models.Refund, error)`
  - `CreateRefund(ctx, refund *models.Refund) error`
  - `FindRefundsByUserID(ctx, userID int64) ([]models.Refund, error)`

### Passo 4: Service — AdminRefundPayment
- Implementar em `subscription_service.go`
- Lógica:
  1. Buscar subscription por chargeID
  2. Verificar se já reembolsado
  3. Chamar `bot.RefundStarPayment()`
  4. Se sucesso: expirar subscription, criar Refund record
  5. Sincronizar features

### Passo 5: Controller — endpoint refund
- `adminSubscriptionController.go`: adicionar handler `Refund`
- Registrar rota `POST /api/admin/subscriptions/refund` (se necessário)

### Passo 6: Frontend API
- `api.ts`: adicionar `adminRefundPayment(userId, chargeId)`

### Passo 7: Frontend AdminSubscriptionsTab
- Expandir cada subscription item:
  - Mostrar `telegramPaymentId` em tooltip/segunda linha
  - Botão "Reembolsar" para subscriptions ativas com charge_id
  - Estado visual de reembolsado
- Adicionar modal de confirmação de reembolso

### Passo 8: Build + testes
- Rodar `go build ./...`
- Rodar `npm run build`
- Validar que não há erros

## Riscos
- `RefundStarPayment` pode falhar se chargeID não existir ou já reembolsado no Telegram — tratar erro
- Duas chamadas concorrentes para refund do mesmo chargeID — usar verificação com `FindRefundByChargeID`
- Atualização do banco (AutoMigrate adiciona tabela) segura com GORM
- Frontend precisa lidar com subscriptions que não têm `telegramPaymentId` (subscriptions antigas)

## Impactos esperados
- Test mode agora custa 1 star (não mais grátis) — usuário vê o fluxo real de pagamento
- Admin pode reembolsar qualquer pagamento via dashboard
- Todas as subscriptions terão charge_id registrado (mesmo em test mode)
- Tabela refunds mantém histórico auditável
- Backend valida channelCount máximo

## Compatibilidade
- Linux ✓
- macOS ✓
- Windows ✓ (via WSL)
- Docker ✓
- CI/CD ✓

## Como testar

### Build
```bash
go build ./...
cd dashboard && npm run build
```

### Testes manuais
1. Set `STARS_TEST_MODE=true` no .env
2. Abrir dashboard, clicar em Premium, selecionar canais, assinar
3. Verificar que invoice é de 1 star (não grátis)
4. Pagar (ou simular)
5. Ir em Admin > Assinaturas, ver charge_id
6. Clicar em "Reembolsar"
7. Verificar que subscription expirou e features foram limpas
8. Verificar que duplicatas não são permitidas (tentar reembolsar de novo)

### Rollback
- Remover tabela `refunds` do banco (ou ignorar)
- Reverter `StarsTestMode` para o comportamento antigo
- `git revert` nos commits

## Observações
- A biblioteca telego já suporta `RefundStarPayment` (v1.9.0)
- O método `RefundStarPayment` do Telegram Bot API requer que o bot seja o receiver dos Stars
- O reembolso no Telegram é instantâneo e não reversível
