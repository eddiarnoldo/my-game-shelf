package main

import "github.com/gin-gonic/gin"

func setupRoutes(router *gin.Engine) {
	router.GET("/hello", handleHello)
}

func handleHello(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello, World!",
	})
}
