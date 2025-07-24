package http

import (
	"p2p_interview/internal/delivery/http/handlers"
	"p2p_interview/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func SetupRouter(auth *handlers.Auth, user *handlers.UserHandlers, logger *zap.Logger) *gin.Engine {
	router := gin.Default()

	router.Use(middleware.ErrorHandler(logger))

	router.POST("/register", auth.Register)

	router.GET("/users", user.GetAllUsers)

	return router
}
