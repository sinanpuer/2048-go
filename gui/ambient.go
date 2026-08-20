package main

import (
	"image/color"
	"math/rand"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
)

// ambientBoard runs a self-playing 2048 game purely for visual flavor. It
// picks a random valid move on a timer, applies the same game logic as real
// gameplay, and resets once stuck.
type ambientBoard struct {
	bw      *boardWidget
	board   Board
	stopCh  chan struct{}
	running bool
}

func newAmbientBoard(cellPx float32) *ambientBoard {
	ab := &ambientBoard{bw: newBoardWidget(cellPx), board: newBoard(ModeRandomizer)}
	ab.bw.render(ab.board)
	return ab
}

func (ab *ambientBoard) canvasObject() fyne.CanvasObject {
	return ab.bw.container
}

func (ab *ambientBoard) start() {
	if ab.running {
		return
	}
	ab.running = true
	ab.stopCh = make(chan struct{})
	stopCh := ab.stopCh

	go func() {
		moves := []func(Board) (Board, bool, int){moveUp, moveDown, moveLeft, moveRight}
		interval := time.Duration(settings.ambientTickIntervalMs()) * time.Millisecond
		// A small random offset per board keeps a whole wall of these from
		// visibly ticking in lockstep.
		time.Sleep(time.Duration(rand.Intn(int(interval))))
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-stopCh:
				return
			case <-ticker.C:
				order := rand.Perm(4)
				var next Board
				moved := false
				for _, idx := range order {
					candidate, ok, _ := moves[idx](ab.board)
					if ok {
						next = candidate
						moved = true
						break
					}
				}
				if moved {
					spawnTile(&next, ModeRandomizer)
				} else {
					next = newBoard(ModeRandomizer)
				}

				fyne.Do(func() {
					ab.board = next
					ab.bw.render(ab.board)
				})
			}
		}
	}()
}

func (ab *ambientBoard) stop() {
	if !ab.running {
		return
	}
	ab.running = false
	close(ab.stopCh)
}

// ambientWall tiles many small, independently self-playing ambient boards
// into a wallpaper-like animated backdrop. Each tile stays at its fixed
// natural size (via container.NewCenter), so it never gets stretched by the
// enclosing layout - the stone and candy themes rely on fixed pixel
// geometry and render incorrectly if stretched.
type ambientWall struct {
	boards []*ambientBoard
	root   fyne.CanvasObject
}

func newAmbientWall(cols, rows int, cellPx float32) *ambientWall {
	aw := &ambientWall{}
	grid := container.New(layout.NewGridLayoutWithColumns(cols))
	for i := 0; i < cols*rows; i++ {
		ab := newAmbientBoard(cellPx)
		aw.boards = append(aw.boards, ab)
		grid.Add(container.NewCenter(ab.canvasObject()))
	}

	backdrop := canvas.NewRectangle(ambientBackdropColor(settings.Theme))
	aw.root = container.NewStack(backdrop, grid)
	return aw
}

func (aw *ambientWall) canvasObject() fyne.CanvasObject { return aw.root }

func (aw *ambientWall) start() {
	for _, b := range aw.boards {
		b.start()
	}
}

func (aw *ambientWall) stop() {
	for _, b := range aw.boards {
		b.stop()
	}
}

func ambientBackdropColor(theme string) color.NRGBA {
	switch theme {
	case ThemeStone:
		return color.NRGBA{R: 0x26, G: 0x24, B: 0x1f, A: 0xff}
	case ThemeCandy:
		return color.NRGBA{R: 0x2b, G: 0x1f, B: 0x26, A: 0xff}
	default:
		return color.NRGBA{R: 0x2b, G: 0x28, B: 0x25, A: 0xff}
	}
}
