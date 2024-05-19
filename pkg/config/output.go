package config

import (
	"errors"
	"fmt"
	"slices"
	"strings"
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
	OutFormatSarif             = "sarif"
)

var AllOutputFormats = []string{
	OutFormatJSON,
	OutFormatLineNumber,
	OutFormatColoredLineNumber,
	OutFormatTab,
	OutFormatColoredTab,
	OutFormatCheckstyle,
	OutFormatCodeClimate,
	OutFormatHTML,
	OutFormatJunitXML,
	OutFormatGithubActions,
	OutFormatTeamCity,
	OutFormatSarif,
}

type Output struct {
	Formats         OutputFormats `mapstructure:"formats"`
	PrintIssuedLine bool          `mapstructure:"print-issued-lines"`
	PrintLinterName bool          `mapstructure:"print-linter-name"`
	UniqByLine      bool          `mapstructure:"uniq-by-line"`
	SortResults     bool          `mapstructure:"sort-results"`
	SortOrder       []string      `mapstructure:"sort-order"`
	PathPrefix      string        `mapstructure:"path-prefix"`
	ShowStats       bool          `mapstructure:"show-stats"`

	// Deprecated: use Formats instead.
	Format string `mapstructure:"format"`
}

func (o *Output) Validate() error {
	if !o.SortResults && len(o.SortOrder) > 0 {
		return errors.New("sort-results should be 'true' to use sort-order")
	}

	validOrders := []string{"linter", "file", "severity"}

	all := strings.Join(o.SortOrder, " ")

	for _, order := range o.SortOrder {
		if strings.Count(all, order) > 1 {
			return fmt.Errorf("the sort-order name %q is repeated several times", order)
		}

		if !slices.Contains(validOrders, order) {
			return fmt.Errorf("unsupported sort-order name %q", order)
		}
	}

	for _, format := range o.Formats {
		err := format.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}

type OutputFormat struct {
	Format string `mapstructure:"format"`
	Path   string `mapstructure:"path"`
}

func (o *OutputFormat) Validate() error {
	if o.Format == "" {
		return errors.New("the format is required")
	}

	if !slices.Contains(AllOutputFormats, o.Format) {
		return fmt.Errorf("unsupported output format %q", o.Format)
	}

	return nil
}

type OutputFormats []OutputFormat

func (p *OutputFormats) UnmarshalText(text []byte) error {
	formats := strings.Split(string(text), ",")

	for _, item := range formats {
		format, path, _ := strings.Cut(item, ":")

		*p = append(*p, OutputFormat{
			Path:   path,
			Format: format,
		})
	}

	return nil
}
