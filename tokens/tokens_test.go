package tokens_test

import (
	"bytes"
	"math"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/ElderlyGate/ui-kit/tokens"
)

func TestShippedCSSIsGeneratedFromTokensJSON(t *testing.T) {
	onDisk, err := os.ReadFile("../css/tokens.css")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(onDisk, tokens.Default().CSS()) {
		t.Fatal("css/tokens.css is stale: run `go run ./cmd/ui-kit-tokens` and commit the result")
	}
}

func TestCSSHasEveryTokenInBothThemes(t *testing.T) {
	css := string(tokens.Default().CSS())
	for _, c := range tokens.Default().Colors {
		if got := strings.Count(css, "--ui-"+c.Name+":"); got != 3 {
			t.Errorf("--ui-%s is written %d times, want 3 (light, system dark, chosen dark)", c.Name, got)
		}
	}
	for _, v := range tokens.Default().Values {
		if !strings.Contains(css, "--ui-"+v.Name+":"+v.Value+";") {
			t.Errorf("--ui-%s is missing", v.Name)
		}
	}
}

func TestParseRefusesBadDocuments(t *testing.T) {
	cases := map[string]string{
		"duplicate name": `{"colors":[{"name":"bg","light":"#fff","dark":"#000"}],"values":[{"name":"bg","value":"1px"}]}`,
		"empty value":    `{"colors":[{"name":"bg","light":"#fff","dark":""}]}`,
		"unknown field":  `{"colours":[]}`,
		"not json":       `{`,
	}
	for name, doc := range cases {
		if _, err := tokens.Parse([]byte(doc)); err == nil {
			t.Errorf("%s: Parse accepted it", name)
		}
	}
}

func TestColorLooksUpBothThemes(t *testing.T) {
	light, dark, ok := tokens.Default().Color("good")
	if !ok || light != "#2F7D5A" || dark != "#6CC59A" {
		t.Fatalf("good = %q %q %v", light, dark, ok)
	}
	if _, _, ok := tokens.Default().Color("nope"); ok {
		t.Fatal("found a token that does not exist")
	}
}

// Text needs 4.5:1, shapes (dots, stars, rings, edges that carry meaning) need 3:1, in both themes.
func TestContrast(t *testing.T) {
	text := [][2]string{
		{"ink", "bg"}, {"ink", "surface"}, {"ink", "sunk"},
		{"muted", "bg"}, {"muted", "surface"}, {"muted", "sunk"},
		{"accent", "bg"}, {"accent", "surface"}, {"accent", "accent-soft"}, {"on-accent", "accent"},
		{"good", "surface"}, {"good", "good-soft"},
		{"warn", "surface"}, {"warn", "warn-soft"},
		{"bad", "surface"}, {"bad", "bad-soft"},
		{"info", "surface"}, {"info", "info-soft"},
		{"good", "bg"}, {"warn", "bg"}, {"bad", "bg"}, {"info", "bg"},
		{"surface", "warn"}, {"surface", "bad"},
		{"ink", "focus-soft"}, {"ink", "good-soft"}, {"ink", "bad-soft"},
		{"code-ink", "code-bg"}, {"code-key", "code-bg"}, {"code-str", "code-bg"},
	}
	shape := [][2]string{
		{"focus", "surface"}, {"focus", "bg"}, {"ring", "surface"}, {"ring", "bg"},
		{"faint", "surface"}, {"faint", "bg"},
		{"code-dim", "code-bg"},
	}
	set := tokens.Default()
	check := func(pairs [][2]string, min float64) {
		for _, p := range pairs {
			fl, fd, _ := set.Color(p[0])
			bl, bd, _ := set.Color(p[1])
			for theme, pair := range map[string][2]string{"light": {fl, bl}, "dark": {fd, bd}} {
				if r := ratio(t, pair[0], pair[1]); r < min {
					t.Errorf("%s: %s on %s is %.2f:1, want %.1f:1", theme, p[0], p[1], r, min)
				}
			}
		}
	}
	check(text, 4.5)
	check(shape, 3)
}

func ratio(t *testing.T, a, b string) float64 {
	la, lb := luminance(t, a), luminance(t, b)
	return (math.Max(la, lb) + 0.05) / (math.Min(la, lb) + 0.05)
}

func luminance(t *testing.T, hex string) float64 {
	t.Helper()
	if len(hex) != 7 || hex[0] != '#' {
		t.Fatalf("not a #RRGGBB colour: %q", hex)
	}
	var ch [3]float64
	for i := range ch {
		n, err := strconv.ParseUint(hex[1+2*i:3+2*i], 16, 8)
		if err != nil {
			t.Fatal(err)
		}
		c := float64(n) / 255
		if c <= 0.03928 {
			ch[i] = c / 12.92
		} else {
			ch[i] = math.Pow((c+0.055)/1.055, 2.4)
		}
	}
	return 0.2126*ch[0] + 0.7152*ch[1] + 0.0722*ch[2]
}
