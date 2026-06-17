package controllers

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/leirbagxis/FreddyBot/internal/api/types"
	"github.com/leirbagxis/FreddyBot/internal/container"
	"github.com/leirbagxis/FreddyBot/pkg/errors"
	"github.com/leirbagxis/FreddyBot/pkg/logger"
	"github.com/mymmrac/telego"
)

type ConnectController struct {
	container *container.AppContainer
}

func NewConnectController(c *container.AppContainer) *ConnectController {
	return &ConnectController{container: c}
}

func (cc *ConnectController) Start(g *gin.Context) {
	userID := g.GetInt64("userID")
	if userID == 0 {
		g.Error(errors.ErrUnauthorized)
		return
	}

	var req types.ConnectStartRequest
	if err := g.ShouldBindJSON(&req); err != nil {
		g.Error(errors.BadRequest("phone é obrigatório"))
		return
	}

	if err := cc.container.TelegramClientService.StartPhoneFlow(g, userID, req.Phone); err != nil {
		if strings.Contains(err.Error(), "auth flow already in progress") {
			g.Error(errors.New(http.StatusConflict, err.Error()))
			return
		}
		g.Error(errors.Internal(err))
		return
	}

	logger.Info("CONNECT", "User %d — code sent to %s", userID, req.Phone)
	g.JSON(http.StatusOK, types.NewSuccessResponse(map[string]string{
		"message": "Código enviado para seu Telegram",
	}, "Código enviado"))
}

func (cc *ConnectController) Verify(g *gin.Context) {
	userID := g.GetInt64("userID")
	if userID == 0 {
		g.Error(errors.ErrUnauthorized)
		return
	}

	var req types.ConnectVerifyRequest
	if err := g.ShouldBindJSON(&req); err != nil {
		g.Error(errors.BadRequest("código é obrigatório"))
		return
	}

	logger.Info("CONNECT", "User %d submitting code", userID)
	if err := cc.container.TelegramClientService.SubmitCode(g, userID, req.Code); err != nil {
		g.Error(errors.Internal(err))
		return
	}

	err := cc.container.TelegramClientService.WaitAuthResult(g, userID)
	if err != nil {
		if err.Error() == "SESSION_PASSWORD_NEEDED" {
			logger.Info("CONNECT", "User %d — 2FA required after code", userID)
			g.JSON(http.StatusOK, types.NewSuccessResponse(map[string]bool{
				"needs2FA": true,
			}, "Senha 2FA necessária"))
			return
		}
		logger.Error("CONNECT", "User %d — verify failed: %v", userID, err)
		g.Error(errors.Internal(err))
		return
	}

	logger.Info("CONNECT", "User %d connected successfully (code verification)", userID)
	cc.sendConnectedMessages(g.Request.Context(), userID)

	g.JSON(http.StatusOK, types.NewSuccessResponse(map[string]bool{
		"connected": true,
	}, "Conta conectada com sucesso"))
}

func (cc *ConnectController) Submit2FA(g *gin.Context) {
	userID := g.GetInt64("userID")
	if userID == 0 {
		g.Error(errors.ErrUnauthorized)
		return
	}

	var req types.Connect2FARequest
	if err := g.ShouldBindJSON(&req); err != nil {
		g.Error(errors.BadRequest("senha é obrigatória"))
		return
	}

	logger.Info("CONNECT", "User %d submitting 2FA password", userID)
	if err := cc.container.TelegramClientService.Submit2FA(g, userID, req.Password); err != nil {
		g.Error(errors.Internal(err))
		return
	}

	err := cc.container.TelegramClientService.WaitAuthResult(g, userID)
	if err != nil {
		logger.Error("CONNECT", "User %d — 2FA failed: %v", userID, err)
		g.Error(errors.Internal(err))
		return
	}

	logger.Info("CONNECT", "User %d connected successfully (with 2FA)", userID)
	cc.sendConnectedMessages(g.Request.Context(), userID)

	g.JSON(http.StatusOK, types.NewSuccessResponse(map[string]bool{
		"connected": true,
	}, "Conta conectada com sucesso"))
}

func (cc *ConnectController) Status(g *gin.Context) {
	userID := g.GetInt64("userID")
	if userID == 0 {
		g.Error(errors.ErrUnauthorized)
		return
	}

	connected, err := cc.container.TelegramClientService.IsConnected(g, userID)
	if err != nil {
		logger.Error("CONNECT", "User %d — status check failed: %v", userID, err)
		g.Error(errors.Internal(err))
		return
	}

	logger.Info("CONNECT", "User %d — status: connected=%v", userID, connected)
	g.JSON(http.StatusOK, types.NewSuccessResponse(types.ConnectStatusResponse{
		Connected: connected,
		UserID:    userID,
	}))
}

func (cc *ConnectController) Disconnect(g *gin.Context) {
	userID := g.GetInt64("userID")
	if userID == 0 {
		g.Error(errors.ErrUnauthorized)
		return
	}

	if err := cc.container.TelegramClientService.DisconnectUser(g, userID); err != nil {
		g.Error(errors.Internal(err))
		return
	}

	logger.Info("CONNECT", "User %d disconnected Telegram account", userID)
	g.JSON(http.StatusOK, types.NewSuccessResponse(map[string]bool{
		"disconnected": true,
	}, "Conta desconectada"))
}

func (cc *ConnectController) sendConnectedMessages(ctx context.Context, userID int64) {
	msg := "✅ Sua conta Telegram foi conectada ao bot com sucesso!\n\n"
	msg += "Agora você pode usar recursos avançados que exigem conexão MTProto."

	if _, err := cc.container.TelegoBot.SendMessage(ctx, &telego.SendMessageParams{
		ChatID:    telego.ChatID{ID: userID},
		Text:      msg,
		ParseMode: telego.ModeHTML,
	}); err != nil {
		logger.Error("CONNECT", "Failed to notify user %d: %v", userID, err)
	}
}
