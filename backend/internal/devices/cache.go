// Package devices holds an in-memory cache of the latest device snapshot,
// kept fresh by a poller that calls the OneStepGPS client on an interval.
package devices

import (
	"sync"

	"uli1step.com/internal/onestepgps"
)

// Cache holds the most recently fetched snapshot of devices.
type Cache struct {
	mu      sync.RWMutex
	devices []onestepgps.Device
}

func NewCache() *Cache {
	return &Cache{}
}

// Get returns the current snapshot of devices.
func (c *Cache) Get() []onestepgps.Device {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.devices
}

func (c *Cache) set(devices []onestepgps.Device) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.devices = devices
}
