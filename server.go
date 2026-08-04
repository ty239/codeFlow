package main

import "net/http"

func registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/health", healthHandler)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"message": "hello backend"})
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
