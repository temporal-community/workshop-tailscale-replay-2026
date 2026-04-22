# Present the slides

The deck lives at `slides/slides.md` and runs under [Slidev](https://sli.dev). It ships with a Temporal-branded theme at `slides/theme-temporal/` (Inter typography, mint + purple palette, grid + planet backgrounds) and a custom Shiki theme for code blocks.

## Prerequisites

- Node.js 20+
- `pnpm` (the lockfile in `slides/pnpm-lock.yaml` pins versions for pnpm specifically; npm will work but you will lose reproducible installs)

```shell
# If you don't have pnpm yet:
corepack enable
corepack prepare pnpm@10.33.0 --activate
```

## Install

```shell
cd slides
pnpm install
```

This pulls in `@slidev/cli` and links the local `slidev-theme-temporal` theme from `slides/theme-temporal/`.

## Run the deck locally

```shell
pnpm dev
```

Slidev starts on `http://localhost:3030`. Useful endpoints:

| URL | What it is |
|---|---|
| `http://localhost:3030/` | The slideshow |
| `http://localhost:3030/presenter/` | Presenter mode (current + next slide, speaker notes, timer) |
| `http://localhost:3030/overview/` | All slides at a glance |

The dev server hot-reloads whenever you edit `slides.md` or anything under `theme-temporal/`.

## Presenting

Run `pnpm dev` on the presenter laptop. Open `http://localhost:3030/presenter/` on your laptop screen and `http://localhost:3030/` on the projector / shared screen. Slidev keeps them in sync. Speaker notes come from HTML comments in `slides.md` (look for the `<!-- ... -->` blocks after each slide).

Handy keyboard shortcuts while presenting:

| Key | Action |
|---|---|
| `space`, `right arrow`, `click` | next step (advances through v-click reveals) |
| `left arrow` | previous step |
| `f` | toggle fullscreen |
| `o` | open overview |
| `d` | toggle dark / light mode |

## Export to PDF

```shell
pnpm export
```

Produces `slides-export.pdf` in `slides/`. Needs Chromium; Slidev pulls in `playwright-chromium` automatically on first run.

For a PDF that shows every v-click reveal expanded (useful for handouts):

```shell
pnpm export -- --with-clicks
```

For PNGs per slide:

```shell
pnpm export -- --format png
```

## Build static HTML

```shell
pnpm build
```

Writes a static site to `slides/dist/`. Useful if you want to host the deck somewhere alongside the docs site, or serve it from a USB stick at a conference with sketchy wifi.

## Tweaking the theme

The theme is local, not installed from npm, so you can edit it in place and see changes instantly in the dev server:

| File | What to change |
|---|---|
| `slides/theme-temporal/styles/layout.css` | Typography, palette custom properties, heading sizes, footer style, background helpers |
| `slides/theme-temporal/styles/code.css` | Code-block padding, inline-code styling |
| `slides/theme-temporal/setup/temporal-dark.json` | Shiki token colors for fenced code blocks |
| `slides/theme-temporal/setup/mermaid.ts` | Mermaid palette (node colors, edge colors, actor styles) |
| `slides/theme-temporal/layouts/*.vue` | Per-layout structure and overrides |
| `slides/theme-temporal/assets/*.webp` | Backgrounds (grid, teal planet, purple planet, rex constellation, glow) |

If you replace a background image, keep it WebP at ~1920 wide and ~200-300 KB for good projector quality without a slow page load. The originals are 16:9.
