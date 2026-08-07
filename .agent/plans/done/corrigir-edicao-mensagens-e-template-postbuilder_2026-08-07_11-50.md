# Plano: Corrigir Edição de Mensagens no Chat e Voltar no "Salvar como Template"

## Pedido do usuário
1. O botão de voltar da função "💾 Salvar como Template" volta para o menu inicial do PostBuilder (`pb-start`), quando deveria voltar para o menu do Post Salvo com o token (`@BotUsername pb <id>`).
2. O PostBuilder está enviando novas mensagens no chat em vez de editar as mensagens existentes, gerando poluição visual e excesso de mensagens (especialmente ao clicar em "🛠️ Post Builder" na mensagem de detecção de mídia).

## Objetivo
- Fazer com que o fluxo do PostBuilder edite a mensagem existente sempre que possível (transformando a mensagem de detecção de mídia no menu e editando prompts/menus no mesmo local).
- Corrigir o botão de voltar no "Salvar como Template" para retornar ao menu da postagem salva (`pb-saved-menu:<sessionID>`).

## Contexto atual
- Ao clicar em `🛠️ Post Builder` (`pb-start`), `state.MenuMessageID` é `0`, fazendo `showMenuTelego` enviar uma nova mensagem abaixo da mensagem de detecção de mídia.
- Ao clicar em `💾 Salvar como Template`, o bot usava `SendMessage` (nova mensagem) e `promptBackKB()` (que apontava para `pb-start`).
- Ao salvar o template via mensagem de texto, o bot enviava uma nova mensagem e tentava abrir o menu de edição rascunho em vez do menu do post salvo.

## Arquivos analisados
- `internal/telegram/handlers/events/postBuilder/postBuilder.go`

## Arquivos que poderão ser modificados
- `internal/telegram/handlers/events/postBuilder/postBuilder.go`

## Estratégia de implementação

1. **Transformar a mensagem de detecção de mídia no Menu Principal**:
   - Quando o usuário clica em `🛠️ Post Builder` (`pb-start`), definir `state.MenuMessageID = update.CallbackQuery.Message.GetMessageID()`.
   - Dessa forma, `showMenuTelego` editará a própria mensagem `"✨ Mídia detectada..."` transformando-a no menu do PostBuilder, sem criar uma nova mensagem no chat.

2. **Corrigir o "Salvar como Template"**:
   - Em `pb-save-template:<sessionID>`:
     - Editar a mensagem atual (`EditMessageText`) com o prompt `"💾 Envie o nome para salvar este post como template:"`.
     - Anexar um botão de retorno limpo: `{Text: "🔙 Voltar ao Post", CallbackData: "pb-saved-menu:" + sessionID}`.
     - Armazenar `sessionID` no estado para saber de qual post salvo veio o template.
   - Em `handleTextInputTelego` (ao receber o nome do template):
     - Salvar o template no banco.
     - Responder com confirmação e botões para voltar ao post salvo (`pb-saved-menu:<sessionID>`) ou ir aos agendamentos.

3. **Manter edição contínua nas edições de campos**:
   - Garantir que prompts e seleções (botões, templates, canais) editem a mensagem existente em vez de acumular novas mensagens no histórico.

---

## Passos detalhados

### Passo 1: Ajustar o handler de `pb-start` em `postBuilder.go`
- Ao receber o callback `pb-start`, se `update.CallbackQuery` não for nil, definir `state.MenuMessageID = update.CallbackQuery.Message.GetMessageID()`.

### Passo 2: Refatorar `pb-save-template:<sessionID>`
- Substituir `SendMessage` por `EditMessageText`.
- Usar teclado com o botão `🔙 Voltar ao Post` (`pb-saved-menu:` + sessionID).
- Salvar `sessionID` no `PostBuilderState`.

### Passo 3: Atualizar o salvamento de template em `handleTextInputTelego`
- Ao salvar o template em `case "awaiting_template_name"`, verificar se o estado possui um `sessionID` de post salvo.
- Se possuir, responder com teclado contendo `🔙 Voltar ao Post` (`pb-saved-menu:` + sessionID).

### Passo 4: Compilação e Testes
- Executar `go build ./cmd/FreddyBot` e `go test ./...`.

---

## Riscos
- Nenhum. Edição de mensagem é suportada nativamente pela Bot API (`editMessageText`).

## Impactos esperados
- Chat limpo e 100% interativo, sem acúmulo de mensagens de menu.
- A mensagem de detecção de mídia é reaproveitada como menu.
- Botão "Voltar ao Post" no template direciona para o post salvo correto.

## Compatibilidade
- Linux / macOS / Windows / Docker / Telego

## Como testar
1. Enviar foto/vídeo -> Clicar em "🛠️ Post Builder" -> Verificar que a mensagem de detecção foi editada em vez de criar uma mensagem duplicada.
2. Salvar post -> Clicar em "💾 Salvar como Template" -> Verificar que a mensagem foi editada e contém o botão "🔙 Voltar ao Post".
3. Clicar em "🔙 Voltar ao Post" -> Verificar que retorna ao card do post salvo com o token inline.
