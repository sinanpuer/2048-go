package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func buildLevelSelect(win fyne.Window) fyne.CanvasObject {
	title := widget.NewLabelWithStyle("Level-Modus", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	hint := widget.NewLabel("Schliesse ein Level ab, um das naechste freizuschalten.")
	hint.Alignment = fyne.TextAlignCenter

	grid := container.New(layout.NewGridLayoutWithColumns(5))
	for _, lvl := range allLevels {
		lvl := lvl
		unlocked := progress.isUnlocked(lvl.Number)
		completed := progress.Completed[lvl.Number]

		label := fmt.Sprintf("%d", lvl.Number)
		if completed {
			label = fmt.Sprintf("%d ✓", lvl.Number)
		}

		btn := widget.NewButton(label, func() {
			win.Canvas().SetOnTypedKey(nil)
			startLevel(win, lvl.Number)
		})
		if !unlocked {
			btn.Disable()
		}
		grid.Add(btn)
	}

	scroll := container.NewVScroll(grid)
	scroll.SetMinSize(fyne.NewSize(420, 420))

	backBtn := widget.NewButton("Zurueck zum Menue", func() {
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
