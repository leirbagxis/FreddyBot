# Plano: Corrigir Validação de Horários com 1 Dígito na Hora (ex: 8:59)

## Pedido do usuário
O usuário relatou o erro:
`❌ Erro ao criar agendamento: horário de início da janela inválido`
ao enviar o formato:
`5`
`8:59-12:00`

## Objetivo
Permitir que horários com 1 dígito na hora (ex: `8:59`, `9:00`) sejam aceitos e formatados automaticamente para `08:59`, `09:00` tanto na validação do backend quanto na criação do agendamento.

## Contexto atual
- `validateScheduleTime` exige estritamente `len(parts[0]) == 2` e `len(parts[1]) == 2`.
- Quando o usuário digita `8:59`, `len("8") == 1`, causando a rejeição do início da janela.

## Arquivos analisados
- `internal/core/services/scheduler.go` (função `validateScheduleTime` e normalização)
- `internal/telegram/handlers/events/postBuilder/postBuilder.go` (normalização de `windowStart` e `windowEnd`)

## Arquivos que poderão ser modificados
- `internal/core/services/scheduler.go`
- `internal/telegram/handlers/events/postBuilder/postBuilder.go`

## Estratégia de implementação
1. Atualizar `validateScheduleTime` em `scheduler.go`:
   - Permitir `len(parts[0])` de 1 ou 2 dígitos (0 a 23).
   - Validar se hora está entre 0 e 23 e minutos entre 0 e 59.
   - Retornar os inteiros válidos de hora e minuto.
2. Atualizar `formatScheduleTime` / normalização em `scheduler.go` e `postBuilder.go` para formatar `windowStart` e `windowEnd` como `HH:MM` (ex: `8:59` -> `08:59`).
3. Compilar o binário em Go (`go build ./cmd/FreddyBot`) para verificar integridade.

## Passos detalhados
1. Editar `validateScheduleTime` em `internal/core/services/scheduler.go`.
2. Normalizar `opts.WindowStart` e `opts.WindowEnd` para `08:59` em `CreateScheduledPost` ou `postBuilder.go`.
3. Executar `go build ./cmd/FreddyBot`.

## Riscos
- **Nenhum**: Torna a entrada de dados mais flexível e tolerante sem alterar o formato interno de armazenamento (`HH:MM`).

## Impactos esperados
- Horários como `8:59`, `9:30` serão aceitos perfeitamente no Telegram e normalizados para `08:59`, `09:30`.

## Compatibilidade
- Linux, macOS, Windows, Docker, CI/CD

## Como testar
```bash
go build ./cmd/FreddyBot
```

## Rollback
```bash
git checkout internal/core/services/scheduler.go
git checkout internal/telegram/handlers/events/postBuilder/postBuilder.go
```
