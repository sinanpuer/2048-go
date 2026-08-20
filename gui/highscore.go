package main

import (
	"encoding/json"
	"fmt"
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

// highscoreKey builds the storage key for a mode+board size combination.
// Board size 4 keeps the plain legacy key (just the mode name) so upgrading
// from before board sizes existed doesn't lose anyone's existing 4x4 record.
func highscoreKey(mode Mode, boardSize int) string {
	if boardSize == 4 {
		return modeName(mode)
	}
	return fmt.Sprintf("%s_%dx%d", modeName(mode), boardSize, boardSize)
}

func (h *Highscores) get(mode Mode, boardSize int) int {
	return h.Scores[highscoreKey(mode, boardSize)]
}

// update stores score as the new record for mode+boardSize if it beats the
// current one, persists immediately, and reports whether it was a new record.
func (h *Highscores) update(mode Mode, score int, boardSize int) bool {
	key := highscoreKey(mode, boardSize)
	if score > h.Scores[key] {
		h.Scores[key] = score
		h.save()
		return true
	}
	return false
}
