// Package resource resolves data/binary files (whisper model, whisper-server) so
// the same relative paths work whether we're running from the repo (dev) or from
// inside a packaged Voice2Prompt.app bundle (Contents/Resources).
package resource

import (
	"os"
	"path/filepath"
)

// searchDirs are the directories we look in, in priority order: current working
// directory, the executable's directory, and (for a .app bundle) Contents/Resources.
func searchDirs() []string {
	dirs := []string{"."}
	if exe, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exe)
		dirs = append(dirs,
			exeDir,                                   // next to the binary
			filepath.Join(exeDir, "..", "Resources"), // .app bundle resources
		)
	}
	return dirs
}

// Find resolves relPath to an existing file. It tries relPath against each search
// dir, then falls back to matching just the basename in each search dir (so a dev
// path like third_party/.../whisper-server maps to Resources/whisper-server in a
// bundle). Returns (resolved, true) if found, else (relPath, false).
func Find(relPath string) (string, bool) {
	if relPath == "" {
		return relPath, false
	}
	if filepath.IsAbs(relPath) {
		_, err := os.Stat(relPath)
		return relPath, err == nil
	}
	for _, d := range searchDirs() {
		cand := filepath.Join(d, relPath)
		if _, err := os.Stat(cand); err == nil {
			return cand, true
		}
	}
	base := filepath.Base(relPath)
	for _, d := range searchDirs() {
		cand := filepath.Join(d, base)
		if _, err := os.Stat(cand); err == nil {
			return cand, true
		}
	}
	return relPath, false
}
