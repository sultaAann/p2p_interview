package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"p2p_interview/internal/delivery/http"
	"p2p_interview/internal/delivery/http/handlers"
	"p2p_interview/internal/infrastructure/config"
	"p2p_interview/internal/infrastructure/database"
	"p2p_interview/internal/infrastructure/logging"
	"p2p_interview/internal/infrastructure/repository"
	ce "p2p_interview/internal/repository/errors"
	"p2p_interview/internal/usecase"
	"syscall"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	err := godotenv.Load("../../config/.env")
	if err != nil {
		fmt.Println(err)
		log.Fatal("Failed load .env file")
	}

	logger := logging.NewLogger(os.Getenv("LOG_TYPE"))
	if logger == nil {
		log.Fatal("Failed to initialize logger")

	}

	config := config.NewConfig(logger)

	DBCredentials, err := config.LoadDatabaseCredentials()
	if err != nil {
		logger.Fatal("Failed to load database credentials", zap.Error(err))
		return
	}

	DB, err := database.NewConnection(logger).ConnectDB(*DBCredentials)
	if err != nil {
		logger.Fatal("Error Creating Connection To Database")
		return
	}

	psqlErrorParser := ce.NewErrorParser(logger)

	repository := repository.NewUserRepository(DB, logger, psqlErrorParser)

	userUseCase := usecase.NewUserCase(repository, logger)

	userHadlers := handlers.NewUserHandler(userUseCase, logger)
	auth := handlers.NewAuth(userUseCase, logger)

	router := http.SetupRouter(auth, userHadlers, logger)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	serverAddr := ":8080"
	logger.Info("Starting server", zap.String("address", serverAddr))
	if err := router.Run(serverAddr); err != nil {
		logger.Fatal("Server failed", zap.Error(err))
	}

	<-ctx.Done()
	logger.Info("Shutting down server")
}
