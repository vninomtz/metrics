package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
)

func run() {
	logger := log.New(os.Stdout, "http: ", log.LstdFlags)

	host := "127.0.0.1"
	port := "8081"

	hdlServer := NewServer(logger)

	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", host, port),
		Handler:      hdlServer,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	atomic.StoreInt32(&healthy, 1)

	go func() {
		<-quit
		logger.Println("Server is shuttingn down")
		atomic.StoreInt32(&healthy, 0)
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		srv.SetKeepAlivesEnabled(false)
		if err := srv.Shutdown(ctx); err != nil {
			logger.Fatalf("Could not gracefully shutdown the server: %v\n", err)
		}
		close(done)
	}()

	logger.Printf("Server is running at %s:%s\n", host, port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("Could not start the server: %v\n", err)
	}
	<-done
	logger.Println("Server stopped")
}

func main() {
	run()
}
