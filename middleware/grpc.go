package middleware

import (
	"context"
	"time"
	"user-management/logger"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

func UnaryLoggingInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()
		zlog := logger.NewZap().With(zap.String("request_id", uuid.New().String()))

		// Get peer info
		p, _ := peer.FromContext(ctx)

		// Call the handler
		resp, err := handler(ctx, req)

		zlog.Info("gRPC request",
			zap.String("method", info.FullMethod),
			zap.Any("request", req),
			zap.Any("response", resp),
			zap.String("peer", p.Addr.String()),
			zap.Duration("duration", time.Since(start)),
			zap.Error(err),
		)

		return resp, err
	}
}
