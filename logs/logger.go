package log

import (
	"go.uber.org/zap"
)

func InitLogger() *zap.Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	return logger
}
