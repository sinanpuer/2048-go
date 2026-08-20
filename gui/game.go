package main

import (
	"image/color"
	"math/rand"
)

const size = 4
const cellSize = 100

// Board is a square grid of tiles. Most game modes (Levels, Puzzle, Duel,
// Party) always use newBoard/newBoardWidget with the default `size` (4) and
// never see any other dimension. Free play is the only mode that lets the
// player pick a different board size (via newBoardSized), which Board
// supports by being slice-based rather than a fixed [4][4] array.
type Board [][]int

// cellPxForSize picks a cell pixel size that keeps the rendered board at
// roughly the same total width regardless of how many cells are in a row.
func cellPxForSize(n int) float32 {
	switch n {
	case 5:
		return 80
	case 6:
		return 66
	default:
		return cellSize
	}
}

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

func modeName(m Mode) string {
	switch m {
	case ModeNormal:
		return tr("mode.normal")
	case ModeRandomizer:
		return tr("mode.randomizer")
	default:
		return tr("mode.endless")
	}
}

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

// newBoardSized creates a fresh n x n board with two starting tiles. Free
// play uses this with the player's chosen size; every other mode goes
// through newBoard, which always passes the default size (4).
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

func hasWon(b Board) bool {
	for r := range b {
		for c := range b[r] {
			if b[r][c] >= 2048 {
				return true
			}
		}
	}
	return false
}

func tileColor(v int) (bg color.Color, fg color.Color) {
	dark := color.NRGBA{R: 0x77, G: 0x6e, B: 0x65, A: 0xff}
	white := color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}

	switch v {
	case 0:
		return color.NRGBA{R: 0xcd, G: 0xc1, B: 0xb4, A: 0xff}, white
	case 2:
		return color.NRGBA{R: 0xee, G: 0xe4, B: 0xda, A: 0xff}, dark
	case 4:
		return color.NRGBA{R: 0xcf, G: 0xc7, B: 0xbb, A: 0xff}, dark
	case 8:
		return color.NRGBA{R: 0xf2, G: 0xb1, B: 0x79, A: 0xff}, white
	case 16:
		return color.NRGBA{R: 0xf5, G: 0x95, B: 0x63, A: 0xff}, white
	case 32:
		return color.NRGBA{R: 0xf6, G: 0x7c, B: 0x5f, A: 0xff}, white
	case 64:
		return color.NRGBA{R: 0xf6, G: 0x5e, B: 0x3b, A: 0xff}, white
	case 128:
		return color.NRGBA{R: 0xed, G: 0xcf, B: 0x72, A: 0xff}, white
	case 256:
		return color.NRGBA{R: 0xed, G: 0xcc, B: 0x61, A: 0xff}, white
	case 512:
		return color.NRGBA{R: 0xed, G: 0xc8, B: 0x50, A: 0xff}, white
	case 1024:
		return color.NRGBA{R: 0xed, G: 0xc5, B: 0x3f, A: 0xff}, white
	case 2048:
		return color.NRGBA{R: 0xed, G: 0xc2, B: 0x2e, A: 0xff}, white
	case 4096:
		return color.NRGBA{R: 0x3c, G: 0x3a, B: 0x32, A: 0xff}, white
	case 8192:
		return color.NRGBA{R: 0x2e, G: 0x2c, B: 0x26, A: 0xff}, white
	case 16384:
		return color.NRGBA{R: 0x1a, G: 0x19, B: 0x16, A: 0xff}, color.NRGBA{R: 0xed, G: 0xc2, B: 0x2e, A: 0xff}
	default:
		return color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xff}, color.NRGBA{R: 0xed, G: 0xc2, B: 0x2e, A: 0xff}
	}
}
