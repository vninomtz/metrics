package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync/atomic"
	"syscall"
	"time"
)

const (
	REQUEST_ID_KEY int = 0
)

func tracing(nextReqId func() string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqId := r.Header.Get("X-Request-Id")
			if reqId == "" {
				reqId = nextReqId()
			}
			ctx := context.WithValue(r.Context(), REQUEST_ID_KEY, reqId)

			w.Header().Set("X-Request-Id", reqId)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func logging(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				reqId, ok := r.Context().Value(REQUEST_ID_KEY).(string)
				if !ok {
					reqId = "unknown"
				}
				logger.Println(reqId, r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())
			}()
			next.ServeHTTP(w, r)
		})
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Ok")
}

var healthy int32

func healthcheck(w http.ResponseWriter, r *http.Request) {
	if atomic.LoadInt32(&healthy) == 1 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.WriteHeader(http.StatusServiceUnavailable)
}

func run() {
	logger := log.New(os.Stdout, "http: ", log.LstdFlags)
	http.HandleFunc("/", handler)
	http.HandleFunc("/healthz", healthcheck)

	host := "127.0.0.1"
	port := "8081"

	nextRequestID := func() string {
		return strconv.FormatInt(time.Now().UnixNano(), 10)
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", host, port),
		Handler:      tracing(nextRequestID)(logging(logger)(http.DefaultServeMux)),
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
