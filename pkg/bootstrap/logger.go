package bootstrap

import "go.uber.org/zap"

var logger *zap.Logger

func GetLogger() *zap.Logger {
	if logger != nil {
		return logger
	}

	cfg := GetConfig()
	switch cfg.Environment {
	case "production":
		logger, _ = zap.NewProduction()
	default:
		logger, _ = zap.NewDevelopment()
	}
	return logger
}

func SyncLogger() {
	if logger != nil {
		_ = logger.Sync()
	}
}
