package telegram

import (
	"context"
	"strings"

	"github.com/leirbagxis/FreddyBot/internal/container"
	"github.com/leirbagxis/FreddyBot/internal/middleware"
	"github.com/leirbagxis/FreddyBot/internal/telegram/events/channelPost"
	callbackAbout "github.com/leirbagxis/FreddyBot/internal/telegram/handlers/callbacks/about"
	callbackClaim "github.com/leirbagxis/FreddyBot/internal/telegram/handlers/callbacks/claimChannel"
	callbackMyChannel "github.com/leirbagxis/FreddyBot/internal/telegram/handlers/callbacks/my_channel"
	callbackProfile "github.com/leirbagxis/FreddyBot/internal/telegram/handlers/callbacks/profile_info"
	callbackStart "github.com/leirbagxis/FreddyBot/internal/telegram/handlers/callbacks/start"
	callbackVote "github.com/leirbagxis/FreddyBot/internal/telegram/handlers/callbacks/vote"
	"github.com/leirbagxis/FreddyBot/internal/telegram/handlers/commands/admin"
	commandConnect "github.com/leirbagxis/FreddyBot/internal/telegram/handlers/commands/connect"
	"github.com/leirbagxis/FreddyBot/internal/telegram/handlers/commands/help"
	commandStart "github.com/leirbagxis/FreddyBot/internal/telegram/handlers/commands/start"
	"github.com/leirbagxis/FreddyBot/internal/telegram/handlers/commands/suporte"
	"github.com/leirbagxis/FreddyBot/internal/telegram/handlers/commands/tutorial"
	"github.com/leirbagxis/FreddyBot/internal/telegram/handlers/events/addChannel"
	"github.com/leirbagxis/FreddyBot/internal/telegram/handlers/events/postBuilder"
	"github.com/leirbagxis/FreddyBot/pkg/config"
	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegohandler"
)

func LoadHandlersTelegoWithBH(bh *telegohandler.BotHandler, c *container.AppContainer) {
	// Middlewares
	bh.Use(middleware.SaveUserMiddlewareTelego(c))
	bh.Use(middleware.CheckBlacklistMiddlewareTelego(c))
	bh.Use(middleware.CheckMaintenanceMiddlewareTelego(c))

	// Channel Post Handler
	bh.Handle(channelpost.HandlerTelego(c), telegohandler.AnyChannelPost())

	// Add Channel Handlers
	addChannelGroup := bh.Group(telegohandler.AnyMyChatMember())
	addChannelGroup.Use(middleware.CheckAddBotMiddlewareTelego(c))
	addChannelGroup.Handle(addchannel.AskAddChannelHandlerTelego(c))

	bh.Handle(addchannel.UpdateChannelInfoHandlerTelego(c), telegohandler.AnyMyChatMember())

	forwardedGroup := bh.Group(matchForwardedChannelTelego())
	forwardedGroup.Use(middleware.CheckAddBotMiddlewareTelego(c))
	forwardedGroup.Handle(addchannel.AskAddChannelHandlerTelego(c))

	// Commands
	bh.Handle(commandStart.HandlerTelego(c), telegohandler.CommandEqual("start"))
	bh.Handle(commandConnect.HandlerTelego(c), telegohandler.CommandEqual("connect"))
	bh.Handle(help.HandlerTelego(c), telegohandler.CommandEqual("help"))
	bh.Handle(suporte.HandlerTelego(c), telegohandler.CommandEqual("ouvidoria"))
	bh.Handle(tutorial.HandlerTelego(c), telegohandler.CommandEqual("tutorial"))

	// Admin Commands (Owner only; /info is owner/admin below)
	adminGroup := bh.Group(matchOwnerTelego())
	adminGroup.Handle(admin.AdminHelpHandlerTelego(c), telegohandler.CommandEqual("admin"))
	adminGroup.Handle(admin.GetAllUsersHandlerTelego(c), telegohandler.CommandEqual("users"))
	adminGroup.Handle(admin.GetAllChannelsHandlerTelego(c), telegohandler.CommandEqual("channels"))
	adminGroup.Handle(admin.GetInfoUserHandlerTelego(c), telegohandler.CommandEqual("user"))
	adminGroup.Handle(admin.NoticeCommandHandlerTelego(c), telegohandler.CommandEqual("notice"))
	adminGroup.Handle(admin.NoticeChannelsHandlerTelego(c), telegohandler.CommandEqual("publi"))
	adminGroup.Handle(admin.SendMessageToIdHandlerTelego(c), telegohandler.CommandEqual("send"))
	adminGroup.Handle(admin.NoticeUsersReplyHandlerTelego(c), telegohandler.CommandEqual("allusers"))
	adminGroup.Handle(admin.NoticeChannelsReplyHandlerTelego(c), telegohandler.CommandEqual("allchannels"))
	adminGroup.Handle(admin.AddChannelCommandHandlerTelego(c), telegohandler.CommandEqual("add"))
	adminGroup.Handle(admin.RemoveChannelHandlerTelego(c), telegohandler.CommandEqual("remove"))
	adminGroup.Handle(admin.RegisterTransferHandlerTelego(c), telegohandler.CommandEqual("transfer"))
	adminGroup.Handle(admin.SetAdminHandlerTelego(c), telegohandler.CommandEqual("setadmin"))
	adminGroup.Handle(admin.ToggleMaintenceHandlerTelego(c), telegohandler.CommandEqual("maintence"))
	adminGroup.Handle(admin.GetBackUpHandlerTelego(c), telegohandler.CommandEqual("backup"))
	adminGroup.Handle(admin.CheckBotAdminHandlerTelego(c), telegohandler.CommandEqual("checkbot"))
	adminGroup.Handle(admin.GetMediaIDHandlerTelego(c), telegohandler.CommandEqual("getid"))

	adminOrOwnerGroup := bh.Group(matchAdminOrOwnerTelego(c))
	adminOrOwnerGroup.Handle(admin.GetInfoChannelHandlerTelego(c), telegohandler.CommandEqual("info"))

	// Message Handlers for active sessions (Text and Sticker inputs)
	bh.Handle(callbackMyChannel.SetStickerSeparatorHandlerTelego(c), matchAwaitingStickerSeparatorTelego(c))
	bh.Handle(callbackMyChannel.SetTransferAccessHandlerTelego(c), matchAwaitingTransferAccessTelego(c))

	// Post Builder - Message Handler (Media and Text Input)
	bh.Handle(postbuilder.HandlerTelego(c), matchPostBuilderTelego(c))

	// Callbacks
	bh.Handle(callbackStart.HandlerTelego(c), telegohandler.CallbackDataEqual("start"))
	bh.Handle(callbackStart.CheckSubscriptionHandlerTelego(c), telegohandler.CallbackDataEqual("check_subscription"))
	bh.Handle(callbackProfile.HandlerTelego(c), telegohandler.CallbackDataEqual("profile-info"))
	bh.Handle(callbackMyChannel.HandlerTelego(c), telegohandler.CallbackDataEqual("profile-user-channels"))
	bh.Handle(callbackMyChannel.ConfigHandlerTelego(c), telegohandler.CallbackDataPrefix("config:"))
	bh.Handle(callbackMyChannel.GroupChannelHandlerTelego(c), telegohandler.CallbackDataPrefix("gc-info:"))
	bh.Handle(callbackVote.HandlerTelego(c), telegohandler.CallbackDataPrefix("vote:"))
	bh.Handle(callbackMyChannel.AskDeleteChannelHandlerTelego(c), telegohandler.CallbackDataEqual("del"))
	bh.Handle(callbackMyChannel.ConfirmDeleteChannelHandlerTelego(c), telegohandler.CallbackDataPrefix("confirm-del:"))

	// Remaining Callbacks
	bh.Handle(callbackAbout.HandlerTelego(c), telegohandler.CallbackDataEqual("about"))
	bh.Handle(callbackClaim.AcceptClaimHandlerTelego(c), telegohandler.CallbackDataPrefix("accept-claim:"))

	// Sticker Separator Callbacks
	bh.Handle(callbackMyChannel.AskStickerSeparatorHandlerTelego(c), telegohandler.CallbackDataEqual("sptc"))
	bh.Handle(callbackMyChannel.RequireStickerSeparatorHandlerTelego(c), telegohandler.CallbackDataEqual("sptc-config"))
	bh.Handle(callbackMyChannel.DeleteSeparatorHandlerTelego(c), telegohandler.CallbackDataEqual("spex"))

	// Transfer Access Callbacks
	bh.Handle(callbackMyChannel.AskTransferAccessHandlerTelego(c), telegohandler.CallbackDataEqual("paccess-info"))
	bh.Handle(callbackMyChannel.TransferAcessHandlerTelego(c), telegohandler.CallbackDataEqual("transfer"))

	// Help Callback
	bh.Handle(help.CallbackHandlerTelego(c), telegohandler.CallbackDataEqual("help"))

	// Telegram Connect Callbacks
	bh.Handle(commandConnect.DisconnectCallbackHandlerTelego(c), telegohandler.CallbackDataEqual("tgconnect:disconnect"))

	// Add Channel Callbacks
	bh.Handle(addchannel.AddYesHandlerTelego(c), telegohandler.CallbackDataPrefix("add-yes:"))
	bh.Handle(addchannel.AddNotHandlerTelego(c), telegohandler.CallbackDataPrefix("add-not:"))

	// Post Builder Callbacks
	bh.Handle(postbuilder.CallbackHandlerTelego(c), telegohandler.CallbackDataPrefix("pb-"))

	// Inline Handlers
	bh.HandleInlineQuery(postbuilder.InlineHandlerTelego(c), telegohandler.InlineQueryPrefix("pb "))
	bh.HandleInlineQuery(callbackClaim.HandlerTelego(c), telegohandler.InlineQueryPrefix("claim "))
	bh.HandleChosenInlineResult(postbuilder.ChosenInlineResultHandlerTelego(c), telegohandler.AnyChosenInlineResult())
}

func LoadHandlersTelego(bot *telego.Bot, c *container.AppContainer) *telegohandler.BotHandler {
	updates, _ := bot.UpdatesViaLongPolling(context.Background(), nil)
	bh, _ := telegohandler.NewBotHandler(bot, updates)
	LoadHandlersTelegoWithBH(bh, c)
	return bh
}

func matchAwaitingStickerSeparatorTelego(c *container.AppContainer) telegohandler.Predicate {
	return func(ctx context.Context, update telego.Update) bool {
		if update.Message == nil || update.Message.From == nil {
			return false
		}
		id, _ := c.CacheService.GetAwaitingStickerSeparator(context.Background(), update.Message.From.ID)
		return id != 0 && update.Message.Sticker != nil
	}
}

func matchAwaitingTransferAccessTelego(c *container.AppContainer) telegohandler.Predicate {
	return func(ctx context.Context, update telego.Update) bool {
		if update.Message == nil || update.Message.From == nil {
			return false
		}
		id, err := c.CacheService.GetTransferChannel(context.Background(), update.Message.From.ID)
		return err == nil && id != 0 && update.Message.Text != ""
	}
}

func matchForwardedChannelTelego() telegohandler.Predicate {
	return func(ctx context.Context, update telego.Update) bool {
		if update.Message == nil || update.Message.ForwardOrigin == nil {
			return false
		}
		_, ok := update.Message.ForwardOrigin.(*telego.MessageOriginChannel)
		return ok
	}
}

func matchOwnerTelego() telegohandler.Predicate {
	return func(ctx context.Context, update telego.Update) bool {
		if update.Message == nil || update.Message.From == nil {
			return false
		}
		return update.Message.From.ID == config.OwnerID
	}
}

func matchAdminOrOwnerTelego(c *container.AppContainer) telegohandler.Predicate {
	return func(ctx context.Context, update telego.Update) bool {
		if update.Message == nil || update.Message.From == nil {
			return false
		}

		userID := update.Message.From.ID
		if userID == config.OwnerID {
			return true
		}

		user, err := c.UserService.GetUserByID(context.Background(), userID)
		return err == nil && user != nil && user.IsAdmin
	}
}

func matchPostBuilderTelego(c *container.AppContainer) telegohandler.Predicate {
	return func(ctx context.Context, update telego.Update) bool {
		if update.Message == nil || update.Message.Chat.Type != telego.ChatTypePrivate {
			return false
		}
		// Se for comando, deixa os comandos tratarem
		if strings.HasPrefix(update.Message.Text, "/") {
			return false
		}

		// Prioridade para sessões ativas
		userId := update.Message.From.ID
		if id, _ := c.CacheService.GetAwaitingStickerSeparator(context.Background(), userId); id != 0 {
			return false
		}
		if id, _ := c.CacheService.GetTransferChannel(context.Background(), userId); id != 0 {
			return false
		}

		// Match if it has media
		if update.Message.Photo != nil || update.Message.Video != nil || update.Message.Animation != nil || update.Message.Audio != nil || update.Message.Document != nil || update.Message.Sticker != nil {
			return true
		}
		// Match if in active session for text input
		state, _ := c.CacheService.GetPostBuilderState(context.Background(), update.Message.From.ID)
		return state != nil && state.Step != ""
	}
}
