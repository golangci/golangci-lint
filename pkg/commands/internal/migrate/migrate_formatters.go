package migrate

import (
	"slices"
	"strings"

	"github.com/golangci/golangci-lint/pkg/commands/internal/migrate/one"
	"github.com/golangci/golangci-lint/pkg/commands/internal/migrate/ptr"
	"github.com/golangci/golangci-lint/pkg/commands/internal/migrate/two"
)

func toFormatters(old *one.Config) two.Formatters {
	enable, _ := ProcessEffectiveLinters(old.Linters)

	formatterNames := onlyFormatterNames(enable)

	var paths []string
	if len(formatterNames) != 0 {
		paths = slices.Concat(old.Issues.ExcludeFiles, old.Issues.ExcludeDirs)
	}

	return two.Formatters{
		Enable: formatterNames,
		Settings: two.FormatterSettings{
			Gci:       toGciSettings(old.LintersSettings.Gci),
			GoFmt:     toGoFmtSettings(old.LintersSettings.GoFmt),
			GoFumpt:   toGoFumptSettings(old.LintersSettings.GoFumpt),
			GoImports: toGoImportsSettings(old.LintersSettings.GoImports),
		},
		Exclusions: two.FormatterExclusions{
			Generated: toExclusionGenerated(old.Issues.ExcludeGenerated),
			Paths:     paths,
		},
	}
}

func toGciSettings(old one.GciSettings) two.GciSettings {
	return two.GciSettings{
		Sections:         old.Sections,
		NoInlineComments: old.NoInlineComments,
		NoPrefixComments: old.NoPrefixComments,
		CustomOrder:      old.CustomOrder,
		NoLexOrder:       old.CustomOrder,
	}
}

func toGoFmtSettings(old one.GoFmtSettings) two.GoFmtSettings {
	settings := two.GoFmtSettings{
		Simplify: old.Simplify,
	}

	for _, rule := range old.RewriteRules {
		settings.RewriteRules = append(settings.RewriteRules, two.GoFmtRewriteRule{
			Pattern:     rule.Pattern,
			Replacement: rule.Replacement,
		})
	}

	return settings
}

func toGoFumptSettings(old one.GoFumptSettings) two.GoFumptSettings {
	return two.GoFumptSettings{
		ModulePath: old.ModulePath,
		ExtraRules: old.ExtraRules,
	}
}

func toGoImportsSettings(old one.GoImportsSettings) two.GoImportsSettings {
	var localPrefixes []string

	if ptr.Deref(old.LocalPrefixes) != "" {
		localPrefixes = strings.Split(ptr.Deref(old.LocalPrefixes), ",")
	}

	return two.GoImportsSettings{
		LocalPrefixes: localPrefixes,
	}
}
