package api

import (
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	db            *pgxpool.Pool
	jwtSecret     []byte
	loginLimiter  *rateLimiter
	signupLimiter *rateLimiter
}

func NewServer(db *pgxpool.Pool, jwtSecret []byte) *Server {
	return &Server{
		db:        db,
		jwtSecret: jwtSecret,
		// Limits are per client IP. Login allows a few attempts a minute to slow
		// password guessing; signup is limited per hour to slow mass account creation.
		loginLimiter:  newRateLimiter(5, time.Minute),
		signupLimiter: newRateLimiter(10, time.Hour),
	}
}

func (s *Server) RegisterRoutes(mux *http.ServeMux) {
	// {$} matches "/" exactly; a bare "/" pattern would catch every unknown path.
	mux.HandleFunc("GET /{$}", homeHandler)
	mux.HandleFunc("GET /health", healthHandler)

	mux.HandleFunc("POST /signup", s.signupLimiter.wrap(s.signupHandler))
	mux.HandleFunc("POST /login", s.loginLimiter.wrap(s.loginHandler))
}
