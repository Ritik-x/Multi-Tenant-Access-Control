package handlers

import (
	"net/http"
	"team-access-control/internal/services"

	"github.com/gin-gonic/gin"
)

type InvitationHandler struct {
	invitationService *services.InviatationService
}
type CreateInvitationRequest struct {
	Email  string `json:"email" binding:"required,email"`
	RoleID string `json:"role_id" binding:"required"`
}
func NewServicesHandler(
	invitationService *services.InviatationService,
) *InvitationHandler{
	return &InvitationHandler{
		invitationService: invitationService,
	}

	
}

func (h *InvitationHandler) CreateInvitation(c *gin.Context) {

	var req CreateInvitationRequest
if err := c.ShouldBindBodyWithJSON(&req); err != nil {
	c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
}

organizationID := c.Param("organizationID")


if organizationID == "" {
	c.JSON(http.StatusBadRequest , gin.H{
		"error": "organization id is required",
	})
	return 
}

	c.JSON(http.StatusOK, gin.H{
	"email":           req.Email,
		"role_id":         req.RoleID,
		"organization_id": organizationID,
	})
}

