package main

import (
	"boredgamz/server"
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	production := os.Getenv("APP_ENV") == "production"
	if !production {
		// Local convenience only; production config comes from the environment.
		_ = godotenv.Load()
	}
	slog.SetDefault(slog.New(newLogHandler(production)))

	if err := checkConfig(production); err != nil {
		slog.Error("invalid configuration", "err", err)
		os.Exit(1)
	}
	s, err := server.CreateServer()
	if err != nil {
		slog.Error("failed to initialize server", "err", err)
		os.Exit(1)
	}
	defer s.DB.Stop()

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	httpServer := &http.Server{
		Addr:              ":" + port,
		Handler:           s.Router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    16 << 10,
	}

	errs := make(chan error, 1)
	go func() {
		slog.Info("boredgamz listening", "addr", httpServer.Addr, "production", production)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errs <- err
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	select {
	case sig := <-stop:
		slog.Info("shutting down", "signal", sig.String())
	case err := <-errs:
		slog.Error("server stopped unexpectedly", "err", err)
		os.Exit(1)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// Shutdown drains HTTP requests. Hijacked WebSockets are not tracked by
	// net/http; they close when the process exits and clients auto-reconnect
	// to the replacement instance (live match state is lost, see docs/DEPLOY.md).
	if err := httpServer.Shutdown(ctx); err != nil {
		slog.Warn("graceful shutdown timed out", "err", err)
	}
}

func newLogHandler(production bool) slog.Handler {
	level := slog.LevelInfo
	if strings.EqualFold(os.Getenv("LOG_LEVEL"), "debug") {
		level = slog.LevelDebug
	}
	options := &slog.HandlerOptions{Level: level}
	if production || os.Getenv("LOG_FORMAT") == "json" {
		return slog.NewJSONHandler(os.Stdout, options)
	}
	return slog.NewTextHandler(os.Stdout, options)
}

func checkConfig(production bool) error {
	if !production {
		return nil
	}
	if os.Getenv("ALLOWED_ORIGINS") == "" {
		return errors.New("ALLOWED_ORIGINS is required in production")
	}
	if os.Getenv("DATABASE_URL") == "" {
		return errors.New("DATABASE_URL is required in production")
	}
	if len(os.Getenv("JWT_SECRET_KEY")) < 32 {
		return errors.New("JWT_SECRET_KEY must contain at least 32 bytes in production")
	}
	return nil
}
