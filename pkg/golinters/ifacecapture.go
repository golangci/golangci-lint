package golinters

import (
	"golang.org/x/tools/go/analysis"

	"github.com/dgunay/ifacecapture/ifacecapture"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewInterfaceCaptureCheck() *goanalysis.Linter {
	ifacecaptureAnalyzer := ifacecapture.Analyzer
	return goanalysis.NewLinter(
		ifacecaptureAnalyzer.Name,
		ifacecaptureAnalyzer.Doc,
		[]*analysis.Analyzer{ifacecaptureAnalyzer},
		nil,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
