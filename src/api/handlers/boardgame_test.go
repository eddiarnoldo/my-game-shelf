package handlers

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/eddiarnoldo/my-game-shelf/src/internal/models"
	"github.com/gin-gonic/gin"
)

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
	handler := NewBoardGameHandler(repo)

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
	handler := NewBoardGameHandler(repo)

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

}

// We need to create the simplest mock repository
// Mocks in Go are about satisfying interfaces, not about test intent.
type mockBoardGameRepo struct {
	createCalled bool
}

func (m *mockBoardGameRepo) Create(ctx context.Context, game *models.BoardGame) error {
	m.createCalled = true
	return nil
}

func (m *mockBoardGameRepo) GetAll(ctx context.Context) ([]*models.BoardGame, error) {
	return nil, nil
}

func (m *mockBoardGameRepo) GetByID(ctx context.Context, id int64) (*models.BoardGame, error) {
	return nil, nil
}

func (m *mockBoardGameRepo) Delete(ctx context.Context, id int64) error {
	return nil
}
