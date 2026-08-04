# Plano: Redesenhar todas as abas do admin com Minimal UI

## Pedido do usuário
Substituir o visual atual do CRM administrativo por uma interface inspirada no Minimal UI e redesenhar todas as abas de administração para que sejam coerentes entre si.

## Objetivo
Transformar `/admin/dash` em um painel administrativo unificado, com a linguagem visual do Minimal UI: fundo neutro, superfícies brancas elevadas, navegação vertical ampla, header compacto, tipografia legível, cor semântica consistente e densidade adequada para operações.

O redesenho abrangerá overview, usuários, canais, broadcast, auditoria, logs, configurações, contas MTProto, features premium e assinaturas. As APIs, permissões, fluxos, textos críticos, confirmações e dados atuais serão preservados.

## Contexto atual
- O projeto usa React 19, Vite, Tailwind v4, shadcn/Base UI e Lucide; não usa Material UI.
- O admin está integrado em `/admin/dash`, protegido pela autenticação atual e composto por `AdminLayout`, `AdminSidebar`, `AdminTopbar` e `AdminDashboard`.
- O redesenho anterior seguiu a referência BizLink: paleta creme fixa, hero com gauge e pipeline. O usuário não aprovou essa direção como resultado final.
- O admin possui dez módulos funcionais: Visão Geral, Usuários, Canais, Broadcast, Auditoria, Logs, Contas MTProto, Features, Assinaturas e Configurações.
- Atualmente, cada módulo usa uma combinação distinta de cards, bordas, espaçamentos, badges, ícones de estado e estilos Tailwind. Broadcast e Configurações, por exemplo, funcionam, mas não pertencem visualmente à mesma família.
- A renderização de Configurações mostrou corretamente os campos e seus estados, mas também exibiu o toast de erro de API no ambiente de mock. Broadcast apresentou formulário, seleção de público, editor, preview e ações funcionando como uma composição visualmente diferente do overview.
- O Minimal UI oficial é um kit React de dashboard construído sobre MUI. Como este projeto já usa shadcn/Base UI e não possui MUI, a implementação adotará a linguagem visual e os princípios de layout sem adicionar MUI, dependências, assets proprietários ou código do kit.
- A documentação oficial do Minimal UI mostra uma paleta primária verde (`#00A76F`), header desktop de 72 px, navegação vertical de 300 px e conteúdo com espaçamento amplo. Essas medidas serão referências, não uma cópia literal.
- O worktree está sujo com alterações anteriores de autenticação e do primeiro CRM. Todo merge deverá ser cirúrgico; não será usado reset, checkout amplo ou limpeza de arquivos alheios.

## Evidências analisadas

### Sistema visual existente
- `dashboard/src/index.css` contém tokens globais e um bloco CRM escopado a `.admin-layout-v2`.
- `dashboard/src/components/ui/` fornece Button, Card, Input, Select, Switch, Badge e Dialog reutilizáveis.
- `@fontsource-variable/geist` já está instalado; não é necessário adicionar fonte ou biblioteca nova.
- `AdminLayout` mantém navegação, URL e drawer mobile; `AdminCrmContext` mantém busca, filtros e ordenação.

### Módulos e estados
- **Visão geral:** hero, gráfico, pipeline, filtro, busca, ordenação, vazio e overflow.
- **Usuários/Canais:** tabelas com detalhes e ações reais.
- **Broadcast:** URL de mídia, editor, seleção de público, IDs específicos, botões inline, preview, validação e confirmação de envio.
- **Auditoria:** disparo de auditoria, carregamento, resultados, remoção em massa e estado vazio.
- **Logs:** filtros, carregamento, vazio, lista expansível e paginação.
- **Configurações:** carregamento, erro, toggles, editores, cache do PostBuilder e salvamento.
- **MTProto:** carregamento, wizard de autenticação, erro, contas conectadas, toggle e remoção.
- **Features premium:** carregamento, preço, toggle, salvar e erro.
- **Assinaturas:** métricas, filtros, busca, seleção, cancelamento, reembolso, confirmação, carregamento e vazio.

### Referência Minimal UI
- O [Minimal UI](https://docs.minimals.cc/) é um dashboard React baseado em MUI, mas fornece princípios de tema, layout e navegação que podem ser reproduzidos sem incorporar sua implementação.
- A documentação define o verde `#00A76F` como cor primária e cita 72 px para header desktop, 300 px para navegação vertical e espaçamento amplo do conteúdo. [Cores](https://docs.minimals.cc/colors/), [layout](https://docs.minimals.cc/layout/).

## Arquivos analisados
- `AGENTS.md`
- `.agent/context.md`
- `.agent/memory/memory.md`
- `.agent/decisions.md`
- `.agent/skills/fable-method/SKILL.md`
- `.agent/skills/fable-method/references/domains/design-ux.md`
- `.agent/plans/done/corrigir-fidelidade-visual-crm_2026-08-03_00-56.md`
- `dashboard/package.json`
- `dashboard/src/index.css`
- `dashboard/src/App.tsx`
- `dashboard/src/components/AdminDashboard.tsx`
- `dashboard/src/components/AdminNoticeTab.tsx`
- `dashboard/src/components/AdminConfigTab.tsx`
- `dashboard/src/components/AdminAuditTab.tsx`
- `dashboard/src/components/AdminLogsTab.tsx`
- `dashboard/src/components/AdminMTProtoAccountsTab.tsx`
- `dashboard/src/components/AdminPremiumFeaturesTab.tsx`
- `dashboard/src/components/AdminSubscriptionsTab.tsx`
- `dashboard/src/components/admin/AdminLayout.tsx`
- `dashboard/src/components/admin/AdminSidebar.tsx`
- `dashboard/src/components/admin/AdminTopbar.tsx`
- `dashboard/src/components/admin/DataTable.tsx`
- `dashboard/src/components/ui/button.tsx`
- `dashboard/src/components/ui/card.tsx`
- `dashboard/src/components/ui/input.tsx`
- `dashboard/src/components/ui/select.tsx`
- `dashboard/src/components/ui/switch.tsx`
- `dashboard/src/components/ui/dialog.tsx`

## Arquivos que poderão ser modificados
- `dashboard/src/index.css` — substituir o bloco visual CRM por tokens Minimal UI escopados ao admin e estilos de estado reutilizáveis.
- `dashboard/src/components/admin/AdminLayout.tsx` — shell, comportamento responsivo e áreas de conteúdo.
- `dashboard/src/components/admin/AdminSidebar.tsx` — navegação, agrupamento, destaque e rodapé.
- `dashboard/src/components/admin/AdminTopbar.tsx` — header, busca contextual, perfil e ações primárias.
- `dashboard/src/components/AdminDashboard.tsx` — títulos, breadcrumbs, toolbars e composição de overview/tabelas.
- `dashboard/src/components/admin/CrmOverviewHero.tsx`
- `dashboard/src/components/admin/CustomerPipeline.tsx`
- `dashboard/src/components/admin/CustomerCard.tsx`
- `dashboard/src/components/admin/DataTable.tsx`
- `dashboard/src/components/AdminNoticeTab.tsx`
- `dashboard/src/components/AdminConfigTab.tsx`
- `dashboard/src/components/AdminAuditTab.tsx`
- `dashboard/src/components/AdminLogsTab.tsx`
- `dashboard/src/components/AdminMTProtoAccountsTab.tsx`
- `dashboard/src/components/AdminPremiumFeaturesTab.tsx`
- `dashboard/src/components/AdminSubscriptionsTab.tsx`
- `dashboard/src/components/ui/button.tsx`, `card.tsx`, `input.tsx`, `select.tsx`, `switch.tsx`, `dialog.tsx` — somente para adicionar variantes e tokens necessários que possam ser reutilizados sem afetar a Mini App.
- `.agent/memory/memory.md`
- `.agent/decisions.md`

## Arquivos que poderão ser criados
- `dashboard/src/components/admin/AdminPageHeader.tsx` — cabeçalho compartilhado com breadcrumb, título, descrição e ações.
- `dashboard/src/components/admin/AdminMetricCard.tsx` — métrica Minimal UI para overview, assinaturas e operações.
- `dashboard/src/components/admin/AdminEmptyState.tsx` — estado vazio consistente, sem placeholders.
- `dashboard/src/components/admin/AdminSectionCard.tsx` — superfície administrativa compartilhada para formulários, resultados e configurações.
- `dashboard/src/components/admin/adminUiTokens.ts` — constantes semânticas de cores/estados, se necessário para não duplicar cores inline.

## Estratégia de implementação
Criar uma camada de design system exclusivamente em `.admin-layout-v2`, preservando o tema da Mini App. A direção escolhida será “Minimal UI operacional”: fundo `#F9FAFB`, superfícies brancas, borda cinza discreta, sombra baixa, raio de 12–16 px, tipografia Geist, destaque verde `#00A76F` e cores semânticas consistentes para sucesso, alerta, erro e informação.

O shell terá navegação vertical de aproximadamente 300 px em desktop, header de 72 px e área de conteúdo com largura máxima legível. As abas deixarão de criar cards ad-hoc: todas passarão a usar cabeçalho, toolbar, superfície, estado vazio, loading, action bar e confirmação no mesmo padrão.

Não será instalado o pacote Minimal UI, Material UI, templates, imagens ou Figma do fornecedor. Essa abordagem evita incompatibilidade de componentes, dependência nova e uso indevido de material licenciado, mantendo apenas uma inspiração visual declarada.

## Passos detalhados

1. Preservar o estado atual e preparar o recorte.
   - Registrar `git status` e os arquivos modificados antes do redesenho.
   - Não tocar em `App.tsx`, `api.ts` ou controladores Go além do estritamente necessário para uma ligação visual já existente.
   - Confirmar que não há dependência MUI a introduzir.

2. Definir os tokens Minimal UI do admin.
   - Criar tokens locais para background, paper, texto, texto secundário, bordas, sombra, raio, altura de header, largura de navegação e escalas de spacing.
   - Definir primário verde `#00A76F`, sucesso, alerta, erro e informação em variantes sólidas e suaves.
   - Remover o creme/gauge/pipeline como linguagem dominante do admin e manter a paleta da Mini App inalterada.
   - Garantir contraste de texto e ação primária.

3. Redesenhar o shell global.
   - Sidebar vertical de 300 px no desktop, com logo, grupos de navegação, item ativo em superfície verde suave e rodapé de administrador.
   - Header de 72 px, background translúcido discreto, breadcrumb/contexto à esquerda e busca/ações contextuais à direita.
   - Drawer no mobile, overlay, foco, tecla Escape quando já suportada e continuidade de navegação por URL.
   - Manter visíveis e navegáveis as dez áreas administrativas.

4. Criar primitivas reutilizáveis do admin.
   - `AdminPageHeader`: breadcrumb, título, descrição e área de ação.
   - `AdminSectionCard`: título, descrição, ação secundária, conteúdo e rodapé.
   - `AdminMetricCard`: ícone semântico, número, legenda e tendência/estado opcional.
   - `AdminEmptyState`: ilustração apenas por ícone Lucide, mensagem, explicação e ação quando existente.
   - Padronizar skeleton, erro inline, tooltips e confirmação destrutiva.

5. Refazer a Visão Geral.
   - Trocar o hero BizLink pelo resumo operacional Minimal UI com métricas verdadeiras, gráfico compacto e cards de atividade/estado.
   - Converter o pipeline para uma lista/grade de usuários e estágios com superfícies elevadas, badges e ações reais, sem dados comerciais inventados.
   - Preservar busca, filtro, ordenação, clique no usuário, estados vazio e overflow.

6. Harmonizar Usuários e Canais.
   - Aplicar page header, toolbar, tabela em card elevado, cabeçalho sticky quando aplicável, badges semânticos, paginação e empty state compartilhado.
   - Preservar ações de abrir detalhes, abrir canal e as operações atuais do usuário.

7. Redesenhar Broadcast.
   - Organizar composição em duas colunas responsivas: editor de mensagem e painel sticky de preview/entrega.
   - Exibir seleção de público como grupos de opções consistentes e os botões inline em seção repetível.
   - Manter contador, validação, estado desabilitado, preview Telegram e modal de confirmação; nenhuma mensagem será enviada durante a validação.

8. Redesenhar Auditoria e Logs.
   - Auditoria: painel de ação com explicação, resultados em lista agrupada e empty state de êxito.
   - Logs: barra de filtros em card, resultados em timeline/lista densa com chips de status, detalhe expansível e paginação consistente.
   - Preservar carregamento, falha de API, resultado vazio e ações de abrir canal/remover resultado.

9. Redesenhar Configurações e MTProto.
   - Configurações: seções Sistema, Legendas e PostBuilder em `AdminSectionCard`, com toggles alinhados, avisos claros e action bar de salvar fixa em telas longas.
   - MTProto: wizard de conexão dentro de um card de progresso, lista de contas em superfícies consistentes e confirmação de remoção claramente destrutiva.
   - Preservar atualização de cache, auto-refresh existente, erros, loading e todos os campos/API atuais.

10. Redesenhar Features e Assinaturas.
    - Features: catálogo de recursos em cards com preço editável, estado ativo/inativo e toggle acessível.
    - Assinaturas: métricas no topo, filtros e busca em toolbar, seleção em massa destacada, cards/tabela de assinatura com status e confirmações de cancelar/reembolsar.
    - Preservar reembolso, cancelamento, cópia de charge ID, estados de seleção e avisos irreversíveis.

11. Remover inconsistências e estilos locais.
    - Substituir cores inline dispersas por tokens semânticos do admin quando elas fizerem parte do layout novo.
    - Evitar sobrescrever a Mini App e componentes de usuário com regras globais.
    - Corrigir apenas rótulos evidentemente errados encontrados durante a revisão visual, sem mudar comportamento.

12. Verificar por observação e interação.
    - Renderizar overview, broadcast, configurações, logs, MTProto e assinaturas em desktop e mobile.
    - Verificar hover, foco, loading, erro, vazio, overflow e confirmação em cada módulo que os possui.
    - Testar busca, filtros, paginação, seleção, toggles, navegação, preview, modais e ações que não produzam efeito externo.
    - Não disparar broadcast, reembolso, cancelamento ou exclusão real durante a validação.

13. Rodar build, revisão adversarial e documentação.
    - Executar Vite build e TypeScript, separando falhas preexistentes de novas falhas.
    - Executar `fable-judge` depois da implementação, comparando alterações ao plano e reproduzindo o fluxo renderizado.
    - Atualizar memória/decisões, registrar resultados e mover o plano para `done` somente quando as evidências existirem.

## Riscos
- **Escopo alto:** são dez módulos e muitos estados; a implementação será organizada por primitives compartilhadas para impedir estilos divergentes.
- **Minimal UI não é dependência atual:** importar MUI ou copiar templates criaria incompatibilidade e risco de licença. Será usada somente a direção visual documentada.
- **Regressão funcional:** Broadcast, reembolso, cancelamento, MTProto e configurações têm efeitos operacionais; seus handlers e confirmações não serão reescritos.
- **Dados reais indisponíveis em desenvolvimento:** algumas abas exibem loading/erro com mock. Layouts serão verificados nesses estados e com dados existentes onde possível; a API autenticada precisa de sessão real para validação total.
- **CSS legado:** `index.css` contém estilos antigos e o bloco CRM anterior. O bloco administrativo deverá ser consolidado em vez de receber mais sobreposições conflitantes.
- **Worktree sujo:** mudanças preexistentes devem ser preservadas e nunca revertidas em bloco.
- **TypeScript preexistente:** há erros fora do escopo administrativo; a conclusão reportará explicitamente qualquer falha não resolvida.

## Impactos esperados
- Um painel administrativo único e coerente, em vez de telas visualmente independentes.
- Menos duplicação de cards, headings, vazios, barras de ação e estilos de formulário.
- Melhor leitura de operações densas como logs, assinaturas e auditoria.
- Responsividade e acessibilidade mais previsíveis em todas as abas.
- Nenhuma migration, endpoint, mudança de autenticação ou dependência nova.

## Compatibilidade
- Linux, macOS e Windows: frontend Vite/React continua suportado.
- Docker e CI/CD: sem dependência nova, alteração de build ou servidor.
- Chrome, Firefox, Safari e Telegram WebView: layout responsivo com fallback para drawer no mobile.
- Tema global: a Mini App continua independente; o admin adotará seu tema visual próprio.

## Como testar

### Build
```bash
cd dashboard && /home/gabriel/.local/share/mise/installs/node/26.5.1/bin/node ./node_modules/vite/bin/vite.js build
```

### Testes
```bash
cd dashboard && /home/gabriel/.local/share/mise/installs/node/26.5.1/bin/node ./node_modules/typescript/bin/tsc --noEmit
```

### Execução
```bash
cd dashboard && /home/gabriel/.local/share/mise/installs/node/26.5.1/bin/node ./node_modules/vite/bin/vite.js --host 127.0.0.1
```

Abrir e verificar:

```txt
/admin/dash
/admin/dash?tab=users
/admin/dash?tab=channels
/admin/dash?tab=notice
/admin/dash?tab=audit
/admin/dash?tab=logs
/admin/dash?tab=accounts
/admin/dash?tab=premium-features
/admin/dash?tab=subscriptions
/admin/dash?tab=config
```

### Cenários obrigatórios
- Desktop em 1440×1024 e 1024×768; mobile em 390×844.
- Navegação entre as dez abas e sincronização de URL.
- Estado loading, erro, vazio e overflow onde o módulo os suporta.
- Foco de teclado visível, labels de botões de ícone e contraste dos estados.
- Busca/filtro/ordenação em overview, usuários, canais, logs e assinaturas.
- Modal sem confirmar ações externas; Broadcast permanece em preview, cancelamento/reembolso/remoção permanecem apenas em confirmação.

## Rollback
Reverter exclusivamente os hunks do redesenho Minimal UI em componentes e estilos administrativos, preservando alterações anteriores de autenticação e CRM. Como não haverá migration, API nova, dependência ou alteração de dados, o rollback será apenas de frontend.

Não usar `git reset --hard`, `git checkout --` amplo, remoção de diretório ou qualquer alteração destrutiva.

## Observações
- O plano recomenda uma direção Minimal UI compatível com o stack atual; não instalará nem redistribuirá o produto Minimal UI/MUI.
- A implementação visual deverá ser comprovada por screenshots e interações; build isolado não será suficiente.
- O objetivo é redesenhar todas as abas em uma mesma família, não apenas trocar cores no overview.

## Resultado da implementação
- Criado sistema visual Minimal UI escopado a `.admin-layout-v2`, sem dependências novas.
- Sidebar desktop ajustada para 300 px, header para 72 px e canvas administrativo para `#f7f8fa`.
- Aplicados tokens de paper, borda, sombra, raio, semântica de estados e ação primária acessível `#007867`.
- Criados `AdminPageHeader`, `AdminMetricCard` e `AdminEmptyState` e aplicados ao shell/overview.
- Todas as abas receberam o mesmo header e superfícies: Usuários, Canais, Broadcast, Auditoria, Logs, Configurações, MTProto, Features e Assinaturas.
- Broadcast, Configurações, logs, auditoria, MTProto, features e assinaturas preservaram seus componentes, handlers e confirmações; a mudança foi de apresentação e composição.

## Evidências de verificação
- Vite build aprovado após a implementação final: 2026 módulos transformados.
- `git diff --check` aprovado.
- Renderização revisada em 1440×1024 para overview, Broadcast e Configurações e em 390×844 para overview mobile.
- Auditoria automatizada local percorreu as dez URLs administrativas: todas renderizaram `.admin-layout-v2`; as nove telas secundárias apresentaram o `AdminPageHeader` esperado.
- Interações não destrutivas verificadas: busca, filtro de blacklist, busca vazia, quatro estados vazios, navegação para Configurações, 41 controles de configuração e drawer mobile com overlay.
- Contraste calculado: `#007867` sobre branco = 5,41:1; texto principal `#1c252e` sobre branco = 15,52:1; texto secundário `#637381` sobre branco = 4,88:1.
- `tsc --noEmit` continua falhando em problemas preexistentes fora do redesenho; não reportou os componentes Minimal UI novos.

## Veredito fable-judge
**VERIFIED WITH CAVEATS**

- Build, diff, renderização e fluxos locais não destrutivos foram reproduzidos.
- Nenhum teste foi enfraquecido, nenhuma dependência foi adicionada e nenhum efeito externo foi executado.
- A API autenticada real não foi exercitada: as telas foram verificadas com o mock do Vite e com estados de erro/loading locais.
- O TypeScript geral da base continua pendente por erros anteriores ao trabalho.
