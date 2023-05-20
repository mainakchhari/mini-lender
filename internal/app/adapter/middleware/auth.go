package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mainakchhari/mini-lender/internal/app/domain"
	"github.com/mainakchhari/mini-lender/internal/app/domain/repository"
)

type UriParams struct {
	UserId int `uri:"uid" binding:"number"`
}

func BasicAuthMiddleware(allowedRoles []domain.Role, userRepository repository.IUser) gin.HandlerFunc {
	return func(c *gin.Context) {
		var uri UriParams
		c.ShouldBindUri(&uri)
		username, password, ok := c.Request.BasicAuth()
		// Check if basic auth is provided
		if ok {
			user, err := userRepository.Authenticate(username, password)
			// Check if authenticated user exists
			if err == nil {

				// Check if route access permitted
				if uri.UserId != 0 && user.ID != uri.UserId {
					c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden route access"})
					return
				}
				c.Set("user", user)

				// RBAC check
				if len(allowedRoles) == 0 {
					c.Next()
					return
				} else {
					for _, role := range allowedRoles {
						if user.Role == string(role) {
							c.Next()
							return
						}
					}
					c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
					return
				}
			}
		}
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	}
}
