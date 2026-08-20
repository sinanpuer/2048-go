package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
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
		return "Normal"
	case ModeRandomizer:
		return "Randomizer"
	default:
		return "Endlos"
	}
}

// ---------- game logic ----------

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

// ---------- colors ----------

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

// ---------- GUI ----------

type gameUI struct {
	win    fyne.Window
	mode   Mode
	board  Board
	score  int
	won    bool
	over   bool
	rects  [size][size]*canvas.Rectangle
	texts  [size][size]*canvas.Text
	scoreL *widget.Label
	modeL  *widget.Label
	msgL   *widget.Label
}

func newGameUI(win fyne.Window) *gameUI {
	return &gameUI{win: win}
}

func (g *gameUI) buildGrid() fyne.CanvasObject {
	grid := container.New(layout.NewGridLayoutWithColumns(size))
	for r := 0; r < size; r++ {
		for c := 0; c < size; c++ {
			rect := canvas.NewRectangle(color.NRGBA{R: 0xcd, G: 0xc1, B: 0xb4, A: 0xff})
			rect.SetMinSize(fyne.NewSize(cellSize, cellSize))
			text := canvas.NewText("", color.White)
			text.Alignment = fyne.TextAlignCenter
			text.TextStyle = fyne.TextStyle{Bold: true}
			text.TextSize = 26
			g.rects[r][c] = rect
			g.texts[r][c] = text
			cell := container.NewStack(rect, container.NewCenter(text))
			grid.Add(cell)
		}
	}
	return grid
}

func (g *gameUI) refresh() {
	for r := 0; r < size; r++ {
		for c := 0; c < size; c++ {
			v := g.board[r][c]
			bg, fg := tileColor(v)
			g.rects[r][c].FillColor = bg
			g.rects[r][c].Refresh()
			t := g.texts[r][c]
			if v == 0 {
				t.Text = ""
			} else {
				t.Text = fmt.Sprintf("%d", v)
			}
			t.Color = fg
			t.Refresh()
		}
	}
	g.scoreL.SetText(fmt.Sprintf("Punkte: %d", g.score))
	g.modeL.SetText(fmt.Sprintf("Modus: %s", modeName(g.mode)))

	switch {
	case g.over:
		g.msgL.SetText("Game Over! Kein Zug mehr moeglich.")
	case g.won && g.mode != ModeEndless:
		g.msgL.SetText("2048 erreicht! Spiel weiter fuer mehr Punkte.")
	default:
		g.msgL.SetText("")
	}
}

func (g *gameUI) applyMove(fn func(Board) (Board, bool, int)) {
	if g.over {
		return
	}
	nb, moved, gained := fn(g.board)
	if !moved {
		return
	}
	g.board = nb
	g.score += gained
	spawnTile(&g.board, g.mode)
	if !g.won && hasWon(g.board) {
		g.won = true
	}
	if !hasMoves(g.board) {
		g.over = true
	}
	g.refresh()
}

func buildMenu(win fyne.Window) fyne.CanvasObject {
	title := canvas.NewText("2048", color.NRGBA{R: 0x77, G: 0x6e, B: 0x65, A: 0xff})
	title.TextSize = 40
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter

	subtitle := widget.NewLabel("Waehle einen Spielmodus")
	subtitle.Alignment = fyne.TextAlignCenter

	start := func(mode Mode) {
		startGame(win, mode)
	}

	btnNormal := widget.NewButton("Normales Spiel", func() { start(ModeNormal) })
	btnRandom := widget.NewButton("Randomizer-Modus (2, 4, 8)", func() { start(ModeRandomizer) })
	btnEndless := widget.NewButton("Endlos-Modus (ueber 2048 hinaus)", func() { start(ModeEndless) })

	buttons := container.NewVBox(btnNormal, btnRandom, btnEndless)

	content := container.NewVBox(
		title,
		subtitle,
		widget.NewSeparator(),
		buttons,
	)

	return container.NewCenter(content)
}

func startGame(win fyne.Window, mode Mode) {
	g := newGameUI(win)
	g.mode = mode
	g.board = newBoard(mode)

	grid := g.buildGrid()

	g.scoreL = widget.NewLabel("")
	g.modeL = widget.NewLabel("")
	g.msgL = widget.NewLabel("")
	g.msgL.Alignment = fyne.TextAlignCenter

	backBtn := widget.NewButton("Zurueck zum Menue", func() {
		win.Canvas().SetOnTypedKey(nil)
		win.SetContent(buildMenu(win))
	})

	top := container.NewHBox(g.scoreL, layout.NewSpacer(), g.modeL, layout.NewSpacer(), backBtn)

	boardBG := canvas.NewRectangle(color.NRGBA{R: 0xbb, G: 0xad, B: 0xa0, A: 0xff})
	boardArea := container.NewStack(boardBG, container.NewPadded(grid))

	content := container.NewBorder(
		container.NewVBox(top, widget.NewSeparator()),
		g.msgL,
		nil, nil,
		container.NewCenter(boardArea),
	)

	win.SetContent(content)
	g.refresh()

	win.Canvas().SetOnTypedKey(func(ev *fyne.KeyEvent) {
		switch ev.Name {
		case fyne.KeyUp:
			g.applyMove(moveUp)
		case fyne.KeyDown:
			g.applyMove(moveDown)
		case fyne.KeyLeft:
			g.applyMove(moveLeft)
		case fyne.KeyRight:
			g.applyMove(moveRight)
		case fyne.KeyEscape:
			win.Canvas().SetOnTypedKey(nil)
			win.SetContent(buildMenu(win))
		}
	})
}

func main() {
	rand.Seed(time.Now().UnixNano())

	a := app.NewWithID("com.sinan.go2048gui")
	w := a.NewWindow("2048")
	w.Resize(fyne.NewSize(460, 620))
	w.SetContent(buildMenu(w))
	w.ShowAndRun()
}
