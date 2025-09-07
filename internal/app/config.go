// Package app contains application configuration, middleware, and server wiring.
package app

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds application configuration loaded from environment variables.
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

// ServerConfig describes HTTP server settings and related middleware configuration.
type ServerConfig struct {
	Port             string
	Host             string
	PublicScheme     string
	EnableSwagger    bool
	EnableLogger     bool
	EnableRateLimit  bool
	AllowedOrigins   []string
	AllowCredentials bool
	GinMode          string
	TrustedProxies   []string
}

// DatabaseConfig describes database connection settings.
type DatabaseConfig struct {
	URL string
}

// Load reads configuration from environment variables and optional .env file.
func Load() (*Config, error) {
	// Load .env file (non-fatal if missing)
	if err := godotenv.Load(); err != nil {
		log.Printf(".env not loaded: %v", err)
	}

	return &Config{
		Server: ServerConfig{
			Port:             getEnv("PORT", "8080"),
			Host:             getEnv("HOST", "localhost"),
			PublicScheme:     getEnv("PUBLIC_SCHEME", "http"),
			EnableSwagger:    getEnvBool("ENABLE_SWAGGER", false),
			EnableLogger:     getEnvBool("ENABLE_LOGGER", true),
			EnableRateLimit:  getEnvBool("ENABLE_RATE_LIMIT", false),
			AllowedOrigins:   splitAndTrim(getEnv("ALLOWED_ORIGINS", "")),
			AllowCredentials: getEnvBool("ALLOW_CREDENTIALS", true),
			GinMode:          getEnv("GIN_MODE", "release"),
			TrustedProxies:   splitAndTrim(getEnv("TRUSTED_PROXIES", "")),
		},
		Database: DatabaseConfig{
			URL: getEnv("DATABASE_URL", "data/app.db"),
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return defaultValue
	}

	switch v {
	case "1", "true", "TRUE", "True", "yes", "YES", "on", "ON":
		return true
	case "0", "false", "FALSE", "False", "no", "NO", "off", "OFF":
		return false
	default:
		return defaultValue
	}
}

func splitAndTrim(s string) []string {
	if s == "" {
		return nil
	}

	parts := strings.Split(s, ",")

	out := make([]string, 0, len(parts))
	for _, p := range parts {
		pp := strings.TrimSpace(p)
		if pp != "" {
			out = append(out, pp)
		}
	}

	return out
}
