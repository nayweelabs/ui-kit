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

A page's own actions (`kit.BarAction`, in `kit.Bar.Actions`) sit in the top bar left of the theme switch, such as Open design. Each is a form that posts to `Action`, in `Target` when set (`_blank` for a new tab), with the action's `Button` as its submit button, so in Go it asks first as any kit button does (`Confirm`), and an optional `Count` beside it. React: `Bar`'s `actions`; React's `Button` has no confirm today.

### A change asks first

Every change asks before it is made (the org rule), and the user picks how, per action type. The types are the same words in every product (`kit.ActionTypes`): create, edit, move, focus, archive, close, delete, settings, connect, run. Each type asks in one of three ways: `name` (type the name to go ahead), `yes-no`, or `never`. By default every type asks yes-no, delete included. A change to the system itself (uninstalling, removing a token, dropping data, deleting a board) asks for the name; the product names those changes, and the user can still pick another way.

The choices live in this app's own config only; a type it does not set takes the change's default. The server reads them and draws each button the way its type asks, passing the default: `kit.AskYesNo`, or `kit.AskName` for a change to the system itself.

```go
style, source := kit.ResolveConfirm(appCfg.Confirm["delete"], kit.AskYesNo) // source: "this app" or "default"
btn := kit.Button{Label: "Delete card", Variant: "danger", Type: "submit"}.Ask(style, "Delete card 20261102-e472?", "20261102-e472")
```

`Ask` sets the button's `Confirm` (the question), `ConfirmName` (only for `name`) and `Remember`; for `never` it leaves them empty and the form is sent at once. `kit.js` asks the question in the kit's dialog before the form is sent: the question, Cancel (it gets first focus), and a button that repeats the verb. With a name, the dialog says "Type 20261102-e472 to go ahead" and the button works only once the text is exactly the name; the text box gets first focus. Escape or Cancel sends nothing; the confirm button sends the form with the pressed button's name and value and `confirm=yes`. Every question offers "Remember my answer"; ticked, the form sends `confirm=remember` instead, and the server sets that type to `never` in this app's config.

`kit.js` stops a marked form in the capture phase, before any handler of the app sees it, so no app code can send it first; only the answered form goes on. A browser without `requestSubmit` still asks, then sends the form itself with the pressed button's value; one without `<dialog>` asks with its own question.

**The server checks the yes.** A change that asks (its style is not `never`) goes ahead only when `kit.Confirmed(r)` says so; otherwise it answers `kit.RefuseUnconfirmed(w)`: `428 Precondition Required` and "Nothing was changed", and makes no change. That is the web twin of "refused without `--yes`": a form sent without JavaScript, or by a script, never acts unasked.

```go
ok, remember := kit.Confirmed(r)
if !ok && style != kit.AskNever {
	kit.RefuseUnconfirmed(w)
	return
}
if remember {
	appCfg.Confirm["delete"] = kit.AskNever
}
```

| Mark | On | Means |
|---|---|---|
| `data-ui-confirm="Save the settings?"` | a form or its submit button | ask this question before the form is sent |
| `data-ui-confirm-name="20261102-e472"` | the same element | the user types this name before the confirm button works |
| `data-ui-confirm-label="Save"` | the same element | the confirm button's words; default, the pressed button's text |
| `data-ui-remember` | the same element | offer "Remember my answer", which sends `confirm=remember` instead of `confirm=yes` |

The Settings page draws `ui/confirm-settings` from a `kit.ConfirmSettings`: `Action`, the product's `Types`, this app's values in `App`, and in `System` the types that have a verb changing the system itself. Each row shows the style in use and whether this app sets it or it is the default (a `this app` or `default` tag); a `System` type says the real default, "Ask y/N · name for changes to the system" (also beside Don't ask, since a remembered `never` never silences a change to the system), and four choices: Type the name, Ask y/N, Don't ask, and Use the default, which drops this app's own value.

The form posts `confirm.<type>`, a style or `default`. `kit.ConfirmChoices(r.PostForm, types)`, with the same `Types` the form was drawn with, returns a choice for exactly those types: keep each style in this app's config, and drop each type that says `default`.

In React: `<ConfirmForm id question label confirmName remember onSubmit>` around a form's fields, `resolveConfirm(app, def)` (this app's value for the type and the change's default, `yes-no` when left out; gives `{style, source}`), `actionTypes` and `<ConfirmSettings action types app system onSubmit>`. Without JavaScript the form is sent at once, without asking. A form that posts to a page listing the exact changes needs no dialog: that page is the question. `/admin` in the preview has Restart (asks yes-no) and Delete (a change to the system itself, asks for the name), each asking until you tick "Remember my answer"; the gallery shows the Settings rows.

### Where things are kept

Every product's Settings page says where it keeps its files, with `ui/locations`. The server passes the real paths, resolved through mankit's `home`; the kit shows the home folder as `~`, and each path is a copy button that copies the full path.

```go
kit.Locations{Items: []kit.Location{
	{Label: "This app’s config", Path: configPath}, // ~/.config/elderlygate/<product>/config.yaml
	{Label: "This app’s data", Path: dataDir},     // ~/.local/share/elderlygate/<product>/
}}
```

`Home` sets the folder shown as `~`; left empty, it is the home of the user the server runs as. In React, `<Locations home items>` draws the same; pass `home`, since a browser does not know it.

### License and features

Every product's Settings page shows its license and its features with `ui/license`, from a `kit.License`: the data of `<app> license show` as it is (its JSON decodes straight into it); what each field means is said once, in the contracts schema (mankit/v1 license-show). The kit draws it and decides nothing: the states, the features in their order, the counts, a countdown of each license's `daysLeft` against its `window.of`, and every refused file with its reason. A build that is not a release draws the dev look beside its state. The headline is the answer's `sentence` word for word; without one it names the state, and an ended license's day is the answer's `endedOn`. The card and the countdown take the warn colour when the answer's `warn` says so. A release-trial's row shows its `trialEnds`, else its `offerEnds`, in place of `until`. The page adds the build's `Version`, its own buttons in `Actions`, and in `Trial` where to ask for a trial key, which the app sets only for none and trial-ended. `ui/license-features` is the features table alone. In React, `<License>` and `<LicenseFeatures>` take the same answer; a test draws the Go cases (`testdata/golden/license-cases.json`) in React and checks each against its golden.

`ui/license-badge` is the one license badge, on the Settings card and in the top bar: Free (none, trial-ended, ended), Trial (trial) or Licensed (licensed, ending), and nothing else; only its colour changes, to warn when the answer's `warn` is true. Every app's top bar shows it (ADR-L3): set `kit.Bar.License` to a `kit.LicenseLight` with the same `kit.License` and `Href`, the app's Settings page. Without JavaScript the badge links there; `kit.js` makes it open the License panel (`ui/license-panel`) under it: the badge, the answer's `sentence` as given, the app's `Actions`, which hold only its add license button, and the features table. The panel builds no text. In React, `<Bar license={{ href, license }}>`, `<LicenseLight>`, `<LicenseBadge>` and `<LicensePanel>` take the same; a test draws every contracts sample answer of license show (`testdata/golden/license-samples.json`, from the pinned contracts module) in Go and in React.

A page that does not use `ui/shell` links `/ui-kit/css/kit.css` and `/ui-kit/js/kit.js` (in `<head>`, without `defer`) and puts `class="ui-page"` on its body.

## In a Vite app

```json
"dependencies": { "@elderlygate/ui-kit": "github:nayweelabs/ui-kit#v0.1.0" }
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
| `js/kit.js` | theme switch, cards/list toggle, fold, tabs, lazy sections, carried links, asking before a form is sent; no app logic |
| `react/` | the React wrappers and their types |
| `test/` | the tests of `kit.js`, in a small stand-in DOM |
| `kit.go`, `types.go`, `confirm.go` | package `kit`; `confirm.go` holds how a change asks |
| `preview/`, `cmd/ui-kit-preview` | the gallery |
| `cmd/ui-kit-tokens` | writes `css/tokens.css` |

### Partials

Shell: `ui/shell`, `ui/assets`, `ui/preview-banner`, `ui/bar`, `ui/bar-action`, `ui/light`, `ui/license-light`, `ui/rail`, `ui/page-header`, `ui/footer`.
Parts: `ui/product-card`, `ui/product-row`, `ui/products`, `ui/pr-item`, `ui/pr-list`, `ui/milestones`, `ui/status`, `ui/pill`, `ui/count`, `ui/section`, `ui/error`, `ui/notice`, `ui/empty`, `ui/table`, `ui/ticket`, `ui/button`, `ui/field`, `ui/toolbar`, `ui/view-toggle`, `ui/tabs`, `ui/crumb`, `ui/dialog`, `ui/confirm-settings`, `ui/locations`, `ui/license`, `ui/license-features`, `ui/license-badge`, `ui/license-panel`, `ui/drawer`, `ui/star`, `ui/skeleton`.
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
