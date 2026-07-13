package devices

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"uli1step.com/internal/onestepgps"
)

type fakeFetcher struct {
	calls   atomic.Int32
	fail    bool
	devices []onestepgps.Device
}

func (f *fakeFetcher) FetchDevices(ctx context.Context) ([]onestepgps.Device, error) {
	f.calls.Add(1)
	if f.fail {
		return nil, errors.New("boom")
	}
	return f.devices, nil
}

func TestStartPoller_RefreshesCache(t *testing.T) {
	cache := NewCache()
	fetcher := &fakeFetcher{devices: []onestepgps.Device{{DeviceID: "a"}}}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go StartPoller(ctx, cache, fetcher, 10*time.Millisecond)

	deadline := time.After(time.Second)
	for {
		if len(cache.Get()) == 1 {
			break
		}
		select {
		case <-deadline:
			t.Fatal("cache was never populated")
		case <-time.After(5 * time.Millisecond):
		}
	}
}

func TestStartPoller_KeepsLastGoodSnapshotOnError(t *testing.T) {
	cache := NewCache()
	cache.set([]onestepgps.Device{{DeviceID: "stale-but-good"}})
	fetcher := &fakeFetcher{fail: true}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go StartPoller(ctx, cache, fetcher, 10*time.Millisecond)

	time.Sleep(50 * time.Millisecond)

	got := cache.Get()
	if len(got) != 1 || got[0].DeviceID != "stale-but-good" {
		t.Fatalf("cache = %+v, want last good snapshot preserved", got)
	}
	if fetcher.calls.Load() < 2 {
		t.Fatalf("expected poller to retry on failure, got %d calls", fetcher.calls.Load())
	}
}
