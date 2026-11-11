// Package config mirrors backend/app/core/config.py.
//
// Directory layout assumptions replicate the Python Settings class, which
// derives paths relative to the location of config.py inside backend/app/core:
//   assets_dir   = backend/assets
//   settings_dir = backend/settings
//   static_dir   = frontend/dist
//
// The Go binary reuses the *existing* backend/assets and backend/settings
// directories (rather than duplicating the large AudioService exe payload),
// so these paths are resolved relative to the repository root, which is
// discovered by walking up from the executable's working directory /
// executable location looking for a `backend` directory sibling.
package config

import (
	"os"
	"path/filepath"
)

type Settings struct {
	AppName    string
	ServerPort int
	AssetsDir  string
	SettingsDir string
	StaticDir  string
}

const defaultPort = 9863

// Load builds Settings, locating the repository root so backend_go can reuse
// backend/assets and backend/settings without copying them, and serve
// frontend/dist as before. Root discovery order:
//  1. MINIPLAY_ROOT env var, if set.
//  2. Walk upward from the executable's directory looking for a `backend`
//     subdirectory (handles running `go run`/the built exe from anywhere).
//  3. Walk upward from the current working directory as a fallback.
func Load() Settings {
	root := os.Getenv("MINIPLAY_ROOT")
	if root == "" {
		root = findRepoRoot()
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
		AssetsDir:   filepath.Join(root, "backend", "assets"),
		SettingsDir: filepath.Join(root, "backend", "settings"),
		StaticDir:   filepath.Join(root, "frontend", "dist"),
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

// findRepoRoot walks upward from both the executable location and the
// current working directory, returning the first ancestor that contains a
// `backend` directory. Falls back to the current working directory.
func findRepoRoot() string {
	candidates := []string{}
	if exe, err := os.Executable(); err == nil {
		candidates = append(candidates, filepath.Dir(exe))
	}
	if wd, err := os.Getwd(); err == nil {
		candidates = append(candidates, wd)
	}

	for _, start := range candidates {
		dir := start
		for i := 0; i < 8; i++ {
			if info, err := os.Stat(filepath.Join(dir, "backend")); err == nil && info.IsDir() {
				return dir
			}
			parent := filepath.Dir(dir)
			if parent == dir {
				break
			}
			dir = parent
		}
	}

	if wd, err := os.Getwd(); err == nil {
		return wd
	}
	return "."
}
