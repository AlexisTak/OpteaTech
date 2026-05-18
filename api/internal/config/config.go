package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port              string
	DatabaseURL       string
	BaseURL           string
	PortalBaseURL     string
	AllowedOrigins    []string
	ResendAPIKey      string
	EmailFrom         string
	AdminEmail        string
	JWTSecret         string
	JWTExpiresSeconds int64
}

func Load() Config {
	return Config{
		Port:              getEnv("PORT", "3001"),
		DatabaseURL:       getEnv("DATABASE_URL", ""),
		BaseURL:           getEnv("BASE_URL", "http://localhost:3000"),
		PortalBaseURL:     getEnv("CLIENT_PORTAL_BASE_URL", getEnv("BASE_URL", "http://localhost:3000")),
		AllowedOrigins:    normalizeOrigins(getEnv("ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:5173,http://localhost:4173")),
		ResendAPIKey:      getEnv("RESEND_API_KEY", ""),
		EmailFrom:         getEnv("FROM_EMAIL", ""),
		AdminEmail:        getEnv("ADMIN_EMAIL", ""),
		JWTSecret:         getEnv("JWT_SECRET", "dev_secret_change_me"),
		JWTExpiresSeconds: getInt64Env("JWT_EXPIRES_SECONDS", 900),
	}
}

func getEnv(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getInt64Env(key string, fallback int64) int64 {
	value := getEnv(key, "")
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return fallback
	}
	return parsed
}

func normalizeOrigins(raw string) []string {
	parts := strings.Split(raw, ",")
	origins := make([]string, 0, len(parts))
	for _, part := range parts {
		origin := strings.TrimSpace(part)
		if origin == "" {
			continue
		}
		origins = append(origins, origin)
	}
	return origins
}
