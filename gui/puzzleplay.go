package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type puzzleUI struct {
	win       fyne.Window
	def       PuzzleLevelDef
	board     Board
	score     int
	combo     int
	movesUsed int
	finished  bool
	bw        *boardWidget
	scoreL    *widget.Label
	movesL    *widget.Label
	comboL    *widget.Label
}

func startPuzzle(win fyne.Window, levelNumber int) {
	var def PuzzleLevelDef
	for _, l := range allPuzzleLevels {
		if l.Number == levelNumber {
			def = l
			break
		}
	}

	pu := &puzzleUI{win: win, def: def, board: def.StartBoard}
	pu.bw = newBoardWidget(cellSize)
	pu.scoreL = widget.NewLabel("")
	pu.movesL = widget.NewLabel("")
	pu.comboL = widget.NewLabel("")
	goalL := widget.NewLabelWithStyle(def.Title(), fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	backBtn := widget.NewButton(tr("puzzle.backSelect"), func() {
		win.Canvas().SetOnTypedKey(nil)
		win.SetContent(buildPuzzleSelect(win))
	})

	header := container.NewVBox(
		widget.NewLabelWithStyle(trf("puzzle.number", def.Number), fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		goalL,
		container.NewHBox(pu.scoreL, layout.NewSpacer(), pu.movesL, layout.NewSpacer(), backBtn),
		container.NewHBox(pu.comboL),
		widget.NewSeparator(),
	)

	boardBG := canvas.NewRectangle(color.NRGBA{R: 0xbb, G: 0xad, B: 0xa0, A: 0xff})
	boardArea := container.NewStack(boardBG, container.NewPadded(pu.bw.container))

	win.SetContent(container.NewBorder(header, nil, nil, nil, container.NewCenter(boardArea)))
	pu.render()

	win.Canvas().SetOnTypedKey(func(ev *fyne.KeyEvent) {
		switch ev.Name {
		case fyne.KeyUp:
			pu.applyMove(moveUp)
		case fyne.KeyDown:
			pu.applyMove(moveDown)
		case fyne.KeyLeft:
			pu.applyMove(moveLeft)
		case fyne.KeyRight:
			pu.applyMove(moveRight)
		case fyne.KeyEscape:
			win.Canvas().SetOnTypedKey(nil)
			win.SetContent(buildPuzzleSelect(win))
		}
	})
}

func (pu *puzzleUI) render() {
	pu.bw.render(pu.board)
	pu.scoreL.SetText(trf("puzzle.score", pu.score, pu.def.ScoreGoal))
	pu.movesL.SetText(trf("puzzle.moves", pu.movesUsed, pu.def.MoveLimit))
	if pu.combo > 0 {
		pu.comboL.SetText(trf("game.combo", pu.combo, comboMultiplier(pu.combo)))
	} else {
		pu.comboL.SetText("")
	}
}

func (pu *puzzleUI) applyMove(fn func(Board) (Board, bool, int)) {
	if pu.finished {
		return
	}
	nb, moved, gained := fn(pu.board)
	if !moved {
		return
	}
	pu.board = nb
	pu.combo, gained = applyCombo(pu.combo, gained)
	pu.score += gained
	pu.movesUsed++
	spawnTile(&pu.board, ModeRandomizer)
	pu.render()

	if pu.score >= pu.def.ScoreGoal {
		pu.finishWin()
		return
	}
	if pu.movesUsed >= pu.def.MoveLimit {
		pu.finishLose(tr("puzzle.moveLimit"))
		return
	}
	if !hasMoves(pu.board) {
		pu.finishLose(tr("game.gameoverMsg"))
	}
}

func (pu *puzzleUI) finishWin() {
	pu.finished = true
	progress.markPuzzleCompleted(pu.def.Number)
	progress.save()
	pu.showResult(true, tr("puzzle.win"))
}

func (pu *puzzleUI) finishLose(reason string) {
	pu.finished = true
	pu.showResult(false, reason)
}

func (pu *puzzleUI) showResult(won bool, message string) {
	win := pu.win
	win.Canvas().SetOnTypedKey(nil)

	resultLabel := widget.NewLabelWithStyle(message, fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	scoreLine := widget.NewLabel(trf("puzzle.score", pu.score, pu.def.ScoreGoal))
	scoreLine.Alignment = fyne.TextAlignCenter

	retryLabel := tr("level.retry")
	if won {
		retryLabel = tr("level.retryWin")
	}
	retryBtn := widget.NewButton(retryLabel, func() {
		win.Canvas().SetOnTypedKey(nil)
		startPuzzle(win, pu.def.Number)
	})
	backBtn := widget.NewButton(tr("puzzle.backSelect"), func() {
		win.Canvas().SetOnTypedKey(nil)
		win.SetContent(buildPuzzleSelect(win))
	})

	buttonRow := []fyne.CanvasObject{retryBtn}
	nextNum := pu.def.Number + 1
	if won && nextNum <= totalPuzzleLevels && progress.isPuzzleUnlocked(nextNum) {
		nextBtn := widget.NewButton(tr("puzzle.next"), func() {
			win.Canvas().SetOnTypedKey(nil)
			startPuzzle(win, nextNum)
		})
		buttonRow = append(buttonRow, nextBtn)
	}
	buttonRow = append(buttonRow, backBtn)

	header := widget.NewLabelWithStyle(trf("puzzle.number", pu.def.Number), fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	boardBG := canvas.NewRectangle(color.NRGBA{R: 0xbb, G: 0xad, B: 0xa0, A: 0xff})
	boardArea := container.NewStack(boardBG, container.NewPadded(pu.bw.container))

	content := container.NewBorder(
		container.NewVBox(header, resultLabel, scoreLine, widget.NewSeparator()),
		container.NewCenter(container.NewHBox(buttonRow...)),
		nil, nil,
		container.NewCenter(boardArea),
	)
	win.SetContent(content)
}
