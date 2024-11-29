package errchkjson

import (
	"github.com/breml/errchkjson"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
)

func New(settings *config.ErrChkJSONSettings) *goanalysis.Linter {
	a := errchkjson.NewAnalyzer()

	cfg := map[string]map[string]any{}
	cfg[a.Name] = map[string]any{
		"omit-safe": true,
	}
	if settings != nil {
		cfg[a.Name] = map[string]any{
			"omit-safe":          !settings.CheckErrorFreeEncoding,
			"report-no-exported": settings.ReportNoExported,
		}
	}

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		cfg,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
