# Plano: Cobrar Canal Extra Premium (Corrigir Bug de Pagamento)

## Pedido do usuário
Corrigir bug onde canais extras são ativados gratuitamente após assinatura. O usuário clica em canal inativo e ele ativa sem cobrar nada.

## Objetivo
Criar fluxo de pagamento para canais extras: ao clicar em canal inativo, criar invoice de35 Stars, cobrar via Telegram, e só após pagamento confirmado ativar o canal.

## Contexto atual
- Assinatura inicial: `CreateInvoice` cobra `basePrice + extras × extraPrice`
- `AddExtraChannel`: simplesmente faz `sub.ExtraChannels++` **sem cobrar**
- `HandlePayment`: não diferencia tipos de pagamento (só ativa assinatura)
- Payload: `"premium_sub:<userID>"` — sem distinção de tipo

## Arquivos analisados
- `internal/core/services/subscription_service.go` — CreateInvoice, HandlePayment, AddExtraChannel
- `internal/api/controllers/subscription_controller.go` — CreateInvoice, AddExtraChannel
- `internal/api/routes/routes.go` — rotas de assinatura
- `dashboard/src/api.ts` — funções frontend
- `dashboard/src/components/PremiumTab.tsx` — UI de gerenciamento

## Arquivos que serão modificados

### Backend
- `internal/core/services/subscription_service.go` — novo prefixo, novo método CreateExtraChannelInvoice, atualizar HandlePayment
- `internal/api/controllers/subscription_controller.go` — novo controller CreateExtraChannelInvoice
- `internal/api/routes/routes.go` — nova rota

### Frontend
- `dashboard/src/api.ts` — nova função createExtraChannelInvoice
- `dashboard/src/components/PremiumTab.tsx` — atualizar handleAddExtra

## Estratégia de implementação

### 1. Novo prefixo de payload
Adicionar constante `InvoiceExtraPayloadPrefix = "premium_extra:"` para diferenciar de assinatura.

### 2. Backend: CreateExtraChannelInvoice
Novo método no `SubscriptionService`:
- Verificar assinatura ativa
- Buscar preço do extra channel via `featureSvc.GetExtraChannelPrice()`
- Criar invoice com payload `"premium_extra:<userID>"`
- Retornar invoiceUrl

### 3. Backend: Atualizar HandlePayment
Adicionar lógica para detectar payload `"premium_extra:"`:
- Se payload começa com `"premium_extra:"` → chamar `AddExtraChannel` (incrementar contador)
- Se payload começa com `"premium_sub:"` → fluxo atual (ativar assinatura)

### 4. Backend: Atualizar HandlePreCheckout
Validar novo payload `"premium_extra:<userID>"` da mesma forma que `"premium_sub:<userID>"`.

### 5. Frontend: Nova API
Adicionar `createExtraChannelInvoice(testMode: boolean)` que chama `POST /api/subscription/channels/add-invoice?test=true`

### 6. Frontend: Atualizar PremiumTab
Substituir `handleAddExtra`:
- Em vez de chamar `addExtraChannel()` diretamente
- Chamar `createExtraChannelInvoice()` → obter invoiceUrl
- Abrir `openInvoice(invoiceUrl)` → callback confirma pagamento
- Após 'paid', chamar `loadStatus()` para atualizar estado

## Passos detalhados

1. Adicionar `InvoiceExtraPayloadPrefix` em subscription_service.go
2. Criar `CreateExtraChannelInvoice()` no SubscriptionService
3. Criar `CreateExtraChannelInvoice()` no SubscriptionController
4. Adicionar rota `POST /api/subscription/channels/add-invoice`
5. Atualizar `HandlePayment` para detectar payload e rotear
6. Atualizar `HandlePreCheckout` para validar novo payload
7. Adicionar `createExtraChannelInvoice()` no api.ts do frontend
8. Atualizar `PremiumTab.tsx` — handleAddExtra com fluxo de pagamento
9. Build e testar

## Riscos
- Usuário paga mas o webhook falha → canal não é ativado (já tratado com idempotência)
- Usuário cancela pagamento → nada acontece (comportamento esperado)
- Teste end-to-end difícil sem Telegram Stars reais → usar StarsTestMode

## Impactos esperados
- Canais extras agora são cobrados (35 Stars cada)
- Fluxo de pagamento consistente com assinatura inicial
- Usuário vê invoice do Telegram ao clicar em canal inativo
- Após pagamento, canal é ativado automaticamente

## Como testar

### Build
```bash
cd dashboard && npm run build
cd .. && go build ./...
```

### Testes manuais
1. Assinar premium com 1 canal (base)
2. Na tela de gerenciamento, clicar em canal inativo
3. Verificar se invoice de35 Stars é criada
4. Pagar via Telegram Stars
5. Verificar que canal é ativado
6. Remover canal extra e verificar que funciona

### Rollback
- Reverter mudanças no `PremiumTab.tsx` (voltar para `addExtraChannel()` direto)
- Manter endpoints novos mas sem uso (não quebra nada)

## Observações
- O `addExtraChannel()` antigo continua existendo para compatibilidade (caso queiram adicionar canais sem cobrar em futuro)
- O preço do extra channel é dinâmico (tabela `premium_features` com key `"extra_channels"`)
- Test mode (StarsTestMode) força preço para1 Star
