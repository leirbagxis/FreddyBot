package channelpost

import (
	"context"
	"fmt"
	"html"
	"log"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/leirbagxis/FreddyBot/internal/container"
	dbmodels "github.com/leirbagxis/FreddyBot/internal/database/models"
	"github.com/leirbagxis/FreddyBot/pkg/config"
)

// MessageQueue manages message processing queue with rate limit and retries
type MessageQueue struct {
	queue       chan QueueItem
	rateLimiter *time.Ticker
	mu          sync.Mutex
	isRunning   bool
}

// QueueItem is an item in the processing queue
type QueueItem struct {
	MessageType        MessageType
	Channel            *dbmodels.Channel
	Post               *models.Message
	Buttons            []dbmodels.Button
	MessageEditAllowed bool
	Processor          *MessageProcessor
	OwnerID            int64 // admin user to notify
}

// groupSeparators to manage separators per media group id
var groupSeparators = sync.Map{} // map[string]bool

var messageQueue *MessageQueue

func init() {
	messageQueue = NewMessageQueue()
}

// NewMessageQueue creates new MessageQueue
func NewMessageQueue() *MessageQueue {
	mq := &MessageQueue{
		queue:       make(chan QueueItem, 1000),
		rateLimiter: time.NewTicker(time.Second),
		isRunning:   true,
	}
	go mq.worker()
	return mq
}

// worker processes items from queue with rate limiting and retry logic
func (mq *MessageQueue) worker() {
	for mq.isRunning {
		select {
		case item := <-mq.queue:
			<-mq.rateLimiter.C
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			err := mq.processWithRetry(ctx, item)
			if err != nil {
				log.Printf("❌ Erro ao processar item da fila: %v", err)
				NotifyOwner(item.Processor.bot, item.OwnerID, fmt.Sprintf("Erro ao processar item da fila: %v", err))
			}
			cancel()
		}
	}
}

// processWithRetry retries processing on rate limit errors
func (mq *MessageQueue) processWithRetry(ctx context.Context, item QueueItem) error {
	maxRetries := 3
	baseDelay := time.Second

	for attempt := 0; attempt < maxRetries; attempt++ {
		err := item.Processor.ProcessMessage(ctx, item.MessageType, item.Channel, item.Post, item.Buttons, item.MessageEditAllowed)
		if err == nil {
			return nil
		}
		if strings.Contains(err.Error(), "Too Many Requests") {
			retryAfter := extractRetryAfter(err.Error())
			if retryAfter == 0 {
				retryAfter = int(baseDelay.Seconds()) * (attempt + 1)
			}
			log.Printf("⏳ Rate limit atingido, aguardando %d segundos (tentativa %d/%d)", retryAfter, attempt+1, maxRetries)
			NotifyOwner(item.Processor.bot, item.OwnerID, fmt.Sprintf("Rate limit atingido. Pausando %d segundos (tentativa %d/%d)", retryAfter, attempt+1, maxRetries))
			time.Sleep(time.Duration(retryAfter) * time.Second)
			continue
		}
		return err
	}

	return fmt.Errorf("falha após %d tentativas", maxRetries)
}

// extractRetryAfter extrai o tempo do retry a partir da mensagem de erro
func extractRetryAfter(errorMsg string) int {
	re := regexp.MustCompile(`retry after (\d+)`)
	matches := re.FindStringSubmatch(errorMsg)
	if len(matches) > 1 {
		if retryAfter, err := strconv.Atoi(matches[1]); err == nil {
			return retryAfter
		}
	}
	return 0
}

// AddToQueue adiciona mensagem na fila ou descarta se fila cheia
func (mq *MessageQueue) AddToQueue(messageType MessageType, channel *dbmodels.Channel, post *models.Message, buttons []dbmodels.Button, messageEditAllowed bool, processor *MessageProcessor, ownerID int64) {
	select {
	case mq.queue <- QueueItem{
		MessageType:        messageType,
		Channel:            channel,
		Post:               post,
		Buttons:            buttons,
		MessageEditAllowed: messageEditAllowed,
		Processor:          processor,
		OwnerID:            ownerID,
	}:
		log.Printf("📥 Mensagem adicionada à fila (tamanho: %d)", len(mq.queue))
	default:
		log.Printf("⚠️ Fila cheia, descartando mensagem")
		NotifyOwner(processor.bot, ownerID, "Fila cheia: mensagem descartada!")
	}
}

func NotifyOwner(v *bot.Bot, ownerID int64, msg string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := v.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    ownerID,
		Text:      fmt.Sprintf("<b>BOT:</b> %s", html.EscapeString(msg)),
		ParseMode: "HTML",
	})
	if err != nil {
		log.Printf("❌ Falha ao notificar o owner: %v", err)
	}
}

// Handler principal para atualizações de ChannelPost
func Handler(c *container.AppContainer) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		dbCtx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		post := update.ChannelPost
		if post == nil {
			// Não processa se não for ChannelPost
			return
		}

		processor := NewMessageProcessor(b)
		chat := post.Chat
		ownerID := config.OwnerID

		channel, err := c.ChannelRepo.GetChannelWithRelations(dbCtx, chat.ID)
		if err != nil {
			log.Printf("Canal %d não encontrado: %v", chat.ID, err)
			NotifyOwner(b, ownerID, fmt.Sprintf("Canal %d não encontrado: %v", chat.ID, err))
			return
		}

		go func() {
			updateCtx, updateCancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer updateCancel()

			updatedChannel, hasChanges := processor.UpdateChannelBasicInfo(updateCtx, chat.ID, channel)
			if hasChanges {
				if err := c.ChannelRepo.UpdateChannelBasicInfoAndFirstButton(updateCtx, updatedChannel); err != nil {
					log.Printf("❌ Erro ao salvar informações do canal %d: %v", chat.ID, err)
					NotifyOwner(b, ownerID, fmt.Sprintf("Erro ao salvar informações do canal %d: %v", chat.ID, err))
				} else {
					log.Printf("✅ Canal %d: informações básicas e primeiro botão atualizados automaticamente", chat.ID)
				}
			}
		}()

		messageType := processor.GetMessageType(post)
		if messageType == "" {
			return
		}

		permissions := processor.CheckPermissions(channel, messageType)
		if !permissions.CanEdit && !permissions.CanAddButtons {
			log.Printf("❌ Sem permissões para processar mensagem no canal %d", channel.ID)
			NotifyOwner(b, ownerID, fmt.Sprintf("Sem permissões para processar mensagem no canal %d", channel.ID))
			return
		}

		var finalButtons []dbmodels.Button
		if permissions.CanAddButtons {
			finalButtons = channel.Buttons
		}

		go func() {
			messageQueue.AddToQueue(messageType, channel, post, finalButtons, permissions.CanEdit, processor, ownerID)
		}()

		if channel.Separator != nil && (permissions.CanEdit || permissions.CanAddButtons) {
			go func() {
				time.Sleep(2 * time.Second)
				processor.HandleSeparator(channel, post, messageType)
			}()
		}
	}
}
