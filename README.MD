# OneStepGPS Device Viewer

A small full-stack app that shows a live map of OneStepGPS devices. A Go
backend polls the OneStepGPS API and streams device updates to a Vue 3
frontend over Server-Sent Events. The map is the primary view; a pullout
panel holds the device list (visibility toggles, sorting) and a per-device
detail view (live fields, custom icon upload). Preferences persist
server-side in SQLite.

## Stack

- **Backend**: Go (stdlib `net/http`), SQLite (preferences), on-disk icon
  storage, SSE streaming.
- **Frontend**: Vue 3 + TypeScript + Pinia + Vite, `vue3-google-map` for
  Google Maps.

## Requirements

- **Go 1.25+** (with Go 1.21+ installed, `go run` auto-downloads the right
  toolchain on first run — that one-time download is normal).
- **Node 20.19+ or 22.12+** (Vite 7's requirement) and npm.
- No C compiler needed — SQLite is the pure-Go `modernc.org/sqlite` driver.
- Two API keys (see below).

## Setup

Two API keys are required:

1. **OneStepGPS API key** — from your OneStepGPS account
   (Account → API Access).
2. **Google Maps JavaScript API key** — from the Google Cloud console, with
   the **Maps JavaScript API** enabled on the project. Without it the app
   still runs; the map pane shows a notice and the device list/detail panel
   works normally.

Run the backend and frontend in two terminals:

### Terminal 1 — backend

```sh
cd backend
cp .env.example .env   # then fill in ONESTEPGPS_API_KEY
go run ./cmd/server
```

Serves on `http://localhost:8080`. On first run it creates `data.db`
(SQLite, preferences) and `uploads/` (device icons) in the `backend/`
directory — no database setup needed. Success looks like `listening on
:8080` followed by `poll devices: refreshed cache with N devices` every
10 seconds; if the API key is missing it exits immediately with
`ONESTEPGPS_API_KEY is required`, and a wrong key shows as repeated
`poll devices:` errors.

### Terminal 2 — frontend

```sh
cd frontend
cp .env.example .env   # then fill in VITE_GOOGLE_MAPS_API_KEY
npm install
npm run dev
```

Open `http://localhost:5173`. The port matters: it's the backend's allowed
CORS origin, so the dev server is pinned to 5173 (`strictPort`) and will
fail loudly rather than drift to another port.

### Quick health checks

- Backend up: `curl http://localhost:8080/healthz` → 200.
- Live data flowing: `curl http://localhost:8080/api/devices` → JSON array
  of devices.
- In the app: green dot in the top-right = SSE feed connected.

## Troubleshooting

| Symptom | Cause / fix |
| --- | --- |
| Backend exits with `ONESTEPGPS_API_KEY is required` | `backend/.env` missing or empty — copy `.env.example` and set the key. |
| Backend fails with `address already in use` | Something else is on port 8080 — stop it (`lsof -i :8080`) and rerun. |
| `npm run dev` fails with "Port 5173 is already in use" | Intentional (CORS origin is fixed) — free the port (`lsof -i :5173`) and rerun. |
| Map area shows "Google Maps API key is missing" | `frontend/.env` missing or `VITE_GOOGLE_MAPS_API_KEY` empty — set it and restart `npm run dev` (Vite only reads `.env` at startup). |
| Map loads but is gray / watermarked "for development purposes only" | The Maps key's Google Cloud project has no billing enabled or the Maps JavaScript API isn't enabled. |
| Devices table is empty, red dot / "reconnecting…" top-right | Backend isn't running or the OneStepGPS key is invalid — check the backend terminal's logs. |

## Architecture notes / tradeoffs

- **SSE over WebSocket**: updates are strictly server → client (device
  positions), so SSE is sufficient and simpler — native `EventSource`
  reconnect handling, plain HTTP, no upgrade dance. Client-initiated
  changes (preferences, icon uploads) go over normal REST calls.
- **SQLite singleton row for preferences**: the app is single-user by
  design (the take-home has no auth), so preferences are one row keyed by a
  fixed ID rather than a users table. Swapping to per-user rows would be a
  schema change, not an architecture change.
- **Icon upload hardening**: uploads are capped at 2 MiB, content-type is
  determined by sniffing the actual bytes (not the client's claimed MIME
  type), only PNG/JPEG/WebP are accepted, device IDs are validated against
  a strict pattern and checked against the live device cache, and files are
  stored under opaque server-generated names.
- **Marker rotation assumption**: markers rotate to the device's heading
  (`angle` from the OneStepGPS API) with plain CSS `rotate()`. This assumes
  the icon art points "up"/north at 0° — true for the built-in arrow, but
  there's no way to know a user-uploaded image's inherent orientation, so
  the same rotation is applied uniformly.
- **Map ID**: the map uses Google's `DEMO_MAP_ID` (required for
  `AdvancedMarkerElement`). Fine for development/demo; a production app
  would create a real Map ID in the Google Cloud console.

## Known limitations

- **Single-user, no auth** — anyone who can reach the backend shares one
  set of preferences.
- **No deployment** — local development only; `data.db` and `uploads/` are
  local files and not committed.
- **Label crowding** — device name labels are always visible on the map and
  can overlap when devices cluster at low zoom. Accepted as a known
  limitation rather than solved with marker clustering / collision
  avoidance, which is out of scope here.
