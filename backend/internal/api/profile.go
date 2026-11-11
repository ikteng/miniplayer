package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"miniplayer-backend/internal/models"
)

// ProfileHandlers mirrors app/api/profile.py.
type ProfileHandlers struct{}

func (h *ProfileHandlers) Register(r chi.Router) {
	r.Get("/api/profile", h.get)
	r.Post("/api/profile", h.post)
}

func (h *ProfileHandlers) get(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, models.NewRespData(200, "success", []interface{}{}))
}

func (h *ProfileHandlers) post(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, models.NewRespData(200, "success", nil))
}
