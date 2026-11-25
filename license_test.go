package kit_test

import (
	"bytes"
	"encoding/json"
	"html/template"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/ElderlyGate/contracts"
	kit "github.com/ElderlyGate/ui-kit"
)

// badgeWords is ADR-L3's word for each license show state.
var badgeWords = map[string]string{"none": "Free", "trial-ended": "Free", "ended": "Free", "trial": "Trial", "licensed": "Licensed", "ending": "Licensed"}

var tag = regexp.MustCompile(`<[^>]*>|\s+`)

// licenseSamples is every good license show answer of the contracts module this kit pins, by file name.
func licenseSamples(t *testing.T) map[string][]byte {
	t.Helper()
	out, err := exec.Command("go", "list", "-m", "-f", "{{.Dir}}", "github.com/ElderlyGate/contracts").Output()
	if err != nil {
		t.Fatalf("go list the contracts module: %v", err)
	}
	files, err := filepath.Glob(filepath.Join(strings.TrimSpace(string(out)), "testdata", "good", "mankit", "v1", "license-show", "*.json"))
	if err != nil || len(files) == 0 {
		t.Fatalf("no license show samples in the contracts module: %v", err)
	}
	samples := map[string][]byte{}
	for _, f := range files {
		raw, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		if err := contracts.Validate("mankit/v1/license-show.schema.json", raw); err != nil {
			t.Fatalf("%s: %v", f, err)
		}
		samples[strings.TrimSuffix(filepath.Base(f), ".json")] = raw
	}
	return samples
}

func draw(t *testing.T, name string, data any) string {
	t.Helper()
	var b bytes.Buffer
	if err := kit.Render(&b, name, data); err != nil {
		t.Fatal(err)
	}
	return b.String()
}

// Every contracts sample answer draws its state's word in the answer's colour, and its panel shows the sentence as given and no text of its own; React draws the same (testdata/golden/license-samples.json).
func TestGoldenLicenseSamplesDrawWordColourAndSentence(t *testing.T) {
	type sample struct {
		Data  json.RawMessage `json:"data"`
		Badge string          `json:"badge"`
		Panel string          `json:"panel"`
	}
	all := map[string]sample{}
	for name, raw := range licenseSamples(t) {
		var answer struct {
			Data json.RawMessage `json:"data"`
		}
		var l kit.License
		if err := json.Unmarshal(raw, &answer); err != nil {
			t.Fatal(err)
		}
		if err := json.Unmarshal(answer.Data, &l); err != nil {
			t.Fatal(err)
		}
		tone := ""
		if l.Warn {
			tone = ` data-tone="warn"`
		}
		badge := draw(t, "ui/license-badge", l)
		if want := `<span class="ui-pill ui-license-badge"` + tone + `>` + badgeWords[l.State] + `</span>`; badge != want {
			t.Errorf("%s: badge %s, want %s", name, badge, want)
		}
		panel := draw(t, "ui/license-panel", l)
		if l.Sentence != "" && strings.Count(panel, `<p class="ui-license-sentence">`+template.HTMLEscapeString(l.Sentence)+`</p>`) != 1 {
			t.Errorf("%s: the panel does not show the sentence %q as given:\n%s", name, l.Sentence, panel)
		}
		bare := l
		bare.Sentence, bare.Features = "", nil
		if text := tag.ReplaceAllString(draw(t, "ui/license-panel", bare), ""); text != "License"+badgeWords[l.State] {
			t.Errorf("%s: without its sentence and features the panel says %q, which it made itself", name, text)
		}
		bar := draw(t, "ui/bar", kit.Bar{App: "Designer", License: &kit.LicenseLight{Href: "/settings", License: l}})
		if !strings.Contains(bar, `aria-haspopup="dialog" title="License"><span class="ui-sr">License: </span>`+badge+`</a>`) || !strings.Contains(bar, panel) {
			t.Errorf("%s: the bar lacks the badge or its panel:\n%s", name, bar)
		}
		all[name] = sample{Data: answer.Data, Badge: badge, Panel: panel}
	}
	got, err := json.MarshalIndent(all, "", " ")
	if err != nil {
		t.Fatal(err)
	}
	got = append(got, '\n')
	path := filepath.Join("testdata", "golden", "license-samples.json")
	if *update {
		if err := os.WriteFile(path, got, 0o644); err != nil {
			t.Fatal(err)
		}
		return
	}
	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("%v (run go test -run Golden -update to write it)", err)
	}
	if !bytes.Equal(got, want) {
		t.Error("testdata/golden/license-samples.json differs from the contracts samples; run go test -run Golden -update")
	}
}
