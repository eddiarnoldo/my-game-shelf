package router

import (
	"github.com/eddiarnoldo/my-game-shelf/src/api/middleware"
	"github.com/gin-gonic/gin"
)

type BoardGameHandlerInterface interface {
	HandleBoardGameCreate(c *gin.Context)
	HandleGetBoardGames(c *gin.Context)
	HandleGetBoardGameByID(c *gin.Context)
	HandleBoardGameDelete(c *gin.Context)
	HandleBoardGameUpdate(c *gin.Context)
	HandleUploadBoardGameImage(c *gin.Context)
	HandleGetBoardGameCoverThumbnailImage(c *gin.Context)
	HandleGetBoardGameImage(c *gin.Context)
	HandleGetBoardGameImageThumbnail(c *gin.Context)
}

type AuthHandlerInterface interface {
	HandleRegister(c *gin.Context)
	HandleLogin(c *gin.Context)
	HandleRefresh(c *gin.Context)
	HandleLogout(c *gin.Context)
	HandleCreateInvite(c *gin.Context)
}

func RegisterRoutes(router *gin.Engine, boardGameHandler BoardGameHandlerInterface, jwtSecret string) {
	api := router.Group("/api")

	// Public read-only routes
	api.GET("/boardgames", boardGameHandler.HandleGetBoardGames)
	api.GET("/boardgames/:id", boardGameHandler.HandleGetBoardGameByID)
	api.GET("/boardgame/:id/images/coverThumbnail", boardGameHandler.HandleGetBoardGameCoverThumbnailImage)
	api.GET("/boardgame/:id/image/:imageId", boardGameHandler.HandleGetBoardGameImage)
	api.GET("/boardgame/:id/image/:imageId/thumbnail", boardGameHandler.HandleGetBoardGameImageThumbnail)

	// Protected write routes — require valid JWT
	protected := api.Group("")
	protected.Use(middleware.AuthRequired(jwtSecret))
	{
		protected.POST("/boardgame", boardGameHandler.HandleBoardGameCreate)
		protected.PUT("/boardgames/:id", boardGameHandler.HandleBoardGameUpdate)
		protected.DELETE("/boardgames/:id", boardGameHandler.HandleBoardGameDelete)
		protected.POST("/boardgame/:id/images", boardGameHandler.HandleUploadBoardGameImage)
	}
}

func RegisterAuthRoutes(router *gin.Engine, authHandler AuthHandlerInterface, jwtSecret string) {
	api := router.Group("/api")

	// Public auth routes
	api.POST("/register", authHandler.HandleRegister)
	api.POST("/login", authHandler.HandleLogin)
	api.POST("/refresh", authHandler.HandleRefresh)
	api.POST("/logout", authHandler.HandleLogout)

	// Admin-only routes
	admin := api.Group("")
	admin.Use(middleware.AuthRequired(jwtSecret), middleware.AdminRequired())
	{
		admin.POST("/invites", authHandler.HandleCreateInvite)
	}
}
