# Takeout Metadata Fixer

Desktop app for **Google Takeout** exports: it finds JSON sidecars next to your photos and videos and writes dates, GPS, and related fields into the files with [ExifTool](https://exiftool.org/). That way imports into iCloud, a NAS, or other tools see sensible dates and locations.

You must install [ExifTool](https://exiftool.org/) yourself. The app resolves it from your `PATH` and common install locations (e.g. Homebrew and `/usr/local/bin` on macOS), which helps when the GUI does not see the same `PATH` as your terminal.

## Screenshots

**Welcome** — pick a Takeout folder and confirm ExifTool is available.

![Welcome screen with folder selection and stepper](https://ecuc8c19st5krt52.public.blob.vercel-storage.com/step1-main.png)

**Review** — folder summary and counts; optional deletion of sidecar JSON after a successful run.

![Review step with statistics and Start](https://ecuc8c19st5krt52.public.blob.vercel-storage.com/step2-review.png)

**Update** — progress while metadata is written; pause or stop anytime.

![Processing step with progress bar](https://ecuc8c19st5krt52.public.blob.vercel-storage.com/step3-progress.png)

**Done** — per-run succeeded, skipped, and failed totals.

![Done step with result summary](https://ecuc8c19st5krt52.public.blob.vercel-storage.com/step4-done.png)

## Download

Prebuilt **macOS (DMG)** and **Windows (.exe)** builds are on [GitHub Releases](https://github.com/MRdevX/takeout-md-fixer/releases). Pushes to `main` trigger a release workflow that bumps the patch version from git tags and publishes artifacts plus a short changelog.

## Resuming

If you stop early or quit, progress is saved in a hidden file **`.takeout-md-fixer-checkpoint.json`** in the folder you selected. The next time you open that folder in the app, it **resumes automatically** and skips files already processed. Delete that file if you want a full pass from scratch on the same folder.

## Develop

You need Go (see [`go.mod`](go.mod)), Node.js, the [Wails v3 CLI](https://v3.wails.io/) (`go install github.com/wailsapp/wails/v3/cmd/wails3@latest`), and ExifTool.

```bash
wails3 dev -config ./build/config.yml
```

[`main.go`](main.go) embeds [`frontend/dist`](frontend/dist), which is gitignored. Before `go build`, `go vet`, or `go test` at the repo root, build the frontend (`npm ci` and `npm run build` in `frontend/`). CI does the same. From the repo root you can run **`task ci:local`** to build the frontend (via Wails bindings + Vite), then `gofmt`, `go vet`, and `go test`.

### Go packages (where to change what)

| Area                                             | Package                                | Notes                                                                                                                                       |
| ------------------------------------------------ | -------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------- |
| Wails API, fix job, pause/resume, checkpoints    | [`internal/service`](internal/service) | Emits app events `fix-progress` and `fix-complete`; checkpoint file `.takeout-md-fixer-checkpoint.json` lives in the chosen Takeout folder. |
| Scanning, sidecar resolution, Takeout JSON types | [`internal/takeout`](internal/takeout) | New filename patterns or JSON fields usually start here.                                                                                    |
| ExifTool invocation and tag mapping              | [`internal/exif`](internal/exif)       | New EXIF / video date tags.                                                                                                                 |
| Path comparison keys (scan vs checkpoint)        | [`internal/pathkey`](internal/pathkey) | Single normalization used across scan and checkpoint maps.                                                                                  |

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
