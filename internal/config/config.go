package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv               string
	Port                 string
	DatabaseURL          string
	CORSAllowedOrigins   []string
	CORSAllowCredentials bool
	JWTSecret            string
	JWTIssuer            string
	JWTAccessTokenTTL    time.Duration
}

func Load() (Config, error) {
	_ = godotenv.Load()

	port := getEnv("PORT", "8080")
	if value, err := strconv.Atoi(port); err != nil || value < 1 || value > 65535 {
		return Config{}, fmt.Errorf("PORT must be a number between 1 and 65535")
	}

	allowCredentials, err := getBoolEnv("CORS_ALLOW_CREDENTIALS", true)
	if err != nil {
		return Config{}, err
	}

	ttl, err := getDurationEnv("JWT_ACCESS_TOKEN_TTL", 15*time.Minute)
	if err != nil {
		return Config{}, err
	}

	return Config{
		AppEnv:               getEnv("APP_ENV", "development"),
		Port:                 port,
		DatabaseURL:          os.Getenv("DATABASE_URL"),
		CORSAllowedOrigins:   getCSVEnv("CORS_ALLOWED_ORIGINS", []string{"http://localhost:3000", "http://localhost:5173"}),
		CORSAllowCredentials: allowCredentials,
		JWTSecret:            os.Getenv("JWT_SECRET"),
		JWTIssuer:            getEnv("JWT_ISSUER", "go-boilerplate"),
		JWTAccessTokenTTL:    ttl,
	}, nil
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}

func getCSVEnv(key string, fallback []string) []string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parts := strings.Split(value, ",")
	items := make([]string, 0, len(parts))
	for _, part := range parts {
		item := strings.TrimSpace(part)
		if item != "" {
			items = append(items, item)
		}
	}
	if len(items) == 0 {
		return fallback
	}

	return items
}

func getBoolEnv(key string, fallback bool) (bool, error) {
	value := os.Getenv(key)
	if value == "" {
		return fallback, nil
	}

	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return false, fmt.Errorf("%s must be a boolean", key)
	}

	return parsed, nil
}

func getDurationEnv(key string, fallback time.Duration) (time.Duration, error) {
	value := os.Getenv(key)
	if value == "" {
		return fallback, nil
	}

	parsed, err := time.ParseDuration(value)
	if err != nil || parsed <= 0 {
		return 0, fmt.Errorf("%s must be a positive duration", key)
	}

	return parsed, nil
}
