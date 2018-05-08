package golinters

import (
	"context"
	"fmt"

	gocycloAPI "github.com/golangci/gocyclo"
	"github.com/golangci/golangci-lint/pkg/result"
)

type Gocyclo struct{}

func (Gocyclo) Name() string {
	return "gocyclo"
}

func (Gocyclo) Desc() string {
	return "Computes and checks the cyclomatic complexity of functions"
}

func (g Gocyclo) Run(ctx context.Context, lintCtx *Context) ([]result.Issue, error) {
	stats := gocycloAPI.Run(lintCtx.Paths.MixedPaths())

	var res []result.Issue
	for _, s := range stats {
		if s.Complexity < lintCtx.RunCfg().Gocyclo.MinComplexity {
			continue
		}

		res = append(res, result.Issue{
			Pos: s.Pos,
			Text: fmt.Sprintf("cyclomatic complexity %d of func %s is high (> %d)",
				s.Complexity, formatCode(s.FuncName, lintCtx.RunCfg()), lintCtx.RunCfg().Gocyclo.MinComplexity),
			FromLinter: g.Name(),
		})
	}

	return res, nil
}
