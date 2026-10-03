package api

import (
	"log"
	"net/http"
	"time"
)

func WithMiddleware(h http.Handler) http.Handler {
	return recoverMiddleware(loggingMiddleware(corsMiddleware(h)))
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}

func recoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tw := &trackingWriter{ResponseWriter: w}
		defer func() {
			err := recover()
			if err == nil {
				return
			}
			// ErrAbortHandler is a deliberate abort; let net/http handle it.
			if err == http.ErrAbortHandler {
				panic(err)
			}
			log.Printf("panic recovered: %v", err)
			// If the handler already started its response, writing an error now
			// would corrupt it, so only send one if nothing was written yet.
			if !tw.wroteHeader {
				writeError(w, http.StatusInternalServerError, "internal server error")
			}
		}()
		next.ServeHTTP(tw, r)
	})
}

// trackingWriter records whether a response has started being written.
type trackingWriter struct {
	http.ResponseWriter
	wroteHeader bool
}

func (t *trackingWriter) WriteHeader(status int) {
	t.wroteHeader = true
	t.ResponseWriter.WriteHeader(status)
}

func (t *trackingWriter) Write(b []byte) (int, error) {
	t.wroteHeader = true
	return t.ResponseWriter.Write(b)
}

// Unwrap lets http.ResponseController reach the underlying writer.
func (t *trackingWriter) Unwrap() http.ResponseWriter {
	return t.ResponseWriter
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
