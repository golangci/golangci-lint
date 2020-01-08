package magic_numbers

import (
	"flag"
	"go/ast"

	"github.com/tommy-muehle/go-mnd/checks"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const Doc = `magic number detector`

var Analyzer = &analysis.Analyzer{
	Name:             "mnd",
	Doc:              Doc,
	Run:              run,
	Flags:            options(),
	Requires:         []*analysis.Analyzer{inspect.Analyzer},
	RunDespiteErrors: true,
}

type Checker interface {
	NodeFilter() []ast.Node
	Check(n ast.Node)
}

func options() flag.FlagSet {
	options := flag.NewFlagSet("", flag.ExitOnError)
	options.String("checks", "", "comma separated list of checks")

	return *options
}

func run(pass *analysis.Pass) (interface{}, error) {
	config := WithOptions(
		WithCustomChecks(pass.Analyzer.Flags.Lookup("checks").Value.String()),
	)

	var checker []Checker
	if config.IsCheckEnabled(checks.ArgumentCheck) {
		checker = append(checker, checks.NewArgumentAnalyzer(pass))
	}
	if config.IsCheckEnabled(checks.CaseCheck) {
		checker = append(checker, checks.NewCaseAnalyzer(pass))
	}
	if config.IsCheckEnabled(checks.ConditionCheck) {
		checker = append(checker, checks.NewConditionAnalyzer(pass))
	}
	if config.IsCheckEnabled(checks.OperationCheck) {
		checker = append(checker, checks.NewOperationAnalyzer(pass))
	}
	if config.IsCheckEnabled(checks.ReturnCheck) {
		checker = append(checker, checks.NewReturnAnalyzer(pass))
	}
	if config.IsCheckEnabled(checks.AssignCheck) {
		checker = append(checker, checks.NewAssignAnalyzer(pass))
	}

	i := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	for _, c := range checker {
		i.Preorder(c.NodeFilter(), func(node ast.Node) {
			c.Check(node)
		})
	}

	return nil, nil
}
