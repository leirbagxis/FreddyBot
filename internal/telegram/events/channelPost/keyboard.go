package channelpost

import (
	"strings"

	"github.com/go-telegram/bot/models"
	dbmodels "github.com/leirbagxis/FreddyBot/internal/database/models"
)

func CreateInlineKeyboard(buttons []dbmodels.Button, customCaption *dbmodels.CustomCaption, channel *dbmodels.Channel, messageType MessageType) *models.InlineKeyboardMarkup {
	var finalButtons []dbmodels.Button
	pm := GetPermissionManager()

	if customCaption != nil && len(customCaption.Buttons) > 0 {
		for _, cb := range customCaption.Buttons {
			finalButtons = append(finalButtons, dbmodels.Button{
				NameButton: cb.NameButton,
				ButtonURL:  cb.ButtonURL,
				PositionY:  cb.PositionY,
				PositionX:  cb.PositionX,
			})
		}
	} else {
		perms := pm.CheckPermissions(channel, messageType)
		if !perms.CanAddButtons {
			return nil
		}
		finalButtons = buttons
	}

	if len(finalButtons) == 0 {
		return nil
	}

	// Construção por linhas
	rows := map[int][]models.InlineKeyboardButton{}

	// Adicionar botões
	for _, b := range finalButtons {
		if b.NameButton == "" || b.ButtonURL == "" {
			continue
		}
		row := b.PositionY
		if row < 0 {
			row = 0
		}
		btn := models.InlineKeyboardButton{Text: b.NameButton, URL: b.ButtonURL}
		rows[row] = append(rows[row], btn)
	}

	// Adicionar reações (se houver) na posição correta
	perms := pm.CheckPermissions(channel, messageType)
	if channel != nil && channel.Reactions != "" && perms.CanAddReactions {
		reactions := strings.Split(channel.Reactions, ",")
		var reactionRow []models.InlineKeyboardButton
		for _, r := range reactions {
			emoji := strings.TrimSpace(r)
			if emoji != "" {
				reactionRow = append(reactionRow, models.InlineKeyboardButton{
					Text:         emoji,
					CallbackData: "vote:" + emoji,
				})
			}
		}
		if len(reactionRow) > 0 {
			row := channel.ReactionPosition
			if row < 0 {
				row = 0
			}

			// Evitar conflito: se a linha das reações já tiver botões, joga para a próxima linha disponível
			hasConflict := false
			maxBtnRow := -1
			for _, b := range finalButtons {
				if b.PositionY > maxBtnRow {
					maxBtnRow = b.PositionY
				}
				if b.PositionY == row {
					hasConflict = true
				}
			}

			if hasConflict {
				row = maxBtnRow + 1
			}

			rows[row] = append(rows[row], reactionRow...)
		}
	}

	// Ordenar por linha
	keyboard := make([][]models.InlineKeyboardButton, 0)
	maxRow := 0
	for r := range rows {
		if r > maxRow {
			maxRow = r
		}
	}

	for r := 0; r <= maxRow; r++ {
		if line, ok := rows[r]; ok && len(line) > 0 {
			keyboard = append(keyboard, line)
		}
	}

	if len(keyboard) == 0 {
		return nil
	}
	return &models.InlineKeyboardMarkup{InlineKeyboard: keyboard}
}
