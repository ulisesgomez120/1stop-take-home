package api

import (
	"time"

	"uli1step.com/internal/onestepgps"
)

// deviceDTO is the JSON shape returned to the frontend for a device.
type deviceDTO struct {
	DeviceID    string    `json:"device_id"`
	DisplayName string    `json:"display_name"`
	ActiveState string    `json:"active_state"`
	Online      bool      `json:"online"`
	Lat         float64   `json:"lat"`
	Lng         float64   `json:"lng"`
	Speed       float64   `json:"speed"`
	DtTracker   time.Time `json:"dt_tracker"`
	DriveStatus string    `json:"drive_status"`
}

func toDeviceDTOs(devices []onestepgps.Device) []deviceDTO {
	dtos := make([]deviceDTO, 0, len(devices))
	for _, d := range devices {
		dtos = append(dtos, deviceDTO{
			DeviceID:    d.DeviceID,
			DisplayName: d.DisplayName,
			ActiveState: d.ActiveState,
			Online:      d.Online,
			Lat:         d.Lat,
			Lng:         d.Lng,
			Speed:       d.Speed,
			DtTracker:   d.DtTracker,
			DriveStatus: d.DriveStatus,
		})
	}
	return dtos
}
