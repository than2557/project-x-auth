package routes

import (
	"auth-service/internal/config"
	"auth-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func ProtectedRoutes(
	router *gin.Engine,
	cfg *config.Config,
) {

	protected := router.Group("/api")

	protected.Use(middleware.JWTAuth(cfg))

	{
		protected.GET("/me", func(c *gin.Context) {

			c.JSON(200, gin.H{
				"user_id": c.GetString("user_id"),
				"email":   c.GetString("email"),
			})
		})
	}
}
