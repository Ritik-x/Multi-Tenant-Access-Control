package handlers

import (
	"log"
	"net/http"
	"strings"
	"team-access-control/internal/repository"
	"team-access-control/internal/services"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService          *services.AuthService
	userRepo             *repository.UserRepository
	membershipRepo       *repository.MemberRepository
	registerationService *services.RegistrationService
}

func NewAuthHandler(
	authService *services.AuthService,
	userRepo *repository.UserRepository,
	membershipRepo *repository.MemberRepository,
	registerationService *services.RegistrationService,
) *AuthHandler {
	return &AuthHandler{
		authService:          authService,
		userRepo:             userRepo,
		membershipRepo:       membershipRepo,
		registerationService: registerationService,
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
	// if err != nil {
	// 	c.JSON(http.StatusUnauthorized, gin.H{
	// 		"error": "invalid email or password",
	// 	})
	// 	return
	// }
	
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

	organizationId, err := h.membershipRepo.GetOrganizationByUserId(c.Request.Context(), user.ID)
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
	c.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
		"token_type":   "Bearer",
	})
}
