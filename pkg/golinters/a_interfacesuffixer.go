package golinters

import (
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	idmustbigint "github.com/tokopedia/mitralib/go/tokolinter/id-must-bigint"
	interfacesuffixer "github.com/tokopedia/mitralib/go/tokolinter/interface-suffix-er"
	missingfuncdoc "github.com/tokopedia/mitralib/go/tokolinter/missing-func-doc"
	noparamreassign "github.com/tokopedia/mitralib/go/tokolinter/no-param-reassign"
	"golang.org/x/tools/go/analysis"
)

func NewInterfaceSuffixEr() *goanalysis.Linter {
	return goanalysis.NewLinter(
		"interfacesuffixer",
		"interfaces are named by the method name plus an -er suffix or similar modification to construct an agent noun: Reader, Writer, Formatter, CloseNotifier etc.",
		[]*analysis.Analyzer{interfacesuffixer.Analyzer},
		nil,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}

func NewMissingFuncDoc() *goanalysis.Linter {
	return goanalysis.NewLinter(
		"missingfuncdoc",
		"every function must have proper documentation",
		[]*analysis.Analyzer{missingfuncdoc.Analyzer},
		nil,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}

func NewIDMustBigInt() *goanalysis.Linter {
	return goanalysis.NewLinter(
		"idmustbigint",
		"ID must use int64 to prevent integer overflow as the data grows",
		[]*analysis.Analyzer{idmustbigint.Analyzer},
		nil,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}

func NewNoParamReassign() *goanalysis.Linter {
	return goanalysis.NewLinter(
		"noparamreassign",
		"function param should not reasiigned or modified in the function body",
		[]*analysis.Analyzer{noparamreassign.Analyzer, interfacesuffixer.Analyzer},
		nil,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
