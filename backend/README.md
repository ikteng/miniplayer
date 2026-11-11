# backend

Go now-playing backend (chi router + gorilla/websocket), replacing the
original Python/FastAPI implementation for lower idle memory and faster
startup as a background service.

## Layout

```
backend/
  cmd/server/main.go        # wiring: config, services, routes, CORS, graceful shutdown
  internal/config/          # settings/port/paths
  internal/models/          # JSON structs (field names/aliases match the frontend contract)
  internal/detectors/       # AudioDetector (shells out to GetMusicStatus.exe) + platform stub
  internal/services/        # settings/lyric/cover/output/now-playing services + on_update port
  internal/ws/              # lyric WebSocket ConnectionManager (fan-out broadcaster)
  internal/api/             # one file per route module
  assets/                   # prebuilt .NET AudioService detector executables (not ported)
  settings/                 # runtime settings.json, lyric cache
```

## Dependencies

- `github.com/go-chi/chi/v5` — lightweight router; chosen over bare
  `net/http.ServeMux` because routes have both `/x` and `/api/x` variants
  and path params (`{path}`, `{profile_id}`), which chi handles cleanly
  without pulling in a full framework.
- `github.com/gorilla/websocket` — WebSocket upgrade/read/write for the
  lyric WS endpoint.
- No SQLite/crypto/image libraries: cross-checking the old Python
  `requirements.txt` against the actual service code showed sqlalchemy,
  aiosqlite, cryptography, pycryptodome, beautifulsoup4, lxml, Pillow,
  pyyaml and jinja2 were never imported by any service/API module that was
  actually wired up, so none of that was ported.

## Assets/settings

`internal/config.Load()` walks upward from the executable's directory (and
the current working directory) looking for a `backend` subdirectory, then
points `AssetsDir`/`SettingsDir` at `backend/assets` and `backend/settings`,
and `StaticDir` at `frontend/dist`. Override with `MINIPLAY_ROOT` (repo
root) or `MINIPLAY_SERVER_PORT` env vars if needed.

## Build & run

```powershell
cd backend
go build -o server.exe ./cmd/server
./server.exe
```

or directly:

```powershell
go run ./cmd/server
```

Listens on `127.0.0.1:9863` by default. Run `go build ./...` / `go vet ./...`
from `backend/` to validate the whole module.

## Known gaps vs. the old Python backend

- `/api/audio/detectPlatform` is stubbed: Python's SMTC platform detection
  used `winsdk` (WinRT bindings), which has no pure-Go equivalent without
  CGo. The Go handler returns the same "platform unavailable" shape Python
  fell back to when the `winsdk` import failed, rather than violating the
  pure-Go build.
- The old `app/api/server.py` and `app/api/system.py` both registered an
  identical `/api/server/info` route; chi panics on duplicate route
  registration, so `system.go` registers it once and `server.go` is a
  documented no-op (same response shape, no duplicate panic).
- `widget.py`'s fallback to `resources/templates/*.html` (files that never
  existed in this repo) is replaced with the same inline fallback HTML used
  by `page.go`, since that Python fallback path would itself throw
  `FileNotFoundError` if ever exercised.
- `GetMusicStatus.exe` is NetEase-specific: it uses UI Automation against the
  `NeteaseMusicHelper` window class, not SMTC, so it misses browser tabs,
  Spotify, VLC, etc. Fixed by `assets/AudioService/GetSmtcStatus/GetSmtcStatus.exe`,
  a small self-contained .NET tool built on the real
  `GlobalSystemMediaTransportControlsSessionManager` SMTC API that reports
  whatever app currently holds the active media session — this is what
  `audio_detector.go` launches now, with `GetMusicStatus.exe` kept only as a
  fallback if the new exe is ever missing.
