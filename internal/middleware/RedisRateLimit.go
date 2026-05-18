package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis_rate/v10"
	"github.com/redis/go-redis/v9"
)

func RedisRateLimit(
	client *redis.Client,
) gin.HandlerFunc {

	limiter := redis_rate.NewLimiter(client)

	return func(c *gin.Context) {

		res, err := limiter.Allow(
			c,
			c.ClientIP(),
			redis_rate.PerMinute(10),
		)

		if err != nil {

			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "rate limit failed",
			})

			c.Abort()

			return
		}

		if res.Allowed == 0 {

			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "too many requests",
			})

			c.Abort()

			return
		}

		c.Header(
			"X-RateLimit-Remaining",
			fmt.Sprintf("%d", res.Remaining),
		)

		c.Next()
	}
}
