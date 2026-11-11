package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// PublicHandlers mirrors app/api/public.py: a catch-all under /public/{path}
// that always returns an empty object.
type PublicHandlers struct{}

func (h *PublicHandlers) Register(r chi.Router) {
	r.Get("/public/{path}", h.get)
}

func (h *PublicHandlers) get(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{})
}
