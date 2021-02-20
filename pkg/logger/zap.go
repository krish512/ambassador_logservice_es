package logger

import (
	"go.uber.org/zap"
)

// Logger is zap logger instance
var logger *zap.Logger

// InitLogger initiates new zap logger instance
func InitLogger() *zap.Logger {
	logger, _ = zap.NewProduction()
	// Logger, _ = zap.NewDevelopment()
	return logger
}
