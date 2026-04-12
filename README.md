# Takeout Metadata Fixer

A small desktop app for **Google Takeout** folders. It reads the JSON files next to your photos and videos and writes dates, GPS, and related info into the media files using [ExifTool](https://exiftool.org/). After that, imports into iCloud, a NAS, or other apps usually show the right dates and places.

**You need ExifTool installed** on your computer. The app looks for it on your `PATH` and in usual install locations (for example Homebrew paths on macOS).

## Screenshots

Pick a Takeout folder and check that ExifTool is found.

![Welcome screen with folder selection and stepper](https://ecuc8c19st5krt52.public.blob.vercel-storage.com/step1-main.png)

See a short summary, then start. You can optionally delete the sidecar JSON after a successful run.

![Review step with statistics and Start](https://ecuc8c19st5krt52.public.blob.vercel-storage.com/step2-review.png)

Watch progress; you can pause or stop.

![Processing step with progress bar](https://ecuc8c19st5krt52.public.blob.vercel-storage.com/step3-progress.png)

When it finishes, you get counts for succeeded, skipped, and failed files.

![Done step with result summary](https://ecuc8c19st5krt52.public.blob.vercel-storage.com/step4-done.png)

## Download

**macOS (DMG)** and **Windows (.exe)** builds: [GitHub Releases](https://github.com/MRdevX/takeout-md-fixer/releases).

## Stopping and resuming

Progress is stored in **`.takeout-md-fixer-checkpoint.json`** inside the folder you chose. If you open that folder again, the app picks up where it left off. Delete that file if you want to process everything from the beginning.

## Develop

You need Go (see [`go.mod`](go.mod)), Node.js, the [Wails v3 CLI](https://v3.wails.io/) (`go install github.com/wailsapp/wails/v3/cmd/wails3@latest`), and ExifTool.

```bash
wails3 dev -config ./build/config.yml
```

The UI is built from `frontend/`; production assets go into `frontend/dist` (gitignored and embedded from [`main.go`](main.go)). Before `go build`, `go vet`, or `go test` at the repo root, run `npm ci` and `npm run build` in `frontend/`. From the repo root, **`task ci:local`** runs the full local check CI uses.

Most of the logic lives under **`internal/service`** (app flow and checkpoints), **`internal/takeout`** (Takeout JSON and file matching), and **`internal/exif`** (ExifTool). If you change Go APIs used by the UI, regenerate bindings with Wails (`wails3 generate bindings` — see Wails docs for flags).

## Build

```bash
wails3 build
```

Output is under `bin/`. For release-style builds (macOS DMG, Windows when supported):

```bash
bash scripts/build-release.sh
```

Optional: `VERSION=x.y.z`.

---

**Mahdi Rashidi** — [contact@mrashidi.me](mailto:contact@mrashidi.me) · [mrashidi.me](https://mrashidi.me) · [GitHub](https://github.com/MRdevX) · [Buy Me a Coffee](https://www.buymeacoffee.com/mrdevx)
