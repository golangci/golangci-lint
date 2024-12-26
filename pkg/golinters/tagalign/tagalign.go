package tagalign

import (
	"github.com/4meepo/tagalign"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
)

func New(settings *config.TagAlignSettings) *goanalysis.Linter {
	var options []tagalign.Option

	if settings != nil {
		options = append(options, tagalign.WithAlign(settings.Align))

		if settings.Sort || len(settings.Order) > 0 {
			options = append(options, tagalign.WithSort(settings.Order...))
		}

		// Strict style will be applied only if Align and Sort are enabled together.
		if settings.Strict && settings.Align && settings.Sort {
			options = append(options, tagalign.WithStrictStyle())
		}
	}

	analyzer := tagalign.NewAnalyzer(options...)

	return goanalysis.NewLinter(
		analyzer.Name,
		analyzer.Doc,
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
