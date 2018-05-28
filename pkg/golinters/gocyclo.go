package golinters

import (
	"context"
	"fmt"
	"sort"

	gocycloAPI "github.com/golangci/gocyclo/pkg/gocyclo"
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
	var stats []gocycloAPI.Stat
	for _, f := range lintCtx.ASTCache.GetAllValidFiles() {
		stats = gocycloAPI.BuildStats(f.F, f.Fset, stats)
	}

	sort.Slice(stats, func(i, j int) bool {
		return stats[i].Complexity > stats[j].Complexity
	})

	var res []result.Issue
	for _, s := range stats {
		if s.Complexity <= lintCtx.Settings().Gocyclo.MinComplexity {
			continue
		}

		res = append(res, result.Issue{
			Pos: s.Pos,
			Text: fmt.Sprintf("cyclomatic complexity %d of func %s is high (> %d)",
				s.Complexity, formatCode(s.FuncName, lintCtx.Cfg), lintCtx.Settings().Gocyclo.MinComplexity),
			FromLinter: g.Name(),
		})
	}

	return res, nil
}
