package api

import (
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"

	"miniplayer-backend/internal/models"
	"miniplayer-backend/internal/services"
)

// WidgetHandlers mirrors app/api/widget.py: serves frontend/dist/widget.html
// for /widget routes, plus a tiny settings endpoint for the "amuse" widget
// profile. Python's fallback to resources/templates/{widget,settings-widget}.html
// is never reachable in this repo (those files don't exist), so instead of
// replicating a dead FileNotFoundError path, this Go version falls back to
// the same inline HTML used by page.go, which is safe and behaviorally
// equivalent for every realistic case (frontend/dist/widget.html exists).
type WidgetHandlers struct {
	StaticDir string
	Settings  *services.SettingsService
}

func (h *WidgetHandlers) Register(r chi.Router) {
	r.Get("/widget", h.widget)
	r.Get("/widget/{profile_id}", h.widget)
	r.Get("/lyric-overlay", h.lyricOverlay)
	r.Get("/settings-widget", h.settingsWidget)
	r.Get("/widgets/amuse/settings", h.getAmuseSettings)
	r.Post("/widgets/amuse/settings", h.postAmuseSettings)
}

func (h *WidgetHandlers) widget(w http.ResponseWriter, r *http.Request) {
	widgetPath := filepath.Join(h.StaticDir, "widget.html")
	if _, err := statExists(widgetPath); err != nil {
		serveFallbackHTML(w)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	http.ServeFile(w, r, widgetPath)
}

// lyricOverlay serves frontend/dist/lyric-overlay.html, an OBS
// browser-source page showing synced (or plain) lyrics for the current
// track, following the same StaticDir-relative pattern as widget().
func (h *WidgetHandlers) lyricOverlay(w http.ResponseWriter, r *http.Request) {
	overlayPath := filepath.Join(h.StaticDir, "lyric-overlay.html")
	if _, err := statExists(overlayPath); err != nil {
		serveFallbackHTML(w)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	http.ServeFile(w, r, overlayPath)
}

func (h *WidgetHandlers) settingsWidget(w http.ResponseWriter, r *http.Request) {
	settingsWidgetPath := filepath.Join(h.StaticDir, "settings-widget.html")
	if _, err := statExists(settingsWidgetPath); err != nil {
		serveFallbackHTML(w)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	http.ServeFile(w, r, settingsWidgetPath)
}

func (h *WidgetHandlers) getAmuseSettings(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, h.Settings.GetSettingsWidget())
}

func (h *WidgetHandlers) postAmuseSettings(w http.ResponseWriter, r *http.Request) {
	data, err := decodeJSONBody(r)
	if err != nil {
		writeJSON(w, http.StatusOK, models.NewRespData(500, err.Error(), nil))
		return
	}
	h.Settings.UpdateSettingsWidget(data)
	writeJSON(w, http.StatusOK, models.NewRespData(200, "success", nil))
}
