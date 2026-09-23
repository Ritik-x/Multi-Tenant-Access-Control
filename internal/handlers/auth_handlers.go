package handlers

import (
	"log"
	"net/http"
	"strings"
	"team-access-control/internal/repository"
	"team-access-control/internal/services"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService          *services.AuthService
	userRepo             *repository.UserRepository
	membershipRepo       *repository.MemberRepository
	registerationService *services.RegistrationService
		sessionRepo        *repository.SessionRepository
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
) *AuthHandler {
	return &AuthHandler{
		authService:          authService,
		userRepo:             userRepo,
		membershipRepo:       membershipRepo,
		registerationService: registerationService,
		sessionRepo :sessionRepo ,
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

func (h *AuthHandler) Refresh (c *gin.Context){
	var req RefreshTRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
			
		})
		return
	}

	refreshTOkenHash := h.authService.HashRefreshToken(req.RefreshToken)


	sessionID , userID , organizationID , expiresAT , err := h.sessionRepo.GetSessionByRefreshTokenHash(c.Request.Context() , refreshTOkenHash)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	if time.Now().After(expiresAT){
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":"refresh token expired",
		})
		return 
	}
	acessToken , err := h.authService.GenerateAcessTokens(userID,
		organizationID,)
		if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to generate access token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{


		"message": "refresh endpoint reached",
		"acess-Token":acessToken,
		"session-ID":sessionID,
	})
}



