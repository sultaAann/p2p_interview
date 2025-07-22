package logging

import (
	"go.uber.org/zap"
)

func NewLogger(logLevel string) *zap.Logger {
	var logger *zap.Logger

	switch logLevel {
	case "production":
		logger, _ = zap.NewProduction()
	case "development":
		logger, _ = zap.NewDevelopment()
	}
	defer logger.Sync()
	return logger
}
