package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpAdapter "github.com/personal-excalidraw/backend/internal/adapter/http"
	"github.com/personal-excalidraw/backend/internal/adapter/http/handler"
	"github.com/personal-excalidraw/backend/internal/adapter/repository/postgres"
	drawingapp "github.com/personal-excalidraw/backend/internal/application/drawing"
	"github.com/personal-excalidraw/backend/internal/infrastructure/config"
	"github.com/personal-excalidraw/backend/internal/infrastructure/database"
	"github.com/personal-excalidraw/backend/internal/infrastructure/logger"
)

func main() {
	// 1. Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. Initialize logger
	appLogger := logger.New(&cfg.Logger)
	appLogger.Info("Starting Personal Excalidraw backend server")

	// 3. Initialize database connection
	db, err := database.NewPostgresDB(&cfg.Database, appLogger)
	if err != nil {
		appLogger.Error("Failed to connect to database", "error", err)
		log.Fatalf("Database connection failed: %v", err)
	}
	defer db.Close()

	// 4. Initialize repositories
	drawingRepo := postgres.NewDrawingRepository(db.Pool)

	// 5. Initialize application services
	drawingService := drawingapp.NewService(drawingRepo, appLogger)

	// 6. Initialize HTTP handlers
	healthHandler := handler.NewHealthHandler()
	drawingHandler := handler.NewDrawingHandler(drawingService, appLogger)

	// 7. Setup router
	router := httpAdapter.NewRouter(cfg, healthHandler, drawingHandler, appLogger)

	// 8. Create HTTP server
	serverAddr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	server := &http.Server{
		Addr:         serverAddr,
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}

	// 9. Start server in goroutine
	go func() {
		appLogger.Info("Server starting",
			"address", serverAddr,
			"read_timeout", cfg.Server.ReadTimeout,
			"write_timeout", cfg.Server.WriteTimeout,
		)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Error("Server error", "error", err)
			os.Exit(1)
		}
	}()

	// 10. Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		appLogger.Error("Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	appLogger.Info("Server exited gracefully")
}
