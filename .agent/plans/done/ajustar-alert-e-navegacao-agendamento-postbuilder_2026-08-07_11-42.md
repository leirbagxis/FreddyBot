# Plano: Alerta NATIVO (AnswerCallbackQuery Alert) sem editar mensagem e Remoção de Cancelar

## Pedido do usuário
1. Ao clicar em "Agendar Envio", se o usuário não tiver nenhum canal cadastrado, NÃO editar nem destruir NENHUMA mensagem. Apenas responder à CallbackQuery com `AnswerCallbackQuery` (`ShowAlert: true`) avisando que ele não possui nenhum canal cadastrado para usar a funcionalidade.
2. Garantir que o botão "Voltar ao Post" restaure o menu da postagem salva com o token sem dar "sessão expirada".
3. Remover o botão "Cancelar" das telas do fluxo de agendamento.

## Objetivo
Emitir um alerta nativo popup via Telegram sem alterar nenhuma mensagem do chat quando o usuário não tiver canais, remover botões "Cancelar" redundantes e assegurar o funcionamento limpo de navegação.

## Contexto atual
- Em `handleSchedulePost`, quando `len(channels) == 0`, a função tentava editar a mensagem.
- O botão `❌ Cancelar` ainda existia no prompt de agendamento.

## Arquivos analisados
- `internal/telegram/handlers/events/postBuilder/postBuilder.go`

## Arquivos que poderão ser modificados
- `internal/telegram/handlers/events/postBuilder/postBuilder.go`

## Estratégia de implementação

1. **AnswerCallbackQuery com `ShowAlert: true` (Sem editar mensagem)**:
   - No callback `pb-schedule:<sessionID>`, verificar os canais do usuário ANTES de chamar a edição.
   - Se `len(channels) == 0` (ou se der erro):
     - Executar `bot.AnswerCallbackQuery(context.Background(), &telego.AnswerCallbackQueryParams{ CallbackQueryID: update.CallbackQuery.ID, Text: "⚠️ Você não possui nenhum canal cadastrado para usar esta funcionalidade!", ShowAlert: true })`
     - Dar `return nil` imediatamente.
     - Nenhuma mensagem será editada, apagada ou enviada no chat.

2. **Remover botão "Cancelar" e padronizar navegação de Voltar**:
   - Na tela de prompt (data/hora): deixar apenas os botões `🔙 Voltar` (`pb-sch:<sessionID>:<channelID>`) e `🔙 Voltar ao Post` (`pb-saved-menu:<sessionID>`). Sem botão "Cancelar".
   - Na seleção de tipo: deixar apenas `🔙 Voltar ao Canal` (`pb-schedule:<sessionID>`) e `🔙 Voltar ao Post` (`pb-saved-menu:<sessionID>`).

3. **Garantir retorno do `pb-saved-menu:<sessionID>`**:
   - Garantir que a chamada a `showSavedPostMenu` resgate o `sessionID` salvo no Redis e edite a mensagem de volta para o menu da postagem com o token `@BotUsername pb <sessionID>`.

---

## Passos detalhados

### Passo 1: Atualizar o handler de `pb-schedule:` em `postBuilder.go`
- Obter os canais do usuário no momento do callback `pb-schedule:<sessionID>`.
- Se `len(channels) == 0`:
  - Disparar `bot.AnswerCallbackQuery` com `ShowAlert: true` e a mensagem: `"⚠️ Você não possui nenhum canal cadastrado para usar esta funcionalidade!"`.
  - Retornar imediatamente `nil` (sem editar nem enviar nada).

### Passo 2: Ajustar `handleScheduleTypeAction`
- Remover o botão `❌ Cancelar` do prompt de data/hora.
- Deixar apenas os botões `🔙 Voltar` e `🔙 Voltar ao Post`.

### Passo 3: Testar e Validar
- Compilar (`go build ./cmd/FreddyBot`) e rodar testes (`go test ./...`).

---

## Riscos
- Nenhum.

## Impactos esperados
- O menu do post permanece 100% intacto quando o usuário não possui canais.
- Surge um popup nativo do Telegram na tela do usuário.
- Navegação limpa e sem botão "Cancelar".

## Compatibilidade
- Linux / macOS / Windows / Docker / Telego

## Como testar
1. Com 0 canais, clicar em "Agendar Envio" -> Verificar popup de alerta na tela sem alteração da mensagem.
2. Cadastrar canal, clicar em "Agendar Envio" -> Navegar e testar os botões `🔙 Voltar` e `🔙 Voltar ao Post`.
