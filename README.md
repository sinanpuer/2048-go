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

### Mit einem Freund spielen (Battle Royale / Party)

Nur **eine Person** braucht die `.exe` (den "Host"), die andere(n) können
über den normalen Browser mitspielen — solange ihr im selben WLAN seid
(z. B. Schul-WLAN reicht):

1. Host öffnet `MergeKingdom.exe` → **"Battle Royale (Party)"** →
   **"Party erstellen"**. Es erscheint ein Link (z. B.
   `http://192.168.x.x:8787`).
2. Diesen Link an die anderen schicken (WhatsApp, AirDrop, vorlesen, …).
3. Alle anderen öffnen den Link einfach in ihrem normalen Browser (auch am
   Schul-PC, ganz ohne Download) — fertig, alle spielen zusammen.

Hinweis: Der Host-PC muss beim ersten Start evtl. die **Windows-Firewall**
für eingehende Verbindungen freigeben (Windows fragt automatisch danach).

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
