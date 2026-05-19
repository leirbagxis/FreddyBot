package controllers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mymmrac/telego"
	"github.com/leirbagxis/FreddyBot/internal/api/auth"
	"github.com/leirbagxis/FreddyBot/internal/api/dto"
	"github.com/leirbagxis/FreddyBot/internal/api/types"
	"github.com/leirbagxis/FreddyBot/internal/container"
	"github.com/leirbagxis/FreddyBot/pkg/errors"
	"github.com/leirbagxis/FreddyBot/pkg/logger"
	"github.com/leirbagxis/FreddyBot/pkg/parser"
)

type UserController struct {
	container *container.AppContainer
}

func NewUserController(container *container.AppContainer) *UserController {
	return &UserController{
		container: container,
	}
}

func (c *UserController) GetUserChannelsController(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.Error(errors.ErrUnauthorized)
		return
	}

	channels, err := c.container.ChannelService.GetUserChannels(ctx, userID.(int64))
	if err != nil {
		ctx.Error(err)
		return
	}

	var dtos []dto.ChannelDTO
	for _, ch := range channels {
		dtos = append(dtos, dto.ToChannelDTO(&ch))
	}

	ctx.JSON(http.StatusOK, types.NewSuccessResponse(dtos))
}

func (c *UserController) GetUserInfo(ctx *gin.Context) {
	userParams := ctx.Param("userParams")
	if len(userParams) < 5 {
		ctx.Error(errors.BadRequest("ID ou Username inválido!"))
		return
	}

	userID, _ := strconv.ParseInt(userParams, 10, 64)
	if userID == 0 {
		user, err := c.container.UserService.GetUserByUsername(context.Background(), userParams)
		if err != nil {
			ctx.Error(err)
			return
		}

		ctx.JSON(http.StatusOK, types.NewSuccessResponse(dto.ToUserDTO(user)))
		return

	}

	user, err := c.container.TelegoBot.GetChat(context.Background(), &telego.GetChatParams{
		ChatID: telego.ChatID{ID: userID},
	})
	if err != nil {
		ctx.Error(errors.BadRequest("Usuario nao encontrado"))
		return
	}

	ctx.JSON(http.StatusOK, types.NewSuccessResponse(user))
}

func (c *UserController) TransferChannelController(ctx *gin.Context) {
	var body *types.TransferChannelRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.Error(errors.BadRequest("payload inválido: " + err.Error()))
		return
	}

	channel, err := c.container.ChannelService.GetChannelByID(ctx, body.ChannelID)
	if err != nil {
		ctx.Error(err)
		return
	}

	if body.NewOwnerID == body.OldOwnerID {
		ctx.Error(errors.BadRequest("O novo dono precisa ser diferente de voce."))
		return
	}

	newOwner, err := c.container.TelegoBot.GetChat(ctx, &telego.GetChatParams{ChatID: telego.ChatID{ID: body.NewOwnerID}})
	if err != nil {
		logger.Error("API", "Erro ao obter chat do novo dono: %v", err)
		ctx.Error(errors.BadRequest("O novo dono precisa iniciar o bot pelo menos uma vez."))
		return
	}

	// Verifica se o novo dono é um bot
	botInfo, _ := c.container.TelegoBot.GetMe(context.Background())
	if body.NewOwnerID == botInfo.ID {
		ctx.Error(errors.BadRequest("O novo dono não pode ser eu."))
		return
	}

	admins, err := c.container.TelegoBot.GetChatAdministrators(ctx, &telego.GetChatAdministratorsParams{
		ChatID: telego.ChatID{ID: body.ChannelID},
	})
	if err != nil {
		logger.Error("API", "Erro ao buscar administradores do canal: %v", err)
		ctx.Error(errors.New(http.StatusInternalServerError, "Erro ao consultar administradores do canal."))
		return
	}

	isAdmin := false
	for _, admin := range admins {
		status := admin.MemberStatus()
		if status == telego.MemberStatusAdministrator {
			if a, ok := admin.(*telego.ChatMemberAdministrator); ok && a.User.ID == body.NewOwnerID {
				isAdmin = true
				break
			}
		}
		if status == telego.MemberStatusCreator {
			if a, ok := admin.(*telego.ChatMemberOwner); ok && a.User.ID == body.NewOwnerID {
				isAdmin = true
				break
			}
		}
	}

	if !isAdmin {
		ctx.Error(errors.BadRequest("O novo dono precisa ser administrador do canal."))
		return
	}

	err = c.container.ChannelService.TransferChannel(ctx, body.ChannelID, body.OldOwnerID, body.NewOwnerID)
	if err != nil {
		ctx.Error(err)
		return
	}

	channelID := fmt.Sprintf("%d", body.ChannelID)
	newOwnerIDStr := fmt.Sprintf("%d", body.NewOwnerID)

	data := map[string]string{
		"channelId":    channelID,
		"channelName":  channel.Title,
		"newOwnerName": newOwner.LastName,
		"newOwnerId":   newOwnerIDStr,
		"miniAppUrl":   auth.GenerateMiniAppUrl(newOwnerIDStr, channelID),
	}

	textOld, buttonOld := parser.GetMessageTelego("success-old-paccess-message", data)
	paramsOld := &telego.SendMessageParams{
		ChatID:    telego.ChatID{ID: body.OldOwnerID},
		Text:      textOld,
		ParseMode: telego.ModeHTML,
	}
	if buttonOld != nil {
		paramsOld.ReplyMarkup = buttonOld
	}
	_, _ = c.container.TelegoBot.SendMessage(context.Background(), paramsOld)

	textNew, buttonNew := parser.GetMessageTelego("success-new-paccess-message", data)
	paramsNew := &telego.SendMessageParams{
		ChatID:    telego.ChatID{ID: body.NewOwnerID},
		Text:      textNew,
		ParseMode: telego.ModeHTML,
	}
	if buttonNew != nil {
		paramsNew.ReplyMarkup = buttonNew
	}
	_, _ = c.container.TelegoBot.SendMessage(context.Background(), paramsNew)

	_, err = c.container.CacheService.DeleteAllUserSessionsBySuffix(ctx, body.OldOwnerID)
	if err != nil {
		logger.Error("API", "Erro ao excluir all sessions: %v", err)
	}

	ctx.JSON(http.StatusOK, types.NewSuccessResponse[any](nil, "Dono migrado com sucesso!"))
}
