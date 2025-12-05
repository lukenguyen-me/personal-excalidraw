package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	Server ServerConfig
	CORS   CORSConfig
	Logger LoggerConfig
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Host         string
	Port         string
	ReadTimeout  int
	WriteTimeout int
}

// CORSConfig holds CORS-related configuration
type CORSConfig struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
}

// LoggerConfig holds logger-related configuration
type LoggerConfig struct {
	Level  string
	Format string // "json" or "text"
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if exists (ignore error if not found)
	_ = godotenv.Load()

	cfg := &Config{
		Server: ServerConfig{
			Host:         getEnv("SERVER_HOST", "0.0.0.0"),
			Port:         getEnv("SERVER_PORT", "8080"),
			ReadTimeout:  getEnvInt("SERVER_READ_TIMEOUT", 15),
			WriteTimeout: getEnvInt("SERVER_WRITE_TIMEOUT", 15),
		},
		CORS: CORSConfig{
			AllowedOrigins: getEnvList("CORS_ALLOWED_ORIGINS", []string{"http://localhost:5173"}),
			AllowedMethods: getEnvList("CORS_ALLOWED_METHODS", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
			AllowedHeaders: getEnvList("CORS_ALLOWED_HEADERS", []string{"Content-Type", "Authorization"}),
		},
		Logger: LoggerConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "text"),
		},
	}

	return cfg, nil
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt retrieves an integer environment variable or returns a default value
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvList retrieves a comma-separated list from environment variable or returns default
func getEnvList(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		// Split by comma and trim spaces
		parts := strings.Split(value, ",")
		result := make([]string, 0, len(parts))
		for _, part := range parts {
			if trimmed := strings.TrimSpace(part); trimmed != "" {
				result = append(result, trimmed)
			}
		}
		if len(result) > 0 {
			return result
		}
	}
	return defaultValue
}
