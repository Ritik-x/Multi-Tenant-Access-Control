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
}
func NewAuthHandler (
	authService  *services.AuthService,
	userRepo *repository.UserRepository,
) *AuthHandler{
	return &AuthHandler{
		authService: authService,
		userRepo:    userRepo,
	}
}

type RegisterRequest struct {
	Name string `json:"name" binding:"required"`
		Email string `json:"email" binding:"required"`
			Password string `json:"password" binding:"required,min=8"`

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