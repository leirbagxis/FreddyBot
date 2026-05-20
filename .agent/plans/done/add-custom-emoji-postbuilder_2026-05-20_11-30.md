# Plano: add-custom-emoji-postbuilder

## Pedido do usuário
O usuário solicitou que o PostBuilder suporte o uso de emojis customizados do Telegram (entidade `custom_emoji`).

## Objetivo
Garantir que emojis customizados enviados no fluxo do PostBuilder sejam corretamente armazenados e renderizados como tags `<tg-emoji emoji-id="...">` sem serem corrompidos pelas rotinas de fallback ou escape HTML.

## Contexto atual
- A rotina `ProcessEntitiesOnlyTelego` em `internal/telegram/events/channelPost/formatting_telego.go` já converte `custom_emoji` para a tag HTML `<tg-emoji emoji-id="...">`.
- A salvaguarda recém-adicionada no PostBuilder (`InlineHandlerTelego` e `sendFinalPostTelego`) chama `DetectParseMode` se houver caracteres Markdown e *nenhuma* tag HTML básica (`<a href=` ou `<b>`).
- O problema: `DetectParseMode` começa com `res := html.EscapeString(text)`. Se a string contiver a tag `<tg-emoji>` e cair na salvaguarda (ex: o usuário misturou um caractere de markdown acidental como `_` com um emoji customizado), a tag HTML do emoji será "escapada" (transformada em `&lt;tg-emoji&gt;`) e quebrará o formato perante o Telegram.

## Arquivos analisados
- `internal/telegram/handlers/events/postBuilder/postBuilder.go`

## Arquivos que poderão ser modificados
- `internal/telegram/handlers/events/postBuilder/postBuilder.go`

## Estratégia de implementação
1. Modificar a condição da salvaguarda (Safeguard) no PostBuilder para considerar a tag `<tg-emoji` como uma tag HTML válida. Se ela estiver presente, não devemos aplicar `DetectParseMode`, evitando assim que o `html.EscapeString` quebre o código.

## Passos detalhados
1. Em `internal/telegram/handlers/events/postBuilder/postBuilder.go`, modificar a condição da salvaguarda em **duas funções** (`InlineHandlerTelego` e `sendFinalPostTelego`):
   **De:**
   ```go
   if channelpost.IsMarkdown(caption) && !strings.Contains(caption, "<a href=") && !strings.Contains(caption, "<b>") {
   ```
   **Para:**
   ```go
   if channelpost.IsMarkdown(caption) && !strings.Contains(caption, "<a href=") && !strings.Contains(caption, "<b>") && !strings.Contains(caption, "<tg-emoji") {
   ```

## Riscos
- O risco é muito baixo. Estamos apenas ampliando a lista de exceções da salvaguarda para evitar o "double escape" de tags de emoji customizado.

## Impactos esperados
- Emojis customizados em mensagens salvas no PostBuilder ou exibidas no modo Inline/Preview serão renderizados corretamente e aceitos pela API do Telegram sem o erro de Bad Request ou formatação quebrada.

## Como testar
1. Criar um PostBuilder com uma mensagem contendo um emoji customizado do Telegram e um texto com `_`.
2. Verificar o Preview.
3. Compartilhar via Inline.
