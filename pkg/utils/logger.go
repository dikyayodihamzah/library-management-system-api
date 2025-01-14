package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"

	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/zap"
)

var (
	DebugMode bool = false
	PrintLog  bool = false
	Log       []interface{}
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

	defer logger.Sync()
	return logger.Sugar()
}

func Debug(v ...any) {
	if GetBool("DEBUG_MODE") || DebugMode {
		_, file, line, ok := runtime.Caller(1)
		if ok {
			fmt.Printf("\033[32m %s line:%d \n \033[0m", file, line)
		}
		if PrintLog {
			Log = append(Log, v...)
		}
		log.Println(v...)
	}
}

func Json(v any) {
	if GetBool("DEBUG_MODE") {
		_, file, line, ok := runtime.Caller(1)
		if ok {
			fmt.Printf("\033[32m %s line:%d \n \033[0m", file, line)
		}
		by, _ := json.MarshalIndent(v, "", "   ")
		log.Println(string(by))
	}
}

func Error(log *zap.SugaredLogger, message string, err error) error {
	Debug(err)
	log.Errorw(message, "error", err)
	return err
}
