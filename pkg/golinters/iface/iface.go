package iface

import (
	"github.com/uudashr/iface/identical"
	"github.com/uudashr/iface/opaque"
	"github.com/uudashr/iface/unused"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
)

func New(settings *config.IfaceSettings) *goanalysis.Linter {
	return goanalysis.NewLinter(
		"iface",
		"Detect the incorrect use of interfaces, helping developers avoid interface pollution.",
		analyzersFromSettings(settings),
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}

func analyzersFromSettings(settings *config.IfaceSettings) []*analysis.Analyzer {
	if settings == nil { // FIXME
		return nil
	}

	var analyzers []*analysis.Analyzer

	if settings.Unused {
		analyzers = append(analyzers, unused.Analyzer)
	}

	if settings.Identical {
		analyzers = append(analyzers, identical.Analyzer)
	}

	if settings.Opaque {
		analyzers = append(analyzers, opaque.Analyzer)
	}

	return analyzers
}
