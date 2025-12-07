package goroutinedeferguard

import (
	"io"
	"log"

	"github.com/status-im/goroutine-defer-guard/pkg/analyzer"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
)

func New(settings *config.GoroutineDeferGuardSettings) *goanalysis.Linter {
	logger := log.New(io.Discard, "", 0)

	var cfg map[string]any
	if settings != nil && settings.Target != "" {
		cfg = map[string]any{
			"target": settings.Target,
		}
	}

	a := analyzer.New(logger)

	return goanalysis.
		NewLinterFromAnalyzer(a).
		WithConfig(cfg).
		WithLoadMode(goanalysis.LoadModeTypesInfo)
}
