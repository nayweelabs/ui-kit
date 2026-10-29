READ THE ORG CONSTITUTION FIRST. Find it by asking hubman (`hubman path constitution`); without hubman it is at ~/.local/share/constitution, and if it is missing there, clone git@github.com:ElderlyGate/constitution.git into it. Then read this project's constitution/. A rule here holds unless the project overrides it by ID.

WORK FROM THE CARDS ON THE BOARD. ARCHIVED CARDS ARE IGNORED.

## This project

ui-kit (product ID `ui-kit`) is the org's shared UI kit, a library: tokens, the shell and the parts every web app takes from one place. The plan is the e-designer design in `.designer/ui-kit/` (G-PLAN-9); its releases are the MVPs.

- Change a token only in `tokens/tokens.json`, then run `go run ./cmd/ui-kit-tokens` and commit `css/tokens.css` with it.
- A new part or state is a case in `preview/cases.go`; `go test -run Golden -update` writes its golden HTML, which is reviewed like code.
- Look at it with `go run ./cmd/ui-kit-preview` (http://127.0.0.1:7790).
- Checks: `GOWORK=off go vet ./... && GOWORK=off go test -race -count=1 ./...` and `npm test`.
