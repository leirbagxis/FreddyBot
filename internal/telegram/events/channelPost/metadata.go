package channelpost

import (
	"context"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	dbmodels "github.com/leirbagxis/FreddyBot/internal/database/models"
	"github.com/leirbagxis/FreddyBot/internal/utils"
	"github.com/leirbagxis/FreddyBot/pkg/logger"
)

// Atualização de infos do canal e primeiro botão
func UpdateChannelBasicInfo(ctx context.Context, b *bot.Bot, chatID int64, channel *dbmodels.Channel, chatObj interface{}) (*dbmodels.Channel, bool) {
	var title, username, inviteLink string

	if chatObj != nil {
		switch c := chatObj.(type) {
		case *models.Chat:
			title = c.Title
			username = c.Username
		case models.Chat:
			title = c.Title
			username = c.Username
		case *models.ChatFullInfo:
			title = c.Title
			username = c.Username
			inviteLink = c.InviteLink
		case models.ChatFullInfo:
			title = c.Title
			username = c.Username
			inviteLink = c.InviteLink
		}
	} else {
		chat, err := b.GetChat(ctx, &bot.GetChatParams{
			ChatID: chatID,
		})
		if err != nil {
			return channel, false
		}
		title = chat.Title
		username = chat.Username
		inviteLink = chat.InviteLink
	}

	updated := false
	if title != "" {
		cleanTitle := utils.RemoveHTMLTags(title)
		if cleanTitle != channel.Title {
			channel.Title = cleanTitle
			updated = true
		}
	}

	// Lógica de URL: Username Público (@) sempre tem prioridade sobre link privado (t.me/+)
	var newInviteURL string
	if username != "" {
		newInviteURL = "@" + username
	} else if inviteLink != "" {
		newInviteURL = inviteLink
	}

	oldInviteURL := channel.InviteURL
	if newInviteURL != "" && newInviteURL != oldInviteURL {
		channel.InviteURL = newInviteURL
		updated = true
	}

	// Atualizar botões que apontam para o link antigo do canal
	if updated && len(channel.Buttons) > 0 {
		buttonsUpdated := syncChannelButtons(ctx, channel, title, newInviteURL, oldInviteURL)
		if buttonsUpdated {
			updated = true
		}
	}

	return channel, updated
}

// syncChannelButtons procura por qualquer botão que aponte para o link antigo do canal
// e o atualiza para o novo nome/url, respeitando a escolha do usuário de ter ou não esse botão.
func syncChannelButtons(ctx context.Context, channel *dbmodels.Channel, title, newURL, oldURL string) bool {
	if newURL == "" {
		return false
	}

	// Formata URLs para comparação (com e sem https://t.me/)
	cleanOldURL := strings.TrimPrefix(oldURL, "@")
	cleanOldURL = strings.TrimPrefix(cleanOldURL, "https://t.me/")

	cleanNewURL := newURL
	if strings.HasPrefix(newURL, "@") {
		cleanNewURL = "https://t.me/" + strings.TrimPrefix(newURL, "@")
	}

	hasChanges := false
	for i := range channel.Buttons {
		btn := &channel.Buttons[i]
		btnURL := strings.TrimPrefix(btn.ButtonURL, "https://t.me/")

		// Se o botão aponta para o link antigo (ou se é o link formatado)
		if btnURL == cleanOldURL || btn.ButtonURL == oldURL {
			if btn.NameButton != title || btn.ButtonURL != cleanNewURL {
				logger.Bot("🔘 Sincronizando botão '%s': '%s' -> '%s'", btn.NameButton, btn.ButtonURL, cleanNewURL)
				btn.NameButton = utils.RemoveHTMLTags(title)
				btn.ButtonURL = cleanNewURL
				hasChanges = true
			}
		}
	}

	return hasChanges
}
