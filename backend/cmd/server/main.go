package main

import (
	"context"
	"log"
	"net/http"

	"uli1step.com/internal/config"
	"uli1step.com/internal/devices"
	"uli1step.com/internal/onestepgps"
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

	pollCtx, stopPolling := context.WithCancel(context.Background())
	defer stopPolling()
	go devices.StartPoller(pollCtx, cache, client, cfg.PollInterval)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	addr := ":" + cfg.Port
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
