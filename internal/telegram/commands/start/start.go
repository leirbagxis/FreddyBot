package start

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/leirbagxis/FreddyBot/internal/container"
	"github.com/leirbagxis/FreddyBot/internal/utils"
	"github.com/leirbagxis/FreddyBot/pkg/config"
	"github.com/leirbagxis/FreddyBot/pkg/parser"
)

func Handler(c *container.AppContainer) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		userID := update.Message.From.ID

		// 0. Verificar Blacklist
		user, err := c.UserRepo.GetUserById(ctx, userID)
		if err == nil && user != nil && user.IsBlacklisted {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:    update.Message.Chat.ID,
				Text:      "🚫 <b>Acesso Bloqueado</b>\n\nVocê está na nossa blacklist e seus comandos estão bloqueados.\n\nCaso tenha dúvidas ou queira solicitar a remoção, acione o comando /ouvidoria.",
				ParseMode: models.ParseModeHTML,
				ReplyParameters: &models.ReplyParameters{
					MessageID: update.Message.ID,
				},
			})
			return
		}

		// 1. Verificar se Force Join está ativado
		configData, err := c.ServerRepo.GetConfig(ctx)
		if err == nil && configData.ForceJoin {

			// Bypass para Owner e Admins
			isAdmin := false
			if userID == config.OwnerID {
				isAdmin = true
			} else {
				user, err := c.UserRepo.GetUserById(ctx, userID)
				if err == nil && user != nil && user.IsAdmin {
					isAdmin = true
				}
			}

			if !isAdmin {
				// ID do canal: -1003767126116
				// Link: https://t.me/LegendasBOTTopic
				const channelID = -1003767126116
				const inviteLink = "https://t.me/LegendasBOTTopic"

				member, err := b.GetChatMember(ctx, &bot.GetChatMemberParams{
					ChatID: channelID,
					UserID: userID,
				})

				isMember := false
				if err == nil && member != nil {
					switch member.Type {
					case "creator", "administrator", "member", "restricted":
						isMember = true
					}
				}

				if !isMember {
					kb := &models.InlineKeyboardMarkup{
						InlineKeyboard: [][]models.InlineKeyboardButton{
							{
								{Text: "📢 Entrar no Canal", URL: inviteLink},
							},
							{
								{Text: "✅ Já entrei", CallbackData: "check_subscription"},
							},
						},
					}

					b.SendMessage(ctx, &bot.SendMessageParams{
						ChatID:      update.Message.Chat.ID,
						Text:        "⚠️ <b>Acesso Restrito</b>\n\nPara utilizar este bot, você precisa estar inscrito em nosso canal oficial.\n\nClique no botão abaixo para entrar e depois tente novamente ou clique em verificar.",
						ParseMode:   models.ParseModeHTML,
						ReplyMarkup: kb,
						ReplyParameters: &models.ReplyParameters{
							MessageID: update.Message.ID,
						},
					})
					return
				}
			}
		}

		firstName := utils.RemoveHTMLTags(update.Message.From.FirstName)

		text, button := parser.GetMessage("start", map[string]string{
			"firstName": firstName,
		})

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      update.Message.Chat.ID,
			Text:        text,
			ReplyMarkup: button,
			ParseMode:   "HTML",
			ReplyParameters: &models.ReplyParameters{
				MessageID: update.Message.ID,
			},
		})
	}
}
