package usecase

import (
	"p2p_interview/internal/domain/models"
	"p2p_interview/internal/repository"

	"go.uber.org/zap"
)

type UserUseCaseImpl struct {
	repository repository.UserRepository
	logger     *zap.Logger
}

func (u *UserUseCaseImpl) GetAllUsers() ([]models.User, error) {
	users, err := u.repository.GetAllUsers()
	if err != nil {
		panic("DADA")
	}
	return users, nil
}
