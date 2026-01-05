package router

import (
	"github.com/gin-gonic/gin"
)

type BoardGameHandlerInterface interface {
	HandleBoardGameCreate(c *gin.Context)
	HandleGetBoardGames(c *gin.Context)
	HandleGetBoardGameByID(c *gin.Context)
	HandleBoardGameDelete(c *gin.Context)
}

func RegisterRoutes(router *gin.Engine, boardGameHandler BoardGameHandlerInterface) {
	api := router.Group("/api")
	{
		api.POST("/boardgame", boardGameHandler.HandleBoardGameCreate)
		api.GET("/boardgames", boardGameHandler.HandleGetBoardGames)
		api.GET("/boardgames/:id", boardGameHandler.HandleGetBoardGameByID)
		api.DELETE("/boardgames/:id", boardGameHandler.HandleBoardGameDelete)
	}
}
