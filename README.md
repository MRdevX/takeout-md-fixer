# Takeout Metadata Fixer

**Google Takeout is great. The metadata on your files afterward? Not always.** This app reads the JSON sidecars next to your photos and videos, then writes dates, GPS, and related fields back into the files with [ExifTool](https://exiftool.org/), so imports to iCloud, a NAS, or elsewhere look right.

[![Go](https://img.shields.io/badge/Go-00ADD8?style=flat-square&logo=go&logoColor=white)](https://go.dev/) [![Wails](https://img.shields.io/badge/Wails_v3-27272a?style=flat-square&logo=wails&logoColor=white)](https://v3.wails.io/) [![Vite](https://img.shields.io/badge/Vite-646CFF?style=flat-square&logo=vite&logoColor=white)](https://vitejs.dev/) [![JavaScript](https://img.shields.io/badge/JavaScript-F7DF1E?style=flat-square&logo=javascript&logoColor=black)](https://developer.mozilla.org/docs/Web/JavaScript) [![ExifTool](https://img.shields.io/badge/ExifTool-system-8B4513?style=flat-square)](https://exiftool.org/)

**You need [ExifTool](https://exiftool.org/) on your `PATH`.** The app checks on startup. Non-standard install? Set `TAKEOUT_EXIFTOOL_PATH` to `exiftool` or `exiftool.exe` (Windows).

## Why I built this

Leaving Google Photos, Takeout files didn’t match what I’d seen in the app. I wrote this to fix metadata before moving the rest of my library.

Stack: Go, [Wails v3](https://v3.wails.io/), [Vite](https://vitejs.dev/) + vanilla JS, [go-exiftool](https://github.com/barasher/go-exiftool). ExifTool is installed by you.

## Before you start

- [Go](https://go.dev/dl/) ([`go.mod`](go.mod))
- [Node.js](https://nodejs.org/)
- [Wails v3 CLI](https://v3.wails.io/) (`go install github.com/wailsapp/wails/v3/cmd/wails3@latest`)
- [ExifTool](https://exiftool.org/) on `PATH` or `TAKEOUT_EXIFTOOL_PATH`

## Develop

```bash
wails3 dev
```

## Build

```bash
wails3 build
```

Binary under `bin/`.

## Local release (macOS DMG + Windows exe)

```bash
chmod +x scripts/build-release.sh
./scripts/build-release.sh
```

Optional: `VERSION=x.y.z`. See `build/config.yml` and `bin/` for outputs.

## Releases

Pushes to `main` run [`.github/workflows/release.yml`](.github/workflows/release.yml): version bump from git tags, changelog, DMG + `takeout-md-fixer.exe`, GitHub Release.

If this helped, fuel is welcome.

<a href="https://www.buymeacoffee.com/mrdevx" title="Buy Me A Coffee"><img src="https://cdn.buymeacoffee.com/buttons/v2/default-yellow.png" alt="Buy me a coffee on buymeacoffee.com" width="217" height="60" /></a>

## Author

**Mahdi Rashidi**

- [contact@mrashidi.me](mailto:contact@mrashidi.me)
- [mrashidi.me](https://mrashidi.me)
- [@MRdevX](https://github.com/MRdevX)
