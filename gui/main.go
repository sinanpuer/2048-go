package main

import (
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
var allPuzzleLevels = generatePuzzleLevels(totalPuzzleLevels)
var progress = loadProgress()
var highscores = loadHighscores()
var settings = loadSettings()

// ---------- free-play GUI ----------

type gameUI struct {
	win     fyne.Window
	mode    Mode
	board   Board
	score   int
	combo   int
	won     bool
	over    bool
	newHigh bool
	bw      *boardWidget
	scoreL  *widget.Label
	modeL   *widget.Label
	highL   *widget.Label
	comboL  *widget.Label
	msgL    *widget.Label
}

func (g *gameUI) refresh() {
	g.bw.render(g.board)
	g.scoreL.SetText(trf("game.score", g.score))
	g.modeL.SetText(trf("game.mode", modeName(g.mode)))
	g.highL.SetText(trf("game.highscore", highscores.get(g.mode)))
	if g.combo > 0 {
		g.comboL.SetText(trf("game.combo", g.combo, comboMultiplier(g.combo)))
	} else {
		g.comboL.SetText("")
	}

	switch {
	case g.newHigh:
		g.msgL.SetText(tr("game.newHigh"))
	case g.won && g.mode != ModeEndless:
		g.msgL.SetText(tr("game.won2048"))
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
	g.combo, gained = applyCombo(g.combo, gained)
	g.score += gained
	spawnTile(&g.board, g.mode)
	if highscores.update(g.mode, g.score) {
		g.newHigh = true
	}
	if !g.won && hasWon(g.board) {
		g.won = true
	}
	g.refresh()

	if !hasMoves(g.board) {
		g.over = true
		g.showGameOver()
	}
}

func (g *gameUI) showGameOver() {
	win := g.win
	win.Canvas().SetOnTypedKey(nil)

	title := widget.NewLabelWithStyle(tr("game.gameover"), fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	scoreLine := widget.NewLabel(trf("game.score", g.score))
	scoreLine.Alignment = fyne.TextAlignCenter

	restartBtn := widget.NewButton(tr("game.restart"), func() {
		win.Canvas().SetOnTypedKey(nil)
		startGame(win, g.mode)
	})
	menuBtn := widget.NewButton(tr("game.mainmenu"), func() {
		win.Canvas().SetOnTypedKey(nil)
		win.SetContent(buildMenu(win))
	})

	panel := container.NewVBox(
		title, scoreLine, widget.NewSeparator(),
		container.NewCenter(container.NewHBox(restartBtn, menuBtn)),
	)

	boardBG := canvas.NewRectangle(color.NRGBA{R: 0xbb, G: 0xad, B: 0xa0, A: 0xff})
	boardArea := container.NewStack(boardBG, container.NewPadded(g.bw.container))

	overlay := canvas.NewRectangle(color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xb0})

	content := container.NewStack(
		container.NewCenter(boardArea),
		overlay,
		container.NewCenter(container.NewPadded(panel)),
	)
	win.SetContent(content)
}

func startGame(win fyne.Window, mode Mode) {
	g := &gameUI{win: win, mode: mode, board: newBoard(mode), bw: newBoardWidget(cellSize)}

	g.scoreL = widget.NewLabel("")
	g.modeL = widget.NewLabel("")
	g.highL = widget.NewLabel("")
	g.comboL = widget.NewLabel("")
	g.msgL = widget.NewLabel("")
	g.msgL.Alignment = fyne.TextAlignCenter

	backBtn := widget.NewButton(tr("game.back"), func() {
		win.Canvas().SetOnTypedKey(nil)
		win.SetContent(buildMenu(win))
	})

	top := container.NewVBox(
		container.NewHBox(g.modeL, layout.NewSpacer(), backBtn),
		container.NewHBox(g.scoreL, layout.NewSpacer(), g.highL),
		container.NewHBox(g.comboL),
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
	ambient := newAmbientWall(4, 5, 34)
	ambient.start()

	title := canvas.NewText("Merge Kingdom", color.White)
	title.TextSize = 38
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter

	subtitle := canvas.NewText(tr("menu.subtitle"), color.NRGBA{R: 0xee, G: 0xee, B: 0xee, A: 0xff})
	subtitle.Alignment = fyne.TextAlignCenter
	subtitle.TextSize = 14

	goTo := func(fn func()) func() {
		return func() {
			ambient.stop()
			fn()
		}
	}

	btnNormal := widget.NewButton(tr("menu.normal"), goTo(func() { startGame(win, ModeNormal) }))
	btnRandom := widget.NewButton(tr("menu.randomizer"), goTo(func() { startGame(win, ModeRandomizer) }))
	btnEndless := widget.NewButton(tr("menu.endless"), goTo(func() { startGame(win, ModeEndless) }))
	btnLevels := widget.NewButton(trf("menu.levels", totalLevels), goTo(func() {
		win.SetContent(buildLevelSelect(win))
	}))
	btnDuel := widget.NewButton(tr("menu.duel"), goTo(func() {
		win.SetContent(buildBotSetup(win))
	}))
	btnPuzzle := widget.NewButton(trf("menu.puzzle", totalPuzzleLevels), goTo(func() {
		win.SetContent(buildPuzzleSelect(win))
	}))
	btnSettings := widget.NewButton(tr("menu.settings"), goTo(func() {
		win.SetContent(buildSettings(win))
	}))
	btnParty := widget.NewButton(tr("menu.party"), goTo(func() {
		win.SetContent(buildPartySetup(win))
	}))

	buttons := container.NewVBox(btnNormal, btnRandom, btnEndless, btnLevels, btnPuzzle, btnDuel, btnParty, btnSettings)

	highscoreLine := canvas.NewText(
		trf("menu.highscores", highscores.get(ModeNormal), highscores.get(ModeRandomizer), highscores.get(ModeEndless)),
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
	w := a.NewWindow("Merge Kingdom")
	w.Resize(fyne.NewSize(720, 820))
	w.SetFixedSize(true)
	w.SetContent(buildMenu(w))
	w.ShowAndRun()
}
