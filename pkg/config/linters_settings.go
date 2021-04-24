package config

import "github.com/pkg/errors"

var defaultLintersSettings = LintersSettings{
	Lll: LllSettings{
		LineLength: 120,
		TabWidth:   1,
	},
	Unparam: UnparamSettings{
		Algo: "cha",
	},
	Nakedret: NakedretSettings{
		MaxFuncLines: 30,
	},
	Prealloc: PreallocSettings{
		Simple:     true,
		RangeLoops: true,
		ForLoops:   false,
	},
	Gocritic: GocriticSettings{
		SettingsPerCheck: map[string]GocriticCheckSettings{},
	},
	Godox: GodoxSettings{
		Keywords: []string{},
	},
	Dogsled: DogsledSettings{
		MaxBlankIdentifiers: 2,
	},
	Gocognit: GocognitSettings{
		MinComplexity: 30,
	},
	WSL: WSLSettings{
		StrictAppend:                     true,
		AllowAssignAndCallCuddle:         true,
		AllowAssignAndAnythingCuddle:     false,
		AllowMultiLineAssignCuddle:       true,
		AllowCuddleDeclaration:           false,
		AllowTrailingComment:             false,
		AllowSeparatedLeadingComment:     false,
		ForceCuddleErrCheckAndAssign:     false,
		ForceExclusiveShortDeclarations:  false,
		ForceCaseTrailingWhitespaceLimit: 0,
	},
	NoLintLint: NoLintLintSettings{
		RequireExplanation: false,
		AllowLeadingSpace:  true,
		RequireSpecific:    false,
		AllowUnused:        false,
	},
	Testpackage: TestpackageSettings{
		SkipRegexp: `(export|internal)_test\.go`,
	},
	Nestif: NestifSettings{
		MinComplexity: 5,
	},
	Exhaustive: ExhaustiveSettings{
		CheckGenerated:             false,
		DefaultSignifiesExhaustive: false,
	},
	Gofumpt: GofumptSettings{
		ExtraRules: false,
	},
	ErrorLint: ErrorLintSettings{
		Errorf: true,
	},
	Ifshort: IfshortSettings{
		MaxDeclLines: 1,
		MaxDeclChars: 30,
	},
	Predeclared: PredeclaredSettings{
		Ignore:    "",
		Qualified: false,
	},
	Forbidigo: ForbidigoSettings{
		ExcludeGodocExamples: true,
	},
}

type LintersSettings struct {
	Gci struct {
		LocalPrefixes string `mapstructure:"local-prefixes"`
	}
	Govet  GovetSettings
	Golint struct {
		MinConfidence float64 `mapstructure:"min-confidence"`
	}
	Gofmt struct {
		Simplify bool
	}
	Goimports struct {
		LocalPrefixes string `mapstructure:"local-prefixes"`
	}
	Gocyclo struct {
		MinComplexity int `mapstructure:"min-complexity"`
	}
	Varcheck struct {
		CheckExportedFields bool `mapstructure:"exported-fields"`
	}
	Structcheck struct {
		CheckExportedFields bool `mapstructure:"exported-fields"`
	}
	Maligned struct {
		SuggestNewOrder bool `mapstructure:"suggest-new"`
	}
	Dupl struct {
		Threshold int
	}
	Goconst struct {
		MatchWithConstants  bool `mapstructure:"match-constant"`
		MinStringLen        int  `mapstructure:"min-len"`
		MinOccurrencesCount int  `mapstructure:"min-occurrences"`
		ParseNumbers        bool `mapstructure:"numbers"`
		NumberMin           int  `mapstructure:"min"`
		NumberMax           int  `mapstructure:"max"`
		IgnoreCalls         bool `mapstructure:"ignore-calls"`
	}
	Gomnd struct {
		Settings map[string]map[string]interface{}
	}
	Depguard struct {
		ListType                 string `mapstructure:"list-type"`
		Packages                 []string
		IncludeGoRoot            bool              `mapstructure:"include-go-root"`
		PackagesWithErrorMessage map[string]string `mapstructure:"packages-with-error-message"`
	}
	Misspell struct {
		Locale      string
		IgnoreWords []string `mapstructure:"ignore-words"`
	}
	Unused struct {
		CheckExported bool `mapstructure:"check-exported"`
	}
	Funlen struct {
		Lines      int
		Statements int
	}
	Whitespace struct {
		MultiIf   bool `mapstructure:"multi-if"`
		MultiFunc bool `mapstructure:"multi-func"`
	}
	RowsErrCheck struct {
		Packages []string
	}
	Gomodguard struct {
		Allowed struct {
			Modules []string `mapstructure:"modules"`
			Domains []string `mapstructure:"domains"`
		} `mapstructure:"allowed"`
		Blocked struct {
			Modules []map[string]struct {
				Recommendations []string `mapstructure:"recommendations"`
				Reason          string   `mapstructure:"reason"`
			} `mapstructure:"modules"`
			Versions []map[string]struct {
				Version string `mapstructure:"version"`
				Reason  string `mapstructure:"reason"`
			} `mapstructure:"versions"`
			LocalReplaceDirectives bool `mapstructure:"local_replace_directives"`
		} `mapstructure:"blocked"`
	}

	WSL              WSLSettings
	Lll              LllSettings
	Unparam          UnparamSettings
	Nakedret         NakedretSettings
	Prealloc         PreallocSettings
	Errcheck         ErrcheckSettings
	Gocritic         GocriticSettings
	Godox            GodoxSettings
	Dogsled          DogsledSettings
	Gocognit         GocognitSettings
	Godot            GodotSettings
	Goheader         GoHeaderSettings
	Testpackage      TestpackageSettings
	Nestif           NestifSettings
	NoLintLint       NoLintLintSettings
	Exhaustive       ExhaustiveSettings
	ExhaustiveStruct ExhaustiveStructSettings
	Gofumpt          GofumptSettings
	ErrorLint        ErrorLintSettings
	Makezero         MakezeroSettings
	Revive           ReviveSettings
	Thelper          ThelperSettings
	Forbidigo        ForbidigoSettings
	Ifshort          IfshortSettings
	Predeclared      PredeclaredSettings
	Cyclop           Cyclop
	ImportAs         ImportAsSettings
	GoModDirectives  GoModDirectivesSettings
	Promlinter       PromlinterSettings
	Tagliatelle      TagliatelleSettings

	Custom map[string]CustomLinterSettings
}

type GoHeaderSettings struct {
	Values       map[string]map[string]string `mapstructure:"values"`
	Template     string                       `mapstructure:"template"`
	TemplatePath string                       `mapstructure:"template-path"`
}

type GovetSettings struct {
	CheckShadowing bool `mapstructure:"check-shadowing"`
	Settings       map[string]map[string]interface{}

	Enable     []string
	Disable    []string
	EnableAll  bool `mapstructure:"enable-all"`
	DisableAll bool `mapstructure:"disable-all"`
}

func (cfg GovetSettings) Validate() error {
	if cfg.EnableAll && cfg.DisableAll {
		return errors.New("enable-all and disable-all can't be combined")
	}
	if cfg.EnableAll && len(cfg.Enable) != 0 {
		return errors.New("enable-all and enable can't be combined")
	}
	if cfg.DisableAll && len(cfg.Disable) != 0 {
		return errors.New("disable-all and disable can't be combined")
	}
	return nil
}

type ErrcheckSettings struct {
	CheckTypeAssertions bool   `mapstructure:"check-type-assertions"`
	CheckAssignToBlank  bool   `mapstructure:"check-blank"`
	Ignore              string `mapstructure:"ignore"`
	Exclude             string `mapstructure:"exclude"`
}

type LllSettings struct {
	LineLength int `mapstructure:"line-length"`
	TabWidth   int `mapstructure:"tab-width"`
}

type UnparamSettings struct {
	CheckExported bool `mapstructure:"check-exported"`
	Algo          string
}

type NakedretSettings struct {
	MaxFuncLines int `mapstructure:"max-func-lines"`
}

type PreallocSettings struct {
	Simple     bool
	RangeLoops bool `mapstructure:"range-loops"`
	ForLoops   bool `mapstructure:"for-loops"`
}

type GodoxSettings struct {
	Keywords []string
}

type DogsledSettings struct {
	MaxBlankIdentifiers int `mapstructure:"max-blank-identifiers"`
}

type GocognitSettings struct {
	MinComplexity int `mapstructure:"min-complexity"`
}

type WSLSettings struct {
	StrictAppend                     bool `mapstructure:"strict-append"`
	AllowAssignAndCallCuddle         bool `mapstructure:"allow-assign-and-call"`
	AllowAssignAndAnythingCuddle     bool `mapstructure:"allow-assign-and-anything"`
	AllowMultiLineAssignCuddle       bool `mapstructure:"allow-multiline-assign"`
	AllowCuddleDeclaration           bool `mapstructure:"allow-cuddle-declarations"`
	AllowTrailingComment             bool `mapstructure:"allow-trailing-comment"`
	AllowSeparatedLeadingComment     bool `mapstructure:"allow-separated-leading-comment"`
	ForceCuddleErrCheckAndAssign     bool `mapstructure:"force-err-cuddling"`
	ForceExclusiveShortDeclarations  bool `mapstructure:"force-short-decl-cuddling"`
	ForceCaseTrailingWhitespaceLimit int  `mapstructure:"force-case-trailing-whitespace"`
}

type GodotSettings struct {
	Scope   string   `mapstructure:"scope"`
	Exclude []string `mapstructure:"exclude"`
	Capital bool     `mapstructure:"capital"`

	// Deprecated: use `Scope` instead
	CheckAll bool `mapstructure:"check-all"`
}

type NoLintLintSettings struct {
	RequireExplanation bool     `mapstructure:"require-explanation"`
	AllowLeadingSpace  bool     `mapstructure:"allow-leading-space"`
	RequireSpecific    bool     `mapstructure:"require-specific"`
	AllowNoExplanation []string `mapstructure:"allow-no-explanation"`
	AllowUnused        bool     `mapstructure:"allow-unused"`
}

type TestpackageSettings struct {
	SkipRegexp string `mapstructure:"skip-regexp"`
}

type NestifSettings struct {
	MinComplexity int `mapstructure:"min-complexity"`
}

type ExhaustiveSettings struct {
	CheckGenerated             bool `mapstructure:"check-generated"`
	DefaultSignifiesExhaustive bool `mapstructure:"default-signifies-exhaustive"`
}

type ExhaustiveStructSettings struct {
	StructPatterns []string `mapstructure:"struct-patterns"`
}

type GofumptSettings struct {
	ExtraRules bool `mapstructure:"extra-rules"`
}

type ErrorLintSettings struct {
	Errorf bool `mapstructure:"errorf"`
}

type MakezeroSettings struct {
	Always bool
}

type ReviveSettings struct {
	IgnoreGeneratedHeader bool `mapstructure:"ignore-generated-header"`
	Confidence            float64
	Severity              string
	Rules                 []struct {
		Name      string
		Arguments []interface{}
		Severity  string
	}
	ErrorCode   int `mapstructure:"error-code"`
	WarningCode int `mapstructure:"warning-code"`
	Directives  []struct {
		Name     string
		Severity string
	}
}

type ThelperSettings struct {
	Test struct {
		First bool `mapstructure:"first"`
		Name  bool `mapstructure:"name"`
		Begin bool `mapstructure:"begin"`
	} `mapstructure:"test"`
	Benchmark struct {
		First bool `mapstructure:"first"`
		Name  bool `mapstructure:"name"`
		Begin bool `mapstructure:"begin"`
	} `mapstructure:"benchmark"`
	TB struct {
		First bool `mapstructure:"first"`
		Name  bool `mapstructure:"name"`
		Begin bool `mapstructure:"begin"`
	} `mapstructure:"tb"`
}

type IfshortSettings struct {
	MaxDeclLines int `mapstructure:"max-decl-lines"`
	MaxDeclChars int `mapstructure:"max-decl-chars"`
}

type ForbidigoSettings struct {
	Forbid               []string `mapstructure:"forbid"`
	ExcludeGodocExamples bool     `mapstructure:"exclude-godoc-examples"`
}

type PredeclaredSettings struct {
	Ignore    string `mapstructure:"ignore"`
	Qualified bool   `mapstructure:"q"`
}

type PromlinterSettings struct {
	Strict          bool     `mapstructure:"strict"`
	DisabledLinters []string `mapstructure:"disabled-linters"`
}

type Cyclop struct {
	MaxComplexity  int     `mapstructure:"max-complexity"`
	PackageAverage float64 `mapstructure:"package-average"`
	SkipTests      bool    `mapstructure:"skip-tests"`
}

type ImportAsSettings map[string]string

type GoModDirectivesSettings struct {
	ReplaceAllowList          []string `mapstructure:"replace-allow-list"`
	ReplaceLocal              bool     `mapstructure:"replace-local"`
	ExcludeForbidden          bool     `mapstructure:"exclude-forbidden"`
	RetractAllowNoExplanation bool     `mapstructure:"retract-allow-no-explanation"`
}

type TagliatelleSettings struct {
	Case struct {
		Rules        map[string]string
		UseFieldName bool `mapstructure:"use-field-name"`
	}
}

type CustomLinterSettings struct {
	Path        string
	Description string
	OriginalURL string `mapstructure:"original-url"`
}
