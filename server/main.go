package main

import (
	"boredgamz/server"
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	if os.Getenv("APP_ENV") == "production" && os.Getenv("ALLOWED_ORIGINS") == "" {
		log.Fatal("ALLOWED_ORIGINS is required in production")
	}
	if os.Getenv("APP_ENV") == "production" && os.Getenv("DATABASE_URL") == "" {
		log.Fatal("DATABASE_URL is required in production")
	}
	if os.Getenv("APP_ENV") == "production" && len(os.Getenv("JWT_SECRET_KEY")) < 32 {
		log.Fatal("JWT_SECRET_KEY must contain at least 32 bytes in production")
	}
	s, err := server.CreateServer()
	if err != nil {
		log.Fatalf("failed to initialize server: %v", err)
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
	}

	go func() {
		log.Printf("boredgamz listening on %s", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server stopped unexpectedly: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Printf("graceful shutdown timed out: %v", err)
	}
}
