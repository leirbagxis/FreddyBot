# Plano: Ajuste Fino dos Componentes, Alinhamento e Tema Telegram

## Pedido do Usuário
1. **Textos escuros / Divs brancas**: Div do botão salvar (aba configurações), item selecionado do menu lateral e seções da aba Broadcast ficaram com fundo branco.
2. **Fila de Revisão desalinhada**: Itens da fila na Visão Geral (`OperationsOverview`) ficaram fora do esquadro.
3. **Filtros na aba Assinaturas**: Botões de filtro muito pequenos e quase invisíveis.
4. **Toggles invisíveis nas Configurações**: Switches/toggles ficaram apagados/invisíveis em temas escuros.

---

## Análise das Causas

### 1. Divs Brancas e Textos Escuros (`index.css`)
- `.admin-layout-v2 [data-slot='card-footer']` possui `background: #fafaf8` / `#fbfcfc` hardcoded no CSS, forçando fundo branco na div do botão salvar.
- `.sidebar-nav-item.active` ou `.sidebar-nav-item:hover` possui `background: #f2f2ef` hardcoded.
- `.broadcast-targets button.is-selected`, `.broadcast-button-row-top`, `.broadcast-submit-row`, `.broadcast-preview-stage` usam `#fafaf8`, `#f7f7f5`, `#efefed` hardcoded no CSS.
- Regras CSS no final do arquivo ainda forçam `color: #1b1b1a !important` em headers e títulos.

### 2. Desalinhamento da Fila de Revisão (`OperationsOverview.tsx` e `index.css`)
- O CSS `.operations-queue-list button` contava com uma estrutura específica de grid (`32px minmax(0, 1fr) auto 18px`). A alteração para a classe do shadcn `<Button className="w-full justify-between...">` alterou o modelo de caixa.
- Solução: Ajustar a estrutura do botão da fila para flexbox com alinhamento vertical centralizado (`align-items: center`), mantendo o avatar, detalhes do usuário, Badge do shadcn e o ícone alinhados perfeitamente.

### 3. Filtros na aba Assinaturas (`AdminSubscriptionsTab.tsx`)
- Os botões de filtro usavam a classe `text-[10px]` com `text-muted-foreground/60` (opacidade de 60% em texto de 10px), tornando-os extremamente pequenos e ilegíveis.
- Solução: Aumentar o tamanho para `text-xs px-3 py-1.5 rounded-lg` com bom contraste para estados ativo e inativo.

### 4. Toggles Invisíveis (`ui/switch.tsx`)
- O componente `Switch` shadcn dependia de `data-unchecked:bg-input` e `SwitchPrimitive.Thumb` com `bg-background`.
- Em temas escuros do Telegram, `var(--input)` e `var(--background)` possuem a mesma cor do fundo ou são transparentes, tornando o botão de toggle e a "bolinha" invisíveis.
- Solução: Atualizar o trilho `data-unchecked` para `bg-muted border border-border` e o knob `SwitchPrimitive.Thumb` para `bg-foreground dark:bg-white` (com sombra e alto contraste).

---

## Arquivos que serão modificados

- `dashboard/src/index.css` (remover todos os fundos `#hex` brancos remanescentes e forçar herança do Telegram)
- `dashboard/src/components/ui/switch.tsx` (toggles com alto contraste)
- `dashboard/src/components/admin/OperationsOverview.tsx` (esquadro e alinhamento perfeito da fila de revisão)
- `dashboard/src/components/AdminSubscriptionsTab.tsx` (filtros com tamanho `text-xs` e legibilidade aprimorada)
- `dashboard/src/components/admin/AdminSidebar.tsx` (estilo ativo transparente/adaptativo ao tema)

---

## Passos Detalhados de Implementação

### Passo 1 — Atualizar `ui/switch.tsx`
- Trilho: `data-checked:bg-primary data-unchecked:bg-muted/80 border border-border`
- Thumb: `bg-white dark:bg-white shadow-sm`
- Tamanho: `h-5 w-9` com thumb `size-4` para visual nítido e padrão.

### Passo 2 — Atualizar `AdminSubscriptionsTab.tsx`
- Ajustar os botões de filtro `all`, `active`, `cancelling`, `expired`, `with_charge`:
  - Tamanho: `text-xs px-3 py-1.5 font-medium rounded-lg`
  - Ativo: `bg-accent text-accent-foreground shadow-sm`
  - Inativo: `bg-muted/40 text-muted-foreground hover:bg-muted hover:text-foreground`

### Passo 3 — Corrigir `OperationsOverview.tsx` e CSS da Fila de Revisão
- Ajustar o container e botão de item da fila:
  - Usar estrutura flex com `align-items: center`, `gap-3`, `py-3 px-4`, preservando os badges e seta alinhados no esquadro.

### Passo 4 — Ajustar `index.css` (Zero Fundos Brancos Hardcoded)
- `.admin-layout-v2 [data-slot='card-footer']`: `background: var(--surface)` (evita div branca no botão Salvar em Configurações).
- `.sidebar-nav-item.active`: `background: var(--accent-soft); color: var(--accent);`
- Aba Broadcast: `.broadcast-targets button`, `.broadcast-button-row-top`, `.broadcast-submit-row`, `.broadcast-preview-stage`: substituir `#fafaf8`, `#f7f7f5`, `#efefed` por `var(--surface)`, `var(--card)` e `var(--accent-soft)`.
- Titulares e Headers: remover `color: #1b1b1a !important`, substituir por `color: var(--text) !important`.

---

## Como Testar
```bash
cd dashboard && npm run build
```
Verificar a compilação do Vite e inspecionar a interface.

## Rollback
```bash
git checkout dashboard/src/index.css dashboard/src/components/ui/switch.tsx dashboard/src/components/admin/OperationsOverview.tsx dashboard/src/components/AdminSubscriptionsTab.tsx
```
