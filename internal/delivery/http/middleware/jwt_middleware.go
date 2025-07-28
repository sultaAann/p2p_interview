package middleware

import (
	"net/http"

	"p2p_interview/internal/infrastructure/security"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AuthMiddleware provides authentication middleware for protected routes.
// It validates the JWT token from the Authorization header and sets the user context if valid.
// Parameters:
//   - logger: Zap logger for authentication logging
//
// Returns:
//   - gin.HandlerFunc: Middleware function to be applied to Gin routes
func JwtAuthMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := security.TokenValid(c)
		if err != nil {
			lastErr := c.Errors.Last()
			logger.Error("Request failed",
				zap.Error(lastErr),
				zap.String("path", c.Request.URL.Path),
				zap.Int("status", c.Writer.Status()),
			)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}
