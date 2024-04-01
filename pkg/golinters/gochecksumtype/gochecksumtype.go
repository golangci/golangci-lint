package gochecksumtype

import (
	"strings"
	"sync"

	gochecksumtype "github.com/alecthomas/go-check-sumtype"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/packages"

	"github.com/golangci/golangci-lint/pkg/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

const name = "gochecksumtype"

func New() *goanalysis.Linter {
	var mu sync.Mutex
	var resIssues []goanalysis.Issue

	analyzer := &analysis.Analyzer{
		Name: name,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
		Run: func(pass *analysis.Pass) (any, error) {
			issues, err := runGoCheckSumType(pass)
			if err != nil {
				return nil, err
			}

			if len(issues) == 0 {
				return nil, nil
			}

			mu.Lock()
			resIssues = append(resIssues, issues...)
			mu.Unlock()

			return nil, nil
		},
	}

	return goanalysis.NewLinter(
		name,
		`Run exhaustiveness checks on Go "sum types"`,
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithIssuesReporter(func(_ *linter.Context) []goanalysis.Issue {
		return resIssues
	}).WithLoadMode(goanalysis.LoadModeTypesInfo)
}

func runGoCheckSumType(pass *analysis.Pass) ([]goanalysis.Issue, error) {
	var resIssues []goanalysis.Issue

	pkg := &packages.Package{
		Fset:      pass.Fset,
		Syntax:    pass.Files,
		Types:     pass.Pkg,
		TypesInfo: pass.TypesInfo,
	}

	var unknownError error
	errors := gochecksumtype.Run([]*packages.Package{pkg})
	for _, err := range errors {
		err, ok := err.(gochecksumtype.Error)
		if !ok {
			unknownError = err
			continue
		}

		resIssues = append(resIssues, goanalysis.NewIssue(&result.Issue{
			FromLinter: name,
			Text:       strings.TrimPrefix(err.Error(), err.Pos().String()+": "),
			Pos:        err.Pos(),
		}, pass))
	}

	return resIssues, unknownError
}
