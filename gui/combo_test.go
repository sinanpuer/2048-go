package main

import "testing"

func TestComboMultiplierTiers(t *testing.T) {
	cases := []struct {
		combo int
		want  float64
	}{
		{0, 1.0}, {1, 1.0}, {2, 1.0},
		{3, 1.1}, {4, 1.1},
		{5, 1.3}, {6, 1.3},
		{7, 1.5}, {8, 1.5},
		{9, 1.8}, {10, 1.8},
		{11, 2.0}, {50, 2.0},
	}
	for _, c := range cases {
		if got := comboMultiplier(c.combo); got != c.want {
			t.Errorf("comboMultiplier(%d) = %.2f, want %.2f", c.combo, got, c.want)
		}
	}
}

func TestApplyComboResetsOnNoMerge(t *testing.T) {
	combo, gained := applyCombo(5, 0)
	if combo != 0 {
		t.Errorf("expected combo to reset to 0 on a non-merging move, got %d", combo)
	}
	if gained != 0 {
		t.Errorf("expected 0 points for a non-merging move, got %d", gained)
	}
}

func TestApplyComboBoostsScoreAtHighStreak(t *testing.T) {
	combo := 0
	var gained int
	for i := 0; i < 5; i++ {
		combo, gained = applyCombo(combo, 100)
	}
	if combo != 5 {
		t.Fatalf("expected combo 5 after 5 merging moves, got %d", combo)
	}
	want := int(100 * comboMultiplier(5))
	if gained != want {
		t.Errorf("expected boosted gain %d at combo 5, got %d", want, gained)
	}
	if gained <= 100 {
		t.Errorf("expected a combo bonus above the base 100 points, got %d", gained)
	}
}
