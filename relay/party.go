// Party/Battle-Royale game logic, adapted by hand from gui/party.go to run
// as one room among many on a shared public relay instead of a single
// locally-hosted game per process.
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const maxPartyPlayers = 9
const partyCountdownSeconds = 5

var partyUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type partyPhase string

const (
	partyPhaseLobby     partyPhase = "lobby"
	partyPhaseCountdown partyPhase = "countdown"
	partyPhasePlaying   partyPhase = "playing"
	partyPhaseGameOver  partyPhase = "gameover"
)

type partyPlayer struct {
	id        string
	name      string
	conn      *websocket.Conn
	send      chan []byte
	board     Board
	score     int
	combo     int
	alive     bool
	connected bool
}

type partyPlayerView struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Board     Board  `json:"board"`
	Score     int    `json:"score"`
	Combo     int    `json:"combo"`
	Alive     bool   `json:"alive"`
	Connected bool   `json:"connected"`
}

type partyStateMsg struct {
	Type      string            `json:"type"`
	Phase     partyPhase        `json:"phase"`
	Countdown int               `json:"countdown"`
	YouID     string            `json:"youId"`
	HostID    string            `json:"hostId"`
	WinnerID  string            `json:"winnerId"`
	Players   []partyPlayerView `json:"players"`
	Error     string            `json:"error,omitempty"`
}

// partyServer holds the state of a single room/match.
type partyServer struct {
	mu                 sync.Mutex
	players            []*partyPlayer
	nextID             int
	hostID             string
	phase              partyPhase
	winnerID           string
	countdownRemaining int
	lastActivity       time.Time
}

func newPartyServer() *partyServer {
	return &partyServer{phase: partyPhaseLobby, lastActivity: time.Now()}
}

func (ps *partyServer) playerCount() int {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	return len(ps.players)
}

func (ps *partyServer) idleSince() time.Time {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	return ps.lastActivity
}

func (ps *partyServer) handleWS(conn *websocket.Conn) {
	var p *partyPlayer

	defer func() {
		if p != nil {
			ps.handleDisconnect(p)
		}
		conn.Close()
	}()

	for {
		_, raw, err := conn.ReadMessage()
		if err != nil {
			return
		}

		var msg struct {
			Type string `json:"type"`
			Name string `json:"name"`
			Dir  string `json:"dir"`
		}
		if err := json.Unmarshal(raw, &msg); err != nil {
			continue
		}

		ps.mu.Lock()
		ps.lastActivity = time.Now()
		ps.mu.Unlock()

		switch msg.Type {
		case "join":
			if p != nil {
				continue
			}
			p, err = ps.addPlayer(conn, msg.Name)
			if err != nil {
				conn.WriteJSON(partyStateMsg{Type: "state", Error: err.Error()})
				return
			}
			go ps.writePump(p)
			ps.broadcastState()
		case "start":
			if p != nil {
				ps.tryStart(p.id)
			}
		case "move":
			if p != nil {
				ps.applyMove(p.id, msg.Dir)
			}
		}
	}
}

func (ps *partyServer) writePump(p *partyPlayer) {
	for data := range p.send {
		if err := p.conn.WriteMessage(websocket.TextMessage, data); err != nil {
			return
		}
	}
}

func (ps *partyServer) addPlayer(conn *websocket.Conn, name string) (*partyPlayer, error) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if ps.phase != partyPhaseLobby {
		return nil, fmt.Errorf("das Spiel hat bereits begonnen")
	}
	if len(ps.players) >= maxPartyPlayers {
		return nil, fmt.Errorf("die Lobby ist voll (max %d Spieler)", maxPartyPlayers)
	}
	if name == "" {
		name = fmt.Sprintf("Spieler %d", len(ps.players)+1)
	}

	ps.nextID++
	p := &partyPlayer{
		id:        fmt.Sprintf("p%d", ps.nextID),
		name:      name,
		conn:      conn,
		send:      make(chan []byte, 16),
		alive:     true,
		connected: true,
	}
	ps.players = append(ps.players, p)
	if ps.hostID == "" {
		ps.hostID = p.id
	}
	return p, nil
}

func (ps *partyServer) handleDisconnect(p *partyPlayer) {
	ps.mu.Lock()
	p.connected = false
	if ps.phase == partyPhaseLobby {
		for i, pl := range ps.players {
			if pl.id == p.id {
				ps.players = append(ps.players[:i], ps.players[i+1:]...)
				break
			}
		}
		if ps.hostID == p.id {
			ps.hostID = ""
			if len(ps.players) > 0 {
				ps.hostID = ps.players[0].id
			}
		}
	} else {
		p.alive = false
	}
	close(p.send)
	ps.checkGameOverLocked()
	ps.mu.Unlock()
	ps.broadcastState()
}

func (ps *partyServer) tryStart(byID string) {
	ps.mu.Lock()
	if ps.phase != partyPhaseLobby || byID != ps.hostID || len(ps.players) < 2 {
		ps.mu.Unlock()
		return
	}
	ps.phase = partyPhaseCountdown
	ps.mu.Unlock()
	ps.broadcastState()

	go ps.runCountdown()
}

func (ps *partyServer) runCountdown() {
	for s := partyCountdownSeconds; s > 0; s-- {
		ps.mu.Lock()
		ps.countdownRemaining = s
		ps.mu.Unlock()
		ps.broadcastState()
		time.Sleep(1 * time.Second)
	}
	ps.beginMatch()
}

func (ps *partyServer) beginMatch() {
	ps.mu.Lock()
	for _, p := range ps.players {
		p.board = newBoard(ModeRandomizer)
		p.score = 0
		p.combo = 0
		p.alive = p.connected
	}
	ps.phase = partyPhasePlaying
	ps.mu.Unlock()
	ps.broadcastState()
}

func (ps *partyServer) applyMove(byID, dir string) {
	var fn func(Board) (Board, bool, int)
	switch dir {
	case "up":
		fn = moveUp
	case "down":
		fn = moveDown
	case "left":
		fn = moveLeft
	case "right":
		fn = moveRight
	default:
		return
	}

	ps.mu.Lock()
	if ps.phase != partyPhasePlaying {
		ps.mu.Unlock()
		return
	}
	var target *partyPlayer
	for _, p := range ps.players {
		if p.id == byID {
			target = p
			break
		}
	}
	if target == nil || !target.alive {
		ps.mu.Unlock()
		return
	}

	nb, moved, gained := fn(target.board)
	if !moved {
		ps.mu.Unlock()
		return
	}
	target.board = nb
	target.combo, gained = applyCombo(target.combo, gained)
	target.score += gained
	spawnTile(&target.board, ModeRandomizer)
	if !hasMoves(target.board) {
		target.alive = false
	}
	ps.checkGameOverLocked()
	ps.mu.Unlock()

	ps.broadcastState()
}

func (ps *partyServer) checkGameOverLocked() {
	if ps.phase != partyPhasePlaying {
		return
	}
	aliveCount := 0
	for _, p := range ps.players {
		if p.alive {
			aliveCount++
		}
	}
	if aliveCount > 1 {
		return
	}
	ps.phase = partyPhaseGameOver
	// The match ends once at most one player is still going, but the
	// winner is whoever scored the most overall - surviving longer isn't
	// worth anything on its own if someone else racked up a higher score
	// before getting stuck. Ties go to whoever is still alive.
	var winner *partyPlayer
	for _, p := range ps.players {
		switch {
		case winner == nil:
			winner = p
		case p.score > winner.score:
			winner = p
		case p.score == winner.score && p.alive && !winner.alive:
			winner = p
		}
	}
	if winner != nil {
		ps.winnerID = winner.id
	}
}

func (ps *partyServer) broadcastState() {
	ps.mu.Lock()
	players := make([]partyPlayerView, len(ps.players))
	for i, p := range ps.players {
		players[i] = partyPlayerView{
			ID: p.id, Name: p.name, Board: p.board,
			Score: p.score, Combo: p.combo, Alive: p.alive, Connected: p.connected,
		}
	}
	base := partyStateMsg{
		Type:      "state",
		Phase:     ps.phase,
		Countdown: ps.countdownRemaining,
		HostID:    ps.hostID,
		WinnerID:  ps.winnerID,
		Players:   players,
	}
	targets := make([]*partyPlayer, len(ps.players))
	copy(targets, ps.players)
	ps.mu.Unlock()

	for _, p := range targets {
		if !p.connected {
			continue
		}
		msg := base
		msg.YouID = p.id
		data, err := json.Marshal(msg)
		if err != nil {
			continue
		}
		select {
		case p.send <- data:
		default:
		}
	}
}
