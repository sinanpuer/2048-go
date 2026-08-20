# Merge Kingdom (2048)

Ein 2048-Klon in Go: eine grafische Version mit Fyne ("Merge Kingdom") und eine
Terminal-Version (CLI).

## Spielen unter Windows (kein Go, kein Compiler nötig)

1. Geh auf die [Releases-Seite](../../releases) dieses Repos.
2. Lade unter dem neuesten Release die Datei **`MergeKingdom.exe`** herunter.
3. Doppelklick auf die `.exe` — fertig, kein Setup, keine Installation, keine
   weiteren Downloads nötig.

Jedes Release wird automatisch per GitHub Actions aus dem Quellcode gebaut
(siehe [`.github/workflows/build-windows.yml`](.github/workflows/build-windows.yml)).
Es gibt dort auch eine `2048-cli.exe` für die Terminal-Version.

Falls es noch kein Release gibt: Im Tab **Actions** den Workflow *"Build
Windows executable"* manuell ausführen (*Run workflow*) — danach liegt die
`.exe` als Artefakt am Ende des Laufs zum Download bereit.

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
