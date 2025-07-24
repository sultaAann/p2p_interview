package handlers

import (
	"net/http"
	"p2p_interview/internal/domain/models"
	"p2p_interview/internal/usecase"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserHandlers struct {
	usecase *usecase.UserUseCase
	logger  *zap.Logger
}

func NewUserHandler(usecase *usecase.UserUseCase, logger *zap.Logger) *UserHandlers {
	return &UserHandlers{usecase: usecase, logger: logger}
}

func (h *UserHandlers) GetAllUsers(c *gin.Context) {
	h.logger.Info("Handling request to get all users")

	users, err := h.usecase.GetAllUsers()
	if err != nil {
		h.logger.Error("Error retrieving users", zap.Error(err))
		c.Error(err)
		return
	}

	if len(users) == 0 {
		h.logger.Info("No users found")
		c.JSON(http.StatusOK, gin.H{"message": "No users found", "data": []models.User{}})
		return
	}

	h.logger.Info("Successfully retrieved users", zap.Int("count", len(users)))
	c.JSON(http.StatusOK, gin.H{"data": users})
}
