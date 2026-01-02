package router

import (
	"github.com/eddiarnoldo/my-game-shelf/src/internal/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, boardGameHandler *handlers.BoardGameHandler) {
	api := router.Group("/api")
	{
		api.POST("/boardgame", boardGameHandler.CreateBoardGame)
		api.GET("/boardgames", boardGameHandler.GetAll)
		api.GET("/boardgames/:id", boardGameHandler.GetByID)
		api.DELETE("/boardgames/:id", boardGameHandler.Delete)
	}
}
