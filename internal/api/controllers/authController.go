package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/leirbagxis/FreddyBot/internal/api/auth"
	"github.com/leirbagxis/FreddyBot/internal/container"
	"github.com/leirbagxis/FreddyBot/pkg/config"
)

type AuthController struct {
	container *container.AppContainer
}

func NewAuthController(c *container.AppContainer) *AuthController {
	return &AuthController{container: c}
}

type LoginRequest struct {
	UserID int64 `json:"userID"`
}

func (ac *AuthController) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Dados inválidos"})
		return
	}

	authData := ctx.GetHeader("x-telegram-init-data")
	if authData == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "InitData ausente"})
		return
	}

	// 1. Validação RIGOROSA do Telegram
	validateResult := auth.ValidateTelegramInitData(authData, 3600)
	if !validateResult.IsValid {
		ctx.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Autenticação do Telegram falhou"})
		return
	}

	// 1.1 Verificar se o userID do InitData condiz com o userID do Request
	// Isso impede que um usuário envie o InitData dele mas peça um token para outro ID (impersonation)
	userDataRaw, ok := validateResult.Data["user"]
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Dados de usuário ausentes no Telegram"})
		return
	}

	// O Telegram envia o campo 'user' como um JSON stringificado
	// Exemplo: {"id":123456,"first_name":"...","username":"..."}
	if !strings.Contains(userDataRaw, fmt.Sprintf("\"id\":%d", req.UserID)) {
		ctx.JSON(http.StatusForbidden, gin.H{"success": false, "message": "ID de usuário não coincide com a autenticação"})
		return
	}

	// 2. Determinar Role
	role := auth.RoleUser
	if req.UserID == config.OwnerID {
		role = auth.RoleOwner
	} else {
		user, err := ac.container.UserRepo.GetUserById(ctx, req.UserID)
		if err == nil && user != nil && user.IsAdmin {
			role = auth.RoleAdmin
		}
	}

	// 3. Gerar Token Seguro (Expirando em 12 horas para UX melhor no Mini App, ou conforme necessário)
	token, err := auth.GenerateToken(req.UserID, role, 1, 12*time.Hour)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Erro ao gerar acesso"})
		return
	}

	// 4. Setar Cookie Seguro
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		MaxAge:   43200, // 12h
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"role":    role,
		"token":   token, // Retornamos o token tbm caso o cliente queira salvar fora do cookie
	})
}
