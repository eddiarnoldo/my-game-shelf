package main

import (
	"github.com/eddiarnoldo/my-game-shelf/internal/handlers"
	"github.com/gin-gonic/gin"
)

func setupRoutes(router *gin.Engine, boardGameHandler *handlers.BoardGameHandler) {
	api := router.Group("/api")
	{
		api.POST("/boardgame", boardGameHandler.CreateBoardGame)
		api.GET("/boardgames", boardGameHandler.GetAll)
		api.GET("/boardgames/:id", boardGameHandler.GetByID)
	}
}
