# Plano: configurar-autodestruicao-no-fluxo-agendamento_2026-08-15_20-56

## Pedido do usuário
Inserir a **Auto-Destruição da Mensagem** como um **parâmetro interno da tela de configuração do Agendamento** (e não como um botão solto no criador de posts), permitindo que o usuário escolha o tempo de expiração ao configurar o horário/canal do agendamento. Na mesma tela de agendamento, oferecer a opção de **Salvar como Rascunho** na biblioteca pessoal.

## Objetivo
Reorganizar a UX do Telegram Bot e do Dashboard para que o parâmetro ⏱️ **Auto-Destruição da Mensagem** (`Desativada / 1h / 6h / 12h / 24h / 48h`) seja ajustado dentro do fluxo de confirmação e configuração do agendamento.

## Contexto atual
- O usuário aciona `📅 Agendar Envio`, seleciona o canal e define a frequência (Pontual, Diário, Semanal, Intervalo).
- Atualmente, a tela de resumo/confirmação do agendamento exibe o canal, tipo e horário.
- Inserindo o seletor de Auto-Destruição e o botão de Salvar Rascunho nesta tela de agendamento, o fluxo fica totalmente fluido, organizado e contextual.

## Arquivos analisados
- `internal/telegram/handlers/events/postBuilder/postBuilder.go`
- `internal/telegram/handlers/callbacks/my_schedules/my_schedules.go`
- `internal/api/controllers/schedulerController.go`
- `dashboard/src/api.ts`
- `dashboard/src/types.ts`

## Arquivos que poderão ser modificados
- `internal/telegram/handlers/events/postBuilder/postBuilder.go`
- `internal/telegram/handlers/callbacks/my_schedules/my_schedules.go`

## Estratégia de implementação

1. **Remoção do Botão Solto no PostBuilder**:
   - Remover a linha isolada de Auto-Destruição do menu principal de criação de post em `showMenuTelego` (mantendo o menu de criação limpo e focado no conteúdo).

2. **Inclusão na Tela de Configuração do Agendamento**:
   - Na tela de confirmação/resumo do agendamento (`handleScheduleTypeAction` / `handleScheduleConfirm`), incluir o seletor interativo:
     - ⏱️ **Auto-Destruição:** `[Desativada / 1h / 6h / 12h / 24h / 48h]` (`pb-sch-autodel:<sessionID>`)
   - Ao clicar no botão de Auto-Destruição durante o agendamento, alterna o tempo (0, 60, 360, 720, 1440, 2880 min) na sessão do agendamento e atualiza o resumo.

3. **Inclusão da Opção de Salvar como Rascunho**:
   - Incluir o botão `💾 Salvar Rascunho` na tela de agendamento e na tela de ações finais da postagem pronta.

4. **Validação & Testes**:
   - Executar suíte de testes Go (`go test ./cmd/... ./internal/... ./pkg/...`).
   - Compilar o binário (`go build ./cmd/FreddyBot`).

## Passos detalhados
1. Salvar o plano em `.agent/plans/pending/configurar-autodestruicao-no-fluxo-agendamento_2026-08-15_20-56.md`.
2. Apresentar o resumo ao usuário e solicitar aprovação explícita.
3. Atualizar `postBuilder.go` para integrar a alternância de Auto-Destruição ao resumo do agendamento.
4. Executar os testes e compilações (`go test ./...`, `go build ./cmd/FreddyBot`).

## Riscos
- Mínimo. Alteração restrita ao fluxo visual e de botões do agendamento no Telegram Bot.

## Compatibilidade
- Linux, macOS, Windows, Docker, CI/CD.

## Como testar

### Testes Go
```bash
go test ./cmd/... ./internal/... ./pkg/...
go build ./cmd/FreddyBot
```

## Rollback
`git checkout internal/telegram/handlers/events/postBuilder/postBuilder.go`
