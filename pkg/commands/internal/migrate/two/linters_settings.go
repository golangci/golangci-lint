package two

type LintersSettings struct {
	FormatterSettings `yaml:"-,omitempty" toml:"-,omitempty"`

	Asasalint       AsasalintSettings       `yaml:"asasalint,omitempty" toml:"asasalint,omitempty"`
	BiDiChk         BiDiChkSettings         `yaml:"bidichk,omitempty" toml:"bidichk,omitempty"`
	CopyLoopVar     CopyLoopVarSettings     `yaml:"copyloopvar,omitempty" toml:"copyloopvar,omitempty"`
	Cyclop          CyclopSettings          `yaml:"cyclop,omitempty" toml:"cyclop,omitempty"`
	Decorder        DecorderSettings        `yaml:"decorder,omitempty" toml:"decorder,omitempty"`
	Depguard        DepGuardSettings        `yaml:"depguard,omitempty" toml:"depguard,omitempty"`
	Dogsled         DogsledSettings         `yaml:"dogsled,omitempty" toml:"dogsled,omitempty"`
	Dupl            DuplSettings            `yaml:"dupl,omitempty" toml:"dupl,omitempty"`
	DupWord         DupWordSettings         `yaml:"dupword,omitempty" toml:"dupword,omitempty"`
	Errcheck        ErrcheckSettings        `yaml:"errcheck,omitempty" toml:"errcheck,omitempty"`
	ErrChkJSON      ErrChkJSONSettings      `yaml:"errchkjson,omitempty" toml:"errchkjson,omitempty"`
	ErrorLint       ErrorLintSettings       `yaml:"errorlint,omitempty" toml:"errorlint,omitempty"`
	Exhaustive      ExhaustiveSettings      `yaml:"exhaustive,omitempty" toml:"exhaustive,omitempty"`
	Exhaustruct     ExhaustructSettings     `yaml:"exhaustruct,omitempty" toml:"exhaustruct,omitempty"`
	Fatcontext      FatcontextSettings      `yaml:"fatcontext,omitempty" toml:"fatcontext,omitempty"`
	Forbidigo       ForbidigoSettings       `yaml:"forbidigo,omitempty" toml:"forbidigo,omitempty"`
	Funlen          FunlenSettings          `yaml:"funlen,omitempty" toml:"funlen,omitempty"`
	GinkgoLinter    GinkgoLinterSettings    `yaml:"ginkgolinter,omitempty" toml:"ginkgolinter,omitempty"`
	Gocognit        GocognitSettings        `yaml:"gocognit,omitempty" toml:"gocognit,omitempty"`
	GoChecksumType  GoChecksumTypeSettings  `yaml:"gochecksumtype,omitempty" toml:"gochecksumtype,omitempty"`
	Goconst         GoConstSettings         `yaml:"goconst,omitempty" toml:"goconst,omitempty"`
	Gocritic        GoCriticSettings        `yaml:"gocritic,omitempty" toml:"gocritic,omitempty"`
	Gocyclo         GoCycloSettings         `yaml:"gocyclo,omitempty" toml:"gocyclo,omitempty"`
	Godot           GodotSettings           `yaml:"godot,omitempty" toml:"godot,omitempty"`
	Godox           GodoxSettings           `yaml:"godox,omitempty" toml:"godox,omitempty"`
	Goheader        GoHeaderSettings        `yaml:"goheader,omitempty" toml:"goheader,omitempty"`
	GoModDirectives GoModDirectivesSettings `yaml:"gomoddirectives,omitempty" toml:"gomoddirectives,omitempty"`
	Gomodguard      GoModGuardSettings      `yaml:"gomodguard,omitempty" toml:"gomodguard,omitempty"`
	Gosec           GoSecSettings           `yaml:"gosec,omitempty" toml:"gosec,omitempty"`
	Gosmopolitan    GosmopolitanSettings    `yaml:"gosmopolitan,omitempty" toml:"gosmopolitan,omitempty"`
	Govet           GovetSettings           `yaml:"govet,omitempty" toml:"govet,omitempty"`
	Grouper         GrouperSettings         `yaml:"grouper,omitempty" toml:"grouper,omitempty"`
	Iface           IfaceSettings           `yaml:"iface,omitempty" toml:"iface,omitempty"`
	ImportAs        ImportAsSettings        `yaml:"importas,omitempty" toml:"importas,omitempty"`
	Inamedparam     INamedParamSettings     `yaml:"inamedparam,omitempty" toml:"inamedparam,omitempty"`
	InterfaceBloat  InterfaceBloatSettings  `yaml:"interfacebloat,omitempty" toml:"interfacebloat,omitempty"`
	Ireturn         IreturnSettings         `yaml:"ireturn,omitempty" toml:"ireturn,omitempty"`
	Lll             LllSettings             `yaml:"lll,omitempty" toml:"lll,omitempty"`
	LoggerCheck     LoggerCheckSettings     `yaml:"loggercheck,omitempty" toml:"loggercheck,omitempty"`
	MaintIdx        MaintIdxSettings        `yaml:"maintidx,omitempty" toml:"maintidx,omitempty"`
	Makezero        MakezeroSettings        `yaml:"makezero,omitempty" toml:"makezero,omitempty"`
	Misspell        MisspellSettings        `yaml:"misspell,omitempty" toml:"misspell,omitempty"`
	Mnd             MndSettings             `yaml:"mnd,omitempty" toml:"mnd,omitempty"`
	MustTag         MustTagSettings         `yaml:"musttag,omitempty" toml:"musttag,omitempty"`
	Nakedret        NakedretSettings        `yaml:"nakedret,omitempty" toml:"nakedret,omitempty"`
	Nestif          NestifSettings          `yaml:"nestif,omitempty" toml:"nestif,omitempty"`
	NilNil          NilNilSettings          `yaml:"nilnil,omitempty" toml:"nilnil,omitempty"`
	Nlreturn        NlreturnSettings        `yaml:"nlreturn,omitempty" toml:"nlreturn,omitempty"`
	NoLintLint      NoLintLintSettings      `yaml:"nolintlint,omitempty" toml:"nolintlint,omitempty"`
	NoNamedReturns  NoNamedReturnsSettings  `yaml:"nonamedreturns,omitempty" toml:"nonamedreturns,omitempty"`
	ParallelTest    ParallelTestSettings    `yaml:"paralleltest,omitempty" toml:"paralleltest,omitempty"`
	PerfSprint      PerfSprintSettings      `yaml:"perfsprint,omitempty" toml:"perfsprint,omitempty"`
	Prealloc        PreallocSettings        `yaml:"prealloc,omitempty" toml:"prealloc,omitempty"`
	Predeclared     PredeclaredSettings     `yaml:"predeclared,omitempty" toml:"predeclared,omitempty"`
	Promlinter      PromlinterSettings      `yaml:"promlinter,omitempty" toml:"promlinter,omitempty"`
	ProtoGetter     ProtoGetterSettings     `yaml:"protogetter,omitempty" toml:"protogetter,omitempty"`
	Reassign        ReassignSettings        `yaml:"reassign,omitempty" toml:"reassign,omitempty"`
	Recvcheck       RecvcheckSettings       `yaml:"recvcheck,omitempty" toml:"recvcheck,omitempty"`
	Revive          ReviveSettings          `yaml:"revive,omitempty" toml:"revive,omitempty"`
	RowsErrCheck    RowsErrCheckSettings    `yaml:"rowserrcheck,omitempty" toml:"rowserrcheck,omitempty"`
	SlogLint        SlogLintSettings        `yaml:"sloglint,omitempty" toml:"sloglint,omitempty"`
	Spancheck       SpancheckSettings       `yaml:"spancheck,omitempty" toml:"spancheck,omitempty"`
	Staticcheck     StaticCheckSettings     `yaml:"staticcheck,omitempty" toml:"staticcheck,omitempty"`
	TagAlign        TagAlignSettings        `yaml:"tagalign,omitempty" toml:"tagalign,omitempty"`
	Tagliatelle     TagliatelleSettings     `yaml:"tagliatelle,omitempty" toml:"tagliatelle,omitempty"`
	Tenv            TenvSettings            `yaml:"tenv,omitempty" toml:"tenv,omitempty"`
	Testifylint     TestifylintSettings     `yaml:"testifylint,omitempty" toml:"testifylint,omitempty"`
	Testpackage     TestpackageSettings     `yaml:"testpackage,omitempty" toml:"testpackage,omitempty"`
	Thelper         ThelperSettings         `yaml:"thelper,omitempty" toml:"thelper,omitempty"`
	Unconvert       UnconvertSettings       `yaml:"unconvert,omitempty" toml:"unconvert,omitempty"`
	Unparam         UnparamSettings         `yaml:"unparam,omitempty" toml:"unparam,omitempty"`
	Unused          UnusedSettings          `yaml:"unused,omitempty" toml:"unused,omitempty"`
	UseStdlibVars   UseStdlibVarsSettings   `yaml:"usestdlibvars,omitempty" toml:"usestdlibvars,omitempty"`
	UseTesting      UseTestingSettings      `yaml:"usetesting,omitempty" toml:"usetesting,omitempty"`
	Varnamelen      VarnamelenSettings      `yaml:"varnamelen,omitempty" toml:"varnamelen,omitempty"`
	Whitespace      WhitespaceSettings      `yaml:"whitespace,omitempty" toml:"whitespace,omitempty"`
	Wrapcheck       WrapcheckSettings       `yaml:"wrapcheck,omitempty" toml:"wrapcheck,omitempty"`
	WSL             WSLSettings             `yaml:"wsl,omitempty" toml:"wsl,omitempty"`

	Custom map[string]CustomLinterSettings `yaml:"custom,omitempty" toml:"custom,omitempty"`
}

type AsasalintSettings struct {
	Exclude              []string `yaml:"exclude,omitempty" toml:"exclude,omitempty"`
	UseBuiltinExclusions *bool    `yaml:"use-builtin-exclusions,omitempty" toml:"use-builtin-exclusions,omitempty"`
}

type BiDiChkSettings struct {
	LeftToRightEmbedding     *bool `yaml:"left-to-right-embedding,omitempty" toml:"left-to-right-embedding,omitempty"`
	RightToLeftEmbedding     *bool `yaml:"right-to-left-embedding,omitempty" toml:"right-to-left-embedding,omitempty"`
	PopDirectionalFormatting *bool `yaml:"pop-directional-formatting,omitempty" toml:"pop-directional-formatting,omitempty"`
	LeftToRightOverride      *bool `yaml:"left-to-right-override,omitempty" toml:"left-to-right-override,omitempty"`
	RightToLeftOverride      *bool `yaml:"right-to-left-override,omitempty" toml:"right-to-left-override,omitempty"`
	LeftToRightIsolate       *bool `yaml:"left-to-right-isolate,omitempty" toml:"left-to-right-isolate,omitempty"`
	RightToLeftIsolate       *bool `yaml:"right-to-left-isolate,omitempty" toml:"right-to-left-isolate,omitempty"`
	FirstStrongIsolate       *bool `yaml:"first-strong-isolate,omitempty" toml:"first-strong-isolate,omitempty"`
	PopDirectionalIsolate    *bool `yaml:"pop-directional-isolate,omitempty" toml:"pop-directional-isolate,omitempty"`
}

type CopyLoopVarSettings struct {
	CheckAlias *bool `yaml:"check-alias,omitempty" toml:"check-alias,omitempty"`
}

type CyclopSettings struct {
	MaxComplexity  *int     `yaml:"max-complexity,omitempty" toml:"max-complexity,omitempty"`
	PackageAverage *float64 `yaml:"package-average,omitempty" toml:"package-average,omitempty"`
}

type DepGuardSettings struct {
	Rules map[string]*DepGuardList `yaml:"rules,omitempty" toml:"rules,omitempty"`
}

type DepGuardList struct {
	ListMode *string        `yaml:"list-mode,omitempty" toml:"list-mode,omitempty"`
	Files    []string       `yaml:"files,omitempty" toml:"files,omitempty"`
	Allow    []string       `yaml:"allow,omitempty" toml:"allow,omitempty"`
	Deny     []DepGuardDeny `yaml:"deny,omitempty" toml:"deny,omitempty"`
}

type DepGuardDeny struct {
	Pkg  *string `yaml:"pkg,omitempty" toml:"pkg,omitempty"`
	Desc *string `yaml:"desc,omitempty" toml:"desc,omitempty"`
}

type DecorderSettings struct {
	DecOrder                  []string `yaml:"dec-order,omitempty" toml:"dec-order,omitempty"`
	IgnoreUnderscoreVars      *bool    `yaml:"ignore-underscore-vars,omitempty" toml:"ignore-underscore-vars,omitempty"`
	DisableDecNumCheck        *bool    `yaml:"disable-dec-num-check,omitempty" toml:"disable-dec-num-check,omitempty"`
	DisableTypeDecNumCheck    *bool    `yaml:"disable-type-dec-num-check,omitempty" toml:"disable-type-dec-num-check,omitempty"`
	DisableConstDecNumCheck   *bool    `yaml:"disable-const-dec-num-check,omitempty" toml:"disable-const-dec-num-check,omitempty"`
	DisableVarDecNumCheck     *bool    `yaml:"disable-var-dec-num-check,omitempty" toml:"disable-var-dec-num-check,omitempty"`
	DisableDecOrderCheck      *bool    `yaml:"disable-dec-order-check,omitempty" toml:"disable-dec-order-check,omitempty"`
	DisableInitFuncFirstCheck *bool    `yaml:"disable-init-func-first-check,omitempty" toml:"disable-init-func-first-check,omitempty"`
}

type DogsledSettings struct {
	MaxBlankIdentifiers *int `yaml:"max-blank-identifiers,omitempty" toml:"max-blank-identifiers,omitempty"`
}

type DuplSettings struct {
	Threshold *int `yaml:"threshold,omitempty" toml:"threshold,omitempty"`
}

type DupWordSettings struct {
	Keywords []string `yaml:"keywords,omitempty" toml:"keywords,omitempty"`
	Ignore   []string `yaml:"ignore,omitempty" toml:"ignore,omitempty"`
}

type ErrcheckSettings struct {
	DisableDefaultExclusions *bool    `yaml:"disable-default-exclusions,omitempty" toml:"disable-default-exclusions,omitempty"`
	CheckTypeAssertions      *bool    `yaml:"check-type-assertions,omitempty" toml:"check-type-assertions,omitempty"`
	CheckAssignToBlank       *bool    `yaml:"check-blank,omitempty" toml:"check-blank,omitempty"`
	ExcludeFunctions         []string `yaml:"exclude-functions,omitempty" toml:"exclude-functions,omitempty"`
}

type ErrChkJSONSettings struct {
	CheckErrorFreeEncoding *bool `yaml:"check-error-free-encoding,omitempty" toml:"check-error-free-encoding,omitempty"`
	ReportNoExported       *bool `yaml:"report-no-exported,omitempty" toml:"report-no-exported,omitempty"`
}

type ErrorLintSettings struct {
	Errorf                *bool                `yaml:"errorf,omitempty" toml:"errorf,omitempty"`
	ErrorfMulti           *bool                `yaml:"errorf-multi,omitempty" toml:"errorf-multi,omitempty"`
	Asserts               *bool                `yaml:"asserts,omitempty" toml:"asserts,omitempty"`
	Comparison            *bool                `yaml:"comparison,omitempty" toml:"comparison,omitempty"`
	AllowedErrors         []ErrorLintAllowPair `yaml:"allowed-errors,omitempty" toml:"allowed-errors,omitempty"`
	AllowedErrorsWildcard []ErrorLintAllowPair `yaml:"allowed-errors-wildcard,omitempty" toml:"allowed-errors-wildcard,omitempty"`
}

type ErrorLintAllowPair struct {
	Err *string `yaml:"err,omitempty" toml:"err,omitempty"`
	Fun *string `yaml:"fun,omitempty" toml:"fun,omitempty"`
}

type ExhaustiveSettings struct {
	Check                      []string `yaml:"check,omitempty" toml:"check,omitempty"`
	DefaultSignifiesExhaustive *bool    `yaml:"default-signifies-exhaustive,omitempty" toml:"default-signifies-exhaustive,omitempty"`
	IgnoreEnumMembers          *string  `yaml:"ignore-enum-members,omitempty" toml:"ignore-enum-members,omitempty"`
	IgnoreEnumTypes            *string  `yaml:"ignore-enum-types,omitempty" toml:"ignore-enum-types,omitempty"`
	PackageScopeOnly           *bool    `yaml:"package-scope-only,omitempty" toml:"package-scope-only,omitempty"`
	ExplicitExhaustiveMap      *bool    `yaml:"explicit-exhaustive-map,omitempty" toml:"explicit-exhaustive-map,omitempty"`
	ExplicitExhaustiveSwitch   *bool    `yaml:"explicit-exhaustive-switch,omitempty" toml:"explicit-exhaustive-switch,omitempty"`
	DefaultCaseRequired        *bool    `yaml:"default-case-required,omitempty" toml:"default-case-required,omitempty"`
}

type ExhaustructSettings struct {
	Include []string `yaml:"include,omitempty" toml:"include,omitempty"`
	Exclude []string `yaml:"exclude,omitempty" toml:"exclude,omitempty"`
}

type FatcontextSettings struct {
	CheckStructPointers *bool `yaml:"check-struct-pointers,omitempty" toml:"check-struct-pointers,omitempty"`
}

type ForbidigoSettings struct {
	Forbid               []ForbidigoPattern `yaml:"forbid,omitempty" toml:"forbid,omitempty"`
	ExcludeGodocExamples *bool              `yaml:"exclude-godoc-examples,omitempty" toml:"exclude-godoc-examples,omitempty"`
	AnalyzeTypes         *bool              `yaml:"analyze-types,omitempty" toml:"analyze-types,omitempty"`
}

type ForbidigoPattern struct {
	Pattern *string `yaml:"pattern,omitempty" toml:"pattern,omitempty"`
	Package *string `yaml:"pkg,omitempty,omitempty" toml:"pkg,omitempty,omitempty"`
	Msg     *string `yaml:"msg,omitempty,omitempty" toml:"msg,omitempty,omitempty"`
}

type FunlenSettings struct {
	Lines          *int  `yaml:"lines,omitempty" toml:"lines,omitempty"`
	Statements     *int  `yaml:"statements,omitempty" toml:"statements,omitempty"`
	IgnoreComments *bool `yaml:"ignore-comments,omitempty" toml:"ignore-comments,omitempty"`
}

type GinkgoLinterSettings struct {
	SuppressLenAssertion       *bool `yaml:"suppress-len-assertion,omitempty" toml:"suppress-len-assertion,omitempty"`
	SuppressNilAssertion       *bool `yaml:"suppress-nil-assertion,omitempty" toml:"suppress-nil-assertion,omitempty"`
	SuppressErrAssertion       *bool `yaml:"suppress-err-assertion,omitempty" toml:"suppress-err-assertion,omitempty"`
	SuppressCompareAssertion   *bool `yaml:"suppress-compare-assertion,omitempty" toml:"suppress-compare-assertion,omitempty"`
	SuppressAsyncAssertion     *bool `yaml:"suppress-async-assertion,omitempty" toml:"suppress-async-assertion,omitempty"`
	SuppressTypeCompareWarning *bool `yaml:"suppress-type-compare-assertion,omitempty" toml:"suppress-type-compare-assertion,omitempty"`
	ForbidFocusContainer       *bool `yaml:"forbid-focus-container,omitempty" toml:"forbid-focus-container,omitempty"`
	AllowHaveLenZero           *bool `yaml:"allow-havelen-zero,omitempty" toml:"allow-havelen-zero,omitempty"`
	ForceExpectTo              *bool `yaml:"force-expect-to,omitempty" toml:"force-expect-to,omitempty"`
	ValidateAsyncIntervals     *bool `yaml:"validate-async-intervals,omitempty" toml:"validate-async-intervals,omitempty"`
	ForbidSpecPollution        *bool `yaml:"forbid-spec-pollution,omitempty" toml:"forbid-spec-pollution,omitempty"`
	ForceSucceedForFuncs       *bool `yaml:"force-succeed,omitempty" toml:"force-succeed,omitempty"`
}

type GoChecksumTypeSettings struct {
	DefaultSignifiesExhaustive *bool `yaml:"default-signifies-exhaustive,omitempty" toml:"default-signifies-exhaustive,omitempty"`
	IncludeSharedInterfaces    *bool `yaml:"include-shared-interfaces,omitempty" toml:"include-shared-interfaces,omitempty"`
}

type GocognitSettings struct {
	MinComplexity *int `yaml:"min-complexity,omitempty" toml:"min-complexity,omitempty"`
}

type GoConstSettings struct {
	IgnoreStrings       *string `yaml:"ignore-strings,omitempty" toml:"ignore-strings,omitempty"`
	MatchWithConstants  *bool   `yaml:"match-constant,omitempty" toml:"match-constant,omitempty"`
	MinStringLen        *int    `yaml:"min-len,omitempty" toml:"min-len,omitempty"`
	MinOccurrencesCount *int    `yaml:"min-occurrences,omitempty" toml:"min-occurrences,omitempty"`
	ParseNumbers        *bool   `yaml:"numbers,omitempty" toml:"numbers,omitempty"`
	NumberMin           *int    `yaml:"min,omitempty" toml:"min,omitempty"`
	NumberMax           *int    `yaml:"max,omitempty" toml:"max,omitempty"`
	IgnoreCalls         *bool   `yaml:"ignore-calls,omitempty" toml:"ignore-calls,omitempty"`
}

type GoCriticSettings struct {
	Go               *string                          `yaml:"-,omitempty" toml:"-,omitempty"`
	DisableAll       *bool                            `yaml:"disable-all,omitempty" toml:"disable-all,omitempty"`
	EnabledChecks    []string                         `yaml:"enabled-checks,omitempty" toml:"enabled-checks,omitempty"`
	EnableAll        *bool                            `yaml:"enable-all,omitempty" toml:"enable-all,omitempty"`
	DisabledChecks   []string                         `yaml:"disabled-checks,omitempty" toml:"disabled-checks,omitempty"`
	EnabledTags      []string                         `yaml:"enabled-tags,omitempty" toml:"enabled-tags,omitempty"`
	DisabledTags     []string                         `yaml:"disabled-tags,omitempty" toml:"disabled-tags,omitempty"`
	SettingsPerCheck map[string]GoCriticCheckSettings `yaml:"settings,omitempty" toml:"settings,omitempty"`
}

type GoCriticCheckSettings map[string]any

type GoCycloSettings struct {
	MinComplexity *int `yaml:"min-complexity,omitempty" toml:"min-complexity,omitempty"`
}

type GodotSettings struct {
	Scope   *string  `yaml:"scope,omitempty" toml:"scope,omitempty"`
	Exclude []string `yaml:"exclude,omitempty" toml:"exclude,omitempty"`
	Capital *bool    `yaml:"capital,omitempty" toml:"capital,omitempty"`
	Period  *bool    `yaml:"period,omitempty" toml:"period,omitempty"`
}

type GodoxSettings struct {
	Keywords []string `yaml:"keywords,omitempty" toml:"keywords,omitempty"`
}

type GoHeaderSettings struct {
	Values       map[string]map[string]string `yaml:"values,omitempty" toml:"values,omitempty"`
	Template     *string                      `yaml:"template,omitempty" toml:"template,omitempty"`
	TemplatePath *string                      `yaml:"template-path,omitempty" toml:"template-path,omitempty"`
}

type GoModDirectivesSettings struct {
	ReplaceAllowList          []string `yaml:"replace-allow-list,omitempty" toml:"replace-allow-list,omitempty"`
	ReplaceLocal              *bool    `yaml:"replace-local,omitempty" toml:"replace-local,omitempty"`
	ExcludeForbidden          *bool    `yaml:"exclude-forbidden,omitempty" toml:"exclude-forbidden,omitempty"`
	RetractAllowNoExplanation *bool    `yaml:"retract-allow-no-explanation,omitempty" toml:"retract-allow-no-explanation,omitempty"`
	ToolchainForbidden        *bool    `yaml:"toolchain-forbidden,omitempty" toml:"toolchain-forbidden,omitempty"`
	ToolchainPattern          *string  `yaml:"toolchain-pattern,omitempty" toml:"toolchain-pattern,omitempty"`
	ToolForbidden             *bool    `yaml:"tool-forbidden,omitempty" toml:"tool-forbidden,omitempty"`
	GoDebugForbidden          *bool    `yaml:"go-debug-forbidden,omitempty" toml:"go-debug-forbidden,omitempty"`
	GoVersionPattern          *string  `yaml:"go-version-pattern,omitempty" toml:"go-version-pattern,omitempty"`
}

type GoModGuardSettings struct {
	Allowed GoModGuardAllowed `yaml:"allowed,omitempty" toml:"allowed,omitempty"`
	Blocked GoModGuardBlocked `yaml:"blocked,omitempty" toml:"blocked,omitempty"`
}

type GoModGuardAllowed struct {
	Modules []string `yaml:"modules,omitempty" toml:"modules,omitempty"`
	Domains []string `yaml:"domains,omitempty" toml:"domains,omitempty"`
}

type GoModGuardBlocked struct {
	Modules                []map[string]GoModGuardModule  `yaml:"modules,omitempty" toml:"modules,omitempty"`
	Versions               []map[string]GoModGuardVersion `yaml:"versions,omitempty" toml:"versions,omitempty"`
	LocalReplaceDirectives *bool                          `yaml:"local-replace-directives,omitempty" toml:"local-replace-directives,omitempty"`
}

type GoModGuardModule struct {
	Recommendations []string `yaml:"recommendations,omitempty" toml:"recommendations,omitempty"`
	Reason          *string  `yaml:"reason,omitempty" toml:"reason,omitempty"`
}

type GoModGuardVersion struct {
	Version *string `yaml:"version,omitempty" toml:"version,omitempty"`
	Reason  *string `yaml:"reason,omitempty" toml:"reason,omitempty"`
}

type GoSecSettings struct {
	Includes    []string       `yaml:"includes,omitempty" toml:"includes,omitempty"`
	Excludes    []string       `yaml:"excludes,omitempty" toml:"excludes,omitempty"`
	Severity    *string        `yaml:"severity,omitempty" toml:"severity,omitempty"`
	Confidence  *string        `yaml:"confidence,omitempty" toml:"confidence,omitempty"`
	Config      map[string]any `yaml:"config,omitempty" toml:"config,omitempty"`
	Concurrency *int           `yaml:"concurrency,omitempty" toml:"concurrency,omitempty"`
}

type GosmopolitanSettings struct {
	AllowTimeLocal  *bool    `yaml:"allow-time-local,omitempty" toml:"allow-time-local,omitempty"`
	EscapeHatches   []string `yaml:"escape-hatches,omitempty" toml:"escape-hatches,omitempty"`
	WatchForScripts []string `yaml:"watch-for-scripts,omitempty" toml:"watch-for-scripts,omitempty"`
}

type GovetSettings struct {
	Go *string `yaml:"-,omitempty" toml:"-,omitempty"`

	Enable     []string `yaml:"enable,omitempty" toml:"enable,omitempty"`
	Disable    []string `yaml:"disable,omitempty" toml:"disable,omitempty"`
	EnableAll  *bool    `yaml:"enable-all,omitempty" toml:"enable-all,omitempty"`
	DisableAll *bool    `yaml:"disable-all,omitempty" toml:"disable-all,omitempty"`

	Settings map[string]map[string]any `yaml:"settings,omitempty" toml:"settings,omitempty"`
}

type GrouperSettings struct {
	ConstRequireSingleConst   *bool `yaml:"const-require-single-const,omitempty" toml:"const-require-single-const,omitempty"`
	ConstRequireGrouping      *bool `yaml:"const-require-grouping,omitempty" toml:"const-require-grouping,omitempty"`
	ImportRequireSingleImport *bool `yaml:"import-require-single-import,omitempty" toml:"import-require-single-import,omitempty"`
	ImportRequireGrouping     *bool `yaml:"import-require-grouping,omitempty" toml:"import-require-grouping,omitempty"`
	TypeRequireSingleType     *bool `yaml:"type-require-single-type,omitempty" toml:"type-require-single-type,omitempty"`
	TypeRequireGrouping       *bool `yaml:"type-require-grouping,omitempty" toml:"type-require-grouping,omitempty"`
	VarRequireSingleVar       *bool `yaml:"var-require-single-var,omitempty" toml:"var-require-single-var,omitempty"`
	VarRequireGrouping        *bool `yaml:"var-require-grouping,omitempty" toml:"var-require-grouping,omitempty"`
}

type IfaceSettings struct {
	Enable   []string                  `yaml:"enable,omitempty" toml:"enable,omitempty"`
	Settings map[string]map[string]any `yaml:"settings,omitempty" toml:"settings,omitempty"`
}

type ImportAsSettings struct {
	Alias          []ImportAsAlias `yaml:"alias,omitempty" toml:"alias,omitempty"`
	NoUnaliased    *bool           `yaml:"no-unaliased,omitempty" toml:"no-unaliased,omitempty"`
	NoExtraAliases *bool           `yaml:"no-extra-aliases,omitempty" toml:"no-extra-aliases,omitempty"`
}

type ImportAsAlias struct {
	Pkg   *string `yaml:"pkg,omitempty" toml:"pkg,omitempty"`
	Alias *string `yaml:"alias,omitempty" toml:"alias,omitempty"`
}

type INamedParamSettings struct {
	SkipSingleParam *bool `yaml:"skip-single-param,omitempty" toml:"skip-single-param,omitempty"`
}

type InterfaceBloatSettings struct {
	Max *int `yaml:"max,omitempty" toml:"max,omitempty"`
}

type IreturnSettings struct {
	Allow  []string `yaml:"allow,omitempty" toml:"allow,omitempty"`
	Reject []string `yaml:"reject,omitempty" toml:"reject,omitempty"`
}

type LllSettings struct {
	LineLength *int `yaml:"line-length,omitempty" toml:"line-length,omitempty"`
	TabWidth   *int `yaml:"tab-width,omitempty" toml:"tab-width,omitempty"`
}

type LoggerCheckSettings struct {
	Kitlog           *bool    `yaml:"kitlog,omitempty" toml:"kitlog,omitempty"`
	Klog             *bool    `yaml:"klog,omitempty" toml:"klog,omitempty"`
	Logr             *bool    `yaml:"logr,omitempty" toml:"logr,omitempty"`
	Slog             *bool    `yaml:"slog,omitempty" toml:"slog,omitempty"`
	Zap              *bool    `yaml:"zap,omitempty" toml:"zap,omitempty"`
	RequireStringKey *bool    `yaml:"require-string-key,omitempty" toml:"require-string-key,omitempty"`
	NoPrintfLike     *bool    `yaml:"no-printf-like,omitempty" toml:"no-printf-like,omitempty"`
	Rules            []string `yaml:"rules,omitempty" toml:"rules,omitempty"`
}

type MaintIdxSettings struct {
	Under *int `yaml:"under,omitempty" toml:"under,omitempty"`
}

type MakezeroSettings struct {
	Always *bool `yaml:"always,omitempty" toml:"always,omitempty"`
}

type MisspellSettings struct {
	Mode        *string              `yaml:"mode,omitempty" toml:"mode,omitempty"`
	Locale      *string              `yaml:"locale,omitempty" toml:"locale,omitempty"`
	ExtraWords  []MisspellExtraWords `yaml:"extra-words,omitempty" toml:"extra-words,omitempty"`
	IgnoreRules []string             `yaml:"ignore-rules,omitempty" toml:"ignore-rules,omitempty"`
}

type MisspellExtraWords struct {
	Typo       *string `yaml:"typo,omitempty" toml:"typo,omitempty"`
	Correction *string `yaml:"correction,omitempty" toml:"correction,omitempty"`
}

type MustTagSettings struct {
	Functions []MustTagFunction `yaml:"functions,omitempty" toml:"functions,omitempty"`
}

type MustTagFunction struct {
	Name   *string `yaml:"name,omitempty" toml:"name,omitempty"`
	Tag    *string `yaml:"tag,omitempty" toml:"tag,omitempty"`
	ArgPos *int    `yaml:"arg-pos,omitempty" toml:"arg-pos,omitempty"`
}

type NakedretSettings struct {
	MaxFuncLines uint `yaml:"max-func-lines,omitempty" toml:"max-func-lines,omitempty"`
}

type NestifSettings struct {
	MinComplexity *int `yaml:"min-complexity,omitempty" toml:"min-complexity,omitempty"`
}

type NilNilSettings struct {
	DetectOpposite *bool    `yaml:"detect-opposite,omitempty" toml:"detect-opposite,omitempty"`
	CheckedTypes   []string `yaml:"checked-types,omitempty" toml:"checked-types,omitempty"`
}

type NlreturnSettings struct {
	BlockSize *int `yaml:"block-size,omitempty" toml:"block-size,omitempty"`
}

type MndSettings struct {
	Checks           []string `yaml:"checks,omitempty" toml:"checks,omitempty"`
	IgnoredNumbers   []string `yaml:"ignored-numbers,omitempty" toml:"ignored-numbers,omitempty"`
	IgnoredFiles     []string `yaml:"ignored-files,omitempty" toml:"ignored-files,omitempty"`
	IgnoredFunctions []string `yaml:"ignored-functions,omitempty" toml:"ignored-functions,omitempty"`
}

type NoLintLintSettings struct {
	RequireExplanation *bool    `yaml:"require-explanation,omitempty" toml:"require-explanation,omitempty"`
	RequireSpecific    *bool    `yaml:"require-specific,omitempty" toml:"require-specific,omitempty"`
	AllowNoExplanation []string `yaml:"allow-no-explanation,omitempty" toml:"allow-no-explanation,omitempty"`
	AllowUnused        *bool    `yaml:"allow-unused,omitempty" toml:"allow-unused,omitempty"`
}

type NoNamedReturnsSettings struct {
	ReportErrorInDefer *bool `yaml:"report-error-in-defer,omitempty" toml:"report-error-in-defer,omitempty"`
}

type ParallelTestSettings struct {
	Go                    *string `yaml:"-,omitempty" toml:"-,omitempty"`
	IgnoreMissing         *bool   `yaml:"ignore-missing,omitempty" toml:"ignore-missing,omitempty"`
	IgnoreMissingSubtests *bool   `yaml:"ignore-missing-subtests,omitempty" toml:"ignore-missing-subtests,omitempty"`
}

type PerfSprintSettings struct {
	IntegerFormat *bool `yaml:"integer-format,omitempty" toml:"integer-format,omitempty"`
	IntConversion *bool `yaml:"int-conversion,omitempty" toml:"int-conversion,omitempty"`

	ErrorFormat *bool `yaml:"error-format,omitempty" toml:"error-format,omitempty"`
	ErrError    *bool `yaml:"err-error,omitempty" toml:"err-error,omitempty"`
	ErrorF      *bool `yaml:"errorf,omitempty" toml:"errorf,omitempty"`

	StringFormat *bool `yaml:"string-format,omitempty" toml:"string-format,omitempty"`
	SprintF1     *bool `yaml:"sprintf1,omitempty" toml:"sprintf1,omitempty"`
	StrConcat    *bool `yaml:"strconcat,omitempty" toml:"strconcat,omitempty"`

	BoolFormat *bool `yaml:"bool-format,omitempty" toml:"bool-format,omitempty"`
	HexFormat  *bool `yaml:"hex-format,omitempty" toml:"hex-format,omitempty"`
}

type PreallocSettings struct {
	Simple     *bool `yaml:"simple,omitempty" toml:"simple,omitempty"`
	RangeLoops *bool `yaml:"range-loops,omitempty" toml:"range-loops,omitempty"`
	ForLoops   *bool `yaml:"for-loops,omitempty" toml:"for-loops,omitempty"`
}

type PredeclaredSettings struct {
	Ignore    []string `yaml:"ignore,omitempty" toml:"ignore,omitempty"`
	Qualified *bool    `yaml:"qualified-name,omitempty" toml:"qualified-name,omitempty"`
}

type PromlinterSettings struct {
	Strict          *bool    `yaml:"strict,omitempty" toml:"strict,omitempty"`
	DisabledLinters []string `yaml:"disabled-linters,omitempty" toml:"disabled-linters,omitempty"`
}

type ProtoGetterSettings struct {
	SkipGeneratedBy         []string `yaml:"skip-generated-by,omitempty" toml:"skip-generated-by,omitempty"`
	SkipFiles               []string `yaml:"skip-files,omitempty" toml:"skip-files,omitempty"`
	SkipAnyGenerated        *bool    `yaml:"skip-any-generated,omitempty" toml:"skip-any-generated,omitempty"`
	ReplaceFirstArgInAppend *bool    `yaml:"replace-first-arg-in-append,omitempty" toml:"replace-first-arg-in-append,omitempty"`
}

type ReassignSettings struct {
	Patterns []string `yaml:"patterns,omitempty" toml:"patterns,omitempty"`
}

type RecvcheckSettings struct {
	DisableBuiltin *bool    `yaml:"disable-builtin,omitempty" toml:"disable-builtin,omitempty"`
	Exclusions     []string `yaml:"exclusions,omitempty" toml:"exclusions,omitempty"`
}

type ReviveSettings struct {
	Go             *string           `yaml:"-,omitempty" toml:"-,omitempty"`
	MaxOpenFiles   *int              `yaml:"max-open-files,omitempty" toml:"max-open-files,omitempty"`
	Confidence     *float64          `yaml:"confidence,omitempty" toml:"confidence,omitempty"`
	Severity       *string           `yaml:"severity,omitempty" toml:"severity,omitempty"`
	EnableAllRules *bool             `yaml:"enable-all-rules,omitempty" toml:"enable-all-rules,omitempty"`
	Rules          []ReviveRule      `yaml:"rules,omitempty" toml:"rules,omitempty"`
	ErrorCode      *int              `yaml:"error-code,omitempty" toml:"error-code,omitempty"`
	WarningCode    *int              `yaml:"warning-code,omitempty" toml:"warning-code,omitempty"`
	Directives     []ReviveDirective `yaml:"directives,omitempty" toml:"directives,omitempty"`
}

type ReviveRule struct {
	Name      *string  `yaml:"name,omitempty" toml:"name,omitempty"`
	Arguments []any    `yaml:"arguments,omitempty" toml:"arguments,omitempty"`
	Severity  *string  `yaml:"severity,omitempty" toml:"severity,omitempty"`
	Disabled  *bool    `yaml:"disabled,omitempty" toml:"disabled,omitempty"`
	Exclude   []string `yaml:"exclude,omitempty" toml:"exclude,omitempty"`
}

type ReviveDirective struct {
	Name     *string `yaml:"name,omitempty" toml:"name,omitempty"`
	Severity *string `yaml:"severity,omitempty" toml:"severity,omitempty"`
}

type RowsErrCheckSettings struct {
	Packages []string `yaml:"packages,omitempty" toml:"packages,omitempty"`
}

type SlogLintSettings struct {
	NoMixedArgs    *bool    `yaml:"no-mixed-args,omitempty" toml:"no-mixed-args,omitempty"`
	KVOnly         *bool    `yaml:"kv-only,omitempty" toml:"kv-only,omitempty"`
	AttrOnly       *bool    `yaml:"attr-only,omitempty" toml:"attr-only,omitempty"`
	NoGlobal       *string  `yaml:"no-global,omitempty" toml:"no-global,omitempty"`
	Context        *string  `yaml:"context,omitempty" toml:"context,omitempty"`
	StaticMsg      *bool    `yaml:"static-msg,omitempty" toml:"static-msg,omitempty"`
	NoRawKeys      *bool    `yaml:"no-raw-keys,omitempty" toml:"no-raw-keys,omitempty"`
	KeyNamingCase  *string  `yaml:"key-naming-case,omitempty" toml:"key-naming-case,omitempty"`
	ForbiddenKeys  []string `yaml:"forbidden-keys,omitempty" toml:"forbidden-keys,omitempty"`
	ArgsOnSepLines *bool    `yaml:"args-on-sep-lines,omitempty" toml:"args-on-sep-lines,omitempty"`
}

type SpancheckSettings struct {
	Checks                   []string `yaml:"checks,omitempty" toml:"checks,omitempty"`
	IgnoreCheckSignatures    []string `yaml:"ignore-check-signatures,omitempty" toml:"ignore-check-signatures,omitempty"`
	ExtraStartSpanSignatures []string `yaml:"extra-start-span-signatures,omitempty" toml:"extra-start-span-signatures,omitempty"`
}

type StaticCheckSettings struct {
	Checks                  []string `yaml:"checks,omitempty" toml:"checks,omitempty"`
	Initialisms             []string `yaml:"initialisms,omitempty" toml:"initialisms,omitempty"`
	DotImportWhitelist      []string `yaml:"dot-import-whitelist,omitempty" toml:"dot-import-whitelist,omitempty"`
	HTTPStatusCodeWhitelist []string `yaml:"http-status-code-whitelist,omitempty" toml:"http-status-code-whitelist,omitempty"`
}

type TagAlignSettings struct {
	Align  *bool    `yaml:"align,omitempty" toml:"align,omitempty"`
	Sort   *bool    `yaml:"sort,omitempty" toml:"sort,omitempty"`
	Order  []string `yaml:"order,omitempty" toml:"order,omitempty"`
	Strict *bool    `yaml:"strict,omitempty" toml:"strict,omitempty"`
}

type TagliatelleSettings struct {
	Case TagliatelleCase `yaml:"case,omitempty" toml:"case,omitempty"`
}

type TagliatelleCase struct {
	TagliatelleBase `yaml:",inline"`
	Overrides       []TagliatelleOverrides `yaml:"overrides,omitempty" toml:"overrides,omitempty"`
}

type TagliatelleOverrides struct {
	TagliatelleBase `yaml:",inline"`
	Package         *string `yaml:"pkg,omitempty" toml:"pkg,omitempty"`
	Ignore          *bool   `yaml:"ignore,omitempty" toml:"ignore,omitempty"`
}

type TagliatelleBase struct {
	Rules         map[string]string                  `yaml:"rules,omitempty" toml:"rules,omitempty"`
	ExtendedRules map[string]TagliatelleExtendedRule `yaml:"extended-rules,omitempty" toml:"extended-rules,omitempty"`
	UseFieldName  *bool                              `yaml:"use-field-name,omitempty" toml:"use-field-name,omitempty"`
	IgnoredFields []string                           `yaml:"ignored-fields,omitempty" toml:"ignored-fields,omitempty"`
}

type TagliatelleExtendedRule struct {
	Case                *string         `yaml:"case,omitempty" toml:"case,omitempty"`
	ExtraInitialisms    *bool           `yaml:"extra-initialisms,omitempty" toml:"extra-initialisms,omitempty"`
	InitialismOverrides map[string]bool `yaml:"initialism-overrides,omitempty" toml:"initialism-overrides,omitempty"`
}

type TestifylintSettings struct {
	EnableAll        *bool    `yaml:"enable-all,omitempty" toml:"enable-all,omitempty"`
	DisableAll       *bool    `yaml:"disable-all,omitempty" toml:"disable-all,omitempty"`
	EnabledCheckers  []string `yaml:"enable,omitempty" toml:"enable,omitempty"`
	DisabledCheckers []string `yaml:"disable,omitempty" toml:"disable,omitempty"`

	BoolCompare          TestifylintBoolCompare          `yaml:"bool-compare,omitempty" toml:"bool-compare,omitempty"`
	ExpectedActual       TestifylintExpectedActual       `yaml:"expected-actual,omitempty" toml:"expected-actual,omitempty"`
	Formatter            TestifylintFormatter            `yaml:"formatter,omitempty" toml:"formatter,omitempty"`
	GoRequire            TestifylintGoRequire            `yaml:"go-require,omitempty" toml:"go-require,omitempty"`
	RequireError         TestifylintRequireError         `yaml:"require-error,omitempty" toml:"require-error,omitempty"`
	SuiteExtraAssertCall TestifylintSuiteExtraAssertCall `yaml:"suite-extra-assert-call,omitempty" toml:"suite-extra-assert-call,omitempty"`
}

type TestifylintBoolCompare struct {
	IgnoreCustomTypes *bool `yaml:"ignore-custom-types,omitempty" toml:"ignore-custom-types,omitempty"`
}

type TestifylintExpectedActual struct {
	ExpVarPattern *string `yaml:"pattern,omitempty" toml:"pattern,omitempty"`
}

type TestifylintFormatter struct {
	CheckFormatString *bool `yaml:"check-format-string,omitempty" toml:"check-format-string,omitempty"`
	RequireFFuncs     *bool `yaml:"require-f-funcs,omitempty" toml:"require-f-funcs,omitempty"`
}

type TestifylintGoRequire struct {
	IgnoreHTTPHandlers *bool `yaml:"ignore-http-handlers,omitempty" toml:"ignore-http-handlers,omitempty"`
}

type TestifylintRequireError struct {
	FnPattern *string `yaml:"fn-pattern,omitempty" toml:"fn-pattern,omitempty"`
}

type TestifylintSuiteExtraAssertCall struct {
	Mode *string `yaml:"mode,omitempty" toml:"mode,omitempty"`
}

type TestpackageSettings struct {
	SkipRegexp    *string  `yaml:"skip-regexp,omitempty" toml:"skip-regexp,omitempty"`
	AllowPackages []string `yaml:"allow-packages,omitempty" toml:"allow-packages,omitempty"`
}

type ThelperSettings struct {
	Test      ThelperOptions `yaml:"test,omitempty" toml:"test,omitempty"`
	Fuzz      ThelperOptions `yaml:"fuzz,omitempty" toml:"fuzz,omitempty"`
	Benchmark ThelperOptions `yaml:"benchmark,omitempty" toml:"benchmark,omitempty"`
	TB        ThelperOptions `yaml:"tb,omitempty" toml:"tb,omitempty"`
}

type ThelperOptions struct {
	First *bool `yaml:"first,omitempty" toml:"first,omitempty"`
	Name  *bool `yaml:"name,omitempty" toml:"name,omitempty"`
	Begin *bool `yaml:"begin,omitempty" toml:"begin,omitempty"`
}

type TenvSettings struct {
	All *bool `yaml:"all,omitempty" toml:"all,omitempty"`
}

type UseStdlibVarsSettings struct {
	HTTPMethod         *bool `yaml:"http-method,omitempty" toml:"http-method,omitempty"`
	HTTPStatusCode     *bool `yaml:"http-status-code,omitempty" toml:"http-status-code,omitempty"`
	TimeWeekday        *bool `yaml:"time-weekday,omitempty" toml:"time-weekday,omitempty"`
	TimeMonth          *bool `yaml:"time-month,omitempty" toml:"time-month,omitempty"`
	TimeLayout         *bool `yaml:"time-layout,omitempty" toml:"time-layout,omitempty"`
	CryptoHash         *bool `yaml:"crypto-hash,omitempty" toml:"crypto-hash,omitempty"`
	DefaultRPCPath     *bool `yaml:"default-rpc-path,omitempty" toml:"default-rpc-path,omitempty"`
	SQLIsolationLevel  *bool `yaml:"sql-isolation-level,omitempty" toml:"sql-isolation-level,omitempty"`
	TLSSignatureScheme *bool `yaml:"tls-signature-scheme,omitempty" toml:"tls-signature-scheme,omitempty"`
	ConstantKind       *bool `yaml:"constant-kind,omitempty" toml:"constant-kind,omitempty"`
}

type UseTestingSettings struct {
	ContextBackground *bool `yaml:"context-background,omitempty" toml:"context-background,omitempty"`
	ContextTodo       *bool `yaml:"context-todo,omitempty" toml:"context-todo,omitempty"`
	OSChdir           *bool `yaml:"os-chdir,omitempty" toml:"os-chdir,omitempty"`
	OSMkdirTemp       *bool `yaml:"os-mkdir-temp,omitempty" toml:"os-mkdir-temp,omitempty"`
	OSSetenv          *bool `yaml:"os-setenv,omitempty" toml:"os-setenv,omitempty"`
	OSTempDir         *bool `yaml:"os-temp-dir,omitempty" toml:"os-temp-dir,omitempty"`
	OSCreateTemp      *bool `yaml:"os-create-temp,omitempty" toml:"os-create-temp,omitempty"`
}

type UnconvertSettings struct {
	FastMath *bool `yaml:"fast-math,omitempty" toml:"fast-math,omitempty"`
	Safe     *bool `yaml:"safe,omitempty" toml:"safe,omitempty"`
}

type UnparamSettings struct {
	CheckExported *bool `yaml:"check-exported,omitempty" toml:"check-exported,omitempty"`
}

type UnusedSettings struct {
	FieldWritesAreUses     *bool `yaml:"field-writes-are-uses,omitempty" toml:"field-writes-are-uses,omitempty"`
	PostStatementsAreReads *bool `yaml:"post-statements-are-reads,omitempty" toml:"post-statements-are-reads,omitempty"`
	ExportedFieldsAreUsed  *bool `yaml:"exported-fields-are-used,omitempty" toml:"exported-fields-are-used,omitempty"`
	ParametersAreUsed      *bool `yaml:"parameters-are-used,omitempty" toml:"parameters-are-used,omitempty"`
	LocalVariablesAreUsed  *bool `yaml:"local-variables-are-used,omitempty" toml:"local-variables-are-used,omitempty"`
	GeneratedIsUsed        *bool `yaml:"generated-is-used,omitempty" toml:"generated-is-used,omitempty"`
}

type VarnamelenSettings struct {
	MaxDistance        *int     `yaml:"max-distance,omitempty" toml:"max-distance,omitempty"`
	MinNameLength      *int     `yaml:"min-name-length,omitempty" toml:"min-name-length,omitempty"`
	CheckReceiver      *bool    `yaml:"check-receiver,omitempty" toml:"check-receiver,omitempty"`
	CheckReturn        *bool    `yaml:"check-return,omitempty" toml:"check-return,omitempty"`
	CheckTypeParam     *bool    `yaml:"check-type-param,omitempty" toml:"check-type-param,omitempty"`
	IgnoreNames        []string `yaml:"ignore-names,omitempty" toml:"ignore-names,omitempty"`
	IgnoreTypeAssertOk *bool    `yaml:"ignore-type-assert-ok,omitempty" toml:"ignore-type-assert-ok,omitempty"`
	IgnoreMapIndexOk   *bool    `yaml:"ignore-map-index-ok,omitempty" toml:"ignore-map-index-ok,omitempty"`
	IgnoreChanRecvOk   *bool    `yaml:"ignore-chan-recv-ok,omitempty" toml:"ignore-chan-recv-ok,omitempty"`
	IgnoreDecls        []string `yaml:"ignore-decls,omitempty" toml:"ignore-decls,omitempty"`
}

type WhitespaceSettings struct {
	MultiIf   *bool `yaml:"multi-if,omitempty" toml:"multi-if,omitempty"`
	MultiFunc *bool `yaml:"multi-func,omitempty" toml:"multi-func,omitempty"`
}

type WrapcheckSettings struct {
	ExtraIgnoreSigs        []string `yaml:"extra-ignore-sigs,omitempty" toml:"extra-ignore-sigs,omitempty"`
	IgnoreSigs             []string `yaml:"ignore-sigs,omitempty" toml:"ignore-sigs,omitempty"`
	IgnoreSigRegexps       []string `yaml:"ignore-sig-regexps,omitempty" toml:"ignore-sig-regexps,omitempty"`
	IgnorePackageGlobs     []string `yaml:"ignore-package-globs,omitempty" toml:"ignore-package-globs,omitempty"`
	IgnoreInterfaceRegexps []string `yaml:"ignore-interface-regexps,omitempty" toml:"ignore-interface-regexps,omitempty"`
}

type WSLSettings struct {
	StrictAppend                     *bool    `yaml:"strict-append,omitempty" toml:"strict-append,omitempty"`
	AllowAssignAndCallCuddle         *bool    `yaml:"allow-assign-and-call,omitempty" toml:"allow-assign-and-call,omitempty"`
	AllowAssignAndAnythingCuddle     *bool    `yaml:"allow-assign-and-anything,omitempty" toml:"allow-assign-and-anything,omitempty"`
	AllowMultiLineAssignCuddle       *bool    `yaml:"allow-multiline-assign,omitempty" toml:"allow-multiline-assign,omitempty"`
	ForceCaseTrailingWhitespaceLimit *int     `yaml:"force-case-trailing-whitespace,omitempty" toml:"force-case-trailing-whitespace,omitempty"`
	AllowTrailingComment             *bool    `yaml:"allow-trailing-comment,omitempty" toml:"allow-trailing-comment,omitempty"`
	AllowSeparatedLeadingComment     *bool    `yaml:"allow-separated-leading-comment,omitempty" toml:"allow-separated-leading-comment,omitempty"`
	AllowCuddleDeclaration           *bool    `yaml:"allow-cuddle-declarations,omitempty" toml:"allow-cuddle-declarations,omitempty"`
	AllowCuddleWithCalls             []string `yaml:"allow-cuddle-with-calls,omitempty" toml:"allow-cuddle-with-calls,omitempty"`
	AllowCuddleWithRHS               []string `yaml:"allow-cuddle-with-rhs,omitempty" toml:"allow-cuddle-with-rhs,omitempty"`
	ForceCuddleErrCheckAndAssign     *bool    `yaml:"force-err-cuddling,omitempty" toml:"force-err-cuddling,omitempty"`
	ErrorVariableNames               []string `yaml:"error-variable-names,omitempty" toml:"error-variable-names,omitempty"`
	ForceExclusiveShortDeclarations  *bool    `yaml:"force-short-decl-cuddling,omitempty" toml:"force-short-decl-cuddling,omitempty"`
}

type CustomLinterSettings struct {
	Type *string `yaml:"type,omitempty" toml:"type,omitempty"`

	Path *string `yaml:"path,omitempty" toml:"path,omitempty"`

	Description *string `yaml:"description,omitempty" toml:"description,omitempty"`

	OriginalURL *string `yaml:"original-url,omitempty" toml:"original-url,omitempty"`

	Settings any `yaml:"settings,omitempty" toml:"settings,omitempty"`
}
