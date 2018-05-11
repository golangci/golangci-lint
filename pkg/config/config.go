package config

import (
	"time"
)

type OutFormat string

const (
	OutFormatJSON              = "json"
	OutFormatLineNumber        = "line-number"
	OutFormatColoredLineNumber = "colored-line-number"
)

var OutFormats = []string{OutFormatColoredLineNumber, OutFormatLineNumber, OutFormatJSON}

var DefaultExcludePatterns = []string{
	// errcheck
	"Error return value of .((os\\.)?std(out|err)\\..*|.*Close|os\\.Remove(All)?|.*printf?|os\\.(Un)?Setenv). is not checked",

	// golint
	"should have comment",
	"comment on exported method",

	// gas
	"G103:", // Use of unsafe calls should be audited
	"G104:", // disable what errcheck does: it reports on Close etc
	"G204:", // Subprocess launching should be audited: too lot false positives
	"G301:", // Expect directory permissions to be 0750 or less
	"G302:", // Expect file permissions to be 0600 or less
	"G304:", // Potential file inclusion via variable: `src, err := ioutil.ReadFile(filename)`

	// govet
	"possible misuse of unsafe.Pointer",
	"should have signature",

	// megacheck
	"ineffective break statement. Did you mean to break out of the outer loop", // developers tend to write in C-style with break in switch
}

type Run struct {
	IsVerbose      bool `mapstructure:"verbose"`
	CPUProfilePath string
	Concurrency    int

	Config string

	Args []string

	BuildTags []string `mapstructure:"build-tags"`

	ExitCodeIfIssuesFound int  `mapstructure:"issues-exit-code"`
	AnalyzeTests          bool `mapstructure:"tests"`
	Deadline              time.Duration
}

type LintersSettings struct {
	Errcheck struct {
		CheckTypeAssertions bool `mapstructure:"check-type-assertions"`
		CheckAssignToBlank  bool `mapstructure:"check-blank"`
	}
	Govet struct {
		CheckShadowing bool `mapstructure:"check-shadowing"`
	}
	Golint struct {
		MinConfidence float64 `mapstructure:"min-confidence"`
	}
	Gofmt struct {
		Simplify bool
	}
	Gocyclo struct {
		MinComplexity int `mapstructure:"min-complexity"`
	}
	Varcheck struct {
		CheckExportedFields bool `mapstructure:"exported-fields"`
	}
	Structcheck struct {
		CheckExportedFields bool `mapstructure:"exported-fields"`
	}
	Maligned struct {
		SuggestNewOrder bool `mapstructure:"suggest-new"`
	}
	Dupl struct {
		Threshold int
	}
	Goconst struct {
		MinStringLen        int `mapstructure:"min-len"`
		MinOccurrencesCount int `mapstructure:"min-occurrences"`
	}
}

type Linters struct {
	Enable     []string
	Disable    []string
	EnableAll  bool `mapstructure:"enable-all"`
	DisableAll bool `mapstructure:"disable-all"`

	Presets []string
}

type Issues struct {
	ExcludePatterns    []string `mapstructure:"exclude"`
	UseDefaultExcludes bool     `mapstructure:"exclude-use-default"`

	MaxIssuesPerLinter int `mapstructure:"max-issues-per-linter"`
	MaxSameIssues      int `mapstructure:"max-same-issues"`

	DiffFromRevision  string `mapstructure:"new-from-rev"`
	DiffPatchFilePath string `mapstructure:"new-from-patch"`
	Diff              bool   `mapstructure:"new"`
}

type Config struct { // nolint:maligned
	Run Run

	Output struct {
		Format              string
		PrintIssuedLine     bool `mapstructure:"print-issued-lines"`
		PrintLinterName     bool `mapstructure:"print-linter-name"`
		PrintWelcomeMessage bool `mapstructure:"print-welcome"`
	}

	LintersSettings LintersSettings `mapstructure:"linters-settings"`
	Linters         Linters
	Issues          Issues
}
