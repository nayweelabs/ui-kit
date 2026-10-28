// Package kit is the org's shared UI kit for Go apps: the CSS, fonts and kit.js as an embedded file system, and html/template partials named "ui/…" that draw the shell and every part with the kit's classes.
package kit

import (
	"embed"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"net/http"
	"strings"
	"sync"
)

//go:generate go run ./cmd/ui-kit-tokens

//go:embed css fonts js partials
var files embed.FS

// Prefix is the path Assets serves the kit's files under.
const Prefix = "/ui-kit/"

// FS returns the kit's static files: css/, fonts/ and js/.
func FS() fs.FS {
	return files
}

// Assets serves css/, fonts/ and js/ under Prefix: mux.Handle(kit.Prefix, kit.Assets()).
func Assets() http.Handler {
	static := http.FileServerFS(files)
	return http.StripPrefix(strings.TrimSuffix(Prefix, "/"), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !servable(r.URL.Path) {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Cache-Control", "no-cache")
		static.ServeHTTP(w, r)
	}))
}

func servable(path string) bool {
	for _, dir := range []string{"/css/", "/fonts/", "/js/"} {
		if strings.HasPrefix(path, dir) && len(path) > len(dir) && !strings.HasSuffix(path, "/") {
			return true
		}
	}
	return false
}

// Templates returns a new template set holding every "ui/…" partial and the kit's functions. An app parses its own templates into it: kit.Templates().ParseFS(appFS, "views/*.html").
func Templates() *template.Template {
	return template.Must(base().Clone())
}

var base = sync.OnceValue(func() *template.Template {
	return template.Must(template.New("ui").Funcs(Funcs()).ParseFS(files, "partials/*.html"))
})

// Funcs are the functions the partials use; an app may use them in its own templates too.
func Funcs() template.FuncMap {
	return template.FuncMap{
		"pill":   func(tone, label string) Pill { return Pill{Tone: tone, Label: label} },
		"tag":    func(label string) Pill { return Pill{Variant: "tag", Label: label} },
		"status": func(tone, label string) Status { return Status{Tone: tone, Label: label} },
		"count":  func(tone string, n int) Count { return Count{Tone: tone, N: n} },
		"star":   func(text string, chip bool) Star { return Star{Text: text, Chip: chip} },
		"skeleton": func(shape string, lines int) Skeleton {
			return Skeleton{Shape: shape, Lines: lines}
		},
		"countText": func(n int) string {
			if n > 99 {
				return "99+"
			}
			return fmt.Sprint(n)
		},
		"last": func(i, n int) bool { return i == n-1 },
		"pct":  func(f float64) string { return fmt.Sprintf("%.2f%%", f*100) },
		"assets": func(p Page) string {
			if p.Assets == "" {
				return Prefix
			}
			return strings.TrimSuffix(p.Assets, "/") + "/"
		},
		"orDefault": func(fallback, v string) string {
			if v == "" {
				return fallback
			}
			return v
		},
	}
}

// Render draws one partial with its data, for apps that build a page from pieces.
func Render(w io.Writer, name string, data any) error {
	return Templates().ExecuteTemplate(w, name, data)
}

// HTML draws one partial into template.HTML, ready to be a Section's or a Page's Body.
func HTML(t *template.Template, name string, data any) (template.HTML, error) {
	var b strings.Builder
	if err := t.ExecuteTemplate(&b, name, data); err != nil {
		return "", err
	}
	return template.HTML(b.String()), nil // #nosec G203 -- output of html/template, already escaped
}
