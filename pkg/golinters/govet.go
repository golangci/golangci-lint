package golinters

import (
	"context"

	"github.com/golangci/golangci-lint/pkg/result"
	govetAPI "github.com/golangci/govet"
)

type Govet struct{}

func (Govet) Name() string {
	return "govet"
}

func (Govet) Desc() string {
	return "Vet examines Go source code and reports suspicious constructs, such as Printf calls whose arguments do not align with the format string"
}

func (g Govet) Run(ctx context.Context, lintCtx *Context) ([]result.Issue, error) {
	// TODO: check .S asm files: govet can do it if pass dirs
	var govetIssues []govetAPI.Issue
	for _, files := range lintCtx.Paths.FilesGrouppedByDirs() {
		issues, err := govetAPI.Run(files, lintCtx.Settings().Govet.CheckShadowing)
		if err != nil {
			return nil, err
		}
		govetIssues = append(govetIssues, issues...)
	}

	var res []result.Issue
	for _, i := range govetIssues {
		res = append(res, result.Issue{
			Pos:        i.Pos,
			Text:       i.Message,
			FromLinter: g.Name(),
		})
	}
	return res, nil
}
