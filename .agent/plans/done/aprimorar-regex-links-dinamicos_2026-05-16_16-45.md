# Plano: aprimorar-regex-links-postbuilder

## Pedido do usuário
O usuário relatou que links com a formatação Markdown `[t.me/legendasbot](https://t.me/FreddyCaptionBot)` inseridos no PostBuilder não estão sendo formatados corretamente (ou seja, o texto continua exibindo a sintaxe crua em vez de virar um link embutido).

## Objetivo
Melhorar a robustez das expressões regulares (Regex) na função `DetectParseMode` para garantir que qualquer sintaxe válida de link Markdown seja capturada e convertida para tags HTML `<a>`, mesmo se contiver espaços em branco ou emojis ao redor.

## Contexto atual
- O PostBuilder usa a função `DetectParseMode` (via `ProcessTextWithFormattingTelego`) para processar Markdown.
- A regex atual de links é `\[([^\]]+)\]\((https?://[^\s)]+)\)`. Ela exige que a URL comece obrigatoriamente com `http` ou `https` e não contenha espaços.
- Caso o usuário digite o link de uma forma ligeiramente diferente ou com caracteres que a regex rejeita, a conversão não ocorre e o texto cru é salvo e exibido.

## Arquivos analisados
- `internal/telegram/events/channelPost/utils_v2.go`
- `internal/telegram/events/channelPost/formatting_telego.go`

## Arquivos que poderão ser modificados
- `internal/telegram/events/channelPost/utils_v2.go`

## Estratégia de implementação
1. **Flexibilizar a Regex de Link:** Alterar a regex em `DetectParseMode` para `\[(.*?)\]\((.*?)\)`. Isso permite capturar de forma "preguiçosa" qualquer conteúdo entre os colchetes e parênteses, delegando ao Telegram a validação final da URL (que já é nativa do ParseMode HTML).
2. **Adicionar Regex Secundária:** Caso o usuário omita o `http://` (ex: `[texto](t.me/canal)`), a regex anterior também falharia. A nova regex mais flexível resolverá esse caso, sendo necessário apenas garantir que adicionamos o prefixo `http://` se ele faltar na hora de montar a tag `<a>`.

## Passos detalhados

1. **Em `internal/telegram/events/channelPost/utils_v2.go`:**
   - Modificar a função `DetectParseMode`.
   - Substituir a regex de links por: `linkRegex := regexp.MustCompile(\`\[(.*?)\]\((.*?)\)\`)`.
   - Substituir a linha do `ReplaceAllString` por uma função `ReplaceAllStringFunc` para permitir tratamento da URL (ex: adicionar `https://` caso a pessoa tenha digitado apenas `t.me/...`).

   ```go
   linkRegex := regexp.MustCompile(`\[(.*?)\]\((.*?)\)`)
   res = linkRegex.ReplaceAllStringFunc(res, func(m string) string {
       matches := linkRegex.FindStringSubmatch(m)
       if len(matches) == 3 {
           text := matches[1]
           url := strings.TrimSpace(matches[2])
           if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
               url = "https://" + url
           }
           return fmt.Sprintf(`<a href="%s">%s</a>`, url, text)
       }
       return m
   })
   ```

2. **Melhorar as demais Regex (Opcional mas recomendado):**
   - Garantir que Negrito, Itálico e Código usem lazy matching `(.*?)` para evitar falhas caso a mensagem seja muito longa ou tenha múltiplas quebras de linha que quebrem a regex atual baseada em negação `[^\n]`.

## Riscos
- Mudar para `(.*?)` sem restrições pode capturar coisas que não são links se o texto for muito mal formatado, mas o escopo dentro de `[` e `(` protege contra falsos positivos.

## Impactos esperados
- O PostBuilder formatará instantaneamente qualquer variação de link Markdown (com ou sem http, com espaços internos, etc).

## Compatibilidade
- Linux
- macOS
- Windows
- Docker
- CI/CD

## Como testar
1. Enviar a mensagem fornecida pelo usuário no PostBuilder.
2. Verificar se o Preview e o envio Final geram o link clicável.
