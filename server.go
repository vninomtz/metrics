package main

import (
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	logger := log.New(os.Stdout, "http: ", log.LstdFlags)

	cfg := NewConfig()
	err := cfg.Load()
	if err != nil {
		logger.Fatalln(err)
	}

	repo, err := NewRepository(cfg.PathDB)
	if err != nil {
		logger.Fatal(err)
	}
	hdlServer := NewServer(logger, repo)

	srv := &http.Server{
		Addr:         cfg.Addr(),
		Handler:      hdlServer,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	logger.Printf("Server is running at %s\n", cfg.Addr())
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("Could not start the server: %v\n", err)
	}
}
