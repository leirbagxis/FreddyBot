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

	ServerRepo *repositories.ServerConfig

	UserRepo      *repositories.UserRepository
	ChannelRepo   *repositories.ChannelRepository
	ButtonRepo    *repositories.ButtonRepository
	SeparatorRepo *repositories.SeparatorRepository
	VoteRepo      *repositories.VoteRepository

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

		ServerRepo: repositories.NewServerConfigRepository(db, cacheService),

		BroadcastQueue: make(chan BroadcastJob, 10000),

		UserRepo:      repositories.NewUserRepository(db),
		ChannelRepo:   repositories.NewChannelRepository(db, cacheService),
		ButtonRepo:    repositories.NewButtonRepository(db),
		SeparatorRepo: repositories.NewSeparatorRepository(db),
		VoteRepo:      repositories.NewVoteRepository(db),

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

func (c *AppContainer) DisconnectChannel(ctx context.Context, userID int64, channelID int64) error {
	// 1. Send farewell message to the channel
	farewellMsg := "Ah, então é assim? Um clique e tudo o que vivemos vira fumaça. Não se preocupe, eu vou embora... mas saiba que meu silêncio será o seu maior arrependimento. Aproveite sua liberdade sem mim. Adeus, ingrato! 🍷"

	// Tenta enviar a mensagem, mas não bloqueia se falhar (o bot pode já ter sido removido manualmente)
	_, _ = c.Bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: channelID,
		Text:   farewellMsg,
	})

	// 2. Leave the channel
	_, _ = c.Bot.LeaveChat(ctx, &bot.LeaveChatParams{
		ChatID: channelID,
	})

	// 3. Delete from DB
	err := c.ChannelRepo.DeleteChannelWithRelations(ctx, userID, channelID)
	if err != nil {
		return err
	}

	// 4. Clean cache
	c.CacheService.SetDeleteChannel(ctx, userID, channelID)

	return nil
}
