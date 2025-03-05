package migrate

import (
	"slices"
	"strings"

	"github.com/golangci/golangci-lint/pkg/commands/internal/migrate/one"
	"github.com/golangci/golangci-lint/pkg/commands/internal/migrate/ptr"
	"github.com/golangci/golangci-lint/pkg/commands/internal/migrate/two"
)

func toLinterSettings(old one.LintersSettings) two.LintersSettings {
	return two.LintersSettings{
		Asasalint:       toAsasalintSettings(old.Asasalint),
		BiDiChk:         toBiDiChkSettings(old.BiDiChk),
		CopyLoopVar:     toCopyLoopVarSettings(old.CopyLoopVar),
		Cyclop:          toCyclopSettings(old.Cyclop),
		Decorder:        toDecorderSettings(old.Decorder),
		Depguard:        toDepGuardSettings(old.Depguard),
		Dogsled:         toDogsledSettings(old.Dogsled),
		Dupl:            toDuplSettings(old.Dupl),
		DupWord:         toDupWordSettings(old.DupWord),
		Errcheck:        toErrcheckSettings(old.Errcheck),
		ErrChkJSON:      toErrChkJSONSettings(old.ErrChkJSON),
		ErrorLint:       toErrorLintSettings(old.ErrorLint),
		Exhaustive:      toExhaustiveSettings(old.Exhaustive),
		Exhaustruct:     toExhaustructSettings(old.Exhaustruct),
		Fatcontext:      toFatcontextSettings(old.Fatcontext),
		Forbidigo:       toForbidigoSettings(old.Forbidigo),
		Funlen:          toFunlenSettings(old.Funlen),
		GinkgoLinter:    toGinkgoLinterSettings(old.GinkgoLinter),
		Gocognit:        toGocognitSettings(old.Gocognit),
		GoChecksumType:  toGoChecksumTypeSettings(old.GoChecksumType),
		Goconst:         toGoConstSettings(old.Goconst),
		Gocritic:        toGoCriticSettings(old.Gocritic),
		Gocyclo:         toGoCycloSettings(old.Gocyclo),
		Godot:           toGodotSettings(old.Godot),
		Godox:           toGodoxSettings(old.Godox),
		Goheader:        toGoHeaderSettings(old.Goheader),
		GoModDirectives: toGoModDirectivesSettings(old.GoModDirectives),
		Gomodguard:      toGoModGuardSettings(old.Gomodguard),
		Gosec:           toGoSecSettings(old.Gosec),
		Gosmopolitan:    toGosmopolitanSettings(old.Gosmopolitan),
		Govet:           toGovetSettings(old.Govet),
		Grouper:         toGrouperSettings(old.Grouper),
		Iface:           toIfaceSettings(old.Iface),
		ImportAs:        toImportAsSettings(old.ImportAs),
		Inamedparam:     toINamedParamSettings(old.Inamedparam),
		InterfaceBloat:  toInterfaceBloatSettings(old.InterfaceBloat),
		Ireturn:         toIreturnSettings(old.Ireturn),
		Lll:             toLllSettings(old.Lll),
		LoggerCheck:     toLoggerCheckSettings(old.LoggerCheck),
		MaintIdx:        toMaintIdxSettings(old.MaintIdx),
		Makezero:        toMakezeroSettings(old.Makezero),
		Misspell:        toMisspellSettings(old.Misspell),
		Mnd:             toMndSettings(old.Mnd),
		MustTag:         toMustTagSettings(old.MustTag),
		Nakedret:        toNakedretSettings(old.Nakedret),
		Nestif:          toNestifSettings(old.Nestif),
		NilNil:          toNilNilSettings(old.NilNil),
		Nlreturn:        toNlreturnSettings(old.Nlreturn),
		NoLintLint:      toNoLintLintSettings(old.NoLintLint),
		NoNamedReturns:  toNoNamedReturnsSettings(old.NoNamedReturns),
		ParallelTest:    toParallelTestSettings(old.ParallelTest),
		PerfSprint:      toPerfSprintSettings(old.PerfSprint),
		Prealloc:        toPreallocSettings(old.Prealloc),
		Predeclared:     toPredeclaredSettings(old.Predeclared),
		Promlinter:      toPromlinterSettings(old.Promlinter),
		ProtoGetter:     toProtoGetterSettings(old.ProtoGetter),
		Reassign:        toReassignSettings(old.Reassign),
		Recvcheck:       toRecvcheckSettings(old.Recvcheck),
		Revive:          toReviveSettings(old.Revive),
		RowsErrCheck:    toRowsErrCheckSettings(old.RowsErrCheck),
		SlogLint:        toSlogLintSettings(old.SlogLint),
		Spancheck:       toSpancheckSettings(old.Spancheck),
		Staticcheck:     toStaticCheckSettings(old),
		TagAlign:        toTagAlignSettings(old.TagAlign),
		Tagliatelle:     toTagliatelleSettings(old.Tagliatelle),
		Tenv:            toTenvSettings(old.Tenv),
		Testifylint:     toTestifylintSettings(old.Testifylint),
		Testpackage:     toTestpackageSettings(old.Testpackage),
		Thelper:         toThelperSettings(old.Thelper),
		Unconvert:       toUnconvertSettings(old.Unconvert),
		Unparam:         toUnparamSettings(old.Unparam),
		Unused:          toUnusedSettings(old.Unused),
		UseStdlibVars:   toUseStdlibVarsSettings(old.UseStdlibVars),
		UseTesting:      toUseTestingSettings(old.UseTesting),
		Varnamelen:      toVarnamelenSettings(old.Varnamelen),
		Whitespace:      toWhitespaceSettings(old.Whitespace),
		Wrapcheck:       toWrapcheckSettings(old.Wrapcheck),
		WSL:             toWSLSettings(old.WSL),
		Custom:          toCustom(old.Custom),
	}
}

func toAsasalintSettings(old one.AsasalintSettings) two.AsasalintSettings {
	return two.AsasalintSettings{
		Exclude:              old.Exclude,
		UseBuiltinExclusions: old.UseBuiltinExclusions,
	}
}

func toBiDiChkSettings(old one.BiDiChkSettings) two.BiDiChkSettings {
	// The values are true be default, but the default are defined after the configuration loading.
	// So the serialization doesn't have good results, but it's complex to do better.
	return two.BiDiChkSettings{
		LeftToRightEmbedding:     old.LeftToRightEmbedding,
		RightToLeftEmbedding:     old.RightToLeftEmbedding,
		PopDirectionalFormatting: old.PopDirectionalFormatting,
		LeftToRightOverride:      old.LeftToRightOverride,
		RightToLeftOverride:      old.RightToLeftOverride,
		LeftToRightIsolate:       old.LeftToRightIsolate,
		RightToLeftIsolate:       old.RightToLeftIsolate,
		FirstStrongIsolate:       old.FirstStrongIsolate,
		PopDirectionalIsolate:    old.PopDirectionalIsolate,
	}
}

func toCopyLoopVarSettings(old one.CopyLoopVarSettings) two.CopyLoopVarSettings {
	return two.CopyLoopVarSettings{
		CheckAlias: old.CheckAlias,
	}
}

func toCyclopSettings(old one.Cyclop) two.CyclopSettings {
	return two.CyclopSettings{
		MaxComplexity:  old.MaxComplexity,
		PackageAverage: old.PackageAverage,
	}
}

func toDecorderSettings(old one.DecorderSettings) two.DecorderSettings {
	return two.DecorderSettings{
		DecOrder:                  old.DecOrder,
		IgnoreUnderscoreVars:      old.IgnoreUnderscoreVars,
		DisableDecNumCheck:        old.DisableDecNumCheck,
		DisableTypeDecNumCheck:    old.DisableTypeDecNumCheck,
		DisableConstDecNumCheck:   old.DisableConstDecNumCheck,
		DisableVarDecNumCheck:     old.DisableVarDecNumCheck,
		DisableDecOrderCheck:      old.DisableDecOrderCheck,
		DisableInitFuncFirstCheck: old.DisableInitFuncFirstCheck,
	}
}

func toDepGuardSettings(old one.DepGuardSettings) two.DepGuardSettings {
	settings := two.DepGuardSettings{}

	for k, r := range old.Rules {
		if settings.Rules == nil {
			settings.Rules = make(map[string]*two.DepGuardList)
		}

		list := &two.DepGuardList{
			ListMode: r.ListMode,
			Files:    r.Files,
			Allow:    r.Allow,
		}

		for _, deny := range r.Deny {
			list.Deny = append(list.Deny, two.DepGuardDeny{
				Pkg:  deny.Pkg,
				Desc: deny.Desc,
			})
		}

		settings.Rules[k] = list
	}

	return settings
}

func toDogsledSettings(old one.DogsledSettings) two.DogsledSettings {
	return two.DogsledSettings{
		MaxBlankIdentifiers: old.MaxBlankIdentifiers,
	}
}

func toDuplSettings(old one.DuplSettings) two.DuplSettings {
	return two.DuplSettings{
		Threshold: old.Threshold,
	}
}

func toDupWordSettings(old one.DupWordSettings) two.DupWordSettings {
	return two.DupWordSettings{
		Keywords: old.Keywords,
		Ignore:   old.Ignore,
	}
}

func toErrcheckSettings(old one.ErrcheckSettings) two.ErrcheckSettings {
	return two.ErrcheckSettings{
		DisableDefaultExclusions: old.DisableDefaultExclusions,
		CheckTypeAssertions:      old.CheckTypeAssertions,
		CheckAssignToBlank:       old.CheckAssignToBlank,
		ExcludeFunctions:         old.ExcludeFunctions,
	}
}

func toErrChkJSONSettings(old one.ErrChkJSONSettings) two.ErrChkJSONSettings {
	return two.ErrChkJSONSettings{
		CheckErrorFreeEncoding: old.CheckErrorFreeEncoding,
		ReportNoExported:       old.ReportNoExported,
	}
}

func toErrorLintSettings(old one.ErrorLintSettings) two.ErrorLintSettings {
	settings := two.ErrorLintSettings{
		Errorf:      old.Errorf,
		ErrorfMulti: old.ErrorfMulti,
		Asserts:     old.Asserts,
		Comparison:  old.Comparison,
	}

	for _, allowedError := range old.AllowedErrors {
		settings.AllowedErrors = append(settings.AllowedErrors, two.ErrorLintAllowPair{
			Err: allowedError.Err,
			Fun: allowedError.Fun,
		})
	}
	for _, allowedError := range old.AllowedErrorsWildcard {
		settings.AllowedErrorsWildcard = append(settings.AllowedErrorsWildcard, two.ErrorLintAllowPair{
			Err: allowedError.Err,
			Fun: allowedError.Fun,
		})
	}

	return settings
}

func toExhaustiveSettings(old one.ExhaustiveSettings) two.ExhaustiveSettings {
	return two.ExhaustiveSettings{
		Check:                      old.Check,
		DefaultSignifiesExhaustive: old.DefaultSignifiesExhaustive,
		IgnoreEnumMembers:          old.IgnoreEnumMembers,
		IgnoreEnumTypes:            old.IgnoreEnumTypes,
		PackageScopeOnly:           old.PackageScopeOnly,
		ExplicitExhaustiveMap:      old.ExplicitExhaustiveMap,
		ExplicitExhaustiveSwitch:   old.ExplicitExhaustiveSwitch,
		DefaultCaseRequired:        old.DefaultCaseRequired,
	}
}

func toExhaustructSettings(old one.ExhaustructSettings) two.ExhaustructSettings {
	return two.ExhaustructSettings{
		Include: old.Include,
		Exclude: old.Exclude,
	}
}

func toFatcontextSettings(old one.FatcontextSettings) two.FatcontextSettings {
	return two.FatcontextSettings{
		CheckStructPointers: old.CheckStructPointers,
	}
}

func toForbidigoSettings(old one.ForbidigoSettings) two.ForbidigoSettings {
	settings := two.ForbidigoSettings{
		ExcludeGodocExamples: old.ExcludeGodocExamples,
		AnalyzeTypes:         old.AnalyzeTypes,
	}

	for _, pattern := range old.Forbid {
		if pattern.Pattern == nil {
			buffer, err := pattern.MarshalString()
			if err != nil {
				// impossible case
				panic(err)
			}

			settings.Forbid = append(settings.Forbid, two.ForbidigoPattern{
				Pattern: ptr.Pointer(string(buffer)),
			})

			continue
		}

		settings.Forbid = append(settings.Forbid, two.ForbidigoPattern{
			Pattern: pattern.Pattern,
			Package: pattern.Package,
			Msg:     pattern.Msg,
		})
	}

	return settings
}

func toFunlenSettings(old one.FunlenSettings) two.FunlenSettings {
	return two.FunlenSettings{
		Lines:          old.Lines,
		Statements:     old.Statements,
		IgnoreComments: old.IgnoreComments,
	}
}

func toGinkgoLinterSettings(old one.GinkgoLinterSettings) two.GinkgoLinterSettings {
	return two.GinkgoLinterSettings{
		SuppressLenAssertion:       old.SuppressLenAssertion,
		SuppressNilAssertion:       old.SuppressNilAssertion,
		SuppressErrAssertion:       old.SuppressErrAssertion,
		SuppressCompareAssertion:   old.SuppressCompareAssertion,
		SuppressAsyncAssertion:     old.SuppressAsyncAssertion,
		SuppressTypeCompareWarning: old.SuppressTypeCompareWarning,
		ForbidFocusContainer:       old.ForbidFocusContainer,
		AllowHaveLenZero:           old.AllowHaveLenZero,
		ForceExpectTo:              old.ForceExpectTo,
		ValidateAsyncIntervals:     old.ValidateAsyncIntervals,
		ForbidSpecPollution:        old.ForbidSpecPollution,
		ForceSucceedForFuncs:       old.ForceSucceedForFuncs,
	}
}

func toGocognitSettings(old one.GocognitSettings) two.GocognitSettings {
	return two.GocognitSettings{
		MinComplexity: old.MinComplexity,
	}
}

func toGoChecksumTypeSettings(old one.GoChecksumTypeSettings) two.GoChecksumTypeSettings {
	return two.GoChecksumTypeSettings{
		DefaultSignifiesExhaustive: old.DefaultSignifiesExhaustive,
		IncludeSharedInterfaces:    old.IncludeSharedInterfaces,
	}
}

func toGoConstSettings(old one.GoConstSettings) two.GoConstSettings {
	return two.GoConstSettings{
		IgnoreStrings:       old.IgnoreStrings,
		MatchWithConstants:  old.MatchWithConstants,
		MinStringLen:        old.MinStringLen,
		MinOccurrencesCount: old.MinOccurrencesCount,
		ParseNumbers:        old.ParseNumbers,
		NumberMin:           old.NumberMin,
		NumberMax:           old.NumberMax,
		IgnoreCalls:         old.IgnoreCalls,
	}
}

func toGoCriticSettings(old one.GoCriticSettings) two.GoCriticSettings {
	settings := two.GoCriticSettings{
		Go:             old.Go,
		DisableAll:     old.DisableAll,
		EnabledChecks:  old.EnabledChecks,
		EnableAll:      old.EnableAll,
		DisabledChecks: old.DisabledChecks,
		EnabledTags:    old.EnabledTags,
		DisabledTags:   old.DisabledTags,
	}

	for k, checkSettings := range settings.SettingsPerCheck {
		if settings.SettingsPerCheck == nil {
			settings.SettingsPerCheck = make(map[string]two.GoCriticCheckSettings)
		}

		settings.SettingsPerCheck[k] = checkSettings
	}

	return settings
}

func toGoCycloSettings(old one.GoCycloSettings) two.GoCycloSettings {
	return two.GoCycloSettings{
		MinComplexity: old.MinComplexity,
	}
}

func toGodotSettings(old one.GodotSettings) two.GodotSettings {
	return two.GodotSettings{
		Scope:   old.Scope,
		Exclude: old.Exclude,
		Capital: old.Capital,
		Period:  old.Period,
	}
}

func toGodoxSettings(old one.GodoxSettings) two.GodoxSettings {
	return two.GodoxSettings{
		Keywords: old.Keywords,
	}
}

func toGoHeaderSettings(old one.GoHeaderSettings) two.GoHeaderSettings {
	return two.GoHeaderSettings{
		Values:       old.Values,
		Template:     old.Template,
		TemplatePath: old.TemplatePath,
	}
}

func toGoModDirectivesSettings(old one.GoModDirectivesSettings) two.GoModDirectivesSettings {
	return two.GoModDirectivesSettings{
		ReplaceAllowList:          old.ReplaceAllowList,
		ReplaceLocal:              old.ReplaceLocal,
		ExcludeForbidden:          old.ExcludeForbidden,
		RetractAllowNoExplanation: old.RetractAllowNoExplanation,
		ToolchainForbidden:        old.ToolchainForbidden,
		ToolchainPattern:          old.ToolchainPattern,
		ToolForbidden:             old.ToolForbidden,
		GoDebugForbidden:          old.GoDebugForbidden,
		GoVersionPattern:          old.GoVersionPattern,
	}
}

func toGoModGuardSettings(old one.GoModGuardSettings) two.GoModGuardSettings {
	blocked := two.GoModGuardBlocked{
		LocalReplaceDirectives: old.Blocked.LocalReplaceDirectives,
	}

	for _, version := range old.Blocked.Modules {
		data := map[string]two.GoModGuardModule{}

		for k, v := range version {
			data[k] = two.GoModGuardModule{
				Recommendations: v.Recommendations,
				Reason:          v.Reason,
			}
		}

		blocked.Modules = append(blocked.Modules, data)
	}

	for _, version := range old.Blocked.Versions {
		data := map[string]two.GoModGuardVersion{}

		for k, v := range version {
			data[k] = two.GoModGuardVersion{
				Version: v.Version,
				Reason:  v.Reason,
			}
		}

		blocked.Versions = append(blocked.Versions, data)
	}

	return two.GoModGuardSettings{
		Allowed: two.GoModGuardAllowed{
			Modules: old.Allowed.Modules,
			Domains: old.Allowed.Domains,
		},
		Blocked: blocked,
	}
}

func toGoSecSettings(old one.GoSecSettings) two.GoSecSettings {
	return two.GoSecSettings{
		Includes:    old.Includes,
		Excludes:    old.Excludes,
		Severity:    old.Severity,
		Confidence:  old.Confidence,
		Config:      old.Config,
		Concurrency: old.Concurrency,
	}
}

func toGosmopolitanSettings(old one.GosmopolitanSettings) two.GosmopolitanSettings {
	return two.GosmopolitanSettings{
		AllowTimeLocal:  old.AllowTimeLocal,
		EscapeHatches:   old.EscapeHatches,
		WatchForScripts: old.WatchForScripts,
	}
}

func toGovetSettings(old one.GovetSettings) two.GovetSettings {
	return two.GovetSettings{
		Go:         old.Go,
		Enable:     old.Enable,
		Disable:    old.Disable,
		EnableAll:  old.EnableAll,
		DisableAll: old.DisableAll,
		Settings:   old.Settings,
	}
}

func toGrouperSettings(old one.GrouperSettings) two.GrouperSettings {
	return two.GrouperSettings{
		ConstRequireSingleConst:   old.ConstRequireSingleConst,
		ConstRequireGrouping:      old.ConstRequireGrouping,
		ImportRequireSingleImport: old.ImportRequireSingleImport,
		ImportRequireGrouping:     old.ImportRequireGrouping,
		TypeRequireSingleType:     old.TypeRequireSingleType,
		TypeRequireGrouping:       old.TypeRequireGrouping,
		VarRequireSingleVar:       old.VarRequireSingleVar,
		VarRequireGrouping:        old.VarRequireGrouping,
	}
}

func toIfaceSettings(old one.IfaceSettings) two.IfaceSettings {
	return two.IfaceSettings{
		Enable:   old.Enable,
		Settings: old.Settings,
	}
}

func toImportAsSettings(old one.ImportAsSettings) two.ImportAsSettings {
	settings := two.ImportAsSettings{
		NoUnaliased:    old.NoUnaliased,
		NoExtraAliases: old.NoExtraAliases,
	}

	for _, alias := range old.Alias {
		settings.Alias = append(settings.Alias, two.ImportAsAlias{
			Pkg:   alias.Pkg,
			Alias: alias.Alias,
		})
	}

	return settings
}

func toINamedParamSettings(old one.INamedParamSettings) two.INamedParamSettings {
	return two.INamedParamSettings{
		SkipSingleParam: old.SkipSingleParam,
	}
}

func toInterfaceBloatSettings(old one.InterfaceBloatSettings) two.InterfaceBloatSettings {
	return two.InterfaceBloatSettings{
		Max: old.Max,
	}
}

func toIreturnSettings(old one.IreturnSettings) two.IreturnSettings {
	return two.IreturnSettings{
		Allow:  old.Allow,
		Reject: old.Reject,
	}
}

func toLllSettings(old one.LllSettings) two.LllSettings {
	return two.LllSettings{
		LineLength: old.LineLength,
		TabWidth:   old.TabWidth,
	}
}

func toLoggerCheckSettings(old one.LoggerCheckSettings) two.LoggerCheckSettings {
	return two.LoggerCheckSettings{
		Kitlog:           old.Kitlog,
		Klog:             old.Klog,
		Logr:             old.Logr,
		Slog:             old.Slog,
		Zap:              old.Zap,
		RequireStringKey: old.RequireStringKey,
		NoPrintfLike:     old.NoPrintfLike,
		Rules:            old.Rules,
	}
}

func toMaintIdxSettings(old one.MaintIdxSettings) two.MaintIdxSettings {
	return two.MaintIdxSettings{
		Under: old.Under,
	}
}

func toMakezeroSettings(old one.MakezeroSettings) two.MakezeroSettings {
	return two.MakezeroSettings{
		Always: old.Always,
	}
}

func toMisspellSettings(old one.MisspellSettings) two.MisspellSettings {
	settings := two.MisspellSettings{
		Mode:        old.Mode,
		Locale:      old.Locale,
		IgnoreRules: old.IgnoreWords,
	}

	for _, word := range old.ExtraWords {
		settings.ExtraWords = append(settings.ExtraWords, two.MisspellExtraWords{
			Typo:       word.Typo,
			Correction: word.Correction,
		})
	}

	return settings
}

func toMndSettings(old one.MndSettings) two.MndSettings {
	return two.MndSettings{
		Checks:           old.Checks,
		IgnoredNumbers:   old.IgnoredNumbers,
		IgnoredFiles:     old.IgnoredFiles,
		IgnoredFunctions: old.IgnoredFunctions,
	}
}

func toMustTagSettings(old one.MustTagSettings) two.MustTagSettings {
	settings := two.MustTagSettings{}

	for _, function := range old.Functions {
		settings.Functions = append(settings.Functions, two.MustTagFunction{
			Name:   function.Name,
			Tag:    function.Tag,
			ArgPos: function.ArgPos,
		})
	}

	return settings
}

func toNakedretSettings(old one.NakedretSettings) two.NakedretSettings {
	return two.NakedretSettings{
		MaxFuncLines: old.MaxFuncLines,
	}
}

func toNestifSettings(old one.NestifSettings) two.NestifSettings {
	return two.NestifSettings{
		MinComplexity: old.MinComplexity,
	}
}

func toNilNilSettings(old one.NilNilSettings) two.NilNilSettings {
	return two.NilNilSettings{
		DetectOpposite: old.DetectOpposite,
		CheckedTypes:   old.CheckedTypes,
	}
}

func toNlreturnSettings(old one.NlreturnSettings) two.NlreturnSettings {
	return two.NlreturnSettings{
		BlockSize: old.BlockSize,
	}
}

func toNoLintLintSettings(old one.NoLintLintSettings) two.NoLintLintSettings {
	return two.NoLintLintSettings{
		RequireExplanation: old.RequireExplanation,
		RequireSpecific:    old.RequireSpecific,
		AllowNoExplanation: old.AllowNoExplanation,
		AllowUnused:        old.AllowUnused,
	}
}

func toNoNamedReturnsSettings(old one.NoNamedReturnsSettings) two.NoNamedReturnsSettings {
	return two.NoNamedReturnsSettings{
		ReportErrorInDefer: old.ReportErrorInDefer,
	}
}

func toParallelTestSettings(old one.ParallelTestSettings) two.ParallelTestSettings {
	return two.ParallelTestSettings{
		Go:                    nil,
		IgnoreMissing:         old.IgnoreMissing,
		IgnoreMissingSubtests: old.IgnoreMissingSubtests,
	}
}

func toPerfSprintSettings(old one.PerfSprintSettings) two.PerfSprintSettings {
	return two.PerfSprintSettings{
		IntegerFormat: old.IntegerFormat,
		IntConversion: old.IntConversion,
		ErrorFormat:   old.ErrorFormat,
		ErrError:      old.ErrError,
		ErrorF:        old.ErrorF,
		StringFormat:  old.StringFormat,
		SprintF1:      old.SprintF1,
		StrConcat:     old.StrConcat,
		BoolFormat:    old.BoolFormat,
		HexFormat:     old.HexFormat,
	}
}

func toPreallocSettings(old one.PreallocSettings) two.PreallocSettings {
	return two.PreallocSettings{
		Simple:     old.Simple,
		RangeLoops: old.RangeLoops,
		ForLoops:   old.ForLoops,
	}
}

func toPredeclaredSettings(old one.PredeclaredSettings) two.PredeclaredSettings {
	var ignore []string
	if ptr.Deref(old.Ignore) != "" {
		ignore = strings.Split(ptr.Deref(old.Ignore), ",")
	}

	return two.PredeclaredSettings{
		Ignore:    ignore,
		Qualified: old.Qualified,
	}
}

func toPromlinterSettings(old one.PromlinterSettings) two.PromlinterSettings {
	return two.PromlinterSettings{
		Strict:          old.Strict,
		DisabledLinters: old.DisabledLinters,
	}
}

func toProtoGetterSettings(old one.ProtoGetterSettings) two.ProtoGetterSettings {
	return two.ProtoGetterSettings{
		SkipGeneratedBy:         old.SkipGeneratedBy,
		SkipFiles:               old.SkipFiles,
		SkipAnyGenerated:        old.SkipAnyGenerated,
		ReplaceFirstArgInAppend: old.ReplaceFirstArgInAppend,
	}
}

func toReassignSettings(old one.ReassignSettings) two.ReassignSettings {
	return two.ReassignSettings{
		Patterns: old.Patterns,
	}
}

func toRecvcheckSettings(old one.RecvcheckSettings) two.RecvcheckSettings {
	return two.RecvcheckSettings{
		DisableBuiltin: old.DisableBuiltin,
		Exclusions:     old.Exclusions,
	}
}

func toReviveSettings(old one.ReviveSettings) two.ReviveSettings {
	settings := two.ReviveSettings{
		MaxOpenFiles:   old.MaxOpenFiles,
		Confidence:     old.Confidence,
		Severity:       old.Severity,
		EnableAllRules: old.EnableAllRules,
		ErrorCode:      old.ErrorCode,
		WarningCode:    old.WarningCode,
	}

	for _, rule := range old.Rules {
		settings.Rules = append(settings.Rules, two.ReviveRule{
			Name:      rule.Name,
			Arguments: rule.Arguments,
			Severity:  rule.Severity,
			Disabled:  rule.Disabled,
			Exclude:   rule.Exclude,
		})
	}

	for _, directive := range old.Directives {
		settings.Directives = append(settings.Directives, two.ReviveDirective{
			Name:     directive.Name,
			Severity: directive.Severity,
		})
	}

	return settings
}

func toRowsErrCheckSettings(old one.RowsErrCheckSettings) two.RowsErrCheckSettings {
	return two.RowsErrCheckSettings{
		Packages: old.Packages,
	}
}

func toSlogLintSettings(old one.SlogLintSettings) two.SlogLintSettings {
	return two.SlogLintSettings{
		NoMixedArgs:    old.NoMixedArgs,
		KVOnly:         old.KVOnly,
		AttrOnly:       old.AttrOnly,
		NoGlobal:       old.NoGlobal,
		Context:        old.Context,
		StaticMsg:      old.StaticMsg,
		NoRawKeys:      old.NoRawKeys,
		KeyNamingCase:  old.KeyNamingCase,
		ForbiddenKeys:  old.ForbiddenKeys,
		ArgsOnSepLines: old.ArgsOnSepLines,
	}
}

func toSpancheckSettings(old one.SpancheckSettings) two.SpancheckSettings {
	return two.SpancheckSettings{
		Checks:                   old.Checks,
		IgnoreCheckSignatures:    old.IgnoreCheckSignatures,
		ExtraStartSpanSignatures: old.ExtraStartSpanSignatures,
	}
}

func toStaticCheckSettings(old one.LintersSettings) two.StaticCheckSettings {
	checks := slices.Concat(old.Staticcheck.Checks, old.Stylecheck.Checks, old.Gosimple.Checks)

	slices.Sort(checks)

	return two.StaticCheckSettings{
		Checks:                  slices.Compact(checks),
		Initialisms:             old.Stylecheck.Initialisms,
		DotImportWhitelist:      old.Stylecheck.DotImportWhitelist,
		HTTPStatusCodeWhitelist: old.Stylecheck.HTTPStatusCodeWhitelist,
	}
}

func toTagAlignSettings(old one.TagAlignSettings) two.TagAlignSettings {
	return two.TagAlignSettings{
		Align:  old.Align,
		Sort:   old.Sort,
		Order:  old.Order,
		Strict: old.Strict,
	}
}

func toTagliatelleSettings(old one.TagliatelleSettings) two.TagliatelleSettings {
	tcase := two.TagliatelleCase{
		TagliatelleBase: two.TagliatelleBase{
			Rules:         old.Case.Rules,
			UseFieldName:  old.Case.UseFieldName,
			IgnoredFields: old.Case.IgnoredFields,
		},
		Overrides: []two.TagliatelleOverrides{},
	}

	for k, rule := range old.Case.ExtendedRules {
		if tcase.ExtendedRules == nil {
			tcase.ExtendedRules = make(map[string]two.TagliatelleExtendedRule)
		}

		tcase.ExtendedRules[k] = two.TagliatelleExtendedRule{
			Case:                rule.Case,
			ExtraInitialisms:    rule.ExtraInitialisms,
			InitialismOverrides: rule.InitialismOverrides,
		}
	}

	return two.TagliatelleSettings{Case: tcase}
}

func toTenvSettings(old one.TenvSettings) two.TenvSettings {
	return two.TenvSettings{
		All: old.All,
	}
}

func toTestifylintSettings(old one.TestifylintSettings) two.TestifylintSettings {
	return two.TestifylintSettings{
		EnableAll:        old.EnableAll,
		DisableAll:       old.DisableAll,
		EnabledCheckers:  old.EnabledCheckers,
		DisabledCheckers: old.DisabledCheckers,
		BoolCompare: two.TestifylintBoolCompare{
			IgnoreCustomTypes: old.BoolCompare.IgnoreCustomTypes,
		},
		ExpectedActual: two.TestifylintExpectedActual{
			ExpVarPattern: old.ExpectedActual.ExpVarPattern,
		},
		Formatter: two.TestifylintFormatter{
			CheckFormatString: old.Formatter.CheckFormatString,
			RequireFFuncs:     old.Formatter.RequireFFuncs,
		},
		GoRequire: two.TestifylintGoRequire{
			IgnoreHTTPHandlers: old.GoRequire.IgnoreHTTPHandlers,
		},
		RequireError: two.TestifylintRequireError{
			FnPattern: old.RequireError.FnPattern,
		},
		SuiteExtraAssertCall: two.TestifylintSuiteExtraAssertCall{
			Mode: old.SuiteExtraAssertCall.Mode,
		},
	}
}

func toTestpackageSettings(old one.TestpackageSettings) two.TestpackageSettings {
	return two.TestpackageSettings{
		SkipRegexp:    old.SkipRegexp,
		AllowPackages: old.AllowPackages,
	}
}

func toThelperSettings(old one.ThelperSettings) two.ThelperSettings {
	return two.ThelperSettings{
		Test: two.ThelperOptions{
			First: old.Test.First,
			Name:  old.Test.Name,
			Begin: old.Test.Begin,
		},
		Fuzz: two.ThelperOptions{
			First: old.Fuzz.First,
			Name:  old.Fuzz.Name,
			Begin: old.Fuzz.Begin,
		},
		Benchmark: two.ThelperOptions{
			First: old.Benchmark.First,
			Name:  old.Benchmark.Name,
			Begin: old.Benchmark.Begin,
		},
		TB: two.ThelperOptions{
			First: old.TB.First,
			Name:  old.TB.Name,
			Begin: old.TB.Begin,
		},
	}
}

func toUnconvertSettings(old one.UnconvertSettings) two.UnconvertSettings {
	return two.UnconvertSettings{
		FastMath: old.FastMath,
		Safe:     old.Safe,
	}
}

func toUnparamSettings(old one.UnparamSettings) two.UnparamSettings {
	return two.UnparamSettings{
		CheckExported: old.CheckExported,
	}
}

func toUnusedSettings(old one.UnusedSettings) two.UnusedSettings {
	return two.UnusedSettings{
		FieldWritesAreUses:     old.FieldWritesAreUses,
		PostStatementsAreReads: old.PostStatementsAreReads,
		ExportedFieldsAreUsed:  old.ExportedFieldsAreUsed,
		ParametersAreUsed:      old.ParametersAreUsed,
		LocalVariablesAreUsed:  old.LocalVariablesAreUsed,
		GeneratedIsUsed:        old.GeneratedIsUsed,
	}
}

func toUseStdlibVarsSettings(old one.UseStdlibVarsSettings) two.UseStdlibVarsSettings {
	return two.UseStdlibVarsSettings{
		HTTPMethod:         old.HTTPMethod,
		HTTPStatusCode:     old.HTTPStatusCode,
		TimeWeekday:        old.TimeWeekday,
		TimeMonth:          old.TimeMonth,
		TimeLayout:         old.TimeLayout,
		CryptoHash:         old.CryptoHash,
		DefaultRPCPath:     old.DefaultRPCPath,
		SQLIsolationLevel:  old.SQLIsolationLevel,
		TLSSignatureScheme: old.TLSSignatureScheme,
		ConstantKind:       old.ConstantKind,
	}
}

func toUseTestingSettings(old one.UseTestingSettings) two.UseTestingSettings {
	return two.UseTestingSettings{
		ContextBackground: old.ContextBackground,
		ContextTodo:       old.ContextTodo,
		OSChdir:           old.OSChdir,
		OSMkdirTemp:       old.OSMkdirTemp,
		OSSetenv:          old.OSSetenv,
		OSTempDir:         old.OSTempDir,
		OSCreateTemp:      old.OSCreateTemp,
	}
}

func toVarnamelenSettings(old one.VarnamelenSettings) two.VarnamelenSettings {
	return two.VarnamelenSettings{
		MaxDistance:        old.MaxDistance,
		MinNameLength:      old.MinNameLength,
		CheckReceiver:      old.CheckReceiver,
		CheckReturn:        old.CheckReturn,
		CheckTypeParam:     old.CheckTypeParam,
		IgnoreNames:        old.IgnoreNames,
		IgnoreTypeAssertOk: old.IgnoreTypeAssertOk,
		IgnoreMapIndexOk:   old.IgnoreMapIndexOk,
		IgnoreChanRecvOk:   old.IgnoreChanRecvOk,
		IgnoreDecls:        old.IgnoreDecls,
	}
}

func toWhitespaceSettings(old one.WhitespaceSettings) two.WhitespaceSettings {
	return two.WhitespaceSettings{
		MultiIf:   old.MultiIf,
		MultiFunc: old.MultiFunc,
	}
}

func toWrapcheckSettings(old one.WrapcheckSettings) two.WrapcheckSettings {
	return two.WrapcheckSettings{
		ExtraIgnoreSigs:        old.ExtraIgnoreSigs,
		IgnoreSigs:             old.IgnoreSigs,
		IgnoreSigRegexps:       old.IgnoreSigRegexps,
		IgnorePackageGlobs:     old.IgnorePackageGlobs,
		IgnoreInterfaceRegexps: old.IgnoreInterfaceRegexps,
	}
}

func toWSLSettings(old one.WSLSettings) two.WSLSettings {
	return two.WSLSettings{
		StrictAppend:                     old.StrictAppend,
		AllowAssignAndCallCuddle:         old.AllowAssignAndCallCuddle,
		AllowAssignAndAnythingCuddle:     old.AllowAssignAndAnythingCuddle,
		AllowMultiLineAssignCuddle:       old.AllowMultiLineAssignCuddle,
		ForceCaseTrailingWhitespaceLimit: old.ForceCaseTrailingWhitespaceLimit,
		AllowTrailingComment:             old.AllowTrailingComment,
		AllowSeparatedLeadingComment:     old.AllowSeparatedLeadingComment,
		AllowCuddleDeclaration:           old.AllowCuddleDeclaration,
		AllowCuddleWithCalls:             old.AllowCuddleWithCalls,
		AllowCuddleWithRHS:               old.AllowCuddleWithRHS,
		ForceCuddleErrCheckAndAssign:     old.ForceCuddleErrCheckAndAssign,
		ErrorVariableNames:               old.ErrorVariableNames,
		ForceExclusiveShortDeclarations:  old.ForceExclusiveShortDeclarations,
	}
}

func toCustom(old map[string]one.CustomLinterSettings) map[string]two.CustomLinterSettings {
	if old == nil {
		return nil
	}

	settings := map[string]two.CustomLinterSettings{}

	for k, s := range old {
		settings[k] = two.CustomLinterSettings{
			Type:        s.Type,
			Path:        s.Path,
			Description: s.Description,
			OriginalURL: s.OriginalURL,
			Settings:    s.Settings,
		}
	}

	return settings
}
