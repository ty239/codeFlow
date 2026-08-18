package api

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	db        *pgxpool.Pool
	jwtSecret []byte
}

func NewServer(db *pgxpool.Pool, jwtSecret []byte) *Server {
	return &Server{db: db, jwtSecret: jwtSecret}
}

func (s *Server) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/health", healthHandler)

	mux.HandleFunc("POST /signup", s.signupHandler)
	mux.HandleFunc("POST /login", s.loginHandler)
}
