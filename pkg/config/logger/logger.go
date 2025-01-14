package logger

import (
	"log"
	"os"
	"strconv"

	"go.uber.org/zap"
)

func NewLogger() *zap.SugaredLogger {
	isProduction, _ := strconv.ParseBool(os.Getenv("IS_PRODUCTION"))
	isTesting, _ := strconv.ParseBool(os.Getenv("IS_TESTING"))

	var cfg zap.Config
	if isTesting {
		return zap.NewNop().Sugar()
	}

	if isProduction {
		cfg = zap.NewProductionConfig()
	} else {
		cfg = zap.NewDevelopmentConfig()
		cfg.DisableStacktrace = true
	}

	logger, err := cfg.Build()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	log.Println("Logger initiated")

	defer logger.Sync()
	return logger.Sugar()
}
