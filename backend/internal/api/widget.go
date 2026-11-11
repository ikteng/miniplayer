package api

import (
	"io/fs"
	"net/http"

	"github.com/go-chi/chi/v5"

	"miniplayer-backend/internal/models"
	"miniplayer-backend/internal/services"
)

// WidgetHandlers mirrors app/api/widget.py: serves the embedded widget.html
// for /widget routes, plus a tiny settings endpoint for the "amuse" widget
// profile.
type WidgetHandlers struct {
	Dist     fs.FS
	Settings *services.SettingsService
}

func (h *WidgetHandlers) Register(r chi.Router) {
	r.Get("/widget", h.widget)
	r.Get("/widget/{profile_id}", h.widget)
	r.Get("/settings-widget", h.settingsWidget)
	r.Get("/widgets/amuse/settings", h.getAmuseSettings)
	r.Post("/widgets/amuse/settings", h.postAmuseSettings)
}

func (h *WidgetHandlers) widget(w http.ResponseWriter, r *http.Request) {
	serveEmbeddedHTML(w, h.Dist, "widget.html")
}

func (h *WidgetHandlers) settingsWidget(w http.ResponseWriter, r *http.Request) {
	serveEmbeddedHTML(w, h.Dist, "settings-widget.html")
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
