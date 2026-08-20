package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

const (
	ThemeClassic = "classic"
	ThemeStone   = "stone"
	ThemeCandy   = "candy"
)

func themeDisplayName(theme string) string {
	switch theme {
	case ThemeStone:
		return "Stein (Castle)"
	case ThemeCandy:
		return "Suessigkeiten"
	default:
		return "Klassisch"
	}
}

func cellFactory(theme string) func(float32) tileCell {
	switch theme {
	case ThemeStone:
		return newStoneCell
	case ThemeCandy:
		return newCandyCell
	default:
		return newFlatCell
	}
}

// ---------- stone: brick-textured tile ----------

var stoneMortar = color.NRGBA{R: 0x33, G: 0x30, B: 0x2a, A: 0xff}

func stoneColors(v int) (base, fg color.Color) {
	dark := color.NRGBA{R: 0x2a, G: 0x27, B: 0x22, A: 0xff}
	cream := color.NRGBA{R: 0xf2, G: 0xef, B: 0xe6, A: 0xff}
	gold := color.NRGBA{R: 0xd4, G: 0xaf, B: 0x37, A: 0xff}
	switch v {
	case 0:
		return color.NRGBA{R: 0xb5, G: 0xb0, B: 0xa6, A: 0xff}, cream
	case 2:
		return color.NRGBA{R: 0xc9, G: 0xc4, B: 0xba, A: 0xff}, dark
	case 4:
		return color.NRGBA{R: 0xa8, G: 0xa2, B: 0x97, A: 0xff}, dark
	case 8:
		return color.NRGBA{R: 0x91, G: 0x87, B: 0x78, A: 0xff}, cream
	case 16:
		return color.NRGBA{R: 0x7c, G: 0x6a, B: 0x58, A: 0xff}, cream
	case 32:
		return color.NRGBA{R: 0x6b, G: 0x54, B: 0x3c, A: 0xff}, cream
	case 64:
		return color.NRGBA{R: 0x5c, G: 0x44, B: 0x2e, A: 0xff}, cream
	case 128:
		return color.NRGBA{R: 0x6a, G: 0x6f, B: 0x76, A: 0xff}, cream
	case 256:
		return color.NRGBA{R: 0x54, G: 0x5b, B: 0x63, A: 0xff}, cream
	case 512:
		return color.NRGBA{R: 0x8a, G: 0x93, B: 0x9c, A: 0xff}, dark
	case 1024:
		return color.NRGBA{R: 0xb0, G: 0x8d, B: 0x3f, A: 0xff}, cream
	case 2048:
		return gold, dark
	default:
		return color.NRGBA{R: 0x1a, G: 0x19, B: 0x16, A: 0xff}, gold
	}
}

type stoneCell struct {
	bricks []*canvas.Rectangle
	text   *canvas.Text
	flash  *flashOverlay
	root   fyne.CanvasObject
}

func newStoneCell(cellPx float32) tileCell {
	sizer := canvas.NewRectangle(color.Transparent)
	sizer.SetMinSize(fyne.NewSize(cellPx, cellPx))

	mortar := canvas.NewRectangle(stoneMortar)
	mortar.Resize(fyne.NewSize(cellPx, cellPx))
	mortar.Move(fyne.NewPos(0, 0))

	gap := cellPx * 0.05
	rowH := (cellPx - gap*3) / 2
	w1 := (cellPx - gap*3) / 2
	half := (w1 - gap) / 2
	y2 := gap*2 + rowH

	mk := func(x, y, w, h float32) *canvas.Rectangle {
		r := canvas.NewRectangle(color.Transparent)
		r.Resize(fyne.NewSize(w, h))
		r.Move(fyne.NewPos(x, y))
		return r
	}

	bricks := []*canvas.Rectangle{
		mk(gap, gap, w1, rowH),
		mk(gap*2+w1, gap, w1, rowH),
		mk(gap, y2, half, rowH),
		mk(gap*2+half, y2, w1, rowH),
		mk(gap*3+half+w1, y2, half, rowH),
	}

	text := canvas.NewText("", color.White)
	text.Alignment = fyne.TextAlignCenter
	text.TextStyle = fyne.TextStyle{Bold: true}
	text.TextSize = cellPx * 0.30

	flash := newFlashOverlay()

	layerObjs := append([]fyne.CanvasObject{mortar}, brickObjs(bricks)...)
	layer := container.NewWithoutLayout(layerObjs...)

	root := container.NewStack(sizer, layer, container.NewCenter(text), flash.object())
	return &stoneCell{bricks: bricks, text: text, flash: flash, root: root}
}

func brickObjs(bricks []*canvas.Rectangle) []fyne.CanvasObject {
	objs := make([]fyne.CanvasObject, len(bricks))
	for i, b := range bricks {
		objs[i] = b
	}
	return objs
}

func (c *stoneCell) object() fyne.CanvasObject { return c.root }

func (c *stoneCell) setValue(v int) {
	base, fg := stoneColors(v)
	baseN := toNRGBA(base)
	for _, b := range c.bricks {
		b.FillColor = baseN
		b.Refresh()
	}
	if v == 0 {
		c.text.Text = ""
	} else {
		c.text.Text = fmt.Sprintf("%d", v)
	}
	c.text.Color = fg
	c.text.Refresh()
}

func (c *stoneCell) flashPulse() { c.flash.trigger() }

// ---------- candy: lollipop-swirl tile ----------

func candyColors(v int) (base, fg color.Color) {
	dark := color.NRGBA{R: 0x3e, G: 0x27, B: 0x23, A: 0xff}
	white := color.White
	switch v {
	case 0:
		return color.NRGBA{R: 0xf3, G: 0xe0, B: 0xea, A: 0xff}, white
	case 2:
		return color.NRGBA{R: 0xff, G: 0xd1, B: 0xdc, A: 0xff}, dark
	case 4:
		return color.NRGBA{R: 0xe0, G: 0xb0, B: 0xff, A: 0xff}, dark
	case 8:
		return color.NRGBA{R: 0xff, G: 0x9e, B: 0xcf, A: 0xff}, white
	case 16:
		return color.NRGBA{R: 0xff, G: 0xb3, B: 0x47, A: 0xff}, white
	case 32:
		return color.NRGBA{R: 0xff, G: 0xf2, B: 0x75, A: 0xff}, dark
	case 64:
		return color.NRGBA{R: 0xb5, G: 0xe6, B: 0x55, A: 0xff}, dark
	case 128:
		return color.NRGBA{R: 0x7e, G: 0xc8, B: 0xe3, A: 0xff}, white
	case 256:
		return color.NRGBA{R: 0xa6, G: 0x85, B: 0xe2, A: 0xff}, white
	case 512:
		return color.NRGBA{R: 0xff, G: 0x6f, B: 0x91, A: 0xff}, white
	case 1024:
		return color.NRGBA{R: 0xff, G: 0xd7, B: 0x00, A: 0xff}, dark
	case 2048:
		return color.NRGBA{R: 0xff, G: 0x3e, B: 0xa5, A: 0xff}, white
	default:
		return dark, color.NRGBA{R: 0xff, G: 0xd7, B: 0x00, A: 0xff}
	}
}

type candyCell struct {
	base  *canvas.Circle
	rings []*canvas.Circle
	text  *canvas.Text
	flash *flashOverlay
	root  fyne.CanvasObject
}

func newCandyCell(cellPx float32) tileCell {
	sizer := canvas.NewRectangle(color.Transparent)
	sizer.SetMinSize(fyne.NewSize(cellPx, cellPx))

	pad := cellPx * 0.06
	d := cellPx - pad*2

	mkCircle := func(inset float32) *canvas.Circle {
		circ := canvas.NewCircle(color.Transparent)
		sz := d - inset*2
		circ.Resize(fyne.NewSize(sz, sz))
		circ.Move(fyne.NewPos(pad+inset, pad+inset))
		return circ
	}

	base := mkCircle(0)
	ring1 := mkCircle(d * 0.16)
	ring2 := mkCircle(d * 0.32)
	ring3 := mkCircle(d * 0.48)
	rings := []*canvas.Circle{ring1, ring2, ring3}

	text := canvas.NewText("", color.White)
	text.Alignment = fyne.TextAlignCenter
	text.TextStyle = fyne.TextStyle{Bold: true}
	text.TextSize = cellPx * 0.28

	flash := newFlashOverlay()

	layer := container.NewWithoutLayout(base, ring1, ring2, ring3)
	root := container.NewStack(sizer, layer, container.NewCenter(text), flash.object())
	return &candyCell{base: base, rings: rings, text: text, flash: flash, root: root}
}

func (c *candyCell) object() fyne.CanvasObject { return c.root }

func (c *candyCell) setValue(v int) {
	base, fg := candyColors(v)
	baseN := toNRGBA(base)
	c.base.FillColor = baseN
	c.base.Refresh()

	if v == 0 {
		for _, r := range c.rings {
			r.FillColor = color.Transparent
			r.Refresh()
		}
		c.text.Text = ""
	} else {
		white := color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
		alt := []color.NRGBA{white, baseN, white}
		for i, r := range c.rings {
			r.FillColor = alt[i]
			r.Refresh()
		}
		c.text.Text = fmt.Sprintf("%d", v)
	}
	c.text.Color = fg
	c.text.Refresh()
}

func (c *candyCell) flashPulse() { c.flash.trigger() }
