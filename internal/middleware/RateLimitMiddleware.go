package middleware

import (
	"time"

	ratelimit "github.com/JGLTechnologies/gin-rate-limit"
	"github.com/gin-gonic/gin"
)

func RateLimitMiddleware() gin.HandlerFunc {

	store := ratelimit.InMemoryStore(&ratelimit.InMemoryOptions{
		Rate:  time.Minute,
		Limit: 10,
	})

	return ratelimit.RateLimiter(store, &ratelimit.Options{
		ErrorHandler: func(c *gin.Context, info ratelimit.Info) {

			c.JSON(429, gin.H{
				"error": "too many requests",
			})
		},

		KeyFunc: func(c *gin.Context) string {
			return c.ClientIP()
		},
	})
}
