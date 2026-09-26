package handlers

import (
	"errors"
	"log"
	"net/http"
	"net/mail"
	"strings"
	"team-access-control/internal/repository"
	"team-access-control/internal/services"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type AuthHandler struct {
	authService          *services.AuthService
	userRepo             *repository.UserRepository
	membershipRepo       *repository.MemberRepository
	registerationService *services.RegistrationService
	sessionService       *services.SessionService

	sessionRepo *repository.SessionRepository
}

type RefreshTRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func NewAuthHandler(
	authService *services.AuthService,
	userRepo *repository.UserRepository,
	membershipRepo *repository.MemberRepository,
	registerationService *services.RegistrationService,
	sessionRepo *repository.SessionRepository,
	sessionService *services.SessionService,
) *AuthHandler {
	return &AuthHandler{
		authService:          authService,
		userRepo:             userRepo,
		membershipRepo:       membershipRepo,
		registerationService: registerationService,
		sessionRepo:          sessionRepo,
		sessionService:       sessionService,
	}
}

type RegisterRequest struct {
	Name             string `json:"name" binding:"required"`
	Email            string `json:"email" binding:"required"`
	Password         string `json:"password" binding:"required,min=8"`
	OrganizationName string `json:"organization_name" binding:"required"`
	OrganizationSlug string `json:"organization_slug" binding:"required"`
}
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	req.Name = strings.TrimSpace(req.Name)
	req.OrganizationName = strings.TrimSpace(req.OrganizationName)
	req.OrganizationSlug = strings.ToLower(
		strings.TrimSpace(req.OrganizationSlug),
	)

if req.Name == ""{
	c.JSON(http.StatusBadRequest , gin.H{
"error":"name is required", 
	})
}
_, err := mail.ParseAddress(req.Email)

if err != nil {
	c.JSON(http.StatusBadRequest, gin.H{
		"error": "invalid email format",
	})
	return
}

if req.Email == ""  {
	c.JSON(http.StatusBadRequest , gin.H{
"error":"email  is required", 
	})
}
if req.Password == "" {
	c.JSON(http.StatusBadRequest, gin.H{
		"error": "password is required",
	})
	return
}
if len(req.Password) < 8 {
	c.JSON(http.StatusBadRequest, gin.H{
		"error": "password must be at least 8 characters",
	})
	return
}

if req.OrganizationName == "" {
	c.JSON(http.StatusBadRequest, gin.H{
		"error": "organization name is required",
	})
	return
}

if req.OrganizationSlug == "" {
	c.JSON(http.StatusBadRequest, gin.H{
		"error": "organization slug is required",
	})
	return
}

	passwordHash, err := h.authService.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to hash password",
		})
		return
	}

	user, organizationID, _, err :=
		h.registerationService.Register(
			c.Request.Context(),
			req.Name,
			req.Email,
			passwordHash,
			req.OrganizationName,
			req.OrganizationSlug,
		)

	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"error": "failed to register user",
	// 	})
	// 	return
	// }
	if err != nil {
		log.Println("REGISTER ERROR:", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":              user.ID,
		"name":            user.Name,
		"email":           user.Email,
		"organization_id": organizationID,
		"created_at":      user.CreatedAt,
	})
}
func (h *AuthHandler) Login(c *gin.Context) {

	var req LoginRequest

	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	user, err := h.userRepo.GetUserByEmail(
		c.Request.Context(),
		req.Email,
	)

	if err != nil {
		log.Println("Login ERROR:", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if !h.authService.CheckPassword(
		req.Password,
		user.PasswordHash,
	) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid email or password",
		})
		return
	}

	organizationId, err :=
		h.membershipRepo.GetOrganizationByUserId(
			c.Request.Context(),
			user.ID,
		)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user is not a member of any organization",
		})
		return
	}

	accessToken, err := h.authService.GenerateAcessTokens(
		user.ID,
			user.Email,
		organizationId,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to generate access token",
		})
		return
	}

	refreshToken, refreshHashToken, err :=
		h.authService.GenerateRefreshToken()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to generate refresh token",
		})
		return
	}

	expiresAt := time.Now().Add(30 * 24 * time.Hour)

	_, err = h.sessionRepo.CreateSession(
		c.Request.Context(),
		user.ID,
		organizationId,
		refreshHashToken,
		expiresAt,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create session",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"token_type":    "Bearer",
	})
}

func (h *AuthHandler) Refresh(c *gin.Context) {

	var req RefreshTRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	// Rotate refresh token
	newRefreshToken, userID, organizationID, err :=
		h.sessionService.RotateRefreshToken(
			c.Request.Context(),
			req.RefreshToken,
		)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid refresh token",
		})
		return
	}

	// Get current user information
	user, err := h.userRepo.GetUserByID(
		c.Request.Context(),
		userID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get user",
		})
		return
	}

	// Generate new access token
	accessToken, err := h.authService.GenerateAcessTokens(
		user.ID,
		user.Email,
		organizationID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to generate access token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": newRefreshToken,
		"token_type":    "Bearer",
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	var req RefreshTRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	refreshTokenHash := h.authService.HashRefreshToken(
		req.RefreshToken,
	)

	sessionId, _, _, _, err := h.sessionRepo.GetSessionByRefreshTokenHash(c.Request.Context(), refreshTokenHash)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid refresh token",
		})
		return
	}
	if err := h.sessionRepo.RevokeSession(
		c.Request.Context(), sessionId,
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to revoke session",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "logged out successfully",
	})

}

func (h *AuthHandler) GetSessions(c *gin.Context) {
	userID := c.GetString("user_id")

	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "authentication context missing",
		})
		return
	}

	sessions, err := h.sessionService.GetActiveSession(c.Request.Context(),
		userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch sessions",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"sessions": sessions,
	})
}




func ( h *AuthHandler) DeleteRevoke(c *gin.Context){
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "authentication context missing",
		})
		return
	}
	sessionID := c.Param("id")

	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "session id is required",
		})
		return
	}
err := h.sessionService.RevokeSession(
    c.Request.Context(),
    sessionID,
    userID,
)
	if err !=nil {
		if errors.Is(err , pgx.ErrNoRows){
			c.JSON(http.StatusNotFound, gin.H{
				"error": "session not found",
			})
			return
		}
			c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to revoke session",
		})
		return
		}
		c.JSON(http.StatusOK, gin.H{
		"message": "session revoked successfully",
	})



	}




	func ( h *AuthHandler) RevokeAllsessions ( c *gin.Context){
		userID := c.GetString("user_id")
		if userID == ""{
			c.JSON(http.StatusUnauthorized , gin.H{
					"error": "authentication context missing",
			})
			return

		}
		err := h.sessionRepo.RevokeAllSessions(
			c.Request.Context(),userID,
		)
		if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to revoke sessions",
		})
		return
	}


	c.JSON(http.StatusOK, gin.H{
		"message": "all sessions revoked successfully",
	})
	}