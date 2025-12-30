package main

import (
	"github.com/eddiarnoldo/my-game-shelf/internal/handlers"
	"github.com/gin-gonic/gin"
)

func setupRoutes(router *gin.Engine, boardGameHandler *handlers.BoardGameHandler) {
	router.POST("/api/game", boardGameHandler.CreateBoardGame)
}
