package golinters

import (
	"github.com/adamdecaf/xmlencoderclose/pkg/analyzer"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewXMLEncoderClose() *goanalysis.Linter {
	return goanalysis.NewLinter(
		"xmlencoderclose",
		"Checks that xml.Encoder is closed",
		[]*analysis.Analyzer{
			analyzer.NewAnalyzer(),
		},
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
