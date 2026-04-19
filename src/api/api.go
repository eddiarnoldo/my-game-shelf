package api

import (
	"fmt"
	"log"

	"github.com/eddiarnoldo/my-game-shelf/src/api/handlers"
	"github.com/eddiarnoldo/my-game-shelf/src/api/middleware"
	"github.com/eddiarnoldo/my-game-shelf/src/api/router"
	"github.com/eddiarnoldo/my-game-shelf/src/config"
	"github.com/eddiarnoldo/my-game-shelf/src/internal/repository"
	"github.com/eddiarnoldo/my-game-shelf/src/internal/services"

	"github.com/gin-gonic/gin"
)

func InitServer(
	boardGameRepo repository.BoardGameRepo,
	imageRepo repository.BoardGameImageRepo,
	userRepo repository.UserRepo,
	sessionRepo repository.SessionRepo,
	refreshTokenRepo repository.RefreshTokenRepo,
	inviteRepo repository.InviteRepo,
	emailService services.EmailService,
) error {
	jwtSecret := config.GetEnv("JWT_SECRET", "")
	if jwtSecret == "" {
		return fmt.Errorf("JWT_SECRET environment variable is required")
	}

	r := gin.Default()

	allowedOrigins := config.GetEnv("ALLOWED_ORIGINS", "*")
	r.Use(middleware.Cors(allowedOrigins))

	boardGameHandler := handlers.NewBoardGameHandler(boardGameRepo, imageRepo)
	authHandler := handlers.NewAuthHandler(
		userRepo, sessionRepo, refreshTokenRepo, inviteRepo, emailService, jwtSecret,
	)

	router.RegisterRoutes(r, boardGameHandler, jwtSecret)
	router.RegisterAuthRoutes(r, authHandler, jwtSecret)

	port := config.GetEnv("APP_PORT", "8080")
	log.Printf("Starting server on port %s...", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
		return err
	}
	return nil
}
