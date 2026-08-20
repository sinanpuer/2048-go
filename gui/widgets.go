package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
)

// boardWidget renders a 4x4 grid of colored tiles and can be reused by the
// interactive game screens and the decorative ambient background alike.
type boardWidget struct {
	rects     [size][size]*canvas.Rectangle
	texts     [size][size]*canvas.Text
	container fyne.CanvasObject
}

func newBoardWidget(cellPx float32) *boardWidget {
	bw := &boardWidget{}
	grid := container.New(layout.NewGridLayoutWithColumns(size))
	for r := 0; r < size; r++ {
		for c := 0; c < size; c++ {
			rect := canvas.NewRectangle(color.NRGBA{R: 0xcd, G: 0xc1, B: 0xb4, A: 0xff})
			rect.SetMinSize(fyne.NewSize(cellPx, cellPx))
			text := canvas.NewText("", color.White)
			text.Alignment = fyne.TextAlignCenter
			text.TextStyle = fyne.TextStyle{Bold: true}
			text.TextSize = cellPx * 0.32
			bw.rects[r][c] = rect
			bw.texts[r][c] = text
			cell := container.NewStack(rect, container.NewCenter(text))
			grid.Add(cell)
		}
	}
	bw.container = grid
	return bw
}

func (bw *boardWidget) render(b Board) {
	for r := 0; r < size; r++ {
		for c := 0; c < size; c++ {
			v := b[r][c]
			bg, fg := tileColor(v)
			bw.rects[r][c].FillColor = bg
			bw.rects[r][c].Refresh()
			t := bw.texts[r][c]
			if v == 0 {
				t.Text = ""
			} else {
				t.Text = fmt.Sprintf("%d", v)
			}
			t.Color = fg
			t.Refresh()
		}
	}
}
