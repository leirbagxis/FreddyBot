package container

import (
	"context"
	"log"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/leirbagxis/FreddyBot/internal/cache"
	"github.com/leirbagxis/FreddyBot/internal/database/repositories"
	adminrepositories "github.com/leirbagxis/FreddyBot/internal/database/repositories/adminRepositories"
	"gorm.io/gorm"
)

type BroadcastButton struct {
	Text  string
	Type  string
	Value string
}

type BroadcastJob struct {
	ChatID  int64
	Text    string
	Buttons []BroadcastButton
}

type AppContainer struct {
	DB  *gorm.DB
	Bot *bot.Bot

	BroadcastQueue chan BroadcastJob

	UserRepo      *repositories.UserRepository
	ChannelRepo   *repositories.ChannelRepository
	ButtonRepo    *repositories.ButtonRepository
	SeparatorRepo *repositories.SeparatorRepository

	// ## CACHE ## \\
	CacheService   *cache.Service
	SessionManager *cache.SessionManager

	// ### ADMIN ### \\
	AdminService *adminrepositories.AdminRepositories
}

func NewAppContainer(db *gorm.DB, bot *bot.Bot) *AppContainer {
	cacheService := cache.NewService()
	container := &AppContainer{
		DB:  db,
		Bot: bot,

		BroadcastQueue: make(chan BroadcastJob, 10000),

		UserRepo:      repositories.NewUserRepository(db),
		ChannelRepo:   repositories.NewChannelRepository(db),
		ButtonRepo:    repositories.NewButtonRepository(db),
		SeparatorRepo: repositories.NewSeparatorRepository(db),

		CacheService:   cacheService,
		SessionManager: cache.NewSessionManager(cacheService),

		AdminService: adminrepositories.NewAdminRepositories(db),
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
		params := &bot.SendMessageParams{
			ChatID:    job.ChatID,
			Text:      job.Text,
			ParseMode: "HTML",
		}

		if len(job.Buttons) > 0 {
			var keyboard [][]models.InlineKeyboardButton
			for _, btn := range job.Buttons {
				button := models.InlineKeyboardButton{
					Text: btn.Text,
				}

				if btn.Type == "url" {
					button.URL = btn.Value
				} else if btn.Type == "callback" {
					button.CallbackData = btn.Value
				}

				// Adiciona o botão em uma nova linha (layout vertical)
				// Se quiser horizontal, teria que agrupar diferente
				keyboard = append(keyboard, []models.InlineKeyboardButton{button})
			}
			params.ReplyMarkup = &models.InlineKeyboardMarkup{
				InlineKeyboard: keyboard,
			}
		}

		_, err := c.Bot.SendMessage(
			context.Background(),
			params,
		)

		if err != nil {
			log.Printf("erro ao enviar para %d: %v\n", job.ChatID, err)
			continue
		}

		// 🔥 Controle de rate limit global
		time.Sleep(35 * time.Millisecond)
	}
}
