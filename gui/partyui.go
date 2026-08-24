package main

import (
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func buildPartySetup(win fyne.Window) fyne.CanvasObject {
	title := widget.NewLabelWithStyle(tr("party.title"), fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	desc := widget.NewLabel(tr("party.desc"))
	desc.Alignment = fyne.TextAlignCenter
	desc.Wrapping = fyne.TextWrapWord

	statusL := widget.NewLabel("")
	statusL.Alignment = fyne.TextAlignCenter
	statusL.Wrapping = fyne.TextWrapWord

	linkBox := container.NewVBox()

	var roomOpen bool

	startBtn := widget.NewButton(tr("party.create"), nil)
	stopBtn := widget.NewButton(tr("party.stop"), nil)
	stopBtn.Hide()

	stopBtn.OnTapped = func() {
		roomOpen = false
		statusL.SetText(tr("party.stopped"))
		linkBox.RemoveAll()
		startBtn.Show()
		stopBtn.Hide()
	}

	startBtn.OnTapped = func() {
		if roomOpen {
			return
		}
		roomOpen = true
		addr := partyRoomURL(newPartyRoomCode())

		statusL.SetText(tr("party.running"))
		linkBox.RemoveAll()
		u, parseErr := url.Parse(addr)
		if parseErr == nil {
			link := widget.NewHyperlink(addr, u)
			linkBox.Add(link)
		} else {
			linkBox.Add(widget.NewLabel(addr))
		}
		startBtn.Hide()
		stopBtn.Show()
	}

	backBtn := widget.NewButton(tr("party.back"), func() {
		win.SetContent(buildMenu(win))
	})

	content := container.NewVBox(
		title, desc, widget.NewSeparator(),
		startBtn, stopBtn,
		statusL, linkBox,
		widget.NewSeparator(),
		backBtn,
	)
	return container.NewCenter(container.NewPadded(content))
}
