// Package preferences provides SQLite-backed CRUD for the single-user
// preferences row (sort order, hidden devices, custom icons).
package preferences

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

// Preferences is the domain representation of the singleton preferences row.
type Preferences struct {
	SortBy          string            `json:"sort_by"`
	SortDir         string            `json:"sort_dir"`
	HiddenDeviceIDs []string          `json:"hidden_device_ids"`
	DeviceIcons     map[string]string `json:"device_icons"`
	UpdatedAt       time.Time         `json:"updated_at"`
}

// Store reads and writes the singleton preferences row.
type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

// Get returns the current preferences.
func (s *Store) Get(ctx context.Context) (*Preferences, error) {
	row := s.db.QueryRowContext(ctx, `SELECT sort_by, sort_dir, hidden_device_ids, device_icons, updated_at
		FROM preferences WHERE id = 1`)

	var (
		p                     Preferences
		hiddenJSON, iconsJSON string
	)
	if err := row.Scan(&p.SortBy, &p.SortDir, &hiddenJSON, &iconsJSON, &p.UpdatedAt); err != nil {
		return nil, fmt.Errorf("get preferences: %w", err)
	}

	if err := json.Unmarshal([]byte(hiddenJSON), &p.HiddenDeviceIDs); err != nil {
		return nil, fmt.Errorf("decode hidden_device_ids: %w", err)
	}
	if err := json.Unmarshal([]byte(iconsJSON), &p.DeviceIcons); err != nil {
		return nil, fmt.Errorf("decode device_icons: %w", err)
	}

	return &p, nil
}

// Update overwrites the singleton preferences row and returns the saved
// result.
func (s *Store) Update(ctx context.Context, p Preferences) (*Preferences, error) {
	hiddenJSON, err := json.Marshal(p.HiddenDeviceIDs)
	if err != nil {
		return nil, fmt.Errorf("encode hidden_device_ids: %w", err)
	}
	iconsJSON, err := json.Marshal(p.DeviceIcons)
	if err != nil {
		return nil, fmt.Errorf("encode device_icons: %w", err)
	}

	_, err = s.db.ExecContext(ctx, `UPDATE preferences
		SET sort_by = ?, sort_dir = ?, hidden_device_ids = ?, device_icons = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = 1`, p.SortBy, p.SortDir, string(hiddenJSON), string(iconsJSON))
	if err != nil {
		return nil, fmt.Errorf("update preferences: %w", err)
	}

	return s.Get(ctx)
}
