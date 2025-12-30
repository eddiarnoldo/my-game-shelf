package handlers

import (
	"net/http"

	"github.com/eddiarnoldo/my-game-shelf/internal/models"
	"github.com/eddiarnoldo/my-game-shelf/internal/repository"
	"github.com/gin-gonic/gin"
)

type BoardGameHandler struct {
	repo *repository.BoardGameRepository
}

// Use this function to create a new BoardGameHandler
func NewBoardGameHandler(repo *repository.BoardGameRepository) *BoardGameHandler {
	return &BoardGameHandler{repo: repo}
}

// This is the function that will handle the creation of a new board game
func (h *BoardGameHandler) CreateBoardGame(c *gin.Context) {
	var game models.BoardGame

	if err := c.ShouldBindJSON(&game); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.Create(c.Request.Context(), &game); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create board game"})
		return
	}

	c.JSON(http.StatusCreated, game)

}
