# Mission

ui-kit is the org's shared UI kit. Every web app of the org draws its frame and its parts from it, so the same thing looks and reads the same everywhere: a status, a count, a card, a list, an error. Change a part once in the kit and every app follows at its next upgrade.

It holds the tokens (colour and size names), the shell (top bar, page header, sections, footer) and the parts, as CSS that every app loads, HTML partials for Go apps, thin React wrappers for Vite apps, and one small script.

## Principles

1. **Written once.** Every value is a token in `tokens/tokens.json`. Every part is one set of `ui-` classes in `css/kit.css`. The Go partials and the React wrappers draw the same markup.
2. **Draws, never decides.** The kit shows what the app's server sends: order, tone, label, whether a thing is in focus (G-STACK-7). It holds no app logic and no data.
3. **Works without JavaScript.** Every part draws on the server. `kit.js` only saves a reload.
4. **Fails in its place.** A part that cannot show its data says so inside itself; the rest of the page stays.
5. **Readable by everyone.** Both themes, text at 4.5:1 and shapes at 3:1, a visible keyboard ring, calm motion, phone width.

## Out of scope

- App logic, data fetching, sorting or filtering.
- An app's own canvas (e-designer keeps its canvas colours where the meaning differs).
- Web components; revisit if a third kind of app appears.

## Words

| Word | Meaning |
|---|---|
| Token | A named value, like `--ui-good`, used instead of writing a colour or size by hand. |
| Shell | The frame of every page: top bar, optional left rail, page header, sections, footer. |
| Part | One reusable piece of a page, like the product card or the PR item. The kit has 20. |
| Partial | The Go `html/template` that draws one part, named `ui/<part>`. |
| Wrapper | The React component that draws one part with the same markup as its partial. |
| Tone | The meaning a colour carries: good, warn, bad, info, accent, or none. |
| Focus | What the org works on now, marked with the yellow star. Not the keyboard ring. |
| Case | One state of one part with sample data. The preview shows every case; the golden tests check every case. |
| Preview | The local gallery of every case in both themes, which the owner looks at before a release. |
