package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/jakabrajadenta/go-explore-api/config"
	"github.com/jakabrajadenta/go-explore-api/internal/handler"
	"github.com/jakabrajadenta/go-explore-api/internal/repository"
	"github.com/jakabrajadenta/go-explore-api/internal/service"
	"github.com/jakabrajadenta/go-explore-api/middleware"
	"github.com/jakabrajadenta/go-explore-api/pkg/logger"
)

func main() {
	// ── Logging ───────────────────────────────────────────────
	logger.Init(os.Getenv("LOG_FORMAT") == "json")

	// ── Database ──────────────────────────────────────────────
	dbCfg := config.NewDBConfig()
	pool, err := config.NewPool(dbCfg)
	if err != nil {
		slog.Error("database connection failed", "error", err)
		os.Exit(1)
	}
	defer pool.Close()
	slog.Info("database connected", "host", dbCfg.Host, "db", dbCfg.DBName, "schema", dbCfg.Schema)

	// ── Dependency wiring ─────────────────────────────────────
	userRepo := repository.NewUserRepository(pool)
	userSvc := service.NewUserService(userRepo)

	// ── Router ────────────────────────────────────────────────
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux, userSvc)

	// ── Server ────────────────────────────────────────────────
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      middleware.Chain(mux, middleware.Logger, middleware.CORS),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	slog.Info("server starting", "addr", "http://localhost:"+port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}
