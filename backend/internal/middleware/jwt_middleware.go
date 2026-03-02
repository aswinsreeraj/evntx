package middleware

import (
	"net/http"
	"strings"

	jwtutil "github.com/aswinsreeraj/evntx/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		userID, err := jwtutil.ParseAccessToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false})
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}
