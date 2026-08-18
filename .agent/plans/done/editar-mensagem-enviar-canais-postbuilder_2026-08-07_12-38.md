# Plano: Editar Mensagem no "Enviar para Canais" e Adicionar Botão Voltar ao Post

## Pedido do usuário
Na ação "📢 Enviar para Canais" (`pb-send-to-channels`), o bot deve EDITAR a mensagem existente no chat (em vez de enviar uma nova mensagem) e incluir o botão de voltar ao post salvo.

## Objetivo
Padronizar a ação "📢 Enviar para Canais" para reutilizar a mensagem via `EditMessageText`, exibir alerta popup nativo (`ShowAlert: true`) se o usuário não possuir canais, e disponibilizar o botão `🔙 Voltar ao Post` (`pb-saved-menu:<sessionID>`).

## Contexto atual
- `handleSendToChannelsTelego` usava `bot.SendMessage` (criando nova mensagem).
- Se não houvesse canais, enviava mensagem de texto no chat em vez de alert popup.
- A lista de canais para envio não possuía o botão `🔙 Voltar ao Post`.

## Arquivos analisados
- `internal/telegram/handlers/events/postBuilder/postBuilder.go`

## Arquivos que poderão ser modificados
- `internal/telegram/handlers/events/postBuilder/postBuilder.go`

## Estratégia de implementação

1. **Alerta nativo sem editar quando 0 canais**:
   - Em `handleSendToChannelsTelego`, se `len(channels) == 0`:
     - Disparar `bot.AnswerCallbackQuery` com `ShowAlert: true` e a mensagem *"⚠️ Você não possui nenhum canal cadastrado para enviar postagens!"*.
     - Retornar imediatamente sem alterar a mensagem do chat.

2. **Edição da mensagem via `EditMessageText`**:
   - Se houver canais, usar `bot.EditMessageText` utilizando o `MessageID` do callback.
   - Adicionar a linha de botão `{Text: "🔙 Voltar ao Post", CallbackData: "pb-saved-menu:" + sessionID}` ao teclado inline.

3. **Botão de retorno na confirmação de envio em `handleSendApplyTelego`**:
   - Na mensagem de confirmação de envio concluído com sucesso (`✅ Postagem enviada com sucesso para o canal!`), anexar teclado inline com o botão `{Text: "🔙 Voltar ao Post", CallbackData: "pb-saved-menu:" + sessionID}`.

---

## Passos detalhados

### Passo 1: Atualizar a chamada em `CallbackHandlerTelego`
- Passar `messageID` e `update.CallbackQuery.ID` para `handleSendToChannelsTelego`.

### Passo 2: Atualizar `handleSendToChannelsTelego`
- Se `len(channels) == 0`: chamar `AnswerCallbackQuery` com `ShowAlert: true`.
- Se houver canais: chamar `EditMessageText` no `messageID` e incluir o botão `🔙 Voltar ao Post`.

### Passo 3: Atualizar `handleSendApplyTelego`
- Incluir o botão `🔙 Voltar ao Post` na mensagem de envio concluído com sucesso.

### Passo 4: Compilação e Testes
- Executar `go build ./cmd/FreddyBot` e `go test ./...`.

---

## Riscos
- Nenhum. Edição de mensagens em callbacks é o padrão estrito do sistema.

## Impactos esperados
- A tela do post salvo é editada diretamente para a seleção de canais sem criar mensagens extras no chat.
- Alerta popup nativo no Telegram quando não houver canais.
- Retorno fácil e limpo ao post salvo.

## Compatibilidade
- Linux / macOS / Windows / Docker / Telego

## Como testar
1. Clicar em "📢 Enviar para Canais" sem canais -> Verificar o alerta popup na tela.
2. Com canais cadastrados, clicar em "📢 Enviar para Canais" -> Verificar que a mensagem é editada diretamente para a lista de canais com o botão "🔙 Voltar ao Post".
3. Clicar em "🔙 Voltar ao Post" -> Verificar que a mensagem é editada de volta para o menu do post salvo com o token inline.
