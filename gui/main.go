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

var allLevels = generateLevels(totalLevels)
var progress = loadProgress()
var highscores = loadHighscores()

// ---------- free-play GUI ----------

type gameUI struct {
	win     fyne.Window
	mode    Mode
	board   Board
	score   int
	won     bool
	over    bool
	newHigh bool
	bw      *boardWidget
	scoreL  *widget.Label
	modeL   *widget.Label
	highL   *widget.Label
	msgL    *widget.Label
}

func (g *gameUI) refresh() {
	g.bw.render(g.board)
	g.scoreL.SetText(fmt.Sprintf("Punkte: %d", g.score))
	g.modeL.SetText(fmt.Sprintf("Modus: %s", modeName(g.mode)))
	g.highL.SetText(fmt.Sprintf("Highscore: %d", highscores.get(g.mode)))

	switch {
	case g.over:
		g.msgL.SetText("Game Over! Kein Zug mehr moeglich.")
	case g.newHigh:
		g.msgL.SetText("Neuer Highscore!")
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
	if highscores.update(g.mode, g.score) {
		g.newHigh = true
	}
	if !g.won && hasWon(g.board) {
		g.won = true
	}
	if !hasMoves(g.board) {
		g.over = true
	}
	g.refresh()
}

func startGame(win fyne.Window, mode Mode) {
	g := &gameUI{win: win, mode: mode, board: newBoard(mode), bw: newBoardWidget(cellSize)}

	g.scoreL = widget.NewLabel("")
	g.modeL = widget.NewLabel("")
	g.highL = widget.NewLabel("")
	g.msgL = widget.NewLabel("")
	g.msgL.Alignment = fyne.TextAlignCenter

	backBtn := widget.NewButton("Zurueck zum Menue", func() {
		win.Canvas().SetOnTypedKey(nil)
		win.SetContent(buildMenu(win))
	})

	top := container.NewVBox(
		container.NewHBox(g.modeL, layout.NewSpacer(), backBtn),
		container.NewHBox(g.scoreL, layout.NewSpacer(), g.highL),
	)

	boardBG := canvas.NewRectangle(color.NRGBA{R: 0xbb, G: 0xad, B: 0xa0, A: 0xff})
	boardArea := container.NewStack(boardBG, container.NewPadded(g.bw.container))

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

// ---------- main menu ----------

func buildMenu(win fyne.Window) fyne.CanvasObject {
	ambient := newAmbientBoard()
	ambient.start()

	title := canvas.NewText("2048", color.White)
	title.TextSize = 44
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter

	subtitle := canvas.NewText("Waehle einen Spielmodus", color.NRGBA{R: 0xee, G: 0xee, B: 0xee, A: 0xff})
	subtitle.Alignment = fyne.TextAlignCenter
	subtitle.TextSize = 14

	goTo := func(fn func()) func() {
		return func() {
			ambient.stop()
			fn()
		}
	}

	btnNormal := widget.NewButton("Normales Spiel", goTo(func() { startGame(win, ModeNormal) }))
	btnRandom := widget.NewButton("Randomizer-Modus (2, 4, 8)", goTo(func() { startGame(win, ModeRandomizer) }))
	btnEndless := widget.NewButton("Endlos-Modus (ueber 2048 hinaus)", goTo(func() { startGame(win, ModeEndless) }))
	btnLevels := widget.NewButton(fmt.Sprintf("Level-Modus (%d Level)", totalLevels), goTo(func() {
		win.SetContent(buildLevelSelect(win))
	}))
	btnDuel := widget.NewButton("KI-Duell (gegen Bot)", goTo(func() {
		win.SetContent(buildBotSetup(win))
	}))

	buttons := container.NewVBox(btnNormal, btnRandom, btnEndless, btnLevels, btnDuel)

	highscoreLine := canvas.NewText(
		fmt.Sprintf("Bestwerte — Normal: %d   Randomizer: %d   Endlos: %d",
			highscores.get(ModeNormal), highscores.get(ModeRandomizer), highscores.get(ModeEndless)),
		color.NRGBA{R: 0xee, G: 0xee, B: 0xee, A: 0xff},
	)
	highscoreLine.Alignment = fyne.TextAlignCenter
	highscoreLine.TextSize = 12

	panel := container.NewVBox(
		title,
		subtitle,
		widget.NewSeparator(),
		buttons,
		highscoreLine,
	)

	overlay := canvas.NewRectangle(color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xa8})

	return container.NewStack(
		ambient.canvasObject(),
		overlay,
		container.NewCenter(container.NewPadded(panel)),
	)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	a := app.NewWithID("com.sinan.go2048gui")
	w := a.NewWindow("2048")
	w.Resize(fyne.NewSize(720, 780))
	w.SetContent(buildMenu(w))
	w.ShowAndRun()
}
