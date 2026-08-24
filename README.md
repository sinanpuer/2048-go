# Merge Kingdom (2048)

Ein 2048-Klon in Go: eine grafische Version mit Fyne ("Merge Kingdom") und eine
Terminal-Version (CLI).

## Im Browser spielen (kein Download, keine Installation)

**<https://sinanpuer.github.io/2048-go/>**

Einfach den Link öffnen — läuft komplett im Browser (als WebAssembly), es wird
nichts heruntergeladen oder installiert. Ideal für Rechner (z. B. in der
Schule), auf denen Downloads/Installationen blockiert sind. Benötigt einen
halbwegs aktuellen Browser mit WebGL-Unterstützung (Chrome/Edge ≥ 106,
Firefox, Safari).

Falls der Link noch nicht funktioniert: In den Repo-Settings unter **Pages**
als Quelle *"Deploy from a branch"* → Branch `main`, Ordner `/docs`
einstellen (einmalig nötig, siehe unten).

Der Build unter `docs/` wird automatisch per GitHub Actions aus `gui/` neu
gebaut, sobald sich der Code ändert (siehe
[`.github/workflows/build-wasm.yml`](.github/workflows/build-wasm.yml)).

## Mit Freunden spielen (Battle Royale / Party)

Funktioniert sowohl aus der Browser- als auch der Windows-Version, und
zwischen beiden gemischt - über das Internet, nicht nur im selben WLAN (das
war früher nötig, ist aber an vielen Netzwerken wie Schul-WLANs unzuverlässig,
weil die oft die direkte Verbindung zwischen Geräten blockieren):

1. Irgendwer öffnet **"Battle Royale (Party)"** → **"Lobby erstellen"**. Es
   erscheint ein Link.
2. Diesen Link an die anderen schicken (WhatsApp, Discord, vorlesen, …).
3. Alle öffnen den Link einfach in ihrem Browser (auch am Schul-PC, ganz ohne
   Download) — fertig, alle spielen zusammen, live gegeneinander.

Das läuft über einen kleinen extern gehosteten Relay-Server (siehe
[`relay/`](relay/)) - der leitet nur Züge/Spielstände zwischen den Spielern
weiter, läuft aber unabhängig von jedem Spieler-PC.

## Spielen unter Windows (kein Go, kein Compiler nötig)

Auch wenn Netzwerke/Filter (z. B. in der Schule) `.exe`-Downloads blockieren,
lässt sich meist eine `.zip`-Datei herunterladen:

1. Geh auf die [Releases-Seite](../../releases) dieses Repos.
2. Lade unter dem neuesten Release die Datei **`MergeKingdom.zip`** herunter
   (nicht die `.exe` direkt — die wird von vielen Filtern blockiert, die
   `.zip` meist nicht).
3. Rechtsklick auf die zip → **"Alle extrahieren"**.
4. Doppelklick auf die extrahierte `MergeKingdom.exe` — kein Setup, keine
   Installation, keine Adminrechte nötig.

Jedes Release wird automatisch per GitHub Actions aus dem Quellcode gebaut
(siehe [`.github/workflows/build-windows.yml`](.github/workflows/build-windows.yml)).
Es gibt dort auch ein `2048-cli.zip` mit der Terminal-Version.

Falls es noch kein Release gibt: Im Tab **Actions** den Workflow *"Build
Windows executable"* manuell ausführen (*Run workflow*) — danach liegen die
`.zip`-Dateien als Artefakt am Ende des Laufs zum Download bereit.

## Neues Release erstellen (für Maintainer)

```bash
git tag v1.0.0
git push origin v1.0.0
```

Das löst den Workflow aus, baut die `.exe` und hängt sie automatisch an ein
neues GitHub-Release an.

## Lokal bauen (optional, nur für Entwickler)

Voraussetzung: [Go](https://go.dev/dl/) sowie ein C-Compiler (z. B.
[MinGW-w64](https://www.mingw-w64.org/)) für die GUI-Variante.

```bash
cd gui
go build -o MergeKingdom.exe .

cd ../cli
go build -o 2048-cli.exe .
```

## Party-Relay-Server deployen (für Maintainer)

`relay/` ist ein eigenständiges Go-Modul (kein Fyne, keine GUI-Abhängigkeiten)
- läuft z. B. kostenlos auf [Render](https://render.com):

1. Auf Render: **New +** → **Web Service** → Repo `2048-go` verbinden.
2. **Root Directory**: `relay`
3. **Runtime**: Docker (das `relay/Dockerfile` wird automatisch erkannt)
4. **Instance Type**: Free
5. Deploy. Bei jedem Push auf `main`, der `relay/**` ändert, baut Render neu.

Die App (`gui/party.go`, Konstante `partyRelayHost`) muss dann auf die
tatsächliche Render-URL zeigen.
