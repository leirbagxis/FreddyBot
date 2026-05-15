package channelpost

import (
	"strings"

	"github.com/go-telegram/bot/models"
	"github.com/leirbagxis/FreddyBot/internal/container"
	"github.com/leirbagxis/FreddyBot/pkg/logger"
)

func StageDecorate(c *container.AppContainer) Stage {
	return func(pCtx *ProcessingContext) error {
		if pCtx.Channel == nil {
			return nil
		}

		rows := map[int][]models.InlineKeyboardButton{}

		// 1. Add Buttons
		for _, b := range pCtx.FinalButtons {
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

		// 2. Add Reactions
		if pCtx.Channel.Reactions != "" && pCtx.Permissions.CanAddReactions {
			reactions := strings.Split(pCtx.Channel.Reactions, ",")
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
				row := pCtx.Channel.ReactionPosition
				if row < 0 {
					row = 0
				}

				// Avoid conflict with existing buttons
				hasConflict := false
				maxBtnRow := -1
				for _, b := range pCtx.FinalButtons {
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

		// 3. Assemble and Sort Keyboard
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

		if len(keyboard) > 0 {
			pCtx.FinalKeyboard = &models.InlineKeyboardMarkup{InlineKeyboard: keyboard}
			logger.Bot("🎹 Teclado construído: %d linhas", len(keyboard))
		} else {
			logger.Bot("🎹 Nenhum botão ou reação a adicionar")
		}

		return nil
	}
}
