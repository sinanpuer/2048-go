package main

import "testing"

func TestGeneratePuzzleLevelsCount(t *testing.T) {
	levels := generatePuzzleLevels(100)
	if len(levels) != 100 {
		t.Fatalf("expected 100 puzzle levels, got %d", len(levels))
	}
}

func TestPuzzleLevelsAreDeterministic(t *testing.T) {
	a := generatePuzzleLevel(37, 100)
	b := generatePuzzleLevel(37, 100)
	if a.StartBoard != b.StartBoard {
		t.Fatalf("expected replaying level 37 to give the same starting board, got %+v vs %+v", a.StartBoard, b.StartBoard)
	}
	if a.MoveLimit != b.MoveLimit || a.ScoreGoal != b.ScoreGoal {
		t.Fatalf("expected identical goals across regenerations of the same level")
	}
}

func TestPuzzleLevel1IsEasy(t *testing.T) {
	l1 := generatePuzzleLevel(1, 100)
	filled := 0
	for r := 0; r < size; r++ {
		for c := 0; c < size; c++ {
			if l1.StartBoard[r][c] != 0 {
				filled++
			}
		}
	}
	if filled == 0 || filled > 3 {
		t.Errorf("expected level 1 to start with a sparse board (1-3 tiles), got %d", filled)
	}
	if l1.ScoreGoal > 100 {
		t.Errorf("expected an easy score goal at level 1, got %d", l1.ScoreGoal)
	}
	if l1.MoveLimit < 10 {
		t.Errorf("expected a generous move budget at level 1, got %d", l1.MoveLimit)
	}
}

func TestPuzzleDifficultyIsNonDecreasing(t *testing.T) {
	levels := generatePuzzleLevels(100)
	lastScore := 0
	for _, l := range levels {
		if l.ScoreGoal < lastScore {
			t.Errorf("level %d: score goal decreased (%d < %d)", l.Number, l.ScoreGoal, lastScore)
		}
		lastScore = l.ScoreGoal
	}
}

func TestPuzzleProgressUnlockLogic(t *testing.T) {
	p := &Progress{Completed: map[int]bool{}, PuzzleCompleted: map[int]bool{}}
	if !p.isPuzzleUnlocked(1) {
		t.Fatal("puzzle 1 must always be unlocked")
	}
	if p.isPuzzleUnlocked(2) {
		t.Fatal("puzzle 2 should be locked before puzzle 1 is completed")
	}
	p.markPuzzleCompleted(1)
	if !p.isPuzzleUnlocked(2) {
		t.Fatal("puzzle 2 should unlock after puzzle 1 is completed")
	}
}
