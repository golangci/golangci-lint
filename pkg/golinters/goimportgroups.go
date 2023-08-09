package golinters

import (
	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"github.com/kmirzavaziri/goimportgroups/pkg/analyzer"
	"golang.org/x/tools/go/analysis"
	"strings"
)

func NewGoImportGroups(settings *config.GoImportGroupsSettings) *goanalysis.Linter {
	a := analyzer.NewAnalyzer()

	cfgMap := map[string]map[string]any{
		a.Name: {
			"groups": strings.Join(settings.Groups, ";"),
		},
	}

	return goanalysis.NewLinter(
		"goimportgroups",
		"Checks if go imports are separated into user-defined groups. "+
			"Define each group as a left associative boolean expression of import path regex patterns, "+
			"i.e.: p1,p2:p3 (comma is logical AND (&&), and colon is logical OR (||))",
		[]*analysis.Analyzer{a},
		cfgMap,
	)
}
