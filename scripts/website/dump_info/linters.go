package main

import (
	"maps"
	"slices"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goformatters"
	"github.com/golangci/golangci-lint/v2/pkg/lint/linter"
	"github.com/golangci/golangci-lint/v2/pkg/lint/lintersdb"
	"github.com/golangci/golangci-lint/v2/scripts/website/types"
)

func saveLinters(dst string) error {
	linters, _ := lintersdb.NewLinterBuilder().Build(config.NewDefault())

	var wraps []types.LinterWrapper
	for _, l := range linters {
		if l.IsDeprecated() && l.Deprecation.Level > linter.DeprecationWarning {
			continue
		}

		if goformatters.IsFormatter(l.Name()) {
			continue
		}

		wrapper := types.LinterWrapper{
			Name:             l.Linter.Name(),
			Desc:             l.Linter.Desc(),
			Groups:           slices.AppendSeq(make([]string, 0, len(l.Groups)), maps.Keys(l.Groups)),
			LoadMode:         l.LoadMode,
			AlternativeNames: l.AlternativeNames,
			OriginalURL:      l.OriginalURL,
			Internal:         l.Internal,
			CanAutoFix:       l.CanAutoFix,
			IsSlow:           l.IsSlow,
			DoesChangeTypes:  l.DoesChangeTypes,
			Since:            l.Since,
		}

		if l.Deprecation != nil {
			wrapper.Deprecation = &types.Deprecation{
				Since:       l.Deprecation.Since,
				Message:     l.Deprecation.Message,
				Replacement: l.Deprecation.Replacement,
			}
		}

		wraps = append(wraps, wrapper)
	}

	return saveToJSONFile(dst, wraps)
}
