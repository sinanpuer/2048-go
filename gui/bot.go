package main

import (
	"math"
	"math/rand"
	"time"
)

type BotDifficulty int

const (
	BotNormal BotDifficulty = iota
	BotExpert
)

func botDifficultyName(d BotDifficulty) string {
	if d == BotExpert {
		return tr("bot.expert")
	}
	return tr("bot.normal")
}

func botMoveInterval(d BotDifficulty) time.Duration {
	if d == BotExpert {
		return 220 * time.Millisecond
	}
	return 480 * time.Millisecond
}

// botPickMove evaluates every legal move and returns the one the bot wants
// to play. Expert always plays the strongest option by a multi-factor
// heuristic; Normal uses a much weaker one and occasionally moves at random,
// making it clearly beatable.
func botPickMove(b Board, d BotDifficulty) (func(Board) (Board, bool, int), bool) {
	moves := []func(Board) (Board, bool, int){moveUp, moveDown, moveLeft, moveRight}

	type candidate struct {
		idx   int
		score float64
	}
	var candidates []candidate
	for i, mv := range moves {
		nb, ok, _ := mv(b)
		if !ok {
			continue
		}
		var sc float64
		if d == BotExpert {
			sc = evaluateBoard(nb)
		} else {
			sc = evaluateBoardSimple(nb)
		}
		candidates = append(candidates, candidate{i, sc})
	}
	if len(candidates) == 0 {
		return nil, false
	}

	if d == BotNormal && rand.Intn(100) < 35 {
		pick := candidates[rand.Intn(len(candidates))]
		return moves[pick.idx], true
	}

	best := candidates[0]
	for _, c := range candidates[1:] {
		if c.score > best.score {
			best = c
		}
	}
	return moves[best.idx], true
}

func evaluateBoardSimple(b Board) float64 {
	empty := 0
	for r := 0; r < size; r++ {
		for c := 0; c < size; c++ {
			if b[r][c] == 0 {
				empty++
			}
		}
	}
	return float64(empty)
}

// evaluateBoard scores a board using the standard 2048-AI ingredients: free
// space, monotonic rows/columns, smooth neighboring tiles, and a bonus for
// keeping the largest tile pinned in a corner.
func evaluateBoard(b Board) float64 {
	empty := 0.0
	for r := 0; r < size; r++ {
		for c := 0; c < size; c++ {
			if b[r][c] == 0 {
				empty++
			}
		}
	}

	smoothness := 0.0
	for r := 0; r < size; r++ {
		for c := 0; c < size-1; c++ {
			if b[r][c] != 0 && b[r][c+1] != 0 {
				smoothness -= math.Abs(logVal(b[r][c]) - logVal(b[r][c+1]))
			}
		}
	}
	for c := 0; c < size; c++ {
		for r := 0; r < size-1; r++ {
			if b[r][c] != 0 && b[r+1][c] != 0 {
				smoothness -= math.Abs(logVal(b[r][c]) - logVal(b[r+1][c]))
			}
		}
	}

	mono := monotonicityScore(b)

	maxTile := highestTile(b)
	cornerBonus := 0.0
	if b[0][0] == maxTile || b[0][size-1] == maxTile || b[size-1][0] == maxTile || b[size-1][size-1] == maxTile {
		cornerBonus = float64(maxTile)
	}

	return empty*270 + smoothness*10 + mono*47 + cornerBonus*3
}

func monotonicityScore(b Board) float64 {
	score := 0.0
	for r := 0; r < size; r++ {
		incr, decr := 0.0, 0.0
		for c := 0; c < size-1; c++ {
			va, vb := logVal(b[r][c]), logVal(b[r][c+1])
			if va > vb {
				decr += va - vb
			} else {
				incr += vb - va
			}
		}
		score -= math.Min(incr, decr)
	}
	for c := 0; c < size; c++ {
		incr, decr := 0.0, 0.0
		for r := 0; r < size-1; r++ {
			va, vb := logVal(b[r][c]), logVal(b[r+1][c])
			if va > vb {
				decr += va - vb
			} else {
				incr += vb - va
			}
		}
		score -= math.Min(incr, decr)
	}
	return score
}

func logVal(v int) float64 {
	if v <= 0 {
		return 0
	}
	return math.Log2(float64(v))
}
