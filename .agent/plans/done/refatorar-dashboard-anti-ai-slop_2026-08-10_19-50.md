# Plano: Refatoração da Interface da Dashboard (Diretrizes Anti-AI-Slop)

## Pedido do usuário
Aplicar integralmente as diretrizes do prompt "Anti-AI-Slop UI Design Prompt" em toda a dashboard do projeto FreddyBot. A interface deve ser limpa, contida, altamente funcional, com tipografia forte, superfícies contínuas e sem elementos decorativos genéricos de SaaS (sem sombras em conteúdo estático, sem gradientes decorativos, sem glassmorphism, sem ícones em quadradinhos coloridos, sem subtítulos óbvios desnecessários e sem cards dentro de cards).

## Objetivo
Refatorar a camada visual e estrutural do React/Tailwind na Dashboard (`dashboard/src`), transformando-a em uma ferramenta operacional de alta densidade, minimalista e focada em utilidade diária (estilo Linear, GitHub, Vercel).

## Contexto atual
- A Dashboard possui estilos globais em `dashboard/src/index.css` que misturam efeitos de sombra (`box-shadow`), gradientes, efeito vidro (`backdrop-filter: blur`), ícones em caixas coloridas com cantos arredondados, cantos `rounded-xl`/`16px` excessivos e subtítulos descritivos em quase todos os cabeçalhos.
- Os componentes de abas e admin (`OperationsOverview.tsx`, `AdminOverview.tsx`, `MetricCard.tsx`, `DashboardInicioTab.tsx`, `CaptionCard.tsx`, `ReactionsCard.tsx`, `CustomCaptionsCard.tsx`, etc.) dependem de grids de cards flutuantes isolados com fundos decorativos.

## Arquivos analisados
- `dashboard/src/index.css`
- `dashboard/src/components/admin/AdminLayout.tsx`
- `dashboard/src/components/admin/AdminSidebar.tsx`
- `dashboard/src/components/admin/AdminTopbar.tsx`
- `dashboard/src/components/admin/AdminOverview.tsx`
- `dashboard/src/components/admin/OperationsOverview.tsx`
- `dashboard/src/components/admin/MetricCard.tsx`
- `dashboard/src/components/admin/StatusBadge.tsx`
- `dashboard/src/components/admin/DataTable.tsx`
- `dashboard/src/components/DashboardInicioTab.tsx`
- `dashboard/src/components/CaptionCard.tsx`
- `dashboard/src/components/CustomCaptionsCard.tsx`
- `dashboard/src/components/ReactionsCard.tsx`
- `dashboard/src/components/ButtonGrid.tsx`
- `dashboard/src/components/TabBar.tsx`
- `dashboard/src/components/ui/card.tsx`
- `dashboard/src/components/ui/button.tsx`
- `dashboard/src/components/ui/input.tsx`

## Arquivos que poderão ser modificados
- `dashboard/src/index.css`
- `dashboard/src/components/admin/AdminOverview.tsx`
- `dashboard/src/components/admin/OperationsOverview.tsx`
- `dashboard/src/components/admin/MetricCard.tsx`
- `dashboard/src/components/admin/StatusBadge.tsx`
- `dashboard/src/components/admin/AdminSidebar.tsx`
- `dashboard/src/components/admin/AdminTopbar.tsx`
- `dashboard/src/components/DashboardInicioTab.tsx`
- `dashboard/src/components/CaptionCard.tsx`
- `dashboard/src/components/CustomCaptionsCard.tsx`
- `dashboard/src/components/ReactionsCard.tsx`
- `dashboard/src/components/ButtonGrid.tsx`
- `dashboard/src/components/TabBar.tsx`
- `dashboard/src/components/ui/card.tsx`

## Estratégia de implementação

1. **Design System & CSS (`index.css`)**:
   - Remover gradientes, efeitos glow, glassmorphic blurs (`backdrop-filter`) e sombras decorativas de conteúdo.
   - Definir superfícies planas neutras com linhas de separação discretas de 1px.
   - Reduzir raio de curvatura excessivo (`rounded-xl` / `16px` -> `rounded-md` / `6px` / `8px`).
   - Eliminar os envoltórios de ícones coloridos (`.section-icon`, `.action-card-icon`, `.metric-card-icon`). Os ícones passam a ser inline ou limpos na cor neutra/secundária.

2. **Métricas e Resumos Operacionais (`MetricCard.tsx`, `OperationsOverview.tsx`, `AdminOverview.tsx`)**:
   - Agrupar métricas em uma linha contínua unificada com divisórias simples em vez de cards soltos e flutuantes.
   - Remover gráficos de fundo tipo sparkline decorativos e subtítulos desnecessários sob cada métrica.
   - Aumentar densidade de informação: números alinhados em tabelas e listas limpas.

3. **Visões de Configuração e Canal (`CaptionCard.tsx`, `ReactionsCard.tsx`, `CustomCaptionsCard.tsx`, `ButtonGrid.tsx`, `DashboardInicioTab.tsx`)**:
   - Remover subtítulos redundantes sob cada título.
   - Substituir "cards flutuantes" por seções planas com bordas de divisão de seção limpas.
   - Manter controles (botões, inputs, switches) compactos e alinhados horizontalmente.

4. **Navegação e Layout (`TabBar.tsx`, `AdminSidebar.tsx`, `AdminTopbar.tsx`)**:
   - Simplificar a barra de navegação inferior e lateral: remoção do efeito glassmorphic blur flutuante exagerado e bordas brilhantes.
   - Alinhamento preciso com tipografia contida e foco nos estados ativo/hover.

5. **Passo de Auditoria e Limpeza Obliterativa**:
   - Revisar visualmente os seletores para garantir que nenhuma sombra ou container decorativo sem propósito permaneça.

## Passos detalhados

1. **Atualizar `dashboard/src/index.css`**:
   - Ajustar tokens de tema (`--radius: 6px`, superfícies planas sem blur/glow).
   - Remover regras CSS decorativas (`.action-card-icon`, `.section-icon`, sombras excessivas).

2. **Refatorar Componentes de UI Base (`dashboard/src/components/ui/card.tsx`, `MetricCard.tsx`, `StatusBadge.tsx`)**:
   - Tornar o Card uma superfície plana minimalista sem sombras padrão.
   - Refatorar `MetricCard.tsx` para exibição densa e direta de valores.

3. **Refatorar Componentes Admin (`OperationsOverview.tsx`, `AdminOverview.tsx`, `AdminSidebar.tsx`, `AdminTopbar.tsx`)**:
   - Eliminar subtítulos desnecessários.
   - Converter blocos dispersos em seções integradas com divisórias sutis.

4. **Refatorar Componentes de Legenda e Canais (`DashboardInicioTab.tsx`, `CaptionCard.tsx`, `CustomCaptionsCard.tsx`, `ReactionsCard.tsx`, `ButtonGrid.tsx`)**:
   - Limpar títulos, subtítulos e envoltórios de ícones decorativos.
   - Garantir alta densidade funcional com formulários limpos.

5. **Validar o Build do Frontend**:
   - Executar `npm run build` na pasta `dashboard/` para assegurar compilação do TypeScript e Vite sem avisos ou erros.

## Riscos
- **Baixo**: Alterações visuais e estruturais no CSS/JSX sem modificar contratos de dados ou APIs do backend.

## Impactos esperados
- A Dashboard passará a ter uma estética sóbria, utilitária e focada em dados (estilo Vercel/Linear/GitHub).
- Maior velocidade de leitura, melhor densidade de informação e eliminação completa de elementos visuais "AI-slop".

## Compatibilidade
- Linux
- macOS
- Windows
- Docker
- CI/CD

## Como testar

### Build Frontend
```bash
cd dashboard && npm run build
```

### Execução em Desenvolvimento
```bash
cd dashboard && npm run dev
```

## Rollback
Em caso de divergência visual, reverter as alterações via git:
```bash
git checkout dashboard/
```

## Observações
O plano respeita todas as restrições impostas pelas regras universais de desenvolvimento.
