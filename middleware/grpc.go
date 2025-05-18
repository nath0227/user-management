package middleware

import (
	"context"
	"log"
	"runtime/debug"
	"time"
	"user-management/logger"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
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
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {

		} else {
			zlog.Info("No gRPC metadata received")
		}
		// Get peer info
		p, _ := peer.FromContext(ctx)

		// Call the handler
		resp, err := handler(ctx, req)

		zlog.Info("gRPC request",
			zap.String("method", info.FullMethod),
			zap.Any("metadata", md),
			zap.Any("request", req),
			zap.Any("response", resp),
			zap.String("peer", p.Addr.String()),
			zap.Duration("duration", time.Since(start)),
			zap.Error(err),
		)

		return resp, err
	}
}

func UnaryInterceptorRecovery() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Recovered from panic: %v\nStack Trace: %s", r, debug.Stack())
				err = status.Errorf(codes.Internal, "internal server error")
			}
		}()
		return handler(ctx, req)
	}
}
