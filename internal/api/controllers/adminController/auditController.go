package admincontroller

import (
	"context"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/leirbagxis/FreddyBot/internal/api/types"
	"github.com/leirbagxis/FreddyBot/internal/container"
	userModes "github.com/leirbagxis/FreddyBot/internal/database/models"
	"github.com/mymmrac/telego"
)

type AuditController struct {
	container *container.AppContainer
}

func NewAuditController(app *container.AppContainer) *AuditController {
	return &AuditController{
		container: app,
	}
}

type AuditResult struct {
	UserID    int64               `json:"userId"`
	FirstName string              `json:"firstName"`
	Channels  []userModes.Channel `json:"channels"`
}

type BulkDeleteRequest struct {
	UserID     int64   `json:"userId"`
	ChannelIDs []int64 `json:"channelIds"`
}

func (c *AuditController) GetCheckBotAudit(ctx *gin.Context) {
	const targetBotID = 5986082367
	bot := c.container.TelegoBot

	channels, err := c.container.ChannelService.GetAllChannels(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	var foundChannels []userModes.Channel
	var mu sync.Mutex

	chQueue := make(chan *userModes.Channel, len(channels))
	for i := range channels {
		chQueue <- &channels[i]
	}
	close(chQueue)

	var wg sync.WaitGroup
	numWorkers := 20 // Mais workers para a API ser responsiva
	if len(channels) < numWorkers {
		numWorkers = len(channels)
	}

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for ch := range chQueue {
				member, err := bot.GetChatMember(context.Background(), &telego.GetChatMemberParams{
					ChatID: telego.ChatID{ID: ch.ID},
					UserID: targetBotID,
				})

				if err == nil {
					status := member.MemberStatus()
					if status == telego.MemberStatusAdministrator || status == telego.MemberStatusCreator {
						mu.Lock()
						foundChannels = append(foundChannels, *ch)
						mu.Unlock()
					}
				}
			}
		}()
	}
	wg.Wait()

	// Agrupar por Usuário
	userGroups := make(map[int64]*AuditResult)
	var userOrder []int64

	for _, ch := range foundChannels {
		if _, ok := userGroups[ch.OwnerID]; !ok {
			owner, _ := c.container.UserService.GetUserByID(context.Background(), ch.OwnerID)
			name := "Desconhecido"
			if owner != nil {
				name = owner.FirstName
			}
			userGroups[ch.OwnerID] = &AuditResult{
				UserID:    ch.OwnerID,
				FirstName: name,
				Channels:  []userModes.Channel{},
			}
			userOrder = append(userOrder, ch.OwnerID)
		}
		userGroups[ch.OwnerID].Channels = append(userGroups[ch.OwnerID].Channels, ch)
	}

	results := make([]*AuditResult, 0, len(userOrder))
	for _, id := range userOrder {
		results = append(results, userGroups[id])
	}

	ctx.JSON(http.StatusOK, types.NewSuccessResponse(results))
}

func (c *AuditController) BulkDeleteUserChannels(ctx *gin.Context) {
	var req BulkDeleteRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, types.NewErrorResponse("Dados inválidos: "+err.Error()))
		return
	}

	if len(req.ChannelIDs) == 0 {
		ctx.JSON(http.StatusBadRequest, types.NewErrorResponse("Nenhum canal selecionado para exclusão"))
		return
	}

	var failedIDs []int64
	deletedCount := 0

	for _, channelID := range req.ChannelIDs {
		err := c.container.ChannelService.DisconnectChannel(ctx, req.UserID, channelID)
		if err != nil {
			failedIDs = append(failedIDs, channelID)
		} else {
			deletedCount++
		}
	}

	if len(failedIDs) > 0 {
		ctx.JSON(http.StatusMultiStatus, types.NewSuccessResponse(gin.H{
			"deletedCount": deletedCount,
			"failedIds":    failedIDs,
		}, "Alguns canais não puderam ser excluídos"))
		return
	}

	ctx.JSON(http.StatusOK, types.NewSuccessResponse(gin.H{
		"deletedCount": deletedCount,
	}, "Todos os canais foram excluídos com sucesso"))
}
