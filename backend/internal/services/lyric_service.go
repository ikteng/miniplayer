package services

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"miniplayer-backend/internal/models"
)

// LyricService mirrors app/services/lyric_service.py, including its
// lyric_cache directory under backend/settings/.
type LyricService struct {
	cacheDir string

	mu            sync.RWMutex
	currentLyric  models.Lyric
	settingsCommon models.SettingsLyricCommon
}

func NewLyricService(settingsDir string) *LyricService {
	cacheDir := filepath.Join(settingsDir, "lyric_cache")
	_ = os.MkdirAll(cacheDir, 0o755)
	return &LyricService{
		cacheDir:       cacheDir,
		currentLyric:   models.NewLyric(),
		settingsCommon: models.NewSettingsLyricCommon(),
	}
}

func (l *LyricService) GetLyric() models.Lyric {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.currentLyric
}

// GetLyricByTitleAuthor mirrors lyric_service.get_lyric_by_title_author:
// looks up (and creates, if missing) a per-(title,author) JSON cache file.
func (l *LyricService) GetLyricByTitleAuthor(title, author string) models.Lyric {
	key := title + "_" + author
	key = strings.ReplaceAll(key, "/", "_")
	key = strings.ReplaceAll(key, "\\", "_")
	cacheFile := filepath.Join(l.cacheDir, key+".json")

	if data, err := os.ReadFile(cacheFile); err == nil {
		var lyric models.Lyric
		if err := json.Unmarshal(data, &lyric); err == nil {
			l.mu.Lock()
			l.currentLyric = lyric
			l.mu.Unlock()
			return lyric
		}
	}

	lyric := models.Lyric{Title: title, Author: author, Source: "netease"}
	if data, err := json.Marshal(lyric); err == nil {
		_ = os.WriteFile(cacheFile, data, 0o644)
	}

	l.mu.Lock()
	l.currentLyric = lyric
	l.mu.Unlock()
	return lyric
}

// GetSettings mirrors lyric_service.get_settings(id): always returns
// defaults, same as Python (the `id` parameter is unused upstream).
func (l *LyricService) GetSettings(id string) models.SettingsLyric {
	return models.NewSettingsLyric()
}

func (l *LyricService) GetCommonSettings() models.SettingsLyricCommon {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.settingsCommon
}

func (l *LyricService) UpdateSettings(id string, data map[string]interface{}) {
	// Mirrors lyric_service.update_settings: a no-op in Python.
}

func (l *LyricService) UpdateCommonSettings(data map[string]interface{}) {
	l.mu.Lock()
	mergeJSON(&l.settingsCommon, data)
	l.mu.Unlock()
}
