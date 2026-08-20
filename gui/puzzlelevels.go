package main

import (
	"math/rand"
)

const totalPuzzleLevels = 100

// PuzzleLevelDef describes a fixed starting board that the player must turn
// into a target score within a limited number of moves.
type PuzzleLevelDef struct {
	Number     int
	StartBoard Board
	MoveLimit  int
	ScoreGoal  int
}

func (p PuzzleLevelDef) Title() string {
	return trf("puzzle.goal", p.ScoreGoal, p.MoveLimit)
}

// generatePuzzleLevels builds n puzzles with a linearly increasing
// difficulty: the starting board gets fuller and its tiles bigger, the
// score goal climbs, and the move budget gets comparatively tighter.
// Each level's starting board is derived from a seed tied to its number,
// so replaying a level always presents the exact same puzzle.
func generatePuzzleLevels(n int) []PuzzleLevelDef {
	levels := make([]PuzzleLevelDef, n)
	for i := 1; i <= n; i++ {
		levels[i-1] = generatePuzzleLevel(i, n)
	}
	return levels
}

func generatePuzzleLevel(level, total int) PuzzleLevelDef {
	rng := rand.New(rand.NewSource(int64(level)*7919 + 104729))

	t := linearGoal(level, total, 0, 1)

	tileCount := 2 + int(t*7) // 2 starting tiles (easy) up to 9 (hard)
	maxExp := 1 + int(t*7)    // tile values up to 2^2 .. 2^8 (4 .. 256)
	moveLimit := int(linearGoal(level, total, 18, 38))
	scoreGoal := roundTo(int(linearGoal(level, total, 30, 6000)), 10)

	var board Board
	placed, attempts := 0, 0
	for placed < tileCount && attempts < 200 {
		attempts++
		r, c := rng.Intn(size), rng.Intn(size)
		if board[r][c] != 0 {
			continue
		}
		exp := 1 + rng.Intn(maxExp)
		board[r][c] = 1 << uint(exp)
		placed++
	}

	return PuzzleLevelDef{Number: level, StartBoard: board, MoveLimit: moveLimit, ScoreGoal: scoreGoal}
}
