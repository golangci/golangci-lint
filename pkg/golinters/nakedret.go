package golinters

import (
	"fmt"
	"go/ast"
	"go/token"
	"sync"

	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

const nakedretName = "nakedret"

//nolint:dupl
func NewNakedret(settings *config.NakedretSettings) *goanalysis.Linter {
	var mu sync.Mutex
	var resIssues []goanalysis.Issue

	analyzer := &analysis.Analyzer{
		Name: nakedretName,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
		Run: func(pass *analysis.Pass) (any, error) {
			issues := runNakedRet(pass, settings)

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
		nakedretName,
		"Finds naked returns in functions greater than a specified function length",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithIssuesReporter(func(*linter.Context) []goanalysis.Issue {
		return resIssues
	}).WithLoadMode(goanalysis.LoadModeSyntax)
}

func runNakedRet(pass *analysis.Pass, settings *config.NakedretSettings) []goanalysis.Issue {
	var issues []goanalysis.Issue

	for _, file := range pass.Files {
		v := nakedretVisitor{
			maxLength: settings.MaxFuncLines,
			f:         pass.Fset,
		}
		v.root = &v

		ast.Walk(&v, file)

		for i := range v.issues {
			issues = append(issues, goanalysis.NewIssue(&v.issues[i], pass))
		}
	}

	return issues
}

type nakedretVisitor struct {
	maxLength int
	f         *token.FileSet
	issues    []result.Issue
	root      *nakedretVisitor

	// Details of the function we're currently dealing with
	funcName    string
	funcLength  int
	reportNaked bool
}

func hasNamedReturns(funcType *ast.FuncType) bool {
	if funcType == nil || funcType.Results == nil {
		return false
	}
	for _, field := range funcType.Results.List {
		for _, ident := range field.Names {
			if ident != nil {
				return true
			}
		}
	}
	return false
}

func (v *nakedretVisitor) Visit(node ast.Node) ast.Visitor {
	var (
		funcType *ast.FuncType
		funcName string
	)
	switch s := node.(type) {
	case *ast.FuncDecl:
		// We've found a function
		funcType = s.Type
		funcName = s.Name.Name
	case *ast.FuncLit:
		// We've found a function literal
		funcType = s.Type
		file := v.f.File(s.Pos())
		funcName = fmt.Sprintf("<func():%v>", file.Position(s.Pos()).Line)
	case *ast.ReturnStmt:
		// We've found a possibly naked return statement
		if v.reportNaked && len(s.Results) == 0 {
			v.root.issues = append(v.root.issues, result.Issue{
				FromLinter: nakedretName,
				Text: fmt.Sprintf("naked return in func `%s` with %d lines of code",
					v.funcName, v.funcLength),
				Pos: v.f.Position(s.Pos()),
			})
		}
	}

	if funcType != nil {
		// Create a new visitor to track returns for this function
		file := v.f.File(node.Pos())
		length := file.Position(node.End()).Line - file.Position(node.Pos()).Line
		return &nakedretVisitor{
			f:           v.f,
			root:        v.root,
			maxLength:   v.maxLength,
			funcName:    funcName,
			funcLength:  length,
			reportNaked: length > v.maxLength && hasNamedReturns(funcType),
		}
	}

	return v
}
