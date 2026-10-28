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

// The lazy page draws each section's placeholder at once, serves each section on its own, and draws
// the whole page on the server for a browser without JavaScript.
func TestTheLazyPageLoadsEachSectionOnItsOwn(t *testing.T) {
	srv := httptest.NewServer(preview.Handler())
	defer srv.Close()
	page := get(t, srv, "/lazy")
	for _, want := range []string{`data-ui-src="/part/prs?"`, `data-ui-src="/part/slow?"`, `data-ui-src="/part/fails?"`, `content="0; url=/lazy?full=1"`} {
		if !strings.Contains(page, want) {
			t.Errorf("/lazy lacks %q", want)
		}
	}
	if part := get(t, srv, "/part/prs?sort=newest"); !strings.Contains(part, `id="lazy-prs"`) || !strings.Contains(part, `<option value="newest" selected>`) {
		t.Errorf("the PR section with its sort: %s", part)
	}
	res, err := http.Get(srv.URL + "/part/fails")
	if err != nil {
		t.Fatal(err)
	}
	res.Body.Close()
	if res.StatusCode != http.StatusInternalServerError {
		t.Errorf("the failing section answers %d", res.StatusCode)
	}
	full := get(t, srv, "/lazy?full=1")
	if strings.Contains(full, "data-ui-src") || !strings.Contains(full, `id="lazy-prs"`) || !strings.Contains(full, "the sample source always fails") || strings.Contains(full, "http-equiv") {
		t.Errorf("the whole page: %s", full)
	}
}
