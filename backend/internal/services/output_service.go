package services

import (
	"os"
	"path/filepath"
	"strings"

	"miniplayer-backend/internal/models"
)

// OutputService mirrors app/services/output_service.py: writes the current
// track's title/author/custom-template text to Outputs/*.txt in the
// process's current working directory, for OBS-style text source overlays.
type OutputService struct {
	settings *SettingsService
}

func NewOutputService(settings *SettingsService) *OutputService {
	return &OutputService{settings: settings}
}

// OutputAsync mirrors output_service.output_async. Each file write is
// best-effort or isolated, same as the Python try/except-per-file.
func (o *OutputService) OutputAsync(track models.Track) {
	cwd, err := os.Getwd()
	if err != nil {
		return
	}
	outDir := filepath.Join(cwd, "Outputs")
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return
	}

	_ = os.WriteFile(filepath.Join(outDir, "title.txt"), []byte(track.Title), 0o644)
	_ = os.WriteFile(filepath.Join(outDir, "author.txt"), []byte(track.Author), 0o644)

	out := o.settings.GetOutput()
	content := out.Template
	content = strings.ReplaceAll(content, "{title}", track.Title)
	content = strings.ReplaceAll(content, "{author}", track.Author)
	content = strings.ReplaceAll(content, "{album}", track.Album)
	_ = os.WriteFile(filepath.Join(outDir, "custom.txt"), []byte(content), 0o644)
}
