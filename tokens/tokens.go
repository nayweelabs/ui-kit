// Package tokens holds the kit's design tokens, written once in tokens.json.
package tokens

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
)

//go:embed tokens.json
var source []byte

// Color is a token with one value per theme.
type Color struct {
	Name  string `json:"name"`
	Light string `json:"light"`
	Dark  string `json:"dark"`
	Use   string `json:"use"`
}

// Value is a token that is the same in both themes.
type Value struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	Use   string `json:"use"`
}

// Set is every token, in the order tokens.json lists them.
type Set struct {
	Colors []Color `json:"colors"`
	Values []Value `json:"values"`
}

// Default returns the tokens shipped in this module's tokens.json.
func Default() Set {
	set, err := Parse(source)
	if err != nil {
		panic(err)
	}
	return set
}

// Parse reads a tokens.json document and checks that every name is set and unique.
func Parse(data []byte) (Set, error) {
	var set Set
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&set); err != nil {
		return Set{}, fmt.Errorf("read tokens: %w", err)
	}
	seen := map[string]bool{}
	check := func(name string, values ...string) error {
		if name == "" || seen[name] {
			return fmt.Errorf("token name %q is empty or used twice", name)
		}
		seen[name] = true
		for _, v := range values {
			if v == "" {
				return fmt.Errorf("token %q has an empty value", name)
			}
		}
		return nil
	}
	for _, c := range set.Colors {
		if err := check(c.Name, c.Light, c.Dark); err != nil {
			return Set{}, err
		}
	}
	for _, v := range set.Values {
		if err := check(v.Name, v.Value); err != nil {
			return Set{}, err
		}
	}
	return set, nil
}

// Color returns the light and dark values of a colour token, and false when there is none.
func (s Set) Color(name string) (light, dark string, ok bool) {
	for _, c := range s.Colors {
		if c.Name == name {
			return c.Light, c.Dark, true
		}
	}
	return "", "", false
}

// CSS writes the tokens as --ui-* custom properties: light first, dark by system or [data-theme].
func (s Set) CSS() []byte {
	var b bytes.Buffer
	b.WriteString("/* Generated from tokens/tokens.json by `go run ./cmd/ui-kit-tokens`. Do not edit. */\n")
	b.WriteString("/* Light first; dark by system or [data-theme=\"dark\"] on :root or on any element. */\n")
	b.WriteString(":root,[data-theme=\"light\"] {\n")
	s.writeColors(&b, false)
	b.WriteString("  color-scheme:light;\n}\n")
	b.WriteString(":root {\n")
	for _, v := range s.Values {
		fmt.Fprintf(&b, "  --ui-%s:%s;\n", v.Name, v.Value)
	}
	b.WriteString("}\n")
	b.WriteString("@media (prefers-color-scheme: dark) {\n  :root:not([data-theme=\"light\"]) {\n")
	for _, c := range s.Colors {
		fmt.Fprintf(&b, "    --ui-%s:%s;\n", c.Name, c.Dark)
	}
	b.WriteString("    color-scheme:dark;\n  }\n}\n")
	b.WriteString(":root[data-theme=\"dark\"],[data-theme=\"dark\"] {\n")
	s.writeColors(&b, true)
	b.WriteString("  color-scheme:dark;\n}\n")
	return b.Bytes()
}

func (s Set) writeColors(b *bytes.Buffer, dark bool) {
	for _, c := range s.Colors {
		v := c.Light
		if dark {
			v = c.Dark
		}
		fmt.Fprintf(b, "  --ui-%s:%s;\n", c.Name, v)
	}
}
