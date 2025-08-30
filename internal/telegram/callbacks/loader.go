package callbacks

import (
	"fmt"
	"strconv"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/leirbagxis/FreddyBot/internal/container"
	"github.com/leirbagxis/FreddyBot/internal/telegram/callbacks/about"
	claimchannel "github.com/leirbagxis/FreddyBot/internal/telegram/callbacks/claimChannel"
	"github.com/leirbagxis/FreddyBot/internal/telegram/callbacks/help"
	mychannel "github.com/leirbagxis/FreddyBot/internal/telegram/callbacks/my_channel"
	profileinfo "github.com/leirbagxis/FreddyBot/internal/telegram/callbacks/profile_info"
	"github.com/leirbagxis/FreddyBot/internal/telegram/callbacks/start"
)

func LoadCallbacksHandlers(b *bot.Bot, c *container.AppContainer) {
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "help", bot.MatchTypeExact, help.Handler())
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "start", bot.MatchTypeExact, start.Handler())
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "about", bot.MatchTypeExact, about.Handler())
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "profile-info", bot.MatchTypeExact, profileinfo.Handler(c))
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "profile-user-channels", bot.MatchTypeExact, mychannel.Handler(c))

	// ## MY CHANNEL HANDLERS ## \\
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "config:", bot.MatchTypePrefix, mychannel.ConfigHandler(c))

	// DELETE CHANNEL
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "del", bot.MatchTypeExact, mychannel.AskDeleteChannelHandler(c))
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "confirm-del:", bot.MatchTypePrefix, mychannel.ConfirmDeleteChannelHandler(c))

	// STICKER SEPARADOR
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "sptc", bot.MatchTypeExact, mychannel.AskStickerSeparatorHandler(c))
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "sptc-config", bot.MatchTypeExact, mychannel.RequireStickerSeparatorHandler(c))
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "spex", bot.MatchTypeExact, mychannel.DeleteSeparatorHandler(c))

	// TRANSFER ACCES
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "paccess-info", bot.MatchTypeExact, mychannel.AskTransferAccessHandler(c))
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "transfer", bot.MatchTypeExact, mychannel.TransferAcessHandler(c))

	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "gc-info:", bot.MatchTypePrefix, mychannel.GroupChannelHandler(c))

	// CHECK MATCH
	b.RegisterHandlerMatchFunc(matchAwaitingSticker, mychannel.SetStickerSeparatorHandler(c))
	b.RegisterHandlerMatchFunc(matchAwaitingNewOwner, mychannel.SetTransferAccessHandler(c))

	b.RegisterHandlerMatchFunc(matchAwaitClaimOwner, claimchannel.Handler(c))
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "accept.claim:", bot.MatchTypePrefix, claimchannel.AcceptClaimHandler(c))
}

func matchAwaitingSticker(update *models.Update) bool {
	fmt.Println("Checking AwaitSticker: ", update.Message != nil && update.Message.From != nil && !update.Message.From.IsBot && update.Message.Sticker != nil)
	return update.Message != nil && update.Message.From != nil && !update.Message.From.IsBot && update.Message.Sticker != nil
}

func matchAwaitingNewOwner(update *models.Update) bool {
	if update.Message != nil && update.Message.From != nil && !update.Message.From.IsBot && update.Message.Text != "" {
		_, err := strconv.Atoi(update.Message.Text)
		if err == nil {
			return true
		}
	}
	return false
}

func matchAwaitClaimOwner(update *models.Update) bool {
	return update.InlineQuery != nil
}
