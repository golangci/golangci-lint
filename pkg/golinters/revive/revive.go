package revive

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/token"
	"os"
	"reflect"
	"sync"

	"github.com/BurntSushi/toml"
	hcversion "github.com/hashicorp/go-version"
	reviveConfig "github.com/mgechev/revive/config"
	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
	"github.com/golangci/golangci-lint/pkg/golinters/internal"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/golangci/golangci-lint/pkg/result"
)

const linterName = "revive"

var debugf = logutils.Debug(logutils.DebugKeyRevive)

// jsonObject defines a JSON object of a failure
type jsonObject struct {
	Severity     lint.Severity
	lint.Failure `json:",inline"`
}

func New(settings *config.ReviveSettings) *goanalysis.Linter {
	var mu sync.Mutex
	var resIssues []goanalysis.Issue

	analyzer := &analysis.Analyzer{
		Name: goanalysis.TheOnlyAnalyzerName,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
		Run:  goanalysis.DummyRun,
	}

	return goanalysis.NewLinter(
		linterName,
		"Fast, configurable, extensible, flexible, and beautiful linter for Go. Drop-in replacement of golint.",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithContextSetter(func(lintCtx *linter.Context) {
		w, err := newWrapper(settings)
		if err != nil {
			lintCtx.Log.Errorf("setup revive: %v", err)
			return
		}

		analyzer.Run = func(pass *analysis.Pass) (any, error) {
			issues, err := w.run(lintCtx, pass)
			if err != nil {
				return nil, err
			}

			if len(issues) == 0 {
				return nil, nil
			}

			mu.Lock()
			resIssues = append(resIssues, issues...)
			mu.Unlock()

			return nil, nil
		}
	}).WithIssuesReporter(func(*linter.Context) []goanalysis.Issue {
		return resIssues
	}).WithLoadMode(goanalysis.LoadModeSyntax)
}

type wrapper struct {
	revive       lint.Linter
	formatter    lint.Formatter
	lintingRules []lint.Rule
	conf         *lint.Config
}

func newWrapper(settings *config.ReviveSettings) (*wrapper, error) {
	conf, err := getConfig(settings)
	if err != nil {
		return nil, err
	}

	conf.GoVersion, err = hcversion.NewVersion(settings.Go)
	if err != nil {
		return nil, err
	}

	formatter, err := reviveConfig.GetFormatter("json")
	if err != nil {
		return nil, err
	}

	lintingRules, err := reviveConfig.GetLintingRules(conf, []lint.Rule{})
	if err != nil {
		return nil, err
	}

	return &wrapper{
		revive:       lint.New(os.ReadFile, settings.MaxOpenFiles),
		formatter:    formatter,
		lintingRules: lintingRules,
		conf:         conf,
	}, nil
}

func (w *wrapper) run(lintCtx *linter.Context, pass *analysis.Pass) ([]goanalysis.Issue, error) {
	packages := [][]string{internal.GetFileNames(pass)}

	failures, err := w.revive.Lint(packages, w.lintingRules, *w.conf)
	if err != nil {
		return nil, err
	}

	formatChan := make(chan lint.Failure)
	exitChan := make(chan bool)

	var output string
	go func() {
		output, err = w.formatter.Format(formatChan, *w.conf)
		if err != nil {
			lintCtx.Log.Errorf("Format error: %v", err)
		}
		exitChan <- true
	}()

	for f := range failures {
		if f.Confidence < w.conf.Confidence {
			continue
		}

		formatChan <- f
	}

	close(formatChan)
	<-exitChan

	var results []jsonObject
	err = json.Unmarshal([]byte(output), &results)
	if err != nil {
		return nil, err
	}

	var issues []goanalysis.Issue
	for i := range results {
		issues = append(issues, toIssue(pass, &results[i]))
	}

	return issues, nil
}

func toIssue(pass *analysis.Pass, object *jsonObject) goanalysis.Issue {
	lineRangeTo := object.Position.End.Line
	if object.RuleName == (&rule.ExportedRule{}).Name() {
		lineRangeTo = object.Position.Start.Line
	}

	return goanalysis.NewIssue(&result.Issue{
		Severity: string(object.Severity),
		Text:     fmt.Sprintf("%s: %s", object.RuleName, object.Failure.Failure),
		Pos: token.Position{
			Filename: object.Position.Start.Filename,
			Line:     object.Position.Start.Line,
			Offset:   object.Position.Start.Offset,
			Column:   object.Position.Start.Column,
		},
		LineRange: &result.Range{
			From: object.Position.Start.Line,
			To:   lineRangeTo,
		},
		FromLinter: linterName,
	}, pass)
}

// This function mimics the GetConfig function of revive.
// This allows to get default values and right types.
// https://github.com/golangci/golangci-lint/issues/1745
// https://github.com/mgechev/revive/blob/v1.5.0/config/config.go#L220
// https://github.com/mgechev/revive/blob/v1.5.0/config/config.go#L172-L178
func getConfig(cfg *config.ReviveSettings) (*lint.Config, error) {
	conf := defaultConfig()

	// Since the Go version is dynamic, this value must be neutralized in order to compare with a "zero value" of the configuration structure.
	zero := &config.ReviveSettings{Go: cfg.Go}

	if !reflect.DeepEqual(cfg, zero) {
		rawRoot := createConfigMap(cfg)
		buf := bytes.NewBuffer(nil)

		err := toml.NewEncoder(buf).Encode(rawRoot)
		if err != nil {
			return nil, fmt.Errorf("failed to encode configuration: %w", err)
		}

		conf = &lint.Config{}
		_, err = toml.NewDecoder(buf).Decode(conf)
		if err != nil {
			return nil, fmt.Errorf("failed to decode configuration: %w", err)
		}
	}

	normalizeConfig(conf)

	for k, r := range conf.Rules {
		err := r.Initialize()
		if err != nil {
			return nil, fmt.Errorf("error in config of rule %q: %w", k, err)
		}
		conf.Rules[k] = r
	}

	debugf("revive configuration: %#v", conf)

	return conf, nil
}

func createConfigMap(cfg *config.ReviveSettings) map[string]any {
	rawRoot := map[string]any{
		"ignoreGeneratedHeader": cfg.IgnoreGeneratedHeader,
		"confidence":            cfg.Confidence,
		"severity":              cfg.Severity,
		"errorCode":             cfg.ErrorCode,
		"warningCode":           cfg.WarningCode,
		"enableAllRules":        cfg.EnableAllRules,
	}

	rawDirectives := map[string]map[string]any{}
	for _, directive := range cfg.Directives {
		rawDirectives[directive.Name] = map[string]any{
			"severity": directive.Severity,
		}
	}

	if len(rawDirectives) > 0 {
		rawRoot["directive"] = rawDirectives
	}

	rawRules := map[string]map[string]any{}
	for _, s := range cfg.Rules {
		rawRules[s.Name] = map[string]any{
			"severity":  s.Severity,
			"arguments": safeTomlSlice(s.Arguments),
			"disabled":  s.Disabled,
			"exclude":   s.Exclude,
		}
	}

	if len(rawRules) > 0 {
		rawRoot["rule"] = rawRules
	}

	return rawRoot
}

func safeTomlSlice(r []any) []any {
	if len(r) == 0 {
		return nil
	}

	if _, ok := r[0].(map[any]any); !ok {
		return r
	}

	var typed []any
	for _, elt := range r {
		item := map[string]any{}
		for k, v := range elt.(map[any]any) {
			item[k.(string)] = v
		}

		typed = append(typed, item)
	}

	return typed
}

// This element is not exported by revive, so we need copy the code.
// Extracted from https://github.com/mgechev/revive/blob/v1.5.0/config/config.go#L16
var defaultRules = []lint.Rule{
	&rule.VarDeclarationsRule{},
	&rule.PackageCommentsRule{},
	&rule.DotImportsRule{},
	&rule.BlankImportsRule{},
	&rule.ExportedRule{},
	&rule.VarNamingRule{},
	&rule.IndentErrorFlowRule{},
	&rule.RangeRule{},
	&rule.ErrorfRule{},
	&rule.ErrorNamingRule{},
	&rule.ErrorStringsRule{},
	&rule.ReceiverNamingRule{},
	&rule.IncrementDecrementRule{},
	&rule.ErrorReturnRule{},
	&rule.UnexportedReturnRule{},
	&rule.TimeNamingRule{},
	&rule.ContextKeysType{},
	&rule.ContextAsArgumentRule{},
	&rule.EmptyBlockRule{},
	&rule.SuperfluousElseRule{},
	&rule.UnusedParamRule{},
	&rule.UnreachableCodeRule{},
	&rule.RedefinesBuiltinIDRule{},
}

var allRules = append([]lint.Rule{
	&rule.ArgumentsLimitRule{},
	&rule.CyclomaticRule{},
	&rule.FileHeaderRule{},
	&rule.ConfusingNamingRule{},
	&rule.GetReturnRule{},
	&rule.ModifiesParamRule{},
	&rule.ConfusingResultsRule{},
	&rule.DeepExitRule{},
	&rule.AddConstantRule{},
	&rule.FlagParamRule{},
	&rule.UnnecessaryStmtRule{},
	&rule.StructTagRule{},
	&rule.ModifiesValRecRule{},
	&rule.ConstantLogicalExprRule{},
	&rule.BoolLiteralRule{},
	&rule.ImportsBlocklistRule{},
	&rule.FunctionResultsLimitRule{},
	&rule.MaxPublicStructsRule{},
	&rule.RangeValInClosureRule{},
	&rule.RangeValAddress{},
	&rule.WaitGroupByValueRule{},
	&rule.AtomicRule{},
	&rule.EmptyLinesRule{},
	&rule.LineLengthLimitRule{},
	&rule.CallToGCRule{},
	&rule.DuplicatedImportsRule{},
	&rule.ImportShadowingRule{},
	&rule.BareReturnRule{},
	&rule.UnusedReceiverRule{},
	&rule.UnhandledErrorRule{},
	&rule.CognitiveComplexityRule{},
	&rule.StringOfIntRule{},
	&rule.StringFormatRule{},
	&rule.EarlyReturnRule{},
	&rule.UnconditionalRecursionRule{},
	&rule.IdenticalBranchesRule{},
	&rule.DeferRule{},
	&rule.UnexportedNamingRule{},
	&rule.FunctionLength{},
	&rule.NestedStructs{},
	&rule.UselessBreak{},
	&rule.UncheckedTypeAssertionRule{},
	&rule.TimeEqualRule{},
	&rule.BannedCharsRule{},
	&rule.OptimizeOperandsOrderRule{},
	&rule.UseAnyRule{},
	&rule.DataRaceRule{},
	&rule.CommentSpacingsRule{},
	&rule.IfReturnRule{},
	&rule.RedundantImportAlias{},
	&rule.ImportAliasNamingRule{},
	&rule.EnforceMapStyleRule{},
	&rule.EnforceRepeatedArgTypeStyleRule{},
	&rule.EnforceSliceStyleRule{},
	&rule.MaxControlNestingRule{},
	&rule.CommentsDensityRule{},
	&rule.FileLengthLimitRule{},
	&rule.FilenameFormatRule{},
}, defaultRules...)

const defaultConfidence = 0.8

// This element is not exported by revive, so we need copy the code.
// Extracted from https://github.com/mgechev/revive/blob/v1.5.0/config/config.go#L183
func normalizeConfig(cfg *lint.Config) {
	// NOTE(ldez): this custom section for golangci-lint should be kept.
	// ---
	if cfg.Confidence == 0 {
		cfg.Confidence = defaultConfidence
	}
	if cfg.Severity == "" {
		cfg.Severity = lint.SeverityWarning
	}
	// ---

	if len(cfg.Rules) == 0 {
		cfg.Rules = map[string]lint.RuleConfig{}
	}
	if cfg.EnableAllRules {
		// Add to the configuration all rules not yet present in it
		for _, r := range allRules {
			ruleName := r.Name()
			_, alreadyInConf := cfg.Rules[ruleName]
			if alreadyInConf {
				continue
			}
			// Add the rule with an empty conf for
			cfg.Rules[ruleName] = lint.RuleConfig{}
		}
	}

	severity := cfg.Severity
	if severity != "" {
		for k, v := range cfg.Rules {
			if v.Severity == "" {
				v.Severity = severity
			}
			cfg.Rules[k] = v
		}
		for k, v := range cfg.Directives {
			if v.Severity == "" {
				v.Severity = severity
			}
			cfg.Directives[k] = v
		}
	}
}

// This element is not exported by revive, so we need copy the code.
// Extracted from https://github.com/mgechev/revive/blob/v1.5.0/config/config.go#L252
func defaultConfig() *lint.Config {
	defaultConfig := lint.Config{
		Confidence: defaultConfidence,
		Severity:   lint.SeverityWarning,
		Rules:      map[string]lint.RuleConfig{},
	}
	for _, r := range defaultRules {
		defaultConfig.Rules[r.Name()] = lint.RuleConfig{}
	}
	return &defaultConfig
}
