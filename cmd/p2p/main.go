package main

import (
	"os"
	"p2p_interview/internal/infrastructure/config"
	"p2p_interview/internal/infrastructure/database"
	"p2p_interview/internal/infrastructure/logging"
	"p2p_interview/internal/infrastructure/repository"
)

func main() {
	logger := logging.NewLogger(os.Getenv("LOG_TYPE"))

	config := config.NewConfig(logger)

	DBCredentials := config.LoadDatabaseCredentials()

	DB := database.NewConnection(logger).ConnectDB(*DBCredentials)

	repository := repository.NewUserRepository(DB, logger)
}
