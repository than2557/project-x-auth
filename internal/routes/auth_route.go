package routes

import (
	"auth-service/internal/handler"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(
	router *gin.Engine,
	authHandler *handler.AuthHandler,
) {

	auth := router.Group("/auth")

	{
		auth.POST("/register", authHandler.Register)

		auth.POST("/login", authHandler.Login)

		auth.POST("/refresh", authHandler.Refresh)

		auth.POST("/logout", authHandler.Logout)
	}
}
