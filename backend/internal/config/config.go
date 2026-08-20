// Package config centralizes reading of environment variables so the rest
// of the backend never hardcodes ports or file paths.
package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Config holds every environment-driven setting used by the backend.
type Config struct {
	WSPort       string
	DBPath       string
	TempImageDir string
	AIModelsPath string
}

// Load reads a .env file (if present) and returns the resolved Config.
// Values already present in the process environment always win over the
// .env file, and sane defaults are used when nothing is set at all.
func Load() Config {
	// The backend is usually started from backend/ (go run cmd/main.go),
	// while .env lives at the repository root, so try both locations.
	// godotenv.Load never overrides variables already set in the environment.
	_ = godotenv.Load(".env")
	_ = godotenv.Load("../.env")

	return Config{
		WSPort:       getEnv("BACKEND_WS_PORT", "8080"),
		DBPath:       getEnv("SQLITE_DB_PATH", "./foto-facil.db"),
		TempImageDir: getEnv("TEMP_IMAGE_DIR", "./tmp_processed_images"),
		AIModelsPath: getEnv("AI_MODELS_PATH", "./backend/models"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
