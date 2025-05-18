package middleware

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"
	"user-management/logger"

	"github.com/google/uuid"
	echo "github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func NewLogging(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		requestId := uuid.New().String()
		zlog := logger.NewZap().With(zap.String("request_id", uuid.New().String()))
		ctx := context.WithValue(c.Request().Context(), logger.LogContext, zlog)
		ctx = context.WithValue(ctx, logger.RequestId, requestId)
		req := c.Request().WithContext(ctx)
		c.SetRequest(req)
		return next(c)
	}
}

func LoggingMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		zlog, err := logger.FromContext(c.Request().Context())
		start := time.Now()
		reqBody := new(bytes.Buffer)
		reqBody.ReadFrom(c.Request().Body)
		zlog.Sugar().With(
			zap.Any("header", c.Request().Header),
			zap.Any("body", reqBody.String()),
		).Infof("[REST API] %s, %s, Request", c.Request().Method, c.Path())
		c.Request().Body = io.NopCloser(reqBody)
		resBody := new(bytes.Buffer)
		multiWriter := io.MultiWriter(c.Response().Writer, resBody)
		writer := &CustomResponseWriter{Writer: multiWriter, ResponseWriter: c.Response().Writer}
		c.Response().Writer = writer
		err = next(c)

		httpPayload := logger.HTTPPayload{
			RequestMethod: c.Request().Method,
			RequestURL:    c.Request().URL.String(),
			Status:        c.Response().Status,
			Latency:       calLatency(start),
			ResponseSize:  fmt.Sprintf("%d", c.Response().Size),
		}
		zlog.Sugar().With(
			zap.Any("header", c.Response().Header()),
			zap.Any("body", resBody.String()),
			zap.Any("httpRequest", httpPayload),
		).Infof("[REST API] %s, %s, Response", c.Request().Method, c.Path())
		return err
	}
}

func calLatency(t time.Time) (output string) {
	diff := time.Since(t)
	return fmt.Sprintf("%.6fs", diff.Seconds())
}
