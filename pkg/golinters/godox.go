package golinters

import (
	"context"
	"go/token"
	"strings"

	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"

	"github.com/matoous/godox"
)

type Godox struct{}

func (Godox) Name() string {
	return "godox"
}

func (Godox) Desc() string {
	return "Tool for detection of FIXME, TODO and other comment keywords"
}

func (f Godox) Run(ctx context.Context, lintCtx *linter.Context) ([]result.Issue, error) {
	var issues []godox.Message
	for _, file := range lintCtx.ASTCache.GetAllValidFiles() {
		issues = append(issues, godox.Run(file.F, file.Fset, lintCtx.Settings().Godox.Keywords...)...)
	}

	if len(issues) == 0 {
		return nil, nil
	}

	res := make([]result.Issue, len(issues))
	for k, i := range issues {
		res[k] = result.Issue{
			Pos: token.Position{
				Filename: i.Pos.Filename,
				Line:     i.Pos.Line,
			},
			Text:       strings.TrimRight(i.Message, "\n"),
			FromLinter: f.Name(),
		}
	}
	return res, nil
}
