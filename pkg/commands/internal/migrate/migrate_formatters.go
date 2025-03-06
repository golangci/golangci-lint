package migrate

import (
	"slices"
	"strings"

	"github.com/golangci/golangci-lint/pkg/commands/internal/migrate/ptr"
	"github.com/golangci/golangci-lint/pkg/commands/internal/migrate/versionone"
	"github.com/golangci/golangci-lint/pkg/commands/internal/migrate/versiontwo"
)

func toFormatters(old *versionone.Config) versiontwo.Formatters {
	enable, _ := ProcessEffectiveLinters(old.Linters)

	formatterNames := onlyFormatterNames(enable)

	var paths []string
	if len(formatterNames) != 0 {
		paths = slices.Concat(old.Issues.ExcludeFiles, old.Issues.ExcludeDirs)
	}

	return versiontwo.Formatters{
		Enable: formatterNames,
		Settings: versiontwo.FormatterSettings{
			Gci:       toGciSettings(old.LintersSettings.Gci),
			GoFmt:     toGoFmtSettings(old.LintersSettings.GoFmt),
			GoFumpt:   toGoFumptSettings(old.LintersSettings.GoFumpt),
			GoImports: toGoImportsSettings(old.LintersSettings.GoImports),
		},
		Exclusions: versiontwo.FormatterExclusions{
			Generated: toExclusionGenerated(old.Issues.ExcludeGenerated),
			Paths:     paths,
		},
	}
}

func toGciSettings(old versionone.GciSettings) versiontwo.GciSettings {
	return versiontwo.GciSettings{
		Sections:         old.Sections,
		NoInlineComments: old.NoInlineComments,
		NoPrefixComments: old.NoPrefixComments,
		CustomOrder:      old.CustomOrder,
		NoLexOrder:       old.CustomOrder,
	}
}

func toGoFmtSettings(old versionone.GoFmtSettings) versiontwo.GoFmtSettings {
	settings := versiontwo.GoFmtSettings{
		Simplify: old.Simplify,
	}

	for _, rule := range old.RewriteRules {
		settings.RewriteRules = append(settings.RewriteRules, versiontwo.GoFmtRewriteRule{
			Pattern:     rule.Pattern,
			Replacement: rule.Replacement,
		})
	}

	return settings
}

func toGoFumptSettings(old versionone.GoFumptSettings) versiontwo.GoFumptSettings {
	return versiontwo.GoFumptSettings{
		ModulePath: old.ModulePath,
		ExtraRules: old.ExtraRules,
	}
}

func toGoImportsSettings(old versionone.GoImportsSettings) versiontwo.GoImportsSettings {
	var localPrefixes []string

	prefixes := ptr.Deref(old.LocalPrefixes)
	if prefixes != "" {
		localPrefixes = strings.Split(prefixes, ",")
	}

	return versiontwo.GoImportsSettings{
		LocalPrefixes: localPrefixes,
	}
}
