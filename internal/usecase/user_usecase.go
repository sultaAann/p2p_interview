package usecase

import (
	"errors"
	"fmt"

	ce "p2p_interview/internal/domain/errors"
	"p2p_interview/internal/domain/models"
	"p2p_interview/internal/infrastructure/security"
	"p2p_interview/internal/repository"
	pserr "p2p_interview/internal/repository/errors"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// TODO: handle custom errors | done
type UserUseCase struct {
	repository repository.UserRepository
	logger     *zap.Logger
}

func NewUserCase(repos repository.UserRepository, logger *zap.Logger) *UserUseCase {
	return &UserUseCase{repository: repos, logger: logger}
}

func (u *UserUseCase) GetAllUsers() ([]models.User, error) {
	u.logger.Info("Starting to retrieve all users")

	users, err := u.repository.GetAllUsers()
	if err != nil {
		var dbErr *pserr.DatabaseError
		if errors.As(err, &dbErr) {
			u.logger.Error("Failed to retrieve users due to database error",
				zap.Error(err),
				zap.String("operation", dbErr.Operation),
				zap.String("details", dbErr.Details),
			)
			return nil, fmt.Errorf("usecase failed to retrieve users due to database issue: %w", err)
		}
		u.logger.Error("Unexpected error in usecase",
			zap.Error(err),
		)
		return nil, fmt.Errorf("unexpected error in usecase: %w", err)
	}

	u.logger.Info("Successfully retrieved all users",
		zap.Int("user_count", len(users)),
	)
	return users, nil
}

func (u *UserUseCase) GetUserById(id string) (*models.User, error) {
	if id == "" {
		u.logger.Warn("Empty ID provided for user retrieval")
		return nil, fmt.Errorf("id cannot be empty")
	}

	u.logger.Info("Starting to retrieve user by ID", zap.String("id", id))
	user, err := u.repository.GetUserById(id)

	if err != nil {
		var nfErr *pserr.NotFoundError
		var dbErr *pserr.DatabaseError
		if errors.As(err, &nfErr) {
			u.logger.Info("User not found", zap.String("resource", nfErr.Resource), zap.String("id", nfErr.ID))
			return nil, err
		} else if errors.As(err, &dbErr) {
			u.logger.Error("Failed to retrieve user due to database error",
				zap.Error(err),
				zap.String("operation", dbErr.Operation),
				zap.String("details", dbErr.Details),
			)
			return nil, fmt.Errorf("usecase failed to retrieve user due to database issue: %w", err)
		}
		u.logger.Error("Unexpected error in usecase", zap.Error(err))
		return nil, fmt.Errorf("unexpected error in usecase: %w", err)
	}

	u.logger.Info("Successfully retrieved user by ID", zap.String("id", id))
	return user, nil
}

func (u *UserUseCase) CreateUser(user models.User) (string, error) {
	u.logger.Info("Starting to create user", zap.String("login", user.Login))

	hashedPassword, err := security.HashPassword(user.PasswordHash)
	if err != nil {
		u.logger.Error("Failed to hash password", zap.Error(err))
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	user.PasswordHash = hashedPassword
	if err := user.Validate(); err != nil {
		var valErr *ce.ValidationError
		if errors.As(err, &valErr) {
			u.logger.Warn("Validation failed for user creation",
				zap.String("field", valErr.Field),
				zap.String("message", valErr.Message),
			)
			return "", err
		}
	}
	id, err := u.repository.CreateUser(user)
	if err != nil {
		var dbErr *pserr.DatabaseError
		var dplcErr *pserr.DuplicateKeyError
		var cnstErr *pserr.ConstraintViolationError
		if errors.As(err, &dbErr) {
			u.logger.Error("Failed to create user due to database error",
				zap.Error(err),
				zap.String("operation", dbErr.Operation),
				zap.String("details", dbErr.Details),
			)
			return "", fmt.Errorf("usecase failed to create user due to database issue: %w", err)
		} else if errors.As(err, &dplcErr) {
			u.logger.Error("Failed to create user due to duplicate error",
				zap.Error(err),
				zap.String("resource", dplcErr.Resource),
				zap.String("filed", dplcErr.Field),
				zap.String("value", dplcErr.Value),
			)
			return "", err
		} else if errors.As(err, &cnstErr) {
			u.logger.Error("Failed to create user due to a row violation error",
				zap.Error(err),
				zap.String("constraint", cnstErr.Constraint),
				zap.String("details", cnstErr.Details),
				zap.String("operation", cnstErr.Operation),
			)
			return "", err
		}
		u.logger.Error("Unexpected error in usecase", zap.Error(err))
		return "", fmt.Errorf("unexpected error in usecase: %w", err)
	}

	u.logger.Info("Successfully created user", zap.String("id", id), zap.String("login", user.Login))
	return id, nil
}

func (u *UserUseCase) LoginCheck(login, password string) (string, error) {
	u.logger.Info("Login checking user", zap.String("login", login))

	user, err := u.repository.GetUserByLogin(login)
	if err != nil {
		var nfErr *pserr.NotFoundError
		var dbErr *pserr.DatabaseError
		if errors.As(err, &nfErr) {
			u.logger.Info("User not found", zap.String("resource", nfErr.Resource), zap.String("login", login))
			return "", err
		} else if errors.As(err, &dbErr) {
			u.logger.Error("Failed to retrieve user due to database error",
				zap.Error(err),
				zap.String("operation", dbErr.Operation),
				zap.String("details", dbErr.Details),
			)
			return "", fmt.Errorf("usecase failed to retrieve user due to database issue: %w", err)
		}
		u.logger.Error("Unexpected error in usecase", zap.Error(err))
		return "", fmt.Errorf("unexpected error in usecase: %w", err)
	}

	err = security.ComparePassword(password, user.PasswordHash)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		u.logger.Error("Failed to compare password", zap.Error(err))
		return "", fmt.Errorf("failed to compare password: %w", err)
	}
	
	u.logger.Info("Successfully checked user", zap.String("login", login))
}
