package logs

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	usermodels "github.com/leirbagxis/FreddyBot/internal/database/models"
	"github.com/leirbagxis/FreddyBot/pkg/config"
	"github.com/leirbagxis/FreddyBot/pkg/logger"
)

func LogAdmin(ctx context.Context, b *bot.Bot, channel *usermodels.Channel) {
	user, err := b.GetChat(ctx, &bot.GetChatParams{
		ChatID: channel.OwnerID,
	})
	if err != nil {
		logger.Error("ADMIN", "Erro ao buscar informacoes do usuario para o log: %v", err)
	}

	ownerName := user.FirstName
	if ownerName == "" {
		ownerName = "Desconhecido"
	}

	text := fmt.Sprintf(
		"📢 <b>Novo canal adicionado!</b>\n\n"+
			"🛰 <b>Canal:</b> %s\n"+
			"🆔 <b>ID:</b> <code>%d</code>\n"+
			"🔗 <b>Link:</b> %s\n"+
			"👤 <b>Owner:</b> %s (<code>%d</code>)\n"+
			"⏰ <b>Criado em:</b> <code>%s</code>",
		channel.Title,
		channel.ID,
		channel.InviteURL,
		ownerName,
		channel.OwnerID,
		channel.CreatedAt.Format("02/01/2006 15:04:05"),
	)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    config.OwnerID,
		Text:      text,
		ParseMode: models.ParseModeHTML,
	})
}
