package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"golang.org/x/term"
)

const size = 4

type Board [size][size]int

type Mode int

const (
	ModeNormal Mode = iota
	ModeRandomizer
)

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
	if mode == ModeNormal {
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

func newBoard(mode Mode) Board {
	var b Board
	spawnTile(&b, mode)
	spawnTile(&b, mode)
	return b
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

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func printBoard(b Board, score int, mode Mode) {
	clearScreen()
	modeName := "Normal"
	if mode == ModeRandomizer {
		modeName = "Randomizer"
	}
	fmt.Printf("2048 [%s] - Score: %d\r\n", modeName, score)
	fmt.Print(strings.Repeat("-", size*7+1) + "\r\n")
	for r := 0; r < size; r++ {
		row := "|"
		for c := 0; c < size; c++ {
			if b[r][c] == 0 {
				row += "      |"
			} else {
				row += fmt.Sprintf("%6d|", b[r][c])
			}
		}
		fmt.Print(row + "\r\n")
		fmt.Print(strings.Repeat("-", size*7+1) + "\r\n")
	}
	fmt.Print("Pfeiltasten = bewegen, q = zurueck zum Menue\r\n")
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

// key represents a single normalized user input.
type key int

const (
	keyUp key = iota
	keyDown
	keyLeft
	keyRight
	keyEnter
	keyQuit
	keyOther
)

func readKey() (key, error) {
	buf := make([]byte, 3)
	n, err := os.Stdin.Read(buf)
	if err != nil {
		return keyOther, err
	}

	if n == 1 {
		switch buf[0] {
		case 'w', 'W':
			return keyUp, nil
		case 'a', 'A':
			return keyLeft, nil
		case 's', 'S':
			return keyDown, nil
		case 'd', 'D':
			return keyRight, nil
		case 'q', 'Q':
			return keyQuit, nil
		case 13, 10:
			return keyEnter, nil
		case 3:
			return keyQuit, nil
		}
		return keyOther, nil
	}

	if n >= 3 && buf[0] == 27 && buf[1] == 91 {
		switch buf[2] {
		case 'A':
			return keyUp, nil
		case 'B':
			return keyDown, nil
		case 'C':
			return keyRight, nil
		case 'D':
			return keyLeft, nil
		}
	}

	return keyOther, nil
}

func playGame(mode Mode) {
	board := newBoard(mode)
	score := 0
	won := false

	for {
		printBoard(board, score, mode)
		if won {
			fmt.Print("Du hast 2048 erreicht! Weiterspielen oder q fuer Menue.\r\n")
		}
		if !hasMoves(board) {
			fmt.Printf("Game Over! Kein Zug mehr moeglich. Endpunktzahl: %d\r\n", score)
			fmt.Print("Beliebige Taste = zurueck zum Menue\r\n")
			readKey()
			return
		}

		k, err := readKey()
		if err != nil {
			return
		}

		var newBoard Board
		var moved bool
		var gained int

		switch k {
		case keyUp:
			newBoard, moved, gained = moveUp(board)
		case keyLeft:
			newBoard, moved, gained = moveLeft(board)
		case keyDown:
			newBoard, moved, gained = moveDown(board)
		case keyRight:
			newBoard, moved, gained = moveRight(board)
		case keyQuit:
			return
		default:
			continue
		}

		if moved {
			board = newBoard
			score += gained
			spawnTile(&board, mode)
			if !won && hasWon(board) {
				won = true
			}
		}
	}
}

func showMenu() (Mode, bool) {
	options := []string{"Normales Spiel", "Randomizer-Modus (2, 4 und 8)"}
	selected := 0

	for {
		clearScreen()
		fmt.Print("=== 2048 - Hauptmenue ===\r\n\r\n")
		for i, opt := range options {
			cursor := "  "
			if i == selected {
				cursor = "> "
			}
			fmt.Printf("%s%s\r\n", cursor, opt)
		}
		fmt.Print("\r\nPfeiltasten hoch/runter = auswaehlen, Enter = bestaetigen, q = beenden\r\n")

		k, err := readKey()
		if err != nil {
			return ModeNormal, false
		}

		switch k {
		case keyUp:
			selected = (selected - 1 + len(options)) % len(options)
		case keyDown:
			selected = (selected + 1) % len(options)
		case keyEnter:
			if selected == 0 {
				return ModeNormal, true
			}
			return ModeRandomizer, true
		case keyQuit:
			return ModeNormal, false
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	if !term.IsTerminal(int(os.Stdin.Fd())) {
		fmt.Println("Dieses Spiel muss in einem echten Terminal ausgefuehrt werden.")
		os.Exit(1)
	}

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println("Konnte Terminal nicht in Raw-Modus versetzen:", err)
		os.Exit(1)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	for {
		mode, ok := showMenu()
		if !ok {
			clearScreen()
			fmt.Print("Bis zum naechsten Mal!\r\n")
			return
		}
		playGame(mode)
	}
}
