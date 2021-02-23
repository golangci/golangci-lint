package golinters

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/token"
	"io/ioutil"

	"github.com/BurntSushi/toml"
	"github.com/mgechev/dots"
	reviveConfig "github.com/mgechev/revive/config"
	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
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

			conf, err := getReviveConfig(cfg)
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
					Text:     fmt.Sprintf("%s: %s", results[i].RuleName, results[i].Failure.Failure),
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

// This function mimics the GetConfig function of revive.
// This allow to get default values and right types.
// https://github.com/golangci/golangci-lint/issues/1745
// https://github.com/mgechev/revive/blob/389ba853b0b3587f0c3b71b5f0c61ea4e23928ec/config/config.go#L155
func getReviveConfig(cfg *config.ReviveSettings) (*lint.Config, error) {
	rawRoot := createConfigMap(cfg)

	buf := bytes.NewBuffer(nil)

	err := toml.NewEncoder(buf).Encode(rawRoot)
	if err != nil {
		return nil, err
	}

	conf := defaultConfig()

	_, err = toml.DecodeReader(buf, conf)
	if err != nil {
		return nil, err
	}

	normalizeConfig(conf)

	// By default golangci-lint ignores missing doc comments, follow same convention by removing this default rule
	// Relevant issue: https://github.com/golangci/golangci-lint/issues/456
	delete(conf.Rules, "package-comments")
	delete(conf.Rules, "exported")

	return conf, nil
}

func createConfigMap(cfg *config.ReviveSettings) map[string]interface{} {
	rawRoot := map[string]interface{}{
		"ignoreGeneratedHeader": cfg.IgnoreGeneratedHeader,
		"confidence":            cfg.Confidence,
		"severity":              cfg.Severity,
		"errorCode":             cfg.ErrorCode,
		"warningCode":           cfg.WarningCode,
	}

	rawDirectives := map[string]map[string]interface{}{}
	for _, directive := range cfg.Directives {
		rawDirectives[directive.Name] = map[string]interface{}{
			"severity": directive.Severity,
		}
	}

	if len(rawDirectives) > 0 {
		rawRoot["directive"] = rawDirectives
	}

	rawRules := map[string]map[string]interface{}{}
	for _, s := range cfg.Rules {
		rawRules[s.Name] = map[string]interface{}{
			"severity":  s.Severity,
			"arguments": s.Arguments,
		}
	}

	if len(rawRules) > 0 {
		rawRoot["rule"] = rawRules
	}

	return rawRoot
}

// This element is not exported by revive, so we need copy the code.
// Extracted from https://github.com/mgechev/revive/blob/389ba853b0b3587f0c3b71b5f0c61ea4e23928ec/config/config.go#L15
var defaultRules = []lint.Rule{
	&rule.VarDeclarationsRule{},
	&rule.PackageCommentsRule{},
	&rule.DotImportsRule{},
	&rule.BlankImportsRule{},
	&rule.ExportedRule{},
	&rule.VarNamingRule{},
	&rule.IndentErrorFlowRule{},
	&rule.IfReturnRule{},
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
}

// This element is not exported by revive, so we need copy the code.
// Extracted from https://github.com/mgechev/revive/blob/389ba853b0b3587f0c3b71b5f0c61ea4e23928ec/config/config.go#L133
func normalizeConfig(cfg *lint.Config) {
	if cfg.Confidence == 0 {
		cfg.Confidence = 0.8
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
// Extracted from https://github.com/mgechev/revive/blob/389ba853b0b3587f0c3b71b5f0c61ea4e23928ec/config/config.go#L182
func defaultConfig() *lint.Config {
	defaultConfig := lint.Config{
		Confidence: 0.0,
		Severity:   lint.SeverityWarning,
		Rules:      map[string]lint.RuleConfig{},
	}
	for _, r := range defaultRules {
		defaultConfig.Rules[r.Name()] = lint.RuleConfig{}
	}
	return &defaultConfig
}
