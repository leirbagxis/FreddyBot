# Plano: Corrigir Navegação e Interatividade no Agendamento do PostBuilder

## Pedido do usuário
Na feature de agendar posts no PostBuilder, ao clicar para agendar, não é possível voltar depois para o menu que contém o token do PostBuilder (`@BotUsername pb <id>`), e o menu durante/após o agendamento deixa de ser interativo.

## Objetivo
Tornar o fluxo de agendamento do PostBuilder totalmente navegável e interativo, permitindo que o usuário alterne entre as etapas de agendamento, cancele/volte ao menu da postagem salva (`@BotUsername pb <id>`), e receba um menu com atalhos de retorno após a criação do agendamento.

## Contexto atual
- `pb-save` gera o menu da postagem salva com o token inline (`@BotUsername pb <id>`) e botões (`🚀 Compartilhar`, `📢 Enviar para Canais`, `📅 Agendar Envio`, `💾 Salvar como Template`).
- Ao clicar em `📅 Agendar Envio` (`pb-schedule:<id>`), o bot faz `EditMessageText` sobrescrevendo a mensagem do token pela seleção de canal.
- Nas telas de agendamento, o botão `❌ Cancelar` dispara `pb-cancel` que exclui todo o estado e encerra o PostBuilder, em vez de voltar ao menu do post salvo.
- Na tela de prompt de texto (data/hora), a mensagem é editada sem nenhum teclado inline (`ReplyMarkup: nil`), tornando a tela estática.
- Após o envio do texto e criação do agendamento em `handleScheduleTextInput`, o bot envia apenas o texto `"✅ Agendamento criado!"` sem botões de navegação.

## Arquivos analisados
- `internal/telegram/handlers/events/postBuilder/postBuilder.go`
- `internal/telegram/loader_telego.go`

## Arquivos que poderão ser modificados
- `internal/telegram/handlers/events/postBuilder/postBuilder.go`
- `internal/telegram/loader_telego.go` (caso seja necessário novo callback handler)

## Estratégia de implementação

1. **Criar a função `showSavedPostMenu`**:
   - Função utilitária centralizada que renderiza/edita a mensagem do post salvo (mostrando o token `@BotUsername pb <id>` e os botões de ação: Compartilhar, Enviar para Canais, Agendar Envio, Salvar como Template, Menu Principal).

2. **Adicionar callback `pb-saved-menu:<sessionID>`**:
   - Permite que qualquer sub-fluxo (agendamento, envio para canais, templates) retorne ao menu da postagem salva sem perder o estado nem o token.

3. **Inserir botão `🔙 Voltar ao Post` nas etapas de agendamento**:
   - Em `handleSchedulePost` (Seleção de Canal): substituir `❌ Cancelar` por `🔙 Voltar ao Post` (`pb-saved-menu:<sessionID>`).
   - Em `handleScheduleTypeSelection` (Tipo de Agendamento): adicionar `🔙 Voltar` (`pb-schedule:<sessionID>`) e `🔙 Voltar ao Post` (`pb-saved-menu:<sessionID>`).
   - Em `handleScheduleTypeAction` (Prompt de Data/Hora): incluir teclado inline interativo com botões `🔙 Voltar` (`pb-sch:<sessionID>:<channelID>`) e `❌ Cancelar` (`pb-saved-menu:<sessionID>`).

4. **Tornar o encerramento do agendamento interativo**:
   - Em `handleScheduleTextInput` (ao concluir o agendamento com sucesso), incluir teclado inline na mensagem de sucesso com botões:
     - `🔙 Voltar ao Post` (`pb-saved-menu:<sessionID>`)
     - `📋 Meus Agendamentos` (`my-schedules`)

---

## Passos detalhados

### Passo 1: Criar `showSavedPostMenu` e o handler `pb-saved-menu:<id>`
Em `postBuilder.go`:
- Criar a função `showSavedPostMenu(ctx, chatID, userID, messageID, sessionID, c)` que carrega a sessão pelo `sessionID` e renderiza/edita a mensagem do menu do post salvo (com o token `@BotUsername pb <id>`).
- Registrar/tratar o prefixo de callback `pb-saved-menu:` para chamar essa função.

### Passo 2: Atualizar `handleSchedulePost` (Seleção de Canal)
- Trocar o botão `❌ Cancelar` (`pb-cancel`) por `🔙 Voltar ao Post` (`pb-saved-menu:` + sessionID).

### Passo 3: Atualizar `handleScheduleTypeSelection` (Seleção do Tipo de Agendamento)
- Adicionar botão `🔙 Voltar` (`pb-schedule:` + sessionID) para retornar à seleção de canal.
- Adicionar botão `🔙 Voltar ao Post` (`pb-saved-menu:` + sessionID).

### Passo 4: Atualizar `handleScheduleTypeAction` (Prompt de Entrada de Texto)
- Adicionar teclado inline ao `EditMessageText` com os botões:
  - `🔙 Voltar` (`pb-sch:` + sessionID + `:` + channelID)
  - `❌ Cancelar` (`pb-saved-menu:` + sessionID)

### Passo 5: Atualizar `handleScheduleTextInput` (Confirmação de Agendamento)
- Ao enviar a mensagem de sucesso (`✅ Agendamento criado!`), anexar `ReplyMarkup` com os botões:
  - `🔙 Voltar ao Post` (`pb-saved-menu:` + sessionID)
  - `📋 Meus Agendamentos` (`my-schedules`)

---

## Riscos
- **Sessão Expirada**: Se o token no Redis expirar (após TTL), o menu exibe mensagem tratada de sessão expirada.
- **Formato da CallbackData**: Respeitar o limite de 64 bytes da Telegram Bot API para callback_data (`pb-saved-menu:<sessionID>`).

## Impactos esperados
- Fluxo de agendamento 100% interativo e navegável.
- O usuário nunca fica preso em uma tela estática sem botões.
- É possível voltar ao menu da postagem salva com seu token inline a qualquer momento.

## Compatibilidade
- Linux / macOS / Windows / Docker
- Telegram Bot API (telego)

## Como testar

### Build
```bash
go build -o FreddyBot ./cmd/FreddyBot/main.go
```

### Execução
```bash
./FreddyBot
```

### Teste Manual no Telegram
1. Enviar mídia/texto no privado do bot para abrir o PostBuilder.
2. Clicar em `✅ Salvar` -> Verificar envio da mensagem com token e menu de ações.
3. Clicar em `📅 Agendar Envio` -> Selecionar canal -> Selecionar tipo de agendamento.
4. Verificar presença do botão `🔙 Voltar` nas telas intermediárias.
5. Na tela de digitação de data/hora, verificar se os botões inline `🔙 Voltar` e `❌ Cancelar` estão visíveis.
6. Digitar a data/hora -> Verificar se a mensagem de confirmação final traz o botão `🔙 Voltar ao Post`.
7. Clicar em `🔙 Voltar ao Post` -> Verificar se a mensagem retorna ao card do post salvo com o token inline.

## Rollback
Desfazer as alterações em `postBuilder.go` via `git checkout`.
