package protogetter

import (
	"github.com/ghostiam/protogetter"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
)

func New(settings *config.ProtoGetterSettings) *goanalysis.Linter {
	var cfg protogetter.Config
	if settings != nil {
		cfg = protogetter.Config{
			SkipGeneratedBy:         settings.SkipGeneratedBy,
			SkipFiles:               settings.SkipFiles,
			SkipAnyGenerated:        settings.SkipAnyGenerated,
			ReplaceFirstArgInAppend: settings.ReplaceFirstArgInAppend,
		}
	}

	cfg.Mode = protogetter.StandaloneMode

	a := protogetter.NewAnalyzer(&cfg)

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
