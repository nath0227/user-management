package logger

import (
	"context"
	"errors"
	"log"

	"go.uber.org/zap"
)

const (
	LogContext = "loggerContext"
	RequestId  = "requestId"
)

func NewZap() *zap.Logger {
	config := zap.NewProductionConfig()
	logger, err := config.Build(zap.AddCaller())
	if err != nil {
		log.Fatal(err)
	}

	return logger
}

func FromContext(context context.Context) (*zap.Logger, error) {
	l, ok := context.Value(LogContext).(*zap.Logger)
	if !ok {
		return nil, errors.New("unable get log from context")
	}
	return l, nil
}
