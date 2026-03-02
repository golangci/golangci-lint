package kubetyped

import (
	"github.com/jplimack-ai/kubetyped"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
)

func New(settings *config.KubetypedSettings) *goanalysis.Linter {
	s := kubetyped.Settings{}
	if settings != nil {
		s.IncludeTestFiles = settings.IncludeTestFiles
		s.IgnoreGVKs = settings.IgnoreGVKs
		s.RejectGVKs = settings.RejectGVKs
		for _, g := range settings.ExtraKnownGVKs {
			s.ExtraKnownGVKs = append(s.ExtraKnownGVKs, kubetyped.GVKEntry{
				APIVersion:   g.APIVersion,
				Kind:         g.Kind,
				TypedPackage: g.TypedPackage,
			})
		}
		for name, enabled := range settings.Checks {
			if s.Checks == nil {
				s.Checks = make(map[string]kubetyped.CheckConfig)
			}
			s.Checks[name] = kubetyped.CheckConfig{Enabled: &enabled}
		}
	}
	return goanalysis.
		NewLinterFromAnalyzer(kubetyped.NewAnalyzer(&s)).
		WithLoadMode(goanalysis.LoadModeTypesInfo)
}
