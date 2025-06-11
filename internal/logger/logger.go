package logger

import (
	"log"

	"go.uber.org/zap"
)

var logger *zap.Logger

func Init() {
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		log.Fatalln(err)
	}
}

func Logger() *zap.Logger {
	return logger
}
