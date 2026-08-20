package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Highscores holds the best score reached per free-play mode (Normal,
// Randomizer, Endless), persisted to disk so it survives restarts.
type Highscores struct {
	Scores map[string]int `json:"scores"`
}

func highscoresFilePath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	appDir := filepath.Join(dir, "go2048gui")
	if err := os.MkdirAll(appDir, 0o755); err != nil {
		return "", err
	}
	return filepath.Join(appDir, "highscores.json"), nil
}

func loadHighscores() *Highscores {
	h := &Highscores{Scores: map[string]int{}}
	path, err := highscoresFilePath()
	if err != nil {
		return h
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return h
	}
	_ = json.Unmarshal(data, h)
	if h.Scores == nil {
		h.Scores = map[string]int{}
	}
	return h
}

func (h *Highscores) save() {
	path, err := highscoresFilePath()
	if err != nil {
		return
	}
	data, err := json.MarshalIndent(h, "", "  ")
	if err != nil {
		return
	}
	_ = os.WriteFile(path, data, 0o644)
}

func (h *Highscores) get(mode Mode) int {
	return h.Scores[modeName(mode)]
}

// update stores score as the new record for mode if it beats the current
// one, persists immediately, and reports whether it was a new record.
func (h *Highscores) update(mode Mode, score int) bool {
	key := modeName(mode)
	if score > h.Scores[key] {
		h.Scores[key] = score
		h.save()
		return true
	}
	return false
}
