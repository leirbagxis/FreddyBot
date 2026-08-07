package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/leirbagxis/FreddyBot/internal/cache"
)

// RateLimit limita a taxa de requisições por IP usando o Redis.
func RateLimit(limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		client := cache.GetRedisClient()
		if client == nil {
			c.Next()
			return
		}

		ip := c.ClientIP()
		key := fmt.Sprintf("rl:%s:%s", c.FullPath(), ip)

		count, err := client.Incr(c.Request.Context(), key).Result()
		if err != nil {
			c.Next()
			return
		}

		if count == 1 {
			client.Expire(c.Request.Context(), key, window)
		}

		if count > int64(limit) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"message": "Muitas requisições. Tente novamente mais tarde.",
			})
			return
		}

		c.Next()
	}
}
