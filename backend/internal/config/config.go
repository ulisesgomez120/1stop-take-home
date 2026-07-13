// Package config loads application configuration from environment variables.
package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	OneStepGPSAPIKey string
	Port             string
	DBPath           string
}

// Load reads a .env file if present (ignored if missing) and returns the
// resolved Config from environment variables, applying defaults where
// appropriate.
func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		OneStepGPSAPIKey: os.Getenv("ONESTEPGPS_API_KEY"),
		Port:             getEnvOrDefault("PORT", "8080"),
		DBPath:           getEnvOrDefault("DB_PATH", "data.db"),
	}

	return cfg, nil
}

func getEnvOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
