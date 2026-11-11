// Package detectors mirrors backend/app/detectors: it shells out to the
// prebuilt .NET helper GetMusicStatus.exe and parses its line-oriented
// stdout protocol exactly like the Python AudioDetector did.
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

// Start launches GetMusicStatus.exe and begins reading its stdout. callback
// is invoked after every parsed line, same as Python's await self._callback().
func (d *AudioDetector) Start(callback func()) {
	exePath := filepath.Join(d.assetsDir, "AudioService", "GetMusicStatus.exe")
	if _, err := os.Stat(exePath); err != nil {
		log.Printf("GetMusicStatus.exe not found at %s", exePath)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	d.cancel = cancel

	cmd := exec.CommandContext(ctx, exePath,
		"--device-id", "default",
		"--platform", "default",
		"--smtc", "true",
		"--poll-interval", "1000",
		"--splayer-next-port", "37685",
	)
	cmd.Dir = filepath.Join(d.assetsDir, "AudioService")

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

	log.Println("Started GetMusicStatus.exe")
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
