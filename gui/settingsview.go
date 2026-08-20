package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func buildSettings(win fyne.Window) fyne.CanvasObject {
	title := widget.NewLabelWithStyle("Einstellungen", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	animCheck := widget.NewCheck("Dynamik (Animationen)", nil)
	animCheck.SetChecked(settings.Animations)

	selectedTheme := settings.Theme
	themes := []string{ThemeClassic, ThemeStone, ThemeCandy}
	cards := make([]*themeCard, len(themes))

	selectTheme := func(theme string) {
		selectedTheme = theme
		for i, t := range themes {
			cards[i].setSelected(t == theme)
		}
	}

	for i, theme := range themes {
		theme := theme
		cell := cellFactory(theme)(52)
		cell.setValue(2048)
		preview := container.NewVBox(
			container.NewCenter(cell.object()),
			widget.NewLabel(themeDisplayName(theme)),
		)
		cards[i] = newThemeCard(preview, func() { selectTheme(theme) })
	}
	selectTheme(selectedTheme)

	previews := container.NewGridWithColumns(3, cards[0], cards[1], cards[2])

	fpsNames := []string{"30 FPS", "60 FPS", "Unbegrenzt"}
	fpsGroup := widget.NewRadioGroup(fpsNames, nil)
	switch settings.FPS {
	case 30:
		fpsGroup.Selected = "30 FPS"
	case 0:
		fpsGroup.Selected = "Unbegrenzt"
	default:
		fpsGroup.Selected = "60 FPS"
	}

	saveBtn := widget.NewButton("Speichern", func() {
		settings.Animations = animCheck.Checked
		settings.Theme = selectedTheme

		switch fpsGroup.Selected {
		case "30 FPS":
			settings.FPS = 30
		case "Unbegrenzt":
			settings.FPS = 0
		default:
			settings.FPS = 60
		}

		settings.save()
		win.SetContent(buildMenu(win))
	})
	backBtn := widget.NewButton("Zurueck", func() {
		win.SetContent(buildMenu(win))
	})

	content := container.NewVBox(
		title,
		widget.NewSeparator(),
		animCheck,
		widget.NewSeparator(),
		widget.NewLabel("Design-Thema"),
		previews,
		widget.NewSeparator(),
		widget.NewLabel("Bildwiederholrate"),
		fpsGroup,
		widget.NewSeparator(),
		container.NewCenter(container.NewHBox(backBtn, saveBtn)),
	)
	return container.NewCenter(container.NewPadded(content))
}
