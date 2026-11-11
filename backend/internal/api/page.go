package api

import (
	"io/fs"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// PageHandlers mirrors app/api/page.py: serves the embedded index.html for
// the root/index/page/settings routes, falling back to an inline HTML
// string if it's missing (matching Python's FileResponse-or-fallback
// behavior without risking a panic on a missing file).
type PageHandlers struct {
	Dist fs.FS
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
	serveEmbeddedHTML(w, h.Dist, "index.html")
}

// serveEmbeddedHTML writes the named file out of the embedded dist
// filesystem, or the inline fallback HTML when it's missing.
func serveEmbeddedHTML(w http.ResponseWriter, dist fs.FS, name string) {
	data, err := fs.ReadFile(dist, name)
	if err != nil {
		serveFallbackHTML(w)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}

// serveFallbackHTML writes the inline fallback HTML when the file is
// missing, matching Python's fallback-to-inline-HTML behavior instead of a
// bare 404.
func serveFallbackHTML(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(fallbackIndexHTML))
}
