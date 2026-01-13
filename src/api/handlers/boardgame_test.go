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
	handler := NewBoardGameHandler(repo, nil)

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

// Helper mock repo and methods
// Mocks in Go are about satisfying interfaces, not about test intent.
type mockBoardGameRepo struct {
	createCalled     bool
	getAllCalled     bool
	getAllError      error
	getByIDCalled    bool
	getByIDError     error
	deleteByIDCalled bool
	deleteError      error
}

func (m *mockBoardGameRepo) Create(ctx context.Context, game *models.BoardGame) error {
	m.createCalled = true
	return nil
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

type mockBoardGameImageRepo struct {
	createCalled     bool
	getAllCalled     bool
	getByIDCalled    bool
	deleteByIDCalled bool
}

func (m *mockBoardGameImageRepo) SaveImage(ctx context.Context, image *models.BoardGameImage) error {
	m.createCalled = true
	return nil
}

func (m *mockBoardGameImageRepo) GetAllImagesForBoardGame(ctx context.Context, boardGameId int64, imageType string) ([]*models.BoardGameImage, error) {
	m.getAllCalled = true
	return []*models.BoardGameImage{}, nil
}

func (m *mockBoardGameImageRepo) GetCoverThumbnail(ctx context.Context, boardGameId int64) (*models.BoardGameImage, error) {
	m.getByIDCalled = true
	return &models.BoardGameImage{}, nil
}

func (m *mockBoardGameImageRepo) DeleteImage(ctx context.Context, id int64) error {
	m.deleteByIDCalled = true
	return nil
}
