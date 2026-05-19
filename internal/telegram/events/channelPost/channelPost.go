package channelpost

import (
	"context"
	"sync"
	"time"

	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegohandler"
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

type PipelineJobTelego struct {
	Ctx      *ProcessingContextTelego
	Pipeline *PipelineTelego
}

func (j PipelineJobTelego) Run() error {
	return j.Pipeline.Execute(j.Ctx)
}

func (j PipelineJobTelego) GetChannelID() int64 {
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

func (mq *MessageQueue) AddTelegoToQueue(pCtx *ProcessingContextTelego, pipeline *PipelineTelego) {
	select {
	case mq.queue <- PipelineJobTelego{Ctx: pCtx, Pipeline: pipeline}:
		logger.Bot("📥 Mensagem Telego adicionada à fila (tamanho: %d)", len(mq.queue))
	default:
		logger.Bot("⚠️ Fila cheia, descartando mensagem Telego")
	}
}

func HandlerTelego(c *container.AppContainer) telegohandler.Handler {
	return func(ctx *telegohandler.Context, update telego.Update) error {
		if update.ChannelPost != nil {
			logger.Bot("🚀 [%d] Novo post recebido no canal %d (Telego)", update.ChannelPost.MessageID, update.ChannelPost.Chat.ID)
		}

		// 1. Execution Pipeline
		executionPipeline := NewPipelineTelego(
			"Execution",
			StageTransformTelego(c),
			StageDecorateTelego(c),
			StageSendTelego(c),
		)

		// 2. Discovery Pipeline
		discoveryPipeline := NewPipelineTelego(
			"Discovery",
			StagePreflightTelego(c),
			StageSpecialFlowsTelego(c),
			StageMediaGroupingTelego(c, executionPipeline),
			StageQueueTelego(c, executionPipeline),
		)

		pCtx := NewProcessingContextTelego(context.Background(), ctx.Bot(), update, discoveryPipeline)
		_ = discoveryPipeline.Execute(pCtx)
		return nil
	}
}
