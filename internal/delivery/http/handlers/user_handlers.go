// internal/delivery/http/handlers/user_handlers.go
package handlers

import (
	"net/http"
	"p2p_interview/internal/usecase"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	_ "p2p_interview/docs"
)

// @title P2P Interview Management API
// @version 1.0
// @description REST API for managing P2P interview system
// @host localhost:8080
// @BasePath /api
// @schemes http
type UserHandlers struct {
	usecase *usecase.UserUseCase
	logger  *zap.Logger
}

func NewUserHandler(usecase *usecase.UserUseCase, logger *zap.Logger) *UserHandlers {
	return &UserHandlers{usecase: usecase, logger: logger}
}

// @Summary Get all users
// @Description Retrieve a list of all users
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {array} models.User
// @Failure 404 {object} object "No users found"
// @Failure 500 {object} object "Internal server error"
// @Security BearerAuth
// @Router /users [get]
func (h *UserHandlers) GetAllUsers(c *gin.Context) {
	h.logger.Info("Handling request to get all users")

	users, err := h.usecase.GetAllUsers()
	if err != nil {
		h.logger.Error("Error retrieving users", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if len(users) == 0 {
		h.logger.Info("No users found")
		c.JSON(http.StatusNotFound, gin.H{"error": "No users found"})
		return
	}

	h.logger.Info("Successfully retrieved users", zap.Int("count", len(users)))
	c.JSON(http.StatusOK, gin.H{"data": users})
}
