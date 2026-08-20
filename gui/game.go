package main

import (
	"image/color"
	"math/rand"
)

const size = 4
const cellSize = 100

type Board [size][size]int

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
	for r := 0; r < size; r++ {
		for c := 0; c < size; c++ {
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
	for r := 0; r < size; r++ {
		for c := 0; c < size; c++ {
			if b[r][c] == 0 {
				empty = append(empty, [2]int{r, c})
			}
		}
	}
	if len(empty) == 0 {
		return false
	}
	pos := empty[rand.Intn(len(empty))]
	b[pos[0]][pos[1]] = randomTileValue(mode, *b)
	return true
}

func newBoard(mode Mode) Board {
	var b Board
	spawnTile(&b, mode)
	spawnTile(&b, mode)
	return b
}

func compressRow(row [size]int) ([size]int, bool, int) {
	var result [size]int
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
	changed := result != row
	return result, changed, gained
}

func reverseRow(row [size]int) [size]int {
	var r [size]int
	for i := 0; i < size; i++ {
		r[i] = row[size-1-i]
	}
	return r
}

func transpose(b Board) Board {
	var t Board
	for r := 0; r < size; r++ {
		for c := 0; c < size; c++ {
			t[c][r] = b[r][c]
		}
	}
	return t
}

func moveLeft(b Board) (Board, bool, int) {
	var nb Board
	moved := false
	totalGained := 0
	for r := 0; r < size; r++ {
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
	var reversed Board
	for r := 0; r < size; r++ {
		reversed[r] = reverseRow(b[r])
	}
	nb, moved, gained := moveLeft(reversed)
	var result Board
	for r := 0; r < size; r++ {
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
	for r := 0; r < size; r++ {
		for c := 0; c < size; c++ {
			if b[r][c] == 0 {
				return true
			}
			if c < size-1 && b[r][c] == b[r][c+1] {
				return true
			}
			if r < size-1 && b[r][c] == b[r+1][c] {
				return true
			}
		}
	}
	return false
}

func hasWon(b Board) bool {
	for r := 0; r < size; r++ {
		for c := 0; c < size; c++ {
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
