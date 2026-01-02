package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/eddiarnoldo/my-game-shelf/src/internal/models"
	"github.com/eddiarnoldo/my-game-shelf/src/internal/repository"
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
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create board game"})
		return
	}

	c.JSON(http.StatusCreated, game)
}

// Return All the board games
func (h *BoardGameHandler) GetAll(c *gin.Context) {
	boardGames, err := h.repo.GetAll(c.Request.Context())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, boardGames)
}

func (h *BoardGameHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")

	// Convert string to int64 (base, bits)
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	game, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Board game not found"})
		return
	}

	c.JSON(http.StatusOK, game)
}

func (h *BoardGameHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	err = h.repo.Delete(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrBoardGameNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Board game not found"})
			return
		}

		// Any other error is internal server error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.Status(http.StatusNoContent)
}
