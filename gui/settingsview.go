package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func buildSettings(win fyne.Window) fyne.CanvasObject {
	title := widget.NewLabelWithStyle(tr("settings.title"), fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	animCheck := widget.NewCheck(tr("settings.dynamik"), nil)
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

	fps30, fps60, fpsUnlim := tr("settings.fps30"), tr("settings.fps60"), tr("settings.fpsUnlim")
	fpsGroup := widget.NewRadioGroup([]string{fps30, fps60, fpsUnlim}, nil)
	switch settings.FPS {
	case 30:
		fpsGroup.Selected = fps30
	case 0:
		fpsGroup.Selected = fpsUnlim
	default:
		fpsGroup.Selected = fps60
	}

	langDE, langEN := tr("settings.langDE"), tr("settings.langEN")
	langGroup := widget.NewRadioGroup([]string{langDE, langEN}, nil)
	if settings.Language == LangEN {
		langGroup.Selected = langEN
	} else {
		langGroup.Selected = langDE
	}

	saveBtn := widget.NewButton(tr("settings.save"), func() {
		settings.Animations = animCheck.Checked
		settings.Theme = selectedTheme

		switch fpsGroup.Selected {
		case fps30:
			settings.FPS = 30
		case fpsUnlim:
			settings.FPS = 0
		default:
			settings.FPS = 60
		}

		if langGroup.Selected == langEN {
			settings.Language = LangEN
		} else {
			settings.Language = LangDE
		}

		settings.save()
		win.SetContent(buildMenu(win))
	})
	backBtn := widget.NewButton(tr("settings.back"), func() {
		win.SetContent(buildMenu(win))
	})

	content := container.NewVBox(
		title,
		widget.NewSeparator(),
		animCheck,
		widget.NewSeparator(),
		widget.NewLabel(tr("settings.theme")),
		previews,
		widget.NewSeparator(),
		widget.NewLabel(tr("settings.fps")),
		fpsGroup,
		widget.NewSeparator(),
		widget.NewLabel(tr("settings.language")),
		langGroup,
		widget.NewSeparator(),
		container.NewCenter(container.NewHBox(backBtn, saveBtn)),
	)
	return container.NewCenter(container.NewPadded(content))
}
