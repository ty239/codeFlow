package main

import (
	"log"
	"os"
)

type Config struct {
	Port        string
	DatabaseURL string
	JWTSecret   string
}

func loadConfig() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/codeflow?sslmode=disable"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Println("WARNING: JWT_SECRET not set, using an insecure default (dev only)")
		jwtSecret = "dev-insecure-secret-change-me"
	}

	return Config{Port: port, DatabaseURL: dbURL, JWTSecret: jwtSecret}
}
