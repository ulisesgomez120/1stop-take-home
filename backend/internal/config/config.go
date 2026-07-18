// Package config loads application configuration from environment variables.
package config

import (
	"errors"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// Non-secret settings with safe defaults. These are plain constants rather
// than env vars: .env is reserved for secrets and values that genuinely
// differ per environment.
const (
	oneStepGPSBaseURL = "https://track.onestepgps.com/v3/api/public/device"
	port              = "8080"
	dbPath            = "data.db"
	uploadsDir        = "uploads"
	pollInterval      = 10 * time.Second
	allowedOrigin     = "http://localhost:5173"
)

type Config struct {
	OneStepGPSAPIKey  string
	OneStepGPSBaseURL string
	Port              string
	DBPath            string
	UploadsDir        string
	PollInterval      time.Duration
	AllowedOrigin     string
}

// Load reads a .env file if present (ignored if missing) and returns the
// resolved Config. ONESTEPGPS_API_KEY is the only value sourced from the
// environment; everything else is a fixed default.
func Load() (*Config, error) {
	_ = godotenv.Load()

	apiKey := os.Getenv("ONESTEPGPS_API_KEY")
	if apiKey == "" {
		return nil, errors.New("ONESTEPGPS_API_KEY is required")
	}

	return &Config{
		OneStepGPSAPIKey:  apiKey,
		OneStepGPSBaseURL: oneStepGPSBaseURL,
		Port:              port,
		DBPath:            dbPath,
		UploadsDir:        uploadsDir,
		PollInterval:      pollInterval,
		AllowedOrigin:     allowedOrigin,
	}, nil
}
