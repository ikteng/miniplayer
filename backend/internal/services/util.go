package services

import (
	"fmt"
	"strings"
)

// GetFormattedDuration mirrors app/utils/timeutil.py get_formatted_duration.
func GetFormattedDuration(seconds int) string {
	if seconds < 0 {
		seconds = 0
	}
	m := seconds / 60
	s := seconds % 60
	return fmt.Sprintf("%d:%02d", m, s)
}

// ParseWindowTitle mirrors app/utils/songutil.py parse_window_title.
// Returns (title, author) where author may be empty if not present.
// Python: splits on " - " -> [title, author] = [parts[1], parts[0]]
// (i.e. "Artist - Title" => title=parts[1], author=parts[0]).
func ParseWindowTitle(windowTitle string) (title string, author string, hasAuthor bool) {
	if windowTitle == "" {
		return "", "", false
	}
	if idx := strings.Index(windowTitle, " - "); idx != -1 {
		left := strings.TrimSpace(windowTitle[:idx])
		right := strings.TrimSpace(windowTitle[idx+3:])
		return right, left, true
	}
	return strings.TrimSpace(windowTitle), "", false
}
