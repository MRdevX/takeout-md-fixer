# Takeout Metadata Fixer

Desktop app for **Google Takeout** exports: it finds JSON sidecars next to your photos and videos and writes dates, GPS, and related fields into the files with [ExifTool](https://exiftool.org/). That way imports into iCloud, a NAS, or other tools see sensible dates and locations.

You must install [ExifTool](https://exiftool.org/) yourself. The app resolves it from your `PATH` and common install locations (e.g. Homebrew and `/usr/local/bin` on macOS), which helps when the GUI does not see the same `PATH` as your terminal.

## Download

Prebuilt **macOS (DMG)** and **Windows (.exe)** builds are on [GitHub Releases](https://github.com/MRdevX/takeout-md-fixer/releases). Pushes to `main` trigger a release workflow that bumps the patch version from git tags and publishes artifacts plus a short changelog.

## Resuming

If you stop early or quit, progress is saved in a hidden file **`.takeout-md-fixer-checkpoint.json`** in the folder you selected. Next time, enable **Continue where you left off** to skip files already processed. Delete that file if you want a full pass from scratch on the same folder.

## Develop

You need Go (see [`go.mod`](go.mod)), Node.js, the [Wails v3 CLI](https://v3.wails.io/) (`go install github.com/wailsapp/wails/v3/cmd/wails3@latest`), and ExifTool.

```bash
wails3 dev -config ./build/config.yml
```

[`main.go`](main.go) embeds [`frontend/dist`](frontend/dist), which is gitignored. Before `go build`, `go vet`, or `go test` at the repo root, build the frontend (`npm ci` and `npm run build` in `frontend/`). CI does the same.

If you change exported Go methods or JSON types used from the UI, regenerate [`frontend/bindings/`](frontend/bindings/) with the Wails version in `go.mod`, then rebuild the frontend:

```bash
wails3 generate bindings -f '-tags production -trimpath -buildvcs=false -ldflags="-w -s"' -clean=true
```

## Build

```bash
wails3 build
```

Binaries land under `bin/`.

Local release script (macOS DMG; Windows cross-compile where supported):

```bash
bash scripts/build-release.sh
```

Optional: `VERSION=x.y.z`.

---

**Mahdi Rashidi** — [contact@mrashidi.me](mailto:contact@mrashidi.me) · [mrashidi.me](https://mrashidi.me) · [GitHub](https://github.com/MRdevX) · [Buy Me a Coffee](https://www.buymeacoffee.com/mrdevx)
