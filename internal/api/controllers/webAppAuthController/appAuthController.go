package webappauthcontroller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/leirbagxis/FreddyBot/internal/api/auth"
	"github.com/leirbagxis/FreddyBot/internal/api/types"
	"github.com/leirbagxis/FreddyBot/internal/container"
)

type WebAppAuthController struct {
	container *container.AppContainer
}

func NewWebAppAuthController(container *container.AppContainer) *WebAppAuthController {
	return &WebAppAuthController{
		container: container,
	}
}

func (c *WebAppAuthController) ReceiveAuthController(ctx *gin.Context) {
	var authData types.WebAPPAuthRequest
	if err := ctx.ShouldBindJSON(&authData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Dados inválidos: " + err.Error(),
		})
		return
	}

	authHeader := ctx.GetHeader("x-telegram-init-data")

	result := auth.ValidateTelegramInitData(authHeader, 86400)
	if !result.IsValid {
		fmt.Println("❌ initData inválido!")
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "❌ initData inválido!",
		})
		return
	}

	user, err := c.container.ChannelRepo.GetChannelByUserID(ctx, authData.User.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Erro ao encontrar usuario: " + err.Error(),
		})
		return
	}

	channel, err := c.container.ChannelRepo.GetChannelByTwoID(ctx, user.OwnerID, authData.ChannelID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Erro ao encontrar Canal: " + err.Error(),
		})
		return
	}

	token, err := auth.GenerateTokenJWT(strconv.FormatInt(channel.ID, 10), strconv.FormatInt(user.ID, 10))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Erro ao gerar token",
			"details": err,
		})
		return
	}

	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode, // 👈 Proteção contra CSRF
	})

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"channel": channel,
	})
}
