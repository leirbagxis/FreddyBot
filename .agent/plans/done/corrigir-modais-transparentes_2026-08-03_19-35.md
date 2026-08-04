# Plano: Corrigir Modais Transparentes (Confirmar Disparo em Massa e Desconectar Bot)

## Pedido do usuário
O modal de "Confirmar Disparo em massa" e o modal de confirmação de "Desconectar o bot" estão transparentes.

## Objetivo
Garantir que todos os modais da aplicação (incluindo `ConfirmModal`, `DialogContent` do shadcn/base-ui) possuam fundo sólido, opaco e com legibilidade perfeita em todos os temas (Light, Dark, Telegram, CRM, etc.), corrigindo a transparência indesejada.

## Contexto atual
- Os modais utilizam o componente `DialogContent` (em `dashboard/src/components/ui/dialog.tsx` e `ConfirmModal.tsx`), o qual aplica a classe `bg-popover`.
- A variável CSS `--popover` estava mapeada para `var(--card-elevated)`.
- No tema Light e Telegram, `--card-elevated` está definido como `transparent`.
- No tema Dark, `--card-elevated` está definido como `rgba(20, 25, 40, 0.75)` (transparência de 75%), e no fallback `.dark` como `rgba(30, 41, 59, 0.6)`.
- Além disso, o `@base-ui/react/dialog` utiliza `DialogPortal`, que renderiza os modais diretamente no `document.body` (fora da div `.admin-layout-v2`). Regras CSS escopadas a `.admin-layout-v2 [data-slot='dialog-content']` não alcançavam os modais no portal.

## Arquivos analisados
- `dashboard/src/index.css`
- `dashboard/src/components/ConfirmModal.tsx`
- `dashboard/src/components/ui/dialog.tsx`
- `dashboard/src/components/AdminNoticeTab.tsx`
- `dashboard/src/components/DashboardInicioTab.tsx`

## Arquivos que poderão ser modificados
- `dashboard/src/index.css`
- `dashboard/src/components/ConfirmModal.tsx`

## Estratégia de implementação
1. Em `dashboard/src/index.css`:
   - Atualizar a variável `--popover` nos blocos de tema `[data-theme="light"]`, `[data-theme="dark"]`, `[data-theme="telegram"]` e `.dark` para utilizar cores sólidas de elevado/card opacas (ex: `#ffffff` no Light/Telegram e `#0f1422` no Dark) em vez de `var(--card-elevated)` transparente.
   - Tornar a regra `[data-slot='dialog-content']` global (remover o prefixo estrito `.admin-layout-v2`) para que diálogos renderizados via Portal no `body` recebam fundo sólido, borda refinada e sombra elevada independentemente de onde são disparados.
2. Em `dashboard/src/components/ConfirmModal.tsx`:
   - Garantir que o `DialogContent` aplique o fundo opaco do tema (`bg-popover`) com `shadow-2xl` e opacidade 100%.

## Passos detalhados
1. Editar `dashboard/src/index.css`:
   - No `[data-theme="light"]`, alterar `--popover: var(--card-elevated);` para `--popover: #ffffff;`.
   - No `[data-theme="dark"]`, alterar `--popover: var(--card-elevated);` para `--popover: #0f1422;`.
   - No `[data-theme="telegram"]`, alterar `--popover: var(--card-elevated);` para `--popover: #ffffff;`.
   - No `.dark` fallback, alterar `--popover: var(--card-elevated);` para `--popover: #0f172a;`.
   - Ajustar o seletor `.admin-layout-v2 [data-slot='dialog-content']` para `[data-slot='dialog-content']`, garantindo `background-color: var(--popover)` e fundo opaco em todos os modais.
2. Verificar `ConfirmModal.tsx` e demais componentes de diálogo.
3. Testar a interface no dashboard (build / lint se aplicável).

## Riscos
- Risco muito baixo. A alteração afeta apenas o fundo de popovers e modais, tornando-os devidamente sólidos e opacos conforme esperado em interfaces modernas.

## Impactos esperados
- Todos os modais (Confirmar Disparo em massa, Desconectar Bot, Excluir Botão, Assinaturas, etc.) passarão a ter fundos 100% opacos e legíveis em qualquer tema.

## Compatibilidade
- Linux
- macOS
- Windows
- Docker
- CI/CD

## Como testar

### Build
```bash
npm --prefix dashboard run build
```

### Execução
Visualizar o dashboard e abrir os modais de Disparo em massa e Desconectar Bot para confirmar que o fundo está sólido e opaco.

## Rollback
Reverter as edições em `dashboard/src/index.css` e `ConfirmModal.tsx` via `git checkout`.

## Observações
Nenhuma.
