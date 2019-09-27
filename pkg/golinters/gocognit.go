// nolint:dupl
package golinters

import (
	"context"
	"fmt"
	"sort"

	"github.com/uudashr/gocognit"

	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

type Gocognit struct{}

func (Gocognit) Name() string {
	return "gocognit"
}

func (Gocognit) Desc() string {
	return "Computes and checks the cognitive complexity of functions"
}

func (g Gocognit) Run(ctx context.Context, lintCtx *linter.Context) ([]result.Issue, error) {
	var stats []gocognit.Stat
	for _, f := range lintCtx.ASTCache.GetAllValidFiles() {
		stats = gocognit.ComplexityStats(f.F, f.Fset, stats)
	}

	if len(stats) == 0 {
		return nil, nil
	}

	sort.Slice(stats, func(i, j int) bool {
		return stats[i].Complexity > stats[j].Complexity
	})

	res := make([]result.Issue, 0, len(stats))
	for _, s := range stats {
		if s.Complexity <= lintCtx.Settings().Gocognit.MinComplexity {
			break // Break as the stats is already sorted from greatest to least
		}

		res = append(res, result.Issue{
			Pos: s.Pos,
			Text: fmt.Sprintf("cognitive complexity %d of func %s is high (> %d)",
				s.Complexity, formatCode(s.FuncName, lintCtx.Cfg), lintCtx.Settings().Gocognit.MinComplexity),
			FromLinter: g.Name(),
		})
	}

	return res, nil
}
