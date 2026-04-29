package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/leirbagxis/FreddyBot/internal/api/auth"
	"github.com/leirbagxis/FreddyBot/internal/container"
)

func GetChannelHandler(app *container.AppContainer) gin.HandlerFunc {
	return func(c *gin.Context) {
		channelIdStr := c.Param("channelId")
		channelId, err := strconv.ParseInt(channelIdStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "channelId inválido"})
			return
		}

		// Obter dados do contexto (injetados pelo middleware)
		ctxUserID, _ := c.Get("userID")
		ctxRole, _ := c.Get("role")

		userID := ctxUserID.(int64)
		role := ctxRole.(auth.Role)

		channel, err := app.ChannelRepo.GetChannelByID(c, channelId)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Canal não encontrado"})
			return
		}

		// --- VERIFICAÇÃO DE PERMISSÃO ---
		// Se não for Admin/Owner, o UserID do token deve ser igual ao OwnerID do canal
		if role != auth.RoleAdmin && role != auth.RoleOwner {
			if channel.OwnerID != userID {
				c.JSON(http.StatusForbidden, gin.H{"error": "Você não tem permissão para acessar este canal"})
				return
			}
		}

		user, err := app.UserRepo.GetUserById(c, channel.OwnerID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Dono do canal não encontrado"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"user":    user,
			"channel": channel,
		})
	}
}
