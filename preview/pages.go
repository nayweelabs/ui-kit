package preview

import (
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"slices"
	"strings"
	"time"

	kit "github.com/ElderlyGate/ui-kit"
)

func fsSub(dir string) (fs.FS, error) { return fs.Sub(files, dir) }

// focus is the control room's focus page, built from the kit as the article draws it.
func focus(tpl *template.Template) (kit.Page, error) {
	prs := SamplePRs()
	var doNow []kit.PR
	for i, verb := range []string{"Approve", "Fix", "Review"} {
		p := prs[i]
		p.Compact, p.Verb = true, verb
		doNow = append(doNow, p)
	}
	open := kit.Section{ID: "open-prs", Title: "Open pull requests", Count: &kit.Count{Tone: "accent", N: len(prs)}, Why: "Yours first, then oldest.",
		Body: mustHTML(tpl, "ui/pr-list", kit.PRList{Items: prs}), More: &kit.Link{Label: "All 12 →", Href: "#"}}
	now := kit.Section{ID: "do-now", Title: "Do now", Why: "In the order the server sends.",
		Body: mustHTML(tpl, "ui/pr-list", kit.PRList{Items: doNow}) + mustHTML(tpl, "preview/waiting", nil)}
	body, err := kit.HTML(tpl, "preview/focus", map[string]template.HTML{
		"Miles": mustHTML(tpl, "ui/milestones", SampleMilestones()),
		"Open":  mustHTML(tpl, "ui/section", open),
		"Now":   mustHTML(tpl, "ui/section", now),
	})
	return kit.Page{Header: kit.Head{Title: "Focus", Lede: "The road to public, what waits for you, and every open pull request."}, Body: body}, err
}

// products is the new Products page; the server filters, the page only draws (G-STACK-7).
func products(tpl *template.Template, r *http.Request) (kit.Page, error) {
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	show := r.URL.Query().Get("show")
	view := r.URL.Query().Get("view")
	var items []kit.Product
	for _, p := range SampleProducts() {
		if q != "" && !strings.Contains(strings.ToLower(p.Name), strings.ToLower(q)) {
			continue
		}
		if show == "needs-me" && !(p.Focus || p.Marker != nil || p.Waiting != nil || p.Problem != nil) {
			continue
		}
		if show == "planned" && p.Variant != "planned" {
			continue
		}
		items = append(items, p)
	}
	tb := productToolbar(view, q, show)
	tb.Action = "/products"
	tb.Count = fmt.Sprintf("%d products", len(items))
	if q != "" || show != "" {
		tb.Clear = &kit.Link{Label: "Clear filters", Href: "/products"}
	}
	data := kit.Products{ID: "products", View: view, Toolbar: tb, Columns: ProductColumns, Items: items,
		Empty: &kit.Empty{Title: "Nothing matches.", Text: "The filter hides every product.", Action: &kit.Button{Label: "Clear filters", Variant: "quiet", Size: "sm", Href: "/products"}}}
	body, err := kit.HTML(tpl, "ui/products", data)
	return kit.Page{
		Header: kit.Head{Title: "Products", Lede: "Is it current, is anything open, is my copy behind.", Actions: []kit.Button{{Label: "Refresh", Href: "/products"}}},
		Notice: &kit.Notice{Tone: "warn", Title: "Update ready: kanban v0.2.0.", Text: "6 checkouts are behind; pull before you start work."},
		Body:   body,
	}, err
}

// admin shows the one exception: an app with more than 6 pages adds a left rail under the top bar.
func admin(tpl *template.Template, r *http.Request) (kit.Page, error) {
	link := func(label string, current bool) kit.Link {
		return kit.Link{Label: label, Href: "/admin", Current: current}
	}
	machines := kit.Table{Label: "Machines", Columns: []kit.Column{{Label: "Machine"}, {Label: "State"}, {Label: "Provider"}, {Label: "Uptime", Num: true}, {Label: "Drift", Num: true}}, Rows: []kit.Row{
		{Cells: []kit.Cell{{Text: "m-01", Href: "#", Sub: "api, worker"}, {Status: &kit.Status{Tone: "good", Label: "running"}}, {Text: "hetzner"}, {Text: "41d", Num: true}, {Text: "0", Num: true}}},
		{Cells: []kit.Cell{{Text: "m-02", Href: "#", Sub: "db"}, {Status: &kit.Status{Tone: "warn", Label: "drifted"}}, {Text: "hetzner"}, {Text: "12d", Num: true}, {Text: "3", Num: true}}},
		{Cells: []kit.Cell{{Text: "m-03", Href: "#", Sub: "edge"}, {Status: &kit.Status{Tone: "info", Label: "starting"}}, {Text: "local"}, {Text: "2m", Num: true}, {Text: "–", Num: true}}},
		{Error: "Can’t read machine m-04: the provider answered 503. The others are current."},
	}}
	panels := map[string]template.HTML{
		"Tabs": mustHTML(tpl, "ui/tabs", kit.Tabs{Label: "Machine views", Items: []kit.Tab{
			{Label: "Machines", Panel: "pane-machines", Selected: true, Count: &kit.Count{N: 4}},
			{Label: "Logs", Panel: "pane-logs"},
			{Label: "Settings", Panel: "pane-settings"},
		}}),
		"Machines": mustHTML(tpl, "ui/table", machines),
		"Actions":  machineActions(tpl, r),
		"Logs":     mustHTML(tpl, "ui/empty", kit.Empty{Title: "No logs in the last hour.", Text: "Machines write logs when a deploy runs.", Command: "fleet logs --since 1h"}),
		"Settings": mustHTML(tpl, "ui/field", kit.Field{ID: "provider-token", Label: "Provider token", Secret: true, Hint: "Kept in the Keychain. Never shown back after saving."}),
	}
	body, err := kit.HTML(tpl, "preview/admin", panels)
	return kit.Page{
		Wide: true,
		Rail: []kit.RailGroup{
			{Label: "Favourites", Links: []kit.Link{link("Deployments", false), link("Machines", true)}},
			{Label: "All services", Links: []kit.Link{link("Overview", false), link("Providers", false), link("Networks", false), link("Secrets", false), link("Backups", false), link("Logs", false), link("Activity", false), {Label: "Alerts", Href: "/admin", Count: &kit.Count{Tone: "bad", N: 2}}}},
		},
		Header: kit.Head{
			Crumb:   []kit.Link{{Label: "fleet", Href: "/admin"}, {Label: "Machines", Href: "/admin"}, {Label: "All machines"}},
			Title:   "Machines",
			Lede:    "An app with more than 6 pages keeps the top bar and adds this rail under it.",
			Actions: []kit.Button{{Label: "Open a PR in a drawer", Opens: "drw-live"}, {Label: "Remove m-03", Variant: "danger", Opens: "dlg-live"}},
		},
		Body: body + must(liveDialogs(tpl)),
	}, err
}

// machineActions restart m-02 (run) and delete m-04 (delete), each asking its type's default way; Remember my answer stops asking for the pressed button's type.
func machineActions(tpl *template.Template, r *http.Request) template.HTML {
	q := r.URL.Query()
	chosen := map[string]string{}
	for _, t := range []string{"run", "delete"} {
		if q.Get("confirm") == "remember" && q.Has(t) {
			chosen[t] = kit.AskNever
		}
	}
	restart := kit.Button{Label: "Restart m-02", Type: "submit", Name: "run", Value: "m-02"}
	remove := kit.Button{Label: "Delete m-04", Variant: "danger", Type: "submit", Name: "delete", Value: "m-04"}
	return mustHTML(tpl, "ui/button", restart.Ask(kit.ConfirmStyle(chosen, "run"), "Restart m-02?", "m-02")) + " " +
		mustHTML(tpl, "ui/button", remove.Ask(kit.ConfirmStyle(chosen, "delete"), "Delete machine m-04?", "m-04"))
}

// lazySections are the sections of the lazy page: one quick with a toolbar, one slow, one that fails.
var lazySections = []kit.Section{
	{ID: "lazy-prs", Title: "Open pull requests", Why: "Quick; its toolbar reloads only this section."},
	{ID: "lazy-slow", Title: "Slow source", Why: "Takes a second and a half."},
	{ID: "lazy-fails", Title: "Failing source", Why: "Its server answers 500."},
}

// lazy is a page whose sections each load on their own; with ?full=1 the server draws it whole, as a
// browser without JavaScript gets it.
func lazy(tpl *template.Template, r *http.Request) (kit.Page, error) {
	full := r.URL.Query().Get("full") == "1"
	var body template.HTML
	for _, sec := range lazySections {
		name := strings.TrimPrefix(sec.ID, "lazy-")
		if full {
			h, _, err := lazyPart(tpl, name, r, false)
			if err != nil {
				h = mustHTML(tpl, "ui/section", kit.Section{ID: sec.ID, Title: sec.Title, Problem: &kit.Problem{Title: "Can’t read the " + strings.ToLower(sec.Title) + ".", Cause: err.Error()}})
			}
			body += h + "\n"
			continue
		}
		q := r.URL.Query()
		sec.Src, sec.Full = "/part/"+name+"?"+q.Encode(), "/lazy?full=1"
		again := r.URL.Query()
		again.Set("fresh", "1")
		sec.Again = "/lazy?" + again.Encode()
		body += mustHTML(tpl, "ui/section", sec) + "\n"
	}
	page := kit.Page{Header: kit.Head{Title: "Lazy sections", Lede: "The shell comes at once; each section loads its own HTML, in parallel, and fails on its own."}, Body: body}
	if !full {
		page.Full = "/lazy?full=1"
	}
	return page, nil
}

// part serves one section of the lazy page as HTML.
func part(tpl *template.Template, w http.ResponseWriter, r *http.Request) {
	h, code, err := lazyPart(tpl, r.PathValue("name"), r, true)
	if err != nil {
		http.Error(w, err.Error(), code)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("UI-Read-At", time.Now().Format(time.RFC3339))
	w.Write([]byte(h))
}

func lazyPart(tpl *template.Template, name string, r *http.Request, wait bool) (template.HTML, int, error) {
	switch name {
	case "prs":
		prs := SamplePRs()
		sort := r.URL.Query().Get("sort")
		if sort == "newest" {
			slices.Reverse(prs)
		}
		tb := kit.Toolbar{Action: "/lazy", Count: fmt.Sprintf("%d open", len(prs)), Filters: []kit.Select{{Name: "sort", Label: "Sort", Active: sort != "", Options: []kit.Option{
			{Value: "", Label: "Oldest first", Selected: sort == ""}, {Value: "newest", Label: "Newest first", Selected: sort == "newest"}}}}}
		body := mustHTML(tpl, "ui/toolbar", tb) + mustHTML(tpl, "ui/pr-list", kit.PRList{Items: prs})
		h, err := kit.HTML(tpl, "ui/section", kit.Section{ID: "lazy-prs", Title: "Open pull requests", Count: &kit.Count{Tone: "accent", N: len(prs)}, Body: body})
		return h, http.StatusInternalServerError, err
	case "slow":
		if wait {
			time.Sleep(1500 * time.Millisecond)
		}
		h, err := kit.HTML(tpl, "ui/section", kit.Section{ID: "lazy-slow", Title: "Slow source", Body: mustHTML(tpl, "ui/milestones", SampleMilestones())})
		return h, http.StatusInternalServerError, err
	case "fails":
		return "", http.StatusInternalServerError, fmt.Errorf("the sample source always fails")
	}
	return "", http.StatusNotFound, fmt.Errorf("no section %s", name)
}
