package middleware

import (
	"net/http"
	"team-access-control/internal/services"

	"github.com/gin-gonic/gin"
)

func RequiredPermission(rbacservice *services.RBACService , requiredPermission string ) gin.HandlerFunc{
	return func(c *gin.Context){
		userID := c.GetString("user_id")
		organizationID := c.GetString("organization_id")
		if userID == "" || organizationID == ""{
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "authentication context missing",
			})
			c.Abort()
			return
		}

		hasPermission ,err := rbacservice.HasPermission(
			c.Request.Context(),
				userID,
			organizationID,
			requiredPermission,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to check permission",
			})
			c.Abort()
			return
		}
		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "permission denied",
			})
			c.Abort()
			return
		}
			c.Next()
	}
}