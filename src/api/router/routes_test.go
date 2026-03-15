package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRegisterRoutes(t *testing.T) {
	tests := []struct {
		name        string
		method      string
		path        string
		checkCalled func(*mockBoardGameHandler) bool
	}{
		{
			name:   "POST /api/boardgame calls HandleBoardGameCreate",
			method: http.MethodPost,
			path:   "/api/boardgame",
			checkCalled: func(m *mockBoardGameHandler) bool {
				return m.handleBoardGameCreateCalled
			},
		},
		{
			name:   "GET /api/boardgames calls HandleGetBoardGames",
			method: http.MethodGet,
			path:   "/api/boardgames",
			checkCalled: func(m *mockBoardGameHandler) bool {
				return m.handleGetBoardGamesCalled
			},
		},
		{
			name:   "GET /api/boardgames/:id calls HandleGetBoardGameByID",
			method: http.MethodGet,
			path:   "/api/boardgames/1",
			checkCalled: func(m *mockBoardGameHandler) bool {
				return m.handleGetBoardGameByIDCalled
			},
		},
		{
			name:   "DELETE /api/boardgames/:id calls HandleBoardGameDelete",
			method: http.MethodDelete,
			path:   "/api/boardgames/1",
			checkCalled: func(m *mockBoardGameHandler) bool {
				return m.handleBoardGameDeleteCalled
			},
		},
		{
			name:   "PUT /api/boardgames/:id calls HandleBoardGameUpdate",
			method: http.MethodPut,
			path:   "/api/boardgames/1",
			checkCalled: func(m *mockBoardGameHandler) bool {
				return m.handleBoardGameUpdateCalled
			},
		},
		{
			name:   "POST /api/boardgame/:id/images calls HandleUploadBoardGameImage",
			method: http.MethodPost,
			path:   "/api/boardgame/1/images",
			checkCalled: func(m *mockBoardGameHandler) bool {
				return m.handleUploadBoardGameImageCalled
			},
		},
		{
			name:   "GET /api/boardgame/:id/images/coverThumbnail calls HandleGetBoardGameCoverThumbnailImage",
			method: http.MethodGet,
			path:   "/api/boardgame/1/images/coverThumbnail",
			checkCalled: func(m *mockBoardGameHandler) bool {
				return m.handleGetBoardGameCoverThumbnailImageCalled
			},
		},
		{
			name:   "GET /api/boardgame/:id/image/:imageId calls HandleGetBoardGameImage",
			method: http.MethodGet,
			path:   "/api/boardgame/1/image/2",
			checkCalled: func(m *mockBoardGameHandler) bool {
				return m.handleGetBoardGameImageCalled
			},
		},
		{
			name:   "GET /api/boardgame/:id/image/:imageId/thumbnail calls HandleGetBoardGameImageThumbnail",
			method: http.MethodGet,
			path:   "/api/boardgame/1/image/2/thumbnail",
			checkCalled: func(m *mockBoardGameHandler) bool {
				return m.handleGetBoardGameImageThumbnailCalled
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			gin.SetMode(gin.TestMode)
			router := gin.New()
			mockHandler := &mockBoardGameHandler{}

			RegisterRoutes(router, mockHandler)

			// Act
			req := httptest.NewRequest(tt.method, tt.path, nil)
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			// Assert
			if rec.Code == http.StatusNotFound {
				t.Fatal("expected route to be registered, got 404")
			}

			if !tt.checkCalled(mockHandler) {
				t.Fatal("expected handler method to be called")
			}
		})
	}
}

type mockBoardGameHandler struct {
	handleBoardGameCreateCalled                bool
	handleGetBoardGamesCalled                  bool
	handleGetBoardGameByIDCalled               bool
	handleBoardGameDeleteCalled                bool
	handleBoardGameUpdateCalled                bool
	handleUploadBoardGameImageCalled           bool
	handleGetBoardGameCoverThumbnailImageCalled bool
	handleGetBoardGameImageCalled              bool
	handleGetBoardGameImageThumbnailCalled     bool
}

func (m *mockBoardGameHandler) HandleBoardGameCreate(c *gin.Context) {
	m.handleBoardGameCreateCalled = true
}
func (m *mockBoardGameHandler) HandleGetBoardGames(c *gin.Context) {
	m.handleGetBoardGamesCalled = true
}
func (m *mockBoardGameHandler) HandleGetBoardGameByID(c *gin.Context) {
	m.handleGetBoardGameByIDCalled = true
}
func (m *mockBoardGameHandler) HandleBoardGameDelete(c *gin.Context) {
	m.handleBoardGameDeleteCalled = true
}
func (m *mockBoardGameHandler) HandleBoardGameUpdate(c *gin.Context) {
	m.handleBoardGameUpdateCalled = true
}
func (m *mockBoardGameHandler) HandleUploadBoardGameImage(c *gin.Context) {
	m.handleUploadBoardGameImageCalled = true
}
func (m *mockBoardGameHandler) HandleGetBoardGameCoverThumbnailImage(c *gin.Context) {
	m.handleGetBoardGameCoverThumbnailImageCalled = true
}
func (m *mockBoardGameHandler) HandleGetBoardGameImage(c *gin.Context) {
	m.handleGetBoardGameImageCalled = true
}
func (m *mockBoardGameHandler) HandleGetBoardGameImageThumbnail(c *gin.Context) {
	m.handleGetBoardGameImageThumbnailCalled = true
}
