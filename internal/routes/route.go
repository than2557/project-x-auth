package routes

import (
	"auth-service/internal/config"
	"auth-service/internal/handler"
	"auth-service/internal/middleware"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func Setup(
	router *gin.Engine,
	authHandler *handler.AuthHandler,
	cfg *config.Config,
	redisClient *redis.Client,
) {

	// =========================
	// Global Middleware
	// =========================

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
		},

		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
			"OPTIONS",
		},

		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Authorization",
		},

		AllowCredentials: true,

		MaxAge: 12 * time.Hour,
	}))

	router.Use(middleware.SecurityHeaders())

	router.Use(
		middleware.RedisRateLimit(redisClient),
	)

	router.SetTrustedProxies(nil)

	// =========================
	// Health
	// =========================

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "auth-service running",
		})
	})

	// =========================
	// Routes
	// =========================

	AuthRoutes(router, authHandler)

	ProtectedRoutes(router, cfg)
}
