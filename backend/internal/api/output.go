package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"miniplayer-backend/internal/models"
	"miniplayer-backend/internal/services"
)

// OutputHandlers mirrors app/api/output.py.
type OutputHandlers struct {
	Settings *services.SettingsService
}

func (h *OutputHandlers) Register(r chi.Router) {
	r.Get("/api/output", h.get)
	r.Post("/api/output", h.post)
}

func (h *OutputHandlers) get(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, models.NewRespData(200, "success", map[string]interface{}{}))
}

// post mirrors app/api/output.py's post_output: replaces the in-memory
// output settings wholesale from the request body (note: unlike
// /api/settings, Python's post_output does NOT persist to settings.json -
// it only mutates settings_service.get_settings().output in memory, so this
// Go handler intentionally mirrors that and does not call save()).
func (h *OutputHandlers) post(w http.ResponseWriter, r *http.Request) {
	data, err := decodeJSONBody(r)
	if err != nil {
		writeJSON(w, http.StatusOK, models.NewRespData(500, err.Error(), nil))
		return
	}
	h.Settings.ReplaceOutputInMemory(data)
	writeJSON(w, http.StatusOK, models.NewRespData(200, "success", nil))
}
