# Decisões Arquiteturais

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
