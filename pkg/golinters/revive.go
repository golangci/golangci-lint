package golinters

import (
	"encoding/json"
	"fmt"
	"go/token"
	"io/ioutil"

	"github.com/mgechev/dots"
	reviveConfig "github.com/mgechev/revive/config"
	"github.com/mgechev/revive/lint"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

const reviveName = "revive"

// jsonObject defines a JSON object of an failure
type jsonObject struct {
	Severity     lint.Severity
	lint.Failure `json:",inline"`
}

// NewNewRevive returns a new Revive linter.
func NewRevive(cfg *config.ReviveSettings) *goanalysis.Linter {
	var issues []goanalysis.Issue

	analyzer := &analysis.Analyzer{
		Name: goanalysis.TheOnlyAnalyzerName,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
	}

	return goanalysis.NewLinter(
		reviveName,
		"Fast, configurable, extensible, flexible, and beautiful linter for Go. Drop-in replacement of golint.",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithContextSetter(func(lintCtx *linter.Context) {
		analyzer.Run = func(pass *analysis.Pass) (interface{}, error) {
			var files []string

			for _, file := range pass.Files {
				files = append(files, pass.Fset.PositionFor(file.Pos(), false).Filename)
			}

			conf, err := setReviveConfig(cfg)
			if err != nil {
				return nil, err
			}

			formatter, err := reviveConfig.GetFormatter("json")
			if err != nil {
				return nil, err
			}

			revive := lint.New(ioutil.ReadFile)

			lintingRules, err := reviveConfig.GetLintingRules(conf)
			if err != nil {
				return nil, err
			}

			packages, err := dots.ResolvePackages(files, []string{})
			if err != nil {
				return nil, err
			}

			failures, err := revive.Lint(packages, lintingRules, *conf)
			if err != nil {
				return nil, err
			}

			formatChan := make(chan lint.Failure)
			exitChan := make(chan bool)

			var output string
			go func() {
				output, err = formatter.Format(formatChan, *conf)
				if err != nil {
					lintCtx.Log.Errorf("Format error: %v", err)
				}
				exitChan <- true
			}()

			for f := range failures {
				if f.Confidence < conf.Confidence {
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

			for i := range results {
				issues = append(issues, goanalysis.NewIssue(&result.Issue{
					Severity: string(results[i].Severity),
					Text:     fmt.Sprintf("%q", results[i].Failure.Failure),
					Pos: token.Position{
						Filename: results[i].Position.Start.Filename,
						Line:     results[i].Position.Start.Line,
						Offset:   results[i].Position.Start.Offset,
						Column:   results[i].Position.Start.Column,
					},
					LineRange: &result.Range{
						From: results[i].Position.Start.Line,
						To:   results[i].Position.End.Line,
					},
					FromLinter: reviveName,
				}, pass))
			}

			return nil, nil
		}
	}).WithIssuesReporter(func(*linter.Context) []goanalysis.Issue {
		return issues
	}).WithLoadMode(goanalysis.LoadModeSyntax)
}

func setReviveConfig(cfg *config.ReviveSettings) (*lint.Config, error) {
	// Get revive default configuration
	conf, err := reviveConfig.GetConfig("")
	if err != nil {
		return nil, err
	}

	// Default is false
	conf.IgnoreGeneratedHeader = cfg.IgnoreGeneratedHeader

	if cfg.Severity != "" {
		conf.Severity = lint.Severity(cfg.Severity)
	}

	if cfg.Confidence != 0 {
		conf.Confidence = cfg.Confidence
	}

	// By default golangci-lint ignores missing doc comments, follow same convention by removing this default rule
	// Relevant issue: https://github.com/golangci/golangci-lint/issues/456
	delete(conf.Rules, "exported")

	if len(cfg.Rules) != 0 {
		// Clear default rules, only use rules defined in config
		conf.Rules = make(map[string]lint.RuleConfig, len(cfg.Rules))
	}
	for _, r := range cfg.Rules {
		conf.Rules[r.Name] = lint.RuleConfig{Arguments: r.Arguments, Severity: lint.Severity(r.Severity)}
	}

	conf.ErrorCode = cfg.ErrorCode
	conf.WarningCode = cfg.WarningCode

	if len(cfg.Directives) != 0 {
		// Clear default Directives, only use Directives defined in config
		conf.Directives = make(map[string]lint.DirectiveConfig, len(cfg.Directives))
	}
	for _, d := range cfg.Directives {
		conf.Directives[d.Name] = lint.DirectiveConfig{Severity: lint.Severity(d.Severity)}
	}

	return conf, nil
}
