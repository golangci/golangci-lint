package golinters

import (
	"fmt"
	"sync"

	malignedAPI "github.com/golangci/maligned"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

const malignedName = "maligned"

//nolint:dupl
func NewMaligned(settings *config.MalignedSettings) *goanalysis.Linter {
	var mu sync.Mutex
	var res []goanalysis.Issue

	analyzer := &analysis.Analyzer{
		Name: malignedName,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
		Run: func(pass *analysis.Pass) (interface{}, error) {
			issues := runMaligned(pass, settings)

			if len(issues) == 0 {
				return nil, nil
			}

			mu.Lock()
			res = append(res, issues...)
			mu.Unlock()

			return nil, nil
		},
	}

	return goanalysis.NewLinter(
		malignedName,
		"Tool to detect Go structs that would take less memory if their fields were sorted",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithIssuesReporter(func(*linter.Context) []goanalysis.Issue {
		return res
	}).WithLoadMode(goanalysis.LoadModeTypesInfo)
}

func runMaligned(pass *analysis.Pass, settings *config.MalignedSettings) []goanalysis.Issue {
	prog := goanalysis.MakeFakeLoaderProgram(pass)

	malignedIssues := malignedAPI.Run(prog)
	if len(malignedIssues) == 0 {
		return nil
	}

	issues := make([]goanalysis.Issue, 0, len(malignedIssues))
	for _, i := range malignedIssues {
		text := fmt.Sprintf("struct of size %d bytes could be of size %d bytes", i.OldSize, i.NewSize)
		if settings.SuggestNewOrder {
			text += fmt.Sprintf(":\n%s", formatCodeBlock(i.NewStructDef, nil))
		}

		issues = append(issues, goanalysis.NewIssue(&result.Issue{
			Pos:        i.Pos,
			Text:       text,
			FromLinter: malignedName,
		}, pass))
	}

	return issues
}
