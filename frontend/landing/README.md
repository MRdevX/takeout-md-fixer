# Takeout Metadata Fixer (landing)

Static marketing page for [takeout.mrashidi.me](https://takeout.mrashidi.me), built with [Vue 3](https://vuejs.org/) and [Vite](https://vite.dev/). Styling matches the desktop app (see `frontend/public/style.css` tokens).

## Setup

```sh
pnpm install
```

## Develop

```sh
pnpm dev
```

## Build

```sh
pnpm build
```

Output: `landing/dist/`. Upload that folder to your host as static files.

## Deploy (takeout.mrashidi.me)

1. Build with `pnpm build`.
2. Point DNS for `takeout.mrashidi.me` to your hosting provider.
3. Serve `dist/` as the site root (no server-side routing needed; single-page).

Examples: nginx `root` to `dist`, Cloudflare Pages, GitHub Pages with a custom domain, or any static host.

## Links

- Downloads: [GitHub Releases](https://github.com/MRdevX/takeout-md-fixer/releases)
- Source: [github.com/MRdevX/takeout-md-fixer](https://github.com/MRdevX/takeout-md-fixer)
