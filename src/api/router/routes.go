package router

import (
	"github.com/eddiarnoldo/my-game-shelf/src/api/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, boardGameHandler *handlers.BoardGameHandler) {
	api := router.Group("/api")
	{
		api.POST("/boardgame", boardGameHandler.HandleBoardGameCreate)
		api.GET("/boardgames", boardGameHandler.HandleGetBoardGames)
		api.GET("/boardgames/:id", boardGameHandler.HandleGetBoardGameByID)
		api.DELETE("/boardgames/:id", boardGameHandler.HandleBoardGameDelete)
	}
}
