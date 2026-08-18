# Plano: corrigir-visibilidade-autodestruicao-e-rascunhos-dashboard-bot_2026-08-15_21-05

## Pedido do usuário
Garantir que a **Auto-Destruição da Mensagem** e o **Salvamento e Carregamento de Rascunhos** fiquem perfeitamente visíveis, funcionais e acessíveis tanto na Dashboard Web (React) quanto no Telegram Bot.

## Objetivo
1. Adicionar o seletor visual de **Auto-Destruição** (`autoDeleteMin`) no modal de criação e edição de agendamentos no Dashboard Web (`ScheduleTab.tsx`).
2. Adicionar o gerenciador e carregador visual de **Rascunhos de Posts** no Dashboard Web (`UserTemplatesManager.tsx` / `ScheduleTab.tsx`).
3. No Telegram Bot, simplificar o salvamento de rascunhos para gerar um nome padrão automático (ex: `Rascunho - DD/MM 15:04` ou Título) com confirmação instantânea de salvamento e exibir a seleção de Auto-Destruição de forma clara no agendamento.

## Contexto atual
- O backend de `AutoDeleteService` e `PostTemplateService` já possui todos os endpoints e métodos no banco e repositórios.
- Na Dashboard Web (`ScheduleTab.tsx`), faltava renderizar o `<select>` de Auto-Destruição no formulário de agendamento e o modal de Rascunhos Salvos.
- No Telegram Bot (`postBuilder.go`), ao clicar em `Salvar como Rascunho`, o bot solicitava o nome manualmente; adicionando geração automática de nome de fallback, o salvamento torna-se 100% instantâneo.

## Arquivos analisados
- `dashboard/src/components/ScheduleTab.tsx`
- `dashboard/src/components/UserTemplatesManager.tsx`
- `dashboard/src/api.ts`
- `dashboard/src/types.ts`
- `internal/telegram/handlers/events/postBuilder/postBuilder.go`
- `internal/telegram/handlers/callbacks/my_schedules/my_schedules.go`

## Arquivos que poderão ser modificados
- `dashboard/src/components/ScheduleTab.tsx`
- `dashboard/src/components/UserTemplatesManager.tsx`
- `internal/telegram/handlers/events/postBuilder/postBuilder.go`

## Estratégia de implementação

1. **Dashboard Web (`ScheduleTab.tsx`)**:
   - No modal de agendamento e edição de agendamento, adicionar o campo de seleção de **Auto-Destruição da Mensagem**:
     - `Desativada` (0 min)
     - `1 hora` (60 min)
     - `6 horas` (360 min)
     - `12 horas` (720 min)
     - `24 horas` (1440 min)
     - `48 horas` (2880 min)
   - Na lista/tabela de agendamentos do Dashboard, exibir o selo `⏱️ Auto-destruição em Xh` para agendamentos temporários.
   - Adicionar o botão/modal **📂 Meus Rascunhos** na Dashboard Web para permitir visualizar, carregar para o agendador e excluir rascunhos salvos.

2. **Telegram Bot (`postBuilder.go`)**:
   - No callback `pb-save-template`, se o usuário não enviar nome ou clicar no botão direto, usar automaticamente o Título do post ou `Rascunho - DD/MM às HH:MM` e salvar na biblioteca com feedback instantâneo!
   - Na tela de agendamento no Telegram Bot, renderizar o botão interativo de Auto-Destruição de forma destacada e intuitiva.

3. **Validação & Testes**:
   - Executar os testes do Go (`go test ./cmd/... ./internal/... ./pkg/...`).
   - Compilar o binário Go (`go build ./cmd/FreddyBot`).
   - Executar o build do Dashboard Frontend (`cd dashboard && npm run build`).

## Passos detalhados
1. Salvar o plano em `.agent/plans/pending/corrigir-visibilidade-autodestruicao-e-rascunhos-dashboard-bot_2026-08-15_21-05.md`.
2. Apresentar o resumo ao usuário e solicitar aprovação explícita.
3. Atualizar `ScheduleTab.tsx` e `UserTemplatesManager.tsx` na Dashboard Web.
4. Atualizar `postBuilder.go` no Telegram Bot.
5. Rodar os testes e compilações (`go test ./...`, `go build ./cmd/FreddyBot`, `npm run build`).

## Riscos
- Nulo. Alteração restrita à exibição dos componentes de UI no Dashboard e simplificação de callback no Bot.

## Compatibilidade
- Linux, macOS, Windows, Docker, CI/CD.

## Como testar

### Testes Go & Build Dashboard
```bash
go test ./cmd/... ./internal/... ./pkg/...
go build ./cmd/FreddyBot
cd dashboard && npm run build
```

## Rollback
`git checkout dashboard/src/components/ScheduleTab.tsx internal/telegram/handlers/events/postBuilder/postBuilder.go`
