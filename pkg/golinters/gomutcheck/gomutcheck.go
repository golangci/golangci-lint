package gomutcheck

import (
	"github.com/BeyCoder/gomutcheck/pkg/analyzer"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
	"golang.org/x/tools/go/analysis"
)

func New() *goanalysis.Linter {
	return goanalysis.NewLinter(
		"gomutcheck",
		"Detect struct field mutations in value receiver methods.",
		[]*analysis.Analyzer{analyzer.New()},
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
