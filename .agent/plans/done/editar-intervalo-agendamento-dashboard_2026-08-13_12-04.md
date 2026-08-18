# Plano: Permitir Edição de Intervalo e Janela de Horário no Dashboard

## Pedido do usuário
Permitir que o usuário edite o intervalo (em minutos) e a janela de horário (início/fim) de um agendamento do tipo `interval` diretamente pela modal de edição no Dashboard.

## Objetivo
Adicionar suporte completo no backend (API DTOs, Repository e Controller) e no frontend (modal de edição no `ScheduleTab.tsx`) para alterar `intervalMin`, `windowStart` e `windowEnd`.

## Contexto atual
- O backend já salva `IntervalMin`, `WindowStart`, `WindowEnd` ao criar agendamentos pelo Telegram.
- A modal de edição no Dashboard (`ScheduleTab.tsx`) atualmente só exibe inputs para "Data" e "Horário".
- O endpoint da API (`EditScheduleRequest`) e o controller ainda não aceitam `intervalMin`, `windowStart`, `windowEnd`.

## Arquivos analisados
- `internal/api/types/scheduler.go`
- `internal/api/controllers/schedulerController.go`
- `internal/core/services/scheduler.go`
- `internal/database/repositories/scheduled_post.go`
- `dashboard/src/api.ts`
- `dashboard/src/components/ScheduleTab.tsx`

## Arquivos que poderão ser modificados
- `internal/api/types/scheduler.go`
- `internal/api/controllers/schedulerController.go`
- `internal/core/services/scheduler.go`
- `internal/database/repositories/scheduled_post.go`
- `dashboard/src/api.ts`
- `dashboard/src/components/ScheduleTab.tsx`

## Estratégia de implementação

1. **Backend DTOs (`internal/api/types/scheduler.go`)**:
   - Atualizar `EditScheduleRequest` para incluir `IntervalMin *int`, `WindowStart *string`, `WindowEnd *string`.
   - Atualizar `CreateScheduleRequest` validação `oneof` para incluir `interval` e adicionar os 3 campos.

2. **Repository (`internal/database/repositories/scheduled_post.go`)**:
   - Adicionar método `UpdateScheduleInterval(ctx, id, intervalMin, windowStart, windowEnd)`.

3. **Service (`internal/core/services/scheduler.go`)**:
   - Adicionar o método `UpdateScheduleInterval(ctx, id, ownerID, intervalMin, windowStart, windowEnd)` no `SchedulerService` com validações.

4. **Controller (`internal/api/controllers/schedulerController.go`)**:
   - Em `UpdateSchedule`, invocar `UpdateScheduleInterval` se algum dos campos de intervalo for enviado.

5. **Dashboard API Client (`dashboard/src/api.ts`)**:
   - Atualizar a função `updateScheduleTime` para aceitar `intervalMin`, `windowStart`, `windowEnd`.

6. **Dashboard UI (`dashboard/src/components/ScheduleTab.tsx`)**:
   - Na modal de edição:
     - Se `schedule.scheduleType === 'interval'`: exibir campos de **Intervalo (minutos)**, **Janela Início** (time/HH:MM) e **Janela Fim** (time/HH:MM).
     - Ao salvar, enviar os valores editados na requisição HTTP.

7. **Validação e Compilação**:
   - Rodar `go build ./cmd/FreddyBot` e `npm run build` na pasta `dashboard`.

## Passos detalhados
1. Atualizar DTOs em `internal/api/types/scheduler.go`.
2. Adicionar método em `internal/database/repositories/scheduled_post.go`.
3. Adicionar método em `internal/core/services/scheduler.go`.
4. Atualizar controller em `internal/api/controllers/schedulerController.go`.
5. Atualizar cliente HTTP em `dashboard/src/api.ts`.
6. Atualizar modal em `dashboard/src/components/ScheduleTab.tsx`.
7. Testar builds backend e frontend.

## Riscos
- **Baixo**: Alteração retrocompatível em endpoints existentes.

## Impactos esperados
- Usuários conseguirão ajustar o tempo do intervalo e a janela de horário de qualquer agendamento `interval` pela Dashboard.

## Compatibilidade
- Linux, macOS, Windows, Docker, CI/CD

## Como testar
```bash
go build ./cmd/FreddyBot
cd dashboard && npm run build
```

## Rollback
```bash
git checkout internal/api/types/scheduler.go
git checkout internal/api/controllers/schedulerController.go
git checkout internal/core/services/scheduler.go
git checkout internal/database/repositories/scheduled_post.go
git checkout dashboard/src/api.ts
git checkout dashboard/src/components/ScheduleTab.tsx
```
