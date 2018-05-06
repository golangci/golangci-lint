package golinters

import (
	"context"
	"fmt"

	gocycloAPI "github.com/golangci/gocyclo"
	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/result"
	"github.com/golangci/golangci-shared/pkg/executors"
)

type gocyclo struct{}

func (gocyclo) Name() string {
	return "gocyclo"
}

func (g gocyclo) Run(ctx context.Context, exec executors.Executor, cfg *config.Run) (*result.Result, error) {
	stats := gocycloAPI.Run(cfg.Paths.MixedPaths())

	res := &result.Result{}
	for _, s := range stats {
		if s.Complexity < cfg.Gocyclo.MinComplexity {
			continue
		}

		res.Issues = append(res.Issues, result.Issue{
			File:       s.Pos.Filename,
			LineNumber: s.Pos.Line,
			Text: fmt.Sprintf("cyclomatic complexity %d of func %s is high (> %d)",
				s.Complexity, formatCode(s.FuncName, cfg), cfg.Gocyclo.MinComplexity),
			FromLinter: g.Name(),
		})
	}

	return res, nil
}
