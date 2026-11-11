package services

import (
	"sync"
	"time"

	"miniplayer-backend/internal/models"
)

// csharpProgressReset marks PrevCSharpProgressSeconds as "no SMTC progress
// applied yet for the current track" (distinct from a real reported value,
// which is always >= 0 since GetSmtcStatus.exe uses -1 for "unknown").
const csharpProgressReset = -2

// NowPlayingService mirrors app/services/nowplaying_service.py: holds the
// current Player/Track state plus the timer bookkeeping (timer_start /
// timer_paused_at) used to derive seekbar progress between polls.
//
// All fields that main.py's on_update() touches directly (PrevStatus,
// PrevWindowTitle, PrevIsPaused, NoneOccursTime, etc.) are exported so the
// detector-driven state machine in internal/api (the Go equivalent of
// on_update) can manipulate them exactly like the Python module-level
// nowplaying_service singleton.
type NowPlayingService struct {
	mu sync.Mutex

	Player models.Player
	Track  models.Track

	timerStart    time.Time
	timerPausedAt *time.Time

	PrevWindowTitle          string
	PrevSmtcTitle            string
	PrevStatus               string
	PrevIsPaused             bool
	NoneOccursTime           float64
	PrevCSharpProgressSeconds int
}

func NewNowPlayingService() *NowPlayingService {
	return &NowPlayingService{
		Player:                    models.NewPlayer(),
		Track:                     models.NewTrack(),
		timerStart:                time.Now(),
		PrevWindowTitle:           "",
		PrevSmtcTitle:             "",
		PrevStatus:                "None",
		PrevIsPaused:              true,
		NoneOccursTime:            0,
		PrevCSharpProgressSeconds: csharpProgressReset,
	}
}

// Lock/Unlock expose the service's mutex so main's on_update-equivalent can
// perform a read-modify-write sequence atomically, mirroring the Python
// code's unsynchronized (but single-threaded asyncio) access pattern safely
// under Go's concurrent scheduler.
func (n *NowPlayingService) Lock()   { n.mu.Lock() }
func (n *NowPlayingService) Unlock() { n.mu.Unlock() }

func (n *NowPlayingService) ResetTimer() {
	n.timerStart = time.Now()
	n.timerPausedAt = nil
}

// SyncTimer rebases the timer so GetTimeMs() reads baselineSec right now,
// then keeps counting up from there - used to align the seekbar with a real
// position SMTC reported, instead of restarting it at 0.
func (n *NowPlayingService) SyncTimer(baselineSec int) {
	n.timerStart = time.Now().Add(-time.Duration(baselineSec) * time.Second)
	n.timerPausedAt = nil
}

func (n *NowPlayingService) StartTimer() {
	if n.timerPausedAt != nil {
		pausedFor := time.Since(*n.timerPausedAt)
		n.timerStart = n.timerStart.Add(pausedFor)
		n.timerPausedAt = nil
	}
}

func (n *NowPlayingService) PauseTimer() {
	if n.timerPausedAt == nil {
		now := time.Now()
		n.timerPausedAt = &now
	}
}

func (n *NowPlayingService) GetTimeMs() int64 {
	if n.timerPausedAt != nil {
		return n.timerPausedAt.Sub(n.timerStart).Milliseconds()
	}
	return time.Since(n.timerStart).Milliseconds()
}

func (n *NowPlayingService) Query() models.Query {
	n.mu.Lock()
	defer n.mu.Unlock()
	return models.Query{Player: n.queryPlayerLocked(), Track: n.Track}
}

func (n *NowPlayingService) QueryPlayer() models.Player {
	n.mu.Lock()
	defer n.mu.Unlock()
	return n.queryPlayerLocked()
}

// queryPlayerLocked mirrors nowplaying_service.query_player; caller must
// hold n.mu.
func (n *NowPlayingService) queryPlayerLocked() models.Player {
	progressMs := n.GetTimeMs()
	progressSec := int(progressMs / 1000)
	duration := n.Track.Duration
	if duration <= 0 {
		duration = 300
	}
	n.Player.SeekbarCurrentPosition = progressSec
	n.Player.SeekbarCurrentPositionHuman = GetFormattedDuration(progressSec)
	if duration > 0 {
		n.Player.StatePercent = float64(progressSec) / float64(duration)
	} else {
		n.Player.StatePercent = 0.0
	}
	return n.Player
}

func (n *NowPlayingService) QueryTrack() models.Track {
	n.mu.Lock()
	defer n.mu.Unlock()
	return n.Track
}

func (n *NowPlayingService) QueryProgress() models.QueryProgress {
	n.mu.Lock()
	defer n.mu.Unlock()
	return models.QueryProgress{ProgressMs: int(n.GetTimeMs())}
}

func (n *NowPlayingService) HasSong() models.RespData {
	n.mu.Lock()
	defer n.mu.Unlock()
	return models.NewRespData(200, "success", n.Player.HasSong)
}

func (n *NowPlayingService) IsConnected() models.RespData {
	n.mu.Lock()
	defer n.mu.Unlock()
	return models.NewRespData(200, "success", n.Player.HasSong)
}
