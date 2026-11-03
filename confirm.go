package kit

import (
	"net/http"
	"net/url"
	"slices"
)

// ActionTypes are the types every change belongs to, the same words in every product.
var ActionTypes = []string{"create", "edit", "move", "focus", "archive", "close", "delete", "settings", "connect", "run"}

// The ways a change can ask; the user picks one per action type.
const (
	AskName    = "name"    // type the name to go ahead
	AskYesNo   = "yes-no"  // Cancel or the verb
	AskNever   = "never"   // no question
	AskDefault = "default" // posted by "ui/confirm-settings": drop this type from the app's config, so its default holds
)

// Where the style in use comes from: this app's config or the default.
const (
	FromApp     = "this app"
	FromDefault = "default"
)

// ResolveConfirm is how a change asks and where that comes from: app, the style this app's config sets for its type, else def.
// def is AskYesNo for every change, delete included, and AskName for a change to the system itself; any other def counts as AskYesNo.
func ResolveConfirm(app, def string) (style, source string) {
	if isStyle(app) {
		return app, FromApp
	}
	if !isStyle(def) {
		def = AskYesNo
	}
	return def, FromDefault
}

func isStyle(s string) bool { return s == AskName || s == AskYesNo || s == AskNever }

// Confirmed reads the confirm field kit.js and ConfirmForm send with a form the user said yes to: ok for "yes" or "remember", remember for "remember".
// A server refuses a change that asks (its style is not AskNever) when ok is false, with RefuseUnconfirmed: the web twin of "refused without --yes".
func Confirmed(r *http.Request) (ok, remember bool) {
	switch r.FormValue("confirm") {
	case "yes":
		return true, false
	case "remember":
		return true, true
	}
	return false, false
}

// StatusUnconfirmed is the answer to a change that asks but came without a yes: 428 Precondition Required.
const StatusUnconfirmed = http.StatusPreconditionRequired

// RefuseUnconfirmed answers a change that came without a yes, such as a form sent without JavaScript, and makes nothing.
func RefuseUnconfirmed(w http.ResponseWriter) {
	http.Error(w, "Nothing was changed: this change asks first, and asking needs JavaScript. Turn it on and try again.", StatusUnconfirmed)
}

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

// ConfirmSettings is the Settings form for how each action type asks ("ui/confirm-settings"); it posts confirm.<type>=<style or default> to Action.
type ConfirmSettings struct {
	Action string
	Types  []string          // the product's action types; empty means all of ActionTypes
	App    map[string]string // this app's values; a type not set asks y/N
	System []string          // the types with a verb that changes the system itself: a change of those asks the name unless this app sets yes-no or name
	Hidden []Hidden
	Save   *Button // nil draws a primary Save
}

// SystemNote follows the words of a System type whose changes to the system still ask the name: by default, and when this app says never.
const SystemNote = " · name for changes to the system"

// ConfirmRow is one action type of the Settings form: the choice checked, and the style in use with where it comes from.
type ConfirmRow struct {
	Type   string
	Value  string // a style, or AskDefault when this app does not set it
	Style  string
	Now    string // Style in words
	Source string // FromApp or FromDefault
}

// Rows are the form's rows, one per action type.
func (s ConfirmSettings) Rows() []ConfirmRow {
	types := typesOr(s.Types)
	words := map[string]string{AskName: "Type the name", AskYesNo: "Ask y/N", AskNever: "Don’t ask"}
	rows := make([]ConfirmRow, len(types))
	for i, t := range types {
		style, source := ResolveConfirm(s.App[t], AskYesNo)
		value := s.App[t]
		if !isStyle(value) {
			value = AskDefault
		}
		now := words[style]
		if slices.Contains(s.System, t) && (source == FromDefault || style == AskNever) {
			now += SystemNote
		}
		rows[i] = ConfirmRow{Type: t, Value: value, Style: style, Now: now, Source: source}
	}
	return rows
}

// ConfirmChoices reads a posted "ui/confirm-settings" form drawn with types (the form's Types; empty means all of ActionTypes): per type, a style to keep in this app's config or AskDefault to drop.
func ConfirmChoices(form url.Values, types []string) map[string]string {
	chosen := map[string]string{}
	for _, t := range typesOr(types) {
		if s := form.Get("confirm." + t); isStyle(s) || s == AskDefault {
			chosen[t] = s
		}
	}
	return chosen
}

// typesOr is types, or all of ActionTypes when it is empty.
func typesOr(types []string) []string {
	if len(types) == 0 {
		return ActionTypes
	}
	return types
}
