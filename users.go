package main

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	ID           string
	Username     string
	Email        string
	Name         string
	PasswordHash string
	CreatedAt    time.Time
}

func createUser(ctx context.Context, db *pgxpool.Pool, username, email, name, passwordHash string) (*User, error) {
	const query = `
		INSERT INTO users (username, email, name, password_hash)
		VALUES ($1, $2, $3, $4)
		RETURNING id, username, email, name, password_hash, created_at`

	var u User
	err := db.QueryRow(ctx, query, username, email, name, passwordHash).Scan(
		&u.ID, &u.Username, &u.Email, &u.Name, &u.PasswordHash, &u.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func getUserByUsername(ctx context.Context, db *pgxpool.Pool, username string) (*User, error) {
	const query = `
		SELECT id, username, email, name, password_hash, created_at
		FROM users
		WHERE username = $1`

	var u User
	err := db.QueryRow(ctx, query, username).Scan(
		&u.ID, &u.Username, &u.Email, &u.Name, &u.PasswordHash, &u.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
