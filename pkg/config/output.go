package config

import (
	"errors"
	"fmt"
	"slices"
)

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
	Format          string   `mapstructure:"format"`
	PrintIssuedLine bool     `mapstructure:"print-issued-lines"`
	PrintLinterName bool     `mapstructure:"print-linter-name"`
	UniqByLine      bool     `mapstructure:"uniq-by-line"`
	SortResults     bool     `mapstructure:"sort-results"`
	SortOrder       []string `mapstructure:"sort-order"`
	PathPrefix      string   `mapstructure:"path-prefix"`
	ShowStats       bool     `mapstructure:"show-stats"`
}

func (o *Output) Validate() error {
	if !o.SortResults && len(o.SortOrder) > 0 {
		return errors.New("sort-results should be 'true' to use sort-order")
	}

	for _, order := range o.SortOrder {
		if !slices.Contains([]string{"linter", "file", "severity"}, order) {
			return fmt.Errorf("unsupported sort-order name %q", order)
		}
	}

	return nil
}
