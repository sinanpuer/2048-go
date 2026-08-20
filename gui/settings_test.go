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
	if s.Language != LangDE {
		t.Errorf("expected German by default, got %s", s.Language)
	}
}

func TestTranslationSwitchesByLanguage(t *testing.T) {
	orig := settings.Language
	defer func() { settings.Language = orig }()

	settings.Language = LangDE
	if got := tr("settings.save"); got != "Speichern" {
		t.Errorf("expected German translation, got %q", got)
	}

	settings.Language = LangEN
	if got := tr("settings.save"); got != "Save" {
		t.Errorf("expected English translation, got %q", got)
	}

	if got := trf("menu.levels", 100); got != "Level mode (100 levels)" {
		t.Errorf("expected formatted English translation, got %q", got)
	}
}

func TestTranslationTablesHaveMatchingKeys(t *testing.T) {
	for k := range i18nDE {
		if _, ok := i18nEN[k]; !ok {
			t.Errorf("key %q exists in German table but not English", k)
		}
	}
	for k := range i18nEN {
		if _, ok := i18nDE[k]; !ok {
			t.Errorf("key %q exists in English table but not German", k)
		}
	}
}

func TestTranslationUnknownKeyFallsBackToKey(t *testing.T) {
	orig := settings.Language
	defer func() { settings.Language = orig }()
	settings.Language = LangDE
	if got := tr("does.not.exist"); got != "does.not.exist" {
		t.Errorf("expected unknown key to fall back to itself, got %q", got)
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
