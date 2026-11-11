package services

import (
	"time"

	"miniplayer-backend/internal/detectors"
	"miniplayer-backend/internal/models"
)

// AppUpdater wires the AudioDetector's callback to NowPlayingService state
// transitions, mirroring main.py's on_update() exactly:
//   - a 3-second "None" grace period before resetting Player/Track
//   - window-title-change detection triggering timer reset/start + output
//   - pause/resume timer handling otherwise
//   - the 300s / "5:00" duration fallback when no duration is known
type AppUpdater struct {
	detector *detectors.AudioDetector
	nowPlaying *NowPlayingService
	output   *OutputService
}

func NewAppUpdater(detector *detectors.AudioDetector, nowPlaying *NowPlayingService, output *OutputService) *AppUpdater {
	return &AppUpdater{detector: detector, nowPlaying: nowPlaying, output: output}
}

// OnUpdate mirrors main.py's on_update() coroutine. It is invoked
// synchronously from the detector's read loop goroutine after each parsed
// stdout line, exactly like Python's `await self._callback()`.
func (a *AppUpdater) OnUpdate() {
	status := a.detector.Status()
	windowTitle := a.detector.WindowTitle()
	smtcTitle := a.detector.SmtcTitle()
	smtcArtist := a.detector.SmtcArtist()
	smtcAlbum := a.detector.SmtcAlbum()
	thumbnail := a.detector.Thumbnail()
	smtcProgressSec, smtcTotalSec := a.detector.Progress()

	n := a.nowPlaying
	n.Lock()
	defer n.Unlock()

	now := float64(time.Now().UnixNano()) / 1e9

	if status == "None" {
		if n.PrevStatus == "None" {
			if now-n.NoneOccursTime > 3 {
				n.Player = models.NewPlayer()
				n.Track = models.NewTrack()
				n.ResetTimer()
				n.PrevWindowTitle = ""
				n.PrevSmtcTitle = ""
				n.PrevIsPaused = true
				n.PrevCSharpProgressSeconds = csharpProgressReset
			}
		} else {
			n.NoneOccursTime = now
		}
		n.PrevStatus = status
		return
	}

	n.Player.HasSong = true
	isPaused := status != "Playing"
	n.Player.IsPaused = isPaused

	identityChanged := (windowTitle != "" && n.PrevWindowTitle != windowTitle) ||
		smtcTitle != n.PrevSmtcTitle

	if identityChanged {
		// SMTC already knows the real position/duration of the new track;
		// sync to that instead of restarting the seekbar at 0 and waiting
		// for Go's own clock to catch up.
		if smtcProgressSec >= 0 {
			n.SyncTimer(smtcProgressSec)
		} else {
			n.ResetTimer()
		}
		n.StartTimer()
		n.PrevCSharpProgressSeconds = smtcProgressSec

		if smtcTitle != "" {
			n.Track.Title = smtcTitle
			n.Track.Author = smtcArtist
			n.Track.Album = smtcAlbum
		} else {
			title, author, hasAuthor := ParseWindowTitle(windowTitle)
			if title != "" {
				n.Track.Title = title
			} else {
				n.Track.Title = windowTitle
			}
			if hasAuthor {
				n.Track.Author = author
			} else {
				n.Track.Author = ""
			}
			n.Track.Album = ""
		}
		n.Track.Cover = thumbnail
		if smtcTotalSec > 0 {
			n.Track.Duration = smtcTotalSec
			n.Track.DurationHuman = GetFormattedDuration(smtcTotalSec)
		} else if n.Track.Duration <= 0 {
			n.Track.Duration = 300
			n.Track.DurationHuman = "5:00"
		}

		trackSnapshot := n.Track
		if a.output != nil {
			go a.output.OutputAsync(trackSnapshot)
		}
	} else {
		if thumbnail != n.Track.Cover {
			n.Track.Cover = thumbnail
		}
		if smtcTitle != "" {
			if smtcArtist != n.Track.Author {
				n.Track.Author = smtcArtist
			}
			if smtcAlbum != n.Track.Album {
				n.Track.Album = smtcAlbum
			}
		}
		if smtcTotalSec > 0 && smtcTotalSec != n.Track.Duration {
			n.Track.Duration = smtcTotalSec
			n.Track.DurationHuman = GetFormattedDuration(smtcTotalSec)
		}
		// Resync whenever SMTC reports a progress value we haven't already
		// applied - catches seeks, and keeps drift from accumulating between
		// a session's own position updates and Go's interpolated clock.
		if smtcProgressSec >= 0 && smtcProgressSec != n.PrevCSharpProgressSeconds {
			n.PrevCSharpProgressSeconds = smtcProgressSec
			n.SyncTimer(smtcProgressSec)
		}
		if isPaused {
			n.PauseTimer()
		} else {
			n.StartTimer()
		}
	}

	n.PrevWindowTitle = windowTitle
	n.PrevSmtcTitle = smtcTitle
	n.PrevStatus = status
	n.PrevIsPaused = isPaused
}
