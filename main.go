package main

import (
	"context"
	"fmt"
	"log"
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

	go startRestServer(ctx, handler, cfg)
	// Start gRPC server
	go startGrpcServer(uc, cfg)

	go countTotalUser(ctx, repo)

	// =================================== //
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-ctx.Done():
		log.Println("terminating: context cancelled")
	case <-sigterm:
		log.Println("terminating: via signal")
	}
}

func countTotalUser(ctx context.Context, repo user.CountUsersRepository) {
	for {
		time.Sleep(10 * time.Second)
		count, _ := repo.CountUsers(ctx)
		log.Printf("User count: %d", count)
	}
}

func startRestServer(ctx context.Context, handler user.Handler, cfg *config.AppConfig) {
	server := echo.New()
	defer server.Shutdown(ctx)
	server.Use(echoMiddleware.Recover())
	server.Use(middleware.HealthCheck)
	server.Use(middleware.NewLogging)
	server.Use(middleware.LoggingMiddleware)
	server.POST("/login", handler.Login)

	g := server.Group("", middleware.JWTMiddleware(cfg.Crypto.JwtKey))
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
		log.Panicf("start http server error: %v", err)
	}

}

func startGrpcServer(usecase user.Usecase, cfg *config.AppConfig) {
	grpcHandler := user.NewGrpcHandler(usecase)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(
			middleware.UnaryLoggingInterceptor(),
		),
	)

	usergrpc.RegisterUserServiceServer(grpcServer, grpcHandler)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GrpcServer.Port))
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v\n", cfg.GrpcServer.Port, err)
	}

	log.Printf("Starting gRPC server on port %s...\n", cfg.GrpcServer.Port)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
