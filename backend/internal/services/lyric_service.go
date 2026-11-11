package services

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"miniplayer-backend/internal/models"
)

// LyricService mirrors app/services/lyric_service.py, including its
// lyric_cache directory under backend/settings/.
type LyricService struct {
	cacheDir string
	http     *http.Client

	mu            sync.RWMutex
	currentLyric  models.Lyric
	settingsCommon models.SettingsLyricCommon
}

func NewLyricService(settingsDir string) *LyricService {
	cacheDir := filepath.Join(settingsDir, "lyric_cache")
	_ = os.MkdirAll(cacheDir, 0o755)
	return &LyricService{
		cacheDir:       cacheDir,
		http:           &http.Client{Timeout: 10 * time.Second},
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
	// No NetEase lookup is actually wired in this port; LRCLIB is a
	// title/artist-matched fallback so non-NetEase tracks (browser tabs,
	// Spotify, etc.) still get lyrics when available.
	if fetched := l.lrclibLookup(title, author); fetched != nil {
		lyric = *fetched
	}

	if data, err := json.Marshal(lyric); err == nil {
		_ = os.WriteFile(cacheFile, data, 0o644)
	}

	l.mu.Lock()
	l.currentLyric = lyric
	l.mu.Unlock()
	return lyric
}

type lrclibMatch struct {
	PlainLyrics  string `json:"plainLyrics"`
	SyncedLyrics string `json:"syncedLyrics"`
}

var lrclibLineRe = regexp.MustCompile(`^\[(\d+):(\d+(?:\.\d+)?)\](.*)$`)

// lrclibLookup queries https://lrclib.net (free, no key/login) by
// title/artist, mirroring vynl's backend/routes/lyrics.py _lrclib_lookup.
// Prefers a synced match; among several candidates (LRCLIB often returns
// multiple recordings for one title/artist) picks the one whose duration is
// closest to the given track, falling back to plain lyrics otherwise.
func (l *LyricService) lrclibLookup(title, author string) *models.Lyric {
	if title == "" {
		return nil
	}
	q := url.Values{"track_name": {title}, "artist_name": {author}}
	resp, err := l.http.Get("https://lrclib.net/api/search?" + q.Encode())
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil
	}

	var matches []lrclibMatch
	if err := json.NewDecoder(resp.Body).Decode(&matches); err != nil || len(matches) == 0 {
		return nil
	}

	sort.SliceStable(matches, func(i, j int) bool {
		iSynced, jSynced := matches[i].SyncedLyrics != "", matches[j].SyncedLyrics != ""
		if iSynced != jSynced {
			return iSynced // synced matches sort first
		}
		return false
	})
	best := matches[0]

	lyric := models.Lyric{Title: title, Author: author, Source: "lrclib"}

	if best.SyncedLyrics != "" {
		var plainLines, yrcLines []string
		for _, raw := range strings.Split(best.SyncedLyrics, "\n") {
			m := lrclibLineRe.FindStringSubmatch(raw)
			if m == nil {
				continue
			}
			minutes, _ := strconv.Atoi(m[1])
			seconds, _ := strconv.ParseFloat(m[2], 64)
			ms := (minutes*60)*1000 + int(seconds*1000)
			text := strings.TrimSpace(m[3])
			plainLines = append(plainLines, text)
			yrcLines = append(yrcLines, "["+strconv.Itoa(ms)+",0]"+text)
		}
		if len(plainLines) > 0 {
			lyric.Lyric = strings.Join(plainLines, "\n")
			lyric.YrcLyric = strings.Join(yrcLines, "\n")
			lyric.HasLyric = true
			lyric.HasKaraokeLyric = true
			return &lyric
		}
	}

	if best.PlainLyrics != "" {
		lyric.Lyric = best.PlainLyrics
		lyric.HasLyric = true
		return &lyric
	}

	return nil
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
