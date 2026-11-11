package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"miniplayer-backend/internal/models"
)

// SystemHandlers mirrors app/api/system.py. Note: app/api/server.py
// registers an identical `/api/server/info` route; FastAPI would use
// whichever router was included first (system before server in main.py),
// so this is the version that wins in Python too.
type SystemHandlers struct{}

func (h *SystemHandlers) Register(r chi.Router) {
	r.Get("/api/system/info", h.systemInfo)
	r.Get("/system/info", h.systemInfo)
	r.Get("/api/server/info", h.serverInfo)
	r.Get("/server/info", h.serverInfo)
	r.Get("/api", h.apiRoot)
}

func (h *SystemHandlers) systemInfo(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, models.NewRespData(200, "success", map[string]interface{}{
		"version":  "1.0.0",
		"platform": "windows",
		"name":     "now-playing-python",
	}))
}

func (h *SystemHandlers) serverInfo(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, models.NewRespData(200, "success", map[string]interface{}{"version": "1.0.0"}))
}

func (h *SystemHandlers) apiRoot(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{"message": "NowPlaying API"})
}
