package repository

import (
	"database/sql"
	"errors"
	"p2p_interview/internal/domain/models"
	"p2p_interview/internal/repository"
	ce "p2p_interview/internal/repository/errors"

	"go.uber.org/zap"
)

var _ repository.UserRepository = &userRepository{}

type userRepository struct {
	db          *sql.DB
	logger      *zap.Logger
	errorParser *ce.ErrorParser
}

func NewUserRepository(db *sql.DB, logger *zap.Logger, errorParser *ce.ErrorParser) repository.UserRepository {
	return &userRepository{db: db, logger: logger, errorParser: errorParser}
}

func (r *userRepository) GetAllUsers() ([]models.User, error) {
	r.logger.Info("Fetching all users")
	res, err := r.db.Query("SELECT id, login, email, name, surname FROM users")
	if err != nil {
		r.logger.Error("Failed to fetch all users", zap.Error(err))
		return nil, r.errorParser.ParsePostgresError(err, "get_all_users")
	}
	defer res.Close()

	users := []models.User{}
	for res.Next() {
		user := models.User{}
		err = res.Scan(&user.Id, &user.Login, &user.Email, &user.Name, &user.Surname)
		if err != nil {
			r.logger.Error("Failed to scan user row", zap.Error(err))
			return nil, r.errorParser.ParsePostgresError(err, "scan_user_row")
		}
		users = append(users, user)
	}

	if err = res.Err(); err != nil {
		r.logger.Error("Error iterating over rows", zap.Error(err))
		return nil, r.errorParser.ParsePostgresError(err, "iterate_user_rows")
	}

	r.logger.Info("Successfully fetched all users", zap.Int("count", len(users)))
	return users, nil
}

func (r *userRepository) GetUserById(id string) (*models.User, error) {
	r.logger.Info("Fetching user by ID", zap.String("id", id))

	var user models.User

	err := r.db.QueryRow(
		"SELECT login, email, name, surname FROM users WHERE id = $1", id,
	).Scan(&user.Login, &user.Email, &user.Name, &user.Surname)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.logger.Info("User not found", zap.String("id", id))
			return nil, &ce.NotFoundError{Resource: "User", ID: id}
		}
		r.logger.Error("Database error while fetching user", zap.Error(err))
		return nil, r.errorParser.ParsePostgresError(err, "get_user_by_id")
	}

	r.logger.Info("Successfully fetched user", zap.String("id", id))
	return &user, nil
}

func (r *userRepository) CreateUser(user models.User) (string, error) {
	r.logger.Info("Creating User", zap.String("login", user.Login))

	var id string
	err := r.db.QueryRow(
		"INSERT INTO users (login, passwordhash, email, name, surname) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		user.Login, user.PasswordHash, user.Email, user.Name, user.Surname,
	).Scan(&id)

	if err != nil {
		return "", r.errorParser.ParsePostgresError(err, "create_user")
	}

	r.logger.Info("Successfully created user", zap.String("id", id), zap.String("login", user.Login))
	return id, nil
}
