// Package detectors shells out to GetSmtcStatus.exe, a generic detector
// built on the Windows SMTC API (GlobalSystemMediaTransportControlsSessionManager)
// that reports whatever app currently holds the active media session —
// browsers, Spotify, VLC, NetEase, etc. — unlike the prebuilt
// GetMusicStatus.exe, which only recognizes NetEase Cloud Music via UI
// Automation. Both emit the same line-oriented stdout protocol.
package detectors

import (
	"bufio"
	"context"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

// AudioDetector polls GetMusicStatus.exe and reports status/window-title
// changes via a callback, mirroring app/detectors/audio_detector.py.
type AudioDetector struct {
	assetsDir string

	mu            sync.RWMutex
	status        string
	windowTitle   string
	progressSec   int
	totalSec      int
	smtcTitle     string
	smtcArtist    string
	smtcAlbum     string
	thumbnail     string

	cmd    *exec.Cmd
	cancel context.CancelFunc
	done   chan struct{}
}

func NewAudioDetector(assetsDir string) *AudioDetector {
	return &AudioDetector{
		assetsDir:   assetsDir,
		status:      "None",
		windowTitle: "",
		progressSec: -1,
		totalSec:    -1,
	}
}

func (d *AudioDetector) Status() string {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.status
}

func (d *AudioDetector) WindowTitle() string {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.windowTitle
}

// SmtcTitle, SmtcArtist, SmtcAlbum, and Thumbnail return real metadata
// reported by SMTC itself (GlobalSystemMediaTransportControlsSessionMediaProperties),
// as opposed to WindowTitle's heuristic parse of the source app's window title.
func (d *AudioDetector) SmtcTitle() string {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.smtcTitle
}

func (d *AudioDetector) SmtcArtist() string {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.smtcArtist
}

func (d *AudioDetector) SmtcAlbum() string {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.smtcAlbum
}

// Thumbnail returns the current track's cover art as a data: URI
// (e.g. "data:image/png;base64,...."), or "" if none is available.
func (d *AudioDetector) Thumbnail() string {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.thumbnail
}

// Start launches GetSmtcStatus.exe and begins reading its stdout. callback
// is invoked after every parsed line, same as Python's await self._callback().
func (d *AudioDetector) Start(callback func()) {
	audioServiceDir := filepath.Join(d.assetsDir, "AudioService")
	exePath := filepath.Join(audioServiceDir, "GetSmtcStatus", "GetSmtcStatus.exe")
	cmdDir := filepath.Join(audioServiceDir, "GetSmtcStatus")
	if _, err := os.Stat(exePath); err != nil {
		// Fall back to the NetEase-only detector if the generic one isn't present.
		exePath = filepath.Join(audioServiceDir, "GetMusicStatus.exe")
		cmdDir = audioServiceDir
		if _, err := os.Stat(exePath); err != nil {
			log.Printf("No audio detector executable found under %s", audioServiceDir)
			return
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	d.cancel = cancel

	cmd := exec.CommandContext(ctx, exePath, "--poll-interval", "1000")
	cmd.Dir = cmdDir

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("Failed to start detector: %v", err)
		return
	}
	// Drain stderr so the pipe never blocks the process.
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		log.Printf("Failed to start detector: %v", err)
		return
	}
	d.cmd = cmd
	d.done = make(chan struct{})

	if stderr != nil {
		go func() { _, _ = io.Copy(io.Discard, stderr) }()
	}

	log.Printf("Started %s", filepath.Base(exePath))
	go d.readLoop(stdout, callback)
}

func (d *AudioDetector) readLoop(stdout io.Reader, callback func()) {
	defer close(d.done)
	scanner := bufio.NewScanner(stdout)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for scanner.Scan() {
		text := strings.TrimRight(scanner.Text(), "\r\n")
		text = strings.TrimSpace(text)
		if text == "" {
			continue
		}
		d.handleLine(text)
		if callback != nil {
			callback()
		}
	}
}

func (d *AudioDetector) handleLine(text string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	switch text {
	case "Playing", "Paused", "None":
		d.status = text
		if text == "None" {
			d.windowTitle = ""
			d.progressSec = -1
			d.totalSec = -1
			d.smtcTitle = ""
			d.smtcArtist = ""
			d.smtcAlbum = ""
			d.thumbnail = ""
		}
		return
	}

	if strings.HasPrefix(text, "Progress:") {
		parts := strings.Split(text[9:], "|")
		if len(parts) == 2 {
			p, errP := strconv.Atoi(strings.TrimSpace(parts[0]))
			t, errT := strconv.Atoi(strings.TrimSpace(parts[1]))
			if errP == nil && errT == nil {
				d.progressSec = p
				d.totalSec = t
			}
		}
		return
	}

	if strings.HasPrefix(text, "Title:") {
		d.smtcTitle = text[len("Title:"):]
		return
	}
	if strings.HasPrefix(text, "Artist:") {
		d.smtcArtist = text[len("Artist:"):]
		return
	}
	if strings.HasPrefix(text, "Album:") {
		d.smtcAlbum = text[len("Album:"):]
		return
	}
	if strings.HasPrefix(text, "Thumbnail:") {
		d.thumbnail = text[len("Thumbnail:"):]
		return
	}

	newTitle := text
	if newTitle != d.windowTitle {
		d.progressSec = -1
		d.totalSec = -1
	}
	d.windowTitle = newTitle
}

// Stop terminates GetMusicStatus.exe, mirroring Python's terminate-then-kill
// fallback with a 2s grace period.
func (d *AudioDetector) Stop() {
	if d.cmd == nil || d.cmd.Process == nil {
		return
	}
	if d.cancel != nil {
		d.cancel()
	}
	done := make(chan struct{})
	go func() {
		_ = d.cmd.Wait()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		_ = d.cmd.Process.Kill()
	}
}
