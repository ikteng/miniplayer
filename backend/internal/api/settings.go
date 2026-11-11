package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"miniplayer-backend/internal/models"
	"miniplayer-backend/internal/services"
)

// SettingsHandlers mirrors app/api/settings.py.
type SettingsHandlers struct {
	Settings *services.SettingsService
}

func (h *SettingsHandlers) Register(r chi.Router) {
	r.Get("/api/settings", h.getSettings)
	r.Post("/api/settings", h.postSettings)
	r.Get("/api/settings/general", h.getGeneral)
	r.Post("/api/settings/general", h.postGeneral)
	r.Put("/api/settings/general", h.postGeneral)
}

func (h *SettingsHandlers) getSettings(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, h.Settings.GetSettings())
}

func (h *SettingsHandlers) postSettings(w http.ResponseWriter, r *http.Request) {
	data, err := decodeJSONBody(r)
	if err != nil {
		writeJSON(w, http.StatusOK, models.NewRespData(500, err.Error(), nil))
		return
	}
	h.Settings.SetSettings(data)
	writeJSON(w, http.StatusOK, models.NewRespData(200, "success", nil))
}

func (h *SettingsHandlers) getGeneral(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, h.Settings.GetSettingsGeneral())
}

func (h *SettingsHandlers) postGeneral(w http.ResponseWriter, r *http.Request) {
	data, err := decodeJSONBody(r)
	if err != nil {
		writeJSON(w, http.StatusOK, models.NewRespData(500, err.Error(), nil))
		return
	}
	h.Settings.UpdateSettingsGeneral(data)
	writeJSON(w, http.StatusOK, models.NewRespData(200, "success", nil))
}
