package middleware

import (
	"net/http"

	"github.com/aswinsreeraj/evntx/internal/domain"
	"github.com/aswinsreeraj/evntx/internal/repository"
	"github.com/gin-gonic/gin"
)

func RBACMiddleware(
	roleRepo repository.UserRoleRepository,
	allowedRoles ...domain.UserRole,
) gin.HandlerFunc {

	return func(c *gin.Context) {

		userIDVal, exists := c.Get("user_id")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false})
			return
		}

		userID := userIDVal.(string)

		roles, err := roleRepo.GetRolesByUserID(userID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"success": false})
			return
		}

		for _, role := range roles {
			for _, allowed := range allowedRoles {
				if role == allowed {
					c.Next()
					return
				}
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"success": false})
	}
}
