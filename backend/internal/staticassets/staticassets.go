// Package staticassets embeds the built frontend (frontend/dist) into the
// Go binary so the server can run as a single standalone executable with no
// loose frontend files on disk.
//
// dist/ here is a copy of frontend/dist, placed by build-exe.ps1 (or
// manually: `npm run build` + copy) before `go build` runs. It must be
// populated with a real `npm run build` output before building the server,
// otherwise this package embeds the checked-in placeholder index.html and
// the server serves that instead of the real frontend.
package staticassets

import "embed"

//go:embed dist
var distFS embed.FS

// Dist returns the embedded frontend/dist filesystem, rooted at "dist" so
// callers see the same tree layout (index.html, widget.html, assets/...)
// that existed on disk.
func Dist() embed.FS {
	return distFS
}
