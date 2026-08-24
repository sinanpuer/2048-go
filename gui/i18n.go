package main

import "fmt"

type Lang string

const (
	LangDE Lang = "de"
	LangEN Lang = "en"
)

var i18nDE = map[string]string{
	"menu.subtitle":   "Waehle einen Spielmodus",
	"menu.normal":     "Normales Spiel",
	"menu.randomizer": "Randomizer-Modus (2, 4, 8)",
	"menu.endless":    "Endlos-Modus (ueber 2048 hinaus)",
	"menu.levels":     "Level-Modus (%d Level)",
	"menu.duel":       "KI-Duell (gegen Bot)",
	"menu.puzzle":     "Raetsel-Modus (%d Raetsel)",
	"menu.settings":   "Einstellungen",
	"menu.highscores": "Bestwerte (%dx%d) — Normal: %d   Randomizer: %d   Endlos: %d",

	"mode.normal":     "Normal",
	"mode.randomizer": "Randomizer",
	"mode.endless":    "Endlos",

	"game.score":       "Punkte: %d",
	"game.mode":        "Modus: %s",
	"game.highscore":   "Highscore: %d",
	"game.combo":       "Combo x%d (%.1fx)",
	"game.newHigh":     "Neuer Highscore!",
	"game.won2048":     "2048 erreicht! Spiel weiter fuer mehr Punkte.",
	"game.back":        "Zurueck zum Menue",
	"game.gameover":    "Game Over!",
	"game.restart":     "Neustart",
	"game.mainmenu":    "Hauptmenue",
	"game.gameoverMsg": "Game Over! Kein Zug mehr moeglich.",

	"level.title":          "Level-Modus",
	"level.hint":           "Schliesse ein Level ab, um das naechste freizuschalten.",
	"level.back":           "Zurueck zum Menue",
	"level.backSelect":     "Zurueck zur Levelauswahl",
	"level.number":         "Level %d",
	"level.time":           "Zeit: %ds",
	"level.timeUp":         "Zeit abgelaufen!",
	"level.win":            "Level geschafft!",
	"level.retry":          "Nochmal versuchen",
	"level.retryWin":       "Nochmal spielen",
	"level.next":           "Naechstes Level",
	"level.goalScore":      "Erreiche %d Punkte",
	"level.goalTile":       "Erreiche die Kachel %d",
	"level.goalTimedScore": "%d Punkte in %ds",
	"level.goalTimedTile":  "Kachel %d in %ds",

	"puzzle.title":      "Raetsel-Modus",
	"puzzle.hint":       "Erreiche die Zielpunktzahl, bevor die Zuege ausgehen.",
	"puzzle.backSelect": "Zurueck zur Raetselauswahl",
	"puzzle.number":     "Raetsel %d",
	"puzzle.score":      "Punkte: %d / %d",
	"puzzle.moves":      "Zuege: %d / %d",
	"puzzle.moveLimit":  "Zuglimit erreicht!",
	"puzzle.win":        "Raetsel geloest!",
	"puzzle.next":       "Naechstes Raetsel",
	"puzzle.goal":       "Erreiche %d Punkte in %d Zuegen",

	"duel.title":         "KI-Duell",
	"duel.desc":          "Du spielst live gegen einen Bot auf einem eigenen Brett. Wer am Ende mehr Punkte hat, gewinnt.",
	"duel.duration":      "Rennlaenge:",
	"duel.bullet":        "Bullet (1 Minute)",
	"duel.normalLength":  "Normal (5 Minuten)",
	"duel.difficulty":    "KI-Schwierigkeit:",
	"duel.start":         "Duell starten",
	"duel.back":          "Zurueck zum Menue",
	"duel.you":           "Du",
	"duel.bot":           "Bot (%s)",
	"duel.cancel":        "Abbrechen",
	"duel.gameOverStuck": "Game Over — kein Zug mehr!",
	"duel.finished":      "KI-Duell beendet",
	"duel.youWin":        "Du gewinnst!",
	"duel.botWins":       "Der Bot gewinnt.",
	"duel.draw":          "Unentschieden!",
	"duel.scoreLine":     "Du: %d   Bot: %d",
	"duel.retry":         "Nochmal",

	"bot.normal": "Normal",
	"bot.expert": "Experte",

	"settings.title":     "Einstellungen",
	"settings.dynamik":   "Dynamik (Animationen)",
	"settings.boardsize": "Feldgroesse (freies Spiel)",
	"settings.theme":     "Design-Thema",
	"settings.fps":       "Bildwiederholrate",
	"settings.fps30":     "30 FPS",
	"settings.fps60":     "60 FPS",
	"settings.fpsUnlim":  "Unbegrenzt",
	"settings.language":  "Sprache",
	"settings.langDE":    "Deutsch",
	"settings.langEN":    "English",
	"settings.save":      "Speichern",
	"settings.back":      "Zurueck",

	"theme.classic": "Klassisch",
	"theme.stone":   "Stein (Castle)",
	"theme.candy":   "Suessigkeiten",

	"menu.party":        "Battle Royale (Party)",
	"party.title":       "Battle Royale (Party)",
	"party.desc":        "Bis zu 9 Spieler treten live gegeneinander an, egal in welchem Netzwerk. Alle spielen ueber ihren Browser mit - kein Download noetig.",
	"party.create":      "Lobby erstellen",
	"party.stop":        "Lobby beenden",
	"party.stopped":     "Lobby beendet.",
	"party.startFailed": "Konnte Lobby nicht starten: %s",
	"party.running":     "Lobby erstellt. Alle Mitspieler oeffnen diesen Link im Browser:",
	"party.back":        "Zurueck zum Menue",
}

var i18nEN = map[string]string{
	"menu.subtitle":   "Choose a game mode",
	"menu.normal":     "Normal game",
	"menu.randomizer": "Randomizer mode (2, 4, 8)",
	"menu.endless":    "Endless mode (beyond 2048)",
	"menu.levels":     "Level mode (%d levels)",
	"menu.duel":       "AI duel (vs bot)",
	"menu.puzzle":     "Puzzle mode (%d puzzles)",
	"menu.settings":   "Settings",
	"menu.highscores": "High scores (%dx%d) — Normal: %d   Randomizer: %d   Endless: %d",

	"mode.normal":     "Normal",
	"mode.randomizer": "Randomizer",
	"mode.endless":    "Endless",

	"game.score":       "Score: %d",
	"game.mode":        "Mode: %s",
	"game.highscore":   "High score: %d",
	"game.combo":       "Combo x%d (%.1fx)",
	"game.newHigh":     "New high score!",
	"game.won2048":     "You reached 2048! Keep playing for more points.",
	"game.back":        "Back to menu",
	"game.gameover":    "Game over!",
	"game.restart":     "Restart",
	"game.mainmenu":    "Main menu",
	"game.gameoverMsg": "Game over! No moves left.",

	"level.title":          "Level mode",
	"level.hint":           "Finish a level to unlock the next one.",
	"level.back":           "Back to menu",
	"level.backSelect":     "Back to level select",
	"level.number":         "Level %d",
	"level.time":           "Time: %ds",
	"level.timeUp":         "Time's up!",
	"level.win":            "Level complete!",
	"level.retry":          "Try again",
	"level.retryWin":       "Play again",
	"level.next":           "Next level",
	"level.goalScore":      "Reach %d points",
	"level.goalTile":       "Reach the %d tile",
	"level.goalTimedScore": "%d points in %ds",
	"level.goalTimedTile":  "Tile %d in %ds",

	"puzzle.title":      "Puzzle mode",
	"puzzle.hint":       "Reach the target score before you run out of moves.",
	"puzzle.backSelect": "Back to puzzle select",
	"puzzle.number":     "Puzzle %d",
	"puzzle.score":      "Score: %d / %d",
	"puzzle.moves":      "Moves: %d / %d",
	"puzzle.moveLimit":  "Move limit reached!",
	"puzzle.win":        "Puzzle solved!",
	"puzzle.next":       "Next puzzle",
	"puzzle.goal":       "Reach %d points in %d moves",

	"duel.title":         "AI duel",
	"duel.desc":          "You play live against a bot on your own board. Whoever has more points at the end wins.",
	"duel.duration":      "Match length:",
	"duel.bullet":        "Bullet (1 minute)",
	"duel.normalLength":  "Normal (5 minutes)",
	"duel.difficulty":    "AI difficulty:",
	"duel.start":         "Start duel",
	"duel.back":          "Back to menu",
	"duel.you":           "You",
	"duel.bot":           "Bot (%s)",
	"duel.cancel":        "Cancel",
	"duel.gameOverStuck": "Game over — no moves left!",
	"duel.finished":      "Duel finished",
	"duel.youWin":        "You win!",
	"duel.botWins":       "The bot wins.",
	"duel.draw":          "Draw!",
	"duel.scoreLine":     "You: %d   Bot: %d",
	"duel.retry":         "Rematch",

	"bot.normal": "Normal",
	"bot.expert": "Expert",

	"settings.title":     "Settings",
	"settings.dynamik":   "Dynamics (animations)",
	"settings.boardsize": "Board size (free play)",
	"settings.theme":     "Design theme",
	"settings.fps":       "Frame rate",
	"settings.fps30":     "30 FPS",
	"settings.fps60":     "60 FPS",
	"settings.fpsUnlim":  "Unlimited",
	"settings.language":  "Language",
	"settings.langDE":    "Deutsch",
	"settings.langEN":    "English",
	"settings.save":      "Save",
	"settings.back":      "Back",

	"theme.classic": "Classic",
	"theme.stone":   "Stone (castle)",
	"theme.candy":   "Candy",

	"menu.party":        "Battle Royale (party)",
	"party.title":       "Battle Royale (party)",
	"party.desc":        "Up to 9 players compete live, on any network. Everyone plays in their browser - no download needed.",
	"party.create":      "Create lobby",
	"party.stop":        "Stop lobby",
	"party.stopped":     "Lobby stopped.",
	"party.startFailed": "Could not start lobby: %s",
	"party.running":     "Lobby created. Everyone opens this link in their browser:",
	"party.back":        "Back to menu",
}

func tr(key string) string {
	table := i18nDE
	if settings.Language == LangEN {
		table = i18nEN
	}
	if v, ok := table[key]; ok {
		return v
	}
	return key
}

func trf(key string, args ...interface{}) string {
	return fmt.Sprintf(tr(key), args...)
}
