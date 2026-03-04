package webappauthcontroller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/leirbagxis/FreddyBot/internal/api/auth"
	"github.com/leirbagxis/FreddyBot/internal/api/types"
	"github.com/leirbagxis/FreddyBot/internal/container"
	"github.com/leirbagxis/FreddyBot/internal/database/models"
	"github.com/leirbagxis/FreddyBot/pkg/config"
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

	result := auth.ValidateTelegramInitData(authHeader, 3600)
	if !result.IsValid {
		fmt.Println("❌ initData inválido!")
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "❌ initData inválido!",
		})
		return
	}

	isAdmin := false
	var channel *models.Channel
	var err error

	if authData.User.ID == config.OwnerID {
		channel, err = c.container.ChannelRepo.GetChannelByID(ctx, authData.ChannelID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Erro ao encontrar Canal: " + err.Error(),
			})
			return
		}
		isAdmin = true
	} else {
		channel, err = c.container.ChannelRepo.GetChannelByTwoID(ctx, authData.User.ID, authData.ChannelID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Erro ao encontrar Canal: " + err.Error(),
			})
			return
		}

	}

	fmt.Println(isAdmin)

	token, err := auth.GenerateChannelToken(fmt.Sprintf("%d", channel.ID), fmt.Sprintf("%d", channel.OwnerID), isAdmin, channel.TokenVersion, 16*time.Minute)
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

func (c *WebAppAuthController) ReceiveAuthMeChannelsController(ctx *gin.Context) {
	var authData types.MeChannelsAuthRequest
	if err := ctx.ShouldBindJSON(&authData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Dados inválidos: " + err.Error(),
		})
		return
	}

	authHeader := ctx.GetHeader("x-telegram-init-data")

	result := auth.ValidateTelegramInitData(authHeader, 3600)
	if !result.IsValid {
		fmt.Println("❌ initData inválido!")
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "❌ initData inválido!",
		})
		return
	}

	channels, err := c.container.ChannelRepo.GetAllChannelsByUserID(ctx, authData.User.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Erro ao encontrar usuario: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":  true,
		"channels": channels,
	})
}

func (c *WebAppAuthController) AdminAuthController(ctx *gin.Context) {
	var authData types.MeChannelsAuthRequest
	if err := ctx.ShouldBindJSON(&authData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Dados inválidos: " + err.Error(),
		})
		return
	}

	if authData.User.ID != config.OwnerID {
		fmt.Println("❌ O usuario nao e admin!")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Falha ao autenticar",
		})
		return
	}

	authHeader := ctx.GetHeader("x-telegram-init-data")

	result := auth.ValidateTelegramInitData(authHeader, 3600)
	if !result.IsValid {
		fmt.Println("❌ initData inválido!")
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "❌ initData inválido!",
		})
		return
	}

	users, err := c.container.AdminService.GetAllUsersAdminRepository(ctx)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"users":   users,
	})
}
