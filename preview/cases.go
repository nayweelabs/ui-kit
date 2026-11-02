package preview

import (
	"html/template"

	kit "github.com/ElderlyGate/ui-kit"
)

// Item is one partial drawn with its data.
type Item struct {
	Template string
	Data     any
}

// Case is one state of one part. The gallery draws every case; the golden tests check every case.
type Case struct {
	Part  string
	State string
	Note  string
	Stage bool
	Items []Item
}

// Part describes one part for the gallery heading.
type Part struct {
	ID      string
	Title   string
	Class   string
	Purpose string
	Rules   []string
}

func one(name string, data any) []Item { return []Item{{name, data}} }

func html(name string, data any) template.HTML {
	h, err := kit.HTML(kit.Templates(), name, data)
	if err != nil {
		panic(err)
	}
	return h
}

// Sample data shaped like the org's real products and pull requests; not live data.

func hubman() kit.Product {
	return kit.Product{
		Name: "hubman", Href: "#hubman", Kind: "tool · github-release", Focus: true,
		Facts: []kit.Fact{
			{Label: "Installed", Tone: "good", Num: &kit.Num{Text: "v0.7.0", Href: "#release", Command: "hubman list --installed"}, Text: "current"},
			{Label: "Board", Num: &kit.Num{Text: "7", Href: "#board", Command: "kanban cards --board hubman --open"}, Text: "open", Pill: &kit.Pill{Tone: "warn", Label: "9 waiting"}},
			{Label: "Pull requests", Num: &kit.Num{Text: "1", Href: "#prs"}, Text: "open · CI", Pill: &kit.Pill{Tone: "good", Label: "passing"}},
			{Label: "My copy", Tone: "warn", Num: &kit.Num{Text: "3", Href: "#copy", Command: "git -C hubman status -sb"}, Text: "behind"},
		},
		Foot: "Next: v0.8.0 · Rename window",
	}
}

func kanbanProduct() kit.Product {
	return kit.Product{
		Name: "kanban", Href: "#kanban", Kind: "tool · github-release", Marker: &kit.Pill{Tone: "warn", Label: "update ready"},
		Facts: []kit.Fact{
			{Label: "Installed", Tone: "warn", Num: &kit.Num{Text: "0.1.0"}, Text: "→ v0.2.0"},
			{Label: "Board", Num: &kit.Num{Text: "1", Href: "#board"}, Text: "open"},
			{Label: "Pull requests", Num: &kit.Num{Text: "0"}, Text: "open · CI", Pill: &kit.Pill{Tone: "good", Label: "passing"}},
			{Label: "My copy", Tone: "bad", Text: "on", Branch: "fix/agent-test-git-identity"},
		},
		Foot: "Not managed by hubman",
	}
}

func checkman() kit.Product {
	return kit.Product{
		Name: "checkman", Href: "#checkman", Kind: "tool · github-release", Waiting: &kit.Count{Tone: "accent", N: 1, Label: "1 pull request waits for you"},
		Facts: []kit.Fact{
			{Label: "Installed", Tone: "good", Num: &kit.Num{Text: "v1.9.2", Href: "#release"}, Text: "current"},
			{Label: "Board", Num: &kit.Num{Text: "2", Href: "#board"}, Text: "open", Pill: &kit.Pill{Tone: "warn", Label: "2 waiting"}},
			{Label: "Pull requests", Num: &kit.Num{Text: "1", Href: "#prs"}, Text: "open · CI", Pill: &kit.Pill{Tone: "good", Label: "passing"}},
			{Label: "My copy", Tone: "warn", Num: &kit.Num{Text: "2"}, Text: "behind"},
		},
		Foot: "Next: v2.0.0 · Private group",
	}
}

func designer() kit.Product {
	return kit.Product{
		Name: "designer", Href: "#designer", Kind: "tool · docker-checkout",
		Problem: &kit.Problem{Title: "Can’t read the board", Cause: "The board folder moved. Other facts below are current."},
		Facts: []kit.Fact{
			{Label: "Installed", Tone: "good", Num: &kit.Num{Text: "v4.11.0"}, Text: "current"},
			{Label: "Board", Error: "can’t read the board"},
			{Label: "Pull requests", Num: &kit.Num{Text: "2", Href: "#prs"}, Text: "open"},
			{Label: "My copy", Tone: "warn", Num: &kit.Num{Text: "3"}, Text: "behind · 2 changes"},
		},
		Foot: "Partial: 1 of 4 facts failed",
	}
}

func planned() kit.Product {
	return kit.Product{
		Name: "fleet", Href: "#fleet", Kind: "app · planned", Variant: "planned",
		Facts: []kit.Fact{
			{Label: "Installed", Tone: "none", Text: "planned"},
			{Label: "Board", Num: &kit.Num{Text: "12", Href: "#board"}, Text: "open"},
			{Label: "Pull requests", Num: &kit.Num{Text: "0"}, Text: "open"},
			{Label: "My copy", Tone: "good", Text: "on main, clean"},
		},
		Foot: "Next: v2 MVP 1 · Private group",
	}
}

func library() kit.Product {
	return kit.Product{
		Name: "mankit", Href: "#mankit", Kind: "library · go module", Variant: "library",
		Facts: []kit.Fact{
			{Label: "Installed", Tone: "info", Text: "library"},
			{Label: "Board", Num: &kit.Num{Text: "3", Href: "#board"}, Text: "open"},
			{Label: "Pull requests", Num: &kit.Num{Text: "1", Href: "#prs"}, Text: "open · CI", Pill: &kit.Pill{Tone: "info", Label: "running"}},
			{Label: "My copy", Tone: "good", Text: "current"},
		},
		Foot: "Latest: v0.9.1",
	}
}

// SampleProducts are the products the gallery and the example pages show.
func SampleProducts() []kit.Product {
	return []kit.Product{hubman(), kanbanProduct(), checkman(), designer(), planned(), library()}
}

// ProductColumns are the four facts, in order.
var ProductColumns = []string{"Installed", "Board", "Pull requests", "My copy"}

func productToolbar(view, q, show string) *kit.Toolbar {
	opt := func(v, l string) kit.Option { return kit.Option{Value: v, Label: l, Selected: v == show} }
	return &kit.Toolbar{
		Search:  &kit.Search{Placeholder: "Filter products", Value: q},
		Filters: []kit.Select{{Name: "show", Label: "Show", Active: show != "", Options: []kit.Option{opt("", "All products"), opt("needs-me", "Needs me"), opt("planned", "Planned")}}},
		Count:   "6 products", Target: "products", View: view,
	}
}

// SamplePRs are open pull requests across the org's repos.
func SamplePRs() []kit.PR {
	return []kit.PR{
		{Repo: "checkman", Number: "5", Href: "#pr-checkman-5", Author: "the org’s bot", Initials: "AI", Age: "opened 2d ago",
			Title: "chore: close the accepted review cards", Branch: "chore/accept-review-cards", Base: "main",
			Files: &kit.Num{Text: "3", Href: "#files"}, Added: "48", Removed: "12",
			Review: &kit.Pill{Tone: "accent", Label: "your review"}, CI: &kit.Pill{Tone: "good", Label: "CI passing"}},
		{Repo: "control-room", Number: "2", Href: "#pr-control-room-2", Author: "the org’s bot", Initials: "AI", Age: "opened 5d ago",
			Title: "feat: MVP 1, see what waits for the owner", Branch: "feat/mvp1", Base: "main",
			Files: &kit.Num{Text: "41", Href: "#files"}, Added: "2,914", Removed: "88", Conflict: "conflicts with main",
			Review: &kit.Pill{Tone: "warn", Label: "changes asked"}, CI: &kit.Pill{Tone: "bad", Label: "CI failing 1/6", Href: "#run"}},
		{Repo: "mankit", Number: "4", Href: "#pr-mankit-4", Draft: true, Author: "owner", Initials: "TA", Age: "updated 3h ago",
			Title: "feat: help layout shared by every tool", Branch: "feat/help-layout", Base: "main",
			Files: &kit.Num{Text: "7", Href: "#files"}, Added: "301", Removed: "54",
			Review: &kit.Pill{Label: "no review yet"}, CI: &kit.Pill{Tone: "info", Label: "CI running"}},
	}
}

// lightGood is the GitHub light of a connected, read-only token.
func lightGood(id string, open bool) kit.Light {
	return kit.Light{ID: id, Tone: "good", Icon: "github", Label: "GitHub", Href: "/settings/github", Open: open,
		Title:   "Connected as the org’s bot, read-only.",
		Details: []kit.Detail{{Label: "Account", Value: "elderlygate-ai"}, {Label: "Access", Value: "read-only"}, {Label: "Last checked", Value: "28 Oct 15:36"}, {Label: "Token from", Value: "the Keychain item control-room-github"}},
		Actions: []kit.Action{{Label: "Enter a new token…", Opens: "dlg-light-live"},
			{Label: "Disconnect", Variant: "danger", Confirm: &kit.Confirm{Text: "Removes the token from the Keychain. It stays valid on GitHub until you delete it there.", Label: "Disconnect now", Action: "#disconnect", Hidden: []kit.Hidden{{Name: "back", Value: "/"}}}}}}
}

// tokenDialog asks for a new token: the steps, the link, the secret field and Connect.
func tokenDialog(id string, open bool) kit.Dialog {
	body := template.HTML(`<p class="ui-dialog-text">Make a fine-grained, read-only token as elderlygate-ai: <a href="https://github.com/settings/personal-access-tokens/new" target="_blank" rel="noopener noreferrer">new token</a>.</p>`) +
		html("ui/field", kit.Field{ID: id + "-token", Name: "token", Label: "Token", Secret: true, Hint: "Checked with GitHub, then kept in the Keychain."})
	return kit.Dialog{ID: id, Title: "Enter a new token", Action: "#connect", Body: body, Confirm: kit.Button{Label: "Connect", Variant: "primary", Type: "submit"}, Open: open}
}

// SampleMilestones is the org's road to public.
func SampleMilestones() kit.Milestones {
	return kit.Milestones{
		Label: "Road to public", Total: "157 days",
		Track: []kit.Mark{{Kind: "today", At: 0}, {Kind: "span", At: .255, Width: .07}, {Kind: "point", At: .255}, {Kind: "point", At: .61}, {Kind: "point", At: 1}},
		Items: []kit.Milestone{
			{Name: "Rename window", When: "7–18 Dec 2026", Days: "40 days to start", Focus: true},
			{Name: "Private group", When: "1 Feb 2027", Days: "96 days"},
			{Name: "Public", When: "3 Apr 2027", Days: "157 days"},
		},
	}
}

// Parts lists the shell pieces and the 20 parts in the order the gallery shows them.
func Parts() []Part {
	return []Part{
		{"preview-banner", "Preview banner", ".ui-preview-banner", "Above the top bar of a build that is not a release from main: amber and black stripes no other part uses, the branch, the commit, and uncommitted changes when there are any.", []string{"The app decides Released from its own build info; the kit only draws it.", "A released build from main draws nothing. It cannot be dismissed.", "It stays on top while the page scrolls; on a phone a long line is cut, and the whole of it, with the version, is on hover."}},
		{"bar", "Top bar", ".ui-bar", "App name, up to 6 page links, global status, theme. On a phone the links scroll sideways in their own strip.", []string{"More than 6 pages of the same weight: add the left rail under the bar; the bar never goes away.", "A status here is a dot plus one word; the full story is a notice on the page that needs it.", "Apps set: the app name becomes a switcher menu (slot for later)."}},
		{"light", "Status light", ".ui-light · .ui-pop", "A traffic light in the top bar for one connection: a dot, an icon and a word. Clicking it opens a popover under it with the facts and an Actions menu.", []string{"Tones: good green, warn amber, bad red, none grey. The word stays the same; the popover says why.", "Without JavaScript it is a link to the app's full page for it, with the same facts and forms.", "The Actions menu is a menu button of plain buttons: one opens a dialog, one asks once more inside the popover before it acts.", "Escape, a click outside or Tab out closes the popover; focus goes back to the light.", "With Src it loads on its own: the bar draws it grey first, so it never slows the page."}},
		{"rail", "Left rail", ".ui-rail", "Only for apps with more than 6 pages, always under the top bar.", []string{"On a phone it becomes a strip that scrolls sideways."}},
		{"page-header", "Page header", ".ui-head", "Breadcrumb when the page is below another, title, one plain sentence, the page's actions.", []string{"On a phone the actions drop below the title, full width."}},
		{"footer", "Footer status bar", ".ui-footer", "Version, branch when not main, when the data was read and reads next, a live dot, links to the JSON, source and help.", []string{"Never fixed, so it never covers content.", "A page that stopped updating turns its dot red and says so in words."}},
		{"product-card", "Product card", ".ui-card", "The few facts about one product that tell you whether to look closer.", []string{"Facts always Installed, Board, Pull requests, My copy; at most two more.", "Every number links to its source; its hover shows the command.", "One marker on the right: the star, or the most urgent pill, or a count waiting for you."}},
		{"products", "Product card and product list", ".ui-cards · .ui-list", "One part, two views, with the toolbar above. Works without JavaScript: the view and filters are in the address.", []string{"Order is the server's order.", "Loading draws card outlines with grey bars, never a spinner over the grid."}},
		{"pr-item", "PR item", ".ui-pr", "One open pull request with everything needed to decide whether to open it.", []string{"Review pill: your review accent, approved good, changes asked warn, no review yet neutral, draft tag.", "CI pill: passing good, failing n/m bad (links to the run), running info, no checks neutral.", "The server sends every field already worked out; the UI does not call GitHub."}},
		{"pr-list", "PR list", ".ui-prs", "PR items in one box with their own loading and empty states.", nil},
		{"milestones", "Milestone strip", ".ui-miles", "Where we are on the road to a date.", []string{"Positions come from the server as fractions.", "On a phone the track hides and the list stacks."}},
		{"status", "Status dot", ".ui-dot", "The smallest status sign: a dot and a word.", []string{"Always with a word. Same five meanings as hubman's ● ↑ ✗ ◌ ○."}},
		{"pill", "Pill", ".ui-pill", "A short state that needs to stand out more than a dot. Also tags.", []string{"One to three words, lower case. At most two per row. A pill is not a button.", "Tags are neutral; a board may map a tag to one of the five tones."}},
		{"count", "Count badge", ".ui-count", "How many of something, next to the thing.", []string{"Zero hides the badge unless zero is news. 100 and more shows 99+, with the true number for screen readers."}},
		{"section", "Section", ".ui-section", "The block every page is made of. It reads its own data and fails on its own.", []string{"States: loaded, loading, empty, partial (content plus an error), failed (error box only).", "Lazy: with Src it draws its skeleton at once and kit.js fetches its HTML, every lazy section in parallel; a failed fetch shows its error box with Retry. Without JavaScript, Full links the page drawn whole. See /lazy."}},
		{"error", "Error box", ".ui-box[data-tone=bad]", "What went wrong, inside the part it affects.", []string{"What failed, why, and the command to check it. Never a raw server answer. No apology."}},
		{"notice", "Notice", ".ui-box[data-tone]", "News about the whole page that is not an error.", []string{"At most one notice under the page header."}},
		{"empty", "Empty state", ".ui-empty", "A part with nothing to show says so and gives the next step.", []string{"Headline is a fact. A filter that hides everything says “Nothing matches” and offers Clear filters."}},
		{"table", "Data list", ".ui-list", "Many items of one kind, compared column by column.", []string{"First column names the item and links to it. Numbers right-aligned in mono. Wide lists scroll inside their box."}},
		{"ticket", "Card / ticket", ".ui-ticket", "A piece of work on a board: kanban cards, decisions, design findings.", []string{"Dragging is the app's business, not the kit's."}},
		{"button", "Buttons", ".ui-btn", "Things that do something. Links go somewhere; buttons act.", []string{"Primary: one per page at most. Label is a verb plus the object. Busy says what it does.", "A button that changes something asks first (Confirm): the question, Cancel, and the button's verb again."}},
		{"field", "Form field, including secrets", ".ui-field", "Label above, input, hint or error below.", []string{"Secrets: password type, a Show toggle, no autocomplete; after saving only “Connected as …”."}},
		{"toolbar", "Toolbar with view toggle", ".ui-toolbar · .ui-seg", "Search, filters, clear · space · count · view toggle · one primary action.", []string{"An active filter gets an accent edge. The view lives in the address and in the browser."}},
		{"tabs", "Tabs", ".ui-tabs", "Views of the same thing on one page.", []string{"Tabs switch views, the top bar switches pages. Arrow keys move between tabs. Each tab has its own address."}},
		{"crumb", "Breadcrumb", ".ui-crumb", "Where this page sits, for pages two or more levels down.", []string{"The last item is the current page and is not a link."}},
		{"dialog", "Modal", "dialog.ui-dialog", "Asks before acting.", []string{"Native <dialog>: traps focus, closes on Escape. Cancel gets first focus. The confirm button repeats the verb and the count."}},
		{"drawer", "Drawer", "dialog.ui-drawer", "Detail next to the list you came from.", []string{"On a phone the drawer becomes a full page."}},
		{"star", "Focus marker", ".ui-star", "A yellow star: this is what we focus on now.", []string{"Only --ui-focus draws it, and it draws nothing else.", "If more than three things carry a star, nothing is in focus.", "Whether a thing is in focus is data (focus: true), never decided by the UI."}},
		{"skeleton", "Loading placeholder", ".ui-skel", "Grey bars in the shape of what is coming.", []string{"Shown only after 300ms. The shimmer stops under reduced motion. Never a full-page spinner."}},
	}
}

// Cases is every part in every state.
func Cases() []Case {
	prs := SamplePRs()
	compact := prs[0]
	compact.Compact, compact.Verb = true, "Approve"
	return []Case{
		{Part: "preview-banner", State: "branch", Items: one("ui/preview-banner", kit.Build{Version: "v0.4.0-3-g758ac48", Branch: "feat/road-focus", Commit: "758ac48"})},
		{Part: "preview-banner", State: "dirty", Note: "Built from a tree with changes that are not committed.", Items: one("ui/preview-banner", kit.Build{Version: "dev", Branch: "feat/preview-banner", Commit: "1c9e2f0", Dirty: true})},
		{Part: "preview-banner", State: "released", Note: "A release from main draws nothing.", Items: one("ui/preview-banner", kit.Build{Version: "v0.4.0", Branch: "main", Commit: "9a1b2c3", Released: true})},
		{Part: "bar", State: "default", Items: one("ui/bar", kit.Bar{App: "Control room", Sub: "the org", Href: "#", Pages: []kit.Link{{Label: "Focus", Href: "#", Current: true}, {Label: "Products", Href: "#"}, {Label: "Designs", Href: "#"}}, Status: []kit.Status{{Tone: "good", Label: "GitHub"}}, Refresh: &kit.Link{Href: "#"}})},
		{Part: "bar", State: "not-connected", Items: one("ui/bar", kit.Bar{App: "Control room", Href: "#", Pages: []kit.Link{{Label: "Focus", Href: "#"}, {Label: "Products", Href: "#", Current: true, Count: &kit.Count{Tone: "accent", N: 2}}, {Label: "Designs", Href: "#"}}, Status: []kit.Status{{Tone: "bad", Label: "GitHub not connected"}}})},
		{Part: "bar", State: "light", Note: "A status light left of the theme switch; it opens its popover.", Items: one("ui/bar", kit.Bar{App: "Control room", Sub: "the org", Href: "#", Pages: []kit.Link{{Label: "Focus", Href: "#", Current: true}, {Label: "Products", Href: "#"}}, Lights: []kit.Light{lightGood("light-bar", false)}, Refresh: &kit.Link{Href: "#"}})},
		{Part: "bar", State: "carried", Note: "Read again and the lazy light's back follow the address when a toolbar changes it without a reload.", Items: one("ui/bar", kit.Bar{App: "Control room", Sub: "the org", Href: "#", Pages: []kit.Link{{Label: "Focus", Href: "/", Current: true}}, Lights: []kit.Light{{ID: "light-carried", Icon: "github", Label: "GitHub", Href: "/settings/github", Src: "/part/light?back=%2F%3Fpr_sort%3Dnewest", Back: "back"}}, Refresh: &kit.Link{Label: "Read again", Href: "/?pr_sort=newest&fresh=1", Carry: "fresh"}, Read: &kit.ReadAt{Text: "read 10:42", Title: "The page was read at 29 Oct 2026, 10:42:13", At: "2026-10-29T10:42:13+01:00"}})},
		{Part: "bar", State: "app-switcher", Note: "The slot for later: it needs an API that lists the running apps.", Items: one("ui/bar", kit.Bar{App: "kanban", Sub: "hubman board", Apps: []kit.Link{{Label: "Control room", Href: "#"}, {Label: "kanban", Href: "#", Current: true}, {Label: "e-designer", Href: "#"}}, Pages: []kit.Link{{Label: "Board", Href: "#", Current: true}, {Label: "Archive", Href: "#"}}})},
		{Part: "light", State: "good", Items: one("ui/light", lightGood("light-good", false))},
		{Part: "light", State: "warn", Items: one("ui/light", kit.Light{ID: "light-warn", Tone: "warn", Icon: "github", Label: "GitHub", Href: "/settings/github", Title: "Connected as the org’s bot, read-only.", Text: "The last read failed: cannot reach GitHub.", Details: lightGood("", false).Details, Actions: lightGood("", false).Actions})},
		{Part: "light", State: "bad", Items: one("ui/light", kit.Light{ID: "light-bad", Tone: "bad", Icon: "github", Label: "GitHub", Href: "/settings/github", Title: "GitHub refuses the token.", Text: "401 Bad credentials: it has expired or been deleted.", Details: []kit.Detail{{Label: "Token from", Value: "the Keychain item control-room-github"}}, Actions: lightGood("", false).Actions})},
		{Part: "light", State: "none", Items: one("ui/light", kit.Light{ID: "light-none", Tone: "none", Icon: "github", Label: "GitHub", Href: "/settings/github", Title: "GitHub is not connected.", Text: "Pull requests and CI are not shown.", Actions: []kit.Action{{Label: "Enter a token…", Opens: "dlg-light-live"}}})},
		{Part: "light", State: "loading", Note: "Lazy: drawn grey at once, then kit.js fetches Src.", Items: one("ui/light", kit.Light{ID: "light-lazy", Icon: "github", Label: "GitHub", Href: "/settings/github", Src: "/part/light"})},
		{Part: "light", State: "open", Note: "The popover open: the facts, then the Actions menu.", Items: one("ui/light", lightGood("light-open", true))},
		{Part: "rail", State: "default", Items: one("ui/rail", []kit.RailGroup{{Label: "Favourites", Links: []kit.Link{{Label: "Deployments", Href: "#", Current: true}, {Label: "Machines", Href: "#", Count: &kit.Count{Tone: "warn", N: 2}}}}, {Label: "All services", Links: []kit.Link{{Label: "Providers", Href: "#"}, {Label: "Secrets", Href: "#"}, {Label: "Logs", Href: "#"}}}})},
		{Part: "page-header", State: "top-level", Items: one("ui/page-header", kit.Head{Title: "Products", Lede: "Is it current, is anything open, is my copy behind.", Actions: []kit.Button{{Label: "Refresh"}}})},
		{Part: "page-header", State: "below-another", Items: one("ui/page-header", kit.Head{Crumb: []kit.Link{{Label: "Products", Href: "#"}, {Label: "kanban", Href: "#"}, {Label: "20261027-a188"}}, Title: "Point every address at the org", Lede: "In Progress · turn: agent · 1.1d", Actions: []kit.Button{{Label: "Open board", Href: "#"}, {Label: "Archive", Variant: "danger", Opens: "dlg-archive"}}})},
		{Part: "footer", State: "live", Items: one("ui/footer", kit.Footer{Version: "control-room v0.4.0", VersionHref: "#release", Fresh: kit.Status{Tone: "good", Label: "read 16:21 from 9 boards · next in 30s"}, Links: []kit.Link{{Label: "JSON", Href: "#"}, {Label: "Source", Href: "#"}, {Label: "Help", Href: "#"}}})},
		{Part: "footer", State: "stopped", Items: one("ui/footer", kit.Footer{Version: "kanban v0.2.0", Branch: "fix/agent-test-git-identity", Fresh: kit.Status{Tone: "bad", Label: "stopped updating at 16:02 · Refresh to try again"}, Links: []kit.Link{{Label: "JSON", Href: "#"}, {Label: "Help", Href: "#"}}})},

		{Part: "product-card", State: "focus", Items: one("ui/product-card", hubman())},
		{Part: "product-card", State: "update-ready", Items: one("ui/product-card", kanbanProduct())},
		{Part: "product-card", State: "waiting-for-you", Items: one("ui/product-card", checkman())},
		{Part: "product-card", State: "partial", Items: one("ui/product-card", designer())},
		{Part: "product-card", State: "planned", Items: one("ui/product-card", planned())},
		{Part: "product-card", State: "library", Items: one("ui/product-card", library())},
		{Part: "products", State: "cards", Items: one("ui/products", kit.Products{ID: "products", View: "cards", Toolbar: productToolbar("cards", "", ""), Columns: ProductColumns, Items: SampleProducts()[:4]})},
		{Part: "products", State: "list", Items: one("ui/products", kit.Products{ID: "products-list", View: "list", Columns: ProductColumns, Items: SampleProducts()[:4]})},
		{Part: "products", State: "loading", Items: one("ui/products", kit.Products{ID: "products-loading", Loading: true})},
		{Part: "products", State: "empty", Items: one("ui/products", kit.Products{ID: "products-empty", Empty: &kit.Empty{Title: "No products in the catalog yet.", Text: "hubman's catalog lists every product of the org.", Command: "hubman catalog add <id>"}})},
		{Part: "products", State: "error", Items: one("ui/products", kit.Products{ID: "products-error", Problem: &kit.Problem{Title: "Can’t read the product catalog.", Cause: "hubman is not installed on this machine. Everything else on this page is current.", Command: "hubman version"}})},

		{Part: "pr-item", State: "your-review", Items: one("ui/pr-item", prs[0])},
		{Part: "pr-item", State: "conflicts-failing", Items: one("ui/pr-item", prs[1])},
		{Part: "pr-item", State: "draft-running", Items: one("ui/pr-item", prs[2])},
		{Part: "pr-item", State: "compact", Note: "Do now: verb first, line 1 and 2.", Items: one("ui/pr-item", compact)},
		{Part: "pr-item", State: "repo-failed", Items: one("ui/pr-item", kit.PR{Repo: "fleet", Problem: &kit.Problem{Title: "Can’t read pull requests of fleet.", Cause: "GitHub answered 404: the token cannot see this repo.", Command: "gh pr list -R <org>/fleet"}})},
		{Part: "pr-list", State: "loaded", Items: one("ui/pr-list", kit.PRList{Items: prs})},
		{Part: "pr-list", State: "loading", Items: one("ui/pr-list", kit.PRList{Loading: true})},
		{Part: "pr-list", State: "empty", Items: one("ui/pr-list", kit.PRList{Empty: &kit.Empty{Title: "No open pull requests.", Text: "Nothing waits for review in the 9 repos the org reads."}})},

		{Part: "milestones", State: "roadmap", Items: one("ui/milestones", SampleMilestones())},
		{Part: "milestones", State: "releases", Note: "One product's releases; a past one is a filled dot and drops from the list.", Items: one("ui/milestones", kit.Milestones{Label: "kanban releases", Total: "61 days", Track: []kit.Mark{{Kind: "past", At: 0}, {Kind: "today", At: .2}, {Kind: "point", At: .55}, {Kind: "point", At: 1}}, Items: []kit.Milestone{{Name: "v0.3.0 · tags map to tones", When: "20 Nov 2026", Days: "23 days"}, {Name: "v0.4.0 · the kit's shell", When: "28 Dec 2026", Days: "61 days"}}})},
		{Part: "milestones", State: "empty", Items: one("ui/milestones", kit.Milestones{Label: "Road to public", Empty: &kit.Empty{Title: "No milestones set.", Text: "Milestones are the releases of the org's family design.", Command: "e-designer get release"}})},

		{Part: "status", State: "tones", Items: []Item{{"ui/status", kit.Status{Tone: "good", Label: "current"}}, {"ui/status", kit.Status{Tone: "warn", Label: "3 behind"}}, {"ui/status", kit.Status{Tone: "bad", Label: "off main"}}, {"ui/status", kit.Status{Tone: "info", Label: "unmanaged"}}, {"ui/status", kit.Status{Tone: "none", Label: "planned"}}}},
		{Part: "pill", State: "tones", Items: []Item{{"ui/pill", kit.Pill{Tone: "good", Label: "CI passing"}}, {"ui/pill", kit.Pill{Tone: "warn", Label: "9 waiting"}}, {"ui/pill", kit.Pill{Tone: "bad", Label: "CI failing"}}, {"ui/pill", kit.Pill{Tone: "info", Label: "running"}}, {"ui/pill", kit.Pill{Tone: "accent", Label: "your review"}}, {"ui/pill", kit.Pill{Label: "no checks"}}}},
		{Part: "pill", State: "tags", Items: []Item{{"ui/pill", kit.Pill{Variant: "tag", Label: "#feature"}}, {"ui/pill", kit.Pill{Variant: "tag", Label: "#cli"}}, {"ui/pill", kit.Pill{Variant: "tag", Tone: "bad", Label: "#bug"}}, {"ui/pill", kit.Pill{Variant: "outline", Tone: "warn", Label: "v0.2.0 ready"}}}},
		{Part: "count", State: "badges", Note: "An icon badge per group: hover or focus it for its words; one with Href jumps to its group.", Items: []Item{
			{"ui/count", kit.Count{Icon: "eye", Tone: "accent", N: 2, Label: "2 waiting for your review: review them on GitHub", Href: "#review"}},
			{"ui/count", kit.Count{Icon: "check", Tone: "good", N: 6, Label: "6 approved, not merged yet: merge them on GitHub"}},
			{"ui/count", kit.Count{Icon: "alert", Tone: "warn", N: 1, Label: "1 needs work: CI failing or changes requested"}},
			{"ui/count", kit.Count{Icon: "edit", N: 1, Label: "1 draft: not ready for review"}}}},
		{Part: "count", State: "tones", Items: []Item{{"ui/count", kit.Count{N: 9}}, {"ui/count", kit.Count{Tone: "accent", N: 5}}, {"ui/count", kit.Count{Tone: "warn", N: 2}}, {"ui/count", kit.Count{Tone: "bad", N: 1}}}},
		{Part: "count", State: "zero-and-many", Note: "Zero hides unless it is news; 100 and more shows 99+.", Items: []Item{{"ui/count", kit.Count{N: 0}}, {"ui/count", kit.Count{N: 0, ShowZero: true, Label: "0 failing"}}, {"ui/count", kit.Count{Tone: "accent", N: 128}}}},

		{Part: "section", State: "loaded", Items: one("ui/section", kit.Section{ID: "sec-prs", Title: "Open pull requests", Count: &kit.Count{Tone: "accent", N: 3}, Why: "Yours first, then by age.", Body: html("ui/pr-list", kit.PRList{Items: prs[:2]}), More: &kit.Link{Label: "All 12 →", Href: "#"}})},
		{Part: "section", State: "loading", Items: one("ui/section", kit.Section{ID: "sec-loading", Title: "Open pull requests", Loading: true})},
		{Part: "section", State: "empty", Items: one("ui/section", kit.Section{ID: "sec-empty", Title: "Open pull requests", Empty: &kit.Empty{Title: "No open pull requests.", Text: "Nothing waits for review in the 9 repos the org reads."}})},
		{Part: "section", State: "partial", Items: one("ui/section", kit.Section{ID: "sec-partial", Title: "Open pull requests", Count: &kit.Count{Tone: "accent", N: 1}, Problem: &kit.Problem{Title: "Can’t read pull requests of fleet.", Cause: "The other 8 repos are current."}, Body: html("ui/pr-list", kit.PRList{Items: prs[:1]})})},
		{Part: "section", State: "lazy", Note: "Loads /part/slow here: the skeleton for a second and a half, then the section.", Items: one("ui/section", kit.Section{ID: "sec-lazy", Title: "Slow source", Why: "Loads on its own.", Src: "/part/slow", Full: "/lazy?full=1"})},
		{Part: "section", State: "lazy-failed", Note: "What a lazy section shows when its fetch fails; Retry fetches it again.", Items: one("ui/section", kit.Section{ID: "sec-lazy-failed", Title: "Failing source", Problem: &kit.Problem{Title: "Can’t load Failing source.", Cause: "The server answered 500 Internal Server Error. Everything else on this page is current.", Retry: "/part/fails"}})},
		{Part: "section", State: "read", Note: "When its sources were read, the exact time on hover, beside Read again.", Items: one("ui/section", kit.Section{ID: "sec-read", Title: "Open pull requests", Why: "Yours first, then by age.", Again: "/?fresh=1", Read: &kit.ReadAt{Text: "read 10:42", Title: "Its sources were read at 29 Oct 2026, 10:42:13", At: "2026-10-29T10:42:13+01:00"}, Body: html("ui/pr-list", kit.PRList{Items: prs[:1]})})},
		{Part: "section", State: "lazy-again", Note: "Read again reloads only this section, with fresh=1; without JavaScript it is a link to the page read again.", Items: one("ui/section", kit.Section{ID: "sec-lazy-again", Title: "Slow source", Why: "Loads on its own.", Src: "/part/slow", Full: "/lazy?full=1", Again: "/lazy?fresh=1"})},
		{Part: "section", State: "failed", Items: one("ui/section", kit.Section{ID: "sec-failed", Title: "Open pull requests", Why: "Yours first, then by age.", Problem: &kit.Problem{Title: "Can’t read pull requests from GitHub.", Cause: "The token was refused (401). Everything else on this page is current.", Command: "control-room github check"}})},

		{Part: "error", State: "default", Items: one("ui/error", kit.Problem{Title: "Can’t read the kanban board for designer.", Cause: "No manifest.json in .kanban/. Register the board or move it back.", Command: "kanban register --board <path>"})},
		{Part: "notice", State: "good", Items: one("ui/notice", kit.Notice{Tone: "good", Title: "GitHub: connected as the org’s bot, read-only."})},
		{Part: "notice", State: "info", Items: one("ui/notice", kit.Notice{Tone: "info", Title: "GitHub is not connected.", Text: "Pull requests and CI are not shown; everything else on this page is current.", Link: &kit.Link{Label: "Connect GitHub", Href: "#field"}})},
		{Part: "notice", State: "warn", Items: one("ui/notice", kit.Notice{Tone: "warn", Title: "Update ready: kanban v0.2.0.", Text: "6 checkouts are behind; pull before you start work."})},

		{Part: "empty", State: "default", Items: one("ui/empty", kit.Empty{Title: "No open pull requests.", Text: "Nothing waits for review in the 9 repos the org reads."})},
		{Part: "empty", State: "with-command", Items: one("ui/empty", kit.Empty{Title: "No designs for this product yet.", Text: "Put a board in the project to start one.", Command: "install.sh <project>"})},
		{Part: "empty", State: "nothing-matches", Items: one("ui/empty", kit.Empty{Title: "Nothing matches “fleet”.", Text: "6 products are hidden by the filter.", Action: &kit.Button{Label: "Clear filters", Variant: "quiet", Size: "sm", Href: "?"}})},
		{Part: "empty", State: "small", Note: "kanban's empty column.", Items: one("ui/empty", kit.Empty{Small: true})},

		{Part: "table", State: "loaded", Items: one("ui/table", kit.Table{Label: "Cards", Small: true, Columns: []kit.Column{{Label: "Card"}, {Label: "Column"}, {Label: "Turn"}, {Label: "Age", Num: true}, {Label: "Id"}}, Rows: []kit.Row{
			{Cells: []kit.Cell{{Text: "See and change checkman’s settings", Href: "#", Star: &kit.Star{}}, {Pill: &kit.Pill{Variant: "tag", Label: "Review"}}, {Text: "tugbay"}, {Text: "1.8d", Num: true}, {ID: "20261026-044d"}}},
			{Cells: []kit.Cell{{Text: "Point every address at the org", Href: "#"}, {Pill: &kit.Pill{Variant: "tag", Label: "In Progress"}}, {Text: "agent"}, {Text: "1.1d", Num: true}, {ID: "20261027-a188"}}},
		}})},
		{Part: "table", State: "error-row", Items: one("ui/table", kit.Table{Label: "Checkouts", Small: true, Columns: []kit.Column{{Label: "Repo"}, {Label: "Branch"}, {Label: "State"}, {Label: "Behind", Num: true}}, Rows: []kit.Row{
			{Cells: []kit.Cell{{Text: "hubman", Href: "#", Sub: "~/…/hubman"}, {Text: "main"}, {Status: &kit.Status{Tone: "warn", Label: "behind"}}, {Text: "3", Num: true}}},
			{Error: "Can’t read the checkout of designer: the folder is gone."},
		}})},
		{Part: "table", State: "loading", Items: one("ui/table", kit.Table{Label: "Cards", Loading: true})},
		{Part: "table", State: "empty", Items: one("ui/table", kit.Table{Label: "Cards", Small: true, Columns: []kit.Column{{Label: "Card"}, {Label: "Column"}, {Label: "Age", Num: true}}, Empty: &kit.Empty{Title: "No cards in Review.", Text: "Cards move here when their work is done."}})},

		{Part: "ticket", State: "focus", Items: one("ui/ticket", kit.Ticket{Title: "See and change checkman’s settings in the CLI and the dashboard", Href: "#", Summary: "settings: get/set config + Settings view", Tags: []kit.Pill{{Variant: "tag", Label: "#feature"}, {Variant: "tag", Label: "#cli"}, {Variant: "tag", Label: "#ui"}}, Turn: "tugbay", Age: "1.8d", ID: "20261026-044d", Focus: true})},
		{Part: "ticket", State: "plain", Items: one("ui/ticket", kit.Ticket{Title: "Decide the licence per product", Summary: "MIT for libraries, a choice for apps", Tags: []kit.Pill{{Variant: "tag", Label: "#decision"}}, Turn: "owner", Age: "4d", ID: "20261024-9b1e"})},

		{Part: "button", State: "variants", Items: []Item{{"ui/button", kit.Button{Label: "New card", Variant: "primary", Kbd: "n"}}, {"ui/button", kit.Button{Label: "Refresh"}}, {"ui/button", kit.Button{Label: "Show 6 more", Variant: "quiet"}}, {"ui/button", kit.Button{Label: "Disconnect GitHub", Variant: "danger", Opens: "dlg-live"}}}},
		{Part: "button", State: "asks-first", Note: "Confirm asks in the kit's dialog before the form is sent; Remember offers Remember my answer.", Items: []Item{{"ui/button", kit.Button{Label: "Save settings", Variant: "primary", Type: "submit", Confirm: "Save the settings?", Remember: true}}, {"ui/button", kit.Button{Label: "Disconnect", Variant: "danger", Type: "submit", Name: "do", Value: "disconnect", Confirm: "Disconnect GitHub?"}}}},
		{Part: "button", State: "asks-for-the-name", Note: "ConfirmName: the confirm button works only once the name is typed exactly.", Items: one("ui/button", kit.Button{Label: "Delete card", Variant: "danger", Type: "submit", Confirm: "Delete card 20261102-e472?", ConfirmName: "20261102-e472", Remember: true})},
		{Part: "button", State: "busy-disabled-small", Items: []Item{{"ui/button", kit.Button{Label: "Saving…", Busy: true}}, {"ui/button", kit.Button{Label: "Archive done (0)", Disabled: true}}, {"ui/button", kit.Button{Label: "Open on GitHub", Size: "sm", Href: "#"}}}},
		{Part: "field", State: "text", Items: one("ui/field", kit.Field{ID: "f-title", Label: "Card title", Value: "Answer help in the shared layout", Hint: "Shown on the card. Keep it under 70 characters."})},
		{Part: "field", State: "secret-invalid", Items: one("ui/field", kit.Field{ID: "f-token", Label: "GitHub token", Secret: true, Error: "GitHub refused this token: it cannot read pull requests. Give it “Pull requests: read”."})},
		{Part: "field", State: "secret-saved", Items: one("ui/field", kit.Field{ID: "f-saved", Label: "GitHub token", Saved: "Connected as the org’s bot", SavedAction: &kit.Button{Label: "Disconnect", Variant: "danger", Size: "sm", Opens: "dlg-live"}})},
		{Part: "toolbar", State: "sort-in-a-tab", Note: "Inside a lazy section it reloads only that section; Hidden keeps the tab.", Items: one("ui/toolbar", kit.Toolbar{Action: "/product/kanban", Hidden: []kit.Hidden{{Name: "tab", Value: "prs"}}, Filters: []kit.Select{{Name: "pr_sort", Label: "Sort", Active: true, Options: []kit.Option{{Value: "review", Label: "Waiting for review first"}, {Value: "newest", Label: "Newest first", Selected: true}}}}, Count: "3 open"})},
		{Part: "toolbar", State: "active-filter", Items: one("ui/toolbar", kit.Toolbar{Search: &kit.Search{Placeholder: "Search anything on a card"}, Filters: []kit.Select{{Name: "tag", Label: "Tag", Active: true, Options: []kit.Option{{Value: "feature", Label: "#feature", Selected: true}, {Value: "", Label: "All tags"}}}}, Clear: &kit.Link{Href: "?"}, Count: "4 cards", Target: "demo-view", View: "list", Primary: &kit.Button{Label: "New card", Variant: "primary", Kbd: "n"}})},
		{Part: "tabs", State: "links", Items: one("ui/tabs", kit.Tabs{Label: "Product views", Items: []kit.Tab{{Label: "Overview", Href: "#", Selected: true}, {Label: "Board", Href: "#", Count: &kit.Count{N: 7}}, {Label: "Pull requests", Href: "#", Count: &kit.Count{Tone: "accent", N: 1}}, {Label: "Releases", Href: "#"}, {Label: "Designs", Href: "#"}}})},
		{Part: "tabs", State: "carried", Note: "Each tab carries the page's other query values, such as the sort, and kit.js keeps them current.", Items: one("ui/tabs", kit.Tabs{Label: "Product views", Items: []kit.Tab{{Label: "Overview", Href: "/product/kanban?pr_sort=newest&tab=overview", Carry: "tab"}, {Label: "Pull requests", Href: "/product/kanban?pr_sort=newest&tab=prs", Carry: "tab", Selected: true}}})},
		{Part: "crumb", State: "default", Items: one("ui/crumb", []kit.Link{{Label: "Products", Href: "#"}, {Label: "kanban", Href: "#"}, {Label: "20261027-a188"}})},
		{Part: "dialog", State: "confirm-danger", Stage: true, Items: one("ui/dialog", kit.Dialog{ID: "dlg-archive-demo", Title: "Archive 4 done cards?", Text: "They move to archive/ and leave the board. You can bring them back with kanban restore.", Confirm: kit.Button{Label: "Archive 4 cards", Variant: "danger", Type: "submit", Value: "confirm"}, Open: true})},
		{Part: "dialog", State: "with-body", Stage: true, Note: "Body holds fields and steps; the form posts to Action.", Items: one("ui/dialog", tokenDialog("dlg-token-demo", true))},
		{Part: "drawer", State: "open", Stage: true, Items: one("ui/drawer", kit.Drawer{ID: "drw-demo", Title: "checkman #5", Open: true, Body: html("ui/pr-item", compact)})},
		{Part: "star", State: "marks", Items: []Item{{"ui/star", kit.Star{Text: "focus"}}, {"ui/star", kit.Star{Text: "focus now", Chip: true}}, {"ui/star", kit.Star{}}}},
		{Part: "skeleton", State: "card", Items: one("ui/skeleton", kit.Skeleton{Shape: "card", Lines: 3})},
		{Part: "skeleton", State: "rows", Items: one("ui/skeleton", kit.Skeleton{Shape: "rows", Lines: 4})},
	}
}
