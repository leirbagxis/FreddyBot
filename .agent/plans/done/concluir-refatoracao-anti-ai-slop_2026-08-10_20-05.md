# Plano: Concluir Refatoração Anti-AI-Slop na Dashboard Inteira

## Pedido do usuário
Aplicar integralmente o "Anti-AI-Slop UI Design Prompt" em TODA a dashboard do FreddyBot, de forma completa: interface limpa, utilitária, tipografia como hierarquia primária, sem cards decorativos, sem sombras em conteúdo, sem gradientes/glassmorphism, sem ícones em caixas coloridas, sem subtítulos óbvios e sem badges/pills decorativos.

## Objetivo
Finalizar a varredura anti-slop em 100% da interface (admin + painel do usuário), eliminando os últimos vestígios de estética "SaaS gerada por IA" que permanecem após as refatorações parciais em andamento no working tree.

## Contexto atual
- Existem 2 planos anteriores (`.agent/plans/done/refatorar-dashboard-anti-ai-slop_2026-08-10_19-50.md` e `refatorar-paineis-admin-anti-ai-slop_2026-08-10_19-53.md`) com implementação PARCIAL não commitada (14 arquivos modificados: `MetricCard`, `OperationsOverview`, `ActivityFeed`, `AlertsPanel`, `FinanceSummary`, `SystemHealth`, `AdminAuditTab`, `DashboardInicioTab`, `CaptionCard`, `ReactionsCard`, `ButtonGrid`, `PermissõesCard`, `ContaTelegramTab`, `ui/card`, `index.css`).
- O build atual passa (`npm run build`), então as mudanças em andamento não quebram nada.
- Auditoria de slop remanescente encontrou violações em ~25 arquivos de componente e no `index.css` (caixas de ícones `size-* rounded-xl/full` com `accent-soft`, `content-card`/`action-card` com radius 16px + sombra em hover, `backdrop-blur` em topbar/toasts/bottom-nav/dialog, `shadow-sm/md/lg` em conteúdo estático, `rounded-2xl`, pills `rounded-full`, gradiente shimmer em grid cells, divider SVG ondulado decorativo, `.dark` fallback com rgba glass).

## Arquivos analisados
- `dashboard/src/index.css` (4096 linhas, classes legadas `.action-card`, `.content-card`, `.section-icon`, `.admin-card`, `.top-bar`, `.bottom-nav`, `.toast-msg`, `.admin-topbar-search`, `.sidebar-nav-item.active`, `.grid-cell.occupied:hover`, `.dark`)
- `dashboard/src/App.tsx` (telas de canais, permissões, empty state, modais de templates)
- `dashboard/src/components/AdminDashboard.tsx` (detalhe do usuário)
- `dashboard/src/components/AdminNoticeTab.tsx`, `AdminAuditTab.tsx`, `AdminLogsTab.tsx`, `AdminConfigTab.tsx`, `AdminSubscriptionsTab.tsx`, `AdminPremiumFeaturesTab.tsx`, `AdminMTProtoAccountsTab.tsx`
- `dashboard/src/components/ScheduleTab.tsx`, `ConfirmModal.tsx`, `PremiumTab.tsx`, `PremiumConfigTab.tsx`, `ConnectedAccountCard.tsx`, `NativeReactionsCard.tsx`, `AuthFlow.tsx`, `AuthDisclaimer.tsx`, `ButtonGrid.tsx`, `NewPackCaptionCard.tsx`, `DashboardInicioTab.tsx`, `Toast.tsx`, `WaveDivider.tsx`
- `dashboard/src/components/ui/dialog.tsx` (overlay com blur, radius xl)

## Arquivos que serão modificados
- `dashboard/src/index.css`
- `dashboard/src/App.tsx`
- `dashboard/src/components/AdminDashboard.tsx`
- `dashboard/src/components/AdminNoticeTab.tsx`
- `dashboard/src/components/AdminAuditTab.tsx`
- `dashboard/src/components/AdminLogsTab.tsx`
- `dashboard/src/components/AdminMTProtoAccountsTab.tsx`
- `dashboard/src/components/AdminPremiumFeaturesTab.tsx`
- `dashboard/src/components/AdminSubscriptionsTab.tsx`
- `dashboard/src/components/ScheduleTab.tsx`
- `dashboard/src/components/ConfirmModal.tsx`
- `dashboard/src/components/PremiumTab.tsx`
- `dashboard/src/components/PremiumConfigTab.tsx`
- `dashboard/src/components/ConnectedAccountCard.tsx`
- `dashboard/src/components/NativeReactionsCard.tsx`
- `dashboard/src/components/AuthFlow.tsx`
- `dashboard/src/components/AuthDisclaimer.tsx`
- `dashboard/src/components/ButtonGrid.tsx`
- `dashboard/src/components/NewPackCaptionCard.tsx`
- `dashboard/src/components/WaveDivider.tsx`
- `dashboard/src/components/ui/dialog.tsx`

## Estratégia de implementação

### 1. CSS global (`index.css`)
- **`.action-card`**: vira linha plana: border-radius 6px, sem `box-shadow`, sem `transform` em hover/active, hover apenas `background: var(--surface-hover)`.
- **`.action-card-icon`**: removido/neutralizado (ícone inline sem fundo colorido).
- **`.content-card`**: radius 8px, sem hover de borda `accent-soft` decorativo.
- **`.content-card-icon`**, **`.section-icon`** (e variantes purple/green/rose/amber + `::after` glow + drop-shadow dark): removidos.
- **`.top-bar`**, **`.bottom-nav`**, **`.toast-msg`**: remover `backdrop-filter: blur`, superfícies opacas planas (`nav-bg` sólido). Toasts/bottom-nav mantêm sombra (camadas flutuantes).
- **`.admin-topbar-search`**: mantém hover/focus (estado interativo), sem blur.
- **`.sidebar-nav-item.active`**: remover `box-shadow` glow, manter `accent-soft` background.
- **`.grid-cell.occupied:hover`**: remover gradient shimmer, hover apenas borda accent.
- **`.dark` fallback**: substituir `rgba(...)` glass por cores sólidas.
- **`.admin-card*` legado**: remover somente se sem uso em TSX (verificar com grep); deixar se ainda usado.
- **`--shadow-sm/--shadow-md`**: garantir `none` em conteúdo; sombras apenas em `--shadow-lg` (camadas flutuantes).

### 2. App.tsx (painel do usuário)
- Greeting/Conta/Templates/lista de canais: trocar `action-card` por `row` plana com divider (`divide-y divide-border`), ícones inline neutros, sem chevron-decorativo exagerado (remover `action-card-icon`).
- Empty state de canais: remover caixa `size-16 rounded-2xl accent-soft`, remover subtítulos óbvios, manter apenas instrução + CTA.
- Aba Permissões (content-card): seções planas com títulos de seção, toggles em linhas com `divide-y`, badges ON/OFF substituídos por texto de estado.
- Header de templates: remover `bg-background/90 backdrop-blur` → superfície sólida.
- Cálculo de greening mantido.

### 3. Admin — detalhe de usuário (`AdminDashboard.tsx`)
- Remover caixa `rounded-xl border p-4` → seção plana com título + ações alinhadas.
- Remover `size-12 rounded-xl` icon box do avatar → avatar simples com inicial.
- Lista "Canais do Usuário": linhas planas com dividers (`divide-y`), sem `rounded-xl border p-3` por item.
- Empty state `rounded-xl border` → linha de texto simples.

### 4. Admin — abas operacionais
- **`AdminNoticeTab`**: bolha de preview do broadcast: `rounded-2xl shadow-sm` → `rounded-md` sem sombra (mantém função de preview); botões preview `rounded-xl border` → linhas de texto/flat; remover `broadcast-section-heading` se decorativo.
- **`AdminAuditTab`**: remover `section-icon` boxes (inclusive variantes), `rounded-xl border` por item → listas `divide-y`; pill `rounded-full` de status → StatusBadge/ou texto.
- **`AdminLogsTab`**: inputs/selects `rounded-xl` → 6px; cards `rounded-xl border overflow-hidden` → linhas planas com separadores; empty state sem caixa.
- **`AdminMTProtoAccountsTab`**: remover `size-10 rounded-full success-soft` box.
- **`AdminPremiumFeaturesTab`**: remover boxes `size-10 rounded-xl accent-soft`.
- **`AdminSubscriptionsTab`**: remover `size-10 rounded-xl` box no modal de confirmação.

### 5. Admin — resto
- **`ScheduleTab`**: Card `hover:shadow-sm` → superfície plana; dialog `shadow-lg` mantém (camada flutuante); remover bordas/radius excessivos em itens de agenda; spinner mantido.
- **`ConfirmModal`**: remover `rounded-2xl` icon box (ícone inline), botões `shadow-sm/md` → flat (`rounded-md`).
- **`AdminTopbar`/`AdminSidebar`**: já estão limpos; apenas garantir estado ativo sem glow (via CSS).

### 6. Usuário — premium/conta/auth
- **`PremiumTab`**, **`PremiumConfigTab`**: remover icon boxes `size-14/11/10 rounded-xl/2xl` e `rgba(168,85,247,.1)`; plan sections.
- **`ConnectedAccountCard`**: círculos `size-8 rounded-full accent-soft` → linhas planas (avatar da conta mantém `rounded-full`, é funcional).
- **`NativeReactionsCard`**: remover `size-10 rounded-xl accent-soft`.
- **`AuthFlow`**, **`AuthDisclaimer`**: remover `rounded-full bg-accent/10` e `rounded-xl` box → layout centralizado simples.
- **`ButtonGrid`**: pill de preview do botão: manter como preview funcional mas sem `shadow-sm`/`drop-shadow-sm` e com superfície opaca; caixas `rounded-xl bg-muted border` de ações → linhas planas.
- **`NewPackCaptionCard`**: remover `content-card-icon` colorido (amarelo warning-soft).
- **`WaveDivider` (PerfLine)**: substituir SVG ondulado por `<hr>` simples 1px (separador sutil), sem componente decorativo de 18px.

### 7. Componentes base
- **`ui/dialog.tsx`**: overlay `bg-black/60 backdrop-blur-md` → `bg-black/60` sem blur; painel `rounded-xl` → `rounded-md` (sombra mantida — camada flutuante).
- **`ui/select.tsx`**: manter sombra no dropdown (flutuante), sem alteração.

## Passos detalhados
1. Auditoria final de classes legadas com grep (garantir que `.admin-card*` não está em uso antes de remover).
2. Editar `index.css` (item 1 da estratégia).
3. Editar `App.tsx`.
4. Editar componentes admin (4 arquivos de detalhe + abas).
5. Editar componentes do usuário (premium/conta/auth/grid/preview).
6. Editar `ui/dialog.tsx` e `WaveDivider.tsx`.
7. Rodar `npm run build` e verificar zero erros de TypeScript.
8. Passe de limpeza final: grep por `shadow-|rounded-2xl|backdrop-blur|bg-gradient` em `src` para confirmar que só restam casos funcionais (diálogos/dropdowns/toasts).

## Riscos
- **Baixo**: alterações puramente visuais; sem mudança de contratos de dados/API.
- **Médio**: arquivo `index.css` grande com classes inter-relacionadas; mitigar rodando build + revisando seletores antes de remover.

## Impactos esperados
- Interface 100% coerente com as regras anti-slop: superfícies planas, tipografia como hierarquia, densidade funcional alta, sem sombras decorativas/blur/gradientes/icon-boxes/pills supérfluos.
- Manutenção visual simplificada (menos classes legadas).

## Compatibilidade
- Linux, macOS, Windows, Docker, CI/CD (mudanças 100% frontend; backend intacto).

## Como testar

### Build
```bash
cd dashboard && npm run build
```

### Execução
```bash
cd dashboard && npm run dev
```

## Rollback
```bash
git checkout dashboard/src
```
(ou remover apenas o arquivo do plano: `git checkout dashboard/src/index.css dashboard/src/App.tsx`)

## Observações
- O trabalho em andamento no working tree (14 arquivos modificados) será PRESERVADO; este plano apenas completa o que falta.
- Componentes 100% funcionais (RichTextEditor, EmojiPicker, Toast lógica, TabBar navegação) não mudam de contrato, apenas de CSS quando necessário.

---

# IMPLEMENTAÇÃO CONCLUÍDA (2026-08-10)

## Resumo da execução
Todos os passos 1-8 do plano foram executados e o build final está verde:

- **index.css**: `.action-card`/`.content-card` planos (radius 6px, hover só background); `.action-card-icon`/`.content-card-icon`/`.section-icon` neutros inline (variantes purple/green/rose/amber + `::after` glow + drop-shadow removidos); `.top-bar`/`bottom-nav`/`.toast-msg`/`.admin-sidebar` sem `backdrop-filter` (bottom-nav radius 100px→8px, toasts 14px→8px); fallback `.dark` com cores sólidas (`--card: #121215`, `--border: #27272a`, etc.); blocos mortos removidos (`.admin-card*`, `.metric-card*`, Activity/Health/Finance/Alert legacy, `.skeleton`, `@keyframes pulse-glow`, regras `.section-icon`/`.rounded-2xl` do audit/notice); sombras de conteúdo sobreviventes removidas (`.admin-audit-page .bg-card`, hover `.admin-features-page`, header `.admin-subscriptions-page`); `[data-slot='dialog-content']` radius 20px→8px sem blur; `.rte-container` 14px→8px; `.cp-bubble` 16px→8px.
- **App.tsx**: header sólido `bg-background`; lista de canais `divide-y divide-border`; empty state sem icon-box `size-16 rounded-2xl`; seções Reações/Links Dinâmicos/Permissões em linhas planas (`rounded-md hover:bg-muted/50`); Badges ON/OFF → texto "Ativado/Desativado"; tela de acesso negado sem caixa `64px danger-soft`; import `Badge` removido.
- **Admin**: `AdminDashboard` (avatars com inicial `rounded-full`, lista de canais `divide-y border rounded-md`, `tabTitles` sem descrição), `AdminAuditTab` (canais flat + empty state plano), `AdminNoticeTab` (preview `rounded-md` sem shadow, badges → span), `AdminLogsTab` (inputs/selects `rounded-md`, eventos `divide-y`, metadata e `<pre>` sem caixas), `AdminMTProtoAccountsTab` (icon-box do "done" removido), `AdminPremiumFeaturesTab` (icon boxes removidos), `AdminSubscriptionsTab` (icon-box do dialog de confirmação removido).
- **Usuário**: `PremiumTab` (trigger flat com Crown inline, benefícios `divide-y`, canal select flat, sem Badges — texto de estado), `PremiumConfigTab` (header inline + status como texto), `ConnectedAccountCard` (linhas `divide-y border-y` com ícones inline), `NativeReactionsCard` (preview/badges em texto, header inline), `AuthFlow` (sucesso sem icon-box `72px`), `AuthDisclaimer` (sem `rounded-[20px]` pills nem icon-box), `ScheduleTab` (Card sem `hover:shadow-sm`), `ConfirmModal` (sem icon-box `64px rounded-2xl`, botões `rounded-md h-10` sem sombra), `ButtonGrid` (pill de reações sem `shadow-sm`/`drop-shadow-sm`, sheets `rounded-md`, dialog sem `section-icon`), `NewPackCaptionCard` (switch rows flat `rounded-md`, ícone neutro), `UserTemplatesManager` (icon-box removido), `WaveDivider` → `PerfLine` virou `<hr>` 1px.
- **ui/dialog.tsx**: overlay `backdrop-blur-md` removido, painel/footer `rounded-md` (sombra mantida — camada flutuante).
- **Casos funcionais mantidos**: avatares circulares com inicial, switches, spinners de loading, barras de progresso, popover/dropdown (`shadow-md`), modais/overlays (`shadow-lg/2xl`), toasts, chips de status (radius 6px), botões de seleção de emoji.
- **Verificação**: `cd dashboard && npm run build` OK (24-25s; único warning = chunk > 500 kB, pré-existente). Fix de sintaxe: JSX desbalanceado no `PremiumTab` após edição do bloco de gestão (div `</div>` órfã removida).