# Plano: Posts Agendados — Agendamentos de Postagens

## Pedido do usuário
Criar sistema de agendamento de posts que permita ao usuário:
1. Agendar um post para uma data/hora específica (uma vez)
2. Agendar um post recorrente (mesmo conteúdo, repete diariamente/semanalmente) — Cenário A
3. Agendar uma fila de posts (múltiplos posts, envia 1 por vez no horário) — Cenário B
4. Gerenciar agendamentos via botão "📅 Meus Agendamentos" no menu de conta do bot
5. Gerenciar agendamentos via nova tab na rota `/me/channels` da dashboard

## Objetivo
- Criar infraestrutura de agendamento (scheduler) que não existe no projeto
- Novo model `ScheduledPost` no banco de dados (persistente, não Redis)
- Integrar com o PostBuilder existente: novo botão "📅 Agendar" após salvar
- Novo botão no menu de perfil do bot para listar agendamentos
- Nova tab na dashboard `/me/channels` com visão global de todos os agendamentos
- API REST completa para CRUD de agendamentos

## Contexto atual

### PostBuilder (reutilizar)
- Coleta completa de dados: mídia, título, corpo, rodapé, reações (votes), botões
- `PostBuilderState` armazenado no Redis com TTL de 24h
- `sendFinalPostTelego()` envia para qualquer `chatID`, suporta todos os tipos de mídia
- Após salvar, mostra 2 botões: "🚀 Compartilhar" e "📢 Enviar para Canais"
- Fluxo de envio para canais: lista canais do usuário → seleciona → envia

### Infraestrutura existente (referência)
- `MessageQueue` em channelPost: chan + workers com rate limiting — padrão para o scheduler
- `BroadcastQueue` em appContainer: chan + workers com delay — outro padrão útil
- `time.AfterFunc` usado no Separator e MediaGroup — para delays pontuais

### O que NÃO existe
- Zero infraestrutura de cron/scheduler/ticker para jobs periódicos
- Nenhum model de agendamento no banco
- `SubscriptionService` tem métodos `ExpireSubscriptions()` e `SendRenewalInvoices()` marcados como "deve ser chamado periodicamente" mas nunca são chamados — evidência de que o scheduler faz falta

### Menu do bot (config/messages.yml)
- `profile-info` (menu de conta) tem botões: "Meus Canais" e "Início"
- Callback handler em `profile_info.go` — adiciona botão de Admin dinamicamente
- Novo botão "📅 Meus Agendamentos" será adicionado aqui

### Dashboard (/me/channels)
- NÃO tem sistema de tabs — é uma lista flat de canais
- Renderiza: greeting card → conta telegram card → premium tab → lista de canais
- Precisará adicionar um sistema de tabs simples (tipo "Canais" | "📅 Agendamentos")

## Arquivos analisados
- [postBuilder.go](file:///home/malbs/Opencode/FreddyBot/internal/telegram/handlers/events/postBuilder/postBuilder.go) — PostBuilder completo (save L681-713, send-to-channels L554-568, L726-780, sendFinalPost L782+)
- [cache/types.go](file:///home/malbs/Opencode/FreddyBot/internal/cache/types.go) — PostBuilderState (L21-32), PostBuilderButton
- [cache/cache.go](file:///home/malbs/Opencode/FreddyBot/internal/cache/cache.go) — SavePostBuilderSession, GetPostBuilderSession
- [messages.yml](file:///home/malbs/Opencode/FreddyBot/config/messages.yml) — menus do bot (start L1-19, profile-info L20-41)
- [profile_info.go](file:///home/malbs/Opencode/FreddyBot/internal/telegram/handlers/callbacks/profile_info/profile_info.go) — handler do menu de conta (L16-92)
- [loader_telego.go](file:///home/malbs/Opencode/FreddyBot/internal/telegram/loader_telego.go) — registro de handlers (L87-134)
- [appContainer.go](file:///home/malbs/Opencode/FreddyBot/internal/container/appContainer.go) — struct (L68-101), NewAppContainer (L103-221)
- [database.go](file:///home/malbs/Opencode/FreddyBot/internal/database/database.go) — AutoMigrate (L92-113)
- [models.go](file:///home/malbs/Opencode/FreddyBot/internal/database/models/models.go) — todos os models
- [channelPost.go](file:///home/malbs/Opencode/FreddyBot/internal/telegram/events/channelPost/channelPost.go) — MessageQueue pattern (L15-114)
- [App.tsx](file:///home/malbs/Opencode/FreddyBot/dashboard/src/App.tsx) — /me/channels rendering (L989-1091), tabs (L43-49)
- [routes.go](file:///home/malbs/Opencode/FreddyBot/internal/api/routes/routes.go) — rotas existentes
- [types.ts](file:///home/malbs/Opencode/FreddyBot/dashboard/src/types.ts) — tipos TypeScript
- [api.ts](file:///home/malbs/Opencode/FreddyBot/dashboard/src/api.ts) — chamadas API

## Arquivos que poderão ser modificados
- `internal/database/database.go` — adicionar ScheduledPost ao AutoMigrate
- `internal/container/appContainer.go` — registrar SchedulerService, iniciar goroutine do scheduler
- `internal/api/routes/routes.go` — novas rotas de agendamento
- `internal/api/dto/dto.go` — ScheduledPostDTO
- `internal/api/dto/mapper.go` — mapper ScheduledPost → DTO
- `internal/telegram/loader_telego.go` — registrar callback handler de schedule
- `internal/telegram/handlers/events/postBuilder/postBuilder.go` — botão "📅 Agendar" após save, callbacks pb-schedule
- `config/messages.yml` — botão "📅 Meus Agendamentos" no profile-info
- `dashboard/src/types.ts` — tipo ScheduledPost
- `dashboard/src/api.ts` — chamadas API de schedule
- `dashboard/src/App.tsx` — tab de agendamentos na rota /me/channels

## Arquivos que serão criados
- `internal/database/models/scheduled_post.go` — model ScheduledPost
- `internal/database/repositories/scheduled_post.go` — repositório GORM
- `internal/core/services/scheduler.go` — service + scheduler (goroutine ticker)
- `internal/api/controllers/schedulerController.go` — controllers REST
- `internal/api/types/scheduler.go` — request/response types
- `internal/telegram/handlers/callbacks/schedule/schedule.go` — callbacks do bot para agendamento
- `dashboard/src/components/ScheduleTab.tsx` — componente da dashboard

## Estratégia de implementação

### Modelo de dados

```
ScheduledPost
├── ID               (uuid PK)
├── OwnerID          (int64 FK → User)
├── ChannelID        (int64 FK → Channel)
├── ChannelTitle     (string — snapshot do título para exibição)
├── PostData         (text — JSON serializado do PostBuilderState)
│
├── ScheduleType     (string — "once" | "daily" | "weekly" | "queue")
├── ScheduleTime     (string — "14:00" HH:MM UTC para recorrente)
├── ScheduledAt      (*time.Time — datetime exato para "once")
├── ScheduleDays     (string — JSON "[1,3,5]" para weekly, vazio = todos)
├── NextRunAt        (time.Time — indexed, usado pelo scheduler)
├── RepeatUntil      (*time.Time — data limite opcional)
│
├── QueueGroupID     (string — agrupa posts de uma fila)
├── QueuePosition    (int — ordem na fila, 0-indexed)
├── LoopQueue        (bool — reiniciar fila ao terminar)
│
├── Status           (string — "pending" | "sent" | "cancelled" | "paused" | "failed")
├── SentAt           (*time.Time — quando foi enviado)
├── SentCount        (int — vezes enviado, para recorrente)
├── LastError        (string — último erro)
├── CreatedAt        (time.Time)
└── UpdatedAt        (time.Time)
```

### Como cada modo funciona

**Modo "once" (uma vez):**
```
ScheduleType: "once"
ScheduledAt:  2026-07-25T14:00:00Z
NextRunAt:    2026-07-25T14:00:00Z
→ Scheduler encontra, envia, marca Status="sent"
```

**Modo "daily" (Cenário A — mesmo post todo dia):**
```
ScheduleType: "daily"
ScheduleTime: "14:00"
NextRunAt:    2026-07-24T14:00:00Z (calculado)
RepeatUntil:  2026-08-24T00:00:00Z (opcional)
→ Scheduler encontra, envia, SentCount++
→ Calcula próximo NextRunAt = amanhã 14:00
→ Se NextRunAt > RepeatUntil → Status="sent" (finalizado)
```

**Modo "weekly" (mesmo post em dias específicos):**
```
ScheduleType: "weekly"
ScheduleTime: "09:00"
ScheduleDays: "[1,3,5]" (seg, qua, sex)
NextRunAt:    2026-07-23T09:00:00Z (próxima segunda)
→ Igual ao daily, mas calcula NextRunAt pulando para o próximo dia válido
```

**Modo "queue" (Cenário B — fila de posts):**
```
Post 1: QueueGroupID="abc", QueuePosition=0, Status="pending", NextRunAt=amanhã 14:00
Post 2: QueueGroupID="abc", QueuePosition=1, Status="pending", NextRunAt não calculado ainda
Post 3: QueueGroupID="abc", QueuePosition=2, Status="pending", LoopQueue=true

→ Scheduler encontra Post 1 (menor QueuePosition pendente com NextRunAt <= now)
→ Envia Post 1, marca Status="sent"
→ Encontra Post 2 (próximo na fila), calcula NextRunAt = próximo dia 14:00
→ Quando Post 3 é enviado:
   → Se LoopQueue=true: reseta todos para "pending", Post 1 recebe NextRunAt
   → Se LoopQueue=false: marca QueueGroupID como completo
```

## Passos detalhados

### Phase 1 — Backend: Model, Repository, Service

**Passo 1.1: Model `ScheduledPost`** (novo arquivo `models/scheduled_post.go`)
- Definir struct conforme modelo de dados acima
- Índices: `idx_schedule_next_run` em `(Status, NextRunAt)` para queries eficientes do scheduler
- Índice: `idx_schedule_owner` em `OwnerID` para listagem por usuário
- Índice: `idx_schedule_queue` em `QueueGroupID` para operações de fila

**Passo 1.2: Registrar no AutoMigrate** (`database.go`)
- Adicionar `&models.ScheduledPost{}` à lista de AutoMigrate

**Passo 1.3: Repository** (novo arquivo `repositories/scheduled_post.go`)
```go
type ScheduledPostRepository struct { db *gorm.DB }

func (r *ScheduledPostRepository) Create(post *models.ScheduledPost) error
func (r *ScheduledPostRepository) GetByID(id string) (*models.ScheduledPost, error)
func (r *ScheduledPostRepository) GetByOwner(ownerID int64) ([]models.ScheduledPost, error)
func (r *ScheduledPostRepository) GetDuePosts(now time.Time) ([]models.ScheduledPost, error)
// → WHERE status = 'pending' AND next_run_at <= now ORDER BY next_run_at ASC
func (r *ScheduledPostRepository) GetQueueGroup(queueGroupID string) ([]models.ScheduledPost, error)
func (r *ScheduledPostRepository) UpdateStatus(id string, status string) error
func (r *ScheduledPostRepository) UpdateNextRunAt(id string, nextRunAt time.Time) error
func (r *ScheduledPostRepository) MarkSent(id string, sentAt time.Time) error
func (r *ScheduledPostRepository) Delete(id string) error
func (r *ScheduledPostRepository) GetByOwnerAndChannel(ownerID, channelID int64) ([]models.ScheduledPost, error)
func (r *ScheduledPostRepository) CountByOwner(ownerID int64) (int64, error)
```

**Passo 1.4: SchedulerService** (novo arquivo `services/scheduler.go`)
```go
type SchedulerService struct {
    repo        *repositories.ScheduledPostRepository
    cacheService *cache.Service
    bot         *telego.Bot
}

// CRUD
func (s *SchedulerService) CreateScheduledPost(ownerID, channelID int64, channelTitle string, postData string, scheduleType string, opts ScheduleOptions) (*models.ScheduledPost, error)
func (s *SchedulerService) CancelScheduledPost(id string, ownerID int64) error
func (s *SchedulerService) PauseScheduledPost(id string, ownerID int64) error
func (s *SchedulerService) ResumeScheduledPost(id string, ownerID int64) error
func (s *SchedulerService) GetUserSchedules(ownerID int64) ([]models.ScheduledPost, error)
func (s *SchedulerService) DeleteScheduledPost(id string, ownerID int64) error

// Scheduler loop
func (s *SchedulerService) Start(ctx context.Context)
// → goroutine com time.Ticker de 30s
// → chama processDuePosts() a cada tick

func (s *SchedulerService) processDuePosts()
// → repo.GetDuePosts(time.Now())
// → para cada post: sendScheduledPost()
// → se recorrente: calculateNextRunAt() e atualiza
// → se queue: avança para o próximo post da fila

func (s *SchedulerService) sendScheduledPost(post *models.ScheduledPost) error
// → desserializa PostData em PostBuilderState
// → chama sendFinalPostTelego() (importado do postBuilder)
// → atualiza status/sentAt/sentCount
// → notifica o dono via DM: "✅ Post agendado enviado no canal X"

func (s *SchedulerService) calculateNextRunAt(post *models.ScheduledPost) *time.Time
// → "daily": amanhã no mesmo horário
// → "weekly": próximo dia da semana válido
// → "once": nil (não repete)
// → verifica RepeatUntil
```

**Passo 1.5: Registrar no AppContainer** (`appContainer.go`)
- Adicionar `SchedulerService *services.SchedulerService` ao struct
- No `NewAppContainer`: criar repo → service → atribuir
- Chamar `go container.SchedulerService.Start(context.Background())` no startup

### Phase 2 — API REST

**Passo 2.1: Request/Response types** (novo arquivo `api/types/scheduler.go`)
```go
type CreateScheduleRequest struct {
    ChannelID    int64  `json:"channelId" binding:"required"`
    SessionID    string `json:"sessionId" binding:"required"`
    ScheduleType string `json:"scheduleType" binding:"required,oneof=once daily weekly"`
    ScheduleTime string `json:"scheduleTime"` // "HH:MM" para recorrente
    ScheduledAt  string `json:"scheduledAt"`  // ISO 8601 para once
    ScheduleDays []int  `json:"scheduleDays"` // [1,3,5] para weekly
    RepeatUntil  string `json:"repeatUntil"`  // ISO 8601 opcional
}

type AddToQueueRequest struct {
    ChannelID    int64  `json:"channelId" binding:"required"`
    SessionID    string `json:"sessionId" binding:"required"`
    QueueGroupID string `json:"queueGroupId"`  // vazio = cria nova fila
    ScheduleTime string `json:"scheduleTime" binding:"required"`
    LoopQueue    bool   `json:"loopQueue"`
}

type UpdateScheduleStatusRequest struct {
    Status string `json:"status" binding:"required,oneof=paused pending cancelled"`
}
```

**Passo 2.2: Controllers** (novo arquivo `api/controllers/schedulerController.go`)
```go
func (ctrl *SchedulerController) CreateSchedule(c *gin.Context)      // POST /schedule
func (ctrl *SchedulerController) AddToQueue(c *gin.Context)          // POST /schedule/queue
func (ctrl *SchedulerController) GetMySchedules(c *gin.Context)      // GET /schedule
func (ctrl *SchedulerController) GetScheduleByID(c *gin.Context)     // GET /schedule/:id
func (ctrl *SchedulerController) UpdateStatus(c *gin.Context)        // PUT /schedule/:id/status
func (ctrl *SchedulerController) DeleteSchedule(c *gin.Context)      // DELETE /schedule/:id
func (ctrl *SchedulerController) GetQueueGroup(c *gin.Context)       // GET /schedule/queue/:groupId
```

**Passo 2.3: Rotas** (`routes.go`)
```go
scheduleRoutes := authorized.Group("/schedule")
{
    scheduleRoutes.POST("/", schedulerController.CreateSchedule)
    scheduleRoutes.POST("/queue", schedulerController.AddToQueue)
    scheduleRoutes.GET("/", schedulerController.GetMySchedules)
    scheduleRoutes.GET("/:id", schedulerController.GetScheduleByID)
    scheduleRoutes.PUT("/:id/status", schedulerController.UpdateStatus)
    scheduleRoutes.DELETE("/:id", schedulerController.DeleteSchedule)
    scheduleRoutes.GET("/queue/:groupId", schedulerController.GetQueueGroup)
}
```

**Passo 2.4: DTO** (`dto.go` e `mapper.go`)
```go
type ScheduledPostDTO struct {
    ID            string `json:"id"`
    OwnerID       int64  `json:"ownerId"`
    ChannelID     int64  `json:"channelId"`
    ChannelTitle  string `json:"channelTitle"`
    MediaType     string `json:"mediaType"` // extraído do PostData para preview
    ScheduleType  string `json:"scheduleType"`
    ScheduleTime  string `json:"scheduleTime"`
    ScheduledAt   string `json:"scheduledAt,omitempty"`
    ScheduleDays  string `json:"scheduleDays,omitempty"`
    NextRunAt     string `json:"nextRunAt"`
    RepeatUntil   string `json:"repeatUntil,omitempty"`
    QueueGroupID  string `json:"queueGroupId,omitempty"`
    QueuePosition int    `json:"queuePosition"`
    LoopQueue     bool   `json:"loopQueue"`
    Status        string `json:"status"`
    SentAt        string `json:"sentAt,omitempty"`
    SentCount     int    `json:"sentCount"`
    LastError     string `json:"lastError,omitempty"`
    CreatedAt     string `json:"createdAt"`
}
```

### Phase 3 — Bot: PostBuilder + Menu

**Passo 3.1: Botão "📅 Agendar" no PostBuilder** (`postBuilder.go`)

No callback `pb-save` (L681-713), adicionar terceira linha de botão:
```go
{Text: "📅 Agendar Envio", CallbackData: "pb-schedule:" + id},
```

**Passo 3.2: Callback handler de agendamento** (novo arquivo `callbacks/schedule/schedule.go`)

Callbacks a tratar:
```
pb-schedule:<sessionID>
  → Mostra lista de canais do usuário (igual pb-send-to-channels)

pb-schedule-channel:<channelID>:<sessionID>
  → Mostra opções de tipo: [📅 Uma vez] [🔁 Diário] [📆 Semanal]

pb-schedule-type:<type>:<channelID>:<sessionID>
  → Se "once": mostra opções rápidas + personalizado
  → Se "daily"/"weekly": pede horário

pb-schedule-quick:<duration>:<channelID>:<sessionID>
  → Atalhos: "1h", "3h", "6h", "amanha-9", "amanha-14"
  → Calcula ScheduledAt, salva no banco, confirma

pb-schedule-time:<channelID>:<sessionID>
  → Pede horário personalizado (Step = "awaiting_schedule_time")
  → Formato: "HH:MM" ou "DD/MM HH:MM"

pb-schedule-confirm:<scheduleID>
  → Confirmação final com resumo
  → Salva no banco via SchedulerService
  → "✅ Agendado para 25/07 às 14:00 no canal @X"

pb-schedule-days:<channelID>:<sessionID>
  → Para weekly: mostra botões dos dias da semana (toggle)
  → [Seg ✅] [Ter] [Qua ✅] [Qui] [Sex ✅] [Sab] [Dom]
```

**Passo 3.3: Botão "📅 Meus Agendamentos" no menu do bot** (`messages.yml`)

No `profile-info` (L32-38), adicionar entre "Meus Canais" e "Início":
```yaml
  - - text: "Meus Agendamentos"
      callback_data: "my-schedules"
      custom_emoji: "5368324170671202286"
```

**Passo 3.4: Handler "my-schedules"** (no mesmo arquivo `schedule/schedule.go`)
```
my-schedules
  → Lista agendamentos do usuário (SchedulerService.GetUserSchedules)
  → Mostra cards inline:
    "📸 @canal — 25/07 14:00 — ⏳ Aguardando"
    [👀 Preview] [⏸ Pausar] [❌ Cancelar]

  → Se nenhum: "Você não tem agendamentos. Use o Post Builder para criar!"
  → Paginação se > 5 agendamentos
```

**Passo 3.5: Registrar handlers** (`loader_telego.go`)
```go
bh.Handle(callbackSchedule.HandlerTelego(c), telegohandler.CallbackDataPrefix("pb-schedule"))
bh.Handle(callbackSchedule.MySchedulesHandlerTelego(c), telegohandler.CallbackDataEqual("my-schedules"))
```

### Phase 4 — Dashboard: Tab de Agendamentos

**Passo 4.1: Tipos TypeScript** (`types.ts`)
```typescript
export interface ScheduledPost {
  id: string;
  ownerId: number;
  channelId: number;
  channelTitle: string;
  mediaType: string;
  scheduleType: 'once' | 'daily' | 'weekly' | 'queue';
  scheduleTime: string;
  scheduledAt?: string;
  scheduleDays?: string;
  nextRunAt: string;
  repeatUntil?: string;
  queueGroupId?: string;
  queuePosition: number;
  loopQueue: boolean;
  status: 'pending' | 'sent' | 'cancelled' | 'paused' | 'failed';
  sentAt?: string;
  sentCount: number;
  lastError?: string;
  createdAt: string;
}
```

**Passo 4.2: API calls** (`api.ts`)
```typescript
getMySchedules(): Promise<ScheduledPost[]>
cancelSchedule(id: string): Promise<void>
pauseSchedule(id: string): Promise<void>
resumeSchedule(id: string): Promise<void>
deleteSchedule(id: string): Promise<void>
```

**Passo 4.3: Componente `ScheduleTab.tsx`** (novo)
- Lista de cards com:
  - Ícone de tipo de mídia (📸📹🎵📄)
  - Título do canal
  - Data/hora (ou "Diário às 14:00", "Seg/Qua/Sex às 09:00")
  - Badge de status (⏳ Aguardando | ✅ Enviado | ⏸ Pausado | ❌ Cancelado | ⚠️ Falhou)
  - Badge de tipo (Uma vez | 🔁 Diário | 📆 Semanal | 📚 Fila [3/10])
  - Botões: [⏸ Pausar / ▶️ Retomar] [❌ Cancelar]
- Filtro por canal (dropdown)
- Filtro por status
- Mensagem vazia: "Nenhum agendamento. Crie posts no bot e use 📅 Agendar!"

**Passo 4.4: Integrar na rota `/me/channels`** (`App.tsx`)
- Adicionar mini tab system na seção de canais:
  ```
  [📋 Canais] [📅 Agendamentos]
  ```
- Tab "Canais" → lista existente (sem alteração)
- Tab "Agendamentos" → `<ScheduleTab />` (novo componente)
- Estado ativo controlado por `useState('canais')`

### Phase 5 — Scheduler (goroutine)

**Passo 5.1: Loop principal** (dentro de `services/scheduler.go`)
```go
func (s *SchedulerService) Start(ctx context.Context) {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()
    
    // Processar posts pendentes no startup (caso bot tenha reiniciado)
    s.processDuePosts()
    
    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            s.processDuePosts()
        }
    }
}
```

**Passo 5.2: Processamento de posts** 
```go
func (s *SchedulerService) processDuePosts() {
    posts, err := s.repo.GetDuePosts(time.Now())
    // Para cada post:
    //   1. Desserializar PostData → PostBuilderState
    //   2. Chamar sendFinalPostTelego()
    //   3. Se sucesso:
    //      - "once" → status = "sent"
    //      - "daily"/"weekly" → SentCount++, NextRunAt = calculateNext()
    //      - "queue" → status = "sent", ativar próximo da fila
    //   4. Se erro:
    //      - LastError = err.Error()
    //      - Retry em 5 min (NextRunAt += 5min), max 3 retries
    //   5. Notificar dono via DM
    //   6. Rate limit: 1s delay entre envios para canais diferentes
}
```

**Passo 5.3: Cálculo de próxima execução**
```go
func calculateNextRunAt(post *models.ScheduledPost) *time.Time {
    // "daily": amanhã no mesmo ScheduleTime
    // "weekly": próximo dia em ScheduleDays no mesmo ScheduleTime
    // "once": nil
    // Respeita RepeatUntil
}
```

**Passo 5.4: Avanço de fila (queue)**
```go
func (s *SchedulerService) advanceQueue(sentPost *models.ScheduledPost) {
    // 1. Buscar todos os posts do QueueGroupID
    // 2. Encontrar o próximo (QueuePosition > sentPost.QueuePosition, Status="pending")
    // 3. Se encontrou: definir NextRunAt = próximo horário
    // 4. Se não encontrou e LoopQueue=true:
    //    - Resetar todos para Status="pending"
    //    - Definir NextRunAt do primeiro (QueuePosition=0)
    // 5. Se não encontrou e LoopQueue=false:
    //    - Fila concluída, notificar dono
}
```

### Phase 6 — Notificações

**Passo 6.1: Notificação ao enviar** (DM para o dono)
```
✅ Post agendado enviado!
📢 Canal: @meucanal
📅 25/07/2026 às 14:00
🔁 Próximo envio: 26/07/2026 às 14:00
```

**Passo 6.2: Notificação de erro**
```
⚠️ Falha ao enviar post agendado
📢 Canal: @meucanal
📅 25/07/2026 às 14:00
❌ Erro: message too long
Nova tentativa em 5 minutos...
```

**Passo 6.3: Notificação de fila concluída**
```
📚 Fila de posts concluída!
📢 Canal: @meucanal
✅ 10/10 posts enviados
```

## Riscos
- **Bot reinicia**: Posts pendentes no banco permanecem — o scheduler pega no startup
- **Rate limit do Telegram**: Delay de 1s entre envios + retry com backoff
- **File ID expira**: Para agendamentos > 30 dias, o file_id pode expirar. Mitigação: log do erro e notifica o dono
- **Fuso horário**: V1 usa UTC. O bot mostra "daqui a Xh" como confirmação. V2 pode adicionar config de timezone por usuário
- **Redis TTL**: PostBuilderState tem 24h TTL no Redis. Ao agendar, os dados são copiados para o banco (PostData JSON) — independente do Redis
- **Muitos agendamentos simultâneos**: O scheduler processa sequencialmente com delay. Para V1 é suficiente
- **PostBuilder sendFinalPostTelego é um método no pacote postBuilder**: Precisará ser exportado ou extraído para um pacote compartilhado

## Impactos esperados
- **Zero impacto** em funcionalidades existentes — PostBuilder ganha um botão extra, menu ganha um item, dashboard ganha uma tab
- **Infraestrutura de scheduler** pode ser reutilizada futuramente (ex: subscription expiration, cleanup jobs)
- **Nova tabela** no banco (ScheduledPost) com migrations automáticas via GORM
- **Goroutine adicional** no startup — leve, poll a cada 30s

## Compatibilidade
- ✅ Linux — Go + SQLite/PostgreSQL
- ✅ macOS — desenvolvimento local
- ✅ Docker — sem dependência extra
- ✅ CI/CD — sem breaking changes

## Como testar

### Build
```bash
cd /home/malbs/Opencode/FreddyBot && go build ./cmd/... ./internal/... ./pkg/...
```

### Dashboard
```bash
cd /home/malbs/Opencode/FreddyBot/dashboard && npm run dev
```

### Testes manuais — Modo "once"
1. Criar post no PostBuilder (enviar mídia no DM)
2. Editar título/corpo, salvar
3. Clicar "📅 Agendar Envio"
4. Selecionar canal → "📅 Uma vez" → "Em 1h"
5. Verificar confirmação
6. Esperar 1h → verificar se o post foi enviado no canal
7. Verificar notificação no DM

### Testes manuais — Modo "daily" (Cenário A)
1. Criar e salvar post
2. Agendar → canal → "🔁 Diário" → "14:00"
3. Verificar na lista de agendamentos (bot + dashboard)
4. Esperar que o horário chegue → verificar envio
5. Verificar que NextRunAt foi recalculado para o dia seguinte
6. Pausar agendamento → verificar que não envia
7. Retomar → verificar que volta a funcionar

### Testes manuais — Modo "queue" (Cenário B)
1. Criar 3 posts diferentes no PostBuilder
2. Agendar cada um para o mesmo canal com "📚 Fila" no mesmo horário
3. Verificar que a fila mostra [1/3], [2/3], [3/3]
4. Esperar envios → verificar que envia 1 por vez
5. Testar com LoopQueue=true → verificar que reinicia

### Testes manuais — Dashboard
1. Abrir /me/channels
2. Clicar na tab "📅 Agendamentos"
3. Verificar que lista os agendamentos
4. Cancelar um agendamento → verificar status
5. Filtrar por canal

## Rollback
1. Parar a goroutine do scheduler (remover `go container.SchedulerService.Start()`)
2. Remover botão "📅 Agendar" do PostBuilder
3. Remover botão "📅 Meus Agendamentos" do messages.yml
4. Remover tab de agendamentos da dashboard
5. A tabela `scheduled_posts` pode ficar no banco — é inofensiva
6. Alternativa: `git revert` do(s) commit(s)

## Observações
- **sendFinalPostTelego** atualmente é uma função não-exportada no pacote `postBuilder`. Para o scheduler chamar, precisará ser exportada (`SendFinalPostTelego`) ou extraída para um pacote compartilhado (ex: `internal/core/services/post_sender.go`). A segunda opção é mais limpa arquiteturalmente.
- **Fuso horário**: V1 trabalha em UTC. O bot pode exibir horários relativos ("daqui a 3 horas") para evitar confusão. Uma config de timezone por usuário pode ser adicionada na V2.
- **Limite de agendamentos**: Considerar um limite por usuário (ex: 20 agendamentos ativos para free, 100 para premium) para evitar abuso.
- **Limpeza de histórico**: Posts com status "sent" ou "cancelled" podem ser limpos após 30 dias via job periódico (aproveitando o scheduler).
- **Custom emoji nos botões de reação (votes)**: O PostBuilderState já suporta `eid:` prefix para custom emoji — o scheduler preserva isso ao desserializar o PostData.
- **Esta feature pode justificar uma flag Premium**: Agendamentos recorrentes e filas são funcionalidades avançadas que agregam valor significativo.
