package handlers

import (
	"net/http"
	"strconv"
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
func (h *AuditLogHandlre) GetAuditLogs(c *gin.Context) {

	organizationID := c.Param("organizationID")

	if organizationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "organization id is required",
		})
		return
	}

	tokenOrganizationID := c.GetString("organization_id")

	if tokenOrganizationID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "authentication context missing",
		})
		return
	}

	if organizationID != tokenOrganizationID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "organization access denied",
		})
		return
	}

	// Pagination defaults
	page := 1
	limit := 20

	// Read page from query parameter
	if pageParam := c.Query("page"); pageParam != "" {

		value, err := strconv.Atoi(pageParam)

		if err != nil || value < 1 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "page must be a positive integer",
			})
			return
		}

		page = value
	}

	// Read limit from query parameter
	if limitParam := c.Query("limit"); limitParam != "" {

		value, err := strconv.Atoi(limitParam)

		if err != nil || value < 1 || value > 100 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "limit must be between 1 and 100",
			})
			return
		}

		limit = value
	}

	logs, err := h.auditLogService.GetLogs(
		c.Request.Context(),
		organizationID,
		page,
		limit,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch audit logs",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"logs": logs,
	})
}