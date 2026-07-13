package api

import (
	"encoding/json"
	"log"
	"net/http"

	"uli1step.com/internal/preferences"
)

// maxPreferencesBodyBytes bounds the PUT request body to prevent a client
// from streaming an unbounded payload into the server.
const maxPreferencesBodyBytes = 1 << 20 // 1 MiB

func handleGetPreferences(store *preferences.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prefs, err := store.Get(r.Context())
		if err != nil {
			log.Printf("get preferences: %v", err)
			http.Error(w, "failed to load preferences", http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, prefs)
	}
}

func handlePutPreferences(store *preferences.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, maxPreferencesBodyBytes)

		var p preferences.Preferences
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		updated, err := store.Update(r.Context(), p)
		if err != nil {
			log.Printf("update preferences: %v", err)
			http.Error(w, "failed to update preferences", http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, updated)
	}
}
