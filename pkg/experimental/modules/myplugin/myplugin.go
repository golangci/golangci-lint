package myplugin

import (
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/experimental/modules/register"
)

func init() {
	register.Plugin("foo", New)
}

type MySettings struct {
	Message string
}

type Foo struct {
	settings MySettings
}

func New(settings any) (register.LinterPlugin, error) {
	s, err := register.DecodeSettings[MySettings](settings)
	if err != nil {
		return nil, err
	}

	return &Foo{settings: s}, nil
}

func (f *Foo) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{
		{
			Name: "myplugin",
			Run: func(_ *analysis.Pass) (any, error) {
				println("I'm running:", f.settings.Message)
				return nil, nil
			},
			ResultType: nil,
		},
	}, nil
}

func (f *Foo) GetLoadMode() string {
	return register.LoadModeSyntax
}
