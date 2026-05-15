# Plano: preservar-formatacao-postbuilder

## Pedido do usuário
Ao definir um Título ou Rodapé com formatação (ex: negrito enviado pelo app do Telegram) ou ao importar dados de outro canal, a formatação está sendo perdida no PostBuilder.

## Objetivo
Garantir que a formatação original (seja via Markdown ou via Entidades do Telegram) seja preservada ao definir Título, Corpo ou Rodapé no PostBuilder, e também ao importar dados de canais.

## Contexto atual
- `handleTextInput` em `internal/telegram/events/postBuilder/postBuilder.go` está salvando o `text` bruto diretamente nos campos `Title`, `Body` e `Footer`.
- A função `ProcessTextWithFormatting` em `internal/telegram/events/channelPost/formatting.go` é a responsável por converter entidades e markdown em HTML, mas ela não está sendo usada para o input do usuário (foi removida em uma iteração anterior para dar lugar ao JIT).
- Ao importar de um canal, o `state.Body` recebe a legenda padrão que já pode estar em HTML, mas o JIT pode tentar re-processar isso.

## Arquivos analisados
- `internal/telegram/events/postBuilder/postBuilder.go`
- `internal/telegram/events/channelPost/formatting.go`

## Arquivos que poderão ser modificados
- `internal/telegram/events/postBuilder/postBuilder.go`

## Estratégia de implementação
1. **Voltar a processar o input com Entidades:** No `handleTextInput`, utilizaremos `channelpost.ProcessTextWithFormatting(text, update.Message.Entities)` para capturar a formatação enviada pelo usuário (negrito, itálico, links feitos pela UI do Telegram).
2. **Armazenamento em HTML:** Como o sistema JIT e a visualização no Menu/Preview esperam HTML (usando `DetectParseMode`), salvar o texto já convertido para HTML no `state` é a forma mais segura de não perder entidades.
3. **Ajuste no JIT:** Garantir que o `DetectParseMode` no momento do Preview/Inline Query não quebre o HTML já existente (a função `isHTML` já deve lidar com isso).

## Passos detalhados

1. Editar `internal/telegram/events/postBuilder/postBuilder.go`:
    - No `handleTextInput`, chamar `channelpost.ProcessTextWithFormatting(text, update.Message.Entities)` e atribuir o resultado a uma variável `formattedText`.
    - Usar `formattedText` ao invés de `text` bruto para `state.Title`, `state.Body` e `state.Footer`.

## Riscos
- **Duplo Escape:** Se o texto já estiver em HTML e passar por um escape acidental. (Mitigado pela lógica de `DetectParseMode` e `isHTML`).
- **Markdown vs HTML:** Usuários que digitam Markdown na mão (`**texto**`) terão isso convertido para HTML no momento do save. Isso é desejado para preservar a intenção.

## Impactos esperados
- Títulos e Rodapés enviados com negrito/itálico pelo app serão preservados.
- Importação de canais manterá a formatação da legenda padrão.

## Como testar

### Execução
1. Abrir o bot.
2. Iniciar o PostBuilder com uma mídia.
3. Definir Título usando a ferramenta de negrito do Telegram.
4. Verificar se o Menu e o Preview mostram o negrito.
5. Importar dados de um canal que possua legenda formatada e verificar se a formatação é mantida.
