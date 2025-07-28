package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ErrorHandler(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			lastErr := c.Errors.Last()
			logger.Error("Request failed",
				zap.Error(lastErr),
				zap.String("path", c.Request.URL.Path),
				zap.Int("status", c.Writer.Status()),
			)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
	}
}
