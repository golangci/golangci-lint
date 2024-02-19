package golinters

import (
	"github.com/rezkam/noniljson"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewNoNilJSON() *goanalysis.Linter {
	a := noniljson.Analyzer

	return goanalysis.NewLinter(
		a.Name,
		"checks that nullable fields in structs used for JSON marshaling use 'omitempty'",
		[]*analysis.Analyzer{a},
		nil,
	)
}
