package handlers

import (
	"net/http"
	"team-access-control/internal/services"

	"github.com/gin-gonic/gin"
)

type InvitationHandler struct {
	invitationService *services.InvitationService
}
type CreateInvitationRequest struct {
	Email  string `json:"email" binding:"required,email"`
	RoleID string `json:"role_id" binding:"required"`
}
type AcceptInvitationRequest struct {
	Token string `json:"token" binding:"required"`
}
func NewServicesHandler(
	invitationService *services.InvitationService,
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

userID :=c.GetString("user_id")
if userID == "" {
	c.JSON(http.StatusUnauthorized, gin.H{
		"error": "user authentication required",
	})
	return
}

ipAddress := c.ClientIP()

tokenOrganizationId := c.GetString("organization_id")
if tokenOrganizationId == ""{
	c.JSON(http.StatusUnauthorized, gin.H{
		"error": "organization context missing",
	})
	return 
}
if organizationID != tokenOrganizationId{
	c.JSON(http.StatusForbidden , gin.H{
	"error": "organization access denied",
	})
	return 
}
	invitationToken, email, err :=
	h.invitationService.CreateInvitation(
		c.Request.Context(),
		organizationID,
		userID,
		req.Email,
		req.RoleID,
		ipAddress,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}


	c.JSON(http.StatusOK, gin.H{
	"email":           email,
		"role_id":         req.RoleID,
		"organization_id": organizationID,
			"token":   invitationToken,
				
	})
}





func ( h *InvitationHandler) AcceptInvitation(c *gin.Context ){
	var req AcceptInvitationRequest

	if err := c.ShouldBindBodyWithJSON(&req) ; err !=nil{
		c.JSON(http.StatusBadRequest , gin.H{
						"error": "invalid request",

		})
		return


	}
	userId := c.GetString("user_id")
	if userId == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "authentication context missing",
		})
		return
	}

	userEmail := c.GetString("user_email")

	if userEmail == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user email missing from authentication context",
		})
		return
	}

	err := h.invitationService.AcceptInvitatio(c.Request.Context() , req.Token , userEmail,userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}


	c.JSON(http.StatusOK, gin.H{
		"message": "invitation accepted successfully",
	})



}