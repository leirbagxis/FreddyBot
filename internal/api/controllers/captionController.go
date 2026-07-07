package controllers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/leirbagxis/FreddyBot/internal/api/types"
	"github.com/leirbagxis/FreddyBot/internal/container"
	"github.com/leirbagxis/FreddyBot/pkg/errors"
)

type CaptionController struct {
	container *container.AppContainer
}

func NewCaptionController(container *container.AppContainer) *CaptionController {
	return &CaptionController{
		container: container,
	}
}

func (c *CaptionController) UpdateDefaultCaptionController(ctx *gin.Context) {
	channelIdStr := ctx.Param("channelId")
	channelId, err := strconv.ParseInt(channelIdStr, 10, 64)
	if err != nil {
		ctx.Error(errors.BadRequest("ID do canal inválido"))
		return
	}

	var captionData types.CaptionDefaultUpdateRequest
	if err := ctx.ShouldBindJSON(&captionData); err != nil {
		ctx.Error(errors.BadRequest("Dados inválidos: " + err.Error()))
		return
	}

	rowsAffected, err := c.container.CaptionService.UpdateDefaultCaption(ctx, channelId, captionData)
	if err != nil {
		ctx.Error(err)
		return
	}

	// ── Baixar emojis personalizados em background ──
	userID, _ := ctx.Get("userID")
	if uid, ok := userID.(int64); ok && c.container.EmojiService.TextContainsEmoji(captionData.Caption) {
		go func(uID int64, caption string) {
			if err := c.container.EmojiService.EnsureEmojisForText(context.Background(), uID, caption); err != nil {
				// Apenas loga, não quebra o fluxo
			}
		}(uid, captionData.Caption)
	}

	ctx.JSON(http.StatusOK, types.NewSuccessResponse(gin.H{"rows_affected": rowsAffected}, "Legenda padrão atualizada com sucesso"))
}

func (c *CaptionController) UpdateNewPackCaptionController(ctx *gin.Context) {
	channelIdStr := ctx.Param("channelId")
	channelId, err := strconv.ParseInt(channelIdStr, 10, 64)
	if err != nil {
		ctx.Error(errors.BadRequest("ID do canal inválido"))
		return
	}

	var captionData types.NewPackCaptionUpdateRequest
	if err := ctx.ShouldBindJSON(&captionData); err != nil {
		ctx.Error(errors.BadRequest("Dados inválidos: " + err.Error()))
		return
	}

	rowsAffected, err := c.container.CaptionService.UpdateNewPackCaption(ctx, channelId, captionData)
	if err != nil {
		ctx.Error(err)
		return
	}

	// ── Baixar emojis personalizados em background ──
	captionText := captionData.Text()
	userID, _ := ctx.Get("userID")
	if uid, ok := userID.(int64); ok && c.container.EmojiService.TextContainsEmoji(captionText) {
		go func(uID int64, caption string) {
			_ = c.container.EmojiService.EnsureEmojisForText(context.Background(), uID, caption)
		}(uid, captionText)
	}

	ctx.JSON(http.StatusOK, types.NewSuccessResponse(gin.H{"rows_affected": rowsAffected}, "Legenda de novos packs atualizada com sucesso"))
}

func (c *CaptionController) UpdateReactionsController(ctx *gin.Context) {
	channelIdStr := ctx.Param("channelId")
	channelId, err := strconv.ParseInt(channelIdStr, 10, 64)
	if err != nil {
		ctx.Error(errors.BadRequest("ID do canal inválido"))
		return
	}

	var reactionsData types.ReactionsUpdateRequest
	if err := ctx.ShouldBindJSON(&reactionsData); err != nil {
		ctx.Error(errors.BadRequest("Dados inválidos: " + err.Error()))
		return
	}

	rowsAffected, err := c.container.CaptionService.UpdateReactions(ctx, channelId, reactionsData)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, types.NewSuccessResponse(gin.H{"rows_affected": rowsAffected}, "Reações atualizadas com sucesso"))
}

func (c *CaptionController) UpdateReactionPositionController(ctx *gin.Context) {
	channelIdStr := ctx.Param("channelId")
	channelId, err := strconv.ParseInt(channelIdStr, 10, 64)
	if err != nil {
		ctx.Error(errors.BadRequest("ID do canal inválido"))
		return
	}

	var posData types.ReactionPositionUpdateRequest
	if err := ctx.ShouldBindJSON(&posData); err != nil {
		ctx.Error(errors.BadRequest("Dados inválidos: " + err.Error()))
		return
	}

	rowsAffected, err := c.container.CaptionService.UpdateReactionPosition(ctx, channelId, posData)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, types.NewSuccessResponse(gin.H{"rows_affected": rowsAffected}, "Posição das reações atualizada com sucesso"))
}
