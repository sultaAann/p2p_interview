package config

import (
	"os"
	"strconv"

	"go.uber.org/zap"
)

type Config struct {
	logger *zap.Logger
}

type ConnectionDB struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

func NewConfig(logger *zap.Logger) *Config {
	return &Config{logger: logger}
}

func (c *Config) LoadDatabaseCredentials() *ConnectionDB {
	c.logger.Info("Loading Database Credentials")

	var CDB ConnectionDB

	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Can not load env fil")
	// }

	CDB.Host = os.Getenv("DB_HOST")
	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		c.logger.Error("Can not load enviroment variable", zap.Error(err))
	}
	CDB.Port = port
	CDB.User = os.Getenv("DB_USER")
	CDB.Password = os.Getenv("DB_PASSWORD")
	CDB.Name = os.Getenv("DB_NAME")

	c.logger.Info("Database Credentials loaded succesfully")

	return &CDB
}
