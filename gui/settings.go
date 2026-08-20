package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Settings holds user-configurable presentation options, persisted to disk
// so they survive restarts.
type Settings struct {
	Animations bool   `json:"animations"`
	Theme      string `json:"theme"`
	FPS        int    `json:"fps"` // 30, 60, or 0 = unlimited
}

func defaultSettings() *Settings {
	return &Settings{Animations: true, Theme: ThemeClassic, FPS: 60}
}

func settingsFilePath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	appDir := filepath.Join(dir, "go2048gui")
	if err := os.MkdirAll(appDir, 0o755); err != nil {
		return "", err
	}
	return filepath.Join(appDir, "settings.json"), nil
}

func loadSettings() *Settings {
	s := defaultSettings()
	path, err := settingsFilePath()
	if err != nil {
		return s
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return s
	}
	_ = json.Unmarshal(data, s)
	if s.Theme != ThemeClassic && s.Theme != ThemeStone && s.Theme != ThemeCandy {
		s.Theme = ThemeClassic
	}
	if s.FPS != 30 && s.FPS != 60 && s.FPS != 0 {
		s.FPS = 60
	}
	return s
}

func (s *Settings) save() {
	path, err := settingsFilePath()
	if err != nil {
		return
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return
	}
	_ = os.WriteFile(path, data, 0o644)
}

// ambientTickInterval derives the menu background's redraw cadence from the
// FPS setting, in milliseconds: a lower setting means less frequent
// redraw work, which matters on weaker machines.
func (s *Settings) ambientTickIntervalMs() int {
	switch s.FPS {
	case 30:
		return 900
	case 0:
		return 400
	default:
		return 650
	}
}
