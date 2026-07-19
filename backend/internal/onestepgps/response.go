package onestepgps

import "time"

// rawResponse mirrors the top-level shape of the OneStepGPS API response.
// Only fields this app uses are declared; everything else is dropped by
// encoding/json during unmarshalling.
type rawResponse struct {
	ResultList []rawDevice `json:"result_list"`
}

type rawDevice struct {
	DeviceID          string         `json:"device_id"`
	DisplayName       string         `json:"display_name"`
	ActiveState       string         `json:"active_state"`
	Online            bool           `json:"online"`
	LatestDevicePoint rawDevicePoint `json:"latest_device_point"`
}

type rawDevicePoint struct {
	Lat         float64        `json:"lat"`
	Lng         float64        `json:"lng"`
	Angle       float64        `json:"angle"`
	Speed       float64        `json:"speed"`
	DtTracker   time.Time      `json:"dt_tracker"`
	DeviceState rawDeviceState `json:"device_state"`
}

type rawDeviceState struct {
	DriveStatus string `json:"drive_status"`
}

func mapDevices(raw rawResponse) []Device {
	devices := make([]Device, 0, len(raw.ResultList))
	for _, rd := range raw.ResultList {
		devices = append(devices, Device{
			DeviceID:    rd.DeviceID,
			DisplayName: rd.DisplayName,
			ActiveState: rd.ActiveState,
			Online:      rd.Online,
			Lat:         rd.LatestDevicePoint.Lat,
			Lng:         rd.LatestDevicePoint.Lng,
			Heading:     rd.LatestDevicePoint.Angle,
			Speed:       rd.LatestDevicePoint.Speed,
			DtTracker:   rd.LatestDevicePoint.DtTracker,
			DriveStatus: rd.LatestDevicePoint.DeviceState.DriveStatus,
		})
	}
	return devices
}
