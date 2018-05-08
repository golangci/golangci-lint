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
	"should have comment",
	"comment on exported method",
	"G104", // disable what errcheck does: it reports on Close etc
	"G204", // Subprocess launching should be audited: too lot false positives
	"G304", // Potential file inclusion via variable: `src, err := ioutil.ReadFile(filename)`
}

type Common struct {
	IsVerbose      bool
	CPUProfilePath string
	Concurrency    int
}

type Run struct { // nolint:maligned
	Args []string

	BuildTags []string

	OutFormat       string
	PrintIssuedLine bool
	PrintLinterName bool

	ExitCodeIfIssuesFound int

	Errcheck struct {
		CheckClose          bool
		CheckTypeAssertions bool
		CheckAssignToBlank  bool
	}
	Govet struct {
		CheckShadowing bool
	}
	Golint struct {
		MinConfidence float64
	}
	Gofmt struct {
		Simplify bool
	}
	Gocyclo struct {
		MinComplexity int
	}
	Varcheck struct {
		CheckExportedFields bool
	}
	Structcheck struct {
		CheckExportedFields bool
	}
	Maligned struct {
		SuggestNewOrder bool
	}
	Megacheck struct {
		EnableStaticcheck bool
		EnableUnused      bool
		EnableGosimple    bool
	}
	Dupl struct {
		Threshold int
	}
	Goconst struct {
		MinStringLen        int
		MinOccurrencesCount int
	}

	EnabledLinters    []string
	DisabledLinters   []string
	EnableAllLinters  bool
	DisableAllLinters bool

	ExcludePatterns    []string
	UseDefaultExcludes bool

	Deadline time.Duration

	MaxIssuesPerLinter int

	DiffFromRevision  string
	DiffPatchFilePath string
	Diff              bool

	AnalyzeTests bool
}

type Config struct {
	Common Common
	Run    Run
}

func NewDefault() *Config {
	return &Config{
		Run: Run{
			OutFormat: OutFormatColoredLineNumber,
		},
	}
}
