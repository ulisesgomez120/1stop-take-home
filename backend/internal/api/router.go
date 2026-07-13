// Package api wires HTTP handlers and routing on top of the domain packages.
package api

import (
	"net/http"
	"time"

	"uli1step.com/internal/devices"
	"uli1step.com/internal/preferences"
)

// Config holds the dependencies and settings needed to build the router.
type Config struct {
	Cache            *devices.Cache
	PreferencesStore *preferences.Store
	StreamInterval   time.Duration
	AllowedOrigin    string
}

// NewRouter builds the HTTP handler for the API surface.
func NewRouter(cfg Config) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("GET /api/devices", handleGetDevices(cfg.Cache))
	mux.HandleFunc("GET /api/devices/stream", handleDeviceStream(cfg.Cache, cfg.StreamInterval))
	mux.HandleFunc("GET /api/preferences", handleGetPreferences(cfg.PreferencesStore))
	mux.HandleFunc("PUT /api/preferences", handlePutPreferences(cfg.PreferencesStore))

	return withCORS(cfg.AllowedOrigin, mux)
}
