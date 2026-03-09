package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/leirbagxis/FreddyBot/internal/api/types"
	"github.com/leirbagxis/FreddyBot/internal/container"
	"github.com/leirbagxis/FreddyBot/pkg/config"
)

func AuthMiddlewareJWT(v *container.AppContainer) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader, err := c.Cookie("token")
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := ValidateChannelToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token invalido ou expirado"})
			return
		}

		if claims.IsAdmin {
			c.Set("channelID", claims.ChannelID)
			c.Set("isAdmin", true)
			c.Next()
			return
		}

		channel, err := v.ChannelRepo.GetChannelByID(c.Request.Context(), claims.ChannelID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "canal não encontrado"})
			return
		}

		if channel.TokenVersion != claims.TV {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token expirado"})
			return
		}

		// Setando dados no contexto
		c.Set("channelID", claims.ChannelID)
		c.Set("ownerID", claims.Subject)
		c.Set("isAdmin", false)
		c.Next()
	}
}

func ValidateTelegramInitData(initData string, secondsToExpire int64) types.ValidateResult {
	params, err := url.ParseQuery(initData)
	if err != nil {
		return types.ValidateResult{IsValid: false}
	}

	data := make(map[string]string)
	var hash string

	for key, val := range params {
		if key == "hash" {
			hash = val[0]
			continue
		}
		data[key] = val[0]
	}

	// ordenar keys
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var parts []string
	for _, k := range keys {
		parts = append(parts, k+"="+data[k])
	}
	dataCheckString := strings.Join(parts, "\n")

	// gerar secret
	step1 := hmac.New(sha256.New, []byte("WebAppData"))
	step1.Write([]byte(config.TelegramBotToken))
	secretKey := step1.Sum(nil)

	h := hmac.New(sha256.New, secretKey)
	h.Write([]byte(dataCheckString))
	expectedHash := hex.EncodeToString(h.Sum(nil))

	if expectedHash != hash {
		return types.ValidateResult{IsValid: false, Data: data}
	}

	authUnix, err := strconv.ParseInt(data["auth_date"], 10, 64)
	if err != nil {
		return types.ValidateResult{IsValid: false, Data: data}
	}

	if secondsToExpire != 0 {
		if time.Now().Unix()-authUnix > secondsToExpire {
			return types.ValidateResult{IsValid: false, Data: data}
		}
	}

	return types.ValidateResult{IsValid: true, Data: data}
}
