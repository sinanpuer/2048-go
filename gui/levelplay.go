package main

import (
	"image/color"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type levelUI struct {
	win       fyne.Window
	def       LevelDef
	board     Board
	score     int
	combo     int
	bw        *boardWidget
	scoreL    *widget.Label
	timeL     *widget.Label
	comboL    *widget.Label
	startedAt time.Time
	finished  bool
	stopTimer chan struct{}
	stopOnce  sync.Once
}

func startLevel(win fyne.Window, levelNumber int) {
	var def LevelDef
	for _, l := range allLevels {
		if l.Number == levelNumber {
			def = l
			break
		}
	}

	lu := &levelUI{win: win, def: def, board: newBoard(ModeRandomizer), startedAt: time.Now()}
	lu.bw = newBoardWidget(cellSize, size)

	lu.scoreL = widget.NewLabel("")
	lu.timeL = widget.NewLabel("")
	lu.comboL = widget.NewLabel("")
	goalL := widget.NewLabelWithStyle(def.Title(), fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	backBtn := widget.NewButton(tr("level.backSelect"), func() {
		lu.stop()
		win.Canvas().SetOnTypedKey(nil)
		win.SetContent(buildLevelSelect(win))
	})

	header := container.NewVBox(
		widget.NewLabelWithStyle(trf("level.number", def.Number), fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		goalL,
		container.NewHBox(lu.scoreL, layout.NewSpacer(), lu.timeL, layout.NewSpacer(), backBtn),
		container.NewHBox(lu.comboL),
		widget.NewSeparator(),
	)

	boardBG := canvas.NewRectangle(color.NRGBA{R: 0xbb, G: 0xad, B: 0xa0, A: 0xff})
	boardArea := container.NewStack(boardBG, container.NewPadded(lu.bw.container))

	content := container.NewBorder(header, nil, nil, nil, container.NewCenter(boardArea))
	win.SetContent(content)
	lu.render()

	win.Canvas().SetOnTypedKey(func(ev *fyne.KeyEvent) {
		switch ev.Name {
		case fyne.KeyUp:
			lu.applyMove(moveUp)
		case fyne.KeyDown:
			lu.applyMove(moveDown)
		case fyne.KeyLeft:
			lu.applyMove(moveLeft)
		case fyne.KeyRight:
			lu.applyMove(moveRight)
		case fyne.KeyEscape:
			lu.stop()
			win.Canvas().SetOnTypedKey(nil)
			win.SetContent(buildLevelSelect(win))
		}
	})

	if def.TimeLimit > 0 {
		lu.timeL.SetText(trf("level.time", def.TimeLimit))
		lu.startTimer()
	}
}

func (lu *levelUI) render() {
	lu.bw.render(lu.board)
	lu.scoreL.SetText(trf("game.score", lu.score))
	if lu.combo > 0 {
		lu.comboL.SetText(trf("game.combo", lu.combo, comboMultiplier(lu.combo)))
	} else {
		lu.comboL.SetText("")
	}
}

func (lu *levelUI) applyMove(fn func(Board) (Board, bool, int)) {
	if lu.finished {
		return
	}
	nb, moved, gained := fn(lu.board)
	if !moved {
		return
	}
	lu.board = nb
	lu.combo, gained = applyCombo(lu.combo, gained)
	lu.score += gained
	spawnTile(&lu.board, ModeRandomizer)
	lu.render()

	if lu.goalReached() {
		lu.finishWin()
		return
	}
	if !hasMoves(lu.board) {
		lu.finishLose(tr("game.gameoverMsg"))
	}
}

func (lu *levelUI) goalReached() bool {
	switch lu.def.Kind {
	case KindScore, KindTimedScore:
		return lu.score >= lu.def.ScoreGoal
	case KindTile, KindTimedTile:
		return highestTile(lu.board) >= lu.def.TileGoal
	}
	return false
}

func (lu *levelUI) startTimer() {
	lu.stopTimer = make(chan struct{})
	stopCh := lu.stopTimer

	go func() {
		ticker := time.NewTicker(200 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-stopCh:
				return
			case <-ticker.C:
				remaining := lu.def.TimeLimit - int(time.Since(lu.startedAt).Seconds())
				if remaining < 0 {
					remaining = 0
				}
				fyne.Do(func() {
					if lu.finished {
						return
					}
					lu.timeL.SetText(trf("level.time", remaining))
					if remaining <= 0 {
						lu.finishLose(tr("level.timeUp"))
					}
				})
				if remaining <= 0 {
					return
				}
			}
		}
	}()
}

func (lu *levelUI) stop() {
	if lu.stopTimer == nil {
		return
	}
	lu.stopOnce.Do(func() { close(lu.stopTimer) })
}

func (lu *levelUI) finishWin() {
	lu.finished = true
	lu.stop()
	progress.markCompleted(lu.def.Number)
	progress.save()
	lu.showResult(true, tr("level.win"))
}

func (lu *levelUI) finishLose(reason string) {
	lu.finished = true
	lu.stop()
	lu.showResult(false, reason)
}

func (lu *levelUI) showResult(won bool, message string) {
	win := lu.win
	resultLabel := widget.NewLabelWithStyle(message, fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	scoreLine := widget.NewLabel(trf("game.score", lu.score))
	scoreLine.Alignment = fyne.TextAlignCenter

	retryLabel := tr("level.retry")
	if won {
		retryLabel = tr("level.retryWin")
	}
	retryBtn := widget.NewButton(retryLabel, func() {
		win.Canvas().SetOnTypedKey(nil)
		startLevel(win, lu.def.Number)
	})
	backBtn := widget.NewButton(tr("level.backSelect"), func() {
		win.Canvas().SetOnTypedKey(nil)
		win.SetContent(buildLevelSelect(win))
	})

	buttonRow := []fyne.CanvasObject{retryBtn}
	nextNum := lu.def.Number + 1
	if won && nextNum <= totalLevels && progress.isUnlocked(nextNum) {
		nextBtn := widget.NewButton(tr("level.next"), func() {
			win.Canvas().SetOnTypedKey(nil)
			startLevel(win, nextNum)
		})
		buttonRow = append(buttonRow, nextBtn)
	}
	buttonRow = append(buttonRow, backBtn)

	header := widget.NewLabelWithStyle(trf("level.number", lu.def.Number), fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	boardBG := canvas.NewRectangle(color.NRGBA{R: 0xbb, G: 0xad, B: 0xa0, A: 0xff})
	boardArea := container.NewStack(boardBG, container.NewPadded(lu.bw.container))

	content := container.NewBorder(
		container.NewVBox(header, resultLabel, scoreLine, widget.NewSeparator()),
		container.NewCenter(container.NewHBox(buttonRow...)),
		nil, nil,
		container.NewCenter(boardArea),
	)
	win.SetContent(content)
}
