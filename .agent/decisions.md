# Decisões Arquiteturais

## Decisão: fronteiras externas falham fechadas e registram intenção persistente

### Data
2026-08-03

### Contexto
A auditoria encontrou webhook Telegram sem autenticação, workers duplicados, pagamento sem vínculo entre invoice e compra, e sessões MTProto simuladas que pareciam válidas.

### Decisão tomada
Exigir secret token no webhook; iniciar o container e seus workers uma única vez; reivindicar postagens agendadas por transição condicional; persistir cada invoice em `payment_intents`; e rejeitar fluxos MTProto sem credenciais ou sessão real.

### Motivo
Essas fronteiras não podem conceder identidade, dinheiro ou capacidade premium por dados que não foram autenticados e persistidos de forma verificável.

### Impacto
Produção em modo webhook precisa definir `TELEGRAM_WEBHOOK_SECRET`. A migração GORM adiciona `payment_intents` e `processing_at` de modo aditivo. O compose de desenvolvimento não foi alterado.

## Decisão: proxy autenticado para foto de canal e dados públicos mínimos

### Data
2026-08-03

### Contexto
A URL de download criada pela Bot API inclui o token do bot. A rota de foto a redirecionava ao navegador e a busca de usuário retornava dados em excesso de um chat Telegram.

### Decisão tomada
A foto é baixada pelo servidor com timeout e limite de 10 MiB, entregue como binário e protegida pela autorização do canal. Busca de usuário usa apenas a base local e retorna `UserLookupDTO` com ID, nome e username.

### Motivo
Evitar exposição de credencial, enumeração de dados de terceiros e vazamento de flags administrativas, sem quebrar o carregamento da imagem pelo dashboard autenticado.

### Impacto
Clientes continuam usando `/api/channel/:channelId/photo`, agora com autorização por canal. IDs que nunca iniciaram o bot não são pesquisáveis pela API.

## Decisão: validação e recálculo explícitos de agendamentos recorrentes

### Data
2026-08-03

### Contexto
Editar apenas `scheduleTime` gravava o horário atual como próximo disparo, e valores de hora inválidos eram normalizados silenciosamente.

### Decisão tomada
Validar `HH:MM` estritamente. Edições de agendamento atualizam apenas os campos solicitados e recalculam o próximo disparo de tipos diário e semanal no serviço.

### Motivo
Impedir publicação antecipada e tornar inválidos de entrada erros visíveis ao cliente.

### Impacto
O banco guarda instantes em UTC, mas a regra de horário recorrente é calculada em `America/Sao_Paulo`.

## Decisão: identidade neutra para o painel administrativo

### Data
2026-08-03

### Contexto
O redesenho anterior adotou verde e composição de CRM/BizLink. O usuário rejeitou explicitamente essa direção e pediu minimalismo real com redesenho das abas administrativas.

### Decisão tomada
Aplicar uma identidade editorial neutra apenas em `.admin-layout-v2`: off-white, branco, grafite, divisores sólidos e ação primária preta. A visão geral passa a ser uma central de operação com fila de revisão, e as abas mantêm seus fluxos de dados sem adotar semântica de CRM.

### Motivo
Atender à direção visual explícita sem introduzir um novo frontend, dependências ou mudanças no backend.

### Impacto
O dashboard comum e a Mini App permanecem inalterados. Busca, filtros, ações sensíveis, confirmações, rotas e chamadas de API do admin são preservados.

## Decisão: Unificação do Núcleo via Core Services
### Data
2026-05-13

### Contexto
O projeto tinha lógica de negócio e acesso a dados duplicados entre os repositórios, controladores da API e handlers do Bot. Repositórios continham lógica complexa que dificultava testes e reutilização.

### Decisão tomada
Implementar uma camada de `Core Services` em `internal/core/services`. Toda a lógica de negócio e acesso a dados (GORM) deve passar por essa camada. Controladores da API e Handlers do Bot tornam-se "cascas" finas que apenas validam entrada/saída e chamam os serviços.

### Motivo
- **DRY (Don't Repeat Yourself):** Reutilização de lógica entre API e Bot.
- **Testabilidade:** Lógica isolada em serviços é mais fácil de testar.
- **Manutenibilidade:** Mudanças na regra de negócio são feitas em um único lugar.
- **Padronização:** Respostas da API via Generics `APIResponse[T]`.

### Impacto
- Removido acesso direto aos repositórios do `AppContainer`.
- Handlers do Bot e Middlewares migrados para usar Serviços.
- API refatorada para usar controladores com serviços e respostas padronizadas.

## Decisão: Mapeamento de Mensagens Inline para Sessões de Postagem
### Data
2026-05-13

### Contexto
Usuários criam postagens no Post Builder e as compartilham via modo inline. Ao votar nessas mensagens, o Telegram não fornece o teclado original, impedindo a atualização visual dos contadores de votos.

### Decisão tomada
Implementar um handler de `ChosenInlineResult` para mapear o `inline_message_id` gerado pelo Telegram para o ID da sessão da postagem no Redis.

### Motivo
Permite reconstruir o teclado original no momento do voto, possibilitando a atualização visual dos contadores enquanto mantém o `CallbackData` limpo (`vote:emoji`).

### Impacto
- Necessário ativar `Inline Feedback` no BotFather.
- Dependência de Redis para o mapeamento temporário (24h).

## Decisão: Uso de MatchFunc Customizado para ChosenInlineResult
### Data
2026-05-13

### Contexto
A biblioteca `go-telegram/bot` v1.19.0 não possui a constante `bot.HandlerTypeChosenInlineResult`, impossibilitando o uso de `RegisterHandler` padrão para esse tipo de update.

### Decisão tomada
Utilizar `RegisterHandlerMatchFunc` com uma função de match manual (`matchChosenInlineResult`) que verifica a presença do campo `ChosenInlineResult` no objeto `models.Update`.

### Motivo
Contornar a limitação da biblioteca sem a necessidade de forks ou atualização imediata da dependência, mantendo a funcionalidade de mapeamento de mensagens inline.

### Impact
- Substituição do registro padrão por match manual em `internal/telegram/events/loader.go`.
- Código permanece compatível com versões futuras caso a constante seja adicionada.

## Decisão: Conversão JIT de Markdown para HTML no PostBuilder
### Data
2026-05-15

### Contexto
O PostBuilder enfrentava problemas ao enviar mensagens via MarkdownV2 devido à rigidez do Telegram com caracteres reservados (como '.', '!', '-'), que causavam erros de "Bad Request". Além disso, formatações enviadas via interface do Telegram (entidades) eram perdidas se não processadas imediatamente.

### Decisão tomada
Padronizar o armazenamento do estado do PostBuilder em HTML. Toda entrada de texto (Título, Corpo, Rodapé) passa por `ProcessTextWithFormatting` no momento do recebimento, convertendo tanto Markdown explícito quanto Entidades do Telegram em HTML seguro.

### Motivo
- **Estabilidade:** O ParseMode HTML do Telegram é muito mais tolerante a caracteres especiais do que o MarkdownV2.
- **Fidelidade:** Permite capturar exatamente o que o usuário formatou no app (negrito/itálico via UI) e o que digitou via Markdown.
- **Simplicidade:** Evita a necessidade de rotinas complexas de escape para MarkdownV2 no lado do servidor.

### Impact
- `handleTextInput` agora salva `formattedText` (HTML).
- `InlineHandler` e `sendFinalPost` (Preview) utilizam `DetectParseMode` para garantir a integridade das tags antes do envio final.
- Melhora na experiência do usuário ao importar canais (legendas já vêm em HTML).

## Decisão: Preservação de Legendas Originais em Mídias
### Data
2026-05-16

### Contexto
Anteriormente, ao enviar uma mídia (foto/vídeo) com legenda, o bot substituía o texto do usuário pela legenda padrão do canal. O comportamento desejado é que a legenda do bot atue como um rodapé (footer), preservando o conteúdo do usuário.

### Decisão tomada
Alterar a lógica de montagem final no `StageTransform` para que, tanto em mensagens de texto quanto em mídias, o bot utilize a função `composeMessage` com a estratégia `append`.

### Motivo
- **UX:** O usuário não perde o contexto que escreveu ao enviar a mídia.
- **Consistência:** Unifica o comportamento entre tipos de mensagem (texto e mídia).
- **Flexibilidade:** Permite usar Links Dinâmicos no texto original enquanto mantém a assinatura padrão do canal abaixo.

### Impact
- `StageTransform` modificado para não sobrescrever `formattedBase` em mídias.
- Legenda padrão é adicionada com duas quebras de linha após o texto original.

## Decisão: Substituição Estrita de Legendas em Áudio
### Data
2026-05-16

### Contexto
Diferente de fotos e vídeos, onde a legenda original deve ser preservada, para arquivos de áudio/música a convenção do projeto é que a legenda original do arquivo seja totalmente descartada em favor da legenda configurada no bot.

### Decisão tomada
Implementar uma exceção no `StageTransform` para `MessageTypeAudio`. Se houver uma legenda configurada (`dbCaption`), ela substituirá completamente o texto original (`formattedBase`).

### Motivo
- **Convenção:** Manter o comportamento esperado para canais de música.
- **Limpeza:** Arquivos de áudio costumam vir com metadados ou legendas de outros bots no arquivo original que devem ser limpos.

### Impacto
- Mensagens de áudio voltam a usar a estratégia "replace".
- Demais mídias permanecem com a estratégia "append".

# Decisão

## Data
2026-05-21

## Contexto
Foi necessário criar uma postagem PostBuilder permanente, editável pela Dashboard Admin, com chave fixa e sem expiração no Redis.

## Decisão tomada
Persistir a configuração no `ServerConfig` (`fixedPostBuilderEnabled`, `fixedPostBuilderKey`, `fixedPostBuilderPayload`) e sincronizar o payload para Redis em `pb_session:<key>` sem TTL. A chave padrão definida foi `legendasbot`.

## Motivo
O banco passa a ser a fonte de verdade editável e durável, enquanto o Redis mantém compatibilidade com o fluxo inline existente `@bot pb <key>` sem expiração.

## Impacto
A Dashboard Admin pode ativar/desativar e editar a postagem fixa. O PostBuilder normal continua usando IDs aleatórios com TTL de 24h.


# Decisão

## Data
2026-05-22

## Contexto
O fluxo de New Pack ganhou dois controles para definir se o botão do pack aparece na mensagem editada pelo bot e/ou no sticker enviado. Canais antigos e entradas antigas em cache Redis não possuem esses campos.

## Decisão tomada
Representar os campos internos do modelo Go como ponteiros booleanos (`*bool`) e tratar `nil` como `true` no mapper e no handler de New Pack.

## Motivo
Preservar o comportamento anterior para bancos e caches existentes sem impedir que o usuário salve `false` explicitamente pela dashboard.

## Impacto
Canais antigos continuam exibindo botões por padrão. Quando o usuário desativa uma opção pela dashboard, o valor `false` passa a ser persistido e respeitado pelo bot.


# Decisão

## Data
2026-05-22

## Contexto
O New Pack ganhou configuração para posicionar a mensagem acima ou abaixo do sticker. Canais antigos e caches antigos não possuem esse campo.

## Decisão tomada
Persistir `newPackMessagePosition` com valores `above` e `below`, usando `above` como default quando o campo estiver vazio ou ausente.

## Motivo
Preservar o comportamento anterior para canais existentes, enquanto permite que o usuário opte pelo envio da mensagem abaixo do sticker pela dashboard.

## Impacto
Canais existentes continuam editando a mensagem de espera. Canais configurados como `below` enviam uma nova mensagem após o sticker e tentam apagar a mensagem de espera.

# Decisão

## Data
2026-05-26

## Contexto
Admins precisam auditar problemas por canal, incluindo postagens ignoradas, erros de edição, permissões ausentes e ações do PostBuilder. Logs apenas em stdout não oferecem histórico filtrável na dashboard.

## Decisão tomada
Salvar eventos estruturados na tabela `channel_events`, com campos indexáveis para filtros e `metadata` JSON textual para detalhes variáveis. A Dashboard Admin passa a consumir `/api/admin/logs` com paginação. Eventos do fluxo de canal e do PostBuilder compartilham a mesma estrutura usando `source`.

## Motivo
PostgreSQL já é a fonte operacional do projeto e atende bem o volume esperado quando os eventos são resumidos, paginados e indexados. Evita introduzir uma stack de observabilidade mais pesada antes de haver necessidade real.

## Impacto
Admins passam a consultar histórico por canal diretamente na dashboard. O logging é best-effort e não bloqueia o bot. A primeira retenção padrão remove eventos com mais de 90 dias durante a inicialização.

# Decisão: Implementação de Contas Conectadas via MTProto

## Data
2026-07-04

## Contexto
Implementar suporte a contas Telegram dos próprios usuários via protocolo MTProto (gotd/td), permitindo recursos Premium. A funcionalidade deve coexistir com a Bot API existente.

## Decisão tomada
1. Criar interface `TelegramExecutor` abstraindo BotAPI e MTProto
2. `BotAPIExecutor` encapsula chamadas telego existentes
3. `MTProtoExecutor` usa gotd/td (ephemeral clients por operação)
4. `UserExecutor` wrapper vincula userID ao executor
5. Sessões criptografadas com AES-256-GCM no PostgreSQL
6. `EditOptions.Entities` carrega JSON de entidades para rich text
7. Factory escolhe implementação por usuário
8. Dados sensíveis (código, senha) apenas em memória

## Motivo
- Modularidade total: nunca misturar MTProto nas regras de negócio
- Segurança: sessões criptografadas, dados sensíveis não persistem
- Suporte a rich text (Message Entities) via gotd/td
- Zero breaking changes: Bot API continua como fallback

## Impacto
- 6 novos endpoints REST para gerenciamento de conta
- 2 novas tabelas no banco (connected_accounts, connected_account_channels)
- 5 novos componentes React
- Pipeline refatorado para usar TelegramExecutor via ExecutorFactory
- Criptografia AES-256-GCM adicionada como dependência

# Decisão: Sistema de Legendas com Message Entities (MTProto)

## Data
2026-07-05

## Contexto
Usuários com contas conectadas precisam configurar legendas com formatação rich text (negrito, itálico, custom_emoji, blockquote, spoiler, etc.) que não são possíveis via HTML da Bot API. O sistema de legacy HTML continua para usuários sem conta conectada.

## Decisão tomada
1. Estender `DefaultCaption` com `Entities (JSON)` e `UseEntities (bool)`
2. Entities armazenados como `MessageEntityDTO` serializado (design library-agnostic)
3. `MTProtoExecutor` usa ephemeral gotd client por operação (mesmo pattern do auth)
4. `UserExecutor` resolve userID para MTProto internamente na factory
5. Post entities + caption entities combinados com offset shift UTF-16
6. `ProcessingContextTelego` ganha `FinalEntities`, `PostEntitiesJSON`, `ExecutorFactory`
7. Fallback automático: sem conta conectada → HTML (BotAPI)
8. `ConnectedAccountChannel` ganha `AccessHash` para MTProto peer resolution

## Motivo
- Legacy HTML continua intacto para usuários sem conta conectada
- Entities DTOs mantêm serialização independente da lib de dispatch
- Design ephemeral evita gerenciamento de pool de conexões MTProto
- Offset shift UTF-16 segue especificação oficial do Telegram

## Impacto
- `AccessHash` adicionado a `connected_account_channels` (campo novo, default 0)
- `MTProtoExecutor` criado com suporte a todos os tipos de entity do Telegram
- Dispatchers migrados para `ExecutorFactory.ForUser(ownerID)`
- StageTransformTelego detecta `UseEntities + HasActiveAccount`
- Build existente 100% preservado, zero breaking changes

# Decisão: Stars Test Mode com Preço de 1 Star e Sistema de Reembolso

## Data
2026-07-10

## Contexto
O modo de teste de assinaturas Stars ativava a assinatura gratuitamente sem passar pelo fluxo real de pagamento. Além disso, não havia sistema de reembolso — o admin só conseguia cancelar assinaturas, sem devolver os Stars pagos.

## Decisão tomada
1. `STARS_TEST_MODE=true` agora faz as invoices custarem 1 star (em vez de ativar grátis)
2. Removeu-se a ativação automática sem pagamento — o usuário sempre paga (1 star em teste, real em produção)
3. Criou-se tabela `refunds` com modelo `Refund` no banco
4. Implementou-se `AdminRefundPayment` que chama `bot.RefundStarPayment()` do Telegram
5. `POST /api/admin/subscriptions/refund` aceita `{userId, telegramPaymentChargeId}`
6. AdminSubscriptionsTab exibiu charge_id e botão "Reembolsar" com confirmação
7. `TelegramPaymentID` (charge_id) é sempre salvo na subscription via `HandlePayment`

## Motivo
- Testar com Stars reais (mesmo que 1) valida o fluxo completo de pagamento
- Admin pode devolver Stars sem precisar de acesso ao Telegram
- Tabela refunds previne reembolso duplicado e mantém auditoria

## Impacto
- Test mode agora custa 1 star (não mais grátis)
- Admin pode reembolsar qualquer pagamento via dashboard
- Nenhuma subscription antiga perde dados
- `bot.RefundStarPayment()` disponível na telego v1.9.0

# Decisão: CRM Admin integrado e pipeline operacional derivado

## Data
2026-08-03

## Contexto
O painel administrativo precisava adotar a linguagem visual das referências BizLink sem perder as configurações e operações já existentes. Um plano anterior propunha criar um segundo frontend, mas isso duplicaria autenticação, build, deploy e manutenção. O domínio também não possui entidade de deal ou estágio comercial persistente.

## Decisão tomada
1. Redesenhar o admin existente em `/admin/dash`, dentro do projeto `dashboard/`.
2. Preservar autenticação cookie-only, autorização `admin`/`owner` e endpoints atuais.
3. Escopar o design system monocromático em `.admin-layout-v2`.
4. Representar o ciclo do usuário como segmentação calculada e somente leitura: Atenção, Ativos, Novos e Em ativação.
5. Não implementar drag-and-drop nem persistência de estágio nesta fase.
6. Substituir controles fictícios da referência por ações reais, como Broadcast.

## Motivo
- Evita duplicação arquitetural e riscos de autenticação divergente.
- Mantém todas as funções administrativas no mesmo fluxo.
- Exibe apenas métricas e estados sustentados pelos dados reais do FreddyBot.
- Permite adicionar um CRM comercial persistente futuramente como mudança de domínio explícita.

## Impacto
- Novo overview com gráfico, taxa de ativação, KPIs e quadro operacional.
- Busca, filtros, ordenação e URL das abas passam a funcionar sem reload.
- Nenhuma migration, endpoint ou regra de negócio foi adicionada.
- A Mini App de usuários não recebe os tokens visuais do CRM.

# Decisão: Fidelidade BizLink com tema claro fixo no admin

## Data
2026-08-03

## Contexto
O primeiro redesenho preservou a arquitetura e a linguagem monocromática, mas ainda se afastava da referência: o tema automático podia abrir o admin escuro, a faixa de analytics tinha quatro KPIs em grade, o gauge era um arco grosso e os espaçamentos lembravam um dashboard genérico.

## Decisão tomada
1. Tornar a paleta clara da referência invariável dentro de `.admin-layout-v2`, sem modificar o tema da Mini App.
2. Tratar topbar e hero como uma única faixa creme.
3. Usar gráfico de barras pareadas, gauge de 45 marcas radiais e somente dois KPIs reais.
4. Iniciar o quadro imediatamente após o hero e limitar o card preto a um único administrador real.
5. Manter todos os módulos e controles administrativos reais, substituindo apenas sua apresentação.

## Motivo
- A referência é explicitamente clara e declara os quatro tons usados.
- A correspondência depende mais de proporção, densidade e composição do que de elementos decorativos.
- Um tema administrativo invariável elimina diferenças causadas por horário, preferência do sistema ou Telegram.
- Dados comerciais fictícios continuariam incompatíveis com o domínio do FreddyBot.

## Impacto
- O admin mantém a mesma aparência em qualquer tema global.
- A primeira dobra em 1024×768 corresponde à estrutura visual da referência.
- Busca, filtros, ordenação, broadcast, configurações e demais módulos continuam funcionais.
- Não há alteração de backend, banco ou deploy.

# Decisão: Sistema visual Minimal UI unificado no admin

## Data
2026-08-03

## Contexto
O usuário não aprovou a direção BizLink. Embora o overview tivesse sido redesenhado, Broadcast, Auditoria, Logs, Configurações, MTProto, Features e Assinaturas ainda empregavam layouts e superfícies visualmente desconectados.

## Decisão tomada
1. Adotar a linguagem do Minimal UI apenas como direção visual, sem instalar Material UI nem copiar componentes ou assets do kit.
2. Escopar os tokens ao admin: canvas neutro, paper branco, bordas discretas, sombras baixas, raio de 16 px, sidebar de 300 px e header de 72 px.
3. Usar `#007867` como ação primária em vez do verde mais claro do kit, pois fornece contraste de 5,41:1 com texto branco em botões pequenos.
4. Reutilizar `AdminPageHeader`, `AdminMetricCard` e `AdminEmptyState` para eliminar variações de hierarquia e estados.
5. Preservar handlers, APIs, confirmações e permissões de todas as dez abas.

## Motivo
- O stack já fornece React, Tailwind, shadcn e Base UI; adicionar MUI aumentaria o custo e criaria um segundo sistema de componentes.
- A coerência operacional depende de padrões compartilhados, não de trocar apenas cores no overview.
- A variante verde escura mantém a identidade Minimal UI sem comprometer legibilidade.

## Impacto
- Todas as abas renderizam no mesmo shell responsivo e com a mesma hierarquia visual.
- Não houve mudança de backend, banco, dependência ou efeito externo.
- A Mini App de usuários permanece sem receber os tokens administrativos.

## Decisão: Autoridade derivada da sessão e propriedade persistida

### Data
2026-08-03

### Contexto
A auditoria identificou que login, transferência de canal e agendamento aceitavam identificadores informados pelo cliente como fonte de autoridade. Isso permitia comparação parcial de IDs Telegram, transferência por terceiro e criação de posts para canais alheios.

### Decisão tomada
1. Validar o `user.id` do initData Telegram por JSON e igualdade numérica exata.
2. Remover `oldOwnerId` do contrato de transferência; o ator vem do JWT e a posse atual vem do banco.
3. Validar posse do canal dentro do `SchedulerService`, cobrindo Dashboard e PostBuilder.
4. Invalidar o cache do canal após transferência e limitar/rate-limit o endpoint público de logs do cliente.

### Motivo
Campos de requisição expressam intenção, não autorização. A fonte de verdade para identidade é a sessão autenticada e a fonte de verdade para propriedade é o banco de dados.

### Impacto
- Usuários comuns não transferem nem agendam em canais de terceiros.
- Admins e owner mantêm o bypass administrativo explícito de transferência.
- O cache deixa de manter autorização desatualizada depois da migração de dono.
- Testes de regressão cobrem prefixos de ID, transferência indevida, agendamento indevido e abuso do endpoint de logs.

# Decisão: Alertas administrativos derivados e Broadcast com alcance estimado

## Data
2026-08-03

## Contexto
O painel administrativo já recebe usuários e canais por `/api/admin/overview`, mas não possui tabela de notificações, API de fila de broadcast, telemetria de entrega ou mecanismo de push para owner/admin.

## Decisão tomada
1. Derivar alertas in-app de blacklist, ativação pendente e novos cadastros a partir da resposta real da visão geral.
2. Exibir essas notificações na visão geral e na barra superior, com navegação para a revisão de usuários.
3. Calcular o alcance do Broadcast no cliente apenas como estimativa da base carregada ou dos IDs fornecidos.
4. Informar que o backend confirma o início do processamento, não a entrega individual.

## Motivo
Evitar métricas, histórico ou notificações externas fictícias enquanto o backend não persistir e não expuser esses eventos.

## Impacto
O owner ganha contexto acionável sem nova migration, endpoint ou efeito externo. Push por Telegram/e-mail e histórico de execução continuam sendo uma evolução de domínio separada.
