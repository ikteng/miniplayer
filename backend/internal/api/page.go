package api

import (
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
)

// PageHandlers mirrors app/api/page.py: serves frontend/dist/index.html for
// the root/index/page/settings routes, falling back to an inline HTML
// string if the file is missing (matching Python's FileResponse-or-fallback
// behavior without risking a panic on a missing file).
type PageHandlers struct {
	StaticDir string
}

const fallbackIndexHTML = `<!DOCTYPE html>
<html>
<head><title>NowPlaying</title></head>
<body><h1>NowPlaying</h1><p>Frontend build not found.</p></body>
</html>`

func (h *PageHandlers) Register(r chi.Router) {
	r.Get("/", h.index)
	r.Get("/index", h.index)
	r.Get("/page/{path}", h.index)
	r.Get("/settings/{path}", h.index)
}

func (h *PageHandlers) index(w http.ResponseWriter, r *http.Request) {
	indexPath := filepath.Join(h.StaticDir, "index.html")
	if _, err := statExists(indexPath); err != nil {
		serveFallbackHTML(w)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	http.ServeFile(w, r, indexPath)
}

// serveFallbackHTML writes the inline fallback HTML when the file is
// missing, matching Python's fallback-to-inline-HTML behavior instead of a
// bare 404 (http.ServeFile would otherwise 404 on a missing file).
func serveFallbackHTML(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(fallbackIndexHTML))
}
