package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/leirbagxis/FreddyBot/internal/api/auth"
	"github.com/leirbagxis/FreddyBot/internal/api/types"
	"github.com/leirbagxis/FreddyBot/internal/container"
	"github.com/leirbagxis/FreddyBot/pkg/config"
	"github.com/leirbagxis/FreddyBot/pkg/errors"
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
		ctx.Error(errors.BadRequest("Dados inválidos"))
		return
	}

	authData := ctx.GetHeader("x-telegram-init-data")
	if authData == "" {
		ctx.Error(errors.New(http.StatusUnauthorized, "InitData ausente"))
		return
	}

	// 1. Validação RIGOROSA do Telegram
	validateResult := auth.ValidateTelegramInitData(authData, 3600)
	if !validateResult.IsValid {
		ctx.Error(errors.New(http.StatusUnauthorized, "Autenticação do Telegram falhou"))
		return
	}

	// 1.1 Verificar se o userID do InitData condiz com o userID do Request
	userDataRaw, ok := validateResult.Data["user"]
	if !ok {
		ctx.Error(errors.New(http.StatusUnauthorized, "Dados de usuário ausentes no Telegram"))
		return
	}

	if !strings.Contains(userDataRaw, fmt.Sprintf("\"id\":%d", req.UserID)) {
		ctx.Error(errors.ErrForbidden)
		return
	}

	// 2. Determinar Role
	role := auth.RoleUser
	isBlacklisted := false
	if req.UserID == config.OwnerID {
		role = auth.RoleOwner
	} else {
		user, err := ac.container.UserService.GetUserByID(ctx, req.UserID)
		if err == nil && user != nil {
			if user.IsAdmin {
				role = auth.RoleAdmin
			}
			if user.IsBlacklisted {
				isBlacklisted = true
			}
		}
	}

	// 3. Gerar Token Seguro
	token, err := auth.GenerateToken(req.UserID, role, 1, 12*time.Hour)
	if err != nil {
		ctx.Error(errors.Internal(err))
		return
	}

	// 4. Setar Cookie Seguro
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		MaxAge:   43200, // 12h
		HttpOnly: true,
		Secure:   config.AppEnv != "dev",
		SameSite: http.SameSiteStrictMode,
	})

	ctx.JSON(http.StatusOK, types.NewSuccessResponse(gin.H{
		"role":          role,
		"token":         token,
		"isBlacklisted": isBlacklisted,
	}, "Login realizado com sucesso"))
}
