package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fikryfahrezy/let-it-go/config"
	healthHandler "github.com/fikryfahrezy/let-it-go/feature/health/handler"
	userHandler "github.com/fikryfahrezy/let-it-go/feature/user/handler"
	"github.com/fikryfahrezy/let-it-go/feature/user/repository"
	"github.com/fikryfahrezy/let-it-go/feature/user/service"
	"github.com/fikryfahrezy/let-it-go/pkg/database"
	"github.com/fikryfahrezy/let-it-go/pkg/logger"
	"github.com/fikryfahrezy/let-it-go/pkg/server"
)

func main() {
	cfg := config.Load()

	log := logger.NewLogger(cfg.Logger)

	log.Info("Starting application",
		slog.String("server_host", cfg.Server.Host),
		slog.Int("server_port", cfg.Server.Port),
	)

	db, err := database.NewDB(cfg.Database)
	if err != nil {
		log.Error("Failed to connect to database",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}

	// Create server configuration
	serverConfig := server.Config{
		Host: cfg.Server.Host,
		Port: cfg.Server.Port,
	}

	// Initialize feature dependencies
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandlerInstance := userHandler.NewUserHandler(userService)
	
	// Initialize health handler
	healthHandlerInstance := healthHandler.NewHealthHandler(db)

	// Create and initialize server
	srv, err := server.New(serverConfig)
	if err != nil {
		log.Error("Failed to create server",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}

	if err := srv.Initialize([]server.RouteHandler{healthHandlerInstance, userHandlerInstance}); err != nil {
		log.Error("Failed to initialize server",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}

	// Start server in goroutine
	go func() {
		if err := srv.Start(); err != nil {
			log.Error("Server error",
				slog.String("error", err.Error()),
			)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down gracefully...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Stop server
	if err := srv.Stop(ctx); err != nil {
		log.Error("Failed to shutdown server gracefully",
			slog.String("error", err.Error()),
		)
	}

	// Close database connection
	if err := db.Close(); err != nil {
		log.Error("Failed to close database connection",
			slog.String("error", err.Error()),
		)
	}

	log.Info("Application shutdown complete")
}