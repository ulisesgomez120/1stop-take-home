// Package config loads application configuration from environment variables.
package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

const defaultOneStepGPSBaseURL = "https://track.onestepgps.com/v3/api/public/device"

type Config struct {
	OneStepGPSAPIKey  string
	OneStepGPSBaseURL string
	Port              string
	DBPath            string
	PollInterval      time.Duration
}

// Load reads a .env file if present (ignored if missing) and returns the
// resolved Config from environment variables, applying defaults where
// appropriate.
func Load() (*Config, error) {
	_ = godotenv.Load()

	pollSecs, err := strconv.Atoi(getEnvOrDefault("POLL_INTERVAL_SECONDS", "10"))
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		OneStepGPSAPIKey:  os.Getenv("ONESTEPGPS_API_KEY"),
		OneStepGPSBaseURL: getEnvOrDefault("ONESTEPGPS_BASE_URL", defaultOneStepGPSBaseURL),
		Port:              getEnvOrDefault("PORT", "8080"),
		DBPath:            getEnvOrDefault("DB_PATH", "data.db"),
		PollInterval:      time.Duration(pollSecs) * time.Second,
	}

	return cfg, nil
}

func getEnvOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
