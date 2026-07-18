package api

import (
	"errors"
	"log"
	"net/http"
	"regexp"

	"uli1step.com/internal/devices"
	"uli1step.com/internal/icons"
)

// maxIconUploadBytes caps icon uploads before they're fully read into
// memory/disk.
const maxIconUploadBytes = 2 << 20 // 2 MiB

// validDeviceID guards against path-traversal-like input reaching the
// filesystem via deviceID, independent of the cache membership check below.
var validDeviceID = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)

func handleUploadIcon(cache *devices.Cache, store *icons.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		deviceID := r.PathValue("deviceID")
		if !validDeviceID.MatchString(deviceID) || !deviceExists(cache, deviceID) {
			http.Error(w, "unknown device", http.StatusNotFound)
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, maxIconUploadBytes)
		if err := r.ParseMultipartForm(maxIconUploadBytes); err != nil {
			http.Error(w, "upload too large or malformed", http.StatusBadRequest)
			return
		}

		file, _, err := r.FormFile("icon")
		if err != nil {
			http.Error(w, "missing icon file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		url, err := store.Save(deviceID, file)
		if err != nil {
			if errors.Is(err, icons.ErrUnsupportedType) {
				http.Error(w, "unsupported file type", http.StatusUnsupportedMediaType)
				return
			}
			log.Printf("save icon: %v", err)
			http.Error(w, "failed to save icon", http.StatusInternalServerError)
			return
		}

		writeJSON(w, http.StatusOK, map[string]string{"url": url})
	}
}

func deviceExists(cache *devices.Cache, deviceID string) bool {
	for _, d := range cache.Get() {
		if d.DeviceID == deviceID {
			return true
		}
	}
	return false
}
