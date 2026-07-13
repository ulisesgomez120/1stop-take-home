package onestepgps

import (
	"encoding/json"
	"os"
	"testing"
	"time"
)

func TestMapDevices(t *testing.T) {
	data, err := os.ReadFile("testdata/devices_response.json")
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}

	var raw rawResponse
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("unmarshal fixture: %v", err)
	}

	devices := mapDevices(raw)
	if len(devices) != 7 {
		t.Fatalf("got %d devices, want 7", len(devices))
	}

	first := devices[0]
	if first.DeviceID != "6kNRGedyQUuiPk81f07-1V" {
		t.Errorf("DeviceID = %q, want %q", first.DeviceID, "6kNRGedyQUuiPk81f07-1V")
	}
	if first.DisplayName != "#5" {
		t.Errorf("DisplayName = %q, want %q", first.DisplayName, "#5")
	}
	if first.ActiveState != "active" {
		t.Errorf("ActiveState = %q, want %q", first.ActiveState, "active")
	}
	if !first.Online {
		t.Error("Online = false, want true")
	}
	if first.Lat != 37.376804299999996 {
		t.Errorf("Lat = %v, want %v", first.Lat, 37.376804299999996)
	}
	if first.Lng != -121.8262862 {
		t.Errorf("Lng = %v, want %v", first.Lng, -121.8262862)
	}
	if first.Speed != 0 {
		t.Errorf("Speed = %v, want 0", first.Speed)
	}
	wantTime := time.Date(2026, 7, 12, 3, 20, 29, 0, time.UTC)
	if !first.DtTracker.Equal(wantTime) {
		t.Errorf("DtTracker = %v, want %v", first.DtTracker, wantTime)
	}
	if first.DriveStatus != "off" {
		t.Errorf("DriveStatus = %q, want %q", first.DriveStatus, "off")
	}
}
