package main

import "testing"

func newTestPlayer(id, name string) *partyPlayer {
	return &partyPlayer{id: id, name: name, send: make(chan []byte, 16), alive: true, connected: true}
}

func TestPartyAddPlayerEnforcesCapacityAndPhase(t *testing.T) {
	ps := newPartyServer()
	for i := 0; i < maxPartyPlayers; i++ {
		if _, err := ps.addPlayer(nil, "P"); err != nil {
			t.Fatalf("unexpected error adding player %d: %v", i, err)
		}
	}
	if _, err := ps.addPlayer(nil, "Overflow"); err == nil {
		t.Fatal("expected an error once the lobby is full")
	}

	ps.phase = partyPhasePlaying
	if _, err := ps.addPlayer(nil, "TooLate"); err == nil {
		t.Fatal("expected an error joining a match already in progress")
	}
}

func TestPartyAddPlayerFirstBecomesHost(t *testing.T) {
	ps := newPartyServer()
	p1, err := ps.addPlayer(nil, "Alice")
	if err != nil {
		t.Fatal(err)
	}
	if ps.hostID != p1.id {
		t.Errorf("expected first joiner to become host, got hostID=%s", ps.hostID)
	}
}

func TestPartyStartRequiresHostAndMinPlayers(t *testing.T) {
	ps := newPartyServer()
	p1 := newTestPlayer("p1", "Alice")
	p2 := newTestPlayer("p2", "Bob")
	ps.players = []*partyPlayer{p1}
	ps.hostID = "p1"

	ps.tryStart("p1")
	if ps.phase != partyPhaseLobby {
		t.Fatal("expected start to be rejected with fewer than 2 players")
	}

	ps.players = []*partyPlayer{p1, p2}
	ps.tryStart("p2")
	if ps.phase != partyPhaseLobby {
		t.Fatal("expected non-host start attempt to be ignored")
	}

	ps.tryStart("p1")
	if ps.phase != partyPhaseCountdown {
		t.Fatalf("expected host start with 2 players to begin countdown, got phase=%s", ps.phase)
	}
}

func TestPartyBeginMatchDealsBoardsAndResetsScores(t *testing.T) {
	ps := newPartyServer()
	p1 := newTestPlayer("p1", "Alice")
	p1.score = 999
	p2 := newTestPlayer("p2", "Bob")
	p2.connected = false
	ps.players = []*partyPlayer{p1, p2}
	ps.phase = partyPhaseCountdown

	ps.beginMatch()

	if ps.phase != partyPhasePlaying {
		t.Fatalf("expected phase=playing after beginMatch, got %s", ps.phase)
	}
	if p1.score != 0 {
		t.Errorf("expected score reset to 0, got %d", p1.score)
	}
	if !p1.alive {
		t.Error("expected a connected player to start alive")
	}
	if p2.alive {
		t.Error("expected a disconnected player to not be revived as alive")
	}
	filled := 0
	for r := 0; r < size; r++ {
		for c := 0; c < size; c++ {
			if p1.board[r][c] != 0 {
				filled++
			}
		}
	}
	if filled == 0 {
		t.Error("expected a fresh board to have starting tiles")
	}
}

func TestPartyApplyMoveNoOpDoesNotChangeState(t *testing.T) {
	ps := newPartyServer()
	p1 := newTestPlayer("p1", "Alice")
	// A fully stuck board: moving in any direction is a no-op.
	p1.board = Board{
		{2, 4, 2, 4},
		{4, 2, 4, 2},
		{2, 4, 2, 4},
		{4, 2, 4, 2},
	}
	before := p1.board
	ps.players = []*partyPlayer{p1}
	ps.phase = partyPhasePlaying

	ps.applyMove("p1", "up")

	if !boardsEqual(p1.board, before) {
		t.Error("a no-op move should not change the board")
	}
	if p1.score != 0 {
		t.Error("a no-op move should not change the score")
	}
	if !p1.alive {
		t.Error("a no-op move should not eliminate the player")
	}
}

func TestPartyApplyMoveMergeIncreasesScore(t *testing.T) {
	ps := newPartyServer()
	p1 := newTestPlayer("p1", "Alice")
	p1.board = Board{
		{2, 2, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
	}
	ps.players = []*partyPlayer{p1}
	ps.phase = partyPhasePlaying

	ps.applyMove("p1", "left")

	if p1.score != 4 {
		t.Errorf("expected merging two 2s to award 4 points, got %d", p1.score)
	}
	if !p1.alive {
		t.Error("player shouldn't be eliminated on a nearly-empty board")
	}
}

func TestPartyGameOverWhenOneRemains(t *testing.T) {
	ps := newPartyServer()
	p1 := newTestPlayer("p1", "Alice")
	p2 := newTestPlayer("p2", "Bob")
	p3 := newTestPlayer("p3", "Carol")
	ps.players = []*partyPlayer{p1, p2, p3}
	ps.phase = partyPhasePlaying

	p1.alive = false
	ps.mu.Lock()
	ps.checkGameOverLocked()
	ps.mu.Unlock()
	if ps.phase != partyPhasePlaying {
		t.Fatalf("match shouldn't end with 2 players still alive, got phase=%s", ps.phase)
	}

	p2.alive = false
	ps.mu.Lock()
	ps.checkGameOverLocked()
	ps.mu.Unlock()
	if ps.phase != partyPhaseGameOver {
		t.Fatalf("expected phase=gameover with 1 player left, got %s", ps.phase)
	}
	if ps.winnerID != "p3" {
		t.Errorf("expected p3 to be declared winner, got %q", ps.winnerID)
	}
}

// TestPartyWinnerIsHighestScoreNotJustSurvivor guards against the bug found
// during real play: the match used to declare whoever was still standing
// the winner, even if another player scored far more before getting stuck.
func TestPartyWinnerIsHighestScoreNotJustSurvivor(t *testing.T) {
	ps := newPartyServer()
	highScorerEliminated := newTestPlayer("p1", "Alice")
	highScorerEliminated.score = 500
	lowScorerSurvivor := newTestPlayer("p2", "Bob")
	lowScorerSurvivor.score = 50
	ps.players = []*partyPlayer{highScorerEliminated, lowScorerSurvivor}
	ps.phase = partyPhasePlaying

	highScorerEliminated.alive = false
	ps.mu.Lock()
	ps.checkGameOverLocked()
	ps.mu.Unlock()

	if ps.phase != partyPhaseGameOver {
		t.Fatalf("expected phase=gameover with 1 player left, got %s", ps.phase)
	}
	if ps.winnerID != "p1" {
		t.Errorf("expected the higher-scoring player p1 to win despite being eliminated first, got %q", ps.winnerID)
	}
}

func TestPartyRestartReturnsToLobbyWithSamePlayers(t *testing.T) {
	ps := newPartyServer()
	host := newTestPlayer("p1", "Alice")
	host.score = 300
	host.alive = false
	other := newTestPlayer("p2", "Bob")
	other.score = 50
	ps.players = []*partyPlayer{host, other}
	ps.hostID = "p1"
	ps.phase = partyPhaseGameOver
	ps.winnerID = "p1"

	ps.tryRestart("p2") // non-host must not be able to restart
	if ps.phase != partyPhaseGameOver {
		t.Fatal("expected a non-host restart attempt to be ignored")
	}

	ps.tryRestart("p1")
	if ps.phase != partyPhaseLobby {
		t.Fatalf("expected phase=lobby after host restarts, got %s", ps.phase)
	}
	if ps.winnerID != "" {
		t.Errorf("expected winnerID to be cleared, got %q", ps.winnerID)
	}
	if len(ps.players) != 2 {
		t.Fatalf("expected both connected players to carry over, got %d", len(ps.players))
	}
	if !host.alive || host.score != 0 {
		t.Errorf("expected host to be reset to alive/score 0, got alive=%v score=%d", host.alive, host.score)
	}
}

func TestPartyRestartDropsDisconnectedPlayers(t *testing.T) {
	ps := newPartyServer()
	host := newTestPlayer("p1", "Alice")
	gone := newTestPlayer("p2", "Bob")
	gone.connected = false
	ps.players = []*partyPlayer{host, gone}
	ps.hostID = "p1"
	ps.phase = partyPhaseGameOver

	ps.tryRestart("p1")

	if len(ps.players) != 1 || ps.players[0].id != "p1" {
		t.Fatalf("expected only the still-connected host to carry over, got %d players", len(ps.players))
	}
}

func TestPartyDisconnectDuringLobbyRemovesPlayerAndReassignsHost(t *testing.T) {
	ps := newPartyServer()
	p1, _ := ps.addPlayer(nil, "Alice")
	p2, _ := ps.addPlayer(nil, "Bob")

	ps.handleDisconnect(p1)

	if len(ps.players) != 1 {
		t.Fatalf("expected the disconnected player to be removed, got %d players", len(ps.players))
	}
	if ps.hostID != p2.id {
		t.Errorf("expected host to be reassigned to remaining player, got hostID=%s", ps.hostID)
	}
}

func TestPartyDisconnectDuringMatchCountsAsElimination(t *testing.T) {
	ps := newPartyServer()
	p1 := newTestPlayer("p1", "Alice")
	p2 := newTestPlayer("p2", "Bob")
	ps.players = []*partyPlayer{p1, p2}
	ps.phase = partyPhasePlaying

	ps.handleDisconnect(p1)

	if len(ps.players) != 2 {
		t.Error("expected players to stay in the roster mid-match so their final board is still visible")
	}
	if p1.alive {
		t.Error("expected a disconnected player to be marked eliminated")
	}
	if ps.phase != partyPhaseGameOver {
		t.Fatalf("expected the match to end once only p2 remains, got phase=%s", ps.phase)
	}
	if ps.winnerID != "p2" {
		t.Errorf("expected p2 to win by the other player disconnecting, got %q", ps.winnerID)
	}
}
