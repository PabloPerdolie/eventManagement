package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port                 string
	DatabaseURL          string
	RedisURL             string
	JWTSecretKey         string
	JWTAccessExpiration  time.Duration
	JWTRefreshExpiration time.Duration
	AllowedOrigin        string
	CoreServiceURL       string
	NotificationServiceURL string
	CommunicationServiceURL string
	PasswordResetExpiration time.Duration
}

func New() (*Config, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is required")
	}

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		return nil, fmt.Errorf("REDIS_URL environment variable is required")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET environment variable is required")
	}

	jwtAccessExpStr := os.Getenv("JWT_ACCESS_EXPIRATION")
	jwtAccessExp := 15 * time.Minute // Default: 15 minutes
	if jwtAccessExpStr != "" {
		exp, err := strconv.Atoi(jwtAccessExpStr)
		if err == nil {
			jwtAccessExp = time.Duration(exp) * time.Minute
		}
	}

	jwtRefreshExpStr := os.Getenv("JWT_REFRESH_EXPIRATION")
	jwtRefreshExp := 7 * 24 * time.Hour // Default: 7 days
	if jwtRefreshExpStr != "" {
		exp, err := strconv.Atoi(jwtRefreshExpStr)
		if err == nil {
			jwtRefreshExp = time.Duration(exp) * time.Hour
		}
	}

	pwResetExpStr := os.Getenv("PASSWORD_RESET_EXPIRATION")
	pwResetExp := 24 * time.Hour // Default: 24 hours
	if pwResetExpStr != "" {
		exp, err := strconv.Atoi(pwResetExpStr)
		if err == nil {
			pwResetExp = time.Duration(exp) * time.Hour
		}
	}

	allowedOrigin := os.Getenv("ALLOWED_ORIGIN")
	if allowedOrigin == "" {
		allowedOrigin = "*"
	}

	coreServiceURL := os.Getenv("CORE_SERVICE_URL")
	if coreServiceURL == "" {
		coreServiceURL = "http://localhost:8080"
	}

	notificationServiceURL := os.Getenv("NOTIFICATION_SERVICE_URL")
	if notificationServiceURL == "" {
		notificationServiceURL = "http://localhost:8082"
	}

	communicationServiceURL := os.Getenv("COMMUNICATION_SERVICE_URL")
	if communicationServiceURL == "" {
		communicationServiceURL = "http://localhost:8083"
	}

	return &Config{
		Port:                    port,
		DatabaseURL:             dbURL,
		RedisURL:                redisURL,
		JWTSecretKey:            jwtSecret,
		JWTAccessExpiration:     jwtAccessExp,
		JWTRefreshExpiration:    jwtRefreshExp,
		AllowedOrigin:           allowedOrigin,
		CoreServiceURL:          coreServiceURL,
		NotificationServiceURL:  notificationServiceURL,
		CommunicationServiceURL: communicationServiceURL,
		PasswordResetExpiration: pwResetExp,
	}, nil
}
