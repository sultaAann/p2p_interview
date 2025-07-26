package repository

import "p2p_interview/internal/domain/models"

type UserRepository interface {
	GetAllUsers() ([]models.User, error)
	GetUserById(id string) (*models.User, error)
	GetUserByLogin(login string) (*models.User, error)
	CreateUser(user models.User) (string, error)
}
