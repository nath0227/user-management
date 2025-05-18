package server

import (
	"context"
	"fmt"
	"net/http"
	"user-management/app/user"
	"user-management/config"
	"user-management/logger"
	"user-management/middleware"

	echo "github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

type HTTP struct {
	server *http.Server
}

func NewEchoHTTPServer(ctx context.Context, zlog *zap.Logger, handler user.Handler, cfg *config.AppConfig) *HTTP {
	server := echo.New()
	server.Server.Addr =fmt.Sprintf(":%s", cfg.HttpServer.Port)
	server.Use(echoMiddleware.Recover())
	server.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins: []string{"*"}, // Allow all origins for development
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))
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


	return &HTTP{server: server.Server}
}

func (h *HTTP) Start() {
	zlog := logger.NewZap()
	zlog.Info("Starting HTTP server on port 8080")
	if err := h.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		zlog.Sugar().Fatalf("Failed to start HTTP server: %v", err)
	}
}

func (h *HTTP) Stop(ctx context.Context) {
	zlog := logger.NewZap()
	zlog.Info("Shutting down HTTP server...")
	if err := h.server.Shutdown(ctx); err != nil {
		zlog.Sugar().Errorf("HTTP server forced to shutdown: %v", err)
	}
	zlog.Info("HTTP server shutdown complete.")
}
