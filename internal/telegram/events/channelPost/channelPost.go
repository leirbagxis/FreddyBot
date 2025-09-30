package channelpost

import (
	"context"
	"fmt"
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
	"github.com/leirbagxis/FreddyBot/internal/utils"
)

// ✅ SISTEMA DE FILA SIMPLIFICADO
type MessageQueue struct {
	queue       chan QueueItem
	rateLimiter *time.Ticker
	mu          sync.Mutex
	isRunning   bool
}

type QueueItem struct {
	MessageType        MessageType
	Channel            *dbmodels.Channel
	Post               *models.Message
	Buttons            []dbmodels.Button
	MessageEditAllowed bool
	Processor          *MessageProcessor
}

// ✅ CONTROLE SIMPLES DE SEPARATORS POR GRUPO
var groupSeparators = sync.Map{} // string -> bool

var messageQueue *MessageQueue

func init() {
	messageQueue = NewMessageQueue()
}

func NewMessageQueue() *MessageQueue {
	mq := &MessageQueue{
		queue:       make(chan QueueItem, 1000),
		rateLimiter: time.NewTicker(time.Second),
		isRunning:   true,
	}
	go mq.worker()
	return mq
}

func (mq *MessageQueue) worker() {
	for mq.isRunning {
		select {
		case item := <-mq.queue:
			<-mq.rateLimiter.C
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			err := mq.processWithRetry(ctx, item)
			if err != nil {
				log.Printf("❌ Erro ao processar item da fila: %v", err)
			}
			cancel()
		}
	}
}

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
			time.Sleep(time.Duration(retryAfter) * time.Second)
			continue
		}
		return err
	}
	return fmt.Errorf("falha após %d tentativas", maxRetries)
}

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

func (mq *MessageQueue) AddToQueue(messageType MessageType, channel *dbmodels.Channel, post *models.Message, buttons []dbmodels.Button, messageEditAllowed bool, processor *MessageProcessor) {
	select {
	case mq.queue <- QueueItem{
		MessageType:        messageType,
		Channel:            channel,
		Post:               post,
		Buttons:            buttons,
		MessageEditAllowed: messageEditAllowed,
		Processor:          processor,
	}:
		log.Printf("📥 Mensagem adicionada à fila (tamanho: %d)", len(mq.queue))
	default:
		log.Printf("⚠️ Fila cheia, descartando mensagem")
	}
}

func Handler(c *container.AppContainer) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		dbCtx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		post := update.ChannelPost
		if post == nil {
			return
		}

		processor := NewMessageProcessor(b)
		chat := post.Chat

		channel, err := c.ChannelRepo.GetChannelWithRelations(dbCtx, chat.ID)
		if err != nil {
			log.Printf("Canal %d não encontrado: %v", chat.ID, err)
			return
		}

		// Atualizar info básica e primeiro botão em background
		go func() {
			updateCtx, updateCancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer updateCancel()
			updatedChannel, hasChanges := processor.UpdateChannelBasicInfo(updateCtx, chat.ID, channel)
			if !hasChanges {
				return
			}
			if err := c.ChannelRepo.UpdateChannelBasicInfoAndFirstButton(updateCtx, updatedChannel); err != nil {
				log.Printf("❌ Erro ao salvar informações do canal %d: %v", chat.ID, err)
			} else {
				log.Printf("✅ Canal %d: informações básicas e primeiro botão atualizados automaticamente", chat.ID)
			}
		}()

		messageType := processor.GetMessageType(post)
		if messageType == "" {
			return
		}

		permissions := processor.CheckPermissions(channel, messageType)
		if !permissions.CanEdit && !permissions.CanAddButtons {
			log.Printf("❌ Sem permissões para processar mensagem no canal %d", channel.ID)
			return
		}

		var finalButtons []dbmodels.Button
		if permissions.CanAddButtons {
			finalButtons = channel.Buttons
		}

		go func() {
			// Note: passa CanEdit apenas como "messageEditAllowed" para controle de fallback
			messageQueue.AddToQueue(messageType, channel, post, finalButtons, permissions.CanEdit, processor)
		}()

		// // Separator será enviado conforme tipo
		// if channel.Separator != nil && (permissions.CanEdit || permissions.CanAddButtons) {
		// 	go func() {
		// 		time.Sleep(2 * time.Second)
		// 		processor.HandleSeparator(channel, post, messageType)
		// 	}()
		// }

		// Envio de separador para mensagens individuais (não grupos)
		if channel.Separator != nil && channel.Separator.SeparatorID != "" {
			// Não enviar separador no início de álbuns; finalizadores de grupo enviarão ao final
			if post.MediaGroupID == "" {
				go func(ch *dbmodels.Channel, p *models.Message) {
					// Contexto próprio para evitar cancelamentos do handler
					ctxSep, cancelSep := context.WithTimeout(context.Background(), 10*time.Second)
					defer cancelSep()
					// Chama a função central de separador, que já faz retry/backoff
					if err := processor.ProcessSeparator(ctxSep, ch, p); err != nil {
						log.Printf("⚠️ Falha ao enviar separador: %v", err)
					}
				}(channel, post)
			}
		}

	}
}

// ✅ ProcessMessage
func (mp *MessageProcessor) ProcessMessage(ctx context.Context, messageType MessageType, channel *dbmodels.Channel, post *models.Message, buttons []dbmodels.Button, messageEditAllowed bool) error {
	switch messageType {
	case MessageTypeText:
		return mp.ProcessTextMessage(ctx, channel, post, buttons, messageEditAllowed)
	case MessageTypeAudio:
		return mp.ProcessAudioMessage(ctx, channel, post, buttons, messageEditAllowed)
	case MessageTypeSticker:
		if len(buttons) > 0 {
			return mp.ProcessStickerMessage(ctx, channel, post, buttons)
		}
		return nil
	case MessageTypePhoto, MessageTypeVideo, MessageTypeAnimation:
		return mp.ProcessMediaMessage(ctx, channel, post, buttons, messageEditAllowed)
	default:
		return nil
	}
}

// ✅ Separator helpers (sem mudanças de permissão aqui)
func (mp *MessageProcessor) HandleSeparator(channel *dbmodels.Channel, post *models.Message, messageType MessageType) {
	if channel.Separator == nil || channel.Separator.SeparatorID == "" {
		return
	}
	mediaGroupID := post.MediaGroupID
	chatID := post.Chat.ID

	// Áudio individual: enviar direto
	if messageType == MessageTypeAudio && mediaGroupID == "" {
		time.Sleep(1 * time.Second)
		mp.sendSeparatorDirect(channel, chatID)
		return
	}

	// Grupo de áudio
	if messageType == MessageTypeAudio && mediaGroupID != "" {
		mp.handleGroupSeparator(channel, mediaGroupID, chatID)
		return
	}

	// Grupos de fotos/vídeos: enviados por finishGroupProcessing
	if mediaGroupID != "" && (messageType == MessageTypePhoto || messageType == MessageTypeVideo || messageType == MessageTypeAnimation) {
		log.Printf("🔄 Separator para grupo de mídia %s será enviado via finishGroupProcessing", mediaGroupID)
		return
	}

	// Outros tipos: direto
	mp.sendSeparatorDirect(channel, chatID)
}

func (mp *MessageProcessor) sendSeparatorDirect(channel *dbmodels.Channel, chatID int64) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := mp.bot.SendSticker(ctx, &bot.SendStickerParams{
		ChatID:  chatID,
		Sticker: &models.InputFileString{Data: channel.Separator.SeparatorID},
	})
	if err != nil {
		log.Printf("❌ Erro ao enviar separator: %v", err)
	} else {
		log.Printf("✅ Separator enviado com sucesso")
	}
}

func (mp *MessageProcessor) handleGroupSeparator(channel *dbmodels.Channel, mediaGroupID string, chatID int64) {
	if _, exists := groupSeparators.LoadOrStore(mediaGroupID, true); exists {
		return
	}
	time.Sleep(3 * time.Second)
	mp.sendSeparatorDirect(channel, chatID)
	time.AfterFunc(10*time.Second, func() {
		groupSeparators.Delete(mediaGroupID)
	})
}

// Atualização de infos do canal e primeiro botão
func (mp *MessageProcessor) UpdateChannelBasicInfo(ctx context.Context, chatID int64, channel *dbmodels.Channel) (*dbmodels.Channel, bool) {
	chat, err := mp.bot.GetChat(ctx, &bot.GetChatParams{
		ChatID: chatID,
	})
	if err != nil {
		return channel, false
	}

	updated := false
	if chat.Title != "" && chat.Title != channel.Title {
		channel.Title = utils.RemoveHTMLTags(chat.Title)
		updated = true
	}
	if chat.Username != "" {
		newUsername := "@" + chat.Username
		if newUsername != channel.InviteURL {
			channel.InviteURL = newUsername
			updated = true
		}
	} else if chat.InviteLink != "" {
		if chat.InviteLink != channel.InviteURL {
			channel.InviteURL = chat.InviteLink
			updated = true
		}
	}

	if len(channel.Buttons) > 0 {
		buttonUpdated := mp.updateFirstButtonFromChannel(ctx, channel)
		if buttonUpdated {
			updated = true
		}
	}

	return channel, updated
}

// Atualizar primeiro botão com title e URL do canal
func (mp *MessageProcessor) updateFirstButtonFromChannel(ctx context.Context, channel *dbmodels.Channel) bool {
	if len(channel.Buttons) == 0 {
		return false
	}
	chat, err := mp.bot.GetChat(ctx, &bot.GetChatParams{
		ChatID: channel.ID,
	})
	if err != nil {
		return false
	}

	novoNome := fmt.Sprintf("%s", chat.Title)
	var novaURL string
	if chat.Username != "" {
		novaURL = "https://t.me/" + chat.Username
	} else if chat.InviteLink != "" {
		novaURL = chat.InviteLink
	} else {
		return false
	}

	firstButton := &channel.Buttons[0]
	if firstButton.NameButton == novoNome && firstButton.ButtonURL == novaURL {
		return false
	}

	log.Printf("🔘 Primeiro botão atualizado: '%s' → '%s' | URL: '%s' → '%s'",
		firstButton.NameButton, novoNome, firstButton.ButtonURL, novaURL)
	firstButton.NameButton = utils.RemoveHTMLTags(novoNome)
	firstButton.ButtonURL = novaURL
	return true
}
