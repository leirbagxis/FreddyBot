package admincontroller

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leirbagxis/FreddyBot/internal/container"
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
	Message string `json:"message"`
	Target  string `json:"target"`
	Buttons []struct {
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

	ctx.JSON(http.StatusOK, users)

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
			log.Println(err)
			return
		}

		for _, user := range users {
			c.container.BroadcastQueue <- container.BroadcastJob{
				ChatID:  user.UserId,
				Text:    notice.Message,
				Buttons: buttons,
			}
		}

	case "channels":
		channels, err := c.container.ChannelRepo.GetAllChannels(ctx)
		if err != nil {
			log.Println(err)
			return
		}

		for _, channel := range channels {
			c.container.BroadcastQueue <- container.BroadcastJob{
				ChatID:  channel.ID,
				Text:    notice.Message,
				Buttons: buttons,
			}
		}

	case "all":
		users, _ := c.container.AdminService.GetAllUsersAdminRepository(ctx)
		channels, _ := c.container.ChannelRepo.GetAllChannels(ctx)

		for _, user := range users {
			c.container.BroadcastQueue <- container.BroadcastJob{
				ChatID:  user.UserId,
				Text:    notice.Message,
				Buttons: buttons,
			}
		}

		for _, channel := range channels {
			c.container.BroadcastQueue <- container.BroadcastJob{
				ChatID:  channel.ID,
				Text:    notice.Message,
				Buttons: buttons,
			}
		}
	}
}
