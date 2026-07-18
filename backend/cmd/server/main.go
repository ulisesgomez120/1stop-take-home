package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"uli1step.com/internal/api"
	"uli1step.com/internal/config"
	"uli1step.com/internal/devices"
	"uli1step.com/internal/icons"
	"uli1step.com/internal/onestepgps"
	"uli1step.com/internal/preferences"
	"uli1step.com/internal/storage/sqlite"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := sqlite.Open(cfg.DBPath)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	cache := devices.NewCache()
	client := onestepgps.NewClient(cfg.OneStepGPSAPIKey, cfg.OneStepGPSBaseURL)
	prefsStore := preferences.NewStore(db)

	iconStore, err := icons.NewStore(cfg.UploadsDir)
	if err != nil {
		log.Fatalf("failed to init icon store: %v", err)
	}

	pollCtx, stopPolling := context.WithCancel(context.Background())
	defer stopPolling()
	go devices.StartPoller(pollCtx, cache, client, cfg.PollInterval)

	handler := api.NewRouter(api.Config{
		Cache:            cache,
		PreferencesStore: prefsStore,
		IconStore:        iconStore,
		StreamInterval:   cfg.PollInterval,
		AllowedOrigin:    cfg.AllowedOrigin,
	})

	// WriteTimeout is intentionally left unset: /api/devices/stream is a
	// long-lived SSE connection, and a global WriteTimeout would cut it off.
	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	log.Printf("listening on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
