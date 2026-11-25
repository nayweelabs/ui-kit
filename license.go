package kit

// License is the License and features section of a Settings page ("ui/license"): the data of `<app> license show` as it is, each field's meaning in the contracts schema (mankit/v1 license-show), drawn as given.
type License struct {
	State    string           `json:"state"`              // state: none, trial, trial-ended, licensed, ending or ended
	EndedOn  string           `json:"endedOn,omitempty"`  // endedOn, the date of the ended headline
	Warn     bool             `json:"warn,omitempty"`     // warn: the badge and the card take the warn colour
	Sentence string           `json:"sentence,omitempty"` // sentence, shown word for word as the headline and in the panel
	Build    LicenseBuild     `json:"build"`              // build; a build that is not a release also draws the dev look
	Licenses []LicenseFile    `json:"licenses"`           // licenses
	Refused  []RefusedFile    `json:"refused"`            // refused
	Features []LicenseFeature `json:"features"`           // features
	Counts   LicenseCounts    `json:"counts"`             // counts
	Version  string           `json:"-"`                  // the version of the running build, as the page already has it
	Trial    string           `json:"-"`                  // where to ask for a trial key, drawn as a button; the app sets it only for none and trial-ended
	Actions  []Button         `json:"-"`                  // the app's own buttons, such as Add a license
}

// LicenseLight is the license badge of the top bar ("ui/license-light"): License's word in its colour, linking to Href, the app's Settings page; kit.js opens the License panel under it instead. Open draws the panel open.
type LicenseLight struct {
	Href    string
	License License
	Open    bool
}

// LicenseBuild is license show build.
type LicenseBuild struct {
	Release bool `json:"release"`
}

// LicenseFile is one of license show licenses.
type LicenseFile struct {
	Kind      string  `json:"kind"`
	For       string  `json:"for"`
	Until     string  `json:"until"`
	TrialEnds string  `json:"trialEnds,omitempty"` // trialEnds, shown in place of Until on a started release-trial
	OfferEnds string  `json:"offerEnds,omitempty"` // offerEnds, shown in place of Until on a release-trial not started
	DaysLeft  int     `json:"daysLeft"`
	Ended     bool    `json:"ended"`
	KeyID     string  `json:"keyId"`
	Window    *Window `json:"window,omitempty"` // with it, a countdown bar of DaysLeft
}

// RefusedFile is one of license show refused: a file the app read but did not take.
type RefusedFile struct {
	File   string `json:"file"`
	Type   string `json:"type"`
	Reason string `json:"reason"`
}

// LicenseCounts is license show counts.
type LicenseCounts struct {
	Free       int `json:"free"`
	Enterprise int `json:"enterprise"`
	On         int `json:"on"`
	Locked     int `json:"locked"`
}

// Window is license show licenses[].window: the countdown's length, which the bar draws DaysLeft against.
type Window struct {
	Of int `json:"of"`
}

// LicenseFeature is one of license show features: a featureState with from and until.
type LicenseFeature struct {
	ID        string         `json:"id"`
	Title     string         `json:"title"`
	Covers    string         `json:"covers"`
	Needs     []FeatureValue `json:"needs,omitempty"`
	Tier      string         `json:"tier"`
	On        bool           `json:"on"`
	By        string         `json:"by,omitempty"`
	StartsOff bool           `json:"startsOff,omitempty"`
	Missing   []FeatureValue `json:"missing,omitempty"`
	From      string         `json:"from"`
	Until     string         `json:"until,omitempty"`
}

// FeatureValue is one value a feature needs: featureValue.
type FeatureValue struct {
	Name    string `json:"name"`
	Purpose string `json:"purpose"`
}

// TrialButton is the button that says how to get a trial key, or nil without Trial.
func (l License) TrialButton() *Button {
	if l.Trial == "" {
		return nil
	}
	return &Button{Label: "How to get a trial key", Variant: "primary", Href: l.Trial}
}
