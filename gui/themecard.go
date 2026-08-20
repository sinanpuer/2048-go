package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// themeCard wraps arbitrary content (a tile preview) in a tappable card
// that highlights with a green border when selected.
type themeCard struct {
	widget.BaseWidget
	border   *canvas.Rectangle
	renderer fyne.WidgetRenderer
	onTapped func()
}

func newThemeCard(content fyne.CanvasObject, onTapped func()) *themeCard {
	border := canvas.NewRectangle(color.Transparent)
	border.StrokeWidth = 3
	border.StrokeColor = color.Transparent
	border.CornerRadius = 10

	t := &themeCard{border: border, onTapped: onTapped}
	t.ExtendBaseWidget(t)

	inner := container.NewStack(border, container.NewPadded(content))
	t.renderer = widget.NewSimpleRenderer(inner)
	return t
}

func (t *themeCard) Tapped(_ *fyne.PointEvent) {
	if t.onTapped != nil {
		t.onTapped()
	}
}

func (t *themeCard) TappedSecondary(_ *fyne.PointEvent) {}

func (t *themeCard) CreateRenderer() fyne.WidgetRenderer {
	return t.renderer
}

func (t *themeCard) setSelected(selected bool) {
	if selected {
		t.border.StrokeColor = color.NRGBA{R: 0x3d, G: 0xb8, B: 0x4a, A: 0xff}
	} else {
		t.border.StrokeColor = color.Transparent
	}
	t.border.Refresh()
}
