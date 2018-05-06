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

func (g Gocyclo) Run(ctx context.Context, lintCtx *Context) (*result.Result, error) {
	stats := gocycloAPI.Run(lintCtx.Paths.MixedPaths())

	res := &result.Result{}
	for _, s := range stats {
		if s.Complexity < lintCtx.RunCfg().Gocyclo.MinComplexity {
			continue
		}

		res.Issues = append(res.Issues, result.Issue{
			File:       s.Pos.Filename,
			LineNumber: s.Pos.Line,
			Text: fmt.Sprintf("cyclomatic complexity %d of func %s is high (> %d)",
				s.Complexity, formatCode(s.FuncName, lintCtx.RunCfg()), lintCtx.RunCfg().Gocyclo.MinComplexity),
			FromLinter: g.Name(),
		})
	}

	return res, nil
}
