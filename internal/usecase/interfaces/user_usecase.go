package interfaces

import "p2p_interview/internal/domain/models"

type UserUseCase interface {
	GetAllUsers() ([]models.User, error)
}
