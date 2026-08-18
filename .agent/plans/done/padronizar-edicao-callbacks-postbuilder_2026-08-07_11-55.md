# Plano: Padronizar Edição em Callbacks e Envio de Novo Menu no Texto do Usuário

## Pedido do usuário
1. **No botão "✅ Salvar"**: O bot estava enviando uma nova mensagem (`SendMessage`) em vez de editar a mensagem existente. Deve SEMPRE editar a mensagem ao receber um callback.
2. **Regra de Comportamento do PostBuilder**:
   - **Ação por Botão/Callback**: O bot deve SEMPRE EDITAR a mensagem existente (`EditMessageText`).
   - **Ação por Entrada de Texto do Usuário**: Quando o usuário enviar o texto no chat (ex: enviou o título digitado), o bot envia uma NOVA mensagem (`SendMessage`) com o menu atualizado abaixo do texto do usuário.

## Objetivo
Garantir que 100% dos botões inline (callbacks) editem a mensagem existente no chat (inclusive "Salvar", "Editar Título", "Editar Corpo", "Editar Rodapé", "Reações" e "Adicionar Botão"), enquanto as entradas de texto do usuário façam o bot responder com o novo menu abaixo da mensagem enviada.

## Contexto atual
- `pb-save` usava `bot.SendMessage` (enviava nova mensagem).
- `pb-edit-title`, `pb-edit-body`, `pb-edit-footer`, `pb-edit-reactions` e `pb-add-button` usavam `bot.SendMessage` (enviavam nova mensagem).
- Ao digitar texto, `handleTextInputTelego` já chama `showMenuTelego(..., replyToMessageID)` (envia novo menu abaixo da resposta do usuário), que está em total sintonia com o pedido.

## Arquivos analisados
- `internal/telegram/handlers/events/postBuilder/postBuilder.go`

## Arquivos que poderão ser modificados
- `internal/telegram/handlers/events/postBuilder/postBuilder.go`

## Estratégia de implementação

1. **Atualizar `pb-save`**:
   - Substituir `bot.SendMessage` por `bot.EditMessageText` utilizando o `MessageID` da callback query (`update.CallbackQuery.Message.GetMessageID()`).
   - O menu do PostBuilder é convertido diretamente no card de confirmação de salvamento (`✅ Postagem salva com sucesso! Utilize o modo inline para enviar: @bot pb id`).

2. **Atualizar prompts de edição em callbacks**:
   - Em `pb-edit-title`, `pb-edit-body`, `pb-edit-footer`, `pb-edit-reactions` e `pb-add-button`:
     - Substituir `bot.SendMessage` por `bot.EditMessageText` usando `messageID := update.CallbackQuery.Message.GetMessageID()`.
     - Definir `state.PromptMessageID = messageID`.

3. **Preservar envio de novo menu nas respostas de texto do usuário**:
   - Em `handleTextInputTelego`: manter o envio da mensagem de menu abaixo da resposta do usuário (`showMenuTelego(..., update.Message.MessageID)`).

---

## Passos detalhados

### Passo 1: Atualizar `case "pb-save"` em `postBuilder.go`
- Obter `messageID := update.CallbackQuery.Message.GetMessageID()`.
- Trocar `bot.SendMessage` por `bot.EditMessageText` utilizando `messageID`.

### Passo 2: Atualizar os callbacks de prompt em `postBuilder.go`
- `pb-edit-title`, `pb-edit-body`, `pb-edit-footer`, `pb-edit-reactions` e `pb-add-button`:
  - Usar `messageID := update.CallbackQuery.Message.GetMessageID()`.
  - Executar `bot.EditMessageText` para exibir o prompt no lugar do menu.
  - Atualizar `state.PromptMessageID = messageID`.

### Passo 3: Compilação e Testes
- Executar `go build ./cmd/FreddyBot` e `go test ./...`.

---

## Riscos
- Nenhum. Edição em callbacks é um padrão estrito e limpo do Telegram Bot API.

## Impactos esperados
- 100% dos cliques em botões editam a mensagem sem criar mensagens duplicadas.
- Ao digitar texto, o menu é enviado como nova mensagem abaixo da digitação do usuário.
- Experiência fluida e interativa.

## Compatibilidade
- Linux / macOS / Windows / Docker / Telego

## Como testar
1. Clicar em "📝 Título" -> Verificar que o menu foi editado para a mensagem de prompt no mesmo balão.
2. Enviar o título via texto -> Verificar que o bot envia o menu atualizado abaixo do texto.
3. Clicar em "✅ Salvar" -> Verificar que a mensagem do menu é editada para "✅ Postagem salva com sucesso!".
