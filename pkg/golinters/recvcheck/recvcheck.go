package recvcheck

import (
	"github.com/raeperd/recvcheck"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/goanalysis"
)

func New() *goanalysis.Linter {
	a := recvcheck.Analyzer

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
