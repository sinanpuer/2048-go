package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func buildPuzzleSelect(win fyne.Window) fyne.CanvasObject {
	title := widget.NewLabelWithStyle(tr("puzzle.title"), fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	hint := widget.NewLabel(tr("puzzle.hint"))
	hint.Alignment = fyne.TextAlignCenter

	grid := container.New(layout.NewGridLayoutWithColumns(5))
	for _, pl := range allPuzzleLevels {
		pl := pl
		unlocked := progress.isPuzzleUnlocked(pl.Number)
		completed := progress.PuzzleCompleted[pl.Number]

		label := fmt.Sprintf("%d", pl.Number)
		if completed {
			label = fmt.Sprintf("%d ✓", pl.Number)
		}

		btn := widget.NewButton(label, func() {
			win.Canvas().SetOnTypedKey(nil)
			startPuzzle(win, pl.Number)
		})
		if !unlocked {
			btn.Disable()
		}
		grid.Add(btn)
	}

	scroll := container.NewVScroll(grid)
	scroll.SetMinSize(fyne.NewSize(420, 420))

	backBtn := widget.NewButton(tr("level.back"), func() {
		win.SetContent(buildMenu(win))
	})

	content := container.NewBorder(
		container.NewVBox(title, hint, widget.NewSeparator()),
		container.NewCenter(backBtn),
		nil, nil,
		container.NewCenter(scroll),
	)
	return content
}
