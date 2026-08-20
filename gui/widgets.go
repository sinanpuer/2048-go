package main

import (
	"fmt"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
)

// tileCell abstracts a single board cell's visual so different design
// themes can render it differently (flat color, brick pattern, candy
// swirl) while the board grid, move logic, and animations stay the same
// regardless of theme.
type tileCell interface {
	object() fyne.CanvasObject
	setValue(v int)
	flashPulse()
}

// flashOverlay is a translucent layer stacked on top of a tile that briefly
// lights up to mark "something just changed here" (a spawn or a merge).
// Shared by every theme so the Dynamik setting behaves identically no
// matter which tile skin is active.
type flashOverlay struct {
	rect *canvas.Rectangle
}

func newFlashOverlay() *flashOverlay {
	return &flashOverlay{rect: canvas.NewRectangle(color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0})}
}

func (f *flashOverlay) object() fyne.CanvasObject { return f.rect }

func (f *flashOverlay) trigger() {
	anim := fyne.NewAnimation(220*time.Millisecond, func(t float32) {
		var a float32
		if t < 0.35 {
			a = t / 0.35
		} else {
			a = 1 - (t-0.35)/0.65
		}
		f.rect.FillColor = color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: uint8(160 * a)}
		f.rect.Refresh()
	})
	anim.Start()
}

func toNRGBA(c color.Color) color.NRGBA {
	return color.NRGBAModel.Convert(c).(color.NRGBA)
}

// ---------- classic: flat colored tile ----------

type flatCell struct {
	rect  *canvas.Rectangle
	text  *canvas.Text
	flash *flashOverlay
	root  fyne.CanvasObject
}

func newFlatCell(cellPx float32) tileCell {
	rect := canvas.NewRectangle(color.NRGBA{R: 0xcd, G: 0xc1, B: 0xb4, A: 0xff})
	rect.SetMinSize(fyne.NewSize(cellPx, cellPx))
	text := canvas.NewText("", color.White)
	text.Alignment = fyne.TextAlignCenter
	text.TextStyle = fyne.TextStyle{Bold: true}
	text.TextSize = cellPx * 0.32
	flash := newFlashOverlay()
	root := container.NewStack(rect, container.NewCenter(text), flash.object())
	return &flatCell{rect: rect, text: text, flash: flash, root: root}
}

func (c *flatCell) object() fyne.CanvasObject { return c.root }

func (c *flatCell) setValue(v int) {
	bg, fg := tileColor(v)
	c.rect.FillColor = bg
	c.rect.Refresh()
	if v == 0 {
		c.text.Text = ""
	} else {
		c.text.Text = fmt.Sprintf("%d", v)
	}
	c.text.Color = fg
	c.text.Refresh()
}

func (c *flatCell) flashPulse() { c.flash.trigger() }

// ---------- shared board grid ----------

// boardWidget renders an n x n grid of themed tiles and can be reused by the
// interactive game screens and the decorative ambient background alike.
// The theme is fixed at construction time from the current settings. Every
// mode except free play always constructs this with n == size (4); free
// play passes the player's chosen board size instead.
type boardWidget struct {
	cells     [][]tileCell
	last      Board
	container fyne.CanvasObject
}

func newBoardWidget(cellPx float32, n int) *boardWidget {
	bw := &boardWidget{last: newEmptyBoard(n)}
	factory := cellFactory(settings.Theme)
	grid := container.New(layout.NewGridLayoutWithColumns(n))
	bw.cells = make([][]tileCell, n)
	for r := 0; r < n; r++ {
		bw.cells[r] = make([]tileCell, n)
		for c := 0; c < n; c++ {
			cell := factory(cellPx)
			bw.cells[r][c] = cell
			grid.Add(cell.object())
		}
	}
	bw.container = grid
	return bw
}

// render updates every cell to match b. When Dynamik is enabled, any cell
// whose value just became a new nonzero value (a spawn or a merge) gets a
// brief highlight flash so the board reads as alive rather than static.
func (bw *boardWidget) render(b Board) {
	for r := range b {
		for c := range b[r] {
			v := b[r][c]
			prev := bw.last[r][c]
			bw.cells[r][c].setValue(v)
			if settings.Animations && v != 0 && v != prev {
				bw.cells[r][c].flashPulse()
			}
		}
	}
	last := newEmptyBoard(len(b))
	for r := range b {
		copy(last[r], b[r])
	}
	bw.last = last
}
