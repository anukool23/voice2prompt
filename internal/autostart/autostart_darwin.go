//go:build darwin

// Package autostart manages launch-at-login via a per-user LaunchAgent plist.
// This needs no extra permission (unlike the System Events / Automation approach).
package autostart

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

const label = "dev.voice2prompt.agent"

func plistPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, "Library", "LaunchAgents", label+".plist"), nil
}

// Enabled reports whether launch-at-login is configured.
func Enabled() bool {
	p, err := plistPath()
	if err != nil {
		return false
	}
	_, err = os.Stat(p)
	return err == nil
}

// SetEnabled turns launch-at-login on or off for the current executable.
func SetEnabled(on bool) error {
	p, err := plistPath()
	if err != nil {
		return err
	}

	if !on {
		_ = exec.Command("launchctl", "unload", p).Run()
		if err := os.Remove(p); err != nil && !os.IsNotExist(err) {
			return err
		}
		return nil
	}

	exe, err := os.Executable()
	if err != nil {
		return err
	}
	plist := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
  <key>Label</key><string>%s</string>
  <key>ProgramArguments</key><array><string>%s</string></array>
  <key>RunAtLoad</key><true/>
  <key>ProcessType</key><string>Interactive</string>
</dict>
</plist>
`, label, exe)

	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(p, []byte(plist), 0o644); err != nil {
		return err
	}
	// Best effort — takes effect at next login regardless.
	_ = exec.Command("launchctl", "load", "-w", p).Run()
	return nil
}
