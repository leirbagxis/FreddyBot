package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/leirbagxis/FreddyBot/internal/api/auth"
	"github.com/leirbagxis/FreddyBot/internal/container"
	"github.com/leirbagxis/FreddyBot/pkg/config"
)

func VerifyJWTHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"valid": false,
				"error": "Token não fornecido",
			})
			return
		}

		tokenString := authHeader
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenString = authHeader[7:]
		}

		claims, err := auth.ValidateTokenJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"valid": false,
				"error": "Token inválido",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"valid":      true,
			"channel_id": claims.ChannelID,
			"owner_id":   claims.OwnerID,
			"expires_at": claims.ExpiresAt.Time,
		})
	}
}

func GenerateJWTHandler(app *container.AppContainer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request struct {
			ChannelID string `json:"channelId" binding:"required"`
			Signature string `json:"signature" binding:"required"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "channel_id e signature são obrigatórios",
			})
			return
		}

		channelIDInt, err := strconv.ParseInt(request.ChannelID, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "channel_id deve ser um número válido",
			})
			return
		}
		channel, err := app.ChannelRepo.GetChannelByID(c, channelIDInt)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Canal não encontrado",
			})
			return
		}
		ownerID := channel.OwnerID
		ownerIDStr := strconv.FormatInt(ownerID, 10)

		if !auth.ValidateSignature(ownerIDStr, request.ChannelID, request.Signature, config.SecreteKey) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Signature inválida",
			})
			return
		}

		token, err := auth.GenerateTokenJWT(request.ChannelID, ownerIDStr)
		fmt.Println(token, err)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Erro ao gerar token",
				"details": err,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":   "Token gerado com sucesso!",
			"token":     token,
			"channelId": request.ChannelID,
			"ownerID":   ownerID,
		})

	}
}
