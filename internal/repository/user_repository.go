package repository

import (
	"database/sql"
	"p2p_interview/internal/domain/models"
	"p2p_interview/internal/domain/repository"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) repository.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetAllUsers() ([]models.User, error) {
	r.db.Query("SELECT * FROM users")
	return nil, nil
}
