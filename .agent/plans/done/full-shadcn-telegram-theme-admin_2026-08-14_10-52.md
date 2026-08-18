# Plano: Full Shadcn + Ícones Correspondentes + Tema Telegram na Dashboard Admin

## Pedido do usuário
1. Trocar todos os componentes nativos por componentes **shadcn/ui** (full shadcn)
2. Trocar todos os ícones para os **correspondentes** (lucide-react)
3. Usar o **tema do Telegram** (`useTheme.ts` → variáveis CSS) ao invés de cores hardcoded

## Objetivo
Eliminar 100% das cores fixas (`--crm-*`, `--minimal-*`, `#hex`) da dashboard admin, herdar o tema do Telegram via `useTheme.ts`, e usar componentes shadcn + ícones lucide-react em todos os componentes.

## Contexto atual

### Problema 1 — Cores hardcoded sobrescrevem o tema do Telegram
O `useTheme.ts` injeta corretamente as variáveis do Telegram (`--bg`, `--card`, `--text`, `--accent`, etc.) na `:root`, mas **3 blocos CSS** no `index.css` sobrescrevem TUDO com valores fixos:

| Bloco | Linhas | O que faz |
|---|---|---|
| `.admin-layout-v2` (1º) | 2752-2791 | Define `--crm-*` fixas (#1f1f1f, #ffffff, etc.) e redefine `--bg`, `--card`, `--text`, `--accent` |
| `.admin-layout-v2` (2º - minimal) | 3565-3619 | Define `--minimal-*` fixas e redefine todas as variáveis novamente |
| `[data-theme="dark"] .admin-layout-v2` | 3621-3667 | Dark mode hardcoded com cores de fallback |
| Overrides finais | 3988-4021 | Duplica overrides com mais `#hex` hardcoded |

### Problema 2 — Componentes nativos em vez de shadcn
- `AdminSidebar.tsx`: `<button>` nativo ao invés de `<Button variant="ghost">`
- `AdminTopbar.tsx`: `<input>` nativo, `<select>` nativo, `<button>` nativo
- `DataTable.tsx`: `<input>` nativo, `<button>` nativo para paginação
- `StatusBadge.tsx`: `<span>` com classes CSS ao invés de shadcn `<Badge>`
- `MetricCard.tsx`: `<div>` ao invés de shadcn `<Card>`
- `OperationsOverview.tsx`: `<article>`, `<section>` sem shadcn

### Problema 3 — Ícones não correspondentes
- `Hash` (canal) → deveria ser `Radio` (broadcast)
- `MessageSquare` (broadcast) → deveria ser `Megaphone`
- `Zap` (auditoria) → deveria ser `ShieldCheck`
- `FileClock` (logs) → deveria ser `ScrollText`
- `Smartphone` (contas MTProto) → deveria ser `KeyRound`
- `Crown` (features) → deveria ser `Sparkles`
- `Star` (assinaturas) → deveria ser `CreditCard`

## Arquivos analisados
- `dashboard/src/index.css` (4182 linhas, 125KB)
- `dashboard/src/hooks/useTheme.ts`
- `dashboard/src/components/admin/AdminSidebar.tsx`
- `dashboard/src/components/admin/AdminTopbar.tsx`
- `dashboard/src/components/admin/AdminLayout.tsx`
- `dashboard/src/components/admin/OperationsOverview.tsx`
- `dashboard/src/components/admin/DataTable.tsx`
- `dashboard/src/components/admin/MetricCard.tsx`
- `dashboard/src/components/admin/StatusBadge.tsx`
- `dashboard/src/components/admin/AdminOverview.tsx`
- `dashboard/src/components/admin/AdminPageHeader.tsx`
- `dashboard/src/components/ui/` (8 componentes shadcn)

## Arquivos que poderão ser modificados
- `dashboard/src/index.css`
- `dashboard/src/components/admin/AdminSidebar.tsx`
- `dashboard/src/components/admin/AdminTopbar.tsx`
- `dashboard/src/components/admin/DataTable.tsx`
- `dashboard/src/components/admin/MetricCard.tsx`
- `dashboard/src/components/admin/StatusBadge.tsx`
- `dashboard/src/components/admin/OperationsOverview.tsx`

## Estratégia de implementação

A estratégia é **substituir todos os overrides hardcoded por herança das variáveis do Telegram** e **usar componentes shadcn consistentemente**. As variáveis `--minimal-*` e `--crm-*` serão **mapeadas** para as variáveis Telegram que o `useTheme.ts` já injeta.

---

## Passos detalhados

### Passo 1 — Eliminar blocos de variáveis hardcoded no CSS (index.css)

**Remover completamente** os 3 blocos de override de variáveis:

1. **Linhas 2752-2791**: Bloco `.admin-layout-v2` com `--crm-*` fixas → REMOVER
2. **Linhas 3565-3619**: Bloco `.admin-layout-v2` com `--minimal-*` fixas → REMOVER  
3. **Linhas 3621-3667**: Bloco `[data-theme="dark"] .admin-layout-v2` → REMOVER

**Substituir por** um bloco bridge que mapeia as variáveis internas usadas pelo CSS para as variáveis do Telegram (que o `useTheme.ts` já injeta):

```css
.admin-layout-v2 {
  /* ── Bridge: map internal tokens to Telegram theme vars ── */
  --crm-sidebar-width: 248px;
  --crm-sidebar-collapsed: 88px;

  /* Map --minimal-* to Telegram vars for CSS that still references them */
  --minimal-primary: var(--accent);
  --minimal-primary-light: var(--accent-hover);
  --minimal-primary-dark: var(--accent);
  --minimal-primary-soft: var(--accent-soft);
  --minimal-canvas: var(--bg);
  --minimal-paper: var(--card, rgba(255, 255, 255, 0.05));
  --minimal-text: var(--text);
  --minimal-muted: var(--hint);
  --minimal-line: var(--border);
  --minimal-shadow: none;
  --minimal-radius: 12px;

  /* shadcn semantic tokens — inherit from Telegram */
  --background: var(--bg);
  --foreground: var(--text);
  --card-foreground: var(--text);
  --popover: var(--card, rgba(255, 255, 255, 0.05));
  --popover-foreground: var(--text);
  --primary: var(--accent);
  --primary-foreground: var(--accent-text);
  --secondary: var(--surface);
  --secondary-foreground: var(--text);
  --muted: var(--surface);
  --muted-foreground: var(--hint);
  --input: var(--border);
  --ring: var(--accent);
  --destructive: #ef4444;

  /* Layout */
  width: 100%;
  min-height: 100dvh;
  height: 100dvh;
  background: var(--bg);
  color: var(--text);
  font-family: 'Montserrat', -apple-system, BlinkMacSystemFont, sans-serif;
}
```

### Passo 2 — Substituir todas as cores `#hex` hardcoded no CSS admin

Trocar em todos os seletores `.admin-layout-v2` e sub-seletores:

| De | Para |
|---|---|
| `#ffffff`, `#fff`, `var(--crm-white)`, `var(--minimal-paper)` | `var(--card)` |
| `#1b1b1b`, `#1b1b1a`, `#1f1f1f`, `var(--crm-ink)`, `var(--minimal-text)` | `var(--text)` |
| `#f7f7f5`, `#fafaf8`, `#fbfcfc`, `#f4f6f8`, `var(--minimal-canvas)`, `var(--crm-cream)` | `var(--bg)` |
| `#efefed`, `#f0f0ed`, `#f1f1ef`, `#f2f2ef`, `var(--crm-soft)`, `var(--minimal-primary-soft)` | `var(--accent-soft)` |
| `#e5e5e1`, `#d8d8d3`, `#dfe3e8`, `#e7e9ec`, `var(--crm-line)`, `var(--minimal-line)` | `var(--border)` |
| `#6e6e68`, `#70706d`, `#637381`, `#92928d`, `#6f6f6b`, `var(--crm-muted)`, `var(--minimal-muted)` | `var(--hint)` |
| `#090909`, `var(--minimal-primary-dark)` | `var(--accent)` |
| `#f9fcfb`, `#fbfefd` | `var(--surface-hover)` |
| `background: #1b1b1b` (botões) | `background: var(--accent)` |
| `color: #fff` (botão texto) | `color: var(--accent-text)` |

### Passo 3 — Remover overrides finais duplicados (linhas 3988-4021)

Estas 33 linhas são cópias do bloco de linhas 15-30 e contêm mais cores hardcoded. Remover e consolidar.

### Passo 4 — Trocar ícones no AdminSidebar.tsx

```diff
 import {
-  LayoutDashboard, Users, Hash, MessageSquare, FileClock,
-  Settings, Zap, Crown, Star, Smartphone, ChevronLeft,
+  LayoutDashboard, Users, Radio, Megaphone, ScrollText,
+  Settings, ShieldCheck, Sparkles, CreditCard, KeyRound, ChevronLeft,
   ChevronRight
 } from 'lucide-react';
```

Atualizar o array `SIDEBAR_ITEMS`:

| Item | De | Para |
|---|---|---|
| Canais | `<Hash>` | `<Radio>` |
| Broadcast | `<MessageSquare>` | `<Megaphone>` |
| Auditoria | `<Zap>` | `<ShieldCheck>` |
| Logs | `<FileClock>` | `<ScrollText>` |
| Contas MTProto | `<Smartphone>` | `<KeyRound>` |
| Features | `<Crown>` | `<Sparkles>` |
| Assinaturas | `<Star>` | `<CreditCard>` |

Trocar `<button>` por shadcn `<Button variant="ghost">` nos itens de navegação.
Aumentar icon `size` de `18` para `20`.

### Passo 5 — Trocar componentes nativos no AdminTopbar.tsx

- `<input>` de busca → `<Input>` do shadcn
- `<select>` de sort/filter → `<Select>` do shadcn (com `SelectTrigger`, `SelectContent`, `SelectItem`)
- `<button>` de ações → `<Button>` do shadcn

### Passo 6 — Trocar componentes nativos no DataTable.tsx

- `<input>` de busca → `<Input>` do shadcn
- `<button>` de paginação → `<Button variant="outline">` do shadcn

### Passo 7 — Usar shadcn Card no MetricCard.tsx

- Wrapping em `<Card>` + `<CardContent>` do shadcn

### Passo 8 — Usar shadcn Badge no StatusBadge.tsx

- Trocar `<span className="status-badge ...">` por `<Badge variant={...}>` do shadcn

---

## Riscos
- **Regressão visual**: mitigado por trocar apenas tokens de cor e manter o layout intacto
- **Conflito shadcn/tailwind**: os componentes shadcn usam `--background`, `--foreground` etc. que agora herdarão corretamente do Telegram
- **CSS que referencia --minimal-\***: o bridge garante compatibilidade — tudo é resolvido para variáveis Telegram

## Impactos esperados
- Dashboard admin com tema dinâmico do Telegram (qualquer tema custom do Telegram será refletido)
- Zero cores fixas — tudo herdado via `useTheme.ts`
- Componentes shadcn padronizados
- Ícones semanticamente corretos

## Compatibilidade
- Linux ✅
- macOS ✅  
- Windows ✅
- Docker ✅
- CI/CD ✅

## Como testar

### Build
```bash
cd dashboard && npm run build
```

### Verificação visual
- DevTools → trocar `--bg`, `--accent` etc. na `:root` e confirmar que a dashboard muda
- Telegram WebApp → tema real aplicado

## Rollback
```bash
git checkout dashboard/src/index.css dashboard/src/components/admin/AdminSidebar.tsx dashboard/src/components/admin/AdminTopbar.tsx dashboard/src/components/admin/DataTable.tsx dashboard/src/components/admin/MetricCard.tsx dashboard/src/components/admin/StatusBadge.tsx
```

## Observações
- O `useTheme.ts` já está correto — injeta todas as variáveis Telegram na `:root`
- O componente shadcn `Card` usa `bg-card` e `text-card-foreground`, que agora herdarão do Telegram
- A fonte Montserrat é mantida no font-family do `.admin-layout-v2`
- Componentes fora da admin dashboard (dashboard do usuário) NÃO são afetados
