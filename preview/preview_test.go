package preview_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ElderlyGate/ui-kit/preview"
)

func get(t *testing.T, srv *httptest.Server, path string) string {
	t.Helper()
	res, err := http.Get(srv.URL + path)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	if res.StatusCode != http.StatusOK {
		t.Fatalf("%s: status %d: %s", path, res.StatusCode, body)
	}
	return string(body)
}

func TestEveryPageDraws(t *testing.T) {
	srv := httptest.NewServer(preview.Handler())
	defer srv.Close()
	for path, want := range map[string]string{
		"/":                    "Every part, every state",
		"/?themes=one":         `class="pv-frame"><div`,
		"/focus":               "Road to public",
		"/products":            "Is it current, is anything open",
		"/admin":               `class="ui-rail"`,
		"/preview/preview.css": ".pv-part",
		"/ui-kit/css/kit.css":  ".ui-card",
	} {
		if body := get(t, srv, path); !strings.Contains(body, want) {
			t.Errorf("%s lacks %q", path, want)
		}
	}
}

func TestGalleryDrawsEveryPartInBothThemes(t *testing.T) {
	srv := httptest.NewServer(preview.Handler())
	defer srv.Close()
	body := get(t, srv, "/")
	for _, p := range preview.Parts() {
		if !strings.Contains(body, `id="`+p.ID+`"`) {
			t.Errorf("gallery lacks part %s", p.ID)
		}
	}
	if l, d := strings.Count(body, `data-theme="light"`), strings.Count(body, `data-theme="dark"`); l != len(preview.Cases()) || d != l {
		t.Errorf("light frames %d, dark frames %d, cases %d", l, d, len(preview.Cases()))
	}
	if strings.Contains(body, `id="dark-dark-`) {
		t.Error("a dark id was prefixed twice")
	}
}

func TestProductsPageFiltersOnTheServer(t *testing.T) {
	srv := httptest.NewServer(preview.Handler())
	defer srv.Close()
	body := get(t, srv, "/products?q=hub&view=list")
	if !strings.Contains(body, "1 products") || strings.Contains(body, `href="#kanban"`) {
		t.Error("filter q=hub did not leave only hubman")
	}
	if !strings.Contains(body, `data-view="cards" hidden`) {
		t.Error("view=list did not hide the cards")
	}
	if body := get(t, srv, "/products?q=nothing"); !strings.Contains(body, "Nothing matches.") {
		t.Error("an empty filter result does not say Nothing matches")
	}
}
