package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type signupRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type userResponse struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}

func toUserResponse(u *User) userResponse {
	return userResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Name:      u.Name,
		CreatedAt: u.CreatedAt.Format(time.RFC3339),
	}
}

func (s *apiServer) signupHandler(w http.ResponseWriter, r *http.Request) {
	var req signupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(req.Email)
	req.Name = strings.TrimSpace(req.Name)

	if req.Username == "" || req.Email == "" || req.Name == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "username, email, name, and password are required")
		return
	}
	if len(req.Password) < 8 {
		writeError(w, http.StatusBadRequest, "password must be at least 8 characters")
		return
	}

	passwordHash, err := hashPassword(req.Password)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to process password")
		return
	}

	user, err := createUser(r.Context(), s.db, req.Username, req.Email, req.Name, passwordHash)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			writeError(w, http.StatusConflict, "username or email already in use")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to create user")
		return
	}

	writeJSON(w, http.StatusCreated, toUserResponse(user))
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string       `json:"token"`
	User  userResponse `json:"user"`
}

func (s *apiServer) loginHandler(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := getUserByUsername(r.Context(), s.db, strings.TrimSpace(req.Username))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeError(w, http.StatusUnauthorized, "invalid username or password")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to log in")
		return
	}

	if !checkPassword(user.PasswordHash, req.Password) {
		writeError(w, http.StatusUnauthorized, "invalid username or password")
		return
	}

	token, err := generateToken(s.jwtSecret, user.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to generate token")
		return
	}

	writeJSON(w, http.StatusOK, loginResponse{Token: token, User: toUserResponse(user)})
}
