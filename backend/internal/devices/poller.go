package devices

import (
	"context"
	"log"
	"time"

	"uli1step.com/internal/onestepgps"
)

// Fetcher is the subset of onestepgps.Client used by the poller, kept as an
// interface so the poller can be tested without a live API call.
type Fetcher interface {
	FetchDevices(ctx context.Context) ([]onestepgps.Device, error)
}

// StartPoller fetches devices immediately, then again on every tick of the
// given interval, updating cache with each successful fetch. Failed polls
// are logged and the cache keeps serving its last good snapshot. Runs until
// ctx is cancelled.
func StartPoller(ctx context.Context, cache *Cache, fetcher Fetcher, interval time.Duration) {
	poll := func() {
		devs, err := fetcher.FetchDevices(ctx)
		if err != nil {
			log.Printf("poll devices: %v", err)
			return
		}
		cache.set(devs)
		log.Printf("poll devices: refreshed cache with %d devices", len(devs))
	}

	poll()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			poll()
		}
	}
}
