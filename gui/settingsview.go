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

	themeNames := []string{themeDisplayName(ThemeClassic), themeDisplayName(ThemeStone), themeDisplayName(ThemeCandy)}
	themeGroup := widget.NewRadioGroup(themeNames, nil)
	themeGroup.Selected = themeDisplayName(settings.Theme)

	preview := func(theme string) fyne.CanvasObject {
		cell := cellFactory(theme)(52)
		cell.setValue(2048)
		return container.NewPadded(cell.object())
	}
	previews := container.NewGridWithColumns(3,
		container.NewCenter(preview(ThemeClassic)),
		container.NewCenter(preview(ThemeStone)),
		container.NewCenter(preview(ThemeCandy)),
	)

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

		switch themeGroup.Selected {
		case themeDisplayName(ThemeStone):
			settings.Theme = ThemeStone
		case themeDisplayName(ThemeCandy):
			settings.Theme = ThemeCandy
		default:
			settings.Theme = ThemeClassic
		}

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
		themeGroup,
		widget.NewSeparator(),
		widget.NewLabel("Bildwiederholrate"),
		fpsGroup,
		widget.NewSeparator(),
		container.NewCenter(container.NewHBox(backBtn, saveBtn)),
	)
	return container.NewCenter(container.NewPadded(content))
}
