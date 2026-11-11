package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"

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

// detectPlatform mirrors app/detectors/platform_detector.py's "winsdk not
// available" branch: this port has no SMTC session enumeration independent
// of the AudioDetector polling loop (see internal/detectors/audio_detector.go),
// so there's nothing to detect here beyond that fixed response shape.
func (h *AudioHandlers) detectPlatform(w http.ResponseWriter, r *http.Request) {
	result := map[string]interface{}{"platform": nil, "reason": "winsdk not available"}
	writeJSON(w, http.StatusOK, models.NewRespData(200, "success", result))
}
