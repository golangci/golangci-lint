package goheader

import (
	"strings"

	goheader "github.com/denis-tingaikin/go-header"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
)

const linterName = "goheader"

func New(settings *config.GoHeaderSettings, replacer *strings.Replacer) *goanalysis.Linter {
	conf := &goheader.Config{}
	if settings != nil {
		conf = &goheader.Config{
			Values:       settings.Values,
			Template:     settings.Template,
			TemplatePath: replacer.Replace(settings.TemplatePath),
		}
	}
	return goanalysis.NewLinter(
		linterName,
		"Checks if file header matches to pattern",
		[]*analysis.Analyzer{goheader.NewAnalyzer(conf)},
		nil,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
