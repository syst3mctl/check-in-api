package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/syst3mctl/check-in-api/internal/adapter/storage/postgres"
	"github.com/syst3mctl/check-in-api/internal/api/handler"
	"github.com/syst3mctl/check-in-api/internal/api/middleware"
	"github.com/syst3mctl/check-in-api/internal/api/router"
	"github.com/syst3mctl/check-in-api/internal/config"
	"github.com/syst3mctl/check-in-api/internal/core/service"
	"github.com/syst3mctl/check-in-api/internal/pkg/logger"
)

// @title Check-In Service API
// @version 1.0
// @description REST API for Check-In/Out Service
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	logger.Init()

	cfg, err := config.Load()
	if err != nil {
		logger.Error("Failed to load config", "error", err)
		return
	}

	db, err := postgres.NewDB(cfg.DatabaseURL)
	if err != nil {
		logger.Error("Failed to connect to database", "error", err)
		return
	}
	defer db.Close()

	// Repositories
	userRepo := postgres.NewUserRepository(db)
	orgRepo := postgres.NewOrgRepository(db)
	attRepo := postgres.NewAttendanceRepository(db)
	reportRepo := postgres.NewReportRepository(db)

	// Services
	authService := service.NewAuthService(userRepo, cfg)
	orgService := service.NewOrgService(orgRepo, userRepo, db)
	attService := service.NewAttendanceService(attRepo)
	reportService := service.NewReportService(reportRepo)

	// Handlers
	authHandler := handler.NewAuthHandler(authService)
	orgHandler := handler.NewOrgHandler(orgService)
	attHandler := handler.NewAttendanceHandler(attService)
	reportHandler := handler.NewReportHandler(reportService)

	// Middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg)

	// Router
	r := router.New(authHandler, orgHandler, attHandler, reportHandler, authMiddleware)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		logger.Info("Starting server", "port", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server failed", "error", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", "error", err)
	}

	logger.Info("Server exiting")
}
