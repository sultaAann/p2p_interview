package http

import (
	"p2p_interview/internal/delivery/http/handlers"
	"p2p_interview/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func SetupRouter(handlers *handlers.UserHandlers, logger *zap.Logger) *gin.Engine {
	router := gin.Default()

	router.Use(middleware.ErrorHandler(logger))

	router.GET("/users", handlers.GetAllUsers)

	return router
}
