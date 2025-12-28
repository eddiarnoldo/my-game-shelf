package main

import (
	"github.com/eddiarnoldo/my-game-shelf/internal/models"
	"github.com/gin-gonic/gin"
)

var games = []models.BoardGame{
	{ID: 1, Name: "One Night Werewolf", MinPlayers: 3, MaxPlayers: 4, PlayTime: 60, AgeRating: 10, Description: "En este juego de roles y deducción, los jugadores asumen identidades secretas."},
	{ID: 2, Name: "Love Letter", MinPlayers: 3, MaxPlayers: 4, PlayTime: 15, AgeRating: 8, Description: "Asume el papel de un pretendiente que intenta entregar una carta de amor a la princesa. Elimina a tus openentes usando cartas con diferentes habilidades."},
}

func setupRoutes(router *gin.Engine) {
	router.GET("/api/hello", handleHello)
	router.GET("/api/games", handleListGames)
}

func handleHello(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello, World!",
	})
}

func handleListGames(c *gin.Context) {
	c.IndentedJSON(200, games)
}
