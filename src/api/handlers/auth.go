package handlers

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/eddiarnoldo/my-game-shelf/src/api/middleware"
	"github.com/eddiarnoldo/my-game-shelf/src/internal/models"
	"github.com/eddiarnoldo/my-game-shelf/src/internal/repository"
	"github.com/eddiarnoldo/my-game-shelf/src/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	userRepo         repository.UserRepo
	sessionRepo      repository.SessionRepo
	refreshTokenRepo repository.RefreshTokenRepo
	inviteRepo       repository.InviteRepo
	emailService     services.EmailService
	jwtSecret        string
}

func NewAuthHandler(
	userRepo repository.UserRepo,
	sessionRepo repository.SessionRepo,
	refreshTokenRepo repository.RefreshTokenRepo,
	inviteRepo repository.InviteRepo,
	emailService services.EmailService,
	jwtSecret string,
) *AuthHandler {
	return &AuthHandler{
		userRepo:         userRepo,
		sessionRepo:      sessionRepo,
		refreshTokenRepo: refreshTokenRepo,
		inviteRepo:       inviteRepo,
		emailService:     emailService,
		jwtSecret:        jwtSecret,
	}
}

// HandleRegister creates a new user account from a valid invite code.
// POST /api/register
func (h *AuthHandler) HandleRegister(c *gin.Context) {
	var req struct {
		InviteCode           string `json:"invite_code" binding:"required"`
		Username             string `json:"username" binding:"required"`
		Password             string `json:"password" binding:"required"`
		PasswordConfirmation string `json:"password_confirmation" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Password != req.PasswordConfirmation {
		c.JSON(http.StatusBadRequest, gin.H{"error": "passwords do not match"})
		return
	}

	invite, err := h.inviteRepo.GetByCode(c, req.InviteCode)
	if errors.Is(err, repository.ErrInviteNotFound) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid invite code"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to validate invite"})
		return
	}
	if invite.Used {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invite has already been used"})
		return
	}
	if time.Now().After(invite.ExpiresAt) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invite has expired"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process password"})
		return
	}

	user := &models.User{
		Username:     req.Username,
		PasswordHash: string(hash),
		Role:         "user",
	}
	if err := h.userRepo.Create(c, user); err != nil {
		if errors.Is(err, repository.ErrDuplicateUsername) {
			c.JSON(http.StatusConflict, gin.H{"error": "username already taken"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	if err := h.inviteRepo.MarkUsed(c, invite.ID, user.ID); err != nil {
		// Non-fatal: user is created, just log
		_ = err
	}

	c.JSON(http.StatusCreated, gin.H{"message": "account created"})
}

// HandleLogin authenticates a user and returns an access token and refresh token.
// POST /api/login
func (h *AuthHandler) HandleLogin(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userRepo.GetByUsername(c, req.Username)
	if errors.Is(err, repository.ErrUserNotFound) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "authentication failed"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	session := &models.Session{
		ID:        newUUIDv4(),
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(3 * 24 * time.Hour),
	}
	if err := h.sessionRepo.Create(c, session); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create session"})
		return
	}

	accessToken, err := h.generateAccessToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	rawRefreshToken := newUUIDv4()
	refreshToken := &models.RefreshToken{
		SessionID: session.ID,
		TokenHash: hashToken(rawRefreshToken),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	if err := h.refreshTokenRepo.Create(c, refreshToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": rawRefreshToken,
	})
}

// HandleRefresh issues a new access token and refresh token using a valid refresh token.
// POST /api/refresh
func (h *AuthHandler) HandleRefresh(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hash := hashToken(req.RefreshToken)
	token, err := h.refreshTokenRepo.GetByTokenHash(c, hash)
	if errors.Is(err, repository.ErrTokenNotFound) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to validate token"})
		return
	}

	// Token theft detection: if already used, delete the session and reject
	if token.Used {
		_ = h.sessionRepo.Delete(c, token.SessionID)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token already used"})
		return
	}

	session, err := h.sessionRepo.GetByID(c, token.SessionID)
	if errors.Is(err, repository.ErrSessionNotFound) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "session not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to validate session"})
		return
	}

	if time.Now().After(session.ExpiresAt) {
		_ = h.sessionRepo.Delete(c, session.ID)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "session expired"})
		return
	}

	if time.Now().After(token.ExpiresAt) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token expired"})
		return
	}

	if err := h.refreshTokenRepo.MarkUsed(c, token.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to rotate token"})
		return
	}

	user, err := h.userRepo.GetByID(c, session.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user"})
		return
	}

	accessToken, err := h.generateAccessToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	rawRefreshToken := newUUIDv4()
	newToken := &models.RefreshToken{
		SessionID: session.ID,
		TokenHash: hashToken(rawRefreshToken),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	if err := h.refreshTokenRepo.Create(c, newToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": rawRefreshToken,
	})
}

// HandleLogout revokes the refresh token and deletes the associated session.
// POST /api/logout
func (h *AuthHandler) HandleLogout(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hash := hashToken(req.RefreshToken)
	token, err := h.refreshTokenRepo.GetByTokenHash(c, hash)
	if errors.Is(err, repository.ErrTokenNotFound) {
		// Idempotent — already gone
		c.Status(http.StatusNoContent)
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process logout"})
		return
	}

	_ = h.sessionRepo.Delete(c, token.SessionID)
	c.Status(http.StatusNoContent)
}

// HandleCreateInvite generates an invite code and emails it to the provided address.
// POST /api/invites — requires auth + admin role
func (h *AuthHandler) HandleCreateInvite(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDVal, _ := c.Get(middleware.ContextKeyUserID)
	userID, _ := userIDVal.(int64)

	// Remove any existing pending invite for this email
	_ = h.inviteRepo.DeletePendingByEmail(c, req.Email)

	code := newUUIDv4()
	invite := &models.Invite{
		Code:      code,
		Email:     req.Email,
		CreatedBy: userID,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}
	if err := h.inviteRepo.Create(c, invite); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create invite"})
		return
	}

	if err := h.emailService.SendInvite(req.Email, code); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to send invite email: %v", err)})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "invite sent"})
}

func (h *AuthHandler) generateAccessToken(user *models.User) (string, error) {
	claims := middleware.AuthClaims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.jwtSecret))
}

func hashToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

func newUUIDv4() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
