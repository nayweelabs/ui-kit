package kit

import "net/url"

// ActionTypes are the types every change belongs to, the same words in every product.
var ActionTypes = []string{"create", "edit", "move", "focus", "archive", "close", "delete", "settings", "connect", "run"}

// The ways a change can ask; the user picks one per action type.
const (
	AskName  = "name"   // type the name to go ahead
	AskYesNo = "yes-no" // Cancel or the verb
	AskNever = "never"  // no question
)

// ConfirmStyle is how a change of actionType asks, from the user's choices; a type not chosen asks its default, delete the name and the rest yes-no.
func ConfirmStyle(chosen map[string]string, actionType string) string {
	switch s := chosen[actionType]; s {
	case AskName, AskYesNo, AskNever:
		return s
	}
	if actionType == "delete" {
		return AskName
	}
	return AskYesNo
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

// ConfirmSettings is the Settings form where the user picks how each action type asks ("ui/confirm-settings"); it posts confirm.<type>=<style> to Action.
type ConfirmSettings struct {
	Action string
	Types  []string          // the product's action types; empty means all of ActionTypes
	Chosen map[string]string // the user's choices; a type not chosen shows its default
	Hidden []Hidden
	Save   *Button // nil draws a primary Save
}

// ConfirmRow is one action type of the Settings form and the style it asks with now.
type ConfirmRow struct {
	Type  string
	Style string
}

// Rows are the form's rows, one per action type, each with its style.
func (s ConfirmSettings) Rows() []ConfirmRow {
	types := s.Types
	if len(types) == 0 {
		types = ActionTypes
	}
	rows := make([]ConfirmRow, len(types))
	for i, t := range types {
		rows[i] = ConfirmRow{Type: t, Style: ConfirmStyle(s.Chosen, t)}
	}
	return rows
}

// ConfirmChoices reads the choices a posted "ui/confirm-settings" form sent, ready to keep in the product's config.
func ConfirmChoices(form url.Values) map[string]string {
	chosen := map[string]string{}
	for _, t := range ActionTypes {
		switch s := form.Get("confirm." + t); s {
		case AskName, AskYesNo, AskNever:
			chosen[t] = s
		}
	}
	return chosen
}
