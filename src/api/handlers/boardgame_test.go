package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/eddiarnoldo/my-game-shelf/src/internal/models"
	"github.com/eddiarnoldo/my-game-shelf/src/internal/repository"
	"github.com/gin-gonic/gin"
)

/*
Implements

	type error interface {
	    Error() string
	}
*/
type ErrMockDBFailureType struct{}

func (e ErrMockDBFailureType) Error() string {
	return "mock database failure"
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func createTestContext(req *http.Request) (*gin.Context, *httptest.ResponseRecorder) {
	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	ctx.Request = req
	return ctx, rec
}

func TestHandleBoardGameCreate_OK(t *testing.T) {
	//Arrange
	repo := &mockBoardGameRepo{}
	handler := NewBoardGameHandler(repo, nil)

	body := []byte(`{
		"name": "Catan",
		"min_players": 4,
		"max_players": 10,
		"play_time": 60,
		"min_age": 8,
		"description": "a fun board game"
	}`)

	//
	req := httptest.NewRequest(http.MethodPost, "/api/boardgames", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	ctx, rec := createTestContext(req)

	// Act
	handler.HandleBoardGameCreate(ctx)

	// Assert
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d %s", rec.Code, rec.Body)
	}

	if !repo.createCalled {
		t.Fatal("expected Create() to be called on repository")
	}
}

func TestHandleBoardGameCreate_BadRequestJSON(t *testing.T) {
	// Arrange
	repo := &mockBoardGameRepo{}
	handler := NewBoardGameHandler(repo, nil)

	//Missing required fields
	body := []byte(`{
		"name": "Catan",
		"min_players": 4
	}`)

	req := httptest.NewRequest(http.MethodPost, "/api/boardgames", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	ctx, rec := createTestContext(req)

	// Act
	handler.HandleBoardGameCreate(ctx)

	// Assert
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}

	if repo.createCalled {
		t.Fatal("Create() should not be called on bad request")
	}
}

func TestHandleGetAllBoardgames_OK(t *testing.T) {
	// Arrange
	repo := &mockBoardGameRepo{}
	handler := NewBoardGameHandler(repo, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/boardgames", nil)
	ctx, rec := createTestContext(req)

	// Act
	handler.HandleGetBoardGames(ctx)

	// Assert
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	if !repo.getAllCalled {
		t.Fatal("expected GetAll() to be called on repository")
	}

	bodyBytes := rec.Body.Bytes()
	var response []*models.BoardGame

	err := json.Unmarshal(bodyBytes, &response)
	if err != nil {
		t.Fatalf("failed to unmarshal response JSON: %v", err)
	}

	//Test Data
	if len(response) != 1 {
		t.Fatalf("expected 1 board game, got %d", len(response))
	}

	bg := response[0]

	if bg.Name != "Honey Buzz" {
		t.Errorf("expected name 'Honey Buzz', got '%s'", bg.Name)
	}

	if bg.MinPlayers != 2 {
		t.Errorf("expected min_players 2, got %d", bg.MinPlayers)
	}

	if bg.MaxPlayers == 0 || bg.MaxPlayers != 4 {
		t.Errorf("expected max_players 4, got %v", bg.MaxPlayers)
	}

}

func TestHandleGetAllBoardGames_errorRepo(t *testing.T) {
	// Arrange
	var ErrMockDBFailure = ErrMockDBFailureType{}
	repo := &mockBoardGameRepo{
		getAllError: ErrMockDBFailure,
	}
	handler := NewBoardGameHandler(repo, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/boardgames", nil)
	ctx, rec := createTestContext(req)

	// Act
	handler.HandleGetBoardGames(ctx)

	// Assert
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", rec.Code)
	}

	if !repo.getAllCalled {
		t.Fatal("expected GetAll() to be called on repository")
	}
}

func TestHandleGetBoardGameByID_OK(t *testing.T) {
	// Arrange
	repo := &mockBoardGameRepo{}
	imageRepo := &mockBoardGameImageRepo{}
	handler := NewBoardGameHandler(repo, imageRepo)

	req := httptest.NewRequest(http.MethodGet, "/api/boardgames/1", nil)
	ctx, rec := createTestContext(req)
	ctx.Params = gin.Params{{Key: "id", Value: "1"}}

	// Act
	handler.HandleGetBoardGameByID(ctx)

	// Assert
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	if !repo.getByIDCalled {
		t.Fatal("expected GetByID() to be called on repository")
	}

	bodyBytes := rec.Body.Bytes()
	var response models.BoardGame

	err := json.Unmarshal(bodyBytes, &response)
	if err != nil {
		t.Fatalf("failed to unmarshal response JSON: %v", err)
	}

	//Test Data
	if response.Name != "Honey Buzz" {
		t.Errorf("expected name 'Honey Buzz', got '%s'", response.Name)
	}

	if response.MinPlayers != 2 {
		t.Errorf("expected min_players 2, got %d", response.MinPlayers)
	}

	if response.MaxPlayers == 0 || response.MaxPlayers != 4 {
		t.Errorf("expected max_players 4, got %v", response.MaxPlayers)
	}
}

func TestHandleGetBoardGameById_errorRepo(t *testing.T) {
	// Arrange
	var ErrMockDBFailure = ErrMockDBFailureType{}
	repo := &mockBoardGameRepo{
		getByIDError: ErrMockDBFailure,
	}
	handler := NewBoardGameHandler(repo, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/boardgames/1", nil)
	ctx, rec := createTestContext(req)
	ctx.Params = gin.Params{{Key: "id", Value: "1"}}

	// Act
	handler.HandleGetBoardGameByID(ctx)

	// Assert
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", rec.Code)
	}

	if !repo.getByIDCalled {
		t.Fatal("expected GetByID() to be called on repository")
	}
}

func TestHandleBoardGameDelete_NoContent(t *testing.T) {
	// Arrange
	repo := &mockBoardGameRepo{}
	handler := NewBoardGameHandler(repo, nil)

	req := httptest.NewRequest(http.MethodDelete, "/api/boardgames/1", nil)
	ctx, rec := createTestContext(req)
	ctx.Params = gin.Params{{Key: "id", Value: "1"}}

	// Act
	handler.HandleBoardGameDelete(ctx)

	// Assert
	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", rec.Code)
	}

	if !repo.deleteByIDCalled {
		t.Fatal("expected Delete() to be called on repository")
	}
}

func TestHandleBoardGameDelete_NotFound(t *testing.T) {
	// Arrange
	var ErrNotFound = repository.ErrBoardGameNotFound
	repo := &mockBoardGameRepo{
		deleteError: ErrNotFound,
	}
	handler := NewBoardGameHandler(repo, nil)

	req := httptest.NewRequest(http.MethodDelete, "/api/boardgames/999", nil)
	ctx, rec := createTestContext(req)
	ctx.Params = gin.Params{{Key: "id", Value: "999"}}

	// Act
	handler.HandleBoardGameDelete(ctx)

	// Assert
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", rec.Code)
	}

	if !repo.deleteByIDCalled {
		t.Fatal("expected Delete() to be called on repository")
	}
}

func TestHandleBoardGameDelete_InternalServerError(t *testing.T) {
	// Arrange
	var ErrMockDBFailure = ErrMockDBFailureType{}
	repo := &mockBoardGameRepo{
		deleteError: ErrMockDBFailure,
	}
	handler := NewBoardGameHandler(repo, nil)

	req := httptest.NewRequest(http.MethodDelete, "/api/boardgames/1", nil)
	ctx, rec := createTestContext(req)
	ctx.Params = gin.Params{{Key: "id", Value: "1"}}

	// Act
	handler.HandleBoardGameDelete(ctx)

	// Assert
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", rec.Code)
	}

	if !repo.deleteByIDCalled {
		t.Fatal("expected Delete() to be called on repository")
	}
}

func TestHandleBoardGameCreate_RepoError(t *testing.T) {
	// Arrange
	repo := &mockBoardGameRepo{createError: ErrMockDBFailureType{}}
	handler := NewBoardGameHandler(repo, nil)

	body := []byte(`{
		"name": "Catan",
		"min_players": 4,
		"max_players": 10,
		"play_time": 60,
		"min_age": 8,
		"description": "a fun board game"
	}`)

	req := httptest.NewRequest(http.MethodPost, "/api/boardgames", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ctx, rec := createTestContext(req)

	// Act
	handler.HandleBoardGameCreate(ctx)

	// Assert
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", rec.Code)
	}

	if !repo.createCalled {
		t.Fatal("expected Create() to be called on repository")
	}
}

func TestHandleGetBoardGameByID_InvalidID(t *testing.T) {
	// Arrange
	repo := &mockBoardGameRepo{}
	imageRepo := &mockBoardGameImageRepo{}
	handler := NewBoardGameHandler(repo, imageRepo)

	req := httptest.NewRequest(http.MethodGet, "/api/boardgames/abc", nil)
	ctx, rec := createTestContext(req)
	ctx.Params = gin.Params{{Key: "id", Value: "abc"}}

	// Act
	handler.HandleGetBoardGameByID(ctx)

	// Assert
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}

	if repo.getByIDCalled {
		t.Fatal("GetByID() should not be called on invalid ID")
	}
}

func TestHandleGetBoardGameByID_ImageRepoError(t *testing.T) {
	// Arrange
	repo := &mockBoardGameRepo{}
	imageRepo := &mockBoardGameImageRepo{getImageDtosError: ErrMockDBFailureType{}}
	handler := NewBoardGameHandler(repo, imageRepo)

	req := httptest.NewRequest(http.MethodGet, "/api/boardgames/1", nil)
	ctx, rec := createTestContext(req)
	ctx.Params = gin.Params{{Key: "id", Value: "1"}}

	// Act
	handler.HandleGetBoardGameByID(ctx)

	// Assert
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", rec.Code)
	}

	if !repo.getByIDCalled {
		t.Fatal("expected GetByID() to be called on repository")
	}

	if !imageRepo.getImageDtosCalled {
		t.Fatal("expected GetBoardGameImageDtos() to be called on image repository")
	}
}

func TestHandleBoardGameDelete_InvalidID(t *testing.T) {
	// Arrange
	repo := &mockBoardGameRepo{}
	handler := NewBoardGameHandler(repo, nil)

	req := httptest.NewRequest(http.MethodDelete, "/api/boardgames/abc", nil)
	ctx, rec := createTestContext(req)
	ctx.Params = gin.Params{{Key: "id", Value: "abc"}}

	// Act
	handler.HandleBoardGameDelete(ctx)

	// Assert
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}

	if repo.deleteByIDCalled {
		t.Fatal("Delete() should not be called on invalid ID")
	}
}

func TestHandleGetBoardGameCoverThumbnailImage_OK(t *testing.T) {
	// Arrange
	imageRepo := &mockBoardGameImageRepo{}
	handler := NewBoardGameHandler(nil, imageRepo)

	req := httptest.NewRequest(http.MethodGet, "/api/boardgame/1/images/cover", nil)
	ctx, rec := createTestContext(req)
	ctx.Params = gin.Params{{Key: "id", Value: "1"}}

	// Act
	handler.HandleGetBoardGameCoverThumbnailImage(ctx)

	// Assert
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	if !imageRepo.getCoverThumbnailCalled {
		t.Fatal("expected GetCoverThumbnail() to be called on image repository")
	}

	if rec.Body.Len() == 0 {
		t.Fatal("expected non-empty response body")
	}
}

func TestHandleGetBoardGameCoverThumbnailImage_InvalidID(t *testing.T) {
	// Arrange
	imageRepo := &mockBoardGameImageRepo{}
	handler := NewBoardGameHandler(nil, imageRepo)

	req := httptest.NewRequest(http.MethodGet, "/api/boardgame/abc/images/cover", nil)
	ctx, rec := createTestContext(req)
	ctx.Params = gin.Params{{Key: "id", Value: "abc"}}

	// Act
	handler.HandleGetBoardGameCoverThumbnailImage(ctx)

	// Assert
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}

	if imageRepo.getCoverThumbnailCalled {
		t.Fatal("GetCoverThumbnail() should not be called on invalid ID")
	}
}

func TestHandleGetBoardGameCoverThumbnailImage_NotFound(t *testing.T) {
	// Arrange
	imageRepo := &mockBoardGameImageRepo{getCoverThumbnailError: ErrMockDBFailureType{}}
	handler := NewBoardGameHandler(nil, imageRepo)

	req := httptest.NewRequest(http.MethodGet, "/api/boardgame/1/images/cover", nil)
	ctx, rec := createTestContext(req)
	ctx.Params = gin.Params{{Key: "id", Value: "1"}}

	// Act
	handler.HandleGetBoardGameCoverThumbnailImage(ctx)

	// Assert
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", rec.Code)
	}

	if !imageRepo.getCoverThumbnailCalled {
		t.Fatal("expected GetCoverThumbnail() to be called on image repository")
	}
}

func TestHandleGetBoardGameImage_OK(t *testing.T) {
	// Arrange
	imageRepo := &mockBoardGameImageRepo{}
	handler := NewBoardGameHandler(nil, imageRepo)

	req := httptest.NewRequest(http.MethodGet, "/api/boardgame/1/image/2", nil)
	ctx, rec := createTestContext(req)
	ctx.Params = gin.Params{
		{Key: "id", Value: "1"},
		{Key: "imageId", Value: "2"},
	}

	// Act
	handler.HandleGetBoardGameImage(ctx)

	// Assert
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	if !imageRepo.getImageDataCalled {
		t.Fatal("expected GetImageDataById() to be called on image repository")
	}

	if rec.Body.Len() == 0 {
		t.Fatal("expected non-empty response body")
	}
}

func TestHandleGetBoardGameImage_InvalidBoardGameID(t *testing.T) {
	// Arrange
	imageRepo := &mockBoardGameImageRepo{}
	handler := NewBoardGameHandler(nil, imageRepo)

	req := httptest.NewRequest(http.MethodGet, "/api/boardgame/abc/image/2", nil)
	ctx, rec := createTestContext(req)
	ctx.Params = gin.Params{
		{Key: "id", Value: "abc"},
		{Key: "imageId", Value: "2"},
	}

	// Act
	handler.HandleGetBoardGameImage(ctx)

	// Assert
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}

	if imageRepo.getImageDataCalled {
		t.Fatal("GetImageDataById() should not be called on invalid board game ID")
	}
}

func TestHandleGetBoardGameImage_InvalidImageID(t *testing.T) {
	// Arrange
	imageRepo := &mockBoardGameImageRepo{}
	handler := NewBoardGameHandler(nil, imageRepo)

	req := httptest.NewRequest(http.MethodGet, "/api/boardgame/1/image/abc", nil)
	ctx, rec := createTestContext(req)
	ctx.Params = gin.Params{
		{Key: "id", Value: "1"},
		{Key: "imageId", Value: "abc"},
	}

	// Act
	handler.HandleGetBoardGameImage(ctx)

	// Assert
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}

	if imageRepo.getImageDataCalled {
		t.Fatal("GetImageDataById() should not be called on invalid image ID")
	}
}

func TestHandleGetBoardGameImage_NotFound(t *testing.T) {
	// Arrange
	imageRepo := &mockBoardGameImageRepo{getImageDataError: ErrMockDBFailureType{}}
	handler := NewBoardGameHandler(nil, imageRepo)

	req := httptest.NewRequest(http.MethodGet, "/api/boardgame/1/image/2", nil)
	ctx, rec := createTestContext(req)
	ctx.Params = gin.Params{
		{Key: "id", Value: "1"},
		{Key: "imageId", Value: "2"},
	}

	// Act
	handler.HandleGetBoardGameImage(ctx)

	// Assert
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", rec.Code)
	}

	if !imageRepo.getImageDataCalled {
		t.Fatal("expected GetImageDataById() to be called on image repository")
	}
}

func TestHandleGetBoardGameImageThumbnail_OK(t *testing.T) {
	// Arrange
	imageRepo := &mockBoardGameImageRepo{}
	handler := NewBoardGameHandler(nil, imageRepo)

	req := httptest.NewRequest(http.MethodGet, "/api/boardgame/1/image/2/thumbnail", nil)
	ctx, rec := createTestContext(req)
	ctx.Params = gin.Params{
		{Key: "id", Value: "1"},
		{Key: "imageId", Value: "2"},
	}

	// Act
	handler.HandleGetBoardGameImageThumbnail(ctx)

	// Assert
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	if !imageRepo.getThumbnailDataCalled {
		t.Fatal("expected GetThumbnailImageDataById() to be called on image repository")
	}

	if rec.Body.Len() == 0 {
		t.Fatal("expected non-empty response body")
	}
}

func TestHandleGetBoardGameImageThumbnail_InvalidBoardGameID(t *testing.T) {
	// Arrange
	imageRepo := &mockBoardGameImageRepo{}
	handler := NewBoardGameHandler(nil, imageRepo)

	req := httptest.NewRequest(http.MethodGet, "/api/boardgame/abc/image/2/thumbnail", nil)
	ctx, rec := createTestContext(req)
	ctx.Params = gin.Params{
		{Key: "id", Value: "abc"},
		{Key: "imageId", Value: "2"},
	}

	// Act
	handler.HandleGetBoardGameImageThumbnail(ctx)

	// Assert
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}

	if imageRepo.getThumbnailDataCalled {
		t.Fatal("GetThumbnailImageDataById() should not be called on invalid board game ID")
	}
}

func TestHandleGetBoardGameImageThumbnail_InvalidImageID(t *testing.T) {
	// Arrange
	imageRepo := &mockBoardGameImageRepo{}
	handler := NewBoardGameHandler(nil, imageRepo)

	req := httptest.NewRequest(http.MethodGet, "/api/boardgame/1/image/abc/thumbnail", nil)
	ctx, rec := createTestContext(req)
	ctx.Params = gin.Params{
		{Key: "id", Value: "1"},
		{Key: "imageId", Value: "abc"},
	}

	// Act
	handler.HandleGetBoardGameImageThumbnail(ctx)

	// Assert
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}

	if imageRepo.getThumbnailDataCalled {
		t.Fatal("GetThumbnailImageDataById() should not be called on invalid image ID")
	}
}

func TestHandleGetBoardGameImageThumbnail_NotFound(t *testing.T) {
	// Arrange
	imageRepo := &mockBoardGameImageRepo{getThumbnailDataError: ErrMockDBFailureType{}}
	handler := NewBoardGameHandler(nil, imageRepo)

	req := httptest.NewRequest(http.MethodGet, "/api/boardgame/1/image/2/thumbnail", nil)
	ctx, rec := createTestContext(req)
	ctx.Params = gin.Params{
		{Key: "id", Value: "1"},
		{Key: "imageId", Value: "2"},
	}

	// Act
	handler.HandleGetBoardGameImageThumbnail(ctx)

	// Assert
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", rec.Code)
	}

	if !imageRepo.getThumbnailDataCalled {
		t.Fatal("expected GetThumbnailImageDataById() to be called on image repository")
	}
}

func TestHandleBoardGameUpdate_OK(t *testing.T) {
	// Arrange
	repo := &mockBoardGameRepo{}
	handler := NewBoardGameHandler(repo, nil)

	body := []byte(`{
		"name": "Catan Updated",
		"min_players": 3,
		"max_players": 6,
		"play_time": 90,
		"min_age": 10,
		"description": "updated description"
	}`)

	req := httptest.NewRequest(http.MethodPut, "/api/boardgames/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ctx, rec := createTestContext(req)
	ctx.Params = gin.Params{{Key: "id", Value: "1"}}

	// Act
	handler.HandleBoardGameUpdate(ctx)

	// Assert
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d %s", rec.Code, rec.Body)
	}

	if !repo.updateCalled {
		t.Fatal("expected Update() to be called on repository")
	}
}

func TestHandleBoardGameUpdate_InvalidID(t *testing.T) {
	// Arrange
	repo := &mockBoardGameRepo{}
	handler := NewBoardGameHandler(repo, nil)

	req := httptest.NewRequest(http.MethodPut, "/api/boardgames/abc", nil)
	req.Header.Set("Content-Type", "application/json")
	ctx, rec := createTestContext(req)
	ctx.Params = gin.Params{{Key: "id", Value: "abc"}}

	// Act
	handler.HandleBoardGameUpdate(ctx)

	// Assert
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}

	if repo.updateCalled {
		t.Fatal("Update() should not be called on invalid ID")
	}
}

func TestHandleBoardGameUpdate_BadRequestJSON(t *testing.T) {
	// Arrange
	repo := &mockBoardGameRepo{}
	handler := NewBoardGameHandler(repo, nil)

	body := []byte(`{"name": "Only Name"}`)

	req := httptest.NewRequest(http.MethodPut, "/api/boardgames/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ctx, rec := createTestContext(req)
	ctx.Params = gin.Params{{Key: "id", Value: "1"}}

	// Act
	handler.HandleBoardGameUpdate(ctx)

	// Assert
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}

	if repo.updateCalled {
		t.Fatal("Update() should not be called on bad request")
	}
}

func TestHandleBoardGameUpdate_NotFound(t *testing.T) {
	// Arrange
	repo := &mockBoardGameRepo{updateError: repository.ErrBoardGameNotFound}
	handler := NewBoardGameHandler(repo, nil)

	body := []byte(`{
		"name": "Catan",
		"min_players": 3,
		"max_players": 6,
		"play_time": 90,
		"min_age": 10,
		"description": "a description"
	}`)

	req := httptest.NewRequest(http.MethodPut, "/api/boardgames/999", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ctx, rec := createTestContext(req)
	ctx.Params = gin.Params{{Key: "id", Value: "999"}}

	// Act
	handler.HandleBoardGameUpdate(ctx)

	// Assert
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", rec.Code)
	}

	if !repo.updateCalled {
		t.Fatal("expected Update() to be called on repository")
	}
}

func TestHandleBoardGameUpdate_InternalServerError(t *testing.T) {
	// Arrange
	repo := &mockBoardGameRepo{updateError: ErrMockDBFailureType{}}
	handler := NewBoardGameHandler(repo, nil)

	body := []byte(`{
		"name": "Catan",
		"min_players": 3,
		"max_players": 6,
		"play_time": 90,
		"min_age": 10,
		"description": "a description"
	}`)

	req := httptest.NewRequest(http.MethodPut, "/api/boardgames/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ctx, rec := createTestContext(req)
	ctx.Params = gin.Params{{Key: "id", Value: "1"}}

	// Act
	handler.HandleBoardGameUpdate(ctx)

	// Assert
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", rec.Code)
	}

	if !repo.updateCalled {
		t.Fatal("expected Update() to be called on repository")
	}
}

// Helper mock repo and methods
// Mocks in Go are about satisfying interfaces, not about test intent.
type mockBoardGameRepo struct {
	createCalled     bool
	createError      error
	getAllCalled      bool
	getAllError       error
	getByIDCalled    bool
	getByIDError     error
	deleteByIDCalled bool
	deleteError      error
	updateCalled     bool
	updateError      error
}

func (m *mockBoardGameRepo) Create(ctx context.Context, game *models.BoardGame) error {
	m.createCalled = true
	return m.createError
}

func (m *mockBoardGameRepo) GetAll(ctx context.Context) ([]*models.BoardGame, error) {
	m.getAllCalled = true

	if m.getAllError != nil {
		return nil, m.getAllError
	}

	return []*models.BoardGame{
		{Name: "Honey Buzz", MinPlayers: 2, MaxPlayers: 4, PlayTime: 30, MinAge: 6, Description: "A sweet game"},
	}, nil
}

func (m *mockBoardGameRepo) GetByID(ctx context.Context, id int64) (*models.BoardGame, error) {
	m.getByIDCalled = true
	if m.getByIDError != nil {
		return nil, m.getByIDError
	}

	dummy := &models.BoardGame{Name: "Honey Buzz", MinPlayers: 2, MaxPlayers: 4, PlayTime: 30, MinAge: 6, Description: "A sweet game"}
	return dummy, nil
}

func (m *mockBoardGameRepo) Delete(ctx context.Context, id int64) error {
	m.deleteByIDCalled = true
	if m.deleteError != nil {
		return m.deleteError
	}
	return nil
}

func (m *mockBoardGameRepo) Update(ctx context.Context, game *models.BoardGame) error {
	m.updateCalled = true
	return m.updateError
}

type mockBoardGameImageRepo struct {
	saveImageCalled          bool
	saveImageError           error
	getImageDtosCalled       bool
	getImageDtosError        error
	getCoverThumbnailCalled  bool
	getCoverThumbnailError   error
	deleteImageCalled        bool
	deleteImageError         error
	getMaxDisplayOrderCalled bool
	getMaxDisplayOrderError  error
	getMaxDisplayOrderResult int
	getImageDataCalled       bool
	getImageDataError        error
	getThumbnailDataCalled   bool
	getThumbnailDataError    error
}

func (m *mockBoardGameImageRepo) SaveImage(ctx context.Context, image *models.BoardGameImage) error {
	m.saveImageCalled = true
	return m.saveImageError
}

func (m *mockBoardGameImageRepo) GetBoardGameImageDtos(ctx context.Context, boardGameId int64) ([]models.BoardGameImageDto, error) {
	m.getImageDtosCalled = true
	if m.getImageDtosError != nil {
		return nil, m.getImageDtosError
	}
	return []models.BoardGameImageDto{}, nil
}

func (m *mockBoardGameImageRepo) GetCoverThumbnail(ctx context.Context, boardGameId int64) (*models.BoardGameImage, error) {
	m.getCoverThumbnailCalled = true
	if m.getCoverThumbnailError != nil {
		return nil, m.getCoverThumbnailError
	}
	return &models.BoardGameImage{
		ImageMimeType: "image/jpeg",
		ThumbnailData: []byte("fake-thumbnail"),
	}, nil
}

func (m *mockBoardGameImageRepo) DeleteImage(ctx context.Context, id int64) error {
	m.deleteImageCalled = true
	return m.deleteImageError
}

func (m *mockBoardGameImageRepo) GetMaxDisplayOrder(ctx context.Context, boardGameId int64) (int, error) {
	m.getMaxDisplayOrderCalled = true
	return m.getMaxDisplayOrderResult, m.getMaxDisplayOrderError
}

func (m *mockBoardGameImageRepo) GetImageDataById(ctx context.Context, boardGameId int64, imageId int64) (*models.BoardGameImageData, error) {
	m.getImageDataCalled = true
	if m.getImageDataError != nil {
		return nil, m.getImageDataError
	}
	return &models.BoardGameImageData{
		Data:          []byte("fake-image"),
		ImageMimeType: "image/jpeg",
	}, nil
}

func (m *mockBoardGameImageRepo) GetThumbnailImageDataById(ctx context.Context, boardGameId int64, imageId int64) (*models.BoardGameImageData, error) {
	m.getThumbnailDataCalled = true
	if m.getThumbnailDataError != nil {
		return nil, m.getThumbnailDataError
	}
	return &models.BoardGameImageData{
		Data:          []byte("fake-thumbnail"),
		ImageMimeType: "image/jpeg",
	}, nil
}
