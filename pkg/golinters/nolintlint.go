package golinters

import (
	"context"
	"fmt"
	"go/ast"

	"github.com/ashanbrown/nolintlint/v2/nolintlint"

	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

type NoLintLint struct{}

func (NoLintLint) Name() string {
	return "nolintlint"
}

func (NoLintLint) Desc() string {
	return "Reports ill-formed or insufficient nolint directives"
}

func (l NoLintLint) Run(ctx context.Context, lintCtx *linter.Context) (results []result.Issue, err error) {
	var needs nolintlint.Needs
	settings := lintCtx.Settings().NoLintLint
	if settings.Explain {
		needs |= nolintlint.NeedsExplanation
	}
	if settings.Machine {
		needs |= nolintlint.NeedsMachine
	}
	if settings.Specific {
		needs |= nolintlint.NeedsSpecific
	}
	lnt, err := nolintlint.NewLinter(
		nolintlint.OptionNeeds(needs),
		nolintlint.OptionExcludes(settings.Exclude),
	)
	if err != nil {
		return nil, err
	}
	for _, pkg := range lintCtx.Packages {
		files, fset, err := getASTFilesForGoPkg(lintCtx, pkg)
		if err != nil {
			return nil, fmt.Errorf("could not load files: %s", err)
		}
		nodes := make([]ast.Node, 0, len(files))
		for _, n := range files {
			nodes = append(nodes, n)
		}
		issues, err := lnt.Run(fset, nodes...)
		if err != nil {
			return nil, fmt.Errorf("linter failed to run: %s", err)
		}
		for _, i := range issues {
			results = append(results, result.Issue{
				FromLinter: l.Name(),
				Text:       i.Details(),
				Pos:        i.Position(),
			})
		}
	}
	return results, nil
}
