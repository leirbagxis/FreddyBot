# Plano: Implementar Agendamento por Intervalo (`interval`)

## Pedido do usuário
Permitir que usuários agendem posts com intervalos menores que 24 horas (ex: a cada 30 minutos, a cada 1 hora).

## Objetivo
Adicionar um novo tipo de agendamento `interval` ao sistema existente, com intervalo configurável em minutos e janela de horário diária opcional.

## Contexto atual
- 4 tipos de agendamento existentes: `once`, `daily`, `weekly`, `queue`
- Todos os recorrentes avançam no mínimo +1 dia
- O ticker do scheduler já roda a cada 30s (suporta granularidade sub-minuto)
- O campo `NextRunAt` já aceita qualquer timestamp UTC

## Arquivos analisados
- `internal/database/models/scheduled_post.go` (modelo)
- `internal/core/services/scheduler.go` (serviço + cálculos de recorrência)
- `internal/telegram/handlers/events/postBuilder/postBuilder.go` (UI Telegram)
- `internal/telegram/handlers/callbacks/my_schedules/my_schedules.go` (listagem)
- `internal/cache/types.go` (cache state)
- `dashboard/src/components/ScheduleTab.tsx` (dashboard UI)

## Arquivos que poderão ser modificados
- `internal/database/models/scheduled_post.go`
- `internal/core/services/scheduler.go`
- `internal/telegram/handlers/events/postBuilder/postBuilder.go`
- `internal/telegram/handlers/callbacks/my_schedules/my_schedules.go`
- `dashboard/src/components/ScheduleTab.tsx`
- `dashboard/src/types.ts` (tipo ScheduledPost no frontend)

## Estratégia de implementação
Alteração cirúrgica e retrocompatível: 3 novos campos opcionais no modelo, 1 nova função de cálculo, novos cases nos switches existentes, novo botão e prompt no Telegram UI.

## Passos detalhados

### Passo 1 — Modelo (`internal/database/models/scheduled_post.go`)

Adicionar 3 campos ao struct `ScheduledPost`:

```go
IntervalMin  int    `gorm:"default:0" json:"intervalMin"`   // Intervalo em minutos (0 = desabilitado)
WindowStart  string `json:"windowStart"`                     // "HH:MM" início da janela diária (opcional)
WindowEnd    string `json:"windowEnd"`                       // "HH:MM" fim da janela diária (opcional)
```

Inserir após a linha do campo `RepeatUntil` (~linha 19). A migration será automática via GORM AutoMigrate.

### Passo 2 — ScheduleOptions (`internal/core/services/scheduler.go`)

Adicionar ao struct `ScheduleOptions` (~linha 25):

```go
IntervalMin  int
WindowStart  string
WindowEnd    string
```

### Passo 3 — Nova função `calculateNextInterval` (`internal/core/services/scheduler.go`)

Adicionar após `calculateNextWeekly` (~após linha 323):

```go
func (s *SchedulerService) calculateNextInterval(post *models.ScheduledPost, sentAt time.Time) *time.Time {
	if post.IntervalMin < 5 {
		post.IntervalMin = 5 // mínimo de segurança
	}
	next := sentAt.Add(time.Duration(post.IntervalMin) * time.Minute)

	// Se tem janela de horário configurada, respeitar limites
	if post.WindowStart != "" && post.WindowEnd != "" {
		brazilTZ := utils.BrazilTZ()
		nextLocal := next.In(brazilTZ)

		startH, startM := parseHHMM(post.WindowStart)
		endH, endM := parseHHMM(post.WindowEnd)

		windowStart := time.Date(nextLocal.Year(), nextLocal.Month(), nextLocal.Day(), startH, startM, 0, 0, brazilTZ)
		windowEnd := time.Date(nextLocal.Year(), nextLocal.Month(), nextLocal.Day(), endH, endM, 0, 0, brazilTZ)

		if nextLocal.Before(windowStart) {
			// Ainda não abriu a janela hoje
			next = windowStart
		} else if nextLocal.After(windowEnd) {
			// Passou da janela hoje, agendar para amanhã no início da janela
			next = windowStart.AddDate(0, 0, 1)
		}

		next = next.UTC()
	}

	return &next
}
```

### Passo 4 — Case `interval` no switch de `sendScheduledPost` (`internal/core/services/scheduler.go`)

No switch em `sendScheduledPost` (~linha 142), adicionar antes do case `"queue"`:

```go
case "interval":
    nextRun := s.calculateNextInterval(post, sentAt)
    if nextRun != nil && (post.RepeatUntil == nil || nextRun.Before(*post.RepeatUntil)) {
        s.repo.UpdateStatus(ctx, post.ID, "pending")
        s.repo.UpdateNextRunAt(ctx, post.ID, *nextRun)
        logger.Info("SCHEDULER", "Post intervalo %s: próximo envio %s (cada %dmin)", post.ID, nextRun.Format("02/01 15:04"), post.IntervalMin)
    } else {
        logger.Info("SCHEDULER", "Post intervalo %s finalizado (repeat_until atingido)", post.ID)
    }
```

### Passo 5 — Cálculo do `NextRunAt` inicial na criação (`internal/core/services/scheduler.go`)

Na função `CreateScheduledPost` (~linha 376), adicionar um bloco para `interval` no cálculo de `nextRunAt`:

```go
} else if opts.ScheduleType == "interval" {
    if opts.WindowStart != "" {
        brazilTZ := utils.BrazilTZ()
        now := time.Now().In(brazilTZ)
        startH, startM := parseHHMM(opts.WindowStart)
        windowStart := time.Date(now.Year(), now.Month(), now.Day(), startH, startM, 0, 0, brazilTZ)
        if now.Before(windowStart) {
            nextRunAt = windowStart.UTC()
        } else {
            // Já estamos dentro da janela, primeiro envio em IntervalMin minutos
            nextRunAt = time.Now().Add(time.Duration(opts.IntervalMin) * time.Minute).UTC()
        }
    } else {
        // Sem janela: primeiro envio em IntervalMin minutos
        nextRunAt = time.Now().Add(time.Duration(opts.IntervalMin) * time.Minute).UTC()
    }
```

E passar os novos campos ao struct do post sendo criado:

```go
IntervalMin:  opts.IntervalMin,
WindowStart:  opts.WindowStart,
WindowEnd:    opts.WindowEnd,
```

### Passo 6 — Validação na criação (`internal/core/services/scheduler.go`)

Adicionar validação no início de `CreateScheduledPost`, junto com as validações existentes:

```go
if opts.ScheduleType == "interval" {
    if opts.IntervalMin < 5 {
        return nil, apperrors.BadRequest("intervalo mínimo é 5 minutos")
    }
    if opts.IntervalMin > 1440 {
        return nil, apperrors.BadRequest("para intervalos de 24h ou mais, use o tipo 'daily'")
    }
    if (opts.WindowStart != "" && opts.WindowEnd == "") || (opts.WindowStart == "" && opts.WindowEnd != "") {
        return nil, apperrors.BadRequest("informe início e fim da janela de horário, ou deixe ambos vazios")
    }
    if opts.WindowStart != "" {
        if _, _, err := validateScheduleTime(opts.WindowStart); err != nil {
            return nil, apperrors.BadRequest("horário de início da janela inválido")
        }
        if _, _, err := validateScheduleTime(opts.WindowEnd); err != nil {
            return nil, apperrors.BadRequest("horário de fim da janela inválido")
        }
    }
}
```

### Passo 7 — Botão no Telegram (`internal/telegram/handlers/events/postBuilder/postBuilder.go`)

Na função `handleScheduleTypeSelection` (~linha 1793), adicionar novo botão antes de "Fila de envio":

```go
{
    {Text: "⏱️ Intervalo fixo", CallbackData: "pb-sch-type:" + sessionID + ":" + channelID + ":interval"},
},
```

### Passo 8 — Prompt no Telegram (`internal/telegram/handlers/events/postBuilder/postBuilder.go`)

Na função `handleScheduleTypeAction` (~linha 1836), adicionar case no switch:

```go
case "interval":
    prompt = "⏱️ <b>Intervalo Fixo</b>\n\nEnvie o intervalo em minutos e, opcionalmente, a janela de horário:\n\n<b>Linha 1:</b> intervalo em minutos\n<b>Linha 2:</b> horário início-fim (opcional)\n\nExemplos:\n<code>30</code>\n→ A cada 30 minutos, 24h/dia\n\n<code>60\n08:00-22:00</code>\n→ A cada 1 hora, das 08:00 às 22:00"
```

### Passo 9 — Parsing de input no Telegram (`internal/telegram/handlers/events/postBuilder/postBuilder.go`)

Na função `handleScheduleTextInput` (~linha 1885), adicionar case `"interval"` antes do fechamento do switch:

```go
case "interval":
    lines := strings.Split(text, "\n")
    intervalMin, err := strconv.Atoi(strings.TrimSpace(lines[0]))
    if err != nil || intervalMin < 5 {
        _, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
            ChatID:    telego.ChatID{ID: chatID},
            Text:      "❌ Intervalo inválido. Envie um número de minutos (mínimo 5).",
            ParseMode: telego.ModeHTML,
        })
        return
    }

    var windowStart, windowEnd string
    if len(lines) >= 2 {
        windowParts := strings.Split(strings.TrimSpace(lines[1]), "-")
        if len(windowParts) == 2 {
            windowStart = strings.TrimSpace(windowParts[0])
            windowEnd = strings.TrimSpace(windowParts[1])
        } else {
            _, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
                ChatID:    telego.ChatID{ID: chatID},
                Text:      "❌ Janela de horário inválida. Use o formato:\n<code>HH:MM-HH:MM</code>",
                ParseMode: telego.ModeHTML,
            })
            return
        }
    }

    opts := services.ScheduleOptions{
        ScheduleType: "interval",
        IntervalMin:  intervalMin,
        WindowStart:  windowStart,
        WindowEnd:    windowEnd,
    }
    schedule, err := c.SchedulerService.CreateScheduledPost(context.Background(), userID, channelID, postData, opts)
    if err != nil {
        _, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
            ChatID: telego.ChatID{ID: chatID},
            Text:   fmt.Sprintf("❌ Erro ao criar agendamento: %s", err.Error()),
        })
        return
    }

    confirmText := fmt.Sprintf("✅ <b>Intervalo criado!</b>\n\nA cada %d minutos", intervalMin)
    if windowStart != "" {
        confirmText += fmt.Sprintf("\nJanela: %s às %s", windowStart, windowEnd)
    }
    _, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
        ChatID:      telego.ChatID{ID: chatID},
        Text:        confirmText,
        ParseMode:   telego.ModeHTML,
        ReplyMarkup: completionKB,
    })
    _ = schedule
```

### Passo 10 — Exibição em "Meus Agendamentos" (`my_schedules.go`)

Na listagem, exibir "Intervalo" como tipo quando `ScheduleType == "interval"`:

```go
text += fmt.Sprintf("%s <b>%s</b> → %s\n⏰ %s | %s\n\n", ...)
// O tipo "interval" já aparecerá como string; melhorar para exibir intervalo:
// Adicionar lógica para exibir "a cada Xmin" quando ScheduleType == "interval"
```

### Passo 11 — Dashboard (`dashboard/src/components/ScheduleTab.tsx`)

Adicionar `interval: 'Intervalo'` ao mapa `scheduleTypeLabels` (~linha 14).

### Passo 12 — Tipo TypeScript (`dashboard/src/types.ts`)

Adicionar os campos opcionais ao tipo `ScheduledPost`:

```ts
intervalMin?: number;
windowStart?: string;
windowEnd?: string;
```

### Passo 13 — Build e Testes

```bash
cd /home/malbs/Opencode/FreddyBot && go build ./cmd/FreddyBot
cd /home/malbs/Opencode/FreddyBot/dashboard && npm run build
```

## Riscos
- **Baixo**: Alteração retrocompatível (campos novos com default 0/"")
- **Mitigação de spam**: Mínimo de 5 minutos no backend
- **Migration automática**: GORM AutoMigrate adiciona as colunas sem downtime

## Impactos esperados
- Novo tipo `interval` disponível no PostBuilder do Telegram
- Compatível com todos os tipos existentes (sem breaking changes)
- Dashboard exibe o novo tipo corretamente

## Compatibilidade
- Linux, macOS, Windows, Docker, CI/CD

## Como testar

### Build
```bash
go build ./cmd/FreddyBot
cd dashboard && npm run build
```

### Teste funcional
1. Abrir o PostBuilder no Telegram
2. Clicar "Agendar Envio" → selecionar canal → "⏱️ Intervalo fixo"
3. Enviar `30` (a cada 30 minutos) ou `60\n08:00-22:00` (a cada 1h das 8h às 22h)
4. Verificar que o post é criado com `NextRunAt` correto
5. Aguardar o scheduler executar e verificar que `NextRunAt` avança corretamente

## Rollback
```bash
git checkout internal/database/models/scheduled_post.go
git checkout internal/core/services/scheduler.go
git checkout internal/telegram/handlers/events/postBuilder/postBuilder.go
git checkout internal/telegram/handlers/callbacks/my_schedules/my_schedules.go
git checkout dashboard/src/components/ScheduleTab.tsx
```

## Observações
- As colunas extras no banco (interval_min, window_start, window_end) ficarão com valores default para posts existentes e não afetam a lógica atual
- O ticker de 30s do scheduler já suporta a granularidade necessária, sem necessidade de alteração
- O limite de 50 agendamentos ativos por usuário permanece inalterado
