# Plano: Corrigir Fundo dos Logs, Footer de Configurações, Grid 3x2 no Broadcast e Header Blur

## Pedido do Usuário
1. **Logs com fundo branco**: A aba de Logs (`AdminLogsTab`) possui blocos/cards com fundo branco.
2. **Botão "Salvar" (Configurações)**: A div do botão salvar (`card-footer`) na aba de Configurações continua com fundo branco.
3. **Público-Alvo no Broadcast (Grid 3x2)**: A seleção de público-alvo (6 opções: Todos, Canais, Usuários, Suporte, IDs Usuários, IDs Canais) deve ser exibida em **3 colunas por linha x 2 linhas**, lado a lado.
4. **Header Transparente com Blur**: O header da admin dashboard (`.admin-topbar`) deve ser transparente com efeito de vidro fosco (`backdrop-filter: blur(16px)`).

---

## Diagnóstico Técnico

### 1. Fundo dos Logs (`AdminLogsTab.tsx` e `index.css`)
- O CSS possui regras `.admin-logs-page > div:first-child`, `.admin-logs-page > div:nth-child(3) > div`, e `.bg-muted/10`, `.bg-muted/20` com fundo `#ffffff` e `#f8fafb` hardcoded.
- Solução: Substituir por `var(--card)`, `var(--bg)` e `var(--surface)` no `index.css` e em `AdminLogsTab.tsx`.

### 2. Div do Botão Salvar nas Configurações (`index.css`)
- Os seletores `.admin-layout-v2 [data-slot='card-footer']`, `.minimal-section-card-footer` e `.admin-config-footer` ainda possuem `background: #fafaf8` ou `background: #fbfcfc` hardcoded no CSS.
- Solução: Alterar para `background: var(--surface)` (para cartões em geral) ou `background: var(--card)` (para a aba de configurações), eliminando a faixa branca.

### 3. Grid 3x2 no Broadcast (`index.css` & `AdminNoticeTab.tsx`)
- A seleção de público possui exatamente 6 alvos (Todos, Canais, Usuários, Suporte, IDs Usuários, IDs Canais).
- Em telas de desktop e tablet, definir `.admin-layout-v2 .broadcast-targets { display: grid; grid-template-columns: repeat(3, 1fr); gap: 10px; }`, criando exatamente **3 colunas x 2 linhas (6 cards)** organizados lado a lado no esquadro.
- Remover fundos `#fafaf8`, `#f7f7f5` remanescentes do composer e do painel de prévia.

### 4. Header Transparente com Blur (`index.css`)
- Ajustar `.admin-layout-v2 .admin-topbar` para:
  - `position: sticky; top: 0; z-index: 40;`
  - `background: var(--nav-bg, rgba(15, 17, 21, 0.75));` (transparente adaptativo)
  - `backdrop-filter: blur(16px); -webkit-backdrop-filter: blur(16px);`

---

## Arquivos que serão modificados

- `dashboard/src/index.css` (header blur, fundo dos logs, card-footer, grid 3x2 do broadcast)
- `dashboard/src/components/AdminLogsTab.tsx` (remover bg hardcoded em expansão de metadata)
- `dashboard/src/components/AdminNoticeTab.tsx` (grid 3x2 dos 6 botões de alvos)

---

## Passos Detalhados

### Passo 1 — Header Transparente com Blur (`index.css`)
- Atualizar `.admin-layout-v2 .admin-topbar`:
  - `position: sticky; top: 0; z-index: 40;`
  - `background: var(--nav-bg); backdrop-filter: blur(16px); -webkit-backdrop-filter: blur(16px);`

### Passo 2 — Eliminar Fundo Branco na Aba de Logs (`index.css` & `AdminLogsTab.tsx`)
- `.admin-logs-page > div:first-child`, `.admin-logs-page > div:nth-child(3) > div`, `.admin-logs-page .bg-muted\/10`, `.admin-logs-page .bg-muted\/20`:
  - Substituir cores fixas por `var(--card)` e `var(--surface)`.
- No `AdminLogsTab.tsx`: substituir `style={{ background: 'var(--muted)' }}` do bloco `<pre>` de metadata por `bg-muted/40 text-foreground border-border`.

### Passo 3 — Eliminar Fundo Branco na Div do Botão Salvar (`index.css`)
- `.admin-layout-v2 [data-slot='card-footer']`, `.minimal-section-card-footer`, `.admin-config-footer`:
  - Definir `background: var(--surface); border-top-color: var(--border);`.

### Passo 4 — Reorganizar Público no Broadcast em Grid 3x2 (`index.css`)
- `.admin-layout-v2 .broadcast-targets`:
  - `display: grid; grid-template-columns: repeat(3, minmax(0, 1fr)); gap: 10px;`
  - Garante exatamente 3 itens por linha em 2 linhas (6 itens no total).
- Eliminar qualquer `background: #fafaf8` em `.broadcast-submit-row`, `.broadcast-button-row-top` e `.broadcast-preview-stage`.

---

## Como Testar
```bash
cd dashboard && npm run build
```
- Verificar a compilação limpa do Vite.
- Testar visualmente a transparência/blur da topbar, a remoção do fundo branco dos logs e do footer de configurações, e o grid 3x2 do Broadcast.

## Rollback
```bash
git checkout dashboard/src/index.css dashboard/src/components/AdminLogsTab.tsx dashboard/src/components/AdminNoticeTab.tsx
```
