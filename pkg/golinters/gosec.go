package golinters

import (
	"fmt"
	"go/token"
	"io"
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/securego/gosec/v2"
	"github.com/securego/gosec/v2/rules"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/packages"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

const gosecName = "gosec"

func NewGosec(settings *config.GoSecSettings) *goanalysis.Linter {
	var mu sync.Mutex
	var resIssues []goanalysis.Issue

	conf := gosec.NewConfig()

	var filters []rules.RuleFilter
	if settings != nil {
		filters = gosecRuleFilters(settings.Includes, settings.Excludes)

		for k, v := range settings.Config {
			if k != gosec.Globals {
				// Uses ToUpper because the parsing of the map's key change the key to lowercase.
				// The value is not impacted by that: the case is respected.
				k = strings.ToUpper(k)
			}
			conf.Set(k, v)
		}
	}

	logger := log.New(io.Discard, "", 0)

	ruleDefinitions := rules.Generate(false, filters...)

	analyzer := &analysis.Analyzer{
		Name: gosecName,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
		Run:  goanalysis.DummyRun,
	}

	return goanalysis.NewLinter(
		gosecName,
		"Inspects source code for security problems",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithContextSetter(func(lintCtx *linter.Context) {
		analyzer.Run = func(pass *analysis.Pass) (any, error) {
			// The `gosecAnalyzer` is here because of concurrency issue.
			gosecAnalyzer := gosec.NewAnalyzer(conf, true, settings.ExcludeGenerated, false, settings.Concurrency, logger)
			gosecAnalyzer.LoadRules(ruleDefinitions.RulesInfo())

			issues := runGoSec(lintCtx, pass, settings, gosecAnalyzer)

			mu.Lock()
			resIssues = append(resIssues, issues...)
			mu.Unlock()

			return nil, nil
		}
	}).WithIssuesReporter(func(*linter.Context) []goanalysis.Issue {
		return resIssues
	}).WithLoadMode(goanalysis.LoadModeTypesInfo)
}

func runGoSec(lintCtx *linter.Context, pass *analysis.Pass, settings *config.GoSecSettings, analyzer *gosec.Analyzer) []goanalysis.Issue {
	pkg := &packages.Package{
		Fset:      pass.Fset,
		Syntax:    pass.Files,
		Types:     pass.Pkg,
		TypesInfo: pass.TypesInfo,
	}

	analyzer.Check(pkg)

	secIssues, _, _ := analyzer.Report()
	if len(secIssues) == 0 {
		return nil
	}

	severity, err := convertToScore(settings.Severity)
	if err != nil {
		lintCtx.Log.Warnf("The provided severity %v", err)
	}

	confidence, err := convertToScore(settings.Confidence)
	if err != nil {
		lintCtx.Log.Warnf("The provided confidence %v", err)
	}

	secIssues = filterIssues(secIssues, severity, confidence)

	issues := make([]goanalysis.Issue, 0, len(secIssues))
	for _, i := range secIssues {
		text := fmt.Sprintf("%s: %s", i.RuleID, i.What) // TODO: use severity and confidence

		var r *result.Range

		line, err := strconv.Atoi(i.Line)
		if err != nil {
			r = &result.Range{}
			if n, rerr := fmt.Sscanf(i.Line, "%d-%d", &r.From, &r.To); rerr != nil || n != 2 {
				lintCtx.Log.Warnf("Can't convert gosec line number %q of %v to int: %s", i.Line, i, err)
				continue
			}
			line = r.From
		}

		column, err := strconv.Atoi(i.Col)
		if err != nil {
			lintCtx.Log.Warnf("Can't convert gosec column number %q of %v to int: %s", i.Col, i, err)
			continue
		}

		issues = append(issues, goanalysis.NewIssue(&result.Issue{
			Pos: token.Position{
				Filename: i.File,
				Line:     line,
				Column:   column,
			},
			Text:       text,
			LineRange:  r,
			FromLinter: gosecName,
		}, pass))
	}

	return issues
}

// based on https://github.com/securego/gosec/blob/569328eade2ccbad4ce2d0f21ee158ab5356a5cf/cmd/gosec/main.go#L170-L188
func gosecRuleFilters(includes, excludes []string) []rules.RuleFilter {
	var filters []rules.RuleFilter

	if len(includes) > 0 {
		filters = append(filters, rules.NewRuleFilter(false, includes...))
	}

	if len(excludes) > 0 {
		filters = append(filters, rules.NewRuleFilter(true, excludes...))
	}

	return filters
}

// code borrowed from https://github.com/securego/gosec/blob/69213955dacfd560562e780f723486ef1ca6d486/cmd/gosec/main.go#L250-L262
func convertToScore(str string) (gosec.Score, error) {
	str = strings.ToLower(str)
	switch str {
	case "", "low":
		return gosec.Low, nil
	case "medium":
		return gosec.Medium, nil
	case "high":
		return gosec.High, nil
	default:
		return gosec.Low, fmt.Errorf("'%s' is invalid, use low instead. Valid options: low, medium, high", str)
	}
}

// code borrowed from https://github.com/securego/gosec/blob/69213955dacfd560562e780f723486ef1ca6d486/cmd/gosec/main.go#L264-L276
func filterIssues(issues []*gosec.Issue, severity, confidence gosec.Score) []*gosec.Issue {
	res := make([]*gosec.Issue, 0)
	for _, issue := range issues {
		if issue.Severity >= severity && issue.Confidence >= confidence {
			res = append(res, issue)
		}
	}
	return res
}
