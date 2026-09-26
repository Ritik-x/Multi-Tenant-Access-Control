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