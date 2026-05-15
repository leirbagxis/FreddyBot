# Plano: corrigir-postbuilder-markdown-entities

## Pedido do usuário
O PostBuilder parou de aplicar a formatação Markdown (ex: links `[texto](url)`) no campo de `body` (corpo) da postagem.

## Objetivo
Corrigir a função `ProcessTextWithFormatting` para garantir que mensagens com marcações Markdown explícitas sejam processadas corretamente, mesmo quando o Telegram envia entidades automáticas (como URLs). Além disso, garantir que o texto seja corretamente escapado para HTML.

## Contexto atual
Atualmente, em `internal/telegram/events/channelPost/formatting.go`, a função `ProcessTextWithFormatting` verifica primeiro se existem `entities` na mensagem (`len(entities) > 0`). Se existirem, ela chama `ProcessEntitiesOnly`, que ignora qualquer sintaxe Markdown e foca apenas nas entidades. Como o Telegram adiciona automaticamente entidades do tipo `url` para links, qualquer texto contendo um link (mesmo formatado como `[texto](url)`) acaba acionando este caminho e ignorando o Markdown.

## Arquivos analisados
- `internal/telegram/events/channelPost/formatting.go`
- `internal/telegram/events/postBuilder/postBuilder.go`

## Arquivos que poderão ser modificados
- `internal/telegram/events/channelPost/formatting.go`

## Estratégia de implementação
1. **Inversão de Prioridade**: Modificar `ProcessTextWithFormatting` para verificar primeiro se o texto contém Markdown explícito (`isMarkdown(text)`). Se contiver, priorizar a conversão de Markdown (`DetectParseMode`) para evitar que entidades automáticas sequestrem o fluxo.
2. **Correção de Escapamento no Markdown**: A função `convertMarkdownToHTML` atualmente não realiza escape HTML em partes do texto que não são formatadas, o que pode quebrar a formatação no Telegram (parse mode HTML). Vamos integrar a lógica de `escapeRemainingText` ou refatorar o parser para garantir que os caracteres especiais (`<`, `>`, `&`) sejam escapados, mantendo apenas as tags HTML geradas pelo parser.

## Passos detalhados
1. Atualizar `ProcessTextWithFormatting`:
   ```go
   func ProcessTextWithFormatting(text string, entities []models.MessageEntity) string {
       if text == "" {
           return ""
       }

       // 1. Se detectarmos markdown explícito, processamos como markdown.
       // Evita que entidades automáticas (url) ignorem a conversão de markdown.
       if isMarkdown(text) {
           return DetectParseMode(text)
       }

       // 2. Se houver entidades (client-side formatting sem markdown explícito)
       if len(entities) > 0 {
           return ProcessEntitiesOnly(text, entities)
       }

       // 3. Fallback
       return DetectParseMode(text)
   }
   ```
2. Refatorar `convertMarkdownToHTML` em `formatting.go` para usar expressões regulares mais seguras e realizar o escape de caracteres HTML puro antes de aplicar as transformações.

## Impactos esperados
- O PostBuilder voltará a aceitar e converter marcações de links, negrito e outras formatações Markdown normalmente.

## Como testar
Rodar o script de reprodução que criamos para garantir que o link seja convertido para `<a href="url">texto</a>` mesmo com a entidade URL presente.
