# ui-kit
The org's shared UI kit: shell, tokens and components for every web app

One set of colour and size names, one page frame, and 20 parts (product card, PR item, milestone strip, section with its own error box, and the small parts they are made of). Go apps take HTML partials, Vite apps take thin React wrappers, and both load the same CSS.

## Look at it

```sh
go run ./cmd/ui-kit-preview            # http://127.0.0.1:7790
go run ./cmd/ui-kit-preview -addr 127.0.0.1:7791
```

The gallery shows every part in every state, light on the left and dark on the right (`/?themes=one` for one theme, `/?part=pr-item` for one part). Three example pages are built only from the kit: `/focus` (the control room), `/products` (cards or list, filtered on the server) and `/admin` (the left rail, tabs and live dialogs).

## In a Go app

```go
import kit "github.com/ElderlyGate/ui-kit"

mux.Handle(kit.Prefix, kit.Assets())                         // /ui-kit/css, /ui-kit/fonts, /ui-kit/js
tpl := template.Must(kit.Templates().ParseFS(views, "views/*.html"))
```

```html
{{template "ui/shell" .Page}}                                 <!-- a whole page -->
{{template "ui/pr-item" .PR}}                                 <!-- one part -->
{{template "ui/pill" (pill .CI.Tone .CI.Label)}}
```

The data each partial takes is a Go type in `types.go` (`kit.Page`, `kit.Section`, `kit.Product`, `kit.PR`, …). The server works out every value, tone and label; the partial only draws it. `kit.HTML` draws a partial into `template.HTML`, ready to be a section's body.

A section that reads something slow loads on its own: give `kit.Section` a `Src`, the address of the section's own HTML, and a `Full`, the page drawn whole on the server. The shell then comes at once with the section's skeleton; `kit.js` fetches every lazy section in parallel and puts each in its place, and a fetch that fails shows that section's error box with Retry. A GET form inside a loaded lazy section, such as its toolbar, reloads only that section and keeps its values in the address. Set `kit.Page.Full` too, and a browser without JavaScript goes straight to the whole page. `/lazy` in the preview shows all of it.

A page that does not use `ui/shell` links `/ui-kit/css/kit.css` and `/ui-kit/js/kit.js` (in `<head>`, without `defer`) and puts `class="ui-page"` on its body.

## In a Vite app

```json
"dependencies": { "@elderlygate/ui-kit": "github:ElderlyGate/ui-kit#v0.1.0" }
```

```js
import "@elderlygate/ui-kit/css/kit.css";
import { Shell, Section, PrList, Products, Pill } from "@elderlygate/ui-kit/react";
```

Put `class="ui-page"` on `<body>`. The wrappers are plain ES modules with types; there is nothing to build. A `Section` catches an error in its children and shows its own error box. Without React, use the classes directly and add `js/kit.js`.

## What is in it

| Path | Holds |
|---|---|
| `tokens/tokens.json` | every token, written once |
| `css/tokens.css` | the `--ui-*` properties, generated from `tokens.json` |
| `css/kit.css` | the shell and every part; imports `tokens.css` and the fonts |
| `fonts/` | Hanken Grotesk and IBM Plex Mono (SIL OFL 1.1, licences beside them) |
| `partials/` | the Go partials, `ui/…` |
| `js/kit.js` | theme switch, cards/list toggle, fold, tabs, lazy sections; no app logic |
| `react/` | the React wrappers and their types |
| `kit.go`, `types.go` | package `kit` |
| `preview/`, `cmd/ui-kit-preview` | the gallery |
| `cmd/ui-kit-tokens` | writes `css/tokens.css` |

### Partials

Shell: `ui/shell`, `ui/assets`, `ui/bar`, `ui/rail`, `ui/page-header`, `ui/footer`.
Parts: `ui/product-card`, `ui/product-row`, `ui/products`, `ui/pr-item`, `ui/pr-list`, `ui/milestones`, `ui/status`, `ui/pill`, `ui/count`, `ui/section`, `ui/error`, `ui/notice`, `ui/empty`, `ui/table`, `ui/ticket`, `ui/button`, `ui/field`, `ui/toolbar`, `ui/view-toggle`, `ui/tabs`, `ui/crumb`, `ui/dialog`, `ui/drawer`, `ui/star`, `ui/skeleton`.
Pieces they share: `ui/icon`, `ui/num`, `ui/fact`, `ui/cell`.

## Rules

- **A link to another site or app opens in a new tab**, with `rel="noopener noreferrer"`; a link inside the app stays in the same tab. The partials do it with the `ext` function (`<a href="{{.Href}}"{{ext .Href}}>`); an app's own templates use the same function, and a form that sends the browser to another app (such as one that starts a board and redirects to it) carries `target="_blank"`. A test checks every golden for it.
- The server works out every value, tone, label, filter and order; a partial only draws it, and `kit.js` only saves a reload.
- Every part works without JavaScript.

## Change it

- A token: edit `tokens/tokens.json`, run `go run ./cmd/ui-kit-tokens`, commit both files. A test checks contrast in both themes.
- A part or a state: change its partial, class and wrapper together, add or change its case in `preview/cases.go`, then `go test -run Golden -update` and review the golden HTML.
- Checks: `go vet ./... && go test -race -count=1 ./...` and `npm test`.

A merge to `main` releases by itself from the commit types. A new token or part is `feat`; renaming or removing a class, partial, wrapper or field is `feat!`.
