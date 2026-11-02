package kit

import "net/url"

// ActionTypes are the types every change belongs to, the same words in every product.
var ActionTypes = []string{"create", "edit", "move", "focus", "archive", "close", "delete", "settings", "connect", "run"}

// The ways a change can ask; the user picks one per action type, for every app or only in one.
const (
	AskName    = "name"    // type the name to go ahead
	AskYesNo   = "yes-no"  // Cancel or the verb
	AskNever   = "never"   // no question
	AskInherit = "inherit" // posted by "ui/confirm-settings": drop this type from the scope, so the layer below holds
)

// Where the style in use comes from: this app's config, the global config every app shares, or the kit's default.
const (
	FromApp     = "this app"
	FromGlobal  = "global"
	FromDefault = "default"
)

// The scopes "ui/confirm-settings" edits, posted as confirm.scope.
const (
	ScopeApp    = "app"
	ScopeGlobal = "global"
)

// ResolveConfirm is how a change of actionType asks and where that comes from: this app's override, else the global value, else the default (delete the name, the rest yes-no).
func ResolveConfirm(app, global map[string]string, actionType string) (style, source string) {
	if s := app[actionType]; isStyle(s) {
		return s, FromApp
	}
	if s := global[actionType]; isStyle(s) {
		return s, FromGlobal
	}
	if actionType == "delete" {
		return AskName, FromDefault
	}
	return AskYesNo, FromDefault
}

func isStyle(s string) bool { return s == AskName || s == AskYesNo || s == AskNever }

// Ask gives a submit button the question its style asks, with Remember my answer; AskName also wants name typed, AskNever asks nothing.
func (b Button) Ask(style, question, name string) Button {
	b.Confirm, b.ConfirmName, b.Remember = "", "", false
	if style == AskNever {
		return b
	}
	b.Confirm, b.Remember = question, true
	if style == AskName {
		b.ConfirmName = name
	}
	return b
}

// ConfirmSettings is the Settings form for how each action type asks, in one scope ("ui/confirm-settings"); it posts confirm.scope and confirm.<type>=<style or inherit> to Action.
type ConfirmSettings struct {
	Action string
	Scope  string            // ScopeApp (the default): this app's overrides; ScopeGlobal: the values every app shares
	Types  []string          // the product's action types; empty means all of ActionTypes
	App    map[string]string // this app's overrides; a type not set inherits
	Global map[string]string // the global values; a type not set takes its default
	Hidden []Hidden
	Save   *Button // nil draws a primary Save
}

// ConfirmRow is one action type of the Settings form: the choice checked in its scope, and the style in use with where it comes from.
type ConfirmRow struct {
	Type   string
	Value  string // a style, or AskInherit when the scope does not set it
	Style  string
	Now    string // Style in words
	Source string // FromApp, FromGlobal or FromDefault
}

// Rows are the form's rows, one per action type.
func (s ConfirmSettings) Rows() []ConfirmRow {
	types := s.Types
	if len(types) == 0 {
		types = ActionTypes
	}
	own, app := s.App, s.App
	if s.Scope == ScopeGlobal {
		own, app = s.Global, nil
	}
	words := map[string]string{AskName: "Type the name", AskYesNo: "Ask y/N", AskNever: "Don’t ask"}
	rows := make([]ConfirmRow, len(types))
	for i, t := range types {
		style, source := ResolveConfirm(app, s.Global, t)
		value := own[t]
		if !isStyle(value) {
			value = AskInherit
		}
		rows[i] = ConfirmRow{Type: t, Value: value, Style: style, Now: words[style], Source: source}
	}
	return rows
}

// ConfirmChoices reads a posted "ui/confirm-settings" form: the scope it edits and, per type, a style to keep or AskInherit to drop.
func ConfirmChoices(form url.Values) (scope string, chosen map[string]string) {
	scope = ScopeApp
	if form.Get("confirm.scope") == ScopeGlobal {
		scope = ScopeGlobal
	}
	chosen = map[string]string{}
	for _, t := range ActionTypes {
		if s := form.Get("confirm." + t); isStyle(s) || s == AskInherit {
			chosen[t] = s
		}
	}
	return scope, chosen
}
