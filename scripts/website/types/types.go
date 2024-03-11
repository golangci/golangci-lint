package types

import (
	"golang.org/x/tools/go/packages"
)

type CLIHelp struct {
	Enable  string `json:"enable"`
	Disable string `json:"disable"`
	Help    string `json:"help"`
}

type ExcludePattern struct {
	ID      string `json:"id,omitempty"`
	Pattern string `json:"pattern,omitempty"`
	Linter  string `json:"linter,omitempty"`
	Why     string `json:"why,omitempty"`
}

type Deprecation struct {
	Since       string `json:"since,omitempty"`
	Message     string `json:"message,omitempty"`
	Replacement string `json:"replacement,omitempty"`
}

// LinterWrapper same fields but with struct tags.
// The field Name and Desc are added to have the information about the linter.
// The field Linter is removed (not serializable).
type LinterWrapper struct {
	Name string `json:"name"` // From linter.
	Desc string `json:"desc"` // From linter.

	EnabledByDefault bool `json:"enabledByDefault,omitempty"`

	LoadMode packages.LoadMode `json:"loadMode,omitempty"`

	InPresets        []string `json:"inPresets,omitempty"`
	AlternativeNames []string `json:"alternativeNames,omitempty"`

	OriginalURL     string `json:"originalURL,omitempty"`
	Internal        bool   `json:"internal"`
	CanAutoFix      bool   `json:"canAutoFix,omitempty"`
	IsSlow          bool   `json:"isSlow"`
	DoesChangeTypes bool   `json:"doesChangeTypes,omitempty"`

	Since       string       `json:"since,omitempty"`
	Deprecation *Deprecation `json:"deprecation,omitempty"`
}
