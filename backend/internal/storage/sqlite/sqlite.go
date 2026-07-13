// Package sqlite manages the SQLite connection and schema setup.
package sqlite

import (
	"database/sql"
	_ "embed"
	"fmt"

	_ "modernc.org/sqlite"
)

//go:embed migrations/0001_init.sql
var initSchema string

// Open opens (creating if necessary) the SQLite database at dbPath, runs
// the embedded migration, and seeds the singleton preferences row.
func Open(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	if _, err := db.Exec(initSchema); err != nil {
		db.Close()
		return nil, fmt.Errorf("run migration: %w", err)
	}

	if _, err := db.Exec(`INSERT OR IGNORE INTO preferences (id, sort_by, sort_dir, hidden_device_ids, device_icons, updated_at)
		VALUES (1, 'display_name', 'asc', '[]', '{}', CURRENT_TIMESTAMP)`); err != nil {
		db.Close()
		return nil, fmt.Errorf("seed preferences row: %w", err)
	}

	return db, nil
}
