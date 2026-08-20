package main

import "testing"

func TestGenerateLevelsMilestones(t *testing.T) {
	levels := generateLevels(100)
	if len(levels) != 100 {
		t.Fatalf("expected 100 levels, got %d", len(levels))
	}
	if levels[49].Number != 50 || levels[49].Kind != KindTile || levels[49].TileGoal != 2048 {
		t.Fatalf("level 50 should require tile 2048, got %+v", levels[49])
	}
	if levels[99].Number != 100 || levels[99].Kind != KindTile || levels[99].TileGoal != 4096 {
		t.Fatalf("level 100 should require tile 4096, got %+v", levels[99])
	}
}

func TestGenerateLevelsDifficultyIsNonDecreasing(t *testing.T) {
	levels := generateLevels(100)

	lastScore := 0
	lastTimedScore := 0
	lastTile := 0
	lastTimedTile := 0

	for _, l := range levels {
		switch l.Kind {
		case KindScore:
			if l.ScoreGoal < lastScore {
				t.Errorf("level %d: score goal decreased (%d < %d)", l.Number, l.ScoreGoal, lastScore)
			}
			lastScore = l.ScoreGoal
		case KindTile:
			if l.TileGoal < lastTile {
				t.Errorf("level %d: tile goal decreased (%d < %d)", l.Number, l.TileGoal, lastTile)
			}
			lastTile = l.TileGoal
		case KindTimedScore:
			if l.ScoreGoal < lastTimedScore {
				t.Errorf("level %d: timed score goal decreased (%d < %d)", l.Number, l.ScoreGoal, lastTimedScore)
			}
			lastTimedScore = l.ScoreGoal
		case KindTimedTile:
			if l.TileGoal < lastTimedTile {
				t.Errorf("level %d: timed tile goal decreased (%d < %d)", l.Number, l.TileGoal, lastTimedTile)
			}
			lastTimedTile = l.TileGoal
		}
		if l.Number <= 1 {
			continue
		}
	}

	if lastTile != 4096 {
		t.Errorf("expected final tile-kind level to reach 4096, got %d", lastTile)
	}
}

func TestLevel1IsEasy(t *testing.T) {
	levels := generateLevels(100)
	l1 := levels[0]
	switch l1.Kind {
	case KindScore:
		if l1.ScoreGoal > 100 {
			t.Errorf("level 1 score goal too high for an easy start: %d", l1.ScoreGoal)
		}
	case KindTile:
		if l1.TileGoal > 8 {
			t.Errorf("level 1 tile goal too high for an easy start: %d", l1.TileGoal)
		}
	}
}

func TestProgressUnlockLogic(t *testing.T) {
	p := &Progress{Completed: map[int]bool{}}
	if !p.isUnlocked(1) {
		t.Fatal("level 1 must always be unlocked")
	}
	if p.isUnlocked(2) {
		t.Fatal("level 2 should be locked before level 1 is completed")
	}
	p.markCompleted(1)
	if !p.isUnlocked(2) {
		t.Fatal("level 2 should unlock after level 1 is completed")
	}
	if p.isUnlocked(3) {
		t.Fatal("level 3 should still be locked")
	}
}
