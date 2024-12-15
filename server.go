package main

import (
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	cfg := NewConfig()
	err := cfg.Load()
	if err != nil {
		log.Fatalf("Error to load config: %v", err)
		os.Exit(1)
	}
	_logger := NewLogger(*cfg)

	repo, err := NewRepository(cfg.PathDB)
	if err != nil {
		_logger.Error("Error to load DB: %v", err)
		os.Exit(1)
	}
	hdlServer := NewServer(_logger, repo)

	srv := &http.Server{
		Addr:         cfg.Addr(),
		Handler:      hdlServer,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	_logger.Info("Server is running at %s\n", cfg.Addr())
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		_logger.Error("Could not start the server: %v\n", err)
	}
}
