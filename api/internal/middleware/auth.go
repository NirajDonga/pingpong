package middleware

import (
	"net/http"

	"github.com/NirajDonga/pingpong/api/internal/auth"
	"github.com/gin-gonic/gin"
)

func Auth(authSvc auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie(auth.SessionCookieName)
		if err != nil || token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing session"})
			return
		}

		claims, err := authSvc.ValidateToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid session"})
			return
		}
		if claims.UserID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid session claims"})
			return
		}

		c.Set("userID", claims.UserID)
		c.Next()
	}
}
