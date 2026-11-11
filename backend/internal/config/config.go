// Package config resolves runtime settings for the standalone server.
//
// The frontend is embedded into the binary (see internal/staticassets), so
// no static-file directory needs to be resolved on disk. AssetsDir and
// SettingsDir, however, must be real on-disk paths: AssetsDir holds the
// AudioService subprocess (exec'd from disk, so it can't be embedded) and
// any loose asset files, and SettingsDir holds persisted user settings.
// Both are resolved relative to the running executable's own directory, so
// a distribution is just "the exe + a sibling AudioService folder" with no
// repo layout (backend/, frontend/) required around it.
package config

import (
	"os"
	"path/filepath"
)

type Settings struct {
	AppName     string
	ServerPort  int
	AssetsDir   string
	SettingsDir string
}

const defaultPort = 9863

// Load builds Settings, resolving AssetsDir/SettingsDir relative to the
// executable's directory (or MINIPLAY_ROOT, if set, for development/testing
// overrides).
func Load() Settings {
	root := os.Getenv("MINIPLAY_ROOT")
	if root == "" {
		root = exeDir()
	}

	port := defaultPort
	if v := os.Getenv("MINIPLAY_SERVER_PORT"); v != "" {
		if p, err := parsePort(v); err == nil {
			port = p
		}
	}

	return Settings{
		AppName:     "Miniplay",
		ServerPort:  port,
		AssetsDir:   filepath.Join(root, "assets"),
		SettingsDir: filepath.Join(root, "settings"),
	}
}

func parsePort(s string) (int, error) {
	n := 0
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0, os.ErrInvalid
		}
		n = n*10 + int(c-'0')
	}
	return n, nil
}

// exeDir returns the directory containing the running executable, falling
// back to the current working directory if it can't be determined.
func exeDir() string {
	if exe, err := os.Executable(); err == nil {
		if resolved, err := filepath.EvalSymlinks(exe); err == nil {
			return filepath.Dir(resolved)
		}
		return filepath.Dir(exe)
	}
	if wd, err := os.Getwd(); err == nil {
		return wd
	}
	return "."
}
