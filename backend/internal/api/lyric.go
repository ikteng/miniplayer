package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"miniplayer-backend/internal/models"
	"miniplayer-backend/internal/services"
)

// LyricHandlers mirrors app/api/lyric.py.
type LyricHandlers struct {
	Lyric *services.LyricService
}

func (h *LyricHandlers) Register(r chi.Router) {
	r.Get("/api/lyric", h.getLyric)
	r.Get("/lyric", h.getLyric)
	r.Get("/api/lyric/settings", h.getLyricSettings)
	r.Get("/api/lyric/settings/common", h.getLyricSettingsCommon)
	r.Put("/api/lyric/settings", h.putLyricSettings)
	r.Put("/api/lyric/settings/common", h.putLyricSettingsCommon)
}

func (h *LyricHandlers) getLyric(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")
	author := r.URL.Query().Get("author")
	if title != "" || author != "" {
		writeJSON(w, http.StatusOK, h.Lyric.GetLyricByTitleAuthor(title, author))
		return
	}
	writeJSON(w, http.StatusOK, h.Lyric.GetLyric())
}

func (h *LyricHandlers) getLyricSettings(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		id = "main"
	}
	writeJSON(w, http.StatusOK, h.Lyric.GetSettings(id))
}

func (h *LyricHandlers) getLyricSettingsCommon(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, h.Lyric.GetCommonSettings())
}

func (h *LyricHandlers) putLyricSettings(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		id = "main"
	}
	data, err := decodeJSONBody(r)
	if err != nil {
		writeJSON(w, http.StatusOK, models.NewRespData(500, err.Error(), nil))
		return
	}
	h.Lyric.UpdateSettings(id, data)
	writeJSON(w, http.StatusOK, models.NewRespData(200, "success", nil))
}

func (h *LyricHandlers) putLyricSettingsCommon(w http.ResponseWriter, r *http.Request) {
	data, err := decodeJSONBody(r)
	if err != nil {
		writeJSON(w, http.StatusOK, models.NewRespData(500, err.Error(), nil))
		return
	}
	h.Lyric.UpdateCommonSettings(data)
	writeJSON(w, http.StatusOK, models.NewRespData(200, "success", nil))
}
