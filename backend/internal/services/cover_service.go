package services

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"miniplayer-backend/internal/models"
)

// CoverService mirrors app/services/cover_service.py.
type CoverService struct {
	assetsDir string

	mu             sync.Mutex
	prevSongTitle  string
	prevSongAuthor string
	videoURLCached bool
	cachedVideoURL interface{}
	havePrevTitle  bool
}

func NewCoverService(assetsDir string) *CoverService {
	return &CoverService{assetsDir: assetsDir}
}

// ConvertToBase64 mirrors cover_service.convert_to_base64: fetches the given
// cover URL and returns it as a data: URI, falling back to the bundled
// no_cover_base64.txt asset (if present) on any failure, exactly like Python.
func (c *CoverService) ConvertToBase64(coverURL string) models.Base64Img {
	if coverURL == "" {
		return c.readDefaultCover()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, coverURL, nil)
	if err == nil {
		resp, err := http.DefaultClient.Do(req)
		if err == nil {
			defer resp.Body.Close()
			if resp.StatusCode == 200 {
				data, readErr := io.ReadAll(resp.Body)
				if readErr == nil {
					contentType := resp.Header.Get("Content-Type")
					if contentType == "" {
						contentType = "image/jpeg"
					}
					b64 := base64.StdEncoding.EncodeToString(data)
					return models.Base64Img{Base64Img: fmt.Sprintf("data:%s;base64,%s", contentType, b64)}
				}
			}
		}
	}

	return c.readDefaultCover()
}

func (c *CoverService) readDefaultCover() models.Base64Img {
	defaultPath := filepath.Join(c.assetsDir, "no_cover_base64.txt")
	if data, err := os.ReadFile(defaultPath); err == nil {
		return models.Base64Img{Base64Img: strings.TrimSpace(string(data))}
	}
	return models.Base64Img{Base64Img: ""}
}

// GetVideoURL mirrors cover_service.get_video_url: a cache keyed on
// (title, author) that always resolves to nil/None since no video source is
// implemented upstream either.
func (c *CoverService) GetVideoURL(songTitle, songAuthor string) interface{} {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.videoURLCached && c.havePrevTitle && c.prevSongTitle == songTitle && c.prevSongAuthor == songAuthor {
		return c.cachedVideoURL
	}
	c.prevSongTitle = songTitle
	c.havePrevTitle = true
	c.prevSongAuthor = songAuthor
	c.cachedVideoURL = nil
	c.videoURLCached = true
	return nil
}
