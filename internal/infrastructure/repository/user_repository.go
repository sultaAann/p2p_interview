package repository

import (
	"database/sql"
	"p2p_interview/internal/domain/models"
	"p2p_interview/internal/repository"

	"go.uber.org/zap"
)

type userRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewUserRepository(db *sql.DB, logger *zap.Logger) repository.UserRepository {
	return &userRepository{db: db, logger: logger}
}

func (r *userRepository) GetAllUsers() ([]models.User, error) {
	r.db.Query("SELECT * FROM users")
	return nil, nil
}
