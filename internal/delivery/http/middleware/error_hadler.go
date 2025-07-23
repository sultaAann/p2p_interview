package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ErrorHandler(logger *zap.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
		if len(ctx.Errors) > 0 {
			lastErr := ctx.Errors.Last()
			logger.Error("Request failed",
				zap.Error(lastErr),
				zap.String("path", ctx.Request.URL.Path),
				zap.Int("status", ctx.Writer.Status()),
			)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
	}
}
