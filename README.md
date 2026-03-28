# Takeout Metadata Fixer

A small desktop app that reads the JSON sidecars Google Takeout ships next to your photos and videos, then writes the real dates and GPS (and related fields) back into the files using [ExifTool](https://exiftool.org/).

**You need [ExifTool](https://exiftool.org/) installed** on your system and available on your `PATH`. The app checks when it starts and will tell you if it cannot find it. If ExifTool is installed in a non-standard location, set `TAKEOUT_EXIFTOOL_PATH` to the full path of the `exiftool` binary (or `exiftool.exe` on Windows).

## Why I built this

One day I wanted to move my photos and memories out of Google Photos and used Takeout. The files came out with metadata that didn’t match what I expected—so bringing them into iCloud or anywhere else felt wrong. I wrote this tool to put that metadata back where it belongs before importing.

## Tech stack

- **Go** — backend and Wails service
- **[Wails v3](https://v3.wails.io/)** — desktop shell, bindings to the UI
- **Vite** + **vanilla JS** — frontend
- **[go-exiftool](https://github.com/barasher/go-exiftool)** — talks to ExifTool
- **[ExifTool](https://exiftool.org/)** — must be installed separately on your machine

## Prerequisites

- [Go](https://go.dev/dl/) (see `go.mod` for the version this repo targets)
- [Node.js](https://nodejs.org/) (for the frontend toolchain)
- [Wails v3 CLI](https://v3.wails.io/) — e.g. `go install github.com/wailsapp/wails/v3/cmd/wails3@latest`
- [ExifTool](https://exiftool.org/) on your `PATH` (or `TAKEOUT_EXIFTOOL_PATH` set to the binary)

## Run (development)

From the project root:

```bash
wails3 dev
```

That builds the frontend, generates bindings, and runs the app with hot reload. First run may install npm dependencies under `frontend/`.

## Build (production binary)

From the project root:

```bash
wails3 build
```

The output binary ends up under `bin/` (exact layout depends on your OS; the Taskfile uses `takeout-md-fixer` as the app name).



## Release builds (local)

From the repo root (macOS recommended for the DMG):

```bash
chmod +x scripts/build-release.sh
./scripts/build-release.sh
```

This updates `build/config.yml` to the chosen version (set `VERSION=x.y.z` to override), runs `wails3 package`, builds a **DMG** on macOS, and cross-compiles **Windows amd64** `takeout-md-fixer.exe` into `bin/`.

## Automated releases (GitHub Actions)

On every push to `main` (for example after merging a PR), [`.github/workflows/release.yml`](.github/workflows/release.yml) runs using official GitHub Actions (`actions/checkout`, `actions/setup-go`, `actions/setup-node`, `actions/upload-artifact`, `actions/download-artifact`) and the GitHub CLI (`gh release create`) to:

1. **Version**: next **patch** from the latest `v*.*.*` git tag (starts at `0.0.1` if no tags exist).
2. **Changelog**: commit subjects since that tag, plus a link to the build commit.
3. **Artifacts**: `takeout-md-fixer.dmg` (macOS) and `takeout-md-fixer.exe` (Windows amd64).
4. **Publish**: a **GitHub Release** with those files attached.

Workflows need **write** permission on **Actions** / **contents** (default `GITHUB_TOKEN` is enough for releases in the same repo).

## Open source

This project is open source and free to use. If you want to say thanks, you can [buy me a coffee](https://www.buymeacoffee.com/mrdevx).

## Author

**Mahdi Rashidi**

- Email: [contact@mrashidi.me](mailto:contact@mrashidi.me)
- Site: [mrashidi.me](https://mrashidi.me)
- GitHub: [@MRdevX](https://github.com/MRdevX)
