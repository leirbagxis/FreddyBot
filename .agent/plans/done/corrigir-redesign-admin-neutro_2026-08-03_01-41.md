# Plano: corrigir redesign admin neutro

## Pedido do usuário

Corrigir o redesenho anterior do painel administrativo: remover a cor verde, tornar a interface realmente minimalista e redesenhar de fato todas as abas do admin, em vez de apenas aplicar uma camada de CSS sobre as telas existentes.

## Objetivo

Entregar um painel administrativo claro, neutro e operacional, com uma linguagem visual única em todas as dez abas. O resultado não deve usar a composição CRM/BizLink, o texto “CRM”, verde como cor primária, bordas pontilhadas, nem cards decorativos em excesso. As funções, rotas e chamadas de API existentes serão preservadas.

## Contexto atual

- O dashboard administrativo está em `dashboard/`, acessível pela rota `/admin/dash`.
- A interface está encapsulada em `.admin-layout-v2`, portanto a correção continuará sem interferir no dashboard dos usuários.
- O redesenho anterior introduziu o tema `Minimal UI` com verde `#007867`, navegação de 300 px, bordas pontilhadas e uma visão geral baseada em CRM/funil. Essa direção foi rejeitada pelo usuário.
- Existem dez superfícies administrativas ativas: visão geral, usuários, canais, broadcast, auditoria, logs, contas MTProto, features, assinaturas e configurações.
- As ações sensíveis existentes — enviar broadcast, remover canais, remover conta MTProto, alterar permissões e salvar configuração — devem manter confirmações e comportamento atuais. O trabalho é de estrutura e apresentação, não de mudança de regra de negócio.

## Arquivos analisados

- `.agent/context.md`
- `.agent/memory/memory.md`
- `docs/superpowers/specs/2026-04-28-admin-refactoring-design.md`
- `dashboard/package.json`
- `dashboard/src/App.tsx`
- `dashboard/src/index.css`
- `dashboard/src/components/AdminDashboard.tsx`
- `dashboard/src/components/AdminNoticeTab.tsx`
- `dashboard/src/components/AdminAuditTab.tsx`
- `dashboard/src/components/AdminLogsTab.tsx`
- `dashboard/src/components/AdminConfigTab.tsx`
- `dashboard/src/components/AdminMTProtoAccountsTab.tsx`
- `dashboard/src/components/AdminPremiumFeaturesTab.tsx`
- `dashboard/src/components/AdminSubscriptionsTab.tsx`
- `dashboard/src/components/admin/AdminLayout.tsx`
- `dashboard/src/components/admin/AdminSidebar.tsx`
- `dashboard/src/components/admin/AdminTopbar.tsx`
- `dashboard/src/components/admin/AdminPageHeader.tsx`
- `dashboard/src/components/admin/AdminMetricCard.tsx`
- `dashboard/src/components/admin/AdminEmptyState.tsx`
- `dashboard/src/components/admin/CrmOverviewHero.tsx`
- `dashboard/src/components/admin/CustomerPipeline.tsx`
- Imagens de referência fornecidas pelo usuário.

## Arquivos que poderão ser modificados

- `dashboard/src/index.css`
- `dashboard/src/components/admin/AdminLayout.tsx`
- `dashboard/src/components/admin/AdminSidebar.tsx`
- `dashboard/src/components/admin/AdminTopbar.tsx`
- `dashboard/src/components/admin/AdminPageHeader.tsx`
- `dashboard/src/components/admin/AdminMetricCard.tsx`
- `dashboard/src/components/admin/AdminEmptyState.tsx`
- `dashboard/src/components/admin/CrmOverviewHero.tsx` ou um substituto neutro
- `dashboard/src/components/admin/CustomerPipeline.tsx` ou um substituto neutro
- `dashboard/src/components/AdminDashboard.tsx`
- `dashboard/src/components/AdminNoticeTab.tsx`
- `dashboard/src/components/AdminAuditTab.tsx`
- `dashboard/src/components/AdminLogsTab.tsx`
- `dashboard/src/components/AdminConfigTab.tsx`
- `dashboard/src/components/AdminMTProtoAccountsTab.tsx`
- `dashboard/src/components/AdminPremiumFeaturesTab.tsx`
- `dashboard/src/components/AdminSubscriptionsTab.tsx`
- Eventuais novos componentes compartilhados pequenos em `dashboard/src/components/admin/`
- `.agent/memory/memory.md`
- `.agent/decisions.md`

## Estratégia de implementação

Adotar uma direção de **minimalismo editorial neutro**, recomendada para corrigir a discrepância apontada:

- Fundo off-white muito sutil, painéis brancos, texto grafite e divisores sólidos cinza-claro.
- Preto/grafite como única cor de ação primária; cores semânticas aparecerão somente em estados que precisam comunicar risco, sucesso ou aviso.
- Sem verde como identidade, gradientes, sombras pesadas, bordas pontilhadas ou grandes blocos de “métricas decorativas”.
- Barra lateral mais compacta e silenciosa, com item ativo em cinza claro e tipografia discreta.
- Cada aba terá uma composição própria para sua tarefa operacional; os componentes compartilhados restringirão consistência a cabeçalho, toolbar, painel, linhas e estados vazios, sem transformar todas as páginas em cópias umas das outras.

Alternativas descartadas:

- Recolorir somente o tema atual: não resolve a ausência de redesenho estrutural.
- Reproduzir novamente a referência BizLink: ela é útil como referência de clareza, mas não representa o admin minimalista solicitado nesta etapa.
- Adotar um kit visual pronto: repetiria o erro de priorizar a aparência de uma biblioteca em vez do direcionamento explícito do usuário.

## Passos detalhados

1. Substituir os tokens verdes e a regra visual `Minimal UI` dentro de `.admin-layout-v2` por tokens neutros centralizados: canvas, papel, grafite, texto secundário, divisor, foco, ação primária e estados semânticos discretos. Remover também o separador pontilhado, a etiqueta `ADMIN` verde e o contexto “FreddyBot CRM”.
2. Reconstruir o shell administrativo: sidebar compacta e estável, topbar com contexto da área atual e ações somente quando relevantes, cabeçalho de página com título, descrição curta e ações. Preservar menu responsivo, `?tab=`, busca, ordenação, filtros e atalho de teclado já existentes.
3. Substituir os primitivos genéricos atuais por painéis neutros, barras de ferramenta, linhas de dados, indicadores compactos e estados de carregamento/vazio/erro coerentes. Eles serão usados onde repetição realmente ajuda, sem impor o mesmo layout a todas as telas.
4. Redesenhar a **Visão geral** de um funil CRM para uma central operacional: resumo compacto de base, atividade/alertas e uma lista de atenção ou atalhos de operação baseados nos dados disponíveis. Remover pipeline de clientes, cartões de “deals” e nomenclatura de CRM.
5. Redesenhar **Usuários** e **Canais** como áreas de dados: toolbar contextual, tabela/linhas densas e legíveis, detalhes laterais ou em painel e ações conservadas. Manter os filtros compartilhados e os fluxos de admin/blacklist/mensagem/canal existentes.
6. Redesenhar **Broadcast** como uma estação de composição: formulário organizado por etapas lógicas, seleção de público sem decoração excessiva, editor preservado e uma prévia fixa em desktop. O disparo e a confirmação continuarão usando o fluxo atual.
7. Redesenhar **Auditoria** como uma tela de diagnóstico: cabeçalho de ação com contexto de risco, resultado em lista operacional por usuário/canal e estados de execução, ausência de resultados e remoção confirmada claramente separados.
8. Redesenhar **Logs** como uma área de investigação: filtros reunidos em uma barra única, contagem/paginação compactas, linhas de evento expansíveis e detalhe técnico legível. Manter filtros, busca, atualização, paginação e navegação ao canal.
9. Redesenhar **Contas MTProto**, **Features** e **Assinaturas** para as tarefas específicas de cada domínio: fluxo guiado de conexão e lista de contas; catálogo/controle de features em linhas; painel financeiro com resumo compacto, filtros e itens de assinatura. Preservar todos os handlers e confirmações existentes.
10. Redesenhar **Configurações** em seções de sistema, conteúdo e PostBuilder com cabeçalhos, controles e ações de salvamento claramente agrupados, sem modificar o payload nem o refresh já implementados.
11. Conferir que nenhum token visual novo vaze para as telas de usuário, que não exista rótulo ou composição CRM no admin e que os estilos antigos não permaneçam prevalecendo sobre a nova estrutura.
12. Construir, renderizar e inspecionar o admin em desktop e mobile; percorrer as dez abas, o menu móvel, busca/filtros, um detalhe expandido, carregamento, vazio, erro quando reproduzível localmente e foco por teclado. Atualizar a memória e registrar a decisão visual depois da implementação aprovada.

## Riscos

- O trabalho sobrepõe arquivos que já possuem alterações locais. Elas serão preservadas e cada edição será limitada à estrutura visual e aos componentes administrativos.
- Algumas abas dependem de respostas autenticadas da API; estados de erro e vazio serão validados por seus caminhos locais sem enviar broadcast, excluir canais, excluir contas ou fazer alterações de produção.
- A quantidade de componentes existentes pode introduzir cascatas de CSS. Os novos estilos permanecerão escopados a `.admin-layout-v2` e a classes novas, com remoção cuidadosa das regras verdes conflitantes.
- A referência visual é uma direção, não um design system completo. A proposta neutra acima será a fonte de decisão; se aprovada, ela substitui explicitamente a direção verde anterior.

## Impactos esperados

- Administração visualmente independente do dashboard comum, porém coerente com o produto.
- Navegação e dados mais fáceis de escanear, com menos ruído visual e menos cards desnecessários.
- Todas as abas passam a ter estrutura intencional para sua operação, não apenas uma nova paleta.
- Nenhuma mudança em endpoints, modelos, permissões ou comportamento do backend.

## Compatibilidade

- Linux
- macOS
- Windows
- Docker
- CI/CD
- Desktop e WebView/mobile do Telegram, com menu responsivo preservado

## Como testar

### Build

```bash
cd dashboard && npm run build
```

### Testes

```bash
cd dashboard && npx tsc --noEmit
```

O diagnóstico de TypeScript será separado entre falhas pré-existentes e qualquer regressão introduzida pelo redesenho.

### Execução

```bash
cd dashboard && npm run dev
```

Com o admin aberto em `/admin/dash`, inspecionar visualmente as dez abas em largura desktop e em 390 px, incluindo menu móvel, foco por teclado, busca/filtro, uma linha expansível, estados vazios e carregamentos. Nenhuma ação destrutiva ou de envio será confirmada durante a validação.

## Rollback

Reverter somente os arquivos listados neste plano para o estado anterior a esta implementação, preservando alterações locais não relacionadas. O plano anterior ficará intacto em `.agent/plans/done/` para rastreabilidade.

## Observações

- Critério principal de aceite: o usuário deve reconhecer que houve redesenho de verdade nas abas e que o painel não tem mais identidade verde ou aparência de CRM/BizLink.
- A implementação não adicionará dependências nem usará imagens geradas; será feita com os componentes React, ícones e estilos já presentes no repositório.
