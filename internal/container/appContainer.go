package container

import (
	"context"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/leirbagxis/FreddyBot/internal/cache"
	"github.com/leirbagxis/FreddyBot/internal/core/services"
	"github.com/leirbagxis/FreddyBot/internal/database/repositories"
	"github.com/leirbagxis/FreddyBot/pkg/logger"
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
	DB  *gorm.DB
	Bot *bot.Bot

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

	// ## CACHE ## \\
	CacheService   *cache.Service
	SessionManager *cache.SessionManager
}

func NewAppContainer(db *gorm.DB, bot *bot.Bot) *AppContainer {
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

	container := &AppContainer{
		DB:  db,
		Bot: bot,

		BroadcastQueue: make(chan BroadcastJob, 10000),

		// Services
		UserService:          services.NewUserService(userRepo),
		ChannelService:       services.NewChannelService(channelRepo, userRepo, separatorRepo, cacheService, bot),
		ButtonService:        services.NewButtonService(buttonRepo, channelRepo, customCaptionRepo, cacheService),
		CaptionService:       services.NewCaptionService(channelRepo, buttonRepo, cacheService),
		PermissionsService:   services.NewPermissionsService(permissionsRepo, channelRepo, cacheService),
		CustomCaptionService: services.NewCustomCaptionService(customCaptionRepo, channelRepo, cacheService),
		SeparatorService:     services.NewSeparatorService(separatorRepo),
		VoteService:          services.NewVoteService(voteRepo),
		ServerService:        services.NewServerService(serverRepo),

		CacheService:   cacheService,
		SessionManager: cache.NewSessionManager(cacheService),
	}

	container.startBroadcastWorkers(5)
	return container
}

func (c *AppContainer) startBroadcastWorkers(workerCount int) {
	for i := 0; i < workerCount; i++ {
		go c.broadcastWorker()
	}
}

func (c *AppContainer) broadcastWorker() {
	for job := range c.BroadcastQueue {
		var keyboard [][]models.InlineKeyboardButton
		var replyMarkup *models.InlineKeyboardMarkup

		if len(job.Buttons) > 0 {
			for _, btn := range job.Buttons {
				button := models.InlineKeyboardButton{
					Text: btn.Text,
				}

				if btn.Type == "url" {
					button.URL = btn.Value
				} else if btn.Type == "callback" {
					button.CallbackData = btn.Value
				}

				keyboard = append(keyboard, []models.InlineKeyboardButton{button})
			}
			replyMarkup = &models.InlineKeyboardMarkup{
				InlineKeyboard: keyboard,
			}
		}

		var err error
		if job.ImageUrl != "" {
			params := &bot.SendPhotoParams{
				ChatID:    job.ChatID,
				Photo:     &models.InputFileString{Data: job.ImageUrl},
				Caption:   job.Text,
				ParseMode: "HTML",
			}
			if replyMarkup != nil {
				params.ReplyMarkup = replyMarkup
			}
			_, err = c.Bot.SendPhoto(context.Background(), params)
		} else {
			params := &bot.SendMessageParams{
				ChatID:    job.ChatID,
				Text:      job.Text,
				ParseMode: "HTML",
			}
			if replyMarkup != nil {
				params.ReplyMarkup = replyMarkup
			}
			_, err = c.Bot.SendMessage(context.Background(), params)
		}

		if err != nil {
			logger.Error("APP", "Erro ao enviar para %d: %v", job.ChatID, err)
			continue
		}

		// 🔥 Controle de rate limit global
		time.Sleep(35 * time.Millisecond)
	}
}
