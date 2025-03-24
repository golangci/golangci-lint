package config

import (
	"errors"
	"fmt"
	"runtime"
)

var defaultLintersSettings = LintersSettings{
	FormatterSettings: defaultFormatterSettings,

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
	ErrorLint: ErrorLintSettings{
		Errorf:      true,
		ErrorfMulti: true,
		Asserts:     true,
		Comparison:  true,
	},
	Exhaustive: ExhaustiveSettings{
		Check:                      []string{"switch"},
		DefaultSignifiesExhaustive: false,
		IgnoreEnumMembers:          "",
		PackageScopeOnly:           false,
		ExplicitExhaustiveMap:      false,
		ExplicitExhaustiveSwitch:   false,
	},
	Forbidigo: ForbidigoSettings{
		ExcludeGodocExamples: true,
	},
	Funlen: FunlenSettings{
		IgnoreComments: true,
	},
	GoChecksumType: GoChecksumTypeSettings{
		DefaultSignifiesExhaustive: true,
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
	Gosec: GoSecSettings{
		Concurrency: runtime.NumCPU(),
	},
	Gosmopolitan: GosmopolitanSettings{
		AllowTimeLocal:  false,
		EscapeHatches:   []string{},
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
		Slog:             true,
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
		IntegerFormat: true,
		IntConversion: true,
		ErrorFormat:   true,
		ErrError:      false,
		ErrorF:        true,
		StringFormat:  true,
		SprintF1:      true,
		StrConcat:     true,
		BoolFormat:    true,
		HexFormat:     true,
	},
	Prealloc: PreallocSettings{
		Simple:     true,
		RangeLoops: true,
		ForLoops:   false,
	},
	Predeclared: PredeclaredSettings{
		Qualified: false,
	},
	SlogLint: SlogLintSettings{
		NoMixedArgs:    true,
		KVOnly:         false,
		AttrOnly:       false,
		NoGlobal:       "",
		Context:        "",
		StaticMsg:      false,
		NoRawKeys:      false,
		KeyNamingCase:  "",
		ForbiddenKeys:  nil,
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
	Unused: UnusedSettings{
		FieldWritesAreUses:     true,
		PostStatementsAreReads: false,
		ExportedFieldsAreUsed:  true,
		ParametersAreUsed:      true,
		LocalVariablesAreUsed:  true,
		GeneratedIsUsed:        true,
	},
	UseStdlibVars: UseStdlibVarsSettings{
		HTTPMethod:     true,
		HTTPStatusCode: true,
	},
	UseTesting: UseTestingSettings{
		ContextBackground: true,
		ContextTodo:       true,
		OSChdir:           true,
		OSMkdirTemp:       true,
		OSSetenv:          true,
		OSTempDir:         false,
		OSCreateTemp:      true,
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
		AllowCuddleUsedInBlock:           false,
		ForceCuddleErrCheckAndAssign:     false,
		ErrorVariableNames:               []string{"err"},
		ForceExclusiveShortDeclarations:  false,
	},
}

type LintersSettings struct {
	FormatterSettings `mapstructure:"-"`

	Asasalint       AsasalintSettings       `mapstructure:"asasalint"`
	BiDiChk         BiDiChkSettings         `mapstructure:"bidichk"`
	CopyLoopVar     CopyLoopVarSettings     `mapstructure:"copyloopvar"`
	Cyclop          CyclopSettings          `mapstructure:"cyclop"`
	Decorder        DecorderSettings        `mapstructure:"decorder"`
	Depguard        DepGuardSettings        `mapstructure:"depguard"`
	Dogsled         DogsledSettings         `mapstructure:"dogsled"`
	Dupl            DuplSettings            `mapstructure:"dupl"`
	DupWord         DupWordSettings         `mapstructure:"dupword"`
	Errcheck        ErrcheckSettings        `mapstructure:"errcheck"`
	ErrChkJSON      ErrChkJSONSettings      `mapstructure:"errchkjson"`
	ErrorLint       ErrorLintSettings       `mapstructure:"errorlint"`
	Exhaustive      ExhaustiveSettings      `mapstructure:"exhaustive"`
	Exhaustruct     ExhaustructSettings     `mapstructure:"exhaustruct"`
	Fatcontext      FatcontextSettings      `mapstructure:"fatcontext"`
	Forbidigo       ForbidigoSettings       `mapstructure:"forbidigo"`
	Funlen          FunlenSettings          `mapstructure:"funlen"`
	GinkgoLinter    GinkgoLinterSettings    `mapstructure:"ginkgolinter"`
	Gocognit        GocognitSettings        `mapstructure:"gocognit"`
	GoChecksumType  GoChecksumTypeSettings  `mapstructure:"gochecksumtype"`
	Goconst         GoConstSettings         `mapstructure:"goconst"`
	Gocritic        GoCriticSettings        `mapstructure:"gocritic"`
	Gocyclo         GoCycloSettings         `mapstructure:"gocyclo"`
	Godot           GodotSettings           `mapstructure:"godot"`
	Godox           GodoxSettings           `mapstructure:"godox"`
	Goheader        GoHeaderSettings        `mapstructure:"goheader"`
	GoModDirectives GoModDirectivesSettings `mapstructure:"gomoddirectives"`
	Gomodguard      GoModGuardSettings      `mapstructure:"gomodguard"`
	Gosec           GoSecSettings           `mapstructure:"gosec"`
	Gosmopolitan    GosmopolitanSettings    `mapstructure:"gosmopolitan"`
	Govet           GovetSettings           `mapstructure:"govet"`
	Grouper         GrouperSettings         `mapstructure:"grouper"`
	Iface           IfaceSettings           `mapstructure:"iface"`
	ImportAs        ImportAsSettings        `mapstructure:"importas"`
	Inamedparam     INamedParamSettings     `mapstructure:"inamedparam"`
	InterfaceBloat  InterfaceBloatSettings  `mapstructure:"interfacebloat"`
	Ireturn         IreturnSettings         `mapstructure:"ireturn"`
	Lll             LllSettings             `mapstructure:"lll"`
	LoggerCheck     LoggerCheckSettings     `mapstructure:"loggercheck"`
	MaintIdx        MaintIdxSettings        `mapstructure:"maintidx"`
	Makezero        MakezeroSettings        `mapstructure:"makezero"`
	Misspell        MisspellSettings        `mapstructure:"misspell"`
	Mnd             MndSettings             `mapstructure:"mnd"`
	MustTag         MustTagSettings         `mapstructure:"musttag"`
	Nakedret        NakedretSettings        `mapstructure:"nakedret"`
	Nestif          NestifSettings          `mapstructure:"nestif"`
	NilNil          NilNilSettings          `mapstructure:"nilnil"`
	Nlreturn        NlreturnSettings        `mapstructure:"nlreturn"`
	NoLintLint      NoLintLintSettings      `mapstructure:"nolintlint"`
	NoNamedReturns  NoNamedReturnsSettings  `mapstructure:"nonamedreturns"`
	ParallelTest    ParallelTestSettings    `mapstructure:"paralleltest"`
	PerfSprint      PerfSprintSettings      `mapstructure:"perfsprint"`
	Prealloc        PreallocSettings        `mapstructure:"prealloc"`
	Predeclared     PredeclaredSettings     `mapstructure:"predeclared"`
	Promlinter      PromlinterSettings      `mapstructure:"promlinter"`
	ProtoGetter     ProtoGetterSettings     `mapstructure:"protogetter"`
	Reassign        ReassignSettings        `mapstructure:"reassign"`
	Recvcheck       RecvcheckSettings       `mapstructure:"recvcheck"`
	Revive          ReviveSettings          `mapstructure:"revive"`
	RowsErrCheck    RowsErrCheckSettings    `mapstructure:"rowserrcheck"`
	SlogLint        SlogLintSettings        `mapstructure:"sloglint"`
	Spancheck       SpancheckSettings       `mapstructure:"spancheck"`
	Staticcheck     StaticCheckSettings     `mapstructure:"staticcheck"`
	TagAlign        TagAlignSettings        `mapstructure:"tagalign"`
	Tagliatelle     TagliatelleSettings     `mapstructure:"tagliatelle"`
	Testifylint     TestifylintSettings     `mapstructure:"testifylint"`
	Testpackage     TestpackageSettings     `mapstructure:"testpackage"`
	Thelper         ThelperSettings         `mapstructure:"thelper"`
	Unconvert       UnconvertSettings       `mapstructure:"unconvert"`
	Unparam         UnparamSettings         `mapstructure:"unparam"`
	Unused          UnusedSettings          `mapstructure:"unused"`
	UseStdlibVars   UseStdlibVarsSettings   `mapstructure:"usestdlibvars"`
	UseTesting      UseTestingSettings      `mapstructure:"usetesting"`
	Varnamelen      VarnamelenSettings      `mapstructure:"varnamelen"`
	Whitespace      WhitespaceSettings      `mapstructure:"whitespace"`
	Wrapcheck       WrapcheckSettings       `mapstructure:"wrapcheck"`
	WSL             WSLSettings             `mapstructure:"wsl"`

	Custom map[string]CustomLinterSettings `mapstructure:"custom"`
}

func (s *LintersSettings) Validate() error {
	if err := s.Govet.Validate(); err != nil {
		return err
	}

	for name, settings := range s.Custom {
		if err := settings.Validate(); err != nil {
			return fmt.Errorf("custom linter %q: %w", name, err)
		}
	}

	return nil
}

type AsasalintSettings struct {
	Exclude              []string `mapstructure:"exclude"`
	UseBuiltinExclusions bool     `mapstructure:"use-builtin-exclusions"`
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
	CheckAlias bool `mapstructure:"check-alias"`
}

type CyclopSettings struct {
	MaxComplexity  int     `mapstructure:"max-complexity"`
	PackageAverage float64 `mapstructure:"package-average"`
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
	Threshold int `mapstructure:"threshold"`
}

type DupWordSettings struct {
	Keywords []string `mapstructure:"keywords"`
	Ignore   []string `mapstructure:"ignore"`
}

type ErrcheckSettings struct {
	DisableDefaultExclusions bool     `mapstructure:"disable-default-exclusions"`
	CheckTypeAssertions      bool     `mapstructure:"check-type-assertions"`
	CheckAssignToBlank       bool     `mapstructure:"check-blank"`
	ExcludeFunctions         []string `mapstructure:"exclude-functions"`
}

type ErrChkJSONSettings struct {
	CheckErrorFreeEncoding bool `mapstructure:"check-error-free-encoding"`
	ReportNoExported       bool `mapstructure:"report-no-exported"`
}

type ErrorLintSettings struct {
	Errorf                bool                 `mapstructure:"errorf"`
	ErrorfMulti           bool                 `mapstructure:"errorf-multi"`
	Asserts               bool                 `mapstructure:"asserts"`
	Comparison            bool                 `mapstructure:"comparison"`
	AllowedErrors         []ErrorLintAllowPair `mapstructure:"allowed-errors"`
	AllowedErrorsWildcard []ErrorLintAllowPair `mapstructure:"allowed-errors-wildcard"`
}

type ErrorLintAllowPair struct {
	Err string `mapstructure:"err"`
	Fun string `mapstructure:"fun"`
}

type ExhaustiveSettings struct {
	Check                      []string `mapstructure:"check"`
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

type FatcontextSettings struct {
	CheckStructPointers bool `mapstructure:"check-struct-pointers"`
}

type ForbidigoSettings struct {
	Forbid               []ForbidigoPattern `mapstructure:"forbid"`
	ExcludeGodocExamples bool               `mapstructure:"exclude-godoc-examples"`
	AnalyzeTypes         bool               `mapstructure:"analyze-types"`
}

type ForbidigoPattern struct {
	Pattern string `yaml:"p" mapstructure:"pattern"`
	Package string `yaml:"pkg,omitempty" mapstructure:"pkg,omitempty"`
	Msg     string `yaml:"msg,omitempty" mapstructure:"msg,omitempty"`
}

type FunlenSettings struct {
	Lines          int  `mapstructure:"lines"`
	Statements     int  `mapstructure:"statements"`
	IgnoreComments bool `mapstructure:"ignore-comments"`
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
	ForceExpectTo              bool `mapstructure:"force-expect-to"`
	ValidateAsyncIntervals     bool `mapstructure:"validate-async-intervals"`
	ForbidSpecPollution        bool `mapstructure:"forbid-spec-pollution"`
	ForceSucceedForFuncs       bool `mapstructure:"force-succeed"`
}

type GoChecksumTypeSettings struct {
	DefaultSignifiesExhaustive bool `mapstructure:"default-signifies-exhaustive"`
	IncludeSharedInterfaces    bool `mapstructure:"include-shared-interfaces"`
}

type GocognitSettings struct {
	MinComplexity int `mapstructure:"min-complexity"`
}

type GoConstSettings struct {
	IgnoreStrings       string `mapstructure:"ignore-strings"`
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
}

type GodoxSettings struct {
	Keywords []string `mapstructure:"keywords"`
}

type GoHeaderSettings struct {
	Values       map[string]map[string]string `mapstructure:"values"`
	Template     string                       `mapstructure:"template"`
	TemplatePath string                       `mapstructure:"template-path"`
}

type GoModDirectivesSettings struct {
	ReplaceAllowList          []string `mapstructure:"replace-allow-list"`
	ReplaceLocal              bool     `mapstructure:"replace-local"`
	ExcludeForbidden          bool     `mapstructure:"exclude-forbidden"`
	RetractAllowNoExplanation bool     `mapstructure:"retract-allow-no-explanation"`
	ToolchainForbidden        bool     `mapstructure:"toolchain-forbidden"`
	ToolchainPattern          string   `mapstructure:"toolchain-pattern"`
	ToolForbidden             bool     `mapstructure:"tool-forbidden"`
	GoDebugForbidden          bool     `mapstructure:"go-debug-forbidden"`
	GoVersionPattern          string   `mapstructure:"go-version-pattern"`
}

type GoModGuardSettings struct {
	Allowed GoModGuardAllowed `mapstructure:"allowed"`
	Blocked GoModGuardBlocked `mapstructure:"blocked"`
}

type GoModGuardAllowed struct {
	Modules []string `mapstructure:"modules"`
	Domains []string `mapstructure:"domains"`
}

type GoModGuardBlocked struct {
	Modules                []map[string]GoModGuardModule  `mapstructure:"modules"`
	Versions               []map[string]GoModGuardVersion `mapstructure:"versions"`
	LocalReplaceDirectives bool                           `mapstructure:"local-replace-directives"`
}

type GoModGuardModule struct {
	Recommendations []string `mapstructure:"recommendations"`
	Reason          string   `mapstructure:"reason"`
}

type GoModGuardVersion struct {
	Version string `mapstructure:"version"`
	Reason  string `mapstructure:"reason"`
}

type GoSecSettings struct {
	Includes    []string       `mapstructure:"includes"`
	Excludes    []string       `mapstructure:"excludes"`
	Severity    string         `mapstructure:"severity"`
	Confidence  string         `mapstructure:"confidence"`
	Config      map[string]any `mapstructure:"config"`
	Concurrency int            `mapstructure:"concurrency"`
}

type GosmopolitanSettings struct {
	AllowTimeLocal  bool     `mapstructure:"allow-time-local"`
	EscapeHatches   []string `mapstructure:"escape-hatches"`
	WatchForScripts []string `mapstructure:"watch-for-scripts"`
}

type GovetSettings struct {
	Go string `mapstructure:"-"`

	Enable     []string `mapstructure:"enable"`
	Disable    []string `mapstructure:"disable"`
	EnableAll  bool     `mapstructure:"enable-all"`
	DisableAll bool     `mapstructure:"disable-all"`

	Settings map[string]map[string]any `mapstructure:"settings"`
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

type IfaceSettings struct {
	Enable   []string                  `mapstructure:"enable"`
	Settings map[string]map[string]any `mapstructure:"settings"`
}

type ImportAsSettings struct {
	Alias          []ImportAsAlias `mapstructure:"alias"`
	NoUnaliased    bool            `mapstructure:"no-unaliased"`
	NoExtraAliases bool            `mapstructure:"no-extra-aliases"`
}

type ImportAsAlias struct {
	Pkg   string `mapstructure:"pkg"`
	Alias string `mapstructure:"alias"`
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
	Slog             bool     `mapstructure:"slog"`
	Zap              bool     `mapstructure:"zap"`
	RequireStringKey bool     `mapstructure:"require-string-key"`
	NoPrintfLike     bool     `mapstructure:"no-printf-like"`
	Rules            []string `mapstructure:"rules"`
}

type MaintIdxSettings struct {
	Under int `mapstructure:"under"`
}

type MakezeroSettings struct {
	Always bool `mapstructure:"always"`
}

type MisspellSettings struct {
	Mode        string               `mapstructure:"mode"`
	Locale      string               `mapstructure:"locale"`
	ExtraWords  []MisspellExtraWords `mapstructure:"extra-words"`
	IgnoreRules []string             `mapstructure:"ignore-rules"`
}

type MisspellExtraWords struct {
	Typo       string `mapstructure:"typo"`
	Correction string `mapstructure:"correction"`
}

type MustTagSettings struct {
	Functions []MustTagFunction `mapstructure:"functions"`
}

type MustTagFunction struct {
	Name   string `mapstructure:"name"`
	Tag    string `mapstructure:"tag"`
	ArgPos int    `mapstructure:"arg-pos"`
}

type NakedretSettings struct {
	MaxFuncLines uint `mapstructure:"max-func-lines"`
}

type NestifSettings struct {
	MinComplexity int `mapstructure:"min-complexity"`
}

type NilNilSettings struct {
	OnlyTwo        *bool    `mapstructure:"only-two"`
	DetectOpposite bool     `mapstructure:"detect-opposite"`
	CheckedTypes   []string `mapstructure:"checked-types"`
}

type NlreturnSettings struct {
	BlockSize int `mapstructure:"block-size"`
}

type MndSettings struct {
	Checks           []string `mapstructure:"checks"`
	IgnoredNumbers   []string `mapstructure:"ignored-numbers"`
	IgnoredFiles     []string `mapstructure:"ignored-files"`
	IgnoredFunctions []string `mapstructure:"ignored-functions"`
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
	IntegerFormat bool `mapstructure:"integer-format"`
	IntConversion bool `mapstructure:"int-conversion"`

	ErrorFormat bool `mapstructure:"error-format"`
	ErrError    bool `mapstructure:"err-error"`
	ErrorF      bool `mapstructure:"errorf"`

	StringFormat bool `mapstructure:"string-format"`
	SprintF1     bool `mapstructure:"sprintf1"`
	StrConcat    bool `mapstructure:"strconcat"`

	BoolFormat bool `mapstructure:"bool-format"`
	HexFormat  bool `mapstructure:"hex-format"`
}

type PreallocSettings struct {
	Simple     bool `mapstructure:"simple"`
	RangeLoops bool `mapstructure:"range-loops"`
	ForLoops   bool `mapstructure:"for-loops"`
}

type PredeclaredSettings struct {
	Ignore    []string `mapstructure:"ignore"`
	Qualified bool     `mapstructure:"qualified-name"`
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

type RecvcheckSettings struct {
	DisableBuiltin bool     `mapstructure:"disable-builtin"`
	Exclusions     []string `mapstructure:"exclusions"`
}

type ReviveSettings struct {
	Go             string            `mapstructure:"-"`
	MaxOpenFiles   int               `mapstructure:"max-open-files"`
	Confidence     float64           `mapstructure:"confidence"`
	Severity       string            `mapstructure:"severity"`
	EnableAllRules bool              `mapstructure:"enable-all-rules"`
	Rules          []ReviveRule      `mapstructure:"rules"`
	ErrorCode      int               `mapstructure:"error-code"`
	WarningCode    int               `mapstructure:"warning-code"`
	Directives     []ReviveDirective `mapstructure:"directives"`
}

type ReviveRule struct {
	Name      string   `mapstructure:"name"`
	Arguments []any    `mapstructure:"arguments"`
	Severity  string   `mapstructure:"severity"`
	Disabled  bool     `mapstructure:"disabled"`
	Exclude   []string `mapstructure:"exclude"`
}

type ReviveDirective struct {
	Name     string `mapstructure:"name"`
	Severity string `mapstructure:"severity"`
}

type RowsErrCheckSettings struct {
	Packages []string `mapstructure:"packages"`
}

type SlogLintSettings struct {
	NoMixedArgs    bool     `mapstructure:"no-mixed-args"`
	KVOnly         bool     `mapstructure:"kv-only"`
	AttrOnly       bool     `mapstructure:"attr-only"`
	NoGlobal       string   `mapstructure:"no-global"`
	Context        string   `mapstructure:"context"`
	StaticMsg      bool     `mapstructure:"static-msg"`
	NoRawKeys      bool     `mapstructure:"no-raw-keys"`
	KeyNamingCase  string   `mapstructure:"key-naming-case"`
	ForbiddenKeys  []string `mapstructure:"forbidden-keys"`
	ArgsOnSepLines bool     `mapstructure:"args-on-sep-lines"`
}

type SpancheckSettings struct {
	Checks                   []string `mapstructure:"checks"`
	IgnoreCheckSignatures    []string `mapstructure:"ignore-check-signatures"`
	ExtraStartSpanSignatures []string `mapstructure:"extra-start-span-signatures"`
}

type StaticCheckSettings struct {
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
	Case TagliatelleCase `mapstructure:"case"`
}

type TagliatelleCase struct {
	TagliatelleBase `mapstructure:",squash"`
	Overrides       []TagliatelleOverrides `mapstructure:"overrides"`
}

type TagliatelleOverrides struct {
	TagliatelleBase `mapstructure:",squash"`
	Package         string `mapstructure:"pkg"`
	Ignore          bool   `mapstructure:"ignore"`
}

type TagliatelleBase struct {
	Rules         map[string]string                  `mapstructure:"rules"`
	ExtendedRules map[string]TagliatelleExtendedRule `mapstructure:"extended-rules"`
	UseFieldName  bool                               `mapstructure:"use-field-name"`
	IgnoredFields []string                           `mapstructure:"ignored-fields"`
}

type TagliatelleExtendedRule struct {
	Case                string          `mapstructure:"case"`
	ExtraInitialisms    bool            `mapstructure:"extra-initialisms"`
	InitialismOverrides map[string]bool `mapstructure:"initialism-overrides"`
}

type TestifylintSettings struct {
	EnableAll        bool     `mapstructure:"enable-all"`
	DisableAll       bool     `mapstructure:"disable-all"`
	EnabledCheckers  []string `mapstructure:"enable"`
	DisabledCheckers []string `mapstructure:"disable"`

	BoolCompare          TestifylintBoolCompare          `mapstructure:"bool-compare"`
	ExpectedActual       TestifylintExpectedActual       `mapstructure:"expected-actual"`
	Formatter            TestifylintFormatter            `mapstructure:"formatter"`
	GoRequire            TestifylintGoRequire            `mapstructure:"go-require"`
	RequireError         TestifylintRequireError         `mapstructure:"require-error"`
	SuiteExtraAssertCall TestifylintSuiteExtraAssertCall `mapstructure:"suite-extra-assert-call"`
}

type TestifylintBoolCompare struct {
	IgnoreCustomTypes bool `mapstructure:"ignore-custom-types"`
}

type TestifylintExpectedActual struct {
	ExpVarPattern string `mapstructure:"pattern"`
}

type TestifylintFormatter struct {
	CheckFormatString *bool `mapstructure:"check-format-string"`
	RequireFFuncs     bool  `mapstructure:"require-f-funcs"`
	RequireStringMsg  bool  `mapstructure:"require-string-msg"`
}

type TestifylintGoRequire struct {
	IgnoreHTTPHandlers bool `mapstructure:"ignore-http-handlers"`
}

type TestifylintRequireError struct {
	FnPattern string `mapstructure:"fn-pattern"`
}

type TestifylintSuiteExtraAssertCall struct {
	Mode string `mapstructure:"mode"`
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

type UseStdlibVarsSettings struct {
	HTTPMethod         bool `mapstructure:"http-method"`
	HTTPStatusCode     bool `mapstructure:"http-status-code"`
	TimeWeekday        bool `mapstructure:"time-weekday"`
	TimeMonth          bool `mapstructure:"time-month"`
	TimeLayout         bool `mapstructure:"time-layout"`
	CryptoHash         bool `mapstructure:"crypto-hash"`
	DefaultRPCPath     bool `mapstructure:"default-rpc-path"`
	SQLIsolationLevel  bool `mapstructure:"sql-isolation-level"`
	TLSSignatureScheme bool `mapstructure:"tls-signature-scheme"`
	ConstantKind       bool `mapstructure:"constant-kind"`
}

type UseTestingSettings struct {
	ContextBackground bool `mapstructure:"context-background"`
	ContextTodo       bool `mapstructure:"context-todo"`
	OSChdir           bool `mapstructure:"os-chdir"`
	OSMkdirTemp       bool `mapstructure:"os-mkdir-temp"`
	OSSetenv          bool `mapstructure:"os-setenv"`
	OSTempDir         bool `mapstructure:"os-temp-dir"`
	OSCreateTemp      bool `mapstructure:"os-create-temp"`
}

type UnconvertSettings struct {
	FastMath bool `mapstructure:"fast-math"`
	Safe     bool `mapstructure:"safe"`
}

type UnparamSettings struct {
	CheckExported bool `mapstructure:"check-exported"`
}

type UnusedSettings struct {
	FieldWritesAreUses     bool `mapstructure:"field-writes-are-uses"`
	PostStatementsAreReads bool `mapstructure:"post-statements-are-reads"`
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
	ExtraIgnoreSigs        []string `mapstructure:"extra-ignore-sigs"`
	IgnoreSigs             []string `mapstructure:"ignore-sigs"`
	IgnoreSigRegexps       []string `mapstructure:"ignore-sig-regexps"`
	IgnorePackageGlobs     []string `mapstructure:"ignore-package-globs"`
	IgnoreInterfaceRegexps []string `mapstructure:"ignore-interface-regexps"`
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
	AllowCuddleUsedInBlock           bool     `mapstructure:"allow-cuddle-used-in-block"`
	ForceCuddleErrCheckAndAssign     bool     `mapstructure:"force-err-cuddling"`
	ErrorVariableNames               []string `mapstructure:"error-variable-names"`
	ForceExclusiveShortDeclarations  bool     `mapstructure:"force-short-decl-cuddling"`
}

// CustomLinterSettings encapsulates the meta-data of a private linter.
type CustomLinterSettings struct {
	// Type plugin type.
	// It can be `goplugin` or `module`.
	Type string `mapstructure:"type"`

	// Path to a plugin *.so file that implements the private linter.
	// Only for Go plugin system.
	Path string `mapstructure:"path"`

	// Description describes the purpose of the private linter.
	Description string `mapstructure:"description"`
	// OriginalURL The URL containing the source code for the private linter.
	OriginalURL string `mapstructure:"original-url"`

	// Settings plugin settings only work with linterdb.PluginConstructor symbol.
	Settings any `mapstructure:"settings"`
}

func (s *CustomLinterSettings) Validate() error {
	if s.Type == "module" {
		if s.Path != "" {
			return errors.New("path not supported with module type")
		}

		return nil
	}

	if s.Path == "" {
		return errors.New("path is required")
	}

	return nil
}
