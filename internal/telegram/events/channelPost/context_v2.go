package channelpost

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	dbmodels "github.com/leirbagxis/FreddyBot/internal/database/models"
)

// ProcessingContext holds the entire state of a message as it moves through the pipeline.
type ProcessingContext struct {
	// Original Request Data
	Ctx    context.Context
	Bot    *bot.Bot
	Update *models.Update

	// Derived Data
	MessageType MessageType
	Channel     *dbmodels.Channel
	Permissions *PermissionCheckResult

	// Transformation State
	OriginalCaption string
	FormattedText        string
	DisableLinkPreview   bool
	FinalButtons         []dbmodels.Button
	FinalKeyboard *models.InlineKeyboardMarkup

	// Media Group State (for albums)
	IsMediaGroup    bool
	MediaGroupID    string
	GroupMessages   []MediaMessage

	// Execution Control
	Pipeline     *Pipeline
	StopPipeline bool // If true, remaining stages are skipped
	Error        error
}

// NewProcessingContext creates a fresh context for an update.
func NewProcessingContext(ctx context.Context, b *bot.Bot, update *models.Update, p *Pipeline) *ProcessingContext {
	return &ProcessingContext{
		Ctx:      ctx,
		Bot:      b,
		Update:   update,
		Pipeline: p,
	}
}
