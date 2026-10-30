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

A section that reads something slow loads on its own: give `kit.Section` a `Src`, the address of the section's own HTML, and a `Full`, the page drawn whole on the server. The shell then comes at once with the section's skeleton; `kit.js` fetches every lazy section in parallel and puts each in its place, and a fetch that fails shows that section's error box with Retry. A GET form inside a loaded lazy section, such as its toolbar, reloads only that section and keeps its values in the address. Set `kit.Page.Full` too, and a browser without JavaScript goes straight to the whole page. Give a section an `Again`, the page with its sources read again (`?fresh=1`), and its head holds a Read again button: `kit.js` reloads only that section, from its `Src` with `fresh=1` added and its current sort or filter kept, so the server reads that section's sources again; without JavaScript it is a link to `Again`. While it reloads, the button spins and the section dims under a thin bar; then "Updated" shows beside the button for a few seconds and is announced (`aria-live`). Beside the button, `ui/read-at` says when the section's sources were last read ("read 10:42", the exact time on hover): a lazy section draws it empty and `kit.js` fills it from the `UI-Read-At` header (RFC 3339) of the section's HTML, so the server sends the time its sources were really read, not the time of the page. A section drawn whole takes a `kit.ReadAt` in `Read`, and the bar shows the page's in `kit.Bar.Read`. `/lazy` in the preview shows all of it.

When a toolbar or the cards/list toggle changes the address without a reload, links drawn with the old address would lose the change. Mark them, and `kit.js` keeps them current:

| Mark | Go field | Drawn as | Kept current |
|---|---|---|---|
| A link that carries the page's query | `Carry` on `kit.Link`, `kit.Tab`, `kit.Button`: its own query values, space-separated | `data-ui-carry="fresh"` | the link takes the page's new query; its own values (here `fresh`) stay as its `Href` had them |
| A lazy light's address of the page | `Back` on `kit.Light`: the name of `Src`'s value | `data-ui-back="back"` on the lazy box | that value of `data-ui-src` becomes the page's address |
| A form's address of the page | `Back` on `kit.Hidden` | `<input type="hidden" … data-ui-back>` | its value becomes the page's address |

The server draws every one of them right for the first load: a carried link's `Href` already holds the page's query (a Read again link `/?pr_sort=newest&fresh=1` with `Carry: "fresh"`; a tab `/product/x?pr_sort=newest&tab=board` with `Carry: "tab"`). A reload (F5) asks the server again and needs none of it.

A link inside a lazy section can change only that section's view, like a toolbar does: give it `data-ui-section-link` and an `href` to the page with the new query (an `<a>` in an inline SVG works too). `kit.js` makes that the address and reloads only the section from its `Src` with the new query; without JavaScript, or with a modifier key, it is a plain link to the page.

A status light (`kit.Light`, in `kit.Bar.Lights`) is a dot, an icon and a word in the top bar for one connection, such as GitHub. It is a link to the app's full page for it; `kit.js` makes it open a popover under it with the facts (`Details`) and an Actions menu. An action links, opens a dialog of `Dialogs` (a `kit.Dialog` whose `Body` holds its fields), or asks once more inside the popover (`Confirm`) before it posts. Give it a `Src` and the bar draws it grey at once and fetches it, like a lazy section.

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
| `js/kit.js` | theme switch, cards/list toggle, fold, tabs, lazy sections, carried links; no app logic |
| `react/` | the React wrappers and their types |
| `test/` | the tests of `kit.js`, in a small stand-in DOM |
| `kit.go`, `types.go` | package `kit` |
| `preview/`, `cmd/ui-kit-preview` | the gallery |
| `cmd/ui-kit-tokens` | writes `css/tokens.css` |

### Partials

Shell: `ui/shell`, `ui/assets`, `ui/preview-banner`, `ui/bar`, `ui/light`, `ui/rail`, `ui/page-header`, `ui/footer`.
Parts: `ui/product-card`, `ui/product-row`, `ui/products`, `ui/pr-item`, `ui/pr-list`, `ui/milestones`, `ui/status`, `ui/pill`, `ui/count`, `ui/section`, `ui/error`, `ui/notice`, `ui/empty`, `ui/table`, `ui/ticket`, `ui/button`, `ui/field`, `ui/toolbar`, `ui/view-toggle`, `ui/tabs`, `ui/crumb`, `ui/dialog`, `ui/drawer`, `ui/star`, `ui/skeleton`.
Pieces they share: `ui/icon`, `ui/num`, `ui/fact`, `ui/cell`, `ui/read-at`.

A count with an `Icon` (`eye`, `check`, `alert`, `edit`, or any `ui/icon`) is a badge: icon and number in its tone, and its `Label` as a tooltip on hover and on keyboard focus, read out through the same hidden label. With `Href` it is a link, such as to a group further down the page (an anchor inside the page; it opens no new tab).

## Rules

- **A link to another site or app opens in a new tab**, with `rel="noopener noreferrer"`; a link inside the app stays in the same tab. The partials do it with the `ext` function (`<a href="{{.Href}}"{{ext .Href}}>`); an app's own templates use the same function, and a form that sends the browser to another app (such as one that starts a board and redirects to it) carries `target="_blank"`. A test checks every golden for it.
- The server works out every value, tone, label, filter and order; a partial only draws it, and `kit.js` only saves a reload.
- Every part works without JavaScript.
- **A status light says one connection's state in its tone**: good green, warn amber, bad red, none grey; its word stays the same and the popover says why. Without JavaScript it is a link to the app's full page with the same facts and forms. Its actions are a menu button of plain buttons, not a `<select>`: choosing an option of a select must not act (it fires while a keyboard user moves through the options), and a menu of buttons needs no script to reach. An action that removes something asks once more inside the popover, never with `confirm()`. Escape, a click outside or Tab out closes the popover, and focus goes back to the light. A light that reads something slow has a `Src`, so it never slows the page.

## Change it

- A token: edit `tokens/tokens.json`, run `go run ./cmd/ui-kit-tokens`, commit both files. A test checks contrast in both themes.
- A part or a state: change its partial, class and wrapper together, add or change its case in `preview/cases.go`, then `go test -run Golden -update` and review the golden HTML.
- Checks: `go vet ./... && go test -race -count=1 ./...` and `npm test`.

A merge to `main` releases by itself from the commit types. A new token or part is `feat`; renaming or removing a class, partial, wrapper or field is `feat!`.
