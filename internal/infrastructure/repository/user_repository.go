package repository

import (
	"database/sql"
	"p2p_interview/internal/domain/models"
	"p2p_interview/internal/repository"

	"go.uber.org/zap"
)

// var _ repository.UserRepository = &userRepository{}

type userRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewUserRepository(db *sql.DB, logger *zap.Logger) repository.UserRepository {
	return &userRepository{db: db, logger: logger}
}

func (r *userRepository) GetAllUsers() ([]models.User, error) {
	r.logger.Info("Fetching all users")
	res, err := r.db.Query("SELECT id, login, email, name, surname FROM users")
	if err != nil {
		r.logger.Error("Failed to fetch all users", zap.Error(err))
		return nil, &repository.DatabaseError{Operation: "query", Details: err.Error()}
	}
	defer res.Close()

	users := []models.User{}
	for res.Next() {
		user := models.User{}
		err = res.Scan(&user.Id, &user.Login, &user.Email, &user.Name, &user.Surname)
		if err != nil {
			r.logger.Error("Failed to scan user row", zap.Error(err))
			return nil, &repository.DatabaseError{Operation: "scan", Details: err.Error()}
		}
		users = append(users, user)
	}

	if err = res.Err(); err != nil {
		r.logger.Error("Error iterating over rows", zap.Error(err))
		return nil, &repository.DatabaseError{Operation: "iteration", Details: err.Error()}
	}

	r.logger.Info("Successfully fetched all users", zap.Int("count", len(users)))
	return users, nil
}
