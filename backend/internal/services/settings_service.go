package services

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	"miniplayer-backend/internal/models"
)

// SettingsService mirrors app/services/settings_service.py: loads/saves
// backend/settings/settings.json, preserving its on-disk shape exactly.
type SettingsService struct {
	mu           sync.RWMutex
	settingsDir  string
	settingsFile string
	wrapper      models.SettingsWrapper
}

func NewSettingsService(settingsDir string) *SettingsService {
	_ = os.MkdirAll(settingsDir, 0o755)
	s := &SettingsService{
		settingsDir:  settingsDir,
		settingsFile: filepath.Join(settingsDir, "settings.json"),
	}
	s.wrapper = s.load()
	return s
}

func (s *SettingsService) load() models.SettingsWrapper {
	if data, err := os.ReadFile(s.settingsFile); err == nil {
		var w models.SettingsWrapper
		// Start from defaults so any fields missing from the file (e.g. an
		// older settings.json) keep their Python-equivalent defaults,
		// mirroring pydantic's behavior of filling in defaults for omitted
		// fields during SettingsWrapper(**data).
		w = models.NewSettingsWrapper()
		if err := json.Unmarshal(data, &w); err == nil {
			return w
		}
	}
	w := models.NewSettingsWrapper()
	s.saveWrapper(w)
	return w
}

func (s *SettingsService) saveWrapper(w models.SettingsWrapper) {
	data, err := json.MarshalIndent(w, "", "  ")
	if err != nil {
		return
	}
	_ = os.WriteFile(s.settingsFile, data, 0o644)
}

func (s *SettingsService) save() {
	s.saveWrapper(s.wrapper)
}

func (s *SettingsService) GetSettings() models.SettingsWrapper {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.wrapper
}

func (s *SettingsService) GetSettingsGeneral() models.SettingsGeneral {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.wrapper.General
}

func (s *SettingsService) GetSettingsWidget() models.SettingsWidget {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.wrapper.Widget
}

func (s *SettingsService) UpdateSettingsWidget(data map[string]interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	applyFieldsToWidget(&s.wrapper.Widget, data)
	s.save()
}

func (s *SettingsService) UpdateSettingsGeneral(data map[string]interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	applyFieldsToGeneral(&s.wrapper.General, data)
	s.save()
}

// SetSettings mirrors settings_service.set_settings: applies each top-level
// section present in the payload, leaving sections not present untouched.
func (s *SettingsService) SetSettings(data map[string]interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if v, ok := data["general"].(map[string]interface{}); ok {
		applyFieldsToGeneral(&s.wrapper.General, v)
	}
	if v, ok := data["lyric"].(map[string]interface{}); ok {
		applyFieldsToLyric(&s.wrapper.Lyric, v)
	}
	if v, ok := data["output"].(map[string]interface{}); ok {
		applyFieldsToOutput(&s.wrapper.Output, v)
	}
	if v, ok := data["widget"].(map[string]interface{}); ok {
		applyFieldsToWidget(&s.wrapper.Widget, v)
	}
	if v, ok := data["lyric_common"].(map[string]interface{}); ok {
		applyFieldsToLyricCommon(&s.wrapper.LyricCommon, v)
	}
	if v, ok := data["plugins"].(map[string]interface{}); ok {
		if s.wrapper.Plugins == nil {
			s.wrapper.Plugins = map[string]interface{}{}
		}
		for k, val := range v {
			s.wrapper.Plugins[k] = val
		}
	}
	s.save()
}

// GetOutput returns the current output settings (used by output_service).
func (s *SettingsService) GetOutput() models.SettingsOutput {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.wrapper.Output
}

// ReplaceOutputInMemory mirrors app/api/output.py's post_output: it merges
// known fields into the in-memory output settings only, WITHOUT persisting
// to settings.json (Python's post_output never calls settings_service.save()
// for this path - only /api/settings and /api/settings/general do).
func (s *SettingsService) ReplaceOutputInMemory(data map[string]interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	applyFieldsToOutput(&s.wrapper.Output, data)
}

// applyFieldsTo* helpers mirror Python's `if hasattr(obj, k): setattr(obj, k, v)`
// pattern: unknown keys are ignored, known keys are overwritten using a
// JSON round-trip so value types coming from a generic map[string]interface{}
// (as decoded from a request body) land correctly in the typed struct.

func applyFieldsToGeneral(g *models.SettingsGeneral, data map[string]interface{}) {
	mergeJSON(g, data)
}

func applyFieldsToLyric(l *models.SettingsLyric, data map[string]interface{}) {
	mergeJSON(l, data)
}

func applyFieldsToLyricCommon(l *models.SettingsLyricCommon, data map[string]interface{}) {
	mergeJSON(l, data)
}

func applyFieldsToOutput(o *models.SettingsOutput, data map[string]interface{}) {
	mergeJSON(o, data)
}

func applyFieldsToWidget(w *models.SettingsWidget, data map[string]interface{}) {
	mergeJSON(w, data)
}

// mergeJSON merges only the keys present in data into dst by marshalling dst,
// overlaying data's keys onto the resulting map, then unmarshalling back.
// This mirrors "for k, v in data.items(): if hasattr(obj, k): setattr(...)"
// without needing per-struct field enumeration via reflection.
func mergeJSON(dst interface{}, data map[string]interface{}) {
	existing, err := json.Marshal(dst)
	if err != nil {
		return
	}
	var existingMap map[string]interface{}
	if err := json.Unmarshal(existing, &existingMap); err != nil {
		return
	}
	for k, v := range data {
		if _, known := existingMap[k]; known {
			existingMap[k] = v
		}
	}
	merged, err := json.Marshal(existingMap)
	if err != nil {
		return
	}
	_ = json.Unmarshal(merged, dst)
}
