package handlers

import (
	"errors"
	"net/http"

	ce "p2p_interview/internal/domain/errors"
	"p2p_interview/internal/domain/models"
	"p2p_interview/internal/usecase"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Auth struct {
	user   *usecase.UserUseCase
	logger *zap.Logger
}

func NewAuth(user *usecase.UserUseCase, logger *zap.Logger) *Auth {
	return &Auth{user: user, logger: logger}
}

type RegisterUserInput struct {
	Login        string `json:"login" binding:"required"`
	PasswordHash string `json:"passwordhash" binding:"required"`
	Email        string `json:"email" binding:"required"`
	Name         string `json:"name" binding:"required"`
	Surname      string `json:"surname"`
}

// TODO: add error handling and returning right error message
func (a Auth) Register(c *gin.Context) {
	var r RegisterUserInput

	if err := c.ShouldBindJSON(&r); err != nil {
		a.logger.Error("Failed to bind JSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user := models.User{}

	user.Login = r.Login
	user.PasswordHash = r.PasswordHash
	user.Email = r.Email
	user.Name = r.Name
	user.Surname = r.Surname

	id, err := a.user.CreateUser(user)
	if err != nil {
		var valErr *ce.ValidationError
		if errors.As(err, &valErr) {
			a.logger.Warn("Validation error in user creation",
				zap.String("field", valErr.Field),
				zap.String("message", valErr.Message),
			)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		a.logger.Error("Error creating user", zap.Error(err))
		c.Error(err)
		return
	}

	a.logger.Info("User created successfully", zap.String("id", id))
	c.JSON(http.StatusCreated, gin.H{"id": id})
}
