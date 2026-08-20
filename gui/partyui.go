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

	var ps *partyServer

	startBtn := widget.NewButton(tr("party.create"), nil)
	stopBtn := widget.NewButton(tr("party.stop"), func() {
		if ps != nil {
			ps.stop()
			ps = nil
		}
		statusL.SetText(tr("party.stopped"))
		linkBox.RemoveAll()
	})
	stopBtn.Hide()

	startBtn.OnTapped = func() {
		if ps != nil {
			return
		}
		ps = newPartyServer()
		addr, err := ps.start()
		if err != nil {
			statusL.SetText(trf("party.startFailed", err.Error()))
			ps = nil
			return
		}

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
		if ps != nil {
			ps.stop()
			ps = nil
		}
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
