package main

import (
	"auth-service/internal/config"
	"auth-service/internal/database"
	"auth-service/internal/handler"
	"auth-service/internal/repository"
	"auth-service/internal/routes"
	"auth-service/internal/service"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {

	cfg := config.LoadConfig()

	db := database.ConnectDB(cfg)

	redisClient := database.ConnectRedis(cfg)

	userRepo := repository.NewUserRepository(db)

	refreshRepo := repository.NewRefreshRepository(db)

	authService := service.NewAuthService(
		userRepo,
		refreshRepo,
		cfg,
	)

	authHandler := handler.NewAuthHandler(authService)

	router := gin.Default()

	routes.Setup(
		router,
		authHandler,
		cfg,
		redisClient,
	)

	log.Printf("Server running on port %s", cfg.AppPort)

	err := router.Run(":" + cfg.AppPort)

	if err != nil {
		log.Fatal(err)
	}
}
