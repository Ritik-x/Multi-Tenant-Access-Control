package middleware

import (
	"net/http"
	"team-access-control/internal/services"

	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(jwtSecret string ) gin.HandlerFunc{
	return func(c *gin.Context){
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "authorization header required",
		})
		c.Abort()
		return 
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) !=2 || parts[0] != "Bearer"{  //Bearer + token tht's why parts 2
		c.JSON(http.StatusUnauthorized ,gin.H{
			"error": "invalid authorization header",
		})
		c.Abort()
			return
	}
tokenStrings := parts[1]

token , err := jwt.ParseWithClaims(
	tokenStrings,
	&services.Claims{},
	func(token *jwt.Token) (interface {  } ,error){
		if token.Method != jwt.SigningMethodHS256{
			return nil , jwt.ErrSignatureInvalid
		}
			return []byte(jwtSecret), nil
	},
)
if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired token",
			})
			c.Abort()
			return


}

claims ,ok := token.Claims.(*services.Claims)

	if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token claims",
			})
			c.Abort()
			return
		}

		if claims.UserID == "" || claims.OrganizationID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "required token claims missing",
		})
		c.Abort()
		return
}
c.Set("user_id"  , claims.UserID)
	c.Set("organization_id", claims.OrganizationID)
	c.Next()
}
}