CREATE TABLE IF NOT EXISTS preferences (
  id INTEGER PRIMARY KEY CHECK (id = 1),
  sort_by TEXT,
  sort_dir TEXT,
  hidden_device_ids TEXT,   -- JSON array
  device_icons TEXT,        -- JSON object: device_id -> icon url
  updated_at DATETIME
);
