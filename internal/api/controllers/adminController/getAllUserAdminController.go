package admincontroller

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/leirbagxis/FreddyBot/internal/container"
	"github.com/leirbagxis/FreddyBot/internal/utils"
	"github.com/leirbagxis/FreddyBot/pkg/logger"
)

type UsersAdminController struct {
	container *container.AppContainer
}

func NewUsersAdminController(app *container.AppContainer) *UsersAdminController {
	return &UsersAdminController{
		container: app,
	}
}

type NoticeRequest struct {
	Message  string `json:"message"`
	Target   string `json:"target"`
	ImageUrl string `json:"imageUrl"`
	Buttons  []struct {
		Text  string `json:"text"`
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"buttons"`
}

func (c *UsersAdminController) GetAllUsersAdminController(ctx *gin.Context) {
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

func (c *UsersAdminController) UpdateUserAdminController(ctx *gin.Context) {
	userIDStr := ctx.Param("userId")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID de usuário inválido"})
		return
	}

	newValue, err := c.container.UserRepo.UpdateUserAdmin(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"isAdmin": newValue,
	})
}

func (c *UsersAdminController) UpdateUserBlacklistController(ctx *gin.Context) {
	userIDStr := ctx.Param("userId")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID de usuário inválido"})
		return
	}

	newValue, err := c.container.UserRepo.UpdateUserBlacklist(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":       true,
		"isBlacklisted": newValue,
	})
}

func (c *UsersAdminController) GetAdminOverview(ctx *gin.Context) {
	users, err := c.container.AdminService.GetAllUsersAdminRepository(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	channels, err := c.container.ChannelRepo.GetAllChannels(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"users":    users,
		"channels": channels,
	})
}

func (c *UsersAdminController) SendNoticeAdminController(ctx *gin.Context) {
	var notice NoticeRequest

	if err := ctx.ShouldBindJSON(&notice); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Broadcast iniciado",
	})

	go c.dispatchNotice(notice)
}

func (c *UsersAdminController) dispatchNotice(notice NoticeRequest) {
	ctx := context.Background()

	var buttons []container.BroadcastButton
	for _, btn := range notice.Buttons {
		buttons = append(buttons, container.BroadcastButton{
			Text:  btn.Text,
			Type:  btn.Type,
			Value: btn.Value,
		})
	}

	switch notice.Target {

	case "users":
		users, err := c.container.AdminService.GetAllUsersAdminRepository(ctx)
		if err != nil {
			logger.Error("API", "Erro ao buscar usuários para broadcast: %v", err)
			return
		}

		for _, user := range users {
			c.container.BroadcastQueue <- container.BroadcastJob{
				ChatID:   user.UserId,
				Text:     utils.MarkdownToTelegramHTML(notice.Message),
				ImageUrl: notice.ImageUrl,
				Buttons:  buttons,
			}
		}

	case "channels":
		channels, err := c.container.ChannelRepo.GetAllChannels(ctx)
		if err != nil {
			logger.Error("API", "Erro ao buscar usuários para broadcast: %v", err)
			return
		}

		for _, channel := range channels {
			c.container.BroadcastQueue <- container.BroadcastJob{
				ChatID:   channel.ID,
				Text:     utils.MarkdownToTelegramHTML(notice.Message),
				ImageUrl: notice.ImageUrl,
				Buttons:  buttons,
			}
		}

	case "all":
		users, _ := c.container.AdminService.GetAllUsersAdminRepository(ctx)
		channels, _ := c.container.ChannelRepo.GetAllChannels(ctx)

		sentMap := make(map[int64]bool)

		for _, user := range users {
			if !sentMap[user.UserId] {
				c.container.BroadcastQueue <- container.BroadcastJob{
					ChatID:   user.UserId,
					Text:     utils.MarkdownToTelegramHTML(notice.Message),
					ImageUrl: notice.ImageUrl,
					Buttons:  buttons,
				}
				sentMap[user.UserId] = true
			}
		}

		for _, channel := range channels {
			if !sentMap[channel.ID] {
				c.container.BroadcastQueue <- container.BroadcastJob{
					ChatID:   channel.ID,
					Text:     utils.MarkdownToTelegramHTML(notice.Message),
					ImageUrl: notice.ImageUrl,
					Buttons:  buttons,
				}
				sentMap[channel.ID] = true
			}
		}
	}
}
