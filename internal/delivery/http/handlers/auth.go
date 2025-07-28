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

	_ "p2p_interview/docs"
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

// @ID register-user
// @Summary User registration
// @Description Creates a new user with the provided data.
// @Tags users
// @Accept json
// @Produce json
// @Param user body RegisterUserInput true "User registration data"
// @Success 201 {object} object{id=string} "Created user ID"
// @Failure 400 {object} object{error=string} "Validation error"
// @Failure 409 {object} object{error=string} "User already exists or constraint violation"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Security BearerAuth
// @Router /register [post]
func (a Auth) Register(c *gin.Context) {
	var r RegisterUserInput

	if err := c.ShouldBindJSON(&r); err != nil {
		a.logger.Error("Failed to bind JSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{
		Login:        r.Login,
		PasswordHash: r.Password,
		Email:        r.Email,
		Name:         r.Name,
		Surname:      r.Surname,
	}

	id, err := a.user.CreateUser(user)
	if err != nil {
		switch {
		case errors.As(err, new(*ce.ValidationError)):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.As(err, new(*pserr.DuplicateKeyError)):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		case errors.As(err, new(*pserr.ConstraintViolationError)):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		case errors.As(err, new(*pserr.DatabaseError)):
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		a.logger.Error("Error creating user", zap.Error(err))
		return
	}

	a.logger.Info("User created successfully", zap.String("id", id))
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

type LoginInput struct {
	Login    string `json:"login" binding:"required" example:"user1"`
	Password string `json:"password" binding:"required" example:"hashedpassword"`
}

// @ID login-user
// @Summary User login
// @Description Authenticates user credentials and returns a JWT token.
// @Tags users
// @Accept json
// @Produce json
// @Param credentials body LoginInput true "User login credentials"
// @Success 200 {object} object{token=string} "JWT access token"
// @Failure 400 {object} object{error=string} "Invalid input or credentials"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /login [post]
func (a Auth) Login(c *gin.Context) {
	var login LoginInput

	if err := c.ShouldBindJSON(&login); err != nil {
		a.logger.Error("Failed to bind JSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := a.user.LoginCheck(login.Login, login.Password)
	if err != nil {
		a.logger.Info("Username or password is incorrect", zap.String("login", login.Login))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username or password is incorrect."})
		return
	}

	a.logger.Info("Successfully logged in user", zap.String("login", login.Login))
	c.JSON(http.StatusOK, gin.H{"token": token})
}
