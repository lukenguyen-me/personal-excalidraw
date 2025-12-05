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
	"github.com/personal-excalidraw/backend/internal/infrastructure/config"
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

	// 3. Initialize HTTP handlers
	healthHandler := handler.NewHealthHandler()

	// 4. Setup router
	router := httpAdapter.NewRouter(cfg, healthHandler, appLogger)

	// 5. Create HTTP server
	serverAddr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	server := &http.Server{
		Addr:         serverAddr,
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}

	// 6. Start server in goroutine
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

	// 7. Graceful shutdown
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
