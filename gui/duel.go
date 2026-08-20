package main

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func buildBotSetup(win fyne.Window) fyne.CanvasObject {
	title := widget.NewLabelWithStyle("KI-Duell", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	desc := widget.NewLabel("Du spielst live gegen einen Bot auf einem eigenen Brett. Wer am Ende mehr Punkte hat, gewinnt.")
	desc.Alignment = fyne.TextAlignCenter
	desc.Wrapping = fyne.TextWrapWord

	duration := 60
	difficulty := BotNormal

	durationRadio := widget.NewRadioGroup([]string{"Bullet (1 Minute)", "Normal (5 Minuten)"}, func(s string) {
		if strings.HasPrefix(s, "Normal") {
			duration = 300
		} else {
			duration = 60
		}
	})
	durationRadio.SetSelected("Bullet (1 Minute)")

	difficultyRadio := widget.NewRadioGroup([]string{"Normal", "Experte"}, func(s string) {
		if s == "Experte" {
			difficulty = BotExpert
		} else {
			difficulty = BotNormal
		}
	})
	difficultyRadio.SetSelected("Normal")

	startBtn := widget.NewButton("Duell starten", func() {
		startDuel(win, duration, difficulty)
	})
	backBtn := widget.NewButton("Zurueck zum Menue", func() {
		win.SetContent(buildMenu(win))
	})

	content := container.NewVBox(
		title, desc, widget.NewSeparator(),
		widget.NewLabel("Rennlaenge:"), durationRadio,
		widget.NewLabel("KI-Schwierigkeit:"), difficultyRadio,
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
	d.playerBW = newBoardWidget(64)
	d.botBW = newBoardWidget(64)
	d.playerScoreL = widget.NewLabel("")
	d.botScoreL = widget.NewLabel("")
	d.playerComboL = widget.NewLabel("")
	d.botComboL = widget.NewLabel("")
	d.playerStatusL = widget.NewLabel("")
	d.botStatusL = widget.NewLabel("")
	d.timeL = widget.NewLabelWithStyle(fmt.Sprintf("%ds", duration), fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	backBtn := widget.NewButton("Abbrechen", func() {
		d.stop()
		win.Canvas().SetOnTypedKey(nil)
		win.SetContent(buildBotSetup(win))
	})

	playerPanel := container.NewVBox(
		widget.NewLabelWithStyle("Du", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		d.playerScoreL,
		d.playerComboL,
		d.playerStatusL,
		container.NewPadded(d.playerBW.container),
	)
	botPanel := container.NewVBox(
		widget.NewLabelWithStyle(fmt.Sprintf("Bot (%s)", botDifficultyName(difficulty)), fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		d.botScoreL,
		d.botComboL,
		d.botStatusL,
		container.NewPadded(d.botBW.container),
	)

	boards := container.NewGridWithColumns(2, playerPanel, botPanel)

	header := container.NewVBox(
		widget.NewLabelWithStyle("KI-Duell", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
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
	d.playerScoreL.SetText(fmt.Sprintf("Punkte: %d", d.playerScore))
	d.botScoreL.SetText(fmt.Sprintf("Punkte: %d", d.botScore))

	if d.playerCombo > 0 {
		d.playerComboL.SetText(fmt.Sprintf("Combo x%d (%.1fx)", d.playerCombo, comboMultiplier(d.playerCombo)))
	} else {
		d.playerComboL.SetText("")
	}
	if d.botCombo > 0 {
		d.botComboL.SetText(fmt.Sprintf("Combo x%d (%.1fx)", d.botCombo, comboMultiplier(d.botCombo)))
	} else {
		d.botComboL.SetText("")
	}

	if d.playerStuck {
		d.playerStatusL.SetText("Game Over — kein Zug mehr!")
	} else {
		d.playerStatusL.SetText("")
	}
	if d.botStuck {
		d.botStatusL.SetText("Game Over — kein Zug mehr!")
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
		result = "Du gewinnst!"
	case d.playerScore < d.botScore:
		result = "Der Bot gewinnt."
	default:
		result = "Unentschieden!"
	}

	win := d.win
	resultLabel := widget.NewLabelWithStyle(result, fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	scoreLine := widget.NewLabel(fmt.Sprintf("Du: %d   Bot: %d", d.playerScore, d.botScore))
	scoreLine.Alignment = fyne.TextAlignCenter

	retryBtn := widget.NewButton("Nochmal", func() {
		win.Canvas().SetOnTypedKey(nil)
		startDuel(win, d.duration, d.difficulty)
	})
	backBtn := widget.NewButton("Zurueck zum Menue", func() {
		win.Canvas().SetOnTypedKey(nil)
		win.SetContent(buildMenu(win))
	})

	header := container.NewVBox(
		widget.NewLabelWithStyle("KI-Duell beendet", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		resultLabel, scoreLine, widget.NewSeparator(),
	)

	playerPanel := container.NewVBox(
		widget.NewLabelWithStyle("Du", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		container.NewPadded(d.playerBW.container),
	)
	botPanel := container.NewVBox(
		widget.NewLabelWithStyle(fmt.Sprintf("Bot (%s)", botDifficultyName(d.difficulty)), fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
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
