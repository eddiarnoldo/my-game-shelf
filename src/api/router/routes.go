package router

import (
	"github.com/gin-gonic/gin"
)

type BoardGameHandlerInterface interface {
	HandleBoardGameCreate(c *gin.Context)
	HandleGetBoardGames(c *gin.Context)
	HandleGetBoardGameByID(c *gin.Context)
	HandleBoardGameDelete(c *gin.Context)
	HandleUploadBoardGameImage(c *gin.Context)
}

func RegisterRoutes(router *gin.Engine, boardGameHandler BoardGameHandlerInterface) {
	api := router.Group("/api")
	{
		api.POST("/boardgame", boardGameHandler.HandleBoardGameCreate)
		api.GET("/boardgames", boardGameHandler.HandleGetBoardGames)
		api.GET("/boardgames/:id", boardGameHandler.HandleGetBoardGameByID)
		api.DELETE("/boardgames/:id", boardGameHandler.HandleBoardGameDelete)
		api.POST("/boardgame/:id/images", boardGameHandler.HandleUploadBoardGameImage)
		/*
			POST /api/boardgame/:id/images   → Upload image
			GET  /api/boardgame/:id/images/cover        → Get cover image (raw bytes)
			GET  /api/boardgame/images/:imageId         → Get any image by ID (raw bytes)
			DELETE /api/boardgame/images/:imageId       → Delete image
		*/
	}
}
