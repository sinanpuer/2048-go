package main

import (
	"math/rand"
	"time"

	"fyne.io/fyne/v2"
)

// ambientBoard runs a self-playing 2048 game in the background of the main
// menu, purely for visual flavor. It picks a random valid move on a timer,
// applies the same game logic as real gameplay, and resets once stuck.
type ambientBoard struct {
	bw      *boardWidget
	board   Board
	stopCh  chan struct{}
	running bool
}

func newAmbientBoard() *ambientBoard {
	ab := &ambientBoard{bw: newBoardWidget(70), board: newBoard(ModeRandomizer)}
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
		ticker := time.NewTicker(650 * time.Millisecond)
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
