package config

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
}

type Run struct {
	Paths     []string
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
}

type Config struct {
	Common Common
	Run    Run
}

func NewDefault() *Config {
	return &Config{
		Run: Run{
			Paths:     []string{"./..."},
			OutFormat: OutFormatColoredLineNumber,
		},
	}
}
