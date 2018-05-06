package config

import "github.com/golangci/golangci-lint/pkg/fsutils"

type OutFormat string

const (
	OutFormatJSON              = "json"
	OutFormatLineNumber        = "line-number"
	OutFormatColoredLineNumber = "colored-line-number"
)

var OutFormats = []string{OutFormatColoredLineNumber, OutFormatLineNumber, OutFormatJSON}

type Common struct {
	IsVerbose      bool
	CPUProfilePath string
	Concurrency    int
}

type Run struct {
	Paths     *fsutils.ProjectPaths
	BuildTags []string

	OutFormat             string
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
