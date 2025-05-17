package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"user-management/app/user"
	usergrpc "user-management/app/user/grpc/gen/go/user/v1"
	"user-management/config"
	"user-management/logger"
	"user-management/middleware"
	"user-management/storage"

	echo "github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	ctx := context.Background()

	cfg, err := config.NewAppConfig()
	if err != nil {
		panic(err)
	}
	zlog := logger.NewZap()
	if err != nil {
		panic(err)
	}
	defer zlog.Sync()

	mongo := storage.InitMongoConnection(ctx, cfg.MongoDB)
	defer mongo.Disconnect(ctx)

	repo := user.NewRepository(mongo, cfg.MongoDB)
	uc := user.NewUsecase(cfg.Crypto, repo)
	handler := user.NewHandler(uc)

	go startRestServer(ctx, zlog, handler, cfg)
	// Start gRPC server
	go startGrpcServer(uc, zlog, cfg)

	go countTotalUserIntervalTicker(ctx, zlog, repo, cfg.UserCountInterval)

	// =================================== //
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-ctx.Done():
		zlog.Error("terminating: context cancelled")
	case <-sigterm:
		zlog.Error("terminating: via signal")
	}
}

func countTotalUserIntervalTicker(ctx context.Context, zlog *zap.Logger, repo user.CountUsersRepository, internval time.Duration) {
	ticker := time.NewTicker(internval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			count, err := repo.CountUsers(ctx)
			if err != nil {
				zlog.Sugar().Errorf("Error getting user count: %v", err)
			} else {
				zlog.Sugar().Infof("Current number of users: %d", count)
			}
		}
	}
}

func startRestServer(ctx context.Context, zlog *zap.Logger, handler user.Handler, cfg *config.AppConfig) {
	server := echo.New()
	defer server.Shutdown(ctx)
	server.Use(echoMiddleware.Recover())
	server.Use(middleware.HealthCheck)
	server.Use(middleware.NewLogging)
	server.Use(middleware.LoggingMiddleware)
	server.POST("/login", handler.Login)

	g := server.Group("", middleware.AuthMiddleware(cfg.Crypto.JwtKey))
	//CreateUser
	g.POST("/register", handler.CreateUser)
	// FindUsers
	g.GET("/users", handler.FindUsers)
	// FindUserById
	g.GET("/users/:id", handler.FindUserById)
	// UpdateUser
	g.PUT("/users/:id", handler.UpdateUser)
	// DeleteUser
	g.DELETE("/users/:id", handler.DeleteUser)

	err := server.Start(fmt.Sprintf(":%s", cfg.HttpServer.Port))
	if err != nil && err != http.ErrServerClosed {
		zlog.Sugar().Panicf("start http server error: %v", err)
	}

}

func startGrpcServer(usecase user.Usecase, zlog *zap.Logger, cfg *config.AppConfig) {
	grpcHandler := user.NewGrpcHandler(usecase)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.UnaryLoggingInterceptor(),
			middleware.GrpcAuthInterceptor(cfg.Crypto.JwtKey),
		),
	)

	usergrpc.RegisterUserServiceServer(grpcServer, grpcHandler)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GrpcServer.Port))
	if err != nil {
		zlog.Sugar().Fatalf("Failed to listen on port %s: %v", cfg.GrpcServer.Port, err)
	}

	zlog.Sugar().Infof("Starting gRPC server on port %s...", cfg.GrpcServer.Port)
	if err := grpcServer.Serve(listener); err != nil {
		zlog.Sugar().Fatalf("Failed to serve gRPC server: %v", err)
	}
}
