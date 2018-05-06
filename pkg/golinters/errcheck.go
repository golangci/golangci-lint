package golinters

import (
	"context"
	"fmt"
	"strings"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/result"
	"github.com/golangci/golangci-shared/pkg/executors"
	errcheckAPI "github.com/kisielk/errcheck/golangci"
)

type errcheck struct{}

func (errcheck) Name() string {
	return "errcheck"
}

func (e errcheck) Run(ctx context.Context, exec executors.Executor, cfg *config.Run) (*result.Result, error) {
	errCfg := &cfg.Errcheck
	issues, err := errcheckAPI.Run(cfg.Paths.MixedPaths(), cfg.BuildTags, errCfg.CheckAssignToBlank, errCfg.CheckTypeAssertions)
	if err != nil {
		return nil, err
	}

	res := &result.Result{}
	for _, i := range issues {
		if !errCfg.CheckClose && strings.HasSuffix(i.FuncName, ".Close") {
			continue
		}

		var text string
		if i.FuncName != "" {
			text = fmt.Sprintf("Error return value of %s is not checked", formatCode(i.FuncName, cfg))
		} else {
			text = "Error return value is not checked"
		}
		res.Issues = append(res.Issues, result.Issue{
			FromLinter: e.Name(),
			Text:       text,
			LineNumber: i.Pos.Line,
			File:       i.Pos.Filename,
		})
	}

	return res, nil
}
