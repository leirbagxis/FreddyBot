# Plano: corrigir-prioridade-markdown-entidades

## Pedido do usuário
O usuário relatou que a formatação de um "link embutido" (criado nativamente na interface do Telegram) não está funcionando no PostBuilder. A mensagem enviada foi: `🐈‍⠀៹ [t.me/legendasbot]  ‹` (com o texto dentro dos colchetes transformado em link clicável no app).

## Objetivo
Corrigir a função `ProcessTextWithFormattingTelego` para que formatações nativas explícitas do Telegram (como `text_link` ou `bold` aplicados pela UI) não sejam ignoradas quando o texto também contiver caracteres que se pareçam com Markdown (como `[`).

## Contexto atual
- O PostBuilder chama `ProcessTextWithFormattingTelego` para converter a entrada do usuário em HTML.
- Atualmente, a função começa verificando `if IsMarkdown(text)`. Se o texto contiver o caractere `[`, ela imediatamente assume que é Markdown cru e descarta completamente o array `entities` fornecido pelo Telegram.
- Ao descartar o array, o link embutido (que era uma entidade `text_link`) é perdido. O parser do Markdown tenta encontrar a sintaxe `[texto](url)` e falha (pois não há `(url)`), resultando no texto exibido sem formatação.

## Arquivos analisados
- `internal/telegram/events/channelPost/formatting_telego.go`

## Arquivos que poderão ser modificados
- `internal/telegram/events/channelPost/formatting_telego.go`

## Estratégia de implementação
1. **Inteligência de Entidades:** Modificar `ProcessTextWithFormattingTelego` para analisar as entidades ANTES de aplicar o fallback de Markdown.
2. **Filtragem de Autodetecção:** O Telegram cria entidades automáticas para URLs e hashtags mesmo quando o usuário não as formatou. Devemos diferenciar essas "entidades automáticas" das "entidades manuais" (como `text_link`, `bold`, `italic`, `custom_emoji`, etc.).
3. **Nova Precedência:**
   - Se houver qualquer entidade "manual" (o usuário explicitly formatou algo), processaremos as entidades e ignoraremos o Markdown.
   - Caso contrário (se só houver entidades automáticas ou nenhuma), verificamos se o usuário tentou usar Markdown (`IsMarkdown`).
   - Se não for Markdown, usamos as entidades automáticas como fallback final.

## Passos detalhados

1. **Em `internal/telegram/events/channelPost/formatting_telego.go`:**
   Substituir a implementação de `ProcessTextWithFormattingTelego` por:
   ```go
   func ProcessTextWithFormattingTelego(text string, entities []telego.MessageEntity) string {
       if text == "" {
           return ""
       }

       hasManualEntities := false
       if len(entities) > 0 {
           for _, e := range entities {
               // Ignora entidades automáticas na verificação
               if e.Type != "url" && e.Type != "email" && e.Type != "phone_number" && e.Type != "hashtag" && e.Type != "cashtag" && e.Type != "bot_command" {
                   hasManualEntities = true
                   break
               }
           }
       }

       // Se o usuário formatou algo explicitamente via UI do Telegram
       if hasManualEntities {
           return ProcessEntitiesOnlyTelego(text, entities)
       }

       // Se usou sintaxe Markdown crua
       if IsMarkdown(text) {
           return DetectParseMode(text)
       }

       // Fallback para entidades automáticas
       if len(entities) > 0 {
           return ProcessEntitiesOnlyTelego(text, entities)
       }

       return DetectParseMode(text)
   }
   ```

## Riscos
- Mínimos. Essa lógica garante que não perderemos a formatação explícita do usuário, seja ela via UI ou via sintaxe Markdown.

## Impactos esperados
- Links embutidos (selecionar texto e criar link) funcionarão perfeitamente, mesmo que o texto contenha colchetes ou sublinhados que antes sequestrariam o parser.

## Compatibilidade
- Linux
- macOS
- Windows
- Docker
- CI/CD

## Como testar
1. Criar um post no PostBuilder usando a interface do Telegram para gerar um link embutido em um texto contendo colchetes.
2. O PostBuilder deve exibir e enviar a mensagem com o link formatado em azul.
