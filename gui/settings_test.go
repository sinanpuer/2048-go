package main

import "testing"

func TestDefaultSettingsAreSane(t *testing.T) {
	s := defaultSettings()
	if !s.Animations {
		t.Error("expected animations on by default")
	}
	if s.Theme != ThemeClassic {
		t.Errorf("expected classic theme by default, got %s", s.Theme)
	}
	if s.FPS != 60 {
		t.Errorf("expected 60 FPS by default, got %d", s.FPS)
	}
}

func TestAmbientTickIntervalScalesWithFPS(t *testing.T) {
	cases := map[int]int{30: 900, 60: 650, 0: 400}
	for fps, want := range cases {
		s := &Settings{FPS: fps}
		if got := s.ambientTickIntervalMs(); got != want {
			t.Errorf("fps=%d: got interval %dms, want %dms", fps, got, want)
		}
	}
}

func TestThemeCellFactoryProducesUsableCells(t *testing.T) {
	for _, theme := range []string{ThemeClassic, ThemeStone, ThemeCandy} {
		factory := cellFactory(theme)
		cell := factory(60)
		if cell.object() == nil {
			t.Fatalf("theme %s: cell.object() is nil", theme)
		}
		for _, v := range []int{0, 2, 4, 2048, 4096, 16384} {
			cell.setValue(v)
		}
		// flashPulse uses fyne.NewAnimation, which needs a running Fyne
		// app context and isn't exercised here; it's covered by the
		// manual runtime smoke test instead.
	}
}
