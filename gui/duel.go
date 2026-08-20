package main

import (
	"fmt"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func buildBotSetup(win fyne.Window) fyne.CanvasObject {
	title := widget.NewLabelWithStyle(tr("duel.title"), fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	desc := widget.NewLabel(tr("duel.desc"))
	desc.Alignment = fyne.TextAlignCenter
	desc.Wrapping = fyne.TextWrapWord

	duration := 60
	difficulty := BotNormal

	bulletLabel, normalLengthLabel := tr("duel.bullet"), tr("duel.normalLength")
	durationRadio := widget.NewRadioGroup([]string{bulletLabel, normalLengthLabel}, func(s string) {
		if s == normalLengthLabel {
			duration = 300
		} else {
			duration = 60
		}
	})
	durationRadio.SetSelected(bulletLabel)

	botNormalLabel, botExpertLabel := tr("bot.normal"), tr("bot.expert")
	difficultyRadio := widget.NewRadioGroup([]string{botNormalLabel, botExpertLabel}, func(s string) {
		if s == botExpertLabel {
			difficulty = BotExpert
		} else {
			difficulty = BotNormal
		}
	})
	difficultyRadio.SetSelected(botNormalLabel)

	startBtn := widget.NewButton(tr("duel.start"), func() {
		startDuel(win, duration, difficulty)
	})
	backBtn := widget.NewButton(tr("duel.back"), func() {
		win.SetContent(buildMenu(win))
	})

	content := container.NewVBox(
		title, desc, widget.NewSeparator(),
		widget.NewLabel(tr("duel.duration")), durationRadio,
		widget.NewLabel(tr("duel.difficulty")), difficultyRadio,
		widget.NewSeparator(),
		startBtn, backBtn,
	)
	return container.NewCenter(container.NewPadded(content))
}

type duelUI struct {
	win           fyne.Window
	duration      int
	difficulty    BotDifficulty
	playerBoard   Board
	botBoard      Board
	playerScore   int
	botScore      int
	playerCombo   int
	botCombo      int
	playerStuck   bool
	botStuck      bool
	playerBW      *boardWidget
	botBW         *boardWidget
	playerScoreL  *widget.Label
	botScoreL     *widget.Label
	playerComboL  *widget.Label
	botComboL     *widget.Label
	playerStatusL *widget.Label
	botStatusL    *widget.Label
	timeL         *widget.Label
	startedAt     time.Time
	finished      bool
	stopCh        chan struct{}
	stopOnce      sync.Once
}

func startDuel(win fyne.Window, duration int, difficulty BotDifficulty) {
	d := &duelUI{
		win:         win,
		duration:    duration,
		difficulty:  difficulty,
		playerBoard: newBoard(ModeRandomizer),
		botBoard:    newBoard(ModeRandomizer),
		startedAt:   time.Now(),
		stopCh:      make(chan struct{}),
	}
	d.playerBW = newBoardWidget(64, size)
	d.botBW = newBoardWidget(64, size)
	d.playerScoreL = widget.NewLabel("")
	d.botScoreL = widget.NewLabel("")
	d.playerComboL = widget.NewLabel("")
	d.botComboL = widget.NewLabel("")
	d.playerStatusL = widget.NewLabel("")
	d.botStatusL = widget.NewLabel("")
	d.timeL = widget.NewLabelWithStyle(fmt.Sprintf("%ds", duration), fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	backBtn := widget.NewButton(tr("duel.cancel"), func() {
		d.stop()
		win.Canvas().SetOnTypedKey(nil)
		win.SetContent(buildBotSetup(win))
	})

	playerPanel := container.NewVBox(
		widget.NewLabelWithStyle(tr("duel.you"), fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		d.playerScoreL,
		d.playerComboL,
		d.playerStatusL,
		container.NewPadded(d.playerBW.container),
	)
	botPanel := container.NewVBox(
		widget.NewLabelWithStyle(trf("duel.bot", botDifficultyName(difficulty)), fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		d.botScoreL,
		d.botComboL,
		d.botStatusL,
		container.NewPadded(d.botBW.container),
	)

	boards := container.NewGridWithColumns(2, playerPanel, botPanel)

	header := container.NewVBox(
		widget.NewLabelWithStyle(tr("duel.title"), fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		container.NewHBox(layout.NewSpacer(), d.timeL, layout.NewSpacer()),
		container.NewCenter(backBtn),
		widget.NewSeparator(),
	)

	win.SetContent(container.NewBorder(header, nil, nil, nil, container.NewCenter(boards)))
	d.render()

	win.Canvas().SetOnTypedKey(func(ev *fyne.KeyEvent) {
		if d.finished {
			return
		}
		switch ev.Name {
		case fyne.KeyUp:
			d.playerMove(moveUp)
		case fyne.KeyDown:
			d.playerMove(moveDown)
		case fyne.KeyLeft:
			d.playerMove(moveLeft)
		case fyne.KeyRight:
			d.playerMove(moveRight)
		case fyne.KeyEscape:
			d.stop()
			win.Canvas().SetOnTypedKey(nil)
			win.SetContent(buildBotSetup(win))
		}
	})

	d.startBotLoop()
	d.startTimerLoop()
}

func (d *duelUI) render() {
	d.playerBW.render(d.playerBoard)
	d.botBW.render(d.botBoard)
	d.playerScoreL.SetText(trf("game.score", d.playerScore))
	d.botScoreL.SetText(trf("game.score", d.botScore))

	if d.playerCombo > 0 {
		d.playerComboL.SetText(trf("game.combo", d.playerCombo, comboMultiplier(d.playerCombo)))
	} else {
		d.playerComboL.SetText("")
	}
	if d.botCombo > 0 {
		d.botComboL.SetText(trf("game.combo", d.botCombo, comboMultiplier(d.botCombo)))
	} else {
		d.botComboL.SetText("")
	}

	if d.playerStuck {
		d.playerStatusL.SetText(tr("duel.gameOverStuck"))
	} else {
		d.playerStatusL.SetText("")
	}
	if d.botStuck {
		d.botStatusL.SetText(tr("duel.gameOverStuck"))
	} else {
		d.botStatusL.SetText("")
	}
}

func (d *duelUI) checkBothStuck() {
	if d.finished {
		return
	}
	if d.playerStuck && d.botStuck {
		d.finish()
	}
}

func (d *duelUI) playerMove(fn func(Board) (Board, bool, int)) {
	if d.finished || d.playerStuck {
		return
	}
	nb, moved, gained := fn(d.playerBoard)
	if !moved {
		return
	}
	d.playerBoard = nb
	d.playerCombo, gained = applyCombo(d.playerCombo, gained)
	d.playerScore += gained
	spawnTile(&d.playerBoard, ModeRandomizer)
	if !hasMoves(d.playerBoard) {
		d.playerStuck = true
	}
	d.render()
	d.checkBothStuck()
}

func (d *duelUI) startBotLoop() {
	stopCh := d.stopCh
	go func() {
		ticker := time.NewTicker(botMoveInterval(d.difficulty))
		defer ticker.Stop()
		for {
			select {
			case <-stopCh:
				return
			case <-ticker.C:
				fyne.Do(func() {
					if d.finished || d.botStuck {
						return
					}
					mv, ok := botPickMove(d.botBoard, d.difficulty)
					if !ok {
						d.botStuck = true
						d.render()
						d.checkBothStuck()
						return
					}
					nb, moved, gained := mv(d.botBoard)
					if !moved {
						d.botStuck = true
						d.render()
						d.checkBothStuck()
						return
					}
					d.botBoard = nb
					d.botCombo, gained = applyCombo(d.botCombo, gained)
					d.botScore += gained
					spawnTile(&d.botBoard, ModeRandomizer)
					if !hasMoves(d.botBoard) {
						d.botStuck = true
					}
					d.render()
					d.checkBothStuck()
				})
			}
		}
	}()
}

func (d *duelUI) startTimerLoop() {
	stopCh := d.stopCh
	go func() {
		ticker := time.NewTicker(200 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-stopCh:
				return
			case <-ticker.C:
				remaining := d.duration - int(time.Since(d.startedAt).Seconds())
				if remaining < 0 {
					remaining = 0
				}
				fyne.Do(func() {
					if d.finished {
						return
					}
					d.timeL.SetText(fmt.Sprintf("%ds", remaining))
					if remaining <= 0 {
						d.finish()
					}
				})
				if remaining <= 0 {
					return
				}
			}
		}
	}()
}

func (d *duelUI) stop() {
	d.stopOnce.Do(func() { close(d.stopCh) })
}

func (d *duelUI) finish() {
	d.finished = true
	d.stop()

	var result string
	switch {
	case d.playerScore > d.botScore:
		result = tr("duel.youWin")
	case d.playerScore < d.botScore:
		result = tr("duel.botWins")
	default:
		result = tr("duel.draw")
	}

	win := d.win
	resultLabel := widget.NewLabelWithStyle(result, fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	scoreLine := widget.NewLabel(trf("duel.scoreLine", d.playerScore, d.botScore))
	scoreLine.Alignment = fyne.TextAlignCenter

	retryBtn := widget.NewButton(tr("duel.retry"), func() {
		win.Canvas().SetOnTypedKey(nil)
		startDuel(win, d.duration, d.difficulty)
	})
	backBtn := widget.NewButton(tr("duel.back"), func() {
		win.Canvas().SetOnTypedKey(nil)
		win.SetContent(buildMenu(win))
	})

	header := container.NewVBox(
		widget.NewLabelWithStyle(tr("duel.finished"), fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		resultLabel, scoreLine, widget.NewSeparator(),
	)

	playerPanel := container.NewVBox(
		widget.NewLabelWithStyle(tr("duel.you"), fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		container.NewPadded(d.playerBW.container),
	)
	botPanel := container.NewVBox(
		widget.NewLabelWithStyle(trf("duel.bot", botDifficultyName(d.difficulty)), fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		container.NewPadded(d.botBW.container),
	)
	boards := container.NewGridWithColumns(2, playerPanel, botPanel)

	content := container.NewBorder(
		header,
		container.NewCenter(container.NewHBox(retryBtn, backBtn)),
		nil, nil,
		container.NewCenter(boards),
	)
	win.SetContent(content)
}
