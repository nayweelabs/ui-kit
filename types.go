package kit

import "html/template"

// The data each partial draws. The server works out every value (G-STACK-7): tone, label, order, age text. A Tone is one of "good", "warn", "bad", "info", "accent", "none", or "" for neutral.

// Link is a place to go. Current marks the page you are on.
//
// Carry makes it a link that carries the page's address: Href holds the page's query, and Carry names
// the link's own query values, space-separated (like "fresh" or "tab"). When the address changes without
// a reload (a toolbar or the view toggle), kit.js gives the link the new query and keeps its own values.
type Link struct {
	Label   string
	Href    string
	Title   string
	Current bool
	Count   *Count
	Carry   string
}

// Status is a dot and a word ("ui/status").
type Status struct {
	Tone  string
	Label string
}

// Pill is a short state or a tag ("ui/pill"). Variant is "", "outline" or "tag".
type Pill struct {
	Tone    string
	Variant string
	Label   string
	Title   string
	Href    string
}

// Count is how many of something ("ui/count"). Zero hides it unless ShowZero; over 99 shows "99+".
// Label says what it counts, to screen readers and on hover. With Icon (a "ui/icon" name) it is a badge
// whose Label shows as a tooltip on hover and on keyboard focus; with Href it is a link, such as to a
// group further down the page.
type Count struct {
	Tone     string
	N        int
	Label    string
	ShowZero bool
	Icon     string
	Href     string
}

// Star is the focus marker ("ui/star"). Label is read by screen readers; Text is shown beside the star.
type Star struct {
	Label  string
	Text   string
	Chip   bool
	Corner bool
}

// Num is a number that links to its source; hovering it shows the command that gives the same number.
type Num struct {
	Text    string
	Href    string
	Command string
}

// Problem is an error box ("ui/error"): what could not be done, why, and the command to check it.
type Problem struct {
	Title   string
	Cause   string
	Command string
	Retry   string // with it, a Retry link to this address; a lazy section's error box reloads the section
}

// Notice is news about the whole page ("ui/notice"). Tone is "good", "info" or "warn".
type Notice struct {
	Tone  string
	Title string
	Text  string
	Link  *Link
}

// Empty is the empty state ("ui/empty"). Small is kanban's one muted word in a column.
type Empty struct {
	Title   string
	Text    string
	Command string
	Action  *Button
	Small   bool
}

// Skeleton is a loading placeholder ("ui/skeleton"). Shape is "card" or "rows".
type Skeleton struct {
	Shape string
	Lines int
}

// Button acts ("ui/button"). With Href it is drawn as a link; with Opens it opens the dialog of that id.
type Button struct {
	Label    string
	Variant  string
	Size     string
	Type     string
	Href     string
	Kbd      string
	Name     string
	Value    string
	Opens    string
	Disabled bool
	Busy     bool
	Carry    string // with Href, the link's own query values; it carries the page's others (see Link)
}

// Fact is one label and value row of a product card, and one cell of its list row ("ui/fact").
type Fact struct {
	Label  string
	Tone   string
	Num    *Num
	Text   string
	Branch string
	Pill   *Pill
	Error  string
}

// Product is one product card or list row ("ui/product-card", "ui/product-row"). Variant is "", "planned" or "library". Marker is the one most urgent pill, Waiting a count waiting for you.
type Product struct {
	Name     string
	Href     string
	Kind     string
	Variant  string
	Focus    bool
	Marker   *Pill
	Waiting  *Count
	Problem  *Problem
	Facts    []Fact
	Foot     string
	OpenHref string
}

// Products is a toolbar with cards or a list of the same products ("ui/products"). View is "cards" or "list", from ?view=; the other view is drawn hidden so kit.js can switch without a reload.
type Products struct {
	ID      string
	View    string
	Toolbar *Toolbar
	Columns []string
	Items   []Product
	Empty   *Empty
	Problem *Problem
	Loading bool
}

// PR is one open pull request ("ui/pr-item"). Compact is Do now's form: a verb, line 1 and 2. A repo that failed is a PR with only Repo and Problem set.
type PR struct {
	Repo     string
	Number   string
	Href     string
	RepoHref string
	Author   string
	Initials string
	Age      string
	Draft    bool
	Title    string
	Verb     string
	Branch   string
	Base     string
	Files    *Num
	Added    string
	Removed  string
	Conflict string
	Review   *Pill
	CI       *Pill
	Open     string
	Compact  bool
	Problem  *Problem
}

// PRList is a list of PR items with its own empty and loading states ("ui/pr-list").
type PRList struct {
	Items   []PR
	Empty   *Empty
	Loading bool
}

// Mark is one point or window on the milestone track. At and Width are fractions 0..1 sent by the server. Kind is "point", "past", "today" or "span".
type Mark struct {
	Kind  string
	At    float64
	Width float64
}

// Milestone is one column under the track.
type Milestone struct {
	Name  string
	When  string
	Days  string
	Focus bool
}

// Milestones is the milestone strip ("ui/milestones").
type Milestones struct {
	Label string
	Total string
	Track []Mark
	Items []Milestone
	Empty *Empty
}

// Section is the block every page is made of ("ui/section"). It fails on its own: Loading draws a skeleton, Problem without Body is failed, Problem with Body is partial, Empty is empty.
//
// With Src it is lazy: it draws its skeleton at once and kit.js fetches Src, the section's own HTML, and puts it in its place; each lazy section on a page loads in parallel and on its own. If the fetch fails, the section shows its error box with a Retry link. A GET form inside a loaded lazy section (a toolbar) reloads only that section, and the address keeps the form's values. Full is the address of the page drawn whole on the server, linked for a browser without JavaScript.
//
// With Again, the head holds a Read again button, a link to Again. On a lazy section kit.js reloads only that section instead, from Src with fresh=1 added, so the server reads the section's sources again; the button stays in the loaded section's head.
type Section struct {
	ID      string
	Title   string
	Count   *Count
	Why     string
	Body    template.HTML
	Problem *Problem
	Empty   *Empty
	Loading bool
	More    *Link
	Src     string
	Full    string
	Again   string // the page with this section's sources read again (fresh=1): a Read again button in the head
	Read    *ReadAt
}

// ReadAt is when a section's or a page's sources were last read ("ui/read-at"): Text is shown ("read 10:42"),
// Title on hover (the exact time), At is the time as RFC 3339. A lazy section draws it empty; kit.js fills it from
// the UI-Read-At header of the section's HTML.
type ReadAt struct {
	Text  string
	Title string
	At    string
}

// Column is one column of a data list. Num right-aligns it in mono; Hidden keeps the label for screen readers only.
type Column struct {
	Label  string
	Num    bool
	Hidden bool
}

// Cell is one cell of a data list. Href makes Text the item's link; Sub is a second line.
type Cell struct {
	Text   string
	Href   string
	Sub    string
	Num    bool
	Status *Status
	Pill   *Pill
	ID     string
	Star   *Star
}

// Row is one row of a data list; Error replaces its cells with one error line.
type Row struct {
	Cells []Cell
	Error string
}

// Table is the data list ("ui/table").
type Table struct {
	Label   string
	Columns []Column
	Rows    []Row
	Empty   *Empty
	Loading bool
	Small   bool
}

// Ticket is a piece of work on a board ("ui/ticket"): a kanban card, a decision, a design finding.
type Ticket struct {
	Title   string
	Href    string
	Summary string
	Tags    []Pill
	Turn    string
	Age     string
	ID      string
	Focus   bool
}

// Field is a labelled input ("ui/field"). A Secret is a password input with Show; once saved, Saved replaces the input and the value is never shown back.
type Field struct {
	ID          string
	Name        string
	Label       string
	Type        string
	Value       string
	Placeholder string
	Hint        string
	Error       string
	Secret      bool
	Saved       string
	SavedAction *Button
}

// Option is one choice of a Select.
type Option struct {
	Value    string
	Label    string
	Selected bool
}

// Select is a filter in a toolbar. Active gives it the accent edge.
type Select struct {
	Name    string
	Label   string
	Options []Option
	Active  bool
}

// Search is the search box of a toolbar.
type Search struct {
	Name        string
	Label       string
	Placeholder string
	Value       string
}

// Toolbar sits above a list or grid ("ui/toolbar"). It is a GET form, so filters work without JavaScript. View is "cards" or "list" when the toolbar carries the view toggle for the view named Target.
type Toolbar struct {
	Action  string
	Search  *Search
	Filters []Select
	Clear   *Link
	Count   string
	Target  string
	View    string
	Primary *Button
	Hidden  []Hidden
}

// Hidden is a value a form sends unchanged, such as the tab it sits on. Back makes Value the page's
// address, which kit.js keeps current when the address changes without a reload.
type Hidden struct {
	Name  string
	Value string
	Back  bool
}

// Tab is one tab. With Href it is a link with its own address; without, it switches the panel of id Panel.
type Tab struct {
	Label    string
	Href     string
	Panel    string
	Count    *Count
	Selected bool
	Carry    string // with Href, the tab's own query values ("tab"); it carries the page's others (see Link)
}

// Tabs switch views of the same thing on one page ("ui/tabs").
type Tabs struct {
	Label string
	Items []Tab
}

// Dialog asks before acting ("ui/dialog"). Confirm repeats the verb and the count; Cancel gets first focus on danger.
type Dialog struct {
	ID      string
	Title   string
	Text    string
	Cancel  string
	Confirm Button
	Action  string
	Open    bool
	Body    template.HTML // between the text and the buttons, inside the form: fields, steps, hidden values
}

// Drawer shows detail next to the list you came from ("ui/drawer").
type Drawer struct {
	ID    string
	Title string
	Body  template.HTML
	Open  bool
}

// Bar is the top bar. Pages are up to 6 links; Apps, when set, turns the app name into a switcher menu. Lights sit left of Refresh and the theme switch.
type Bar struct {
	App     string
	Sub     string
	Href    string
	Pages   []Link
	Apps    []Link
	Status  []Status
	Lights  []Light
	Refresh *Link
	Read    *ReadAt // when the page was read, left of Refresh
	NoTheme bool
}

// Light is a status light in the top bar ("ui/light"): a traffic-light dot (Tone "good", "warn", "bad" or "none"), an icon and a short Label. It is a link to Href, the app's full page for it; kit.js makes it open a popover under it instead, holding Title, Text, the Details and the Actions menu. Open draws the popover open.
//
// With Src it is lazy: it is drawn neutral at once, and kit.js fetches Src, the light's own HTML, and puts it in its place.
type Light struct {
	ID      string
	Tone    string
	Icon    string
	Label   string
	Href    string
	Title   string
	Text    string
	Details []Detail
	Actions []Action
	Dialogs []Dialog // the dialogs its actions open
	Open    bool
	Src     string
	Back    string // the name of Src's query value that holds the page's address ("back"); kit.js keeps it current
	Retry   string // with it, the popover says why the light could not load (data-ui-cause) and a Retry link reloads it from this address
}

// Detail is one title and value line of a light's popover.
type Detail struct {
	Label string
	Value string
}

// Action is one choice of a light's Actions menu: a link (Href), a dialog it opens (Opens, a dialog's ID), or a form it posts after asking once more inside the popover (Confirm).
type Action struct {
	Label   string
	Variant string // "danger" for one that removes something
	Href    string
	Opens   string
	Confirm *Confirm
}

// Confirm is the step an action asks inside the popover: Text says what will happen, and the button Label posts the form to Action with Hidden.
type Confirm struct {
	Text   string
	Label  string
	Action string
	Hidden []Hidden
}

// Head is the page header: breadcrumb (pages two or more levels down), title, one line, actions.
type Head struct {
	Crumb   []Link
	Title   string
	Lede    string
	Actions []Button
}

// Build is what the running build was made from, as the app's server knows it (usually from -ldflags at build time).
// Released is the app's own answer to "is this a release from main"; when it is false, "ui/preview-banner" says so above the top bar.
// Commit is drawn as sent, so send it short.
type Build struct {
	Version  string
	Branch   string
	Commit   string
	Dirty    bool
	Released bool
}

// Footer is the status bar: version, branch when not main, freshness with a live dot, links.
type Footer struct {
	Version     string
	VersionHref string
	Branch      string
	Fresh       Status
	Links       []Link
}

// RailGroup is a labelled group of rail links.
type RailGroup struct {
	Label string
	Links []Link
}

// Page is a whole page in the shell ("ui/shell"). Rail is only for apps with more than 6 pages. Assets is where the kit's files are served, default "/ui-kit/".
type Page struct {
	Lang   string
	Title  string
	Assets string
	Head   template.HTML
	Wide   bool
	Bar    Bar
	Rail   []RailGroup
	Header Head
	Notice *Notice
	Body   template.HTML
	Footer Footer
	Build  *Build // nil draws no preview banner; a build that is not Released draws it above the top bar
	Full   string // the address of this page drawn whole on the server: a browser without JavaScript goes there, for pages with lazy sections
}
