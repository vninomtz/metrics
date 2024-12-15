package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func addRoutes(mux *http.ServeMux, logger *log.Logger, repo *Repository) {
	mux.Handle("/api/views", handleViews(logger, repo))
	mux.HandleFunc("/healthz", handleHealthz)
}

func handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}
func handleViews(logger *log.Logger, repo *Repository) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
			return
		}

		var data View
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			logger.Printf("views: Error decoding body  %v\n", err)
			http.Error(w, "Error to process metrics", http.StatusBadRequest)
			return
		}

		data.Created = time.Now()
		err := repo.SaveView(data)
		if err != nil {
			logger.Printf("Error:  %v\n", err)
			http.Error(w, "Error to process metrics", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
	})
}
