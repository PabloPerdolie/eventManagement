package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config содержит все настройки приложения
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	CORS     CORSConfig
}

// ServerConfig содержит настройки HTTP-сервера
type ServerConfig struct {
	Port         string
	ReadTimeout  int
	WriteTimeout int
}

// DatabaseConfig содержит настройки подключения к БД
type DatabaseConfig struct {
	Host            string
	Port            string
	Username        string
	Password        string
	DBName          string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime int
}

// JWTConfig содержит настройки для JWT
type JWTConfig struct {
	SecretKey string
	ExpiresIn int // в часах
}

// CORSConfig содержит настройки CORS
type CORSConfig struct {
	AllowedOrigins string
}

// LoadConfig загружает конфигурацию из переменных окружения
func LoadConfig() (*Config, error) {
	// Попытка загрузить .env файл, если он существует
	_ = godotenv.Load()

	// Настройки сервера
	port := getEnv("SERVER_PORT", "8080")
	readTimeout, _ := strconv.Atoi(getEnv("SERVER_READ_TIMEOUT", "10"))
	writeTimeout, _ := strconv.Atoi(getEnv("SERVER_WRITE_TIMEOUT", "10"))

	// Настройки базы данных
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "event_manager")
	dbSSLMode := getEnv("DB_SSLMODE", "disable")
	dbMaxOpenConns, _ := strconv.Atoi(getEnv("DB_MAX_OPEN_CONNS", "20"))
	dbMaxIdleConns, _ := strconv.Atoi(getEnv("DB_MAX_IdLE_CONNS", "5"))
	dbConnMaxLifetime, _ := strconv.Atoi(getEnv("DB_CONN_MAX_LIFETIME", "300"))

	// Настройки JWT
	jwtSecret := getEnv("JWT_SECRET", "")
	if jwtSecret == "" {
		// Для разработки можно использовать значение по умолчанию,
		// но в production это должен быть надежный секретный ключ
		jwtSecret = "dev_secret_key"
	}
	jwtExpiresIn, _ := strconv.Atoi(getEnv("JWT_EXPIRES_IN", "24"))

	// Настройки CORS
	corsAllowedOrigins := getEnv("CORS_ALLOWED_ORIGINS", "*")

	return &Config{
		Server: ServerConfig{
			Port:         port,
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
		},
		Database: DatabaseConfig{
			Host:            dbHost,
			Port:            dbPort,
			Username:        dbUser,
			Password:        dbPassword,
			DBName:          dbName,
			SSLMode:         dbSSLMode,
			MaxOpenConns:    dbMaxOpenConns,
			MaxIdleConns:    dbMaxIdleConns,
			ConnMaxLifetime: dbConnMaxLifetime,
		},
		JWT: JWTConfig{
			SecretKey: jwtSecret,
			ExpiresIn: jwtExpiresIn,
		},
		CORS: CORSConfig{
			AllowedOrigins: corsAllowedOrigins,
		},
	}, nil
}

// getEnv получает значение переменной окружения или возвращает значение по умолчанию
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
