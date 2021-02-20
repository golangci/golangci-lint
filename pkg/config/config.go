package config

import (
	"errors"
	"fmt"
	"regexp"
	"time"
)

const (
	OutFormatJSON              = "json"
	OutFormatLineNumber        = "line-number"
	OutFormatColoredLineNumber = "colored-line-number"
	OutFormatTab               = "tab"
	OutFormatCheckstyle        = "checkstyle"
	OutFormatCodeClimate       = "code-climate"
	OutFormatJunitXML          = "junit-xml"
	OutFormatGithubActions     = "github-actions"
)

var OutFormats = []string{
	OutFormatColoredLineNumber,
	OutFormatLineNumber,
	OutFormatJSON,
	OutFormatTab,
	OutFormatCheckstyle,
	OutFormatCodeClimate,
	OutFormatJunitXML,
	OutFormatGithubActions,
}

type ExcludePattern struct {
	ID      string
	Pattern string
	Linter  string
	Why     string
}

var DefaultExcludePatterns = []ExcludePattern{
	{
		ID: "EXC0001",
		Pattern: "Error return value of .((os\\.)?std(out|err)\\..*|.*Close" +
			"|.*Flush|os\\.Remove(All)?|.*print(f|ln)?|os\\.(Un)?Setenv). is not checked",
		Linter: "errcheck",
		Why:    "Almost all programs ignore errors on these functions and in most cases it's ok",
	},
	{
		ID: "EXC0002",
		Pattern: "(comment on exported (method|function|type|const)|" +
			"should have( a package)? comment|comment should be of the form)",
		Linter: "golint",
		Why:    "Annoying issue about not having a comment. The rare codebase has such comments",
	},
	{
		ID:      "EXC0003",
		Pattern: "func name will be used as test\\.Test.* by other packages, and that stutters; consider calling this",
		Linter:  "golint",
		Why:     "False positive when tests are defined in package 'test'",
	},
	{
		ID:      "EXC0004",
		Pattern: "(possible misuse of unsafe.Pointer|should have signature)",
		Linter:  "govet",
		Why:     "Common false positives",
	},
	{
		ID:      "EXC0005",
		Pattern: "ineffective break statement. Did you mean to break out of the outer loop",
		Linter:  "staticcheck",
		Why:     "Developers tend to write in C-style with an explicit 'break' in a 'switch', so it's ok to ignore",
	},
	{
		ID:      "EXC0006",
		Pattern: "Use of unsafe calls should be audited",
		Linter:  "gosec",
		Why:     "Too many false-positives on 'unsafe' usage",
	},
	{
		ID:      "EXC0007",
		Pattern: "Subprocess launch(ed with variable|ing should be audited)",
		Linter:  "gosec",
		Why:     "Too many false-positives for parametrized shell calls",
	},
	{
		ID:      "EXC0008",
		Pattern: "(G104|G307)",
		Linter:  "gosec",
		Why:     "Duplicated errcheck checks",
	},
	{
		ID:      "EXC0009",
		Pattern: "(Expect directory permissions to be 0750 or less|Expect file permissions to be 0600 or less)",
		Linter:  "gosec",
		Why:     "Too many issues in popular repos",
	},
	{
		ID:      "EXC0010",
		Pattern: "Potential file inclusion via variable",
		Linter:  "gosec",
		Why:     "False positive is triggered by 'src, err := ioutil.ReadFile(filename)'",
	},
	{
		ID: "EXC0011",
		Pattern: "(comment on exported (method|function|type|const)|" +
			"should have( a package)? comment|comment should be of the form)",
		Linter: "stylecheck",
		Why:    "Annoying issue about not having a comment. The rare codebase has such comments",
	},
}

func GetDefaultExcludePatternsStrings() []string {
	ret := make([]string, len(DefaultExcludePatterns))
	for i, p := range DefaultExcludePatterns {
		ret[i] = p.Pattern
	}
	return ret
}

func GetExcludePatterns(include []string) []ExcludePattern {
	includeMap := make(map[string]bool, len(include))
	for _, inc := range include {
		includeMap[inc] = true
	}

	var ret []ExcludePattern
	for _, p := range DefaultExcludePatterns {
		if !includeMap[p.ID] {
			ret = append(ret, p)
		}
	}

	return ret
}

type Run struct {
	IsVerbose           bool `mapstructure:"verbose"`
	Silent              bool
	CPUProfilePath      string
	MemProfilePath      string
	TracePath           string
	Concurrency         int
	PrintResourcesUsage bool `mapstructure:"print-resources-usage"`

	Config   string
	NoConfig bool

	Args []string

	BuildTags           []string `mapstructure:"build-tags"`
	ModulesDownloadMode string   `mapstructure:"modules-download-mode"`

	ExitCodeIfIssuesFound int  `mapstructure:"issues-exit-code"`
	AnalyzeTests          bool `mapstructure:"tests"`

	// Deprecated: Deadline exists for historical compatibility
	// and should not be used. To set run timeout use Timeout instead.
	Deadline time.Duration
	Timeout  time.Duration

	PrintVersion       bool
	SkipFiles          []string `mapstructure:"skip-files"`
	SkipDirs           []string `mapstructure:"skip-dirs"`
	UseDefaultSkipDirs bool     `mapstructure:"skip-dirs-use-default"`

	AllowParallelRunners bool `mapstructure:"allow-parallel-runners"`
	AllowSerialRunners   bool `mapstructure:"allow-serial-runners"`
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

type Cyclop struct {
	MaxComplexity  int     `mapstructure:"max-complexity"`
	PackageAverage float64 `mapstructure:"package-average"`
	SkipTests      bool    `mapstructure:"skip-tests"`
}

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

type CustomLinterSettings struct {
	Path        string
	Description string
	OriginalURL string `mapstructure:"original-url"`
}

type Linters struct {
	Enable     []string
	Disable    []string
	EnableAll  bool `mapstructure:"enable-all"`
	DisableAll bool `mapstructure:"disable-all"`
	Fast       bool

	Presets []string
}

type BaseRule struct {
	Linters []string
	Path    string
	Text    string
	Source  string
}

func (b BaseRule) Validate(minConditionsCount int) error {
	if err := validateOptionalRegex(b.Path); err != nil {
		return fmt.Errorf("invalid path regex: %v", err)
	}
	if err := validateOptionalRegex(b.Text); err != nil {
		return fmt.Errorf("invalid text regex: %v", err)
	}
	if err := validateOptionalRegex(b.Source); err != nil {
		return fmt.Errorf("invalid source regex: %v", err)
	}
	nonBlank := 0
	if len(b.Linters) > 0 {
		nonBlank++
	}
	if b.Path != "" {
		nonBlank++
	}
	if b.Text != "" {
		nonBlank++
	}
	if b.Source != "" {
		nonBlank++
	}
	if nonBlank < minConditionsCount {
		return fmt.Errorf("at least %d of (text, source, path, linters) should be set", minConditionsCount)
	}
	return nil
}

const excludeRuleMinConditionsCount = 2

type ExcludeRule struct {
	BaseRule `mapstructure:",squash"`
}

func validateOptionalRegex(value string) error {
	if value == "" {
		return nil
	}
	_, err := regexp.Compile(value)
	return err
}

func (e ExcludeRule) Validate() error {
	return e.BaseRule.Validate(excludeRuleMinConditionsCount)
}

const severityRuleMinConditionsCount = 1

type SeverityRule struct {
	BaseRule `mapstructure:",squash"`
	Severity string
}

func (s *SeverityRule) Validate() error {
	return s.BaseRule.Validate(severityRuleMinConditionsCount)
}

type Issues struct {
	IncludeDefaultExcludes []string      `mapstructure:"include"`
	ExcludeCaseSensitive   bool          `mapstructure:"exclude-case-sensitive"`
	ExcludePatterns        []string      `mapstructure:"exclude"`
	ExcludeRules           []ExcludeRule `mapstructure:"exclude-rules"`
	UseDefaultExcludes     bool          `mapstructure:"exclude-use-default"`

	MaxIssuesPerLinter int `mapstructure:"max-issues-per-linter"`
	MaxSameIssues      int `mapstructure:"max-same-issues"`

	DiffFromRevision  string `mapstructure:"new-from-rev"`
	DiffPatchFilePath string `mapstructure:"new-from-patch"`
	Diff              bool   `mapstructure:"new"`

	NeedFix bool `mapstructure:"fix"`
}

type Severity struct {
	Default       string         `mapstructure:"default-severity"`
	CaseSensitive bool           `mapstructure:"case-sensitive"`
	Rules         []SeverityRule `mapstructure:"rules"`
}

type Version struct {
	Format string `mapstructure:"format"`
}

type Config struct {
	Run Run

	Output struct {
		Format              string
		Color               string
		PrintIssuedLine     bool   `mapstructure:"print-issued-lines"`
		PrintLinterName     bool   `mapstructure:"print-linter-name"`
		UniqByLine          bool   `mapstructure:"uniq-by-line"`
		SortResults         bool   `mapstructure:"sort-results"`
		PrintWelcomeMessage bool   `mapstructure:"print-welcome"`
		PathPrefix          string `mapstructure:"path-prefix"`
	}

	LintersSettings LintersSettings `mapstructure:"linters-settings"`
	Linters         Linters
	Issues          Issues
	Severity        Severity
	Version         Version

	InternalTest bool // Option is used only for testing golangci-lint code, don't use it
}

func NewDefault() *Config {
	return &Config{
		LintersSettings: defaultLintersSettings,
	}
}
