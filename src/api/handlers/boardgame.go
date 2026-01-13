package handlers

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/eddiarnoldo/my-game-shelf/src/internal/helpers"
	"github.com/eddiarnoldo/my-game-shelf/src/internal/models"
	"github.com/eddiarnoldo/my-game-shelf/src/internal/repository"
	"github.com/gin-gonic/gin"
)

type BoardGameHandler struct {
	repo      repository.BoardGameRepo
	imageRepo repository.BoardGameImageRepo
}

// Use this function to create a new BoardGameHandler
func NewBoardGameHandler(repo repository.BoardGameRepo, imageRepo repository.BoardGameImageRepo) *BoardGameHandler {
	return &BoardGameHandler{repo: repo, imageRepo: imageRepo}
}

// This is the function that will handle the creation of a new board game
func (h *BoardGameHandler) HandleBoardGameCreate(c *gin.Context) {
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

func (h *BoardGameHandler) HandleGetBoardGames(c *gin.Context) {
	boardGames, err := h.repo.GetAll(c.Request.Context())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, boardGames)
}

func (h *BoardGameHandler) HandleGetBoardGameByID(c *gin.Context) {
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

func (h *BoardGameHandler) HandleBoardGameDelete(c *gin.Context) {
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

	/*
		In Gin, when you call c.Status(), it sets the status code,
		but if the response body gets written afterward
		(which can happen automatically in some cases),
		Gin may override it with 200 OK.
	*/
	c.Status(http.StatusNoContent)
	c.Writer.WriteHeaderNow() // Force Gin to write the header immediately
}

// Image handlers
func (h *BoardGameHandler) HandleUploadBoardGameImage(c *gin.Context) {
	boardGameIDParam := c.Param("id")
	boardGameID, err := strconv.ParseInt(boardGameIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid board game ID"})
		return
	}

	// 2. Get the uploaded file
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No image provided"})
		return
	}

	// 3. Get image type from form
	imageType := c.PostForm("imageType") //TODO create a constants file
	if imageType != "cover" && imageType != "gameplay" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image type"})
		return
	}

	// 4. Validate file size (e.g., max 10MB)
	const maxFileSize = 10 * 1024 * 1024 // 10MB
	if file.Size > maxFileSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File too large (max 10MB)"})
		return
	}

	// 5. Validate MIME type
	if !strings.HasPrefix(file.Header.Get("Content-Type"), "image/") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File must be an image"})
		return
	}

	// 6. Open and read the file
	openedFile, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read image"})
		return
	}
	defer openedFile.Close()

	// 7. Read file bytes
	imageData, err := io.ReadAll(openedFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read image data"})
		return
	}

	// 8. Generate thumbnail
	thumbnailData, err := helpers.GenerateThumbnail(imageData, file.Header.Get("Content-Type"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate thumbnail"})
		return
	}

	// 9. Create image model
	image := &models.BoardGameImage{
		BoardGameID:   boardGameID,
		ImageData:     imageData,
		ImageMimeType: file.Header.Get("Content-Type"),
		ThumbnailData: thumbnailData,
		ImageType:     imageType,
		DisplayOrder:  0, // TODO: Calculate this
	}

	// 10. Save to database
	err = h.imageRepo.SaveImage(c.Request.Context(), image)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
		return
	}

	// 11. Return success with image ID
	c.JSON(http.StatusCreated, gin.H{
		"message": "Image uploaded successfully",
		"imageId": image.ID,
	})

}

func (h *BoardGameHandler) HandleGetBoardGameCoverImage(c *gin.Context) {
	boardGameIDParam := c.Param("id")
	boardGameID, err := strconv.ParseInt(boardGameIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid board game ID"})
		return
	}

	// 2. Get the cover thumbnail from repository
	image, err := h.imageRepo.GetCoverThumbnail(c.Request.Context(), boardGameID)
	if err != nil {
		// If no cover image exists, return 404
		c.JSON(http.StatusNotFound, gin.H{"error": "Cover image not found"})
		return
	}

	// 3. Set the Content-Type header (crucial!)
	c.Header("Content-Type", image.ImageMimeType)

	// 4. Optional: Add cache headers for better performance
	c.Header("Cache-Control", "public, max-age=86400") // Cache for 24 hours

	// 5. Write the thumbnail bytes directly to response
	c.Data(http.StatusOK, image.ImageMimeType, image.ThumbnailData)
}
