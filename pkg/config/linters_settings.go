package config

import (
	"encoding"
	"errors"
	"runtime"

	"gopkg.in/yaml.v3"
)

var defaultLintersSettings = LintersSettings{
	Asasalint: AsasalintSettings{
		UseBuiltinExclusions: true,
	},
	Decorder: DecorderSettings{
		DecOrder:                  []string{"type", "const", "var", "func"},
		DisableDecNumCheck:        true,
		DisableDecOrderCheck:      true,
		DisableInitFuncFirstCheck: true,
	},
	Dogsled: DogsledSettings{
		MaxBlankIdentifiers: 2,
	},
	Dupl: DuplSettings{
		Threshold: 150,
	},
	Errcheck: ErrcheckSettings{
		Ignore: "fmt:.*",
	},
	ErrorLint: ErrorLintSettings{
		Errorf:      true,
		ErrorfMulti: true,
		Asserts:     true,
		Comparison:  true,
	},
	Exhaustive: ExhaustiveSettings{
		Check:                      []string{"switch"},
		CheckGenerated:             false,
		DefaultSignifiesExhaustive: false,
		IgnoreEnumMembers:          "",
		PackageScopeOnly:           false,
		ExplicitExhaustiveMap:      false,
		ExplicitExhaustiveSwitch:   false,
	},
	Forbidigo: ForbidigoSettings{
		ExcludeGodocExamples: true,
	},
	Gci: GciSettings{
		Sections:      []string{"standard", "default"},
		SkipGenerated: true,
	},
	Gocognit: GocognitSettings{
		MinComplexity: 30,
	},
	Goconst: GoConstSettings{
		MatchWithConstants:  true,
		MinStringLen:        3,
		MinOccurrencesCount: 3,
		NumberMin:           3,
		NumberMax:           3,
		IgnoreCalls:         true,
	},
	Gocritic: GoCriticSettings{
		SettingsPerCheck: map[string]GoCriticCheckSettings{},
	},
	Gocyclo: GoCycloSettings{
		MinComplexity: 30,
	},
	Godox: GodoxSettings{
		Keywords: []string{},
	},
	Godot: GodotSettings{
		Scope:  "declarations",
		Period: true,
	},
	Gofmt: GoFmtSettings{
		Simplify: true,
	},
	Gofumpt: GofumptSettings{
		LangVersion: "",
		ModulePath:  "",
		ExtraRules:  false,
	},
	Gosec: GoSecSettings{
		Concurrency: runtime.NumCPU(),
	},
	Gosmopolitan: GosmopolitanSettings{
		AllowTimeLocal:  false,
		EscapeHatches:   []string{},
		IgnoreTests:     true,
		WatchForScripts: []string{"Han"},
	},
	Inamedparam: INamedParamSettings{
		SkipSingleParam: false,
	},
	InterfaceBloat: InterfaceBloatSettings{
		Max: 10,
	},
	Lll: LllSettings{
		LineLength: 120,
		TabWidth:   1,
	},
	LoggerCheck: LoggerCheckSettings{
		Kitlog:           true,
		Klog:             true,
		Logr:             true,
		Zap:              true,
		RequireStringKey: false,
		NoPrintfLike:     false,
		Rules:            nil,
	},
	MaintIdx: MaintIdxSettings{
		Under: 20,
	},
	Nakedret: NakedretSettings{
		MaxFuncLines: 30,
	},
	Nestif: NestifSettings{
		MinComplexity: 5,
	},
	NoLintLint: NoLintLintSettings{
		RequireExplanation: false,
		RequireSpecific:    false,
		AllowUnused:        false,
	},
	PerfSprint: PerfSprintSettings{
		IntConversion: true,
		ErrError:      false,
		ErrorF:        true,
		SprintF1:      true,
		StrConcat:     true,
	},
	Prealloc: PreallocSettings{
		Simple:     true,
		RangeLoops: true,
		ForLoops:   false,
	},
	Predeclared: PredeclaredSettings{
		Ignore:    "",
		Qualified: false,
	},
	SlogLint: SlogLintSettings{
		NoMixedArgs:    true,
		KVOnly:         false,
		AttrOnly:       false,
		ContextOnly:    false,
		StaticMsg:      false,
		NoRawKeys:      false,
		KeyNamingCase:  "",
		ArgsOnSepLines: false,
	},
	TagAlign: TagAlignSettings{
		Align:  true,
		Sort:   true,
		Order:  nil,
		Strict: false,
	},
	Testpackage: TestpackageSettings{
		SkipRegexp:    `(export|internal)_test\.go`,
		AllowPackages: []string{"main"},
	},
	Unparam: UnparamSettings{
		Algo: "cha",
	},
	Unused: UnusedSettings{
		FieldWritesAreUses:     true,
		PostStatementsAreReads: false,
		ExportedIsUsed:         true,
		ExportedFieldsAreUsed:  true,
		ParametersAreUsed:      true,
		LocalVariablesAreUsed:  true,
		GeneratedIsUsed:        true,
	},
	UseStdlibVars: UseStdlibVarsSettings{
		HTTPMethod:     true,
		HTTPStatusCode: true,
	},
	Varnamelen: VarnamelenSettings{
		MaxDistance:   5,
		MinNameLength: 3,
	},
	WSL: WSLSettings{
		StrictAppend:                     true,
		AllowAssignAndCallCuddle:         true,
		AllowAssignAndAnythingCuddle:     false,
		AllowMultiLineAssignCuddle:       true,
		ForceCaseTrailingWhitespaceLimit: 0,
		AllowTrailingComment:             false,
		AllowSeparatedLeadingComment:     false,
		AllowCuddleDeclaration:           false,
		AllowCuddleWithCalls:             []string{"Lock", "RLock"},
		AllowCuddleWithRHS:               []string{"Unlock", "RUnlock"},
		ForceCuddleErrCheckAndAssign:     false,
		ErrorVariableNames:               []string{"err"},
		ForceExclusiveShortDeclarations:  false,
	},
}

type LintersSettings struct {
	Asasalint       AsasalintSettings
	BiDiChk         BiDiChkSettings
	CopyLoopVar     CopyLoopVarSettings
	Cyclop          Cyclop
	Decorder        DecorderSettings
	Depguard        DepGuardSettings
	Dogsled         DogsledSettings
	Dupl            DuplSettings
	DupWord         DupWordSettings
	Errcheck        ErrcheckSettings
	ErrChkJSON      ErrChkJSONSettings
	ErrorLint       ErrorLintSettings
	Exhaustive      ExhaustiveSettings
	Exhaustruct     ExhaustructSettings
	Forbidigo       ForbidigoSettings
	Funlen          FunlenSettings
	Gci             GciSettings
	GinkgoLinter    GinkgoLinterSettings
	Gocognit        GocognitSettings
	Goconst         GoConstSettings
	Gocritic        GoCriticSettings
	Gocyclo         GoCycloSettings
	Godot           GodotSettings
	Godox           GodoxSettings
	Gofmt           GoFmtSettings
	Gofumpt         GofumptSettings
	Goheader        GoHeaderSettings
	Goimports       GoImportsSettings
	Gomnd           GoMndSettings
	GoModDirectives GoModDirectivesSettings
	Gomodguard      GoModGuardSettings
	Gosec           GoSecSettings
	Gosimple        StaticCheckSettings
	Gosmopolitan    GosmopolitanSettings
	Govet           GovetSettings
	Grouper         GrouperSettings
	ImportAs        ImportAsSettings
	Inamedparam     INamedParamSettings
	InterfaceBloat  InterfaceBloatSettings
	Ireturn         IreturnSettings
	Lll             LllSettings
	LoggerCheck     LoggerCheckSettings
	MaintIdx        MaintIdxSettings
	Makezero        MakezeroSettings
	Misspell        MisspellSettings
	MustTag         MustTagSettings
	Nakedret        NakedretSettings
	Nestif          NestifSettings
	NilNil          NilNilSettings
	Nlreturn        NlreturnSettings
	NoLintLint      NoLintLintSettings
	NoNamedReturns  NoNamedReturnsSettings
	ParallelTest    ParallelTestSettings
	PerfSprint      PerfSprintSettings
	Prealloc        PreallocSettings
	Predeclared     PredeclaredSettings
	Promlinter      PromlinterSettings
	ProtoGetter     ProtoGetterSettings
	Reassign        ReassignSettings
	Revive          ReviveSettings
	RowsErrCheck    RowsErrCheckSettings
	SlogLint        SlogLintSettings
	Spancheck       SpancheckSettings
	Staticcheck     StaticCheckSettings
	Stylecheck      StaticCheckSettings
	TagAlign        TagAlignSettings
	Tagliatelle     TagliatelleSettings
	Tenv            TenvSettings
	Testifylint     TestifylintSettings
	Testpackage     TestpackageSettings
	Thelper         ThelperSettings
	Unconvert       UnconvertSettings
	Unparam         UnparamSettings
	Unused          UnusedSettings
	UseStdlibVars   UseStdlibVarsSettings
	Varnamelen      VarnamelenSettings
	Whitespace      WhitespaceSettings
	Wrapcheck       WrapcheckSettings
	WSL             WSLSettings

	Custom map[string]CustomLinterSettings
}

func (s *LintersSettings) Validate() error {
	return s.Govet.Validate()
}

type AsasalintSettings struct {
	Exclude              []string `mapstructure:"exclude"`
	UseBuiltinExclusions bool     `mapstructure:"use-builtin-exclusions"`
	IgnoreTest           bool     `mapstructure:"ignore-test"`
}

type BiDiChkSettings struct {
	LeftToRightEmbedding     bool `mapstructure:"left-to-right-embedding"`
	RightToLeftEmbedding     bool `mapstructure:"right-to-left-embedding"`
	PopDirectionalFormatting bool `mapstructure:"pop-directional-formatting"`
	LeftToRightOverride      bool `mapstructure:"left-to-right-override"`
	RightToLeftOverride      bool `mapstructure:"right-to-left-override"`
	LeftToRightIsolate       bool `mapstructure:"left-to-right-isolate"`
	RightToLeftIsolate       bool `mapstructure:"right-to-left-isolate"`
	FirstStrongIsolate       bool `mapstructure:"first-strong-isolate"`
	PopDirectionalIsolate    bool `mapstructure:"pop-directional-isolate"`
}

type CopyLoopVarSettings struct {
	IgnoreAlias bool `mapstructure:"ignore-alias"`
}

type Cyclop struct {
	MaxComplexity  int     `mapstructure:"max-complexity"`
	PackageAverage float64 `mapstructure:"package-average"`
	SkipTests      bool    `mapstructure:"skip-tests"`
}

type DepGuardSettings struct {
	Rules map[string]*DepGuardList `mapstructure:"rules"`
}

type DepGuardList struct {
	ListMode string         `mapstructure:"list-mode"`
	Files    []string       `mapstructure:"files"`
	Allow    []string       `mapstructure:"allow"`
	Deny     []DepGuardDeny `mapstructure:"deny"`
}

type DepGuardDeny struct {
	Pkg  string `mapstructure:"pkg"`
	Desc string `mapstructure:"desc"`
}

type DecorderSettings struct {
	DecOrder                  []string `mapstructure:"dec-order"`
	IgnoreUnderscoreVars      bool     `mapstructure:"ignore-underscore-vars"`
	DisableDecNumCheck        bool     `mapstructure:"disable-dec-num-check"`
	DisableTypeDecNumCheck    bool     `mapstructure:"disable-type-dec-num-check"`
	DisableConstDecNumCheck   bool     `mapstructure:"disable-const-dec-num-check"`
	DisableVarDecNumCheck     bool     `mapstructure:"disable-var-dec-num-check"`
	DisableDecOrderCheck      bool     `mapstructure:"disable-dec-order-check"`
	DisableInitFuncFirstCheck bool     `mapstructure:"disable-init-func-first-check"`
}

type DogsledSettings struct {
	MaxBlankIdentifiers int `mapstructure:"max-blank-identifiers"`
}

type DuplSettings struct {
	Threshold int
}

type DupWordSettings struct {
	Keywords []string `mapstructure:"keywords"`
	Ignore   []string `mapstructure:"ignore"`
}

type ErrcheckSettings struct {
	DisableDefaultExclusions bool     `mapstructure:"disable-default-exclusions"`
	CheckTypeAssertions      bool     `mapstructure:"check-type-assertions"`
	CheckAssignToBlank       bool     `mapstructure:"check-blank"`
	Ignore                   string   `mapstructure:"ignore"`
	ExcludeFunctions         []string `mapstructure:"exclude-functions"`

	// Deprecated: use ExcludeFunctions instead
	Exclude string `mapstructure:"exclude"`
}

type ErrChkJSONSettings struct {
	CheckErrorFreeEncoding bool `mapstructure:"check-error-free-encoding"`
	ReportNoExported       bool `mapstructure:"report-no-exported"`
}

type ErrorLintSettings struct {
	Errorf      bool `mapstructure:"errorf"`
	ErrorfMulti bool `mapstructure:"errorf-multi"`
	Asserts     bool `mapstructure:"asserts"`
	Comparison  bool `mapstructure:"comparison"`
}

type ExhaustiveSettings struct {
	Check                      []string `mapstructure:"check"`
	CheckGenerated             bool     `mapstructure:"check-generated"`
	DefaultSignifiesExhaustive bool     `mapstructure:"default-signifies-exhaustive"`
	IgnoreEnumMembers          string   `mapstructure:"ignore-enum-members"`
	IgnoreEnumTypes            string   `mapstructure:"ignore-enum-types"`
	PackageScopeOnly           bool     `mapstructure:"package-scope-only"`
	ExplicitExhaustiveMap      bool     `mapstructure:"explicit-exhaustive-map"`
	ExplicitExhaustiveSwitch   bool     `mapstructure:"explicit-exhaustive-switch"`
	DefaultCaseRequired        bool     `mapstructure:"default-case-required"`
}

type ExhaustructSettings struct {
	Include []string `mapstructure:"include"`
	Exclude []string `mapstructure:"exclude"`
}

type ForbidigoSettings struct {
	Forbid               []ForbidigoPattern `mapstructure:"forbid"`
	ExcludeGodocExamples bool               `mapstructure:"exclude-godoc-examples"`
	AnalyzeTypes         bool               `mapstructure:"analyze-types"`
}

var _ encoding.TextUnmarshaler = &ForbidigoPattern{}

// ForbidigoPattern corresponds to forbidigo.pattern and adds mapstructure support.
// The YAML field names must match what forbidigo expects.
type ForbidigoPattern struct {
	// patternString gets populated when the config contains a string as entry in ForbidigoSettings.Forbid[]
	// because ForbidigoPattern implements encoding.TextUnmarshaler
	// and the reader uses the mapstructure.TextUnmarshallerHookFunc as decoder hook.
	//
	// If the entry is a map, then the other fields are set as usual by mapstructure.
	patternString string

	Pattern string `yaml:"p" mapstructure:"p"`
	Package string `yaml:"pkg,omitempty" mapstructure:"pkg,omitempty"`
	Msg     string `yaml:"msg,omitempty" mapstructure:"msg,omitempty"`
}

func (p *ForbidigoPattern) UnmarshalText(text []byte) error {
	// Validation happens when instantiating forbidigo.
	p.patternString = string(text)
	return nil
}

// MarshalString converts the pattern into a string as needed by forbidigo.NewLinter.
//
// MarshalString is intentionally not called MarshalText,
// although it has the same signature
// because implementing encoding.TextMarshaler led to infinite recursion when yaml.Marshal called MarshalText.
func (p *ForbidigoPattern) MarshalString() ([]byte, error) {
	if p.patternString != "" {
		return []byte(p.patternString), nil
	}

	return yaml.Marshal(p)
}

type FunlenSettings struct {
	Lines          int
	Statements     int
	IgnoreComments bool `mapstructure:"ignore-comments"`
}

type GciSettings struct {
	LocalPrefixes string   `mapstructure:"local-prefixes"` // Deprecated
	Sections      []string `mapstructure:"sections"`
	SkipGenerated bool     `mapstructure:"skip-generated"`
	CustomOrder   bool     `mapstructure:"custom-order"`
}

type GinkgoLinterSettings struct {
	SuppressLenAssertion       bool `mapstructure:"suppress-len-assertion"`
	SuppressNilAssertion       bool `mapstructure:"suppress-nil-assertion"`
	SuppressErrAssertion       bool `mapstructure:"suppress-err-assertion"`
	SuppressCompareAssertion   bool `mapstructure:"suppress-compare-assertion"`
	SuppressAsyncAssertion     bool `mapstructure:"suppress-async-assertion"`
	SuppressTypeCompareWarning bool `mapstructure:"suppress-type-compare-assertion"`
	ForbidFocusContainer       bool `mapstructure:"forbid-focus-container"`
	AllowHaveLenZero           bool `mapstructure:"allow-havelen-zero"`
}

type GocognitSettings struct {
	MinComplexity int `mapstructure:"min-complexity"`
}

type GoConstSettings struct {
	IgnoreStrings       string `mapstructure:"ignore-strings"`
	IgnoreTests         bool   `mapstructure:"ignore-tests"`
	MatchWithConstants  bool   `mapstructure:"match-constant"`
	MinStringLen        int    `mapstructure:"min-len"`
	MinOccurrencesCount int    `mapstructure:"min-occurrences"`
	ParseNumbers        bool   `mapstructure:"numbers"`
	NumberMin           int    `mapstructure:"min"`
	NumberMax           int    `mapstructure:"max"`
	IgnoreCalls         bool   `mapstructure:"ignore-calls"`
}

type GoCriticSettings struct {
	Go               string                           `mapstructure:"-"`
	DisableAll       bool                             `mapstructure:"disable-all"`
	EnabledChecks    []string                         `mapstructure:"enabled-checks"`
	EnableAll        bool                             `mapstructure:"enable-all"`
	DisabledChecks   []string                         `mapstructure:"disabled-checks"`
	EnabledTags      []string                         `mapstructure:"enabled-tags"`
	DisabledTags     []string                         `mapstructure:"disabled-tags"`
	SettingsPerCheck map[string]GoCriticCheckSettings `mapstructure:"settings"`
}

type GoCriticCheckSettings map[string]any

type GoCycloSettings struct {
	MinComplexity int `mapstructure:"min-complexity"`
}

type GodotSettings struct {
	Scope   string   `mapstructure:"scope"`
	Exclude []string `mapstructure:"exclude"`
	Capital bool     `mapstructure:"capital"`
	Period  bool     `mapstructure:"period"`

	// Deprecated: use `Scope` instead
	CheckAll bool `mapstructure:"check-all"`
}

type GodoxSettings struct {
	Keywords []string
}

type GoFmtSettings struct {
	Simplify     bool
	RewriteRules []GoFmtRewriteRule `mapstructure:"rewrite-rules"`
}

type GoFmtRewriteRule struct {
	Pattern     string
	Replacement string
}

type GofumptSettings struct {
	ModulePath string `mapstructure:"module-path"`
	ExtraRules bool   `mapstructure:"extra-rules"`

	// Deprecated: use the global `run.go` instead.
	LangVersion string `mapstructure:"lang-version"`
}

type GoHeaderSettings struct {
	Values       map[string]map[string]string `mapstructure:"values"`
	Template     string                       `mapstructure:"template"`
	TemplatePath string                       `mapstructure:"template-path"`
}

type GoImportsSettings struct {
	LocalPrefixes string `mapstructure:"local-prefixes"`
}

type GoMndSettings struct {
	Settings         map[string]map[string]any // Deprecated
	Checks           []string                  `mapstructure:"checks"`
	IgnoredNumbers   []string                  `mapstructure:"ignored-numbers"`
	IgnoredFiles     []string                  `mapstructure:"ignored-files"`
	IgnoredFunctions []string                  `mapstructure:"ignored-functions"`
}

type GoModDirectivesSettings struct {
	ReplaceAllowList          []string `mapstructure:"replace-allow-list"`
	ReplaceLocal              bool     `mapstructure:"replace-local"`
	ExcludeForbidden          bool     `mapstructure:"exclude-forbidden"`
	RetractAllowNoExplanation bool     `mapstructure:"retract-allow-no-explanation"`
}

type GoModGuardSettings struct {
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

type GoSecSettings struct {
	Includes         []string       `mapstructure:"includes"`
	Excludes         []string       `mapstructure:"excludes"`
	Severity         string         `mapstructure:"severity"`
	Confidence       string         `mapstructure:"confidence"`
	ExcludeGenerated bool           `mapstructure:"exclude-generated"`
	Config           map[string]any `mapstructure:"config"`
	Concurrency      int            `mapstructure:"concurrency"`
}

type GosmopolitanSettings struct {
	AllowTimeLocal  bool     `mapstructure:"allow-time-local"`
	EscapeHatches   []string `mapstructure:"escape-hatches"`
	IgnoreTests     bool     `mapstructure:"ignore-tests"`
	WatchForScripts []string `mapstructure:"watch-for-scripts"`
}

type GovetSettings struct {
	Go             string `mapstructure:"-"`
	CheckShadowing bool   `mapstructure:"check-shadowing"`
	Settings       map[string]map[string]any

	Enable     []string
	Disable    []string
	EnableAll  bool `mapstructure:"enable-all"`
	DisableAll bool `mapstructure:"disable-all"`
}

func (cfg *GovetSettings) Validate() error {
	if cfg.EnableAll && cfg.DisableAll {
		return errors.New("govet: enable-all and disable-all can't be combined")
	}
	if cfg.EnableAll && len(cfg.Enable) != 0 {
		return errors.New("govet: enable-all and enable can't be combined")
	}
	if cfg.DisableAll && len(cfg.Disable) != 0 {
		return errors.New("govet: disable-all and disable can't be combined")
	}
	return nil
}

type GrouperSettings struct {
	ConstRequireSingleConst   bool `mapstructure:"const-require-single-const"`
	ConstRequireGrouping      bool `mapstructure:"const-require-grouping"`
	ImportRequireSingleImport bool `mapstructure:"import-require-single-import"`
	ImportRequireGrouping     bool `mapstructure:"import-require-grouping"`
	TypeRequireSingleType     bool `mapstructure:"type-require-single-type"`
	TypeRequireGrouping       bool `mapstructure:"type-require-grouping"`
	VarRequireSingleVar       bool `mapstructure:"var-require-single-var"`
	VarRequireGrouping        bool `mapstructure:"var-require-grouping"`
}

type ImportAsSettings struct {
	Alias          []ImportAsAlias
	NoUnaliased    bool `mapstructure:"no-unaliased"`
	NoExtraAliases bool `mapstructure:"no-extra-aliases"`
}

type ImportAsAlias struct {
	Pkg   string
	Alias string
}

type INamedParamSettings struct {
	SkipSingleParam bool `mapstructure:"skip-single-param"`
}

type InterfaceBloatSettings struct {
	Max int `mapstructure:"max"`
}

type IreturnSettings struct {
	Allow  []string `mapstructure:"allow"`
	Reject []string `mapstructure:"reject"`
}

type LllSettings struct {
	LineLength int `mapstructure:"line-length"`
	TabWidth   int `mapstructure:"tab-width"`
}

type LoggerCheckSettings struct {
	Kitlog           bool     `mapstructure:"kitlog"`
	Klog             bool     `mapstructure:"klog"`
	Logr             bool     `mapstructure:"logr"`
	Zap              bool     `mapstructure:"zap"`
	RequireStringKey bool     `mapstructure:"require-string-key"`
	NoPrintfLike     bool     `mapstructure:"no-printf-like"`
	Rules            []string `mapstructure:"rules"`
}

type MaintIdxSettings struct {
	Under int `mapstructure:"under"`
}

type MakezeroSettings struct {
	Always bool
}

type MisspellSettings struct {
	Mode       string               `mapstructure:"mode"`
	Locale     string               `mapstructure:"locale"`
	ExtraWords []MisspellExtraWords `mapstructure:"extra-words"`
	// TODO(ldez): v2 the option must be renamed to `IgnoredRules`.
	IgnoreWords []string `mapstructure:"ignore-words"`
}

type MisspellExtraWords struct {
	Typo       string `mapstructure:"typo"`
	Correction string `mapstructure:"correction"`
}

type MustTagSettings struct {
	Functions []struct {
		Name   string `mapstructure:"name"`
		Tag    string `mapstructure:"tag"`
		ArgPos int    `mapstructure:"arg-pos"`
	} `mapstructure:"functions"`
}

type NakedretSettings struct {
	MaxFuncLines int `mapstructure:"max-func-lines"`
}

type NestifSettings struct {
	MinComplexity int `mapstructure:"min-complexity"`
}

type NilNilSettings struct {
	CheckedTypes []string `mapstructure:"checked-types"`
}

type NlreturnSettings struct {
	BlockSize int `mapstructure:"block-size"`
}

type NoLintLintSettings struct {
	RequireExplanation bool     `mapstructure:"require-explanation"`
	RequireSpecific    bool     `mapstructure:"require-specific"`
	AllowNoExplanation []string `mapstructure:"allow-no-explanation"`
	AllowUnused        bool     `mapstructure:"allow-unused"`
}

type NoNamedReturnsSettings struct {
	ReportErrorInDefer bool `mapstructure:"report-error-in-defer"`
}

type ParallelTestSettings struct {
	Go                    string `mapstructure:"-"`
	IgnoreMissing         bool   `mapstructure:"ignore-missing"`
	IgnoreMissingSubtests bool   `mapstructure:"ignore-missing-subtests"`
}

type PerfSprintSettings struct {
	IntConversion bool `mapstructure:"int-conversion"`
	ErrError      bool `mapstructure:"err-error"`
	ErrorF        bool `mapstructure:"errorf"`
	SprintF1      bool `mapstructure:"sprintf1"`
	StrConcat     bool `mapstructure:"strconcat"`
}

type PreallocSettings struct {
	Simple     bool
	RangeLoops bool `mapstructure:"range-loops"`
	ForLoops   bool `mapstructure:"for-loops"`
}

type PredeclaredSettings struct {
	Ignore    string `mapstructure:"ignore"`
	Qualified bool   `mapstructure:"q"`
}

type PromlinterSettings struct {
	Strict          bool     `mapstructure:"strict"`
	DisabledLinters []string `mapstructure:"disabled-linters"`
}

type ProtoGetterSettings struct {
	SkipGeneratedBy         []string `mapstructure:"skip-generated-by"`
	SkipFiles               []string `mapstructure:"skip-files"`
	SkipAnyGenerated        bool     `mapstructure:"skip-any-generated"`
	ReplaceFirstArgInAppend bool     `mapstructure:"replace-first-arg-in-append"`
}

type ReassignSettings struct {
	Patterns []string `mapstructure:"patterns"`
}

type ReviveSettings struct {
	MaxOpenFiles          int  `mapstructure:"max-open-files"`
	IgnoreGeneratedHeader bool `mapstructure:"ignore-generated-header"`
	Confidence            float64
	Severity              string
	EnableAllRules        bool `mapstructure:"enable-all-rules"`
	Rules                 []struct {
		Name      string
		Arguments []any
		Severity  string
		Disabled  bool
		Exclude   []string
	}
	ErrorCode   int `mapstructure:"error-code"`
	WarningCode int `mapstructure:"warning-code"`
	Directives  []struct {
		Name     string
		Severity string
	}
}

type RowsErrCheckSettings struct {
	Packages []string
}

type SlogLintSettings struct {
	NoMixedArgs    bool   `mapstructure:"no-mixed-args"`
	KVOnly         bool   `mapstructure:"kv-only"`
	AttrOnly       bool   `mapstructure:"attr-only"`
	ContextOnly    bool   `mapstructure:"context-only"`
	StaticMsg      bool   `mapstructure:"static-msg"`
	NoRawKeys      bool   `mapstructure:"no-raw-keys"`
	KeyNamingCase  string `mapstructure:"key-naming-case"`
	ArgsOnSepLines bool   `mapstructure:"args-on-sep-lines"`
}

type SpancheckSettings struct {
	Checks                []string `mapstructure:"checks"`
	IgnoreCheckSignatures []string `mapstructure:"ignore-check-signatures"`
}

type StaticCheckSettings struct {
	// Deprecated: use the global `run.go` instead.
	GoVersion string `mapstructure:"go"`

	Checks                  []string `mapstructure:"checks"`
	Initialisms             []string `mapstructure:"initialisms"`                // only for stylecheck
	DotImportWhitelist      []string `mapstructure:"dot-import-whitelist"`       // only for stylecheck
	HTTPStatusCodeWhitelist []string `mapstructure:"http-status-code-whitelist"` // only for stylecheck
}

func (s *StaticCheckSettings) HasConfiguration() bool {
	return len(s.Initialisms) > 0 || len(s.HTTPStatusCodeWhitelist) > 0 || len(s.DotImportWhitelist) > 0 || len(s.Checks) > 0
}

type TagAlignSettings struct {
	Align  bool     `mapstructure:"align"`
	Sort   bool     `mapstructure:"sort"`
	Order  []string `mapstructure:"order"`
	Strict bool     `mapstructure:"strict"`
}

type TagliatelleSettings struct {
	Case struct {
		Rules        map[string]string
		UseFieldName bool `mapstructure:"use-field-name"`
	}
}

type TestifylintSettings struct {
	EnableAll        bool     `mapstructure:"enable-all"`
	DisableAll       bool     `mapstructure:"disable-all"`
	EnabledCheckers  []string `mapstructure:"enable"`
	DisabledCheckers []string `mapstructure:"disable"`

	BoolCompare struct {
		IgnoreCustomTypes bool `mapstructure:"ignore-custom-types"`
	} `mapstructure:"bool-compare"`

	ExpectedActual struct {
		ExpVarPattern string `mapstructure:"pattern"`
	} `mapstructure:"expected-actual"`

	RequireError struct {
		FnPattern string `mapstructure:"fn-pattern"`
	} `mapstructure:"require-error"`

	SuiteExtraAssertCall struct {
		Mode string `mapstructure:"mode"`
	} `mapstructure:"suite-extra-assert-call"`
}

type TestpackageSettings struct {
	SkipRegexp    string   `mapstructure:"skip-regexp"`
	AllowPackages []string `mapstructure:"allow-packages"`
}

type ThelperSettings struct {
	Test      ThelperOptions `mapstructure:"test"`
	Fuzz      ThelperOptions `mapstructure:"fuzz"`
	Benchmark ThelperOptions `mapstructure:"benchmark"`
	TB        ThelperOptions `mapstructure:"tb"`
}

type ThelperOptions struct {
	First *bool `mapstructure:"first"`
	Name  *bool `mapstructure:"name"`
	Begin *bool `mapstructure:"begin"`
}

type TenvSettings struct {
	All bool `mapstructure:"all"`
}

type UseStdlibVarsSettings struct {
	HTTPMethod         bool `mapstructure:"http-method"`
	HTTPStatusCode     bool `mapstructure:"http-status-code"`
	TimeWeekday        bool `mapstructure:"time-weekday"`
	TimeMonth          bool `mapstructure:"time-month"`
	TimeLayout         bool `mapstructure:"time-layout"`
	CryptoHash         bool `mapstructure:"crypto-hash"`
	DefaultRPCPath     bool `mapstructure:"default-rpc-path"`
	OSDevNull          bool `mapstructure:"os-dev-null"`
	SQLIsolationLevel  bool `mapstructure:"sql-isolation-level"`
	TLSSignatureScheme bool `mapstructure:"tls-signature-scheme"`
	ConstantKind       bool `mapstructure:"constant-kind"`
	SyslogPriority     bool `mapstructure:"syslog-priority"`
}

type UnconvertSettings struct {
	FastMath bool `mapstructure:"fast-math"`
	Safe     bool `mapstructure:"safe"`
}

type UnparamSettings struct {
	CheckExported bool `mapstructure:"check-exported"`
	Algo          string
}

type UnusedSettings struct {
	FieldWritesAreUses     bool `mapstructure:"field-writes-are-uses"`
	PostStatementsAreReads bool `mapstructure:"post-statements-are-reads"`
	ExportedIsUsed         bool `mapstructure:"exported-is-used"`
	ExportedFieldsAreUsed  bool `mapstructure:"exported-fields-are-used"`
	ParametersAreUsed      bool `mapstructure:"parameters-are-used"`
	LocalVariablesAreUsed  bool `mapstructure:"local-variables-are-used"`
	GeneratedIsUsed        bool `mapstructure:"generated-is-used"`
}

type VarnamelenSettings struct {
	MaxDistance        int      `mapstructure:"max-distance"`
	MinNameLength      int      `mapstructure:"min-name-length"`
	CheckReceiver      bool     `mapstructure:"check-receiver"`
	CheckReturn        bool     `mapstructure:"check-return"`
	CheckTypeParam     bool     `mapstructure:"check-type-param"`
	IgnoreNames        []string `mapstructure:"ignore-names"`
	IgnoreTypeAssertOk bool     `mapstructure:"ignore-type-assert-ok"`
	IgnoreMapIndexOk   bool     `mapstructure:"ignore-map-index-ok"`
	IgnoreChanRecvOk   bool     `mapstructure:"ignore-chan-recv-ok"`
	IgnoreDecls        []string `mapstructure:"ignore-decls"`
}

type WhitespaceSettings struct {
	MultiIf   bool `mapstructure:"multi-if"`
	MultiFunc bool `mapstructure:"multi-func"`
}

type WrapcheckSettings struct {
	// TODO(ldez): v2 the options must be renamed to use hyphen.
	IgnoreSigs             []string `mapstructure:"ignoreSigs"`
	IgnoreSigRegexps       []string `mapstructure:"ignoreSigRegexps"`
	IgnorePackageGlobs     []string `mapstructure:"ignorePackageGlobs"`
	IgnoreInterfaceRegexps []string `mapstructure:"ignoreInterfaceRegexps"`
}

type WSLSettings struct {
	StrictAppend                     bool     `mapstructure:"strict-append"`
	AllowAssignAndCallCuddle         bool     `mapstructure:"allow-assign-and-call"`
	AllowAssignAndAnythingCuddle     bool     `mapstructure:"allow-assign-and-anything"`
	AllowMultiLineAssignCuddle       bool     `mapstructure:"allow-multiline-assign"`
	ForceCaseTrailingWhitespaceLimit int      `mapstructure:"force-case-trailing-whitespace"`
	AllowTrailingComment             bool     `mapstructure:"allow-trailing-comment"`
	AllowSeparatedLeadingComment     bool     `mapstructure:"allow-separated-leading-comment"`
	AllowCuddleDeclaration           bool     `mapstructure:"allow-cuddle-declarations"`
	AllowCuddleWithCalls             []string `mapstructure:"allow-cuddle-with-calls"`
	AllowCuddleWithRHS               []string `mapstructure:"allow-cuddle-with-rhs"`
	ForceCuddleErrCheckAndAssign     bool     `mapstructure:"force-err-cuddling"`
	ErrorVariableNames               []string `mapstructure:"error-variable-names"`
	ForceExclusiveShortDeclarations  bool     `mapstructure:"force-short-decl-cuddling"`
}

// CustomLinterSettings encapsulates the meta-data of a private linter.
// For example, a private linter may be added to the golangci config file as shown below.
//
//	linters-settings:
//	 custom:
//	   example:
//	     path: /example.so
//	     description: The description of the linter
//	     original-url: github.com/golangci/example-linter
type CustomLinterSettings struct {
	// Path to a plugin *.so file that implements the private linter.
	Path string
	// Description describes the purpose of the private linter.
	Description string
	// OriginalURL The URL containing the source code for the private linter.
	OriginalURL string `mapstructure:"original-url"`

	// Settings plugin settings only work with linterdb.PluginConstructor symbol.
	Settings any
}
