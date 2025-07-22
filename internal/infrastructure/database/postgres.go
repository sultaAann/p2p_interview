package database

import (
	"fmt"

	"database/sql"

	_ "github.com/lib/pq"

	"go.uber.org/zap"

	"p2p_interview/internal/infrastructure/config"
)

type Connection struct {
	logger *zap.Logger
}

func NewConnection(logger *zap.Logger) *Connection {
	return &Connection{logger: logger}
}

func (c *Connection) ConnectDB(info config.ConnectionDB) *sql.DB {
	c.logger.Info("Try to connect to DB")

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable", info.Host, info.Port, info.User, info.Password, info.Name)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		c.logger.Fatal("Can not open connection to database", zap.Error(err))
	}

	err = db.Ping()
	if err != nil {
		c.logger.Fatal("Can not to ping database", zap.Error(err))
	}

	c.logger.Info("Connected to DB")

	return db
}
