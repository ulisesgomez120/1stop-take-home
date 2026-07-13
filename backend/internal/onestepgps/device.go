// Package onestepgps provides a client for the OneStepGPS device API and
// maps its response shape into this app's domain Device type.
package onestepgps

import "time"

// Device is the domain representation of a tracked vehicle/device, trimmed
// down to only the fields this app needs.
type Device struct {
	DeviceID    string
	DisplayName string
	ActiveState string
	Online      bool
	Lat         float64
	Lng         float64
	Speed       float64
	DtTracker   time.Time
	DriveStatus string
}
