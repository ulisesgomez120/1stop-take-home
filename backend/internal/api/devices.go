package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"uli1step.com/internal/devices"
)

func handleGetDevices(cache *devices.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, toDeviceDTOs(cache.Get()))
	}
}

// handleDeviceStream is a stub SSE endpoint: it writes the current cache
// snapshot as a "data:" event on an interval. Reconnect/resilience polish is
// deferred to a later phase.
func handleDeviceStream(cache *devices.Cache, interval time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "streaming unsupported", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		writeSnapshot := func() bool {
			payload, err := json.Marshal(toDeviceDTOs(cache.Get()))
			if err != nil {
				return false
			}
			if _, err := fmt.Fprintf(w, "data: %s\n\n", payload); err != nil {
				return false
			}
			flusher.Flush()
			return true
		}

		if !writeSnapshot() {
			return
		}

		for {
			select {
			case <-r.Context().Done():
				return
			case <-ticker.C:
				if !writeSnapshot() {
					return
				}
			}
		}
	}
}
