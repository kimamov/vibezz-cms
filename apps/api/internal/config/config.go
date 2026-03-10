package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port             string
	GinMode          string
	DatabaseURL      string
	JWTSecret        string
	SessionSecret    string
	UploadDir        string
	MaxUploadSizeMB  int64
	AdminURL         string
	WebURL           string
}

func Load() *Config {
	return &Config{
		Port:            getEnv("PORT", "8080"),
		GinMode:         getEnv("GIN_MODE", "debug"),
		DatabaseURL:     getEnv("DATABASE_URL", "postgres://vibezz:vibezz@localhost:5432/vibezz_cms?sslmode=disable"),
		JWTSecret:       getEnv("JWT_SECRET", "change-me-in-production"),
		SessionSecret:   getEnv("SESSION_SECRET", "change-me-in-production"),
		UploadDir:       getEnv("UPLOAD_DIR", "./uploads"),
		MaxUploadSizeMB: getEnvInt("MAX_UPLOAD_SIZE_MB", 50),
		AdminURL:        getEnv("ADMIN_URL", "http://localhost:3001"),
		WebURL:          getEnv("WEB_URL", "http://localhost:3000"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int64) int64 {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.ParseInt(v, 10, 64); err == nil {
			return i
		}
	}
	return fallback
}
