package logger

import (
	"context"
	"errors"
	"log"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	LogContext = "loggerContext"
	RequestId  = "requestId"
)

func NewZap() *zap.Logger {
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
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
