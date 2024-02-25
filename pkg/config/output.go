package config

const (
	OutFormatJSON              = "json"
	OutFormatLineNumber        = "line-number"
	OutFormatColoredLineNumber = "colored-line-number"
	OutFormatTab               = "tab"
	OutFormatColoredTab        = "colored-tab"
	OutFormatCheckstyle        = "checkstyle"
	OutFormatCodeClimate       = "code-climate"
	OutFormatHTML              = "html"
	OutFormatJunitXML          = "junit-xml"
	OutFormatGithubActions     = "github-actions"
	OutFormatTeamCity          = "teamcity"
)

var OutFormats = []string{
	OutFormatColoredLineNumber,
	OutFormatLineNumber,
	OutFormatJSON,
	OutFormatTab,
	OutFormatCheckstyle,
	OutFormatCodeClimate,
	OutFormatHTML,
	OutFormatJunitXML,
	OutFormatGithubActions,
	OutFormatTeamCity,
}

type Output struct {
	Format          string `mapstructure:"format"`
	PrintIssuedLine bool   `mapstructure:"print-issued-lines"`
	PrintLinterName bool   `mapstructure:"print-linter-name"`
	UniqByLine      bool   `mapstructure:"uniq-by-line"`
	SortResults     bool   `mapstructure:"sort-results"`
	PathPrefix      string `mapstructure:"path-prefix"`
}
