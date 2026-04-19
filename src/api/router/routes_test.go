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
			name:   "POST /api/boardgame route is registered (protected)",
			method: http.MethodPost,
			path:   "/api/boardgame",
			// Auth middleware blocks without a valid JWT → 401, not 404.
			checkCalled: nil,
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
			name:        "DELETE /api/boardgames/:id route is registered (protected)",
			method:      http.MethodDelete,
			path:        "/api/boardgames/1",
			checkCalled: nil,
		},
		{
			name:        "PUT /api/boardgames/:id route is registered (protected)",
			method:      http.MethodPut,
			path:        "/api/boardgames/1",
			checkCalled: nil,
		},
		{
			name:        "POST /api/boardgame/:id/images route is registered (protected)",
			method:      http.MethodPost,
			path:        "/api/boardgame/1/images",
			checkCalled: nil,
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

			RegisterRoutes(router, mockHandler, "test-secret")

			// Act
			req := httptest.NewRequest(tt.method, tt.path, nil)
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			// Assert
			if rec.Code == http.StatusNotFound {
				t.Fatal("expected route to be registered, got 404")
			}

			if tt.checkCalled != nil && !tt.checkCalled(mockHandler) {
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

func TestRegisterAuthRoutes(t *testing.T) {
	tests := []struct {
		name        string
		method      string
		path        string
		checkCalled func(*mockAuthHandler) bool
	}{
		{
			name:   "POST /api/register calls HandleRegister",
			method: http.MethodPost,
			path:   "/api/register",
			checkCalled: func(m *mockAuthHandler) bool {
				return m.handleRegisterCalled
			},
		},
		{
			name:   "POST /api/login calls HandleLogin",
			method: http.MethodPost,
			path:   "/api/login",
			checkCalled: func(m *mockAuthHandler) bool {
				return m.handleLoginCalled
			},
		},
		{
			name:   "POST /api/refresh calls HandleRefresh",
			method: http.MethodPost,
			path:   "/api/refresh",
			checkCalled: func(m *mockAuthHandler) bool {
				return m.handleRefreshCalled
			},
		},
		{
			name:   "POST /api/logout calls HandleLogout",
			method: http.MethodPost,
			path:   "/api/logout",
			checkCalled: func(m *mockAuthHandler) bool {
				return m.handleLogoutCalled
			},
		},
		{
			name:   "POST /api/invites route is registered (admin-only)",
			method: http.MethodPost,
			path:   "/api/invites",
			// Auth middleware blocks without a valid JWT → 401, not 404.
			// checkCalled is nil: we only assert the route exists.
			checkCalled: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			router := gin.New()
			mockHandler := &mockAuthHandler{}

			RegisterAuthRoutes(router, mockHandler, "test-secret")

			req := httptest.NewRequest(tt.method, tt.path, nil)
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			if rec.Code == http.StatusNotFound {
				t.Fatal("expected route to be registered, got 404")
			}

			if tt.checkCalled != nil && !tt.checkCalled(mockHandler) {
				t.Fatal("expected handler method to be called")
			}
		})
	}
}

type mockAuthHandler struct {
	handleRegisterCalled     bool
	handleLoginCalled        bool
	handleRefreshCalled      bool
	handleLogoutCalled       bool
	handleCreateInviteCalled bool
}

func (m *mockAuthHandler) HandleRegister(c *gin.Context) {
	m.handleRegisterCalled = true
}
func (m *mockAuthHandler) HandleLogin(c *gin.Context) {
	m.handleLoginCalled = true
}
func (m *mockAuthHandler) HandleRefresh(c *gin.Context) {
	m.handleRefreshCalled = true
}
func (m *mockAuthHandler) HandleLogout(c *gin.Context) {
	m.handleLogoutCalled = true
}
func (m *mockAuthHandler) HandleCreateInvite(c *gin.Context) {
	m.handleCreateInviteCalled = true
}
