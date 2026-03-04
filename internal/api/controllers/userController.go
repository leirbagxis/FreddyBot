package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-telegram/bot"
	"github.com/leirbagxis/FreddyBot/internal/api/auth"
	"github.com/leirbagxis/FreddyBot/internal/api/types"
	"github.com/leirbagxis/FreddyBot/internal/container"
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

func (c *UserController) GetUserInfo(ctx *gin.Context) {
	userParams := ctx.Param("userParams")
	if len(userParams) < 5 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"succes":  false,
			"message": "ID ou Username inválido!",
		})
		return
	}

	userID, _ := strconv.ParseInt(userParams, 10, 64)
	if userID == 0 {
		user, err := c.container.UserRepo.GetUserByUsername(context.Background(), userParams)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"succes":  false,
				"message": "Usuario nao encontrado",
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"succes": true,
			"user":   user,
		})
		return

	}

	user, err := c.container.Bot.GetChat(context.Background(), &bot.GetChatParams{
		ChatID: userID,
	})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"succes":  false,
			"message": "Usuario nao encontrado",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"succes": true,
		"user":   user,
	})
}

func (c *UserController) TransferChannelController(ctx *gin.Context) {
	var body *types.TransferChannelRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "payload inválido", "details": err.Error()})
		return
	}

	channel, err := c.container.ChannelRepo.GetChannelByTwoID(ctx, body.OldOwnerID, body.ChannelID)
	if err != nil {
		log.Printf("Erro ao buscar canal: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Canal não encontrado ou você não tem permissão para alterá-lo.",
		})
		return
	}
	if channel == nil {
		log.Printf("Canal retornado é nil para channelId=%d e userId=%d", body.ChannelID, body.OldOwnerID)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Erro interno: canal não encontrado.",
		})
		return
	}

	if body.NewOwnerID == body.OldOwnerID {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "O novo dono precisa ser diferente de voce.",
		})
		return
	}

	newOwner, err := c.container.Bot.GetChat(ctx, &bot.GetChatParams{ChatID: body.NewOwnerID})
	if err != nil {
		log.Println("Erro ao obter chat do novo dono:", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "O novo dono precisa iniciar o bot pelo menos uma vez. Peça para ele mandar uma mensagem no bot antes de transferir o canal.",
		})
		return
	}

	// Verifica se o novo dono é um bot
	if body.NewOwnerID == c.container.Bot.ID() {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "O novo dono não pode ser eu.",
		})
		return
	}

	admins, err := c.container.Bot.GetChatAdministrators(ctx, &bot.GetChatAdministratorsParams{
		ChatID: body.ChannelID,
	})
	if err != nil {
		log.Println("Erro ao buscar administradores do canal:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Erro ao consultar administradores do canal.",
		})
		return
	}

	isAdmin := false
	for _, admin := range admins {
		if admin.Administrator != nil && admin.Administrator.User.ID == body.NewOwnerID {
			isAdmin = true
			break
		}
		if admin.Owner != nil && admin.Owner.User != nil && admin.Owner.User.ID == body.NewOwnerID {
			isAdmin = true
			break
		}
	}

	if !isAdmin {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "O novo dono precisa ser administrador do canal.",
		})
		return
	}

	// Deletar dados vinculados ao antigo dono
	_ = c.container.SeparatorRepo.DeleteSeparatorByOwnerChannelId(ctx, body.OldOwnerID)

	err = c.container.ChannelRepo.UpdateOwnerChannel(ctx, body.ChannelID, body.OldOwnerID, body.NewOwnerID)
	if err != nil {
		log.Printf("Erro ao transferir posse do canal: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Erro ao passar a posse para o novo usuário.",
		})
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

	textOld, buttonOld := parser.GetMessage("success-old-paccess-message", data)
	c.container.Bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      body.OldOwnerID,
		Text:        textOld,
		ReplyMarkup: buttonOld,
		ParseMode:   "HTML",
	})

	textNew, buttonNew := parser.GetMessage("success-new-paccess-message", data)
	c.container.Bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      body.NewOwnerID,
		Text:        textNew,
		ReplyMarkup: buttonNew,
		ParseMode:   "HTML",
	})

	_, err = c.container.CacheService.DeleteAllUserSessionsBySuffix(ctx, body.OldOwnerID)
	if err != nil {
		log.Printf("Erro ao excluir all sessions: %v", err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"succes":  true,
		"message": "Dono migrado com sucesso!",
	})
}
