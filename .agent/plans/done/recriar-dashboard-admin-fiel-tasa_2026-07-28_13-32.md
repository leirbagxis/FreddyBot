# Plano: Recriar Dashboard Admin (Fiel às Imagens TASA Explorer)

## Pedido do usuário
"Voce nao redeenho/refez a dashboard de acordo com as imagens anexadas" — O usuário deseja a recriação da Dashboard Admin exatamente com a mesma estrutura visual, composição de tela e componentes exibidos nas imagens anexadas (TASA Explorer design).

## Objetivo
Refazer a interface da Dashboard Admin (`AdminDashboard.tsx`, `AdminLayout.tsx`, `AdminOverview.tsx`) para replicar fielmente a composição das telas do mockup TASA Explorer:
- **Header Superior**:
  - Botão de voltar (`<`), Título centralizado (`Design` / `Statistic`), Ícone de opções (`⋮`).
- **Filtro de Período Pill Capsule**:
  - Cápsula arredondada `[ Today | Week | Month | Year ]` com o estado ativo preenchido em Laranja quentinho (`#FA5B2E`).
- **Grid de Métricas 2x2 (Cards brancos estilo TASA)**:
  - Card 1: `Workload` (porcentagem grande ex: `60%`, `+9% ↗` verde).
  - Card 2: `Hours` (valor numérico grande ex: `35`, `+12% ↗` verde).
  - Card 3: `Employees` (valor numérico grande ex: `8`, `+2% ↗` verde).
  - Card 4: `Budget` (valor em moeda ex: `$748.00`, `-10% ↘` vermelho).
- **Sessões Duplas Alternáveis ou Combinadas (Design & Statistic views)**:
  - **View Design (Próximos Eventos / Atividades)**:
    - Cabeçalho `Upcoming events` com botão `View more` em laranja.
    - Lista de cards brancos arredondados com ícone geométrico colorido à esquerda, Título e Data/Hora ao centro, e Stack de Avatares com badge numérico (`+4`, `+2`, `+5`) à direita.
  - **View Statistic (Estatísticas Gerais & Atividade)**:
    - Card `Overall sales` com valor destacado (`$3.284`), badge de tendência (`-10% ↗`), caixa de data (`09 Jun $278.00`) e o gráfico fluido de onda em tom laranja com preenchimento em gradiente.
    - Lista `Employee activity` com avatar da foto/ícone do usuário, Nome, Cargo/Função e badge de porcentagem em pílula (`44% ↗`, `21% ↘`, `37% ↗`).
- **Navegação Inferior Flutuante (Floating Pill Nav)**:
  - Barra de navegação flutuante preta cápsula (`#050505`) centralizada na parte inferior da tela com os 4 ícones principais (Home / Dashboard, Usuários / Equipe, Estatísticas, Configurações) com o item ativo destacado em círculo laranja (`#FA5B2E`).

## Contexto atual
- `dashboard/src/components/admin/AdminLayout.tsx`
- `dashboard/src/components/AdminDashboard.tsx`
- `dashboard/src/components/admin/AdminOverview.tsx`
- `dashboard/src/index.css`

## Arquivos que poderão ser modificados
- `dashboard/src/components/admin/AdminLayout.tsx`
- `dashboard/src/components/AdminDashboard.tsx`
- `dashboard/src/components/admin/AdminOverview.tsx`
- `dashboard/src/index.css`

## Estratégia de implementação
1. **Reformular `AdminOverview.tsx`**: Construir a estrutura exatamente conforme a tela `Design` e `Statistic` das imagens.
2. **Implementar a Barra Flutuante de Navegação (Pill Navbar)**: Posicionar a barra flutuante preta na base com os ícones e estados ativos idênticos aos mockups.
3. **Restruturar os Cards e Listas**:
   - Ajustar o layout do Grid 2x2 para ser exibido exatamente no estilo visual dos mockups.
   - Criar a lista de `Upcoming events` com stack de avatares (`+4`, `+2`, `+5`).
   - Criar o card `Overall sales` com a curva suavizada exata em gradiente laranja e a lista `Employee activity`.
4. **Adequar `AdminLayout.tsx`**: Permitir que a visualização fique limpa, sem barra lateral tradicional obstruindo o layout mobile/desktop inspirado na imagem.

## Passos detalhados
1. Registrar o plano em `.agent/plans/pending/recriar-dashboard-admin-fiel-tasa_2026-07-28_13-32.md`.
2. Mostrar o resumo e pedir autorização do usuário.
3. Após aprovação, mover para `.agent/plans/approved/`.
4. Atualizar os componentes `AdminOverview.tsx`, `AdminDashboard.tsx` e `AdminLayout.tsx`.
5. Compilar o painel (`npm run build`) para assegurar zero erros de execução.
6. Mover o plano para `.agent/plans/done/`.

## Riscos
- Mudar a disposição estrutural sem perder dados existentes de usuários/canais.
- Mitigação: Mapear os dados reais da API (usuários, canais, logs, estatísticas) para os novos componentes visuais do TASA Explorer.

## Impactos esperados
- A Dashboard Admin será uma réplica visual e estética perfeita das imagens fornecidas.

## Como testar
```bash
cd dashboard && npm run build
```

## Rollback
```bash
git checkout dashboard/
```
