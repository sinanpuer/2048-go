// Board/move logic kept in sync by hand with gui/game.go - duplicated here
// (trimmed of Fyne-only bits like tileColor) so this module has no
// dependency on the GUI package and can be deployed standalone.
package main

import "math/rand"

const size = 4

// Board is a square grid of tiles.
type Board [][]int

// newEmptyBoard allocates an n x n grid of zeroed cells.
func newEmptyBoard(n int) Board {
	b := make(Board, n)
	for i := range b {
		b[i] = make([]int, n)
	}
	return b
}

// boardsEqual reports whether two boards have the same dimensions and cell
// values. Board is slice-based, so it can't be compared with == or !=.
func boardsEqual(a, b Board) bool {
	if len(a) != len(b) {
		return false
	}
	for r := range a {
		if len(a[r]) != len(b[r]) {
			return false
		}
		for c := range a[r] {
			if a[r][c] != b[r][c] {
				return false
			}
		}
	}
	return true
}

func equalRows(a, b []int) bool {
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

type Mode int

const (
	ModeNormal Mode = iota
	ModeRandomizer
	ModeEndless
)

func highestTile(b Board) int {
	max := 0
	for r := range b {
		for c := range b[r] {
			if b[r][c] > max {
				max = b[r][c]
			}
		}
	}
	return max
}

func randomTileValue(mode Mode, b Board) int {
	if mode != ModeRandomizer {
		if rand.Intn(10) == 0 {
			return 4
		}
		return 2
	}

	highest := highestTile(b)
	roll := rand.Intn(100)

	switch {
	case highest < 32:
		if roll < 80 {
			return 2
		}
		return 4
	case highest < 128:
		switch {
		case roll < 60:
			return 2
		case roll < 90:
			return 4
		default:
			return 8
		}
	case highest < 512:
		switch {
		case roll < 50:
			return 2
		case roll < 80:
			return 4
		default:
			return 8
		}
	default:
		switch {
		case roll < 40:
			return 2
		case roll < 75:
			return 4
		default:
			return 8
		}
	}
}

func spawnTile(b *Board, mode Mode) bool {
	var empty [][2]int
	for r := range *b {
		for c := range (*b)[r] {
			if (*b)[r][c] == 0 {
				empty = append(empty, [2]int{r, c})
			}
		}
	}
	if len(empty) == 0 {
		return false
	}
	pos := empty[rand.Intn(len(empty))]
	(*b)[pos[0]][pos[1]] = randomTileValue(mode, *b)
	return true
}

// newBoardSized creates a fresh n x n board with two starting tiles.
func newBoardSized(mode Mode, n int) Board {
	b := newEmptyBoard(n)
	spawnTile(&b, mode)
	spawnTile(&b, mode)
	return b
}

func newBoard(mode Mode) Board {
	return newBoardSized(mode, size)
}

func compressRow(row []int) ([]int, bool, int) {
	result := make([]int, len(row))
	values := []int{}
	for _, v := range row {
		if v != 0 {
			values = append(values, v)
		}
	}
	merged := []int{}
	gained := 0
	for i := 0; i < len(values); i++ {
		if i < len(values)-1 && values[i] == values[i+1] {
			mergedVal := values[i] * 2
			merged = append(merged, mergedVal)
			gained += mergedVal
			i++
		} else {
			merged = append(merged, values[i])
		}
	}
	for i, v := range merged {
		result[i] = v
	}
	changed := !equalRows(result, row)
	return result, changed, gained
}

func reverseRow(row []int) []int {
	n := len(row)
	r := make([]int, n)
	for i := 0; i < n; i++ {
		r[i] = row[n-1-i]
	}
	return r
}

func transpose(b Board) Board {
	n := len(b)
	t := newEmptyBoard(n)
	for r := 0; r < n; r++ {
		for c := 0; c < n; c++ {
			t[c][r] = b[r][c]
		}
	}
	return t
}

func moveLeft(b Board) (Board, bool, int) {
	nb := newEmptyBoard(len(b))
	moved := false
	totalGained := 0
	for r := range b {
		newRow, changed, gained := compressRow(b[r])
		nb[r] = newRow
		if changed {
			moved = true
		}
		totalGained += gained
	}
	return nb, moved, totalGained
}

func moveRight(b Board) (Board, bool, int) {
	reversed := newEmptyBoard(len(b))
	for r := range b {
		reversed[r] = reverseRow(b[r])
	}
	nb, moved, gained := moveLeft(reversed)
	result := newEmptyBoard(len(b))
	for r := range nb {
		result[r] = reverseRow(nb[r])
	}
	return result, moved, gained
}

func moveUp(b Board) (Board, bool, int) {
	t := transpose(b)
	nb, moved, gained := moveLeft(t)
	return transpose(nb), moved, gained
}

func moveDown(b Board) (Board, bool, int) {
	t := transpose(b)
	nb, moved, gained := moveRight(t)
	return transpose(nb), moved, gained
}

func hasMoves(b Board) bool {
	n := len(b)
	for r := 0; r < n; r++ {
		for c := 0; c < n; c++ {
			if b[r][c] == 0 {
				return true
			}
			if c < n-1 && b[r][c] == b[r][c+1] {
				return true
			}
			if r < n-1 && b[r][c] == b[r+1][c] {
				return true
			}
		}
	}
	return false
}
