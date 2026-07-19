# OneStepGPS Device Viewer

A small full-stack app that shows a live map of OneStepGPS devices. A Go
backend polls the OneStepGPS API and streams device updates to a Vue 3
frontend over Server-Sent Events. The map is the primary view; a pullout
panel holds the device list (hide/show, sort, custom icon upload) and a
per-device detail view. Preferences persist server-side in SQLite.

## Stack

- **Backend**: Go (stdlib `net/http`), SQLite (preferences), on-disk icon
  storage, SSE streaming.
- **Frontend**: Vue 3 + TypeScript + Pinia + Vite, `vue3-google-map` for
  Google Maps.

## Setup

Two API keys are required:

1. **OneStepGPS API key** — from your OneStepGPS account.
2. **Google Maps JavaScript API key** — from Google Cloud console, with the
   Maps JavaScript API enabled.

### Backend

```sh
cd backend
echo 'ONESTEPGPS_API_KEY=<your key>' > .env
go run ./cmd/server
```

Serves on `http://localhost:8080`. Creates `data.db` (SQLite, preferences)
and `uploads/` (device icons) in the working directory on first run.

### Frontend

```sh
cd frontend
cp .env.example .env   # then fill in VITE_GOOGLE_MAPS_API_KEY
npm install
npm run dev
```

Serves on `http://localhost:5173` (the backend's allowed CORS origin).

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
