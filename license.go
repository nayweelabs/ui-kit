package kit

// License is the License and features section of a Settings page ("ui/license"): `<app> license show` and the feature list, each field's meaning in the contracts schema (mankit/v1), drawn as given.
type License struct {
	State     string           // licenseShow state: none, trial, trial-ended, licensed, ending or dev
	Release   bool             // licenseShow build.release
	Version   string           // the version of the running build, as the page already has it
	Licenses  []LicenseFile    // licenseShow licenses, in its order
	Features  []LicenseFeature // in the app's order, enterprise features first
	Countdown *Countdown       // nil draws no bar
	Trial     string           // where to ask for a trial key; with it, no license or an ended trial draws a button that goes there
	Actions   []Button         // the app's own buttons, such as Add a license
}

// LicenseFile is one license file the app read: licenseShow licenses[].
type LicenseFile struct {
	Kind  string // org, product or trial
	For   string
	Until string
	Ended bool
	KeyID string
}

// LicenseFeature is one row of the features table: a featureState of the feature list with from and until of licenseShow features[].
type LicenseFeature struct {
	ID        string
	Covers    string
	Tier      string // free or enterprise; empty is free
	On        bool
	By        string
	StartsOff bool
	Missing   []FeatureValue
	From      string // org, product, both, trial or none
	Until     string
}

// FeatureValue is one value a feature needs: featureValue.
type FeatureValue struct {
	Name    string
	Purpose string
}

// Countdown is the bar under the license's headline: Left of Of days, as the app counts them.
type Countdown struct {
	Left int
	Of   int
}
