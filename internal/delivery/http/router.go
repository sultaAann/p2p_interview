// internal/delivery/http/router.go
package http

import (
	"p2p_interview/internal/delivery/http/handlers"
	"p2p_interview/internal/delivery/http/middleware"

	_ "p2p_interview/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	
	"go.uber.org/zap"
)

// @title P2P Interview Management API
// @version 1.0
// @description REST API for managing P2P interview system
// @host localhost:8080
// @BasePath /api
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func SetupRouter(auth *handlers.Auth, user *handlers.UserHandlers, logger *zap.Logger) *gin.Engine {
	router := gin.Default()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler)) // docs

	router.Use(middleware.ErrorHandler(logger))

	router.POST("/register", auth.Register) // @Router /register [post]
	router.POST("/login", auth.Login)       // @Router /login [post]

	protected := router.Group("/api")
	{
		protected.Use(middleware.JwtAuthMiddleware(logger))
		protected.GET("/users", user.GetAllUsers) // @Router /users [get]
	}

	logger.Info("Router configured with routes: /register, /login, /api/users")
	return router
}
