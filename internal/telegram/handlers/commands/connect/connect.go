package connect

import (
	"strings"

	"github.com/leirbagxis/FreddyBot/internal/container"
	"github.com/leirbagxis/FreddyBot/pkg/config"
	"github.com/leirbagxis/FreddyBot/pkg/logger"
	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegohandler"
)

func HandlerTelego(c *container.AppContainer) telegohandler.Handler {
	return func(ctx *telegohandler.Context, update telego.Update) error {
		userID := update.Message.From.ID
		chatID := update.Message.Chat.ID

		connected, err := c.TelegramClientService.IsConnected(ctx, userID)
		if err != nil {
			connected = false
		}

		var text string
		var keyboard [][]telego.InlineKeyboardButton

		if connected {
			text = "🔗 <b>Conta Telegram conectada</b>\n\n"
			text += "Sua conta pessoal está conectada ao bot via MTProto.\n\n"
			text += "Caso precise, você pode desconectar abaixo:"

			keyboard = [][]telego.InlineKeyboardButton{
				{
					{
						Text:         "🔌 Desconectar Conta",
						CallbackData: "tgconnect:disconnect",
					},
				},
			}
		} else {
			text = "🔗 <b>Conectar Conta Telegram</b>\n\n"
			text += "Conecte sua conta pessoal do Telegram ao bot para usar recursos avançados.\n\n"
			text += "📱 <b>Como funciona:</b>\n"
			text += "1. Clique no botão abaixo para abrir o Mini App\n"
			text += "2. Digite seu número de telefone\n"
			text += "3. Insira o código enviado pelo Telegram\n"
			text += "4. Se tiver 2FA ativado, digite sua senha\n\n"
			text += "🔒 Seus dados são criptografados e armazenados com segurança."

			webAppURL := strings.TrimRight(config.WebAppURL, "/") + "/connect"
			keyboard = [][]telego.InlineKeyboardButton{
				{
					{
						Text:   "🔗 Conectar Conta",
						WebApp: &telego.WebAppInfo{URL: webAppURL},
					},
				},
			}
		}

		params := &telego.SendMessageParams{
			ChatID:      telego.ChatID{ID: chatID},
			Text:        text,
			ParseMode:   telego.ModeHTML,
			ReplyMarkup: &telego.InlineKeyboardMarkup{InlineKeyboard: keyboard},
		}

		if _, err := c.TelegoBot.SendMessage(ctx, params); err != nil {
			logger.Error("CONNECT", "Error sending connect message to %d: %v", userID, err)
		}

		return nil
	}
}

func DisconnectCallbackHandlerTelego(c *container.AppContainer) telegohandler.Handler {
	return func(ctx *telegohandler.Context, update telego.Update) error {
		userID := update.CallbackQuery.From.ID
		chatID := update.CallbackQuery.Message.GetChat().ID

		if err := c.TelegramClientService.DisconnectUser(ctx, userID); err != nil {
			logger.Error("CONNECT", "Error disconnecting user %d: %v", userID, err)
			_ = c.TelegoBot.AnswerCallbackQuery(ctx, &telego.AnswerCallbackQueryParams{
				CallbackQueryID: update.CallbackQuery.ID,
				Text:            "❌ Erro ao desconectar",
				ShowAlert:       true,
			})
			return nil
		}

		_ = c.TelegoBot.AnswerCallbackQuery(ctx, &telego.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
			Text:            "✅ Conta desconectada",
		})

		_, _ = c.TelegoBot.EditMessageText(ctx, &telego.EditMessageTextParams{
			ChatID:    telego.ChatID{ID: chatID},
			MessageID: update.CallbackQuery.Message.GetMessageID(),
			Text:      "🔌 Sua conta foi desconectada do bot.\n\nUse /connect para conectar novamente.",
			ParseMode: telego.ModeHTML,
		})

		return nil
	}
}
