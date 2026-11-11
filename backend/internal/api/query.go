// Package api mirrors backend/app/api/*.py: one file per Python router
// module, registering the same paths/methods and producing the same JSON
// shapes.
package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"miniplayer-backend/internal/models"
	"miniplayer-backend/internal/services"
)

// QueryHandlers mirrors app/api/query.py.
type QueryHandlers struct {
	NowPlaying *services.NowPlayingService
	Cover      *services.CoverService
}

func (h *QueryHandlers) Register(r chi.Router) {
	r.Get("/query", h.query)
	r.Get("/api/query", h.query)
	r.Get("/query/progress", h.queryProgress)
	r.Get("/api/query/progress", h.queryProgress)
	r.Get("/query/player", h.queryPlayer)
	r.Get("/api/query/player", h.queryPlayer)
	r.Get("/query/track", h.queryTrack)
	r.Get("/api/query/track", h.queryTrack)
	r.Post("/cover/convert", h.coverConvert)
	r.Post("/api/cover/convert", h.coverConvert)
	r.Post("/cover/videoUrl", h.coverVideo)
	r.Post("/api/cover/videoUrl", h.coverVideo)
	r.Get("/api/query/hasSong", h.hasSong)
	r.Get("/api/query/isConnected", h.isConnected)
}

func (h *QueryHandlers) query(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, h.NowPlaying.Query())
}

func (h *QueryHandlers) queryProgress(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, h.NowPlaying.QueryProgress())
}

func (h *QueryHandlers) queryPlayer(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, h.NowPlaying.QueryPlayer())
}

func (h *QueryHandlers) queryTrack(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, h.NowPlaying.QueryTrack())
}

func (h *QueryHandlers) coverConvert(w http.ResponseWriter, r *http.Request) {
	var dto models.CoverConvertDTO
	defer r.Body.Close()
	_ = json.NewDecoder(r.Body).Decode(&dto)
	result := h.Cover.ConvertToBase64(dto.CoverURL)
	writeJSON(w, http.StatusOK, result)
}

func (h *QueryHandlers) coverVideo(w http.ResponseWriter, r *http.Request) {
	var dto models.CoverVideoDTO
	defer r.Body.Close()
	_ = json.NewDecoder(r.Body).Decode(&dto)
	result := h.Cover.GetVideoURL(dto.SongTitle, dto.SongAuthor)
	// Python returns the raw result (None) directly, not wrapped in RespData.
	writeJSON(w, http.StatusOK, result)
}

func (h *QueryHandlers) hasSong(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, h.NowPlaying.HasSong())
}

func (h *QueryHandlers) isConnected(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, h.NowPlaying.IsConnected())
}
