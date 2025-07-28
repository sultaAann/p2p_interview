// internal/delivery/http/handlers/user_handlers.go
package handlers

import (
	"net/http"
	response "p2p_interview/internal/delivery/http/models"

	"p2p_interview/internal/usecase"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	_ "p2p_interview/docs"
)

type UserHandlers struct {
	usecase *usecase.UserUseCase
	logger  *zap.Logger
}

func NewUserHandler(usecase *usecase.UserUseCase, logger *zap.Logger) *UserHandlers {
	return &UserHandlers{usecase: usecase, logger: logger}
}

// @ID get-all-users
// @Summary Get all users
// @Description Returns a list of all users in the system.
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} object{data=[]models.UserResponse} "Successful response with user list"
// @Failure 404 {object} object{error=string} "No users found"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Security BearerAuth
// @Router /api/users [get]
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
	result := []response.UserResponse{}
	for _, user := range users {
		result = append(result, response.UserToUserResponse(user))
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}
