# Plano: Aplicar Regras de Design na Dashboard Admin

## Pedido do usuário
"Agora vamos aplicar todas as regras de designer na dashboard admin. Fonte, div com fundo branco leve e sem o contorno da div, thema de acordo com o telegram, cantos arredondados e os icones correspondentes"

## Objetivo
Unificar o sistema de design da Dashboard Admin (`/admin/dash`) com o padrão visual do resto da aplicação:
1. **Fonte**: Garantir o uso da fonte **Montserrat** em toda a Dashboard Admin.
2. **Cards e Divs**: Aplicar o fundo translúcido suave (`bg-white/5` / `rgba(255, 255, 255, 0.05)`) em todos os cards, tabelas, containers e seções da admin, removendo todas as bordas e contornos (`border: none` / `border-color: transparent`).
3. **Cantos Arredondados**: Padronizar cantos arredondados (`rounded-2xl` / `16px`) em cards, tabelas, inputs e botões.
4. **Tema de acordo com o Telegram**: Remover estilos/cores legados hardcoded (como `#1b1b1b`, `#fafaf8`, `#f7f7f5`, `#deded9`) em `.admin-layout-v2`, integrando a Admin completamente com as variáveis dinâmicas do tema Telegram.
5. **Ícones**: Manter ícones Lucide consistentes e limpos em caixas suaves (`section-icon`).

## Arquivos analisados
- `dashboard/src/index.css` (regras `.admin-layout-v2`, `.data-table-wrapper`, `.admin-notice-page`, etc.)
- `dashboard/src/components/AdminDashboard.tsx`
- `dashboard/src/components/admin/DataTable.tsx`
- `dashboard/src/components/admin/AdminTopbar.tsx`
- `dashboard/src/components/admin/AdminSidebar.tsx`

## Arquivos que poderão ser modificados
- `dashboard/src/index.css`
- `dashboard/src/components/admin/DataTable.tsx`

## Estratégia de implementação

1. **`dashboard/src/index.css`**:
   - Remover as variáveis CSS legadas hardcoded de `.admin-layout-v2` (`--crm-white`, `--crm-cream`, `--crm-soft`, `--crm-line`, `'General Sans'`).
   - Configurar `.admin-layout-v2` para usar `font-family: 'Montserrat', sans-serif`, `background: var(--bg)`, e `color: var(--text)`.
   - Atualizar todos os elementos cards/divs da admin (`[data-slot='card']`, `.data-table-wrapper`, `.admin-notice-page > div`, `.admin-audit-page [data-slot='card']`, `.admin-logs-page > div`, `.admin-accounts-page .cfg-account-card`, `.admin-features-page [data-slot='card']`, `.admin-subscriptions-page [data-slot='card']`, `.admin-config-page [data-slot='card']`):
     - `background: rgba(255, 255, 255, 0.05)`
     - `border: none`
     - `border-radius: 16px`
     - `box-shadow: none`
   - Atualizar a tabela de dados (`.data-table-wrapper`, `.data-table th`, `.data-table td`, `.data-table-toolbar`):
     - `border: none`
     - `background: transparent` ou `rgba(255, 255, 255, 0.03)`
     - `border-bottom: 1px solid rgba(255, 255, 255, 0.05)` nas linhas
     - Hover nas linhas: `background: rgba(255, 255, 255, 0.08)`
   - Remover os blocos de overrides legados em `index.css` que forçavam bordas cinzas e fundos brancos opacos na admin.

2. **`dashboard/src/components/admin/DataTable.tsx`**:
   - Remover bordas hardcoded e ajustar classes para seguir cantos arredondados `rounded-2xl` e o tema transparente/translúcido.

3. **Validação**:
   - Executar `npm run build` na pasta `dashboard`.

## Riscos
- **Baixo**: Apenas ajustes visuais de CSS e classes sem alteração na lógica de negócio da admin.

## Impactos esperados
- A Dashboard Admin terá a mesma estética moderna, translúcida (`bg-white/5`), sem bordas, com fonte Montserrat e cantos arredondados, perfeitamente sintonizada com o tema do Telegram.

## Compatibilidade
- Linux, macOS, Windows, Docker, CI/CD

## Como testar
```bash
cd dashboard && npm run build
```

## Rollback
```bash
git checkout dashboard/src/index.css dashboard/src/components/admin/DataTable.tsx
```
