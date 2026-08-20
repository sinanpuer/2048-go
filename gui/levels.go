package main

import (
	"encoding/json"
	"math"
	"os"
	"path/filepath"
)

const totalLevels = 100

type LevelKind int

const (
	KindScore LevelKind = iota
	KindTile
	KindTimedScore
	KindTimedTile
)

type LevelDef struct {
	Number    int
	Kind      LevelKind
	ScoreGoal int
	TileGoal  int
	TimeLimit int // seconds; 0 means untimed
}

func (l LevelDef) Title() string {
	switch l.Kind {
	case KindScore:
		return trf("level.goalScore", l.ScoreGoal)
	case KindTile:
		return trf("level.goalTile", l.TileGoal)
	case KindTimedScore:
		return trf("level.goalTimedScore", l.ScoreGoal, l.TimeLimit)
	case KindTimedTile:
		return trf("level.goalTimedTile", l.TileGoal, l.TimeLimit)
	}
	return ""
}

// generateLevels builds a difficulty curve of n levels, cycling through the
// four level kinds. Level 50 and the final level are pinned to the 2048 and
// 4096 milestones requested explicitly, overriding whatever the cycle gives.
func generateLevels(n int) []LevelDef {
	levels := make([]LevelDef, n)
	for i := 1; i <= n; i++ {
		kind := LevelKind((i - 1) % 4)
		def := LevelDef{Number: i, Kind: kind, TimeLimit: 0}

		switch kind {
		case KindScore:
			def.ScoreGoal = scoreGoalForLevel(i)
		case KindTile:
			def.TileGoal = tileGoalForLevel(i)
		case KindTimedScore:
			def.ScoreGoal = timedScoreGoalForLevel(i)
			def.TimeLimit = timeLimitForLevel(i)
		case KindTimedTile:
			def.TileGoal = timedTileGoalForLevel(i)
			def.TimeLimit = timeLimitForLevel(i)
		}

		levels[i-1] = def
	}

	if n >= 50 {
		levels[49] = LevelDef{Number: 50, Kind: KindTile, TileGoal: 2048}
	}
	if n >= 1 {
		levels[n-1] = LevelDef{Number: n, Kind: KindTile, TileGoal: 4096}
	}

	return levels
}

// linearGoal interpolates a value between start (at level 1) and end (at
// the final level, out of total levels) in a straight line, so difficulty
// ramps up steadily instead of exploding late-game.
func linearGoal(level, total int, start, end float64) float64 {
	if total <= 1 {
		return start
	}
	t := float64(level-1) / float64(total-1)
	return start + (end-start)*t
}

func scoreGoalForLevel(level int) int {
	return roundTo(int(linearGoal(level, totalLevels, 60, 9200)), 10)
}

func tileGoalForLevel(level int) int {
	var exp float64
	if level <= 50 {
		t := float64(level-1) / 49.0
		exp = 3 + 8*math.Pow(t, 1.15)
	} else {
		t := float64(level-50) / 50.0
		exp = 11 + t*t
	}
	e := int(math.Round(exp))
	if e < 1 {
		e = 1
	}
	return 1 << uint(e)
}

func timeLimitForLevel(level int) int {
	t := 90 - int(float64(level)/100.0*60)
	if t < 30 {
		t = 30
	}
	return t
}

func timedScoreGoalForLevel(level int) int {
	return roundTo(int(linearGoal(level, totalLevels, 30, 1400)), 5)
}

func timedTileGoalForLevel(level int) int {
	exp := 4 + 5*float64(level-1)/99.0
	e := int(math.Round(exp))
	if e < 1 {
		e = 1
	}
	return 1 << uint(e)
}

func roundTo(v, step int) int {
	if step <= 0 {
		return v
	}
	return ((v + step/2) / step) * step
}

// ---------- progress persistence ----------

type Progress struct {
	Completed       map[int]bool `json:"completed"`
	PuzzleCompleted map[int]bool `json:"puzzleCompleted"`
}

func progressFilePath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	appDir := filepath.Join(dir, "go2048gui")
	if err := os.MkdirAll(appDir, 0o755); err != nil {
		return "", err
	}
	return filepath.Join(appDir, "progress.json"), nil
}

func loadProgress() *Progress {
	p := &Progress{Completed: map[int]bool{}, PuzzleCompleted: map[int]bool{}}
	path, err := progressFilePath()
	if err != nil {
		return p
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return p
	}
	_ = json.Unmarshal(data, p)
	if p.Completed == nil {
		p.Completed = map[int]bool{}
	}
	if p.PuzzleCompleted == nil {
		p.PuzzleCompleted = map[int]bool{}
	}
	return p
}

func (p *Progress) save() {
	path, err := progressFilePath()
	if err != nil {
		return
	}
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return
	}
	_ = os.WriteFile(path, data, 0o644)
}

func (p *Progress) isUnlocked(levelNumber int) bool {
	if levelNumber <= 1 {
		return true
	}
	return p.Completed[levelNumber-1]
}

func (p *Progress) markCompleted(levelNumber int) {
	p.Completed[levelNumber] = true
}

func (p *Progress) isPuzzleUnlocked(levelNumber int) bool {
	if levelNumber <= 1 {
		return true
	}
	return p.PuzzleCompleted[levelNumber-1]
}

func (p *Progress) markPuzzleCompleted(levelNumber int) {
	p.PuzzleCompleted[levelNumber] = true
}
