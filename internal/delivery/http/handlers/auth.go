package handlers

import (
	"errors"
	"net/http"

	ce "p2p_interview/internal/domain/errors"
	"p2p_interview/internal/domain/models"
	pserr "p2p_interview/internal/repository/errors"
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
	Login    string `json:"login" binding:"required" example:"user1"`
	Password string `json:"password" binding:"required" example:"hashedpassword"`
	Email    string `json:"email" binding:"required" example:"user@example.com"`
	Name     string `json:"name" binding:"required" example:"John"`
	Surname  string `json:"surname" example:"Doe"`
}

func (a Auth) Register(c *gin.Context) {
	var r RegisterUserInput

	if err := c.ShouldBindJSON(&r); err != nil {
		a.logger.Error("Failed to bind JSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user := models.User{}

	user.Login = r.Login
	user.PasswordHash = r.Password
	user.Email = r.Email
	user.Name = r.Name
	user.Surname = r.Surname

	id, err := a.user.CreateUser(user)
	if err != nil {
		var valErr *ce.ValidationError
		var dbErr *pserr.DatabaseError
		var dplcErr *pserr.DuplicateKeyError
		var cnstErr *pserr.ConstraintViolationError

		if errors.As(err, &valErr) {
			a.logger.Warn("Validation error in user creation",
				zap.String("field", valErr.Field),
				zap.String("message", valErr.Message),
			)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		} else if errors.As(err, &dbErr) {
			a.logger.Error("Failed to create user due to database error",
				zap.Error(err),
				zap.String("operation", dbErr.Operation),
				zap.String("details", dbErr.Details),
			)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		} else if errors.As(err, &dplcErr) {
			a.logger.Error("Failed to create user due to duplicate error",
				zap.Error(err),
				zap.String("resource", dplcErr.Resource),
				zap.String("filed", dplcErr.Field),
				zap.String("value", dplcErr.Value),
			)
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		} else if errors.As(err, &cnstErr) {
			a.logger.Error("Failed to create user due to a row violation error",
				zap.Error(err),
				zap.String("constraint", cnstErr.Constraint),
				zap.String("details", cnstErr.Details),
				zap.String("operation", cnstErr.Operation),
			)
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		a.logger.Error("Error creating user", zap.Error(err))
		c.Error(err)
		return
	}

	a.logger.Info("User created successfully", zap.String("id", id))
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

type LoginInput struct {
	Login    string `json:"login" binding:"required" example:"user1"`
	Password string `json:"password" binding:"required" example:"hashedpassword"`
}

func (a Auth) Login(c *gin.Context) {
	var login LoginInput

	if err := c.ShouldBindJSON(&login); err != nil {
		a.logger.Error("Failed to bind JSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := a.user.LoginCheck(login.Login, login.Password)
	if err != nil {
		a.logger.Info("username or password is incorrect", zap.String("login", login.Login))
		c.JSON(http.StatusBadRequest, gin.H{"error": "username or password is incorrect."})
		return
	}
	a.logger.Info("Successfully login user", zap.String("login", login.Login))
	c.JSON(http.StatusOK, gin.H{"token": token})
}
