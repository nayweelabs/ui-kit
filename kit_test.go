package kit_test

import (
	"bytes"
	"flag"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
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
