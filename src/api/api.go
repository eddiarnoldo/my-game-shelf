package api

import (
	"log"

	"github.com/eddiarnoldo/my-game-shelf/src/api/middleware"
	"github.com/eddiarnoldo/my-game-shelf/src/api/router"
	"github.com/eddiarnoldo/my-game-shelf/src/config"
	"github.com/eddiarnoldo/my-game-shelf/src/internal/handlers"
	"github.com/eddiarnoldo/my-game-shelf/src/internal/repository"
	"github.com/gin-gonic/gin"
)

func InitServer(boardGameRepo *repository.BoardGameRepository) error {
	//Create gin router
	r := gin.Default()

	allowedOrigins := config.GetEnv("ALLOWED_ORIGINS", "*")
	r.Use(middleware.Cors(allowedOrigins))

	// Initialize handlers
	boardGameHandler := handlers.NewBoardGameHandler(boardGameRepo)

	router.RegisterRoutes(r, boardGameHandler)

	// Start server
	port := config.GetEnv("APP_PORT", "8080")
	log.Printf("Starting server on port %s...", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
		return err
	}
	return nil
}
