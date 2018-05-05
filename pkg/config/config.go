package config

type OutFormat string

const (
	OutFormatJSON              = "json"
	OutFormatLineNumber        = "line-number"
	OutFormatColoredLineNumber = "colored-line-number"
)

var OutFormats = []string{OutFormatColoredLineNumber, OutFormatLineNumber, OutFormatJSON}

type Common struct {
	IsVerbose bool
}

type Run struct {
	Paths                 []string
	OutFormat             string
	ExitCodeIfIssuesFound int
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
