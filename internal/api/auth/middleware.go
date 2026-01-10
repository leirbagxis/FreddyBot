package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/leirbagxis/FreddyBot/internal/api/types"
	"github.com/leirbagxis/FreddyBot/pkg/config"
)

func AuthMiddlewareJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader, err := c.Cookie("token")
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := ValidateTokenJWT(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token invalido ou expirado"})
			return
		}

		// Setando dados no contexto
		c.Set("channelID", claims.ChannelID)
		c.Set("ownerID", claims.OwnerID)

		c.Next()
	}
}

func ValidateTelegramInitData(initData string, secondsToExpire int64) types.ValidateResult {
	params, _ := url.ParseQuery(initData)
	data := make(map[string]string)
	var hash string

	for key, val := range params {
		if key == "hash" {
			hash = val[0]
		} else {
			data[key] = val[0]
		}
	}

	for k, v := range data {
		log.Printf("   %s = %s\n", k, v)
	}

	// 1) ordenar keys
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 2) montar data_check_string
	var parts []string
	for _, k := range keys {
		parts = append(parts, k+"="+data[k])
	}
	dataCheckString := strings.Join(parts, "\n")

	// 3) gerar secret key (Key = HMAC_SHA256("WebAppData", botToken))
	step1 := hmac.New(sha256.New, []byte("WebAppData"))
	step1.Write([]byte(config.TelegramBotToken))
	secretKey := step1.Sum(nil)

	// 4) gerar HMAC com o secretKey
	h := hmac.New(sha256.New, secretKey)
	h.Write([]byte(dataCheckString))
	expectedHash := hex.EncodeToString(h.Sum(nil))

	// 5) validar tempo (auth_date)
	authDateStr := data["auth_date"]
	authUnix, err := strconv.ParseInt(authDateStr, 10, 64)
	if err != nil {
		log.Printf("❌ Erro ao converter auth_date: %s\n", err)
		return types.ValidateResult{IsValid: false, Data: data}
	}

	now := time.Now().Unix()
	diff := now - authUnix

	// Se secondsToExpire == 0 → ignorar expiração
	if secondsToExpire != 0 && diff > secondsToExpire {
		log.Println("❌ Expiração inválida: auth_date muito antigo")
		return types.ValidateResult{IsValid: false, Data: data}
	}

	// Comparar hash
	if expectedHash != hash {
		log.Println("❌ Hash inválido! initData foi alterado.")
		return types.ValidateResult{IsValid: false, Data: data}
	}

	log.Println("✅ initData válido!")

	return types.ValidateResult{IsValid: true, Data: data}
}
