package handlers

import (
	"net/http"
	"strings"
	"team-access-control/internal/models"
	"team-access-control/internal/repository"
	"team-access-control/internal/services"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *services.AuthService
	userRepo *repository.UserRepository
	membershipRepo   *repository.MemberRepository
}
func NewAuthHandler (
	authService  *services.AuthService,
	userRepo *repository.UserRepository,
	membershipRepo   *repository.MemberRepository,
) *AuthHandler{
	return &AuthHandler{
		authService: authService,
		userRepo:    userRepo,
		membershipRepo: membershipRepo,
	}
}

type RegisterRequest struct {
	Name string `json:"name" binding:"required"`
		Email string `json:"email" binding:"required"`
			Password string `json:"password" binding:"required,min=8"`

}
type LoginRequest struct {
	Email string  `json:"email" binding:"required,email"`
	Password string  `json:"password" binding:"required"`

}

func ( h *AuthHandler ) Register ( c *gin.Context){
	var req RegisterRequest

	if err := c.ShouldBindBodyWithJSON(&req); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
		req.Name = strings.TrimSpace(req.Name)
		passwordHash , err  := h.authService.HashPassword(req.Password)
		if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to hash password",
		})
		return
}
user := &models.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: passwordHash,
}

if err := h.userRepo.CreateUser(
	c.Request.Context(),user,
);err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create user",
		})
		return
	}
		c.JSON(http.StatusCreated , gin.H{

		 	"id":         user.ID,
		"name":       user.Name,
		"email":      user.Email,
		"created_at": user.CreatedAt,})
}
func ( h *AuthHandler) Login ( c *gin.Context){

	var req LoginRequest
	if err := c.ShouldBindBodyWithJSON(&req) ; err !=nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	req.Email =strings.ToLower(strings.TrimSpace(req.Email))

	user , err := h.userRepo.GetUserByMail(
		c.Request.Context(),
		req.Email,
	)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid email or password",
		})
		return
	}
	if !h.authService.CheckPassword(
		req.Password,
		user.PasswordHash,
	){c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid email or password",
		})
		return
	}



	organizationId , err := h.membershipRepo.GetOrganizationByUserId(c.Request.Context(),user.ID,
)
if err != nil {
	c.JSON(http.StatusUnauthorized, gin.H{
		"error": "user is not a member of any organization",
	})
	return
}
	accessToken , err := h.authService.GenerateAcessTokens(
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