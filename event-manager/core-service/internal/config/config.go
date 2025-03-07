package config

import (
	"fmt"
	"os"
)

// Config contains all the configuration for the application
type Config struct {
	Port          string
	DatabaseURL   string
	JWTSecretKey  string
	AllowedOrigin string
}

// New returns a new config loaded from environment variables
func New() (*Config, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is required")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET environment variable is required")
	}

	allowedOrigin := os.Getenv("ALLOWED_ORIGIN")
	if allowedOrigin == "" {
		allowedOrigin = "*"
	}

	return &Config{
		Port:          port,
		DatabaseURL:   dbURL,
		JWTSecretKey:  jwtSecret,
		AllowedOrigin: allowedOrigin,
	}, nil
}
