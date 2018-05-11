package golinters

import (
	"context"

	megacheckAPI "github.com/golangci/go-tools/cmd/megacheck"
	"github.com/golangci/golangci-lint/pkg/result"
)

type Megacheck struct {
	UnusedEnabled      bool
	GosimpleEnabled    bool
	StaticcheckEnabled bool
}

func (m Megacheck) Name() string {
	if m.UnusedEnabled && !m.GosimpleEnabled && !m.StaticcheckEnabled {
		return "unused"
	}
	if m.GosimpleEnabled && !m.UnusedEnabled && !m.StaticcheckEnabled {
		return "gosimple"
	}
	if m.StaticcheckEnabled && !m.UnusedEnabled && !m.GosimpleEnabled {
		return "staticcheck"
	}

	return "megacheck" // all enabled
}

func (m Megacheck) Desc() string {
	descs := map[string]string{
		"unused":      "Checks Go code for unused constants, variables, functions and types",
		"gosimple":    "Linter for Go source code that specialises on simplifying code",
		"staticcheck": "Staticcheck is go vet on steroids, applying a ton of static analysis checks",
		"megacheck":   "3 sub-linters in one: unused, gosimple and staticcheck",
	}

	return descs[m.Name()]
}

func (m Megacheck) Run(ctx context.Context, lintCtx *Context) ([]result.Issue, error) {
	issues := megacheckAPI.Run(lintCtx.Program, lintCtx.LoaderConfig, lintCtx.SSAProgram,
		m.StaticcheckEnabled, m.GosimpleEnabled, m.UnusedEnabled)

	var res []result.Issue
	for _, i := range issues {
		res = append(res, result.Issue{
			Pos:        i.Position,
			Text:       i.Text,
			FromLinter: m.Name(),
		})
	}
	return res, nil
}
