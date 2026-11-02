package kit_test

import (
	"bytes"
	"flag"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	kit "github.com/ElderlyGate/ui-kit"
	"github.com/ElderlyGate/ui-kit/preview"
)

var update = flag.Bool("update", false, "rewrite the golden files in testdata/golden")

// Every part in every state draws exactly its golden HTML.
func TestPartialsMatchGoldenHTML(t *testing.T) {
	seen := map[string]bool{}
	for _, c := range preview.Cases() {
		name := c.Part + "--" + c.State
		if seen[name] {
			t.Fatalf("case %s is listed twice", name)
		}
		seen[name] = true
		t.Run(name, func(t *testing.T) {
			var b bytes.Buffer
			for _, it := range c.Items {
				if err := kit.Render(&b, it.Template, it.Data); err != nil {
					t.Fatal(err)
				}
				b.WriteByte('\n')
			}
			path := filepath.Join("testdata", "golden", name+".html")
			if *update {
				if err := os.WriteFile(path, b.Bytes(), 0o644); err != nil {
					t.Fatal(err)
				}
				return
			}
			want, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("%v (run go test -run Golden -update to write it)", err)
			}
			if !bytes.Equal(b.Bytes(), want) {
				t.Errorf("%s differs from its golden file; got:\n%s", name, b.String())
			}
		})
	}
}

// Every part the gallery lists has at least one case, and every case belongs to a listed part.
func TestEveryPartHasCases(t *testing.T) {
	parts := map[string]int{}
	for _, p := range preview.Parts() {
		parts[p.ID] = 0
	}
	for _, c := range preview.Cases() {
		if _, ok := parts[c.Part]; !ok {
			t.Errorf("case %s--%s names no listed part", c.Part, c.State)
		}
		parts[c.Part]++
	}
	for id, n := range parts {
		if n == 0 {
			t.Errorf("part %s has no case", id)
		}
	}
}

func TestShellDrawsAWholePage(t *testing.T) {
	var b bytes.Buffer
	err := kit.Render(&b, "ui/shell", kit.Page{
		Title:  "Control room",
		Bar:    kit.Bar{App: "Control room", Pages: []kit.Link{{Label: "Focus", Href: "/", Current: true}}},
		Header: kit.Head{Title: "Focus"},
		Body:   "<p>body</p>",
		Footer: kit.Footer{Version: "v0.4.0"},
		Assets: "/static/kit",
	})
	if err != nil {
		t.Fatal(err)
	}
	page := b.String()
	for _, want := range []string{
		"<!doctype html>", `<title>Focus · Control room</title>`, `href="/static/kit/css/kit.css"`,
		`src="/static/kit/js/kit.js"`, `<body class="ui-page">`, `aria-current="page">Focus</a>`, "<p>body</p>", `class="ui-footer`,
	} {
		if !strings.Contains(page, want) {
			t.Errorf("page lacks %q", want)
		}
	}
}

func TestShellDrawsThePreviewBannerOnlyForABuildThatIsNotARelease(t *testing.T) {
	page := func(b *kit.Build) string {
		var w bytes.Buffer
		if err := kit.Render(&w, "ui/shell", kit.Page{Bar: kit.Bar{App: "x"}, Build: b}); err != nil {
			t.Fatal(err)
		}
		return w.String()
	}
	preview := page(&kit.Build{Version: "v0.4.0-3-g758ac48", Branch: "feat/road-focus", Commit: "758ac48", Dirty: true})
	for _, want := range []string{
		`<body class="ui-page">` + "\n" + `<div class="ui-preview-banner" role="note" aria-label="Preview build" title="PREVIEW · branch feat/road-focus · 758ac48 · not a release · uncommitted changes · version v0.4.0-3-g758ac48">`,
		`<b>PREVIEW</b> · branch feat/road-focus · 758ac48 · not a release · uncommitted changes</span></div>` + "\n" + `<header class="ui-bar`,
	} {
		if !strings.Contains(preview, want) {
			t.Errorf("preview page lacks %q:\n%s", want, preview)
		}
	}
	for name, b := range map[string]*kit.Build{"no build": nil, "released": {Version: "v0.4.0", Branch: "main", Commit: "9a1b2c3", Released: true}} {
		if got := page(b); strings.Contains(got, "ui-preview-banner") || !strings.Contains(got, `<body class="ui-page">`+"\n"+`<header class="ui-bar`) {
			t.Errorf("%s: page draws the preview banner or a gap where it would be:\n%s", name, got)
		}
	}
}

func TestAppTemplatesCanUseThePartials(t *testing.T) {
	tpl, err := kit.Templates().Parse(`{{define "app"}}{{template "ui/pill" (pill .Tone .Label)}}{{end}}`)
	if err != nil {
		t.Fatal(err)
	}
	var b bytes.Buffer
	if err := tpl.ExecuteTemplate(&b, "app", map[string]string{"Tone": "warn", "Label": "3 behind"}); err != nil {
		t.Fatal(err)
	}
	if got, want := b.String(), `<span class="ui-pill" data-tone="warn">3 behind</span>`; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if _, err := kit.Templates().Parse(`{{define "x"}}{{end}}`); err != nil {
		t.Fatalf("a second app cannot use the partials after the first: %v", err)
	}
}

func TestDataIsEscaped(t *testing.T) {
	var b bytes.Buffer
	if err := kit.Render(&b, "ui/pill", kit.Pill{Label: `<script>alert(1)</script>`, Href: "javascript:alert(1)"}); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(b.String(), "<script>") || strings.Contains(b.String(), "javascript:") {
		t.Fatalf("not escaped: %s", b.String())
	}
}

func TestAssetsServeOnlyTheKitFiles(t *testing.T) {
	srv := httptest.NewServer(http.StripPrefix("", kit.Assets()))
	defer srv.Close()
	for path, want := range map[string]int{
		"/ui-kit/css/kit.css":    200,
		"/ui-kit/css/tokens.css": 200,
		"/ui-kit/js/kit.js":      200,
		"/ui-kit/fonts/hanken-grotesk-latin-wght-normal.woff2": 200,
		"/ui-kit/fonts/OFL-HankenGrotesk.txt":                  200,
		"/ui-kit/partials/basics.html":                         404,
		"/ui-kit/css/":                                         404,
		"/ui-kit/../kit.go":                                    404,
	} {
		res, err := http.Get(srv.URL + path)
		if err != nil {
			t.Fatal(err)
		}
		res.Body.Close()
		if res.StatusCode != want {
			t.Errorf("%s: status %d, want %d", path, res.StatusCode, want)
		}
	}
}

func TestEveryFontHasItsLicence(t *testing.T) {
	entries, err := os.ReadDir("fonts")
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		n := e.Name()
		if !strings.HasSuffix(n, ".woff2") {
			continue
		}
		licence := map[string]string{"hanken-grotesk": "OFL-HankenGrotesk.txt", "ibm-plex-mono": "OFL-IBMPlexMono.txt"}
		found := false
		for prefix, file := range licence {
			if strings.HasPrefix(n, prefix) {
				if _, err := os.Stat(filepath.Join("fonts", file)); err == nil {
					found = true
				}
			}
		}
		if !found {
			t.Errorf("%s has no licence file beside it", n)
		}
	}
}

// A lazy section draws its skeleton with the address of its own HTML, the error box it shows when that
// fails, and a link to the whole page for a browser without JavaScript; the page sends such a browser there.
func TestALazySectionDrawsItsPlaceholderAndItsFailedState(t *testing.T) {
	var b bytes.Buffer
	if err := kit.Render(&b, "ui/section", kit.Section{ID: "prs", Title: "Open pull requests", Src: "/part/prs?pr_sort=newest", Full: "/?full=1", Body: "<p>never drawn</p>"}); err != nil {
		t.Fatal(err)
	}
	got := b.String()
	for _, want := range []string{`<div class="ui-lazy" data-ui-src="/part/prs?pr_sort=newest">`, `aria-busy="true"`, `class="ui-skel-group"`,
		`<noscript><p class="ui-hint">This section loads with JavaScript. <a href="/?full=1">Show the whole page</a>.</p></noscript>`,
		`<template data-ui-failed>`, `Can’t load Open pull requests.`, `<p data-ui-cause></p>`, `href="/part/prs?pr_sort=newest" data-ui-retry>Retry</a>`} {
		if !strings.Contains(got, want) {
			t.Errorf("lazy section lacks %q:\n%s", want, got)
		}
	}
	if strings.Contains(got, "never drawn") {
		t.Error("a lazy section draws no body of its own")
	}
	b.Reset()
	if err := kit.Render(&b, "ui/shell", kit.Page{Bar: kit.Bar{App: "x"}, Full: "/?full=1"}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(b.String(), `<noscript><meta http-equiv="refresh" content="0; url=/?full=1"></noscript>`) {
		t.Errorf("the page sends a browser without JavaScript to the whole page:\n%s", b.String())
	}
}

var anchor = regexp.MustCompile(`<a\s[^>]*>`)

// Every link to another site or app opens in a new tab; links inside the app stay in the same tab.
func TestLinksToOtherSitesOpenInANewTab(t *testing.T) {
	files, err := filepath.Glob(filepath.Join("testdata", "golden", "*.html"))
	if err != nil || len(files) == 0 {
		t.Fatalf("no goldens: %v", err)
	}
	for _, f := range files {
		b, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		for _, a := range anchor.FindAllString(string(b), -1) {
			out := strings.Contains(a, `href="http://`) || strings.Contains(a, `href="https://`)
			if out != strings.Contains(a, `target="_blank"`) {
				t.Errorf("%s: %s", filepath.Base(f), a)
			}
			if out && !strings.Contains(a, `rel="noopener noreferrer"`) {
				t.Errorf("%s: a new tab without noopener: %s", filepath.Base(f), a)
			}
		}
	}
}

// A light is a link to its page; its popover, drawn hidden, holds the facts, a menu of the actions, the
// step an action asks before it posts, and the dialogs its actions open.
func TestALightIsALinkWithAHiddenPopoverOfFactsAndActions(t *testing.T) {
	light := kit.Light{ID: "gh", Tone: "good", Icon: "github", Label: "GitHub", Href: "/settings/github", Title: "Connected as bot, read-only.", Text: "<b>x</b>",
		Details: []kit.Detail{{Label: "Account", Value: "bot"}, {Label: "Token from", Value: "the Keychain"}},
		Actions: []kit.Action{{Label: "Enter a new token…", Opens: "gh-token"}, {Label: "Settings", Href: "/settings/github"},
			{Label: "Disconnect", Variant: "danger", Confirm: &kit.Confirm{Text: "Removes the token.", Label: "Disconnect now", Action: "/github/disconnect", Hidden: []kit.Hidden{{Name: "back", Value: "/roads"}}}}},
		Dialogs: []kit.Dialog{{ID: "gh-token", Title: "Enter a new token", Action: "/github/connect", Body: "<p>steps</p>", Confirm: kit.Button{Label: "Connect", Type: "submit"}}}}
	var b bytes.Buffer
	if err := kit.Render(&b, "ui/light", light); err != nil {
		t.Fatal(err)
	}
	got := b.String()
	for _, want := range []string{`<div class="ui-light" id="gh">`, `<a class="ui-light-btn" href="/settings/github" data-ui-light aria-controls="gh-pop" aria-expanded="false" aria-haspopup="dialog" title="Connected as bot, read-only.">`,
		`<span class="ui-dot" data-tone="good"></span>`, `<span class="ui-light-word">GitHub</span><span class="ui-sr">: Connected as bot, read-only.</span>`,
		`<div class="ui-pop" id="gh-pop" role="dialog" aria-labelledby="gh-pop-title" tabindex="-1" hidden>`, `&lt;b&gt;x&lt;/b&gt;`,
		`<dt>Account</dt><dd>bot</dd><dt>Token from</dt><dd>the Keychain</dd>`, `<details class="ui-menu"><summary class="ui-btn" data-size="sm">Actions</summary>`,
		`<button class="ui-menu-item" type="button" data-ui-open="gh-token" aria-haspopup="dialog">Enter a new token…</button>`, `<a class="ui-menu-item" href="/settings/github">Settings</a>`,
		`<button class="ui-menu-item" type="button" data-variant="danger" data-ui-step aria-controls="gh-step-2" aria-expanded="false">Disconnect</button>`,
		`<form class="ui-pop-step" id="gh-step-2" method="post" action="/github/disconnect" hidden><p>Removes the token.</p><input type="hidden" name="back" value="/roads">`,
		`<button class="ui-btn" type="button" data-size="sm" data-ui-step-cancel>Cancel</button><button class="ui-btn" type="submit" data-variant="danger" data-size="sm">Disconnect now</button>`,
		`<dialog class="ui-dialog" id="gh-token"`, `<form method="post" action="/github/connect">`, `<p>steps</p><div class="ui-dialog-acts">`} {
		if !strings.Contains(got, want) {
			t.Errorf("light lacks %q:\n%s", want, got)
		}
	}
	light.Open = true
	b.Reset()
	if err := kit.Render(&b, "ui/light", light); err != nil {
		t.Fatal(err)
	}
	if got := b.String(); !strings.Contains(got, `aria-expanded="true"`) || !strings.Contains(got, `tabindex="-1">`) || strings.Contains(got, `tabindex="-1" hidden>`) {
		t.Errorf("open, the popover is still hidden:\n%s", got)
	}
}

// A lazy light is drawn grey at once, a plain link to its page, with the light it shows when its fetch fails.
func TestALazyLightIsDrawnNeutralFirst(t *testing.T) {
	var b bytes.Buffer
	if err := kit.Render(&b, "ui/light", kit.Light{ID: "gh", Tone: "good", Icon: "github", Label: "GitHub", Href: "/settings/github", Title: "never drawn", Src: "/part/light/github?back=%2F"}); err != nil {
		t.Fatal(err)
	}
	got := b.String()
	for _, want := range []string{`<div class="ui-lazy" data-ui-src="/part/light/github?back=%2F">`, `<a class="ui-light-btn" href="/settings/github" aria-busy="true">`, `<span class="ui-dot" data-tone="none"></span>`,
		`<template data-ui-failed>`, `Can’t load GitHub.`, `<p class="ui-pop-text" data-ui-cause></p>`, `href="/part/light/github?back=%2F" data-ui-retry>Retry</a>`} {
		if !strings.Contains(got, want) {
			t.Errorf("lazy light lacks %q:\n%s", want, got)
		}
	}
	before, _, _ := strings.Cut(got, "<template")
	if strings.Contains(before, "never drawn") || strings.Contains(before, "ui-pop") || strings.Contains(before, `data-tone="good"`) {
		t.Errorf("the placeholder shows what is not loaded yet:\n%s", got)
	}
}

func TestAHiddenBackValueIsMarkedForKitJS(t *testing.T) {
	var b bytes.Buffer
	l := kit.Light{ID: "gh", Label: "GitHub", Href: "/s", Actions: []kit.Action{{Label: "Disconnect", Confirm: &kit.Confirm{Text: "x", Label: "Disconnect now", Action: "/d",
		Hidden: []kit.Hidden{{Name: "back", Value: "/?pr_sort=newest", Back: true}, {Name: "tab", Value: "prs"}}}}}}
	if err := kit.Render(&b, "ui/light", l); err != nil {
		t.Fatal(err)
	}
	got := b.String()
	if !strings.Contains(got, `<input type="hidden" name="back" value="/?pr_sort=newest" data-ui-back>`) || !strings.Contains(got, `<input type="hidden" name="tab" value="prs">`) {
		t.Errorf("only a Back value carries data-ui-back:\n%s", got)
	}
}

// A change asks the way this app says, else the global file, else its type's default, and says where that came from.
func TestAButtonAsksTheWayTheLayersSay(t *testing.T) {
	app := map[string]string{"move": kit.AskNever, "run": "bogus"}
	global := map[string]string{"move": kit.AskName, "edit": kit.AskName, "run": kit.AskNever}
	del := kit.Button{Label: "Delete", Type: "submit"}
	asked := kit.Button{Label: "Delete", Type: "submit", Confirm: "Delete 7?", Remember: true}
	named := asked
	named.ConfirmName = "7"
	for typ, want := range map[string]struct {
		source string
		button kit.Button
	}{
		"move":   {kit.FromApp, del},
		"edit":   {kit.FromGlobal, named},
		"run":    {kit.FromGlobal, del},
		"delete": {kit.FromDefault, named},
		"focus":  {kit.FromDefault, asked},
	} {
		style, source := kit.ResolveConfirm(app, global, typ)
		if got := del.Ask(style, "Delete 7?", "7"); source != want.source || got != want.button {
			t.Errorf("%s: got %s %+v, want %s %+v", typ, source, got, want.source, want.button)
		}
	}
}

// Each scope of the Settings form posts back what it shows, with its scope, and only real types and choices are read.
func TestConfirmSettingsPostsBackWhatItShows(t *testing.T) {
	app, global := map[string]string{"move": kit.AskNever}, map[string]string{"edit": kit.AskName}
	for scope, want := range map[string]map[string]string{
		kit.ScopeApp:    {"move": kit.AskNever, "edit": kit.AskInherit, "delete": kit.AskInherit},
		kit.ScopeGlobal: {"move": kit.AskInherit, "edit": kit.AskName, "delete": kit.AskInherit},
	} {
		var b bytes.Buffer
		if err := kit.Render(&b, "ui/confirm-settings", kit.ConfirmSettings{Action: "/settings", Scope: scope, App: app, Global: global}); err != nil {
			t.Fatal(err)
		}
		form := url.Values{}
		for _, m := range regexp.MustCompile(`name="(confirm\.[a-z]+)" value="([a-z-]+)"( checked)?`).FindAllStringSubmatch(b.String(), -1) {
			if m[3] != "" || m[1] == "confirm.scope" {
				form.Set(m[1], m[2])
			}
		}
		form.Set("confirm.bogus", kit.AskNever)
		gotScope, got := kit.ConfirmChoices(form)
		if gotScope != scope || len(got) != len(kit.ActionTypes) || got["move"] != want["move"] || got["edit"] != want["edit"] || got["delete"] != want["delete"] {
			t.Errorf("%s: scope %s, choices %v", scope, gotScope, got)
		}
	}
}
