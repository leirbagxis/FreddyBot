# Plano: Redesenhar Dashboard Admin (Estilo TASA Explorer)

## Pedido do usuário
Recriar a dashboard admin usando o mesmo layout/design das 4 imagens fornecidas (TASA Explorer design system), mantendo as funcionalidades e integrações de dados existentes.

## Objetivo
Aplicar fielmente a identidade visual do design system "TASA Explorer" na Dashboard Admin (`AdminDashboard.tsx`, `MetricCard.tsx`, `AdminOverview.tsx`, `AdminSidebar.tsx`, `AdminTopbar.tsx`, etc.), incorporando:
- **Paleta de Cores TASA**:
  - Background principal: `#F4F4F4` (Clean light mode) / Dark mode equivalente
  - Cards e containers: `#FFFFFF` com sombras suaves e bordas de 16px a 24px (`rounded-2xl` / `rounded-3xl`)
  - Texto principal e botões primários: `#050505`
  - Acento vibrante: `#FA5B2E` (Laranja quente) para abas ativas, botões de ação e destaques
- **Controles de Período em Pílula**:
  - Abas em formato de cápsula arredondada (`Today`, `Week`, `Month`, `Year`), onde o item selecionado ganha fundo `#FA5B2E` e texto branco.
- **Cards de Métricas TASA Explorer**:
  - Títulos em tipografia limpa (ex: Workload, Hours, Employees, Budget, Overall sales)
  - Valores numéricos em grande destaque em negrito
  - Porcentagens de variação com setas coloridas (`↗` verde para positivo/crescimento, `↘` vermelho para decréscimo)
- **Gráfico de Estatísticas / Vendas Gerais**:
  - Gráfico de onda/linha fluida com preenchimento em gradiente suave laranja (`#FA5B2E`).
- **Listas de Atividades & Usuários**:
  - Cards brancos arredondados, ícones estilizados com formas coloridas, avatares agrupados (`+4`, `+2`), badges de porcentagem em pílulas com setas de tendência.
- **Barra de Navegação Estilo Pílula Flutuante**:
  - Menu de navegação moderno em formato de pílula preta com ícones limpos (Home, Usuários, Estatísticas, Menu) e destaque laranja no item ativo.

## Contexto atual
A dashboard admin atual possui componentes em `dashboard/src/components/AdminDashboard.tsx`, `dashboard/src/components/admin/` (`MetricCard.tsx`, `AdminOverview.tsx`, `AdminSidebar.tsx`, etc.) que utilizam o tema base shadcn/tailwind padrão.

## Arquivos analisados
- `/home/malbs/Opencode/FreddyBot/dashboard/src/components/AdminDashboard.tsx`
- `/home/malbs/Opencode/FreddyBot/dashboard/src/components/admin/MetricCard.tsx`
- `/home/malbs/Opencode/FreddyBot/dashboard/src/components/admin/AdminOverview.tsx`
- `/home/malbs/Opencode/FreddyBot/dashboard/src/components/admin/AdminSidebar.tsx`
- `/home/malbs/Opencode/FreddyBot/dashboard/src/components/admin/AdminTopbar.tsx`
- `/home/malbs/Opencode/FreddyBot/dashboard/src/index.css`

## Arquivos que poderão ser modificados
- `dashboard/src/components/AdminDashboard.tsx`
- `dashboard/src/components/admin/MetricCard.tsx`
- `dashboard/src/components/admin/AdminOverview.tsx`
- `dashboard/src/components/admin/AdminSidebar.tsx`
- `dashboard/src/components/admin/AdminTopbar.tsx`
- `dashboard/src/index.css`

## Estratégia de implementação
1. **Design System & Tokens (CSS)**: Adicionar utilitários e variáveis do TASA Explorer em `index.css` (`#F4F4F4`, `#FFFFFF`, `#050505`, `#FA5B2E`).
2. **Componente de Filtro por Período**: Implementar a barra de pílulas `[Today | Week | Month | Year]` com estados ativos em laranja.
3. **Redesenho dos Cards de Métricas (`MetricCard.tsx`)**: Atualizar o visual para corresponder exatamente às imagens do TASA Explorer (Workload, Hours, Employees, Budget).
4. **Gráficos & Estatísticas**: Adicionar/estilizar gráfico de linha com gradiente laranja suave para dados numéricos da dashboard.
5. **Redesenho de Tabelas e Atividades**: Estilizar a lista de usuários e atividades recentes com ícones coloridos, avatares agrupados e badges em pílulas com porcentagens `44% ↗`.
6. **Barra Flutuante de Navegação (Pill Navbar)**: Criar o menu inferior/superior no formato de cápsula com os ícones principais.

## Passos detalhados
1. Criar e registrar este plano de ação em `.agent/plans/pending/redesenhar-dashboard-admin-tasa-explorer_2026-07-28_13-24.md`.
2. Apresentar o plano detalhado ao usuário e solicitar autorização explícita para prosseguir.
3. Após o aceite, mover o plano para `.agent/plans/approved/`.
4. Atualizar o CSS (`index.css`) com as variáveis e classes do TASA Explorer.
5. Implementar a reformulação visual de `MetricCard.tsx`, `AdminOverview.tsx`, `AdminDashboard.tsx` e componentes auxiliares da admin.
6. Executar o build do Vite (`npm run build`) para verificar a integridade da aplicação.
7. Mover o plano para `.agent/plans/done/`.

## Riscos
- Alterações visuais afetassem a responsividade ou a exibição de dados reais da API.
- Mitigação: Manter intacta a passagem de dados via props e handlers, aplicando as mudanças estritamente na camada visual (JSX / Tailwind / CSS).

## Impactos esperados
- Dashboard Admin com design moderno e premium, perfeitamente fiel às imagens do TASA Explorer.
- Preservação de todas as rotas e funções do painel admin.

## Compatibilidade
- Linux
- Navegadores modernos (Chrome, Firefox, Safari)

## Como testar

### Build
```bash
cd dashboard && npm run build
```

## Rollback
```bash
git checkout dashboard/
```
