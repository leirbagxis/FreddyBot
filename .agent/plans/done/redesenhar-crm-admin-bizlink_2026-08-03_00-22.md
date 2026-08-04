# Plano: Redesenhar CRM Admin no estilo BizLink

## Pedido do usuário
Realizar uma análise técnica das três imagens anexadas e criar um CRM administrativo visualmente fiel à referência, integrado às configurações e operações administrativas que já existem no FreddyBot.

## Objetivo
Redesenhar o painel existente em `/admin/dash` com a linguagem visual das referências BizLink — tipografia sans geométrica, paleta monocromática quente, sidebar fixa, barra de busca, faixa de indicadores e quadro em colunas — preservando autenticação, permissões, APIs e todas as funções administrativas atuais.

O resultado será um único dashboard integrado ao projeto React/Vite existente. Não será criado um segundo frontend `admin-crm/`.

## Contexto atual
- O dashboard é React 19 + TypeScript + Vite + Tailwind v4 + shadcn, em `dashboard/`.
- O admin já funciona na rota `/admin/dash` e usa o mesmo bundle servido pelo binário Go.
- A autenticação atual é cookie-only, com `AuthMiddlewareJWT` e autorização de papel `admin`/`owner` no backend.
- A API administrativa já oferece overview, usuários, canais, broadcast, configurações, auditoria, logs, contas MTProto, features premium e assinaturas.
- O admin já possui layout com sidebar, topbar e componentes compartilhados, mas sua linguagem visual atual usa indigo, transparências e cards genéricos; ela não corresponde à referência monocromática.
- A busca do topbar atual mantém estado local, mas não filtra nenhuma tela. Notificação e perfil também aparecem como controles sem fluxo associado.
- O endpoint `/api/admin/overview` entrega usuários e canais completos, incluindo `created_at`, flags de admin/blacklist e canais por usuário. Esses dados permitem indicadores e segmentação operacional sem inventar receita ou negócios.
- Não existe no domínio uma entidade de CRM/deal, estágio persistente ou API de drag-and-drop. Portanto, o quadro será inicialmente informativo e derivado de dados reais, não um funil comercial editável.
- Há alterações locais do usuário em `dashboard/src/App.tsx`, `dashboard/src/api.ts` e `internal/api/controllers/authController.go`; qualquer implementação deverá preservá-las e fazer merge cirúrgico.
- Já existe um plano pendente anterior, `crm-admin-freddybot_2026-08-03_00-13.md`, que propõe um frontend independente. Este plano novo o substitui conceitualmente, sem apagar nem sobrescrever o histórico.

## Análise técnica das imagens

### Direção visual observada
- Tipografia declarada na referência: **General Sans**, com peso regular no corpo e semibold/bold em títulos e números.
- Paleta declarada na própria imagem:
  - `#1f1f1f`: texto principal, botão primário e card selecionado.
  - `#f6f7ed`: painel de indicadores em tom creme.
  - `#f4f4f4`: superfícies secundárias, item ativo e cards alternativos.
  - `#ffffff`: fundo principal e cards neutros.
- Bordas finas em cinza quente, raios pequenos/médios e sombras praticamente ausentes.
- Hierarquia baseada em espaço em branco, contraste tipográfico e blocos planos, não em gradientes ou glassmorphism.

### Estrutura observada
- Sidebar fixa e branca, dividida em navegação principal, grupos secundários e identidade do usuário no rodapé.
- Topbar horizontal com busca à esquerda e ordenação, filtros, perfil e ação primária à direita.
- Faixa de analytics em creme, dividida em gráfico de barras, indicador semicircular e KPIs numéricos.
- Quadro principal em quatro colunas com contagem, ordenação e cards compactos.
- Card selecionado em fundo `#1f1f1f`, texto branco e metadados adicionais.
- Densidade alta, mas legível: textos pequenos, espaçamento consistente e metadados reduzidos.

### Adaptação semântica ao FreddyBot
- `Dashboard` → **Visão Geral**.
- `Customers` → **Usuários**.
- `Projects/Tasks/Activity` → módulos reais de **Canais**, **Broadcast**, **Auditoria** e **Logs**.
- A seção inferior da sidebar será usada para **MTProto**, **Features**, **Assinaturas** e **Configurações**.
- O botão “Add customer” não será copiado, pois usuários são criados pelo fluxo do Telegram. Ele será substituído por **Enviar broadcast**, uma ação real já suportada.
- O indicador “Successful deals” será **Taxa de ativação**: usuários com pelo menos um canal dividido pelo total de usuários.
- O gráfico “New customers” será **Novos usuários**, agrupado pelos últimos cinco dias a partir de `created_at`.
- Os KPIs monetários da referência não serão inventados. A faixa exibirá métricas suportadas: total de usuários, canais, administradores e bloqueados.
- O quadro em colunas terá estágios mutuamente exclusivos e calculados nesta prioridade:
  1. **Atenção**: usuário em blacklist.
  2. **Ativos**: usuário com um ou mais canais e fora da blacklist.
  3. **Novos**: usuário sem canais, fora da blacklist e criado nos últimos sete dias.
  4. **Em ativação**: demais usuários sem canais e fora da blacklist.
- O quadro será somente leitura na primeira versão. Clicar em um card abrirá o detalhe/ações do usuário já existentes. Drag-and-drop exigiria modelo, migration e endpoints novos e não faz parte deste redesenho visual.

## Critérios de pronto
- `/admin/dash` renderiza o novo layout CRM usando a paleta e composição das referências no tema claro.
- Todas as dez áreas administrativas atuais continuam acessíveis e funcionais.
- Busca, ordenação e filtros exibidos no CRM têm comportamento real.
- Indicadores e colunas usam apenas dados existentes e regras documentadas.
- Loading, erro, vazio, hover, focus, seleção e overflow possuem estados visuais definidos.
- A interface é observada em desktop, tablet e mobile, com screenshots de comparação.
- TypeScript e build do Vite passam, sem regressão no fluxo de autenticação ou nas rotas administrativas.

## Arquivos analisados
- `AGENTS.md`
- `.agent/context.md`
- `.agent/memory/memory.md`
- `.agent/decisions.md`
- `.agent/plans/pending/crm-admin-freddybot_2026-08-03_00-13.md`
- `.agent/plans/done/redesenhar-admin-dashboard_2026-07-05_12-30.md`
- `.agent/plans/done/redesign-config-tab-theme-tabs_2026-07-24_23-52.md`
- `dashboard/package.json`
- `dashboard/tsconfig.json`
- `dashboard/vite.config.ts`
- `dashboard/src/App.tsx`
- `dashboard/src/api.ts`
- `dashboard/src/types.ts`
- `dashboard/src/mockData.ts`
- `dashboard/src/index.css`
- `dashboard/src/components/AdminDashboard.tsx`
- `dashboard/src/components/AdminConfigTab.tsx`
- `dashboard/src/components/admin/AdminLayout.tsx`
- `dashboard/src/components/admin/AdminSidebar.tsx`
- `dashboard/src/components/admin/AdminTopbar.tsx`
- `dashboard/src/components/admin/AdminOverview.tsx`
- `dashboard/src/components/admin/MetricCard.tsx`
- `dashboard/src/components/admin/DataTable.tsx`
- `dashboard/src/components/admin/SystemHealth.tsx`
- `internal/api/routes/routes.go`
- `internal/api/controllers/adminController/getAllUserAdminController.go`
- `internal/api/controllers/adminController/configController.go`
- `internal/api/controllers/adminController/channelEventsController.go`
- `internal/api/controllers/adminController/adminSubscriptionController.go`
- `internal/api/controllers/adminController/premiumFeaturesController.go`
- `internal/api/controllers/adminController/adminAccountController.go`
- `internal/core/services/user.go`
- `internal/database/repositories/user.go`
- `internal/database/models/models.go`
- `internal/api/dto/dto.go`
- `internal/api/dto/mapper.go`
- `Makefile`
- `Dockerfile`
- `nginx/default.conf`
- `original-67e6a497c042e3e5028704313e0738c1.webp`
- `original-5618f2daa44a40756d821e94d38b6c67.webp`
- `original-68641804ff1b3e20d9f1fb6d0a521878.webp`

## Arquivos que poderão ser modificados
- `dashboard/src/index.css` — tokens e estilos CRM escopados ao admin.
- `dashboard/src/App.tsx` — sincronização de aba/URL e ligação das ações, preservando as mudanças locais de autenticação.
- `dashboard/src/types.ts` — tipos derivados do CRM e correções de contrato necessárias.
- `dashboard/src/mockData.ts` — mock administrativo coerente com `User` (`first_name`, flags e datas) para validar os estados em desenvolvimento.
- `dashboard/src/components/AdminDashboard.tsx` — composição das abas e ligação com overview/quadro.
- `dashboard/src/components/admin/AdminLayout.tsx` — shell responsivo.
- `dashboard/src/components/admin/AdminSidebar.tsx` — navegação no estilo da referência.
- `dashboard/src/components/admin/AdminTopbar.tsx` — busca, filtros, ordenação, perfil e ação real.
- `dashboard/src/components/admin/MetricCard.tsx` — variante monocromática/compacta.
- `dashboard/src/components/admin/DataTable.tsx` — toolbar e tabela coerentes com o CRM.
- Componentes administrativos existentes em `dashboard/src/components/Admin*Tab.tsx`, somente quando necessário para aplicar o shell visual e estados compartilhados sem alterar lógica.

## Arquivos que poderão ser criados
- `dashboard/src/components/admin/CrmOverviewHero.tsx` — gráfico de novos usuários, gauge de ativação e KPIs.
- `dashboard/src/components/admin/CustomerPipeline.tsx` — colunas, filtros, ordenação e overflow.
- `dashboard/src/components/admin/CustomerCard.tsx` — card de usuário neutro/selecionado.
- `dashboard/src/components/admin/crmSelectors.ts` — funções puras para métricas, segmentação, busca e ordenação.
- `dashboard/src/components/admin/crmSelectors.test.ts` — testes unitários se o projeto adotar um runner disponível; caso contrário, as funções serão cobertas por validação TypeScript e cenários manuais documentados.

## Estratégia de implementação
Aplicar o redesenho dentro do dashboard atual, com tokens CSS escopados em `.admin-layout-v2` para não alterar a Mini App dos usuários. A composição seguirá a referência no tema claro e manterá equivalentes semânticos nos temas escuro e Telegram.

O overview usará os dados já carregados por `/api/admin/overview`; as regras de métricas e pipeline ficarão em funções puras, evitando chamadas adicionais e lógica visual espalhada. As demais abas continuarão consumindo seus endpoints atuais.

General Sans não existe no repositório. Para não introduzir download ou licença não verificada, a implementação usará primeiro `Geist Variable`, já instalado e visualmente próximo, com a pilha `"General Sans", "Geist Variable", sans-serif`. Fidelidade tipográfica pixel-perfect dependerá do fornecimento/validação dos arquivos General Sans.

## Passos detalhados

1. Preservar o estado atual do worktree.
   - Revalidar `git status` e o diff de `App.tsx`, `api.ts` e `authController.go`.
   - Não sobrescrever as correções cookie-only/fail-closed existentes.

2. Consolidar o design system administrativo.
   - Criar tokens escopados para cores, bordas, raios, larguras, espaçamentos e tipografia.
   - Remover do admin claro o aspecto indigo/glass sem alterar os temas da Mini App.
   - Resolver regras duplicadas de `.admin-sidebar` pela ordem/cascade, sem limpeza ampla de CSS fora do escopo.

3. Refazer o shell do admin.
   - Sidebar fixa com grupos Principal, Operações, Premium e Sistema.
   - Estado ativo em `#f4f4f4`, sem barra indigo.
   - Identidade do administrador no rodapé e colapso acessível.
   - Mobile com drawer e overlay; desktop com largura fixa e conteúdo fluido.

4. Tornar o topbar funcional.
   - Busca compartilhada com o overview/usuários.
   - Ordenação por recentes, nome e quantidade de canais.
   - Filtros por admin, blacklist, com canais e sem canais.
   - Ação primária “Enviar broadcast” navegando para a aba já existente.
   - Remover ou desabilitar visualmente controles sem função real até que possuam fluxo.
   - Sincronizar a aba ativa com `?tab=` via History API, evitando reload completo.

5. Implementar o hero de analytics.
   - Barras de novos usuários nos cinco dias mais recentes.
   - Gauge SVG/CSS da taxa de ativação.
   - KPIs de usuários, canais, admins e blacklist.
   - Tratar divisão por zero, datas inválidas e listas vazias.

6. Implementar o pipeline operacional.
   - Aplicar a regra de prioridade Atenção → Ativos → Novos → Em ativação.
   - Exibir contagem por coluna e cards com nome, username/ID, canais, data de entrada e badges.
   - Card selecionado usa fundo `#1f1f1f`; clique abre o detalhe real do usuário.
   - Colunas vazias exibem estado informativo; grandes listas usam limite inicial e “ver todos”, sem renderizar milhares de cards de uma vez.
   - Sem drag-and-drop ou estágio persistido nesta fase.

7. Harmonizar as telas existentes.
   - Manter Users e Channels em tabela para gestão precisa.
   - Aplicar mesma densidade, bordas, inputs, botões, estados e títulos às abas Broadcast, Auditoria, Logs, Configurações, MTProto, Features e Assinaturas.
   - Preservar toggles, saves, modais, toasts, paginação e chamadas de API atuais.

8. Corrigir o ambiente de demonstração local.
   - Tipar `mockAdminData` como `AdminDashboardData`.
   - Trocar campos inconsistentes como `firstName` por `first_name` e preencher flags obrigatórias.
   - Incluir casos para novo, ativo, blacklist, admin, vazio e overflow.

9. Verificar acessibilidade e responsividade.
   - Focus visible, navegação por teclado, `aria-label` em botões de ícone e ordem semântica de headings.
   - Conferir contraste dos quatro tons e dos estados de perigo/sucesso.
   - Respeitar `prefers-reduced-motion`.
   - Validar 1440×1024, 1024×768 e 390×844.

10. Executar build e revisão visual.
    - Rodar TypeScript/Vite.
    - Abrir `/admin/dash` com mock e com API real quando disponível.
    - Capturar screenshots dos três tamanhos e comparar hierarquia, alinhamento, densidade e estados com as referências.
    - Percorrer todas as abas e ações críticas antes de declarar concluído.

11. Executar `fable-judge` após a implementação.
    - Confirmar que o trabalho foi renderizado, que os estados existem e que nenhum controle falso ou dado inventado foi entregue.
    - Só então mover este plano para `.agent/plans/done/` e registrar decisões/memória relevantes.

## Riscos
- **Fidelidade de fonte:** General Sans não está disponível localmente; Geist será o fallback até haver ativo/licença confirmados.
- **Referências apenas desktop:** comportamento tablet/mobile precisa ser inferido e validado por observação.
- **Sem entidade de CRM:** o quadro é uma segmentação operacional calculada, não um funil comercial persistente.
- **Escalabilidade:** `/api/admin/overview` carrega todos os usuários e canais. O limite visual reduz o custo de render, mas paginação/aggregates no backend será necessária se o volume crescer muito.
- **Worktree sujo:** `App.tsx` e `api.ts` já possuem alterações do usuário e exigem merge cuidadoso.
- **CSS legado:** há seletores administrativos duplicados em `index.css`; mudanças globais podem afetar a Mini App se não forem escopadas.
- **Mock inconsistente:** o mock atual usa campos incompatíveis com o tipo `User`, o que pode mascarar ou gerar erro de build.
- **Ambiente atual:** `npm`/`node` não estão disponíveis nesta sessão, portanto o build não foi observado durante a análise; a implementação não poderá ser considerada pronta sem esse comando em um ambiente Node funcional.
- **Temas:** a correspondência exata com a referência vale para o tema claro; escuro e Telegram serão adaptações semânticas.

## Impactos esperados
- Um único admin CRM integrado, sem duplicar autenticação, dependências, deploy ou manutenção.
- Todas as configurações administrativas permanecem no mesmo fluxo e com as mesmas APIs.
- Melhor leitura operacional de aquisição/ativação de usuários.
- Busca, filtros e ordenação deixam de ser apenas decorativos.
- Nenhuma migration ou mudança de regra de negócio é necessária nesta fase.
- A Mini App de usuários permanece visualmente isolada do redesenho por meio do escopo CSS.

## Compatibilidade
- Linux: suportado; desenvolvimento requer Node compatível com o Dockerfile (Node 24) e Go 1.25 para build integral.
- macOS: suportado por Vite e navegador moderno.
- Windows: suportado pelo frontend; comandos de Makefile podem exigir WSL/Git Bash.
- Docker: preservado pelo build multi-stage atual.
- CI/CD: sem nova etapa; continua usando `npm ci`, `npm run build` e `go build`.
- Navegadores: Chrome/Edge/Firefox/Safari modernos; Telegram WebView continua usando o dashboard de usuário existente.

## Como testar

### Build
```bash
cd dashboard && npm run build
```

```bash
make build
```

### Testes
```bash
cd dashboard && npx tsc --noEmit
```

Se for adicionado um runner de testes já aprovado/disponível:

```bash
cd dashboard && npm test -- crmSelectors
```

### Execução
```bash
cd dashboard && npm run dev -- --host 0.0.0.0
```

Abrir:

```txt
http://localhost:5173/admin/dash
```

### Cenários funcionais
- Admin/owner autenticado acessa todas as abas.
- Usuário sem papel admin recebe acesso negado pela API.
- Busca por nome, username e ID.
- Filtros e ordenação combinados.
- Usuário novo, ativo, admin, blacklist e sem canais aparecem na coluna correta.
- Clique no card abre detalhe; ações admin/blacklist/broadcast continuam funcionando.
- Configuração carrega, salva e informa erro sem perder valores.
- Listas vazias, falha de API, loading, conteúdo longo e milhares de registros não quebram o layout.

## Rollback
Reverter somente os arquivos listados neste plano, preservando as alterações preexistentes do usuário. Componentes novos podem ser removidos após retirar seus imports. Como não há migration nem novo serviço, o rollback não altera banco, Redis, autenticação ou API.

Não usar `git reset --hard`, `git checkout --` amplo ou exclusão de diretório; o worktree já contém mudanças não relacionadas.

## Observações
- Este plano substitui a recomendação de criar `admin-crm/`, mas não apaga o plano anterior por razões de rastreabilidade.
- “Idêntico” será tratado como fidelidade de composição, paleta, tipografia, densidade e interação, com conteúdo adaptado ao domínio real do FreddyBot.
- Não serão exibidos receita, pagamentos, “deals” ou tarefas fictícias só para copiar números da referência.
- Drag-and-drop persistente pode ser planejado depois como funcionalidade separada, com modelo de estágio, migration, serviço, endpoints e auditoria.

## Resultado da implementação
- CRM integrado implementado em `/admin/dash`, sem criar projeto frontend separado.
- Hero com novos usuários, gauge de ativação e quatro KPIs alimentados por dados reais.
- Pipeline operacional com Novos, Em ativação, Ativos e Atenção.
- Busca, filtros, ordenação, atalho `Ctrl/⌘+K`, broadcast e URL das abas validados por interação no Firefox.
- Sidebar e topbar redesenhados; colisão com CSS legado de largura foi identificada e corrigida.
- Mock administrativo tipado e cobrindo os quatro estágios.
- Layout observado em 1440×1024 e 390×844, incluindo menu mobile, estado vazio e erro de API.
- Contrastes principais medidos entre 5,13:1 e 16,48:1.
- `vite build` passou após a implementação.
- `tsc --noEmit` permanece vermelho por erros preexistentes em componentes fora do CRM; nenhum erro novo foi apontado nos arquivos adicionados.
- `fable-judge`: **VERIFIED WITH CAVEATS** — build e interações reproduzidos; API real autenticada e General Sans exata não estavam disponíveis no ambiente.
