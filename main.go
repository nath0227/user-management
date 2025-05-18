package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
	"user-management/app/user"
	"user-management/config"
	"user-management/logger"
	"user-management/server"
	"user-management/storage"

	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	cfg, err := config.NewAppConfig()
	if err != nil {
		panic(err)
	}
	zlog := logger.NewZap()
	defer zlog.Sync()

	mongo := storage.InitMongoConnection(ctx, cfg.MongoDB)
	defer mongo.Disconnect(ctx)

	repo := user.NewRepository(mongo, cfg.MongoDB)
	uc := user.NewUsecase(cfg.Crypto, repo)
	handler := user.NewHandler(uc)

	grpcServer, err := server.NewGRPCServer(uc, zlog, cfg)
	if err != nil {
		panic(err)
	}
	httpServer := server.NewEchoHTTPServer(ctx, zlog, handler, cfg)
	// Start HTTP server
	go httpServer.Start()
	// Start gRPC server
	go grpcServer.Start()

	go countTotalUserIntervalTicker(ctx, zlog, repo, cfg.UserCountInterval)

	// =================================== //
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	<-ctx.Done() // Wait for termination signal
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	httpServer.Stop(shutdownCtx)
	grpcServer.Stop(shutdownCtx)
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
