package golinters

import (
	"context"
	"fmt"
	"strings"

	megacheckAPI "github.com/golangci/go-tools/cmd/megacheck"
	"github.com/golangci/golangci-lint/pkg/lint"
	"github.com/golangci/golangci-lint/pkg/result"
)

type Megacheck struct {
	UnusedEnabled      bool
	GosimpleEnabled    bool
	StaticcheckEnabled bool
}

func (m Megacheck) Name() string {
	names := []string{}
	if m.UnusedEnabled {
		names = append(names, "unused")
	}
	if m.GosimpleEnabled {
		names = append(names, "gosimple")
	}
	if m.StaticcheckEnabled {
		names = append(names, "staticcheck")
	}

	if len(names) == 1 {
		return names[0] // only one sublinter is enabled
	}

	if len(names) == 3 {
		return "megacheck" // all enabled
	}

	return fmt.Sprintf("megacheck.{%s}", strings.Join(names, ","))
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

func (m Megacheck) Run(ctx context.Context, lintCtx *lint.Context) ([]result.Issue, error) {
	issues := megacheckAPI.Run(lintCtx.Program, lintCtx.LoaderConfig, lintCtx.SSAProgram,
		m.StaticcheckEnabled, m.GosimpleEnabled, m.UnusedEnabled)
	if len(issues) == 0 {
		return nil, nil
	}

	res := make([]result.Issue, 0, len(issues))
	for _, i := range issues {
		res = append(res, result.Issue{
			Pos:        i.Position,
			Text:       i.Text,
			FromLinter: m.Name(),
		})
	}
	return res, nil
}
