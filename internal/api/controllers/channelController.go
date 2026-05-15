package controllers

import (
	"io"
	"mime"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/leirbagxis/FreddyBot/internal/api/auth"
	"github.com/leirbagxis/FreddyBot/internal/api/dto"
	"github.com/leirbagxis/FreddyBot/internal/api/types"
	"github.com/leirbagxis/FreddyBot/internal/container"
	"github.com/leirbagxis/FreddyBot/pkg/errors"
)

type ChannelController struct {
	container *container.AppContainer
}

func NewChannelController(container *container.AppContainer) *ChannelController {
	return &ChannelController{
		container: container,
	}
}

func (c *ChannelController) GetAllChannelsController(ctx *gin.Context) {
	channels, err := c.container.ChannelService.GetAllChannels(ctx)
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

func (c *ChannelController) GetChannelByIDController(ctx *gin.Context) {
	channelIdStr := ctx.Param("channelId")
	channelId, err := strconv.ParseInt(channelIdStr, 10, 64)
	if err != nil {
		ctx.Error(errors.BadRequest("channelId inválido"))
		return
	}

	// Obter dados do contexto (injetados pelo middleware)
	ctxUserID, _ := ctx.Get("userID")
	ctxRole, _ := ctx.Get("role")

	userID := ctxUserID.(int64)
	role := ctxRole.(auth.Role)

	channel, err := c.container.ChannelService.GetChannelByID(ctx, channelId)
	if err != nil {
		ctx.Error(err)
		return
	}

	// --- VERIFICAÇÃO DE PERMISSÃO ---
	// Se não for Admin/Owner, o UserID do token deve ser igual ao OwnerID do canal
	if role != auth.RoleAdmin && role != auth.RoleOwner {
		if channel.OwnerID != userID {
			ctx.Error(errors.ErrForbidden)
			return
		}
	}

	ctx.JSON(http.StatusOK, types.NewSuccessResponse(gin.H{
		"user":    dto.ToUserDTO(channel.Owner),
		"channel": dto.ToChannelDTO(channel),
	}))
}

func (c *ChannelController) DisconectChannel(ctx *gin.Context) {
	channelIDStr, exists := ctx.Get("channelID")
	channelID, ok := channelIDStr.(int64)
	if !exists || !ok {
		ctx.Error(errors.ErrUnauthorized)
		return
	}

	channel, err := c.container.ChannelService.GetChannelByID(ctx, channelID)
	if err != nil {
		ctx.Error(err)
		return
	}

	err = c.container.ChannelService.DisconnectChannel(ctx, channel.OwnerID, channelID)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *ChannelController) GetSeparator(ctx *gin.Context) {
	channelIdStr := ctx.Param("channelId")
	channelId, err := strconv.ParseInt(channelIdStr, 10, 64)
	if err != nil {
		ctx.Error(errors.BadRequest("ID do canal inválido"))
		return
	}

	separatorId := ctx.Param("separatorId")
	if separatorId == "" {
		ctx.Error(errors.BadRequest("Separator ID é obrigatório"))
		return
	}

	stickerData, err := c.container.SeparatorService.GetSeparatorByTwoID(ctx, channelId, separatorId)
	if err != nil {
		ctx.Error(err)
		return
	}
	if stickerData == nil {
		ctx.Error(errors.ErrNotFound)
		return
	}

	ext := strings.ToLower(filepath.Ext(stickerData.SeparatorURL))

	if ext == ".tgs" {
		ctx.Error(errors.New(http.StatusNotImplemented, "Formato TGS ainda não suportado"))
		return
	}

	resp, err := http.Get(stickerData.SeparatorURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		ctx.Error(errors.New(http.StatusInternalServerError, "Erro ao buscar conteúdo do sticker"))
		return
	}
	defer resp.Body.Close()

	contentType := mime.TypeByExtension(ext)
	if contentType == "" {
		if ext == ".webm" {
			contentType = "video/webm"
		} else {
			contentType = "application/octet-stream"
		}
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		ctx.Error(errors.Internal(err))
		return
	}

	// Adiciona headers explícitos
	ctx.Header("Content-Type", contentType)
	ctx.Header("Content-Disposition", "inline; filename=sticker"+ext)
	ctx.Data(http.StatusOK, contentType, content)
}
