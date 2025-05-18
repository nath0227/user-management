package server

import (
	"context"
	"fmt"
	"net"
	"user-management/app/user"
	usergrpc "user-management/app/user/grpc/gen/go/user/v1"
	"user-management/config"
	"user-management/logger"
	"user-management/middleware"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type GRPC struct {
	server   *grpc.Server
	listener net.Listener
}

func NewGRPCServer(usecase user.Usecase, zlog *zap.Logger, cfg *config.AppConfig) (*GRPC, error) {
	grpcHandler := user.NewGrpcHandler(usecase)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.UnaryInterceptorRecovery(),
			middleware.UnaryLoggingInterceptor(),
			middleware.GrpcAuthInterceptor(cfg.Crypto.JwtKey),
		),
	)

	usergrpc.RegisterUserServiceServer(grpcServer, grpcHandler)
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GrpcServer.Port))
	if err != nil {
		zlog.Sugar().Errorf("Failed to listen on port %s: %v", cfg.GrpcServer.Port, err)
		return nil, err
	}
	
	return &GRPC{
		server:   grpcServer,
		listener: listener,
	}, nil
}

func (g *GRPC) Start() {
	zlog := logger.NewZap()
	zlog.Info("Starting gRPC server on port 50051")
	if err := g.server.Serve(g.listener); err != nil {
		zlog.Sugar().Fatalf("Failed to start gRPC server: %v", err)
	}
}

func (g *GRPC) Stop(ctx context.Context) {
	zlog := logger.NewZap()
	zlog.Info("Shutting down gRPC server...")
	done := make(chan struct{})
	go func() {
		g.server.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
		zlog.Info("gRPC server shutdown complete.")
	case <-ctx.Done():
		zlog.Error("gRPC shutdown timeout exceeded, forcing stop.")
		g.server.Stop()
	}
}
