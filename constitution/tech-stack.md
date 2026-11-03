# Tech stack

Only what is this project's own. The org rules hold as written, except the overrides at the end.

## Stack

- **CSS** for every part: `css/tokens.css` (generated) and `css/kit.css`. Classes start with `ui-`; values are `--ui-*` tokens only.
- **Go** library, module `github.com/ElderlyGate/ui-kit`, package `kit` at the module root, so it can embed `css/`, `fonts/`, `js/` and `partials/`. Partials are `html/template`.
- **JavaScript**: one file, `js/kit.js`, plain, no build. It only switches theme and view, folds, moves between tabs, and runs the small helpers of the field, id and dialog parts.
- **React** wrappers in `react/index.js`: plain ES modules with `createElement`, types in `react/index.d.ts`. No build step for anyone.
- **Fonts**: Hanken Grotesk and IBM Plex Mono, latin subsets as woff2, bundled under the SIL Open Font License with their licence files.

## Architecture

- `tokens/tokens.json` is the one source of every token. `cmd/ui-kit-tokens` writes `css/tokens.css` from it; package `tokens` gives Go programs the same values.
- `[data-theme]` works on `:root` and on any element, so light and dark can show side by side.
- Each part is one partial, one set of classes and one wrapper, with the same markup.
- `preview/cases.go` lists every part in every state. The preview draws the cases; the golden tests in `testdata/golden/` check them.
- Go apps mount `kit.Assets()` at `/ui-kit/` and parse their templates into `kit.Templates()`. Vite apps install the package from a release tag and import `css/kit.css` and `react`.

## Decisions

- A change to a token or a new part is `feat`. Renaming or removing a class, a partial, a wrapper or a field is `feat!`.
- No version file: `package.json` has no `version`, so the tag is the only version (V-1, V-6).
- Golden HTML changes are reviewed like code.

## Overrides

- **H-1** 2026-10-29 · ui-kit is a library with no `ui-kit` command, so it answers no `help`. Its two developer commands, `go run ./cmd/ui-kit-tokens` and `go run ./cmd/ui-kit-preview`, are listed in the README.
