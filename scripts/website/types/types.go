package types

import (
	"golang.org/x/tools/go/packages"
)

type CLIHelp struct {
	Enable            string `json:"enable"`
	RootCmdHelp       string `json:"rootHelp"`
	RunCmdHelp        string `json:"runHelp"`
	LintersCmdHelp    string `json:"lintersHelp"`
	FmtCmdHelp        string `json:"fmtHelp"`
	FormattersCmdHelp string `json:"formattersHelp"`
	HelpCmdHelp       string `json:"help"`
	MigrateCmdHelp    string `json:"migrateHelp"`
	ConfigCmdHelp     string `json:"configHelp"`
	CustomCmdHelp     string `json:"customHelp"`
	CacheCmdHelp      string `json:"cacheHelp"`
	VersionCmdHelp    string `json:"versionHelp"`
	CompletionCmdHelp string `json:"completionHelp"`
}

type ExcludeRule struct {
	Linters    []string `json:"linters,omitempty"`
	Path       string   `json:"path,omitempty"`
	PathExcept string   `json:"path-except,omitempty"`
	Text       string   `json:"text,omitempty"`
	Source     string   `json:"source,omitempty"`
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

	Groups map[string]struct{}

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

func (l *LinterWrapper) IsDeprecated() bool {
	return l.Deprecation != nil
}
