package main

import "math"

// comboMultiplier rewards consecutive merging moves. A move that fails to
// merge anything resets the streak; each further merging move raises the
// multiplier applied to the points earned by that move.
func comboMultiplier(combo int) float64 {
	switch {
	case combo >= 11:
		return 2.0
	case combo >= 9:
		return 1.8
	case combo >= 7:
		return 1.5
	case combo >= 5:
		return 1.3
	case combo >= 3:
		return 1.1
	default:
		return 1.0
	}
}

// applyCombo advances the combo streak based on whether the move that just
// happened produced a merge (gained > 0) and returns the new combo count
// together with the (possibly boosted) points to add to the score.
func applyCombo(combo, gained int) (newCombo, boostedGained int) {
	if gained <= 0 {
		return 0, gained
	}
	newCombo = combo + 1
	boostedGained = int(math.Round(float64(gained) * comboMultiplier(newCombo)))
	return newCombo, boostedGained
}
