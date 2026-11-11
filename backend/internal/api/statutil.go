package api

import "os"

// statExists is a tiny os.Stat wrapper shared by page.go/widget.go fallback
// logic, kept as a named helper for readability at call sites.
func statExists(path string) (os.FileInfo, error) {
	return os.Stat(path)
}
