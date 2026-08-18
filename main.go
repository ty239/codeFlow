package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := loadConfig()

	ctx, cancelConnect := context.WithTimeout(context.Background(), 10*time.Second)
	db, err := connectDB(ctx, cfg.DatabaseURL)
	cancelConnect()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	api := &apiServer{db: db, jwtSecret: []byte(cfg.JWTSecret)}

	mux := http.NewServeMux()
	registerRoutes(mux, api)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      withMiddleware(mux),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("server listening on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("graceful shutdown failed: %v", err)
	}
	log.Println("server stopped")
}
