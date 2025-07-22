package repository

import "p2p_interview/internal/domain/models"

type UserRepository interface {
	GetAllUsers() ([]models.User, error)
}
