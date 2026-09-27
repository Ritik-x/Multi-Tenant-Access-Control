package handlers

import (
	"net/http"
	"team-access-control/internal/services"

	"github.com/gin-gonic/gin"
)
type AuditLogHandlre struct {
	auditLogService *services.AuditLogService

}


func NewAuditLogHandler(auditLogService *services.AuditLogService,) *AuditLogHandlre{
	return &AuditLogHandlre{
		auditLogService: auditLogService,
}
}
// @Summary Get audit logs
// @Description Returns audit logs for the authenticated user's organization.
// @Tags Audit Logs
// @Produce json
// @Security BearerAuth
// @Param organizationID path string true "Organization ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /organizations/{organizationID}/audit-logs [get]
func( h *AuditLogHandlre) GetAuditLogs( c *gin.Context){
	organizationId := c.Param("organizationID")

	if organizationId == ""{
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "organization id is required",
		})
		return

	}


	tokenOrganizationId := c.GetString("organization_id")
	if tokenOrganizationId == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "authentication context missing",
		})
		return
	}


	if organizationId != tokenOrganizationId {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "organization access denied",
		})
		return
	}


	logs , err:= h.auditLogService.GetLogs(c.Request.Context() , organizationId )


	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch audit logs",
		})
		return
	}

	c.JSON(http.StatusOK , gin.H{
		"logs": logs,
	})
}