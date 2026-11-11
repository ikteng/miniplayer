package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"miniplayer-backend/internal/models"
)

// PluginsHandlers mirrors app/api/plugins.py.
type PluginsHandlers struct{}

func (h *PluginsHandlers) Register(r chi.Router) {
	r.Get("/api/plugins", h.get)
	r.Post("/api/plugins", h.post)
}

func (h *PluginsHandlers) get(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, models.NewRespData(200, "success", []interface{}{}))
}

func (h *PluginsHandlers) post(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, models.NewRespData(200, "success", nil))
}
