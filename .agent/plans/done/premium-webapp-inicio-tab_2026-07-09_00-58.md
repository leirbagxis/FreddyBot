# Plano: premium-webapp-inicio-tab

## Pedido do usuario
1. Mover a UI de compra do plano premium para a aba "Inicio" (DashboardInicioTab e/ou pagina de lista de canais), removendo a aba "Premium" separada.
2. Fazer todo o fluxo de pagamento via WebApp (`Telegram.WebApp.openInvoice()`), sem precisar o usuario voltar ao chat do bot.
3. Ter uma flag `?test=true` na hora de criar a invoice para testar o pagamento por Stars direto do Telegram.

## Objetivo
- Refatorar o fluxo de assinatura premium para ser 100% WebApp
- Integrar a UI de premium na aba Inicio
- Remover dependencia de `sendInvoice` (que envia msg pro chat do bot) em favor de `createInvoiceLink` (que retorna URL para WebApp)
- Adicionar `?test=true` no endpoint de create: em dev, ativa a assinatura sem pagamento real para testes

## Contexto atual
- **PremiumTab** (`dashboard/src/components/PremiumTab.tsx`) — componente separado com UI de assinar/gerenciar premium, renderizado em aba propria (`activeTab === 'premium'`)
- **DashboardInicioTab** (`dashboard/src/components/DashboardInicioTab.tsx`) — componente da aba "Inicio", mostra saudacao + informacoes do canal + transferencia de posse
- **App.tsx** — tabs incluem `{ id: 'premium', label: 'Premium', icon: <Crown> }`; renderiza `<PremiumTab>` quando `activeTab === 'premium'`
- **SubscriptionService.CreateInvoice** — usa `bot.SendInvoice()` que envia invoice como mensagem no chat do bot; usuario precisa clicar la e pagar, depois voltar ao WebApp
- **SubscriptionController.CreateInvoice** — retorna `payload` e `totalStars`; frontend mostra toast "Invoice enviada! Verifique seu chat com o bot"
- **api.ts** — `createSubscriptionInvoice()` chama `POST /api/subscription/create`
- **Telego v1.9.0** ja tem `CreateInvoiceLink` que retorna URL para usar com `WebApp.openInvoice()`

## Arquivos analisados
- `dashboard/src/App.tsx`
- `dashboard/src/components/PremiumTab.tsx`
- `dashboard/src/components/DashboardInicioTab.tsx`
- `dashboard/src/api.ts`
- `dashboard/src/types.ts`
- `internal/core/services/subscription_service.go`
- `internal/api/controllers/subscription_controller.go`
- `internal/api/routes/routes.go`
- `internal/container/appContainer.go`

## Arquivos que serao modificados

### Backend
- `internal/core/services/subscription_service.go` — `CreateInvoice` passa a usar `CreateInvoiceLink` + aceita `testMode` para ativar direto em dev
- `internal/api/controllers/subscription_controller.go` — le query param `test`, chama service com flag, retorna invoiceUrl

### Frontend
- `dashboard/src/App.tsx` — remove aba Premium, renderiza premium dentro da channels list
- `dashboard/src/components/PremiumTab.tsx` — refatorado: integrado como card reutilizavel, usa `WebApp.openInvoice()`, mostra botao "Testar" em dev
- `dashboard/src/components/DashboardInicioTab.tsx` — adiciona secao premium com status + compra + gerenciamento
- `dashboard/src/api.ts` — `createSubscriptionInvoice` aceita `test` param, retorna `invoiceUrl`
- `dashboard/src/types.ts` — estende `Telegram.WebApp` com `openInvoice` + tipo `InvoiceStatus`

## Estrategia de implementacao

### Parte 1: Backend — CreateInvoiceLink + flag test
1. Em `SubscriptionService`, mudar `CreateInvoice` para aceitar `testMode bool`
2. Se `testMode && config.AppEnv == "dev"`: ativar assinatura direto (sem invoice)
3. Se normal: usar `bot.CreateInvoiceLink()` em vez de `bot.SendInvoice()`
4. Retornar `InvoiceURL string` em vez de `Message *telego.Message`
5. Atualizar `InvoiceResult` DTO

### Parte 2: Backend — Controller lê query param
1. `CreateInvoice` handler le `ctx.Query("test")`
2. Passa `testMode=true` para o service se `?test=true`
3. Retorna `invoiceUrl` no JSON

### Parte 3: Frontend — WebApp.openInvoice + flag test
1. `createSubscriptionInvoice(test?)` envia `POST /api/subscription/create?test=true` se test mode
2. Em modo normal: usa `window.Telegram.WebApp.openInvoice(invoiceUrl, callback)`
3. Callback: se `status === 'paid'`, recarrega status
4. Se falhou (`cancelled`/`failed`): mostra toast
5. Em modo test (dev): chama API com `?test=true`, API ja ativa direto, frontend so recarrega

### Parte 4: Frontend — Integrar na aba Inicio
1. Em `DashboardInicioTab`, adicionar card premium:
   - Se nao tem assinatura: mostrar beneficios + botoes "Assinar com Stars" + "Testar" (dev)
   - Se tem ativa: mostrar status + canais extras + cancelar
2. Na channels list page (isChannels), adicionar card premium apos "Conta Telegram"
3. Remover `PremiumTab` da renderizacao em `activeTab === 'premium'`
4. Remover tab "Premium" do array `tabs`
5. Manter logica do PremiumTab como componente interno (ou importado inline)

### Parte 5: types.ts — Extender WebApp interface
Adicionar:
```typescript
openInvoice: (url: string, callback: (status: InvoiceStatus) => void) => void;
```
onde `InvoiceStatus = 'paid' | 'cancelled' | 'failed' | 'pending'`

## Passos detalhados

1. **Backend: InvoiceResult** — trocar `Message *telego.Message` por `InvoiceURL string`
2. **Backend: CreateInvoice** — aceitar `testMode bool`; se testMode+dev ativar direto, se nao usar `CreateInvoiceLink`
3. **Backend: Controller** — ler `?test` query param, passar ao service, retornar invoiceUrl
4. **Frontend: types.ts** — estender `Telegram.WebApp` com `openInvoice`
5. **Frontend: api.ts** — `createSubscriptionInvoice(test?)` envia query param; retorna `{ invoiceUrl, payload, totalStars }`
6. **Frontend: PremiumTab** — refatorar: exportar funcoes compartilhadas; `handleCreateInvoice` usa `WebApp.openInvoice()`; fallback pra dev (abre link em nova aba); botao "Testar ⚡" visivel em dev
7. **Frontend: DashboardInicioTab** — adicionar secao premium reutilizando logica do PremiumTab
8. **Frontend: App.tsx** — remover tab `premium`; remover renderizacao `activeTab === 'premium'`; adicionar premium na channels list page
9. **Build & Test** — `go build ./...` + `npm run build` no frontend

## Riscos
- `CreateInvoiceLink` requer Bot API 7.5+ (telego v1.9.0 suporta)
- `WebApp.openInvoice` so funciona dentro do Telegram; em navegador (dev), fazer fallback abrindo URL em nova aba
- Test flag `?test=true` nunca deve funcionar em producao (verificar `config.AppEnv`)
- Usuario fechar WebApp antes do callback — proxima abertura recarrega status via `fetchSubscriptionStatus`

## Impactos esperados
- Fluxo de pagamento 100% dentro do WebApp (melhor UX)
- Nao envia mais invoice pro chat do bot
- Aba Inicio mais rica
- Test flag facilita desenvolvimento/testes sem gastar Stars

## Compatibilidade
- Linux OK
- macOS OK
- Windows OK
- Docker OK
- CI/CD OK

## Como testar

### Build
```bash
cd /home/malbs/Opencode/FreddyBot && go build ./...
```

### Testes
```bash
cd /home/malbs/Opencode/FreddyBot && go test ./internal/...
```

### Frontend
```bash
cd /home/malbs/Opencode/FreddyBot/dashboard && npm run build
```

### Execucao (modo dev)
```bash
make run
# Abrir WebApp no Telegram
# Navegar ate "Meus Canais" -> ver card premium
# Clicar "Assinar com Stars" -> deve abrir overlay de pagamento via openInvoice
# Clicar "Testar ⚡" -> ativa direto sem pagar (dev only)
```

## Rollback
- Reverter `subscription_service.go` para usar `SendInvoice`
- Restaurar tab Premium em `App.tsx`
- Restaurar `DashboardInicioTab.tsx` sem secao premium

## Observacoes
- A flag `?test=true` so funciona em modo dev (`config.AppEnv == "dev"`). Em producao, ignorada.
- `InvoiceStatus` possivel: `'paid'`, `'cancelled'`, `'failed'`, `'pending'`
- `CreateInvoiceLinkParams` aceita os mesmos campos de `SendInvoiceParams` exceto `ChatID`
