package golinters

import (
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/packages"
	"honnef.co/go/tools/unused"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

func NewUnused() *goanalysis.Linter {
	u := unused.NewChecker(false)
	analyzers := []*analysis.Analyzer{u.Analyzer()}
	setAnalyzersGoVersion(analyzers)

	const name = "unused"
	lnt := goanalysis.NewLinter(
		name,
		"Checks Go code for unused constants, variables, functions and types",
		analyzers,
		nil,
	).WithIssuesReporter(func(lintCtx *linter.Context) []goanalysis.Issue {
		typesToPkg := map[*types.Package]*packages.Package{}
		for _, pkg := range lintCtx.OriginalPackages {
			typesToPkg[pkg.Types] = pkg
		}

		var issues []goanalysis.Issue
		for _, ur := range u.Result() {
			p := u.ProblemObject(lintCtx.Packages[0].Fset, ur)
			pkg := typesToPkg[ur.Pkg()]
			i := &result.Issue{
				FromLinter: name,
				Text:       p.Message,
				Pos:        p.Pos,
				Pkg:        pkg,
				LineRange: &result.Range{
					From: p.Pos.Line,
					To:   p.End.Line,
				},
			}
			// See https://github.com/golangci/golangci-lint/issues/1048
			// If range is invalid, this will break `--fix` mode.
			if i.LineRange.To >= i.LineRange.From {
				i.Replacement = &result.Replacement{
					// Suggest deleting unused stuff.
					NeedOnlyDelete: true,
				}
			}
			issues = append(issues, goanalysis.NewIssue(i, nil))
		}
		return issues
	}).WithContextSetter(func(lintCtx *linter.Context) {
		u.WholeProgram = lintCtx.Settings().Unused.CheckExported
	}).WithLoadMode(goanalysis.LoadModeWholeProgram)
	lnt.UseOriginalPackages()
	return lnt
}
