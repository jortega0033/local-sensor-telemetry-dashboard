# Design tokens

The page uses the Pipenzo "Flow look" tokens: a black page, Flow-style alpha layers and one accent, Pipenzo yellow `#FECE02`. Red exists only as the danger colour.

## Source

`web/tokens.css` is a byte-for-byte copy of `extension/sidepanel/tokens.css` in `jortega0033/clip-to-reel-studio`, commit `a13378ce73e5b27213c9a200fc960ec93afef054` (2026-11-21). Do not edit it here.

## Resync

Copy the file over `web/tokens.css`. Nothing else changes.

## Loading

The page embeds `web/tokens.css` with the other assets (same-origin) and links it with `<link rel="stylesheet" href="tokens.css">` before any other stylesheet. The first UI ticket (#13) wires this and checks it.

## Rule

New UI uses these tokens. No new colour literals in CSS or JavaScript. Script that draws or sets inline colours reads them with `getComputedStyle(document.documentElement).getPropertyValue('--token')`. The header comment of the file holds the contrast rules: yellow only on the primary action and the current-step marker, and text never on a bare alpha layer.

## Light and dark

Dark is the primary design. The `:root` block is the matching light theme, and the dark block inside the `prefers-color-scheme: dark` media query follows the operating system colour scheme. There is no manual switch. `color-scheme: light dark` is set.

## Token names by role

- Surfaces: `--bg`, `--bg-raised`, `--bg-sunken`, `--surface`, `--ink`
- Text: `--fg`, `--muted`, `--faint`, `--on-primary`, `--on-ink`, `--on-ink-muted`, `--on-media`
- Lines: `--border`, `--hairline`, `--control-border`, `--input-border`, `--ink-border`
- Accent and actions: `--primary`, `--primary-edge`, `--primary-press`, `--ink-hover`, `--btn-secondary`, `--btn-secondary-border`, `--focus`, `--ring-accent`
- State layers: `--state-hover`, `--state-selected`
- Status: `--danger`, `--danger-bg`, `--ok`, `--ok-bg`, `--success`, `--info`, `--warn-fg`, `--warn-bg`
- Hero: `--hero-bg`, `--hero-border`, `--hero-chip-border`, `--hero-shadow`, `--on-hero`, `--on-hero-muted`
- Glass and gradients (not used by this dashboard, see `docs/ui-quality.md` rules 3 and 4): `--glass*`, `--grad-*`, `--scrim`
- Type: `--font`, `--font-display`, `--mono`. This repo bundles no font files (no `assets/fonts`, no `fonts.css`), so the stack falls back to the system fonts (Segoe UI / system-ui), which is intended. Any font added later must be a same-origin file embedded in `web/`, never a CDN or a hosted font.
- Shape and space: `--radius-*`, `--space-1`, `--shadow-*`
- Motion: `--dur-1`, `--ease-in-out`, `--ease-out`

Check the file itself for the exact list; some tokens exist only for the extension and may not fit this dashboard.
