package db

import (
	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/lint/lintersdb"
	"github.com/golangci/golangci-lint/pkg/logutils"
)

type Builder interface {
	Build(cfg *config.Config) []*linter.Config
}

type PluginBuilderWrapper struct {
	builders []Builder
}

func NewPluginBuilderWrapper(log logutils.Log) *PluginBuilderWrapper {
	return &PluginBuilderWrapper{
		builders: []Builder{
			NewPluginModuleBuilder(log),
			lintersdb.NewPluginBuilder(log),
		},
	}
}

func (b *PluginBuilderWrapper) Build(cfg *config.Config) []*linter.Config {
	var linters []*linter.Config

	for _, builder := range b.builders {
		linters = append(linters, builder.Build(cfg)...)
	}

	return linters
}
