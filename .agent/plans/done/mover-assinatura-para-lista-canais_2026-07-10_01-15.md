# Plano: mover-assinatura-para-lista-canais

## Pedido do usuario
Mover a parte de comprar/assinar premium para a pagina de lista de canais (`/me/channels`), remover o card de assinatura de cada canal individual e colocar apenas um status (ativado/desativado). Antes de assinar, o usuario deve escolher quais canais quer incluir no premium.

## Objetivo
- **Channel dashboard** (`/dashboard/:id`): Remover `PremiumTab`, substituir por badge simples
- **Channels list** (`/me/channels`): Manter assinatura com seletor de canais
- **Selecao de canais**: Antes de pagar, usuario escolhe quais canais entram no premium
- **Preco dinamico**: basePrice + (nCanaisSelecionados - 1) * extraChannelPrice

## Contexto atual
- `PremiumTab` aparece em 2 lugares: dentro de `DashboardInicioTab` (por canal) + lista de canais
- Assinatura e user-level (nao por canal)
- `ExtraChannels` e um contador, nao tracking por canal
- Preco = basePrice (ex: 80) + extraChannels * extraChannelPrice (ex: 35)
- SubscriptionController aceita `?test=true` para bypass

## Arquivos analisados
- `dashboard/src/components/PremiumTab.tsx` — componente atual
- `dashboard/src/components/DashboardInicioTab.tsx` — onde PremiumTab aparece por canal
- `dashboard/src/App.tsx` — onde PremiumTab aparece na lista de canais + routes
- `dashboard/src/api.ts` — chamadas de API de subscription
- `dashboard/src/types.ts` — tipos (SubscriptionStatus, Channel, etc.)
- `internal/api/controllers/subscription_controller.go` — endpoints
- `internal/core/services/subscription_service.go` — logica

## Arquivos que poderao ser modificados

### Frontend
- `dashboard/src/components/PremiumTab.tsx` — REFACTOR: adicionar seletor de canais
- `dashboard/src/components/DashboardInicioTab.tsx` — remover PremiumTab, adicionar badge
- `dashboard/src/App.tsx` — mover fluxo de assinatura
- `dashboard/src/types.ts` — tipagens se necessario
- `dashboard/src/api.ts` — se precisar de novos endpoints

### Backend (minimo possivel)
- `internal/api/controllers/subscription_controller.go` — aceitar channelIds opcional
- `internal/core/services/subscription_service.go` — calcular preco com base em channelIds

## Estrategia de implementacao

### Fase 1 — Backend (minimo)
1. Modificar `CreateInvoice` para aceitar `channelCount` extra (alem do existing.ExtraChannels)
2. SubscriptionController: ler `?channels=N` query param
3. Calcular preco: totalStars = basePrice + max(0, channelCount - 1) * extraChannelPrice
4. Subscription: salvar ExtraChannels = channelCount - 1 (ou manter como esta)

Na verdade, para manter a compatibilidade e simplicidade, vou fazer:
- O frontend calcula o numero de canais extras
- Passa como query param `?channels=3` (total de canais selecionados)
- Backend calcula: extras = max(0, channels - 1), total = basePrice + extras * extraChannelPrice
- Salva no subscription.ExtraChannels

### Fase 2 — Frontend: PremiumTab refactor
1. Adicionar props: `channels: Channel[]` e `onStatusChange`
2. Estado "nao assinante" agora mostra:
   - Lista de canais do usuario com checkboxes
   - Preco base + (selecionados - 1) * extraPrice
   - Botao "Assinar N canais por X Stars"
3. Estado "assinante" mantido, mas com indicacao de quantos canais

### Fase 3 — Frontend: DashboardInicioTab
1. Remover `<PremiumTab>` do `DashboardInicioTab.tsx`
2. Adicionar badge/card pequeno mostrando status:
   - Se usuario tem subscription ativa → "Premium Ativo" com check verde
   - Senao → "Premium Inativo" com indicacao de onde assinar

### Fase 4 — Frontend: App.tsx channels list
1. Manter PremiumTab na lista de canais
2. Passar `channels={user?.channels}` como prop
3. Integrar channel selection no fluxo de assinatura

## Passos detalhados

### Passo 1 — Backend: modificar CreateInvoice
Arquivo: `internal/core/services/subscription_service.go`
- Adicionar parametro `channelCount int` em `CreateInvoice`
- Calcular: `extraChannels := max(0, channelCount-1)`
  - Se `existing != nil`: `extraChannels = max(existing.ExtraChannels, channelCount-1)`
- Atualizar `totalStars`:
  - `totalStars = basePrice + extraChannels * extraChannelPrice`
- Passar `extraChannels` para `activateSubscription`

### Passo 2 — Backend: SubscriptionController
Arquivo: `internal/api/controllers/subscription_controller.go`
- Ler query param `?channels=N` (opcional, default 1)
- Passar para `svc.CreateInvoice(ctx, userID, testMode, channelCount)`

### Passo 3 — Frontend: types.ts
- Sem alteracoes por enquanto (Channel ja tem os campos necessarios)

### Passo 4 — Frontend: PremiumTab.tsx
- Adicionar prop `channels: Channel[]`
- Estado "nao assinante":
  - Renderizar lista de canais com checkbox cada
  - Calcular: `selectedCount = channelsSelecionados.length`
  - Calcular: `extras = max(0, selectedCount - 1)`
  - Calcular: `total = basePrice + extras * extraChannelPrice`
  - Botao: "Assinar {selectedCount} canais por {total} Stars"
- Estado "assinante":
  - Manter visual atual (periodo, canais extras, cancelar)
  - Mostrar quantos canais estao inclusos

### Passo 5 — Frontend: DashboardInicioTab.tsx
- Remover secao `PremiumTab`
- Adicionar badge no header ou card:
  - Se subscription ativa: `<Badge variant="default">Premium ✓</Badge>`
  - Senao: `<Badge variant="outline">Premium —</Badge>` com link para /me/channels

### Passo 6 — Frontend: App.tsx
- Na secao de lista de canais, passar `channels` para PremiumTab
- Ajustar layout para comportar a lista com checkboxes

## Riscos
- **Quebrar fluxo existente**: Usuarios atuais podem perder acesso a assinatura se nao encontrarem
- **UX complexa**: Selecao de canais + precificacao pode confundir usuarios
- **Dados inconsistentes**: Se o backend nao tracking por canal, badge de status pode ser impreciso
- **ExtraChannels**: Atualmente e so um contador, usar como "canais alem do primeiro"

## Impactos esperados
- UI mais limpa em cada canal (sem poluicao de assinatura)
- Fluxo de compra mais claro (seleciona antes de pagar)
- Possivel aumento de conversao (usuario ve o valor por canal)
- Badge de status em cada canal da visibilidade

## Compatibilidade
- Linux ✓
- macOS ✓
- Windows ✓ (frontend)
- Docker ✓
- CI/CD ✓

## Como testar

### Build backend
```bash
go build ./internal/... ./cmd/...
```

### Build frontend
```bash
cd dashboard && npx tsc --noEmit
```

### Teste manual
1. Acessar `/me/channels` — ver lista de canais + card de assinatura
2. Clicar "Assinar" — ver canais com checkboxes, preco dinamico
3. Selecionar canais — preco atualiza
4. Confirmar — invoice criada (ou modo teste)
5. Apos ativar — badge aparece nos canais
6. Acessar `/dashboard/:id` — ver badge de status no lugar do card grande

## Rollback
1. Reverter alteracoes em `DashboardInicioTab.tsx` — restaurar `<PremiumTab>`
2. Reverter `App.tsx` — restaurar layout anterior
3. Reverter `subscription_controller.go` — remover `channels` param
4. Reverter `subscription_service.go` — restaurar assinatura original

## Observacoes
- O backend NAO vai tracking per-channel (sem alteracao no model Channel)
- O `ExtraChannels` count representa "canais alem do primeiro"
- O badge no channel dashboard mostra "Premium Ativo" se o usuario tem subscription ativa (independente de selecao)
- Futuramente podemos adicionar tracking real por canal
