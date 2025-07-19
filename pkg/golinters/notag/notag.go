package notag

import (
	"github.com/guerinoni/notag/pkg/analyzer"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
)

func New(settings *config.NoTagSettings) *goanalysis.Linter {
	c := analyzer.Setting{
		GlobalTagsDenied: settings.Denied,
		Pkg:              settings.DeniedPkg,
		PkgPath:          settings.DeniedPkgPath,
	}
	a := analyzer.NewAnalyzerWithConfig(c)
	return goanalysis.NewLinterFromAnalyzer(a)
}
