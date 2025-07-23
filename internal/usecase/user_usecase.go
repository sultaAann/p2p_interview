package usecase

import (
	"errors"
	"fmt"

	"p2p_interview/internal/domain/models"
	"p2p_interview/internal/repository"
	"p2p_interview/internal/usecase/interfaces"

	"go.uber.org/zap"
)

// var _ interfaces.UserUseCase = &userUseCaseImpl{}

type userUseCaseImpl struct {
	repository repository.UserRepository
	logger     *zap.Logger
}

func NewUserCase(repos repository.UserRepository, logger *zap.Logger) interfaces.UserUseCase {
	return &userUseCaseImpl{repository: repos, logger: logger}
}

func (u *userUseCaseImpl) GetAllUsers() ([]models.User, error) {
	u.logger.Info("Starting to retrieve all users")

	users, err := u.repository.GetAllUsers()
	if err != nil {
		var dbErr *repository.DatabaseError
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
