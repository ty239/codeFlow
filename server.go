package main

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type apiServer struct {
	db        *pgxpool.Pool
	jwtSecret []byte
}

func registerRoutes(mux *http.ServeMux, s *apiServer) {
	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/health", healthHandler)

	mux.HandleFunc("POST /signup", s.signupHandler)
	mux.HandleFunc("POST /login", s.loginHandler)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"message": "hello backend"})
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
