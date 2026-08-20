package main

import "testing"

func TestBotPickMoveOnFullBoardWithNoMoves(t *testing.T) {
	// A board with no possible merges and no empty cells: no legal move exists.
	b := Board{
		{2, 4, 2, 4},
		{4, 2, 4, 2},
		{2, 4, 2, 4},
		{4, 2, 4, 2},
	}
	if _, ok := botPickMove(b, BotExpert); ok {
		t.Fatal("expected no legal move on a fully stuck board")
	}
}

func TestBotPickMoveAppliesCleanly(t *testing.T) {
	b := newBoard(ModeRandomizer)
	for _, d := range []BotDifficulty{BotNormal, BotExpert} {
		mv, ok := botPickMove(b, d)
		if !ok {
			t.Fatalf("expected a legal move on a fresh board for difficulty %v", d)
		}
		nb, moved, _ := mv(b)
		if !moved {
			t.Fatalf("bot picked a move that didn't actually change the board (difficulty %v)", d)
		}
		_ = nb
	}
}

// playFullGame runs a bot against itself until the board is stuck and
// returns the highest tile it reached, used to compare difficulty tiers.
func playFullGame(d BotDifficulty) int {
	b := newBoard(ModeRandomizer)
	for {
		mv, ok := botPickMove(b, d)
		if !ok {
			return highestTile(b)
		}
		nb, moved, _ := mv(b)
		if !moved {
			return highestTile(b)
		}
		spawnTile(&nb, ModeRandomizer)
		b = nb
	}
}

func TestExpertOutperformsNormalOnAverage(t *testing.T) {
	const runs = 25
	normalTotal, expertTotal := 0, 0
	for i := 0; i < runs; i++ {
		normalTotal += playFullGame(BotNormal)
		expertTotal += playFullGame(BotExpert)
	}
	normalAvg := float64(normalTotal) / runs
	expertAvg := float64(expertTotal) / runs

	if expertAvg <= normalAvg {
		t.Fatalf("expected Expert bot to reach a higher average tile than Normal, got expert=%.1f normal=%.1f", expertAvg, normalAvg)
	}
}
