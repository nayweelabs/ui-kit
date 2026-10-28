// Package preview is the kit's gallery: the shell and every part in every state, light and dark, plus three example pages built only from the kit. It is what the owner looks at before a release.
package preview

import (
	"embed"
	"html/template"
	"io/fs"
	"net/http"
	"strings"

	kit "github.com/ElderlyGate/ui-kit"
)

//go:embed templates static
var files embed.FS

// Version is shown in the preview's footer.
var Version = "dev"

// Handler serves the gallery at /, the example pages, the kit's files under /ui-kit/ and the preview's own CSS.
func Handler() http.Handler {
	tpl := template.Must(kit.Templates().ParseFS(files, "templates/*.html"))
	mux := http.NewServeMux()
	mux.Handle(kit.Prefix, kit.Assets())
	mux.Handle("GET /preview/", http.StripPrefix("/preview", http.FileServerFS(sub("static"))))
	pages := map[string]func(*http.Request) (kit.Page, error){
		"/":         func(r *http.Request) (kit.Page, error) { return gallery(tpl, r) },
		"/focus":    func(r *http.Request) (kit.Page, error) { return focus(tpl) },
		"/products": func(r *http.Request) (kit.Page, error) { return products(tpl, r) },
		"/admin":    func(r *http.Request) (kit.Page, error) { return admin(tpl, r) },
	}
	for path, build := range pages {
		pattern := "GET " + path
		if path == "/" {
			pattern = "GET /{$}"
		}
		mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
			page, err := build(r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			page.Head = `<link rel="stylesheet" href="/preview/preview.css">`
			page.Bar = bar(path)
			page.Footer = footer()
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			if err := tpl.ExecuteTemplate(w, "ui/shell", page); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		})
	}
	return mux
}

func sub(dir string) fs.FS {
	return must(fsSub(dir))
}

func bar(current string) kit.Bar {
	link := func(label, href string) kit.Link { return kit.Link{Label: label, Href: href, Current: href == current} }
	return kit.Bar{
		App: "UI kit", Sub: "preview", Href: "/",
		Pages:  []kit.Link{link("Gallery", "/"), link("Control room", "/focus"), link("Products", "/products"), link("Left rail", "/admin")},
		Status: []kit.Status{{Tone: "good", Label: "local"}},
	}
}

func footer() kit.Footer {
	return kit.Footer{
		Version: "ui-kit " + Version,
		Fresh:   kit.Status{Tone: "good", Label: "sample data, shaped like the real thing"},
		Links:   []kit.Link{{Label: "Tokens", Href: "/ui-kit/css/tokens.css"}, {Label: "kit.css", Href: "/ui-kit/css/kit.css"}, {Label: "kit.js", Href: "/ui-kit/js/kit.js"}},
	}
}

// galleryCase is one case drawn once per theme.
type galleryCase struct {
	Case
	Light template.HTML
	Dark  template.HTML
}

type galleryPart struct {
	Part
	Cases []galleryCase
}

type galleryData struct {
	Both  bool
	Parts []galleryPart
	Live  template.HTML
}

func gallery(tpl *template.Template, r *http.Request) (kit.Page, error) {
	data := galleryData{Both: r.URL.Query().Get("themes") != "one"}
	byPart := map[string][]galleryCase{}
	for _, c := range Cases() {
		var b strings.Builder
		for _, it := range c.Items {
			h, err := kit.HTML(tpl, it.Template, it.Data)
			if err != nil {
				return kit.Page{}, err
			}
			b.WriteString(string(h))
			b.WriteString("\n")
		}
		light := b.String()
		byPart[c.Part] = append(byPart[c.Part], galleryCase{Case: c, Light: template.HTML(light), Dark: template.HTML(darkIDs(light))})
	}
	only := r.URL.Query().Get("part")
	for _, p := range Parts() {
		if only == "" || only == p.ID {
			data.Parts = append(data.Parts, galleryPart{Part: p, Cases: byPart[p.ID]})
		}
	}
	live, err := liveDialogs(tpl)
	if err != nil {
		return kit.Page{}, err
	}
	data.Live = live
	body, err := kit.HTML(tpl, "preview/gallery", data)
	if err != nil {
		return kit.Page{}, err
	}
	return kit.Page{
		Title: "UI kit",
		Header: kit.Head{
			Title: "Every part, every state",
			Lede:  "The shell and the 20 parts of the kit, each drawn by its Go partial with sample data. Left light, right dark; the theme button switches the whole page.",
			Actions: []kit.Button{
				{Label: "Both themes", Href: "/", Variant: primaryIf(data.Both)},
				{Label: "One theme", Href: "/?themes=one", Variant: primaryIf(!data.Both)},
			},
		},
		Body: body,
	}, nil
}

func primaryIf(on bool) string {
	if on {
		return "primary"
	}
	return ""
}

// darkIDs gives the dark copy of a case its own ids, so labels and controls still point at their own copy.
func darkIDs(h string) string {
	for _, attr := range []string{`id="`, `for="`, `aria-controls="`, `aria-describedby="`, `aria-labelledby="`, `data-ui-reveal="`} {
		h = strings.ReplaceAll(h, attr, attr+"dark-")
	}
	return h
}

func liveDialogs(tpl *template.Template) (template.HTML, error) {
	var b strings.Builder
	for _, d := range []kit.Dialog{
		{ID: "dlg-live", Title: "Disconnect GitHub?", Text: "Pull requests and CI stop showing. The token is deleted from the Keychain.", Confirm: kit.Button{Label: "Disconnect GitHub", Variant: "danger", Type: "submit", Value: "confirm"}},
		{ID: "dlg-archive", Title: "Archive this card?", Text: "It moves to archive/ and leaves the board. kanban restore brings it back.", Confirm: kit.Button{Label: "Archive 1 card", Variant: "danger", Type: "submit", Value: "confirm"}},
	} {
		h, err := kit.HTML(tpl, "ui/dialog", d)
		if err != nil {
			return "", err
		}
		b.WriteString(string(h))
	}
	h, err := kit.HTML(tpl, "ui/drawer", kit.Drawer{ID: "drw-live", Title: "control-room #2", Body: mustHTML(tpl, "ui/pr-item", SamplePRs()[1])})
	if err != nil {
		return "", err
	}
	b.WriteString(string(h))
	return template.HTML(b.String()), nil
}

func mustHTML(tpl *template.Template, name string, data any) template.HTML {
	return must(kit.HTML(tpl, name, data))
}

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
