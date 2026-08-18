# Plano: corrigir-detalhes-agendamento-telegram_2026-08-15_00-30

## Pedido do usuário
1. Corrigir o botão dos agendamentos no bot do Telegram que atualmente envia o callback `schedule-detail:...` mas não responde ao ser clicado.
2. Reformular as mensagens dos agendamentos no Telegram para que fiquem mais detalhadas, amigáveis e totalmente em português (substituindo termos como "pending" por "Ativo" e formatando horários e tipos em PT-BR).
3. Adicionar telas de detalhes do agendamento com opções para Pausar, Retomar e Excluir o agendamento diretamente pelo bot.

## Objetivo
Implementar o handler de callback `schedule-detail:`, implementar ações de gerenciamento de agendamento (Pausar, Retomar, Excluir) e redesenhar o menu de agendamentos no Telegram com textos ricos, traduzidos para o Português e com fuso horário do Brasil (`BRT`).

## Contexto atual
- `my_schedules.go` exibe uma lista básica onde o status aparece em inglês (ex: `pending`) e renderiza botões com CallbackData `schedule-detail:<ID>`.
- `loader_telego.go` não possuía nenhum handler registrado para tratar o prefixo `schedule-detail:`, fazendo com que clicar no botão não disparasse nenhuma ação.
- O serviço `SchedulerService` já possui os métodos de backend necessários: `GetScheduleByID`, `PauseScheduledPost`, `ResumeScheduledPost` e `DeleteScheduledPost`.

## Arquivos analisados
- `internal/telegram/handlers/callbacks/my_schedules/my_schedules.go`
- `internal/telegram/loader_telego.go`
- `internal/core/services/scheduler.go`
- `internal/database/models/scheduled_post.go`
- `internal/utils/utils.go`

## Arquivos que poderão ser modificados
- `internal/telegram/handlers/callbacks/my_schedules/my_schedules.go`
- `internal/telegram/loader_telego.go`

## Estratégia de implementação

1. **Redesenho da Lista (`my_schedules.go` - `HandlerTelego`)**:
   - Traduzir os status para Português:
     - `pending` $\rightarrow$ `🟢 Ativo`
     - `paused` $\rightarrow$ `🟡 Pausado`
     - `sent` / `completed` $\rightarrow$ `✅ Concluído`
     - `failed` / `cancelled` $\rightarrow$ `🔴 Cancelado / Falhou`
   - Formatar o tipo de agendamento em Português:
     - `once` $\rightarrow$ `Envio Único`
     - `daily` $\rightarrow$ `Diário`
     - `weekly` $\rightarrow$ `Semanal`
     - `interval` $\rightarrow$ `Intervalo (<N> min)`
     - `queue` $\rightarrow$ `Fila`
   - Formatar horários e janelas com clareza.
   - Renderizar o texto em HTML elegante com emojis descritivos.

2. **Implementação da Tela de Detalhes (`DetailHandlerTelego`)**:
   - Capturar o ID do agendamento a partir do callback `schedule-detail:<ID>`.
   - Buscar os dados do agendamento via `c.SchedulerService.GetScheduleByID`.
   - Formatar as informações detalhadas em Português:
     - Canal e ID
     - Tipo de Agendamento
     - Horário / Janela / Frequência
     - Próxima execução formatada no fuso horário do Brasil (`utils.BrazilTZ()`)
     - Opção de Fixar Mensagem (Sim/Não)
     - Total de Envios Realizados
     - Status atual traduzido
     - Erro recente (se houver)
   - Adicionar botões interativos de ação:
     - Se `pending`: Botão `🟡 Pausar Agendamento` (`schedule-pause:<ID>`)
     - Se `paused`: Botão `🟢 Retomar Agendamento` (`schedule-resume:<ID>`)
     - Botão `🗑️ Excluir Agendamento` (`schedule-delete:<ID>`)
     - Botão `🔙 Voltar para Lista` (`my-schedules`)

3. **Implementação dos Callbacks de Ação**:
   - `PauseHandlerTelego`: invoca `c.SchedulerService.PauseScheduledPost`, exibe alerta de confirmação no Telegram ("🟡 Agendamento pausado com sucesso!") e atualiza a tela de detalhes.
   - `ResumeHandlerTelego`: invoca `c.SchedulerService.ResumeScheduledPost`, exibe alerta de confirmação ("🟢 Agendamento ativado com sucesso!") e atualiza a tela de detalhes.
   - `DeleteHandlerTelego`: invoca `c.SchedulerService.DeleteScheduledPost`, exibe alerta de confirmação ("🗑️ Agendamento excluído!") e retorna para a lista de agendamentos (`HandlerTelego`).

4. **Registro no Router Telego (`loader_telego.go`)**:
   - Registrar no `BotHandler`:
     ```go
     bh.Handle(callbackMySchedules.DetailHandlerTelego(c), telegohandler.CallbackDataPrefix("schedule-detail:"))
     bh.Handle(callbackMySchedules.PauseHandlerTelego(c), telegohandler.CallbackDataPrefix("schedule-pause:"))
     bh.Handle(callbackMySchedules.ResumeHandlerTelego(c), telegohandler.CallbackDataPrefix("schedule-resume:"))
     bh.Handle(callbackMySchedules.DeleteHandlerTelego(c), telegohandler.CallbackDataPrefix("schedule-delete:"))
     ```

## Passos detalhados
1. Criar o arquivo de plano em `.agent/plans/pending/corrigir-detalhes-agendamento-telegram_2026-08-15_00-30.md`.
2. Obter aprovação explícita do usuário.
3. Atualizar `internal/telegram/handlers/callbacks/my_schedules/my_schedules.go` com a lista aprimorada e os novos handlers de detalhe, pausa, retomada e exclusão.
4. Atualizar `internal/telegram/loader_telego.go` para registrar as novas rotas de callback.
5. Executar os testes unitários (`go test ./...`) e compilação (`go build ./cmd/FreddyBot`).

## Riscos
- Nulo. As alterações são restritas aos handlers de callback do bot do Telegram de agendamentos.

## Compatibilidade
- Linux, macOS, Windows, Docker, CI/CD.

## Como testar

### Build Backend
```bash
go build ./cmd/FreddyBot
```

### Testes Go
```bash
go test ./internal/telegram/...
```

## Rollback
`git checkout internal/telegram/handlers/callbacks/my_schedules/my_schedules.go internal/telegram/loader_telego.go`
