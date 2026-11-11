package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"miniplayer-backend/internal/detectors"
	"miniplayer-backend/internal/models"
)

// AudioHandlers mirrors app/api/audio.py.
type AudioHandlers struct{}

func (h *AudioHandlers) Register(r chi.Router) {
	r.Get("/api/audio/devices", h.devices)
	r.Get("/audio/devices", h.devices)
	r.Get("/api/audio/volume", h.volume)
	r.Get("/audio/volume", h.volume)
	r.Post("/api/audio/device", h.setDevice)
	r.Get("/api/audio/detectPlatform", h.detectPlatform)
}

func (h *AudioHandlers) devices(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, models.NewRespData(200, "success", []interface{}{}))
}

func (h *AudioHandlers) volume(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, models.NewRespData(200, "success", map[string]interface{}{"volume": 0}))
}

func (h *AudioHandlers) setDevice(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, models.NewRespData(200, "success", nil))
}

func (h *AudioHandlers) detectPlatform(w http.ResponseWriter, r *http.Request) {
	result := detectors.DetectPlatform()
	writeJSON(w, http.StatusOK, models.NewRespData(200, "success", result))
}
