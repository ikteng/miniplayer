package detectors

import "strings"

// AumidPlatformMap mirrors app/detectors/platform_detector.py's
// AUMID_PLATFORM_MAP: maps a session's source app user model ID to our
// platform id. Browser AUMIDs are intentionally left unmapped.
var AumidPlatformMap = map[string]string{
	"spotify.exe":         "spotify",
	"tencent.qqmusic":     "qq",
	"kugou.exe":           "kugou",
	"kuwo.exe":            "kuwo",
	"cloudmusic.exe":      "netease",
	"cider":               "cider",
	"potplayermini64.exe": "potplayer",
	"potplayer.exe":       "potplayer",
	"foobar2000.exe":      "foobar",
	"aimp.exe":            "aimp",
}

// DetectPlatformResult mirrors the dict returned by detect_platform().
type DetectPlatformResult struct {
	Platform *string `json:"platform"`
	Reason   string  `json:"reason,omitempty"`
	Aumid    string  `json:"aumid,omitempty"`
	Title    string  `json:"title,omitempty"`
}

// DetectPlatform mirrors app/detectors/platform_detector.py's detect_platform().
//
// The Python implementation used the `winsdk` GlobalSystemMediaTransportControlsSessionManager
// bindings to enumerate SMTC sessions directly, independent of the
// GetMusicStatus.exe polling loop. Go has no equivalent first-party SMTC
// binding in the standard library, and pulling in a CGo/WinRT binding
// contradicts the pure-Go build goal of this port. Since the main polling
// loop already gets window title/status from GetMusicStatus.exe (which
// itself uses SMTC/WindowsMediaController internally - see
// backend/assets/AudioService/WindowsMediaController.dll), this endpoint is
// kept as a best-effort stub that reports SMTC session enumeration as
// unavailable, exactly matching Python's "winsdk not available" branch when
// winsdk cannot be imported. This is the one place behavior cannot be fully
// reproduced without a native WinRT binding; see README for details.
func DetectPlatform() DetectPlatformResult {
	return DetectPlatformResult{Platform: nil, Reason: "winsdk not available"}
}

// MatchAumid is a helper retained for parity with the Python map lookup,
// useful if a future native SMTC binding is wired in.
func MatchAumid(aumid string) (string, bool) {
	p, ok := AumidPlatformMap[strings.ToLower(aumid)]
	return p, ok
}
