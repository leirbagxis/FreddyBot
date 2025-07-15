package controllers

import (
	"io"
	"mime"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/leirbagxis/FreddyBot/internal/container"
)

type SeparatorController struct {
	container *container.AppContainer
}

func NewSeparatorController(container *container.AppContainer) *SeparatorController {
	return &SeparatorController{
		container: container,
	}
}

func (ctrl *SeparatorController) GetSeparator(ctx *gin.Context) {
	channelIdStr := ctx.Param("channelId")
	channelId, err := strconv.ParseInt(channelIdStr, 10, 54)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID do canal inválido",
		})
		return
	}

	separatorId := ctx.Param("separatorId")
	if separatorId == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Separator ID é obrigatório"})
		return
	}

	stickerData, err := ctrl.container.SeparatorRepo.GetSeparatorByTwoID(ctx, channelId, separatorId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar separator"})
		return
	}
	if stickerData == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Separator não encontrado"})
		return
	}

	ext := strings.ToLower(filepath.Ext(stickerData.SeparatorURL))

	if ext == ".tgs" {
		ctx.JSON(http.StatusNotImplemented, gin.H{"error": "Formato TGS ainda não suportado"})
		return
	}

	resp, err := http.Get(stickerData.SeparatorURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar conteúdo do sticker"})
		return
	}
	defer resp.Body.Close()

	contentType := mime.TypeByExtension(ext)
	if contentType == "" {
		if ext == ".webm" {
			contentType = "video/webm"
		} else {
			contentType = "application/octet-stream"
		}
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao ler conteúdo"})
		return
	}

	// Adiciona headers explícitos
	ctx.Header("Content-Type", contentType)
	ctx.Header("Content-Disposition", "inline; filename=sticker"+ext)
	ctx.Data(http.StatusOK, contentType, content)
}

// func (ctrl *SeparatorController) GetSeparator(ctx *gin.Context) {
// 	channelIdStr := ctx.Param("channelId")

// 	channelId, err := strconv.ParseInt(channelIdStr, 10, 54)
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{
// 			"success": false,
// 			"message": "ID do canal inválido",
// 		})
// 		return
// 	}

// 	separatorId := ctx.Param("separatorId")
// 	if separatorId == "" {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Separator ID é obrigatório"})
// 		return
// 	}

// 	stickerData, err := ctrl.container.SeparatorRepo.GetSeparatorByTwoID(ctx, channelId, separatorId)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar separator"})
// 		return
// 	}
// 	if stickerData == nil {
// 		ctx.JSON(http.StatusNotFound, gin.H{"error": "Separator não encontrado"})
// 		return
// 	}

// 	ext := strings.ToLower(filepath.Ext(stickerData.SeparatorURL))

// 	var content []byte
// 	if ext == ".tgs" {
// 		return
// 	}

// 	// Content-Type padrão
// 	contentType := mime.TypeByExtension(ext)
// 	if contentType == "" {
// 		contentType = "application/octet-stream"
// 	}

// 	// Faz o download do conteúdo da URL
// 	resp, err := http.Get(stickerData.SeparatorURL)
// 	if err != nil || resp.StatusCode != http.StatusOK {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar imagem do sticker"})
// 		return
// 	}
// 	defer resp.Body.Close()

// 	content, err = io.ReadAll(resp.Body)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao ler imagem"})
// 		return
// 	}

// 	ctx.Data(http.StatusOK, contentType, content)
// }
