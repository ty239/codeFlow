package main

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func registerRoutes(mux *http.ServeMux, db *pgxpool.Pool) {
	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/health", healthHandler)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"message": "hello backend"})
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
