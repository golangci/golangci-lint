package multisplit

import (
	"github.com/kenyoni-software/go-multisplit/multisplit"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
	"github.com/golangci/golangci-lint/v2/pkg/golinters/internal"
)

func New(settings *config.MultiSplitSettings) *goanalysis.Linter {
	analyzer := multisplit.NewAnalyzer()
	analyzer.Settings = toMultiSplitSettings(settings)
	return goanalysis.
		NewLinterFromAnalyzer(analyzer.Analyzer).
		WithLoadMode(goanalysis.LoadModeSyntax)
}

func toMultiSplitSettings(settings *config.MultiSplitSettings) multisplit.Settings {
	cfg := multisplit.DefaultSettings()
	if settings == nil {
		return cfg
	}

	ruleMap := map[string]*bool{
		"assign":             &cfg.Assign,
		"const-decl-func":    &cfg.ConstDeclFunc,
		"const-decl-pkg":     &cfg.ConstDeclPkg,
		"func-params":        &cfg.FuncParams,
		"func-return-values": &cfg.FuncReturnValues,
		"short-var-decl":     &cfg.ShortVarDecl,
		"struct-fields":      &cfg.StructFields,
		"var-decl-func":      &cfg.VarDeclFunc,
		"var-decl-pkg":       &cfg.VarDeclPkg,
		"var-decl-init-func": &cfg.VarDeclInitFunc,
		"var-decl-init-pkg":  &cfg.VarDeclInitPkg,
	}

	for _, rule := range settings.Rules {
		if target, ok := ruleMap[rule]; ok {
			*target = true
		} else {
			internal.LinterLogger.Fatalf("multisplit: unknown rule '%s'", rule)
		}
	}

	if settings.ConstDeclFuncToBlock != nil {
		cfg.ConstDeclFuncToBlock = *settings.ConstDeclFuncToBlock
	}
	if settings.ConstDeclPkgToBlock != nil {
		cfg.ConstDeclPkgToBlock = *settings.ConstDeclPkgToBlock
	}
	if settings.VarDeclFuncToBlock != nil {
		cfg.VarDeclFuncToBlock = *settings.VarDeclFuncToBlock
	}
	if settings.VarDeclPkgToBlock != nil {
		cfg.VarDeclPkgToBlock = *settings.VarDeclPkgToBlock
	}
	if settings.VarDeclInitFuncToBlock != nil {
		cfg.VarDeclInitFuncToBlock = *settings.VarDeclInitFuncToBlock
	}
	if settings.VarDeclInitFuncToShort != nil {
		cfg.VarDeclInitFuncToShort = *settings.VarDeclInitFuncToShort
	}
	if settings.VarDeclInitPkgToBlock != nil {
		cfg.VarDeclInitPkgToBlock = *settings.VarDeclInitPkgToBlock
	}

	return cfg
}
