package revive

import (
	"bytes"
	"cmp"
	"fmt"
	"go/token"
	"os"
	"reflect"
	"strings"
	"sync"

	"github.com/BurntSushi/toml"
	hcversion "github.com/hashicorp/go-version"
	reviveConfig "github.com/mgechev/revive/config"
	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
	"github.com/golangci/golangci-lint/v2/pkg/golinters/internal"
	"github.com/golangci/golangci-lint/v2/pkg/lint/linter"
	"github.com/golangci/golangci-lint/v2/pkg/logutils"
	"github.com/golangci/golangci-lint/v2/pkg/result"
)

const linterName = "revive"

var (
	debugf  = logutils.Debug(logutils.DebugKeyRevive)
	isDebug = logutils.HaveDebugTag(logutils.DebugKeyRevive)
)

func New(settings *config.ReviveSettings) *goanalysis.Linter {
	var mu sync.Mutex
	var resIssues []*goanalysis.Issue

	analyzer := &analysis.Analyzer{
		Name: linterName,
		Doc:  "Fast, configurable, extensible, flexible, and beautiful linter for Go. Drop-in replacement of golint.",
		Run:  goanalysis.DummyRun,
	}

	return goanalysis.
		NewLinterFromAnalyzer(analyzer).
		WithContextSetter(func(lintCtx *linter.Context) {
			w, err := newWrapper(settings)
			if err != nil {
				lintCtx.Log.Errorf("setup revive: %v", err)
				return
			}

			analyzer.Run = func(pass *analysis.Pass) (any, error) {
				issues, err := w.run(pass)
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
		}).
		WithIssuesReporter(func(*linter.Context) []*goanalysis.Issue {
			return resIssues
		}).
		WithLoadMode(goanalysis.LoadModeSyntax)
}

type wrapper struct {
	revive       lint.Linter
	lintingRules []lint.Rule
	conf         *lint.Config
}

func newWrapper(settings *config.ReviveSettings) (*wrapper, error) {
	conf, err := getConfig(settings)
	if err != nil {
		return nil, err
	}

	displayRules(conf)

	conf.GoVersion, err = hcversion.NewVersion(settings.Go)
	if err != nil {
		return nil, err
	}

	lintingRules, err := reviveConfig.GetLintingRules(conf, []lint.Rule{})
	if err != nil {
		return nil, err
	}

	return &wrapper{
		revive:       lint.New(os.ReadFile, settings.MaxOpenFiles),
		lintingRules: lintingRules,
		conf:         conf,
	}, nil
}

func (w *wrapper) run(pass *analysis.Pass) ([]*goanalysis.Issue, error) {
	packages := [][]string{internal.GetGoFileNames(pass)}

	failures, err := w.revive.Lint(packages, w.lintingRules, *w.conf)
	if err != nil {
		return nil, err
	}

	var issues []*goanalysis.Issue
	for failure := range failures {
		if failure.Confidence < w.conf.Confidence {
			continue
		}

		issues = append(issues, w.toIssue(pass, &failure))
	}

	return issues, nil
}

func (w *wrapper) toIssue(pass *analysis.Pass, failure *lint.Failure) *goanalysis.Issue {
	lineRangeTo := failure.Position.End.Line
	if failure.RuleName == (&rule.ExportedRule{}).Name() {
		lineRangeTo = failure.Position.Start.Line
	}

	issue := &result.Issue{
		Severity: string(failure.SeverityFor(w.conf)),
		Text:     fmt.Sprintf("%s: %s", failure.RuleName, failure.Failure),
		Pos: token.Position{
			Filename: failure.Position.Start.Filename,
			Line:     failure.Position.Start.Line,
			Offset:   failure.Position.Start.Offset,
			Column:   failure.Position.Start.Column,
		},
		LineRange: &result.Range{
			From: failure.Position.Start.Line,
			To:   lineRangeTo,
		},
		FromLinter: linterName,
	}

	if failure.ReplacementLine != "" {
		f := pass.Fset.File(token.Pos(failure.Position.Start.Offset))

		// Skip cgo files because the positions are wrong.
		if failure.Filename() == f.Name() {
			issue.SuggestedFixes = []analysis.SuggestedFix{{
				TextEdits: []analysis.TextEdit{{
					Pos: f.LineStart(failure.Position.Start.Line),
					End: goanalysis.EndOfLinePos(f, failure.Position.End.Line),
					// ReplacementLine doesn't contain the full line (missing newline), so we have to add a newline.
					// Also `failure.Position.End.Offset` is at the end of the node but not the line.
					NewText: []byte(failure.ReplacementLine + "\n"),
				}},
			}}
		}
	}

	return goanalysis.NewIssue(issue, pass)
}

// This function mimics the GetConfig function of revive.
// This allows to get default values and right types.
// https://github.com/golangci/golangci-lint/issues/1745
// TODO: update the following links:
// https://github.com/mgechev/revive/blob/v1.13.0/config/config.go#L249
// https://github.com/mgechev/revive/blob/v1.13.0/config/config.go#L198-L204
func getConfig(cfg *config.ReviveSettings) (*lint.Config, error) {
	conf := reviveConfig.Default()

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

	return conf, nil
}

func createConfigMap(cfg *config.ReviveSettings) map[string]any {
	const severity = "severity"

	rawRoot := map[string]any{
		"confidence":         cfg.Confidence,
		severity:             cfg.Severity,
		"errorCode":          cfg.ErrorCode,
		"warningCode":        cfg.WarningCode,
		"enableAllRules":     cfg.EnableAllRules,
		"enableDefaultRules": cfg.EnableDefaultRules,

		// Should be managed with `linters.exclusions.generated`.
		"ignoreGeneratedHeader": false,
	}

	rawDirectives := map[string]map[string]any{}
	for _, directive := range cfg.Directives {
		rawDirectives[directive.Name] = map[string]any{
			severity: directive.Severity,
		}
	}

	if len(rawDirectives) > 0 {
		rawRoot["directive"] = rawDirectives
	}

	rawRules := map[string]map[string]any{}
	for _, s := range cfg.Rules {
		rawRules[s.Name] = map[string]any{
			severity:    s.Severity,
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

func normalizeConfig(cfg *lint.Config) {
	// NOTE(ldez): this custom section for golangci-lint should be kept.
	// ---
	cfg.Confidence = cmp.Or(cfg.Confidence, reviveConfig.DefaultConfidence)
	cfg.Severity = cmp.Or(cfg.Severity, lint.SeverityWarning)
	// ---

	reviveConfig.Normalize(cfg)
}

func displayRules(conf *lint.Config) {
	if !isDebug {
		return
	}

	allRules := reviveConfig.AllRuleNames()
	defaultRules := reviveConfig.DefaultRuleNames()
	enabledRules := reviveConfig.EnabledRuleNames(conf)

	debugf("All available rules (%d): %s.", len(allRules), strings.Join(allRules, ", "))
	debugf("Default rules (%d): %s.", len(defaultRules), strings.Join(defaultRules, ", "))
	debugf("Enabled by config rules (%d): %s.", len(enabledRules), strings.Join(enabledRules, ", "))

	debugf("revive configuration: %#v", conf)
}
