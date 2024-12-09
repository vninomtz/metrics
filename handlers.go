package main

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"
)

const (
	REQUEST_ID_KEY int = 0
)

func mTracing(nextReqId func() string, next http.Handler) http.Handler {
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

func mLogging(logger *log.Logger, next http.Handler) http.Handler {
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

func NewServer(logger *log.Logger) http.Handler {
	mux := http.NewServeMux()

	addRoutes(mux)

	var handler http.Handler = mux

	nextRequestID := func() string {
		return strconv.FormatInt(time.Now().UnixNano(), 10)
	}
	handler = mLogging(logger, handler)
	handler = mTracing(nextRequestID, handler)

	return handler
}

var healthy int32

func handleHealthz(w http.ResponseWriter, r *http.Request) {
	if atomic.LoadInt32(&healthy) == 1 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.WriteHeader(http.StatusServiceUnavailable)
}

func addRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/healthz", handleHealthz)
}
