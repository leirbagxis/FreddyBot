package container

import (
	"context"
	"encoding/json"
	"time"

	"github.com/leirbagxis/FreddyBot/internal/cache"
	"github.com/leirbagxis/FreddyBot/internal/core/services"
	"github.com/leirbagxis/FreddyBot/internal/database"
	"github.com/leirbagxis/FreddyBot/internal/database/repositories"
	"github.com/leirbagxis/FreddyBot/pkg/logger"
	"github.com/mymmrac/telego"
	"gorm.io/gorm"
)

type BroadcastButton struct {
	Text  string
	Type  string
	Value string
}

type BroadcastJob struct {
	ChatID   int64
	Text     string
	ImageUrl string
	Buttons  []BroadcastButton
}

type AppContainer struct {
	DB        *gorm.DB
	TelegoBot *telego.Bot

	BroadcastQueue chan BroadcastJob

	// ## SERVICES ## \\
	UserService          *services.UserService
	ChannelService       *services.ChannelService
	ButtonService        *services.ButtonService
	CaptionService       *services.CaptionService
	PermissionsService   *services.PermissionsService
	CustomCaptionService *services.CustomCaptionService
	SeparatorService     *services.SeparatorService
	VoteService          *services.VoteService
	ServerService        *services.ServerService
	ChannelEventService  *services.ChannelEventService

	// ## CACHE ## \\
	CacheService   *cache.Service
	SessionManager *cache.SessionManager
}

func NewAppContainer(db *gorm.DB, telegoClient *telego.Bot) *AppContainer {
	cacheService := cache.NewService()

	// Repositories (Removed cache from repositories)
	userRepo := repositories.NewUserRepository(db)
	channelRepo := repositories.NewChannelRepository(db)
	buttonRepo := repositories.NewButtonRepository(db)
	separatorRepo := repositories.NewSeparatorRepository(db)
	voteRepo := repositories.NewVoteRepository(db)
	customCaptionRepo := repositories.NewCustomCaptionRepository(db)
	permissionsRepo := repositories.NewPermissionsRepository(db)
	serverRepo := repositories.NewServerConfigRepository(db)
	channelEventRepo := repositories.NewChannelEventRepository(db)

	container := &AppContainer{
		DB:        db,
		TelegoBot: telegoClient,

		BroadcastQueue: make(chan BroadcastJob, 10000),

		// Services
		UserService:          services.NewUserService(userRepo),
		ChannelService:       services.NewChannelService(channelRepo, userRepo, separatorRepo, cacheService, telegoClient),
		ButtonService:        services.NewButtonService(buttonRepo, channelRepo, customCaptionRepo, cacheService),
		CaptionService:       services.NewCaptionService(channelRepo, buttonRepo, cacheService),
		PermissionsService:   services.NewPermissionsService(permissionsRepo, channelRepo, cacheService),
		CustomCaptionService: services.NewCustomCaptionService(customCaptionRepo, channelRepo, cacheService),
		SeparatorService:     services.NewSeparatorService(separatorRepo),
		VoteService:          services.NewVoteService(voteRepo),
		ServerService:        services.NewServerService(serverRepo),
		ChannelEventService:  services.NewChannelEventService(channelEventRepo),

		CacheService:   cacheService,
		SessionManager: cache.NewSessionManager(cacheService),
	}

	container.syncFixedPostBuilderSession(context.Background())
	go container.ChannelEventService.CleanupOld(context.Background(), services.ChannelEventRetentionDays)
	container.startBroadcastWorkers(5)
	return container
}

func (c *AppContainer) syncFixedPostBuilderSession(ctx context.Context) {
	config, err := c.ServerService.GetConfig(ctx)
	if err != nil {
		logger.Error("APP", "Erro ao carregar PostBuilder fixo: %v", err)
		return
	}

	if config.FixedPostBuilderKey == "" {
		return
	}

	if !config.FixedPostBuilderEnabled {
		_ = c.CacheService.DeletePostBuilderSession(ctx, config.FixedPostBuilderKey)
		return
	}

	var state cache.PostBuilderState
	if err := json.Unmarshal([]byte(config.FixedPostBuilderPayload), &state); err != nil {
		logger.Warn("APP", "Payload do PostBuilder fixo inválida, restaurando padrão: %v", err)
		config.FixedPostBuilderPayload = database.DefaultFixedPostBuilderPayload()
		config.FixedPostBuilderEnabled = true
		if err := c.ServerService.SaveConfig(ctx, config); err != nil {
			logger.Error("APP", "Erro ao salvar reparo do PostBuilder fixo: %v", err)
			return
		}
		if err := json.Unmarshal([]byte(config.FixedPostBuilderPayload), &state); err != nil {
			logger.Error("APP", "Payload padrão do PostBuilder fixo inválida: %v", err)
			return
		}
	}

	if err := c.CacheService.SetPostBuilderSession(ctx, config.FixedPostBuilderKey, state, 0); err != nil {
		logger.Error("APP", "Erro ao sincronizar PostBuilder fixo no Redis: %v", err)
	}
}

func (c *AppContainer) startBroadcastWorkers(workerCount int) {
	for i := 0; i < workerCount; i++ {
		go c.broadcastWorker()
	}
}

func (c *AppContainer) broadcastWorker() {
	for job := range c.BroadcastQueue {
		var keyboard [][]telego.InlineKeyboardButton
		var replyMarkup *telego.InlineKeyboardMarkup

		if len(job.Buttons) > 0 {
			for _, btn := range job.Buttons {
				button := telego.InlineKeyboardButton{
					Text: btn.Text,
				}

				if btn.Type == "url" {
					button.URL = btn.Value
				} else if btn.Type == "callback" {
					button.CallbackData = btn.Value
				}

				keyboard = append(keyboard, []telego.InlineKeyboardButton{button})
			}
			replyMarkup = &telego.InlineKeyboardMarkup{
				InlineKeyboard: keyboard,
			}
		}

		var err error
		if job.ImageUrl != "" {
			params := &telego.SendPhotoParams{
				ChatID:    telego.ChatID{ID: job.ChatID},
				Photo:     telego.InputFile{URL: job.ImageUrl},
				Caption:   job.Text,
				ParseMode: telego.ModeHTML,
			}
			if replyMarkup != nil {
				params.ReplyMarkup = replyMarkup
			}
			_, err = c.TelegoBot.SendPhoto(context.Background(), params)
		} else {
			params := &telego.SendMessageParams{
				ChatID:    telego.ChatID{ID: job.ChatID},
				Text:      job.Text,
				ParseMode: telego.ModeHTML,
			}
			if replyMarkup != nil {
				params.ReplyMarkup = replyMarkup
			}
			_, err = c.TelegoBot.SendMessage(context.Background(), params)
		}

		if err != nil {
			logger.Error("APP", "Erro ao enviar para %d: %v", job.ChatID, err)
			continue
		}

		// 🔥 Controle de rate limit global
		time.Sleep(35 * time.Millisecond)
	}
}
