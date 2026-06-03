package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/jakabrajadenta/go-explore-api/handler"
	"github.com/jakabrajadenta/go-explore-api/middleware"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

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
