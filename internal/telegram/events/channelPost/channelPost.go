package channelpost

import (
	"context"
	"sync"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/leirbagxis/FreddyBot/internal/container"
	"github.com/leirbagxis/FreddyBot/pkg/logger"
)

// ✅ SISTEMA DE FILA UNIFICADO
type Job interface {
	Run() error
	GetChannelID() int64
}

type MessageQueue struct {
	queue       chan Job
	mu          sync.Mutex
	isRunning   bool
	lastProcess sync.Map // map[int64]time.Time
}

type PipelineJob struct {
	Ctx      *ProcessingContext
	Pipeline *Pipeline
}

func (j PipelineJob) Run() error {
	return j.Pipeline.Execute(j.Ctx)
}

func (j PipelineJob) GetChannelID() int64 {
	if j.Ctx.Channel != nil {
		return j.Ctx.Channel.ID
	}
	return 0
}

var (
	groupSeparators = sync.Map{} // string -> bool
)

var messageQueue *MessageQueue

func init() {
	messageQueue = NewMessageQueue()
}

func NewMessageQueue() *MessageQueue {
	mq := &MessageQueue{
		queue:     make(chan Job, 5000),
		isRunning: true,
	}
	// Iniciar 20 workers para processamento paralelo (ideal para 2 vCPUs + I/O)
	for i := 0; i < 20; i++ {
		go mq.worker()
	}
	return mq
}

func (mq *MessageQueue) worker() {
	for mq.isRunning {
		select {
		case job := <-mq.queue:
			// Controle per-chat: evita processar mensagens do mesmo chat muito rápido
			channelID := job.GetChannelID()
			if channelID != 0 {
				last, ok := mq.lastProcess.Load(channelID)
				if ok {
					elapsed := time.Since(last.(time.Time))
					if elapsed < 500*time.Millisecond {
						time.Sleep(500*time.Millisecond - elapsed)
					}
				}
				mq.lastProcess.Store(channelID, time.Now())
			}

			if err := job.Run(); err != nil {
				logger.Error("BOT", "❌ Erro ao processar job da fila: %v", err)
			}
		}
	}
}

func (mq *MessageQueue) AddV2ToQueue(pCtx *ProcessingContext, pipeline *Pipeline) {
	select {
	case mq.queue <- PipelineJob{Ctx: pCtx, Pipeline: pipeline}:
		logger.Bot("📥 Mensagem V2 adicionada à fila (tamanho: %d)", len(mq.queue))
	default:
		logger.Bot("⚠️ Fila cheia, descartando mensagem V2")
	}
}

func Handler(c *container.AppContainer) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		if update.ChannelPost != nil {
			logger.Bot("🚀 [%d] Novo post recebido no canal %d", update.ChannelPost.ID, update.ChannelPost.Chat.ID)
		}

		// 1. Execution Pipeline: Transformation -> Decoration -> Dispatch (Runs in Worker)
		executionPipeline := NewPipeline(
			"Execution",
			StageTransform(c),
			StageDecorate(c),
			StageSend(c),
		)

		// 2. Discovery Pipeline: Filters -> Metadata -> Grouping -> Queue (Runs in Handler thread)
		discoveryPipeline := NewPipeline(
			"Discovery",
			StagePreflight(c),
			StageSpecialFlows(c),
			StageMediaGrouping(c, executionPipeline),
			StageQueue(c, executionPipeline),
		)

		pCtx := NewProcessingContext(ctx, b, update, discoveryPipeline)
		_ = discoveryPipeline.Execute(pCtx)
	}
}
