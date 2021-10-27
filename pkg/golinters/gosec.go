package golinters

import (
	"fmt"
	"go/token"
	"io/ioutil"
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

	gasConfig := gosec.NewConfig()

	var filters []rules.RuleFilter
	if settings != nil {
		filters = gosecRuleFilters(settings.Includes, settings.Excludes)

		for k, v := range settings.Config {
			// Uses ToUpper because the parsing of the map's key change the key to lowercase.
			// The value is not impacted by that: the case is respected.
			gasConfig.Set(strings.ToUpper(k), v)
		}
	}

	ruleDefinitions := rules.Generate(filters...)

	logger := log.New(ioutil.Discard, "", 0)

	analyzer := &analysis.Analyzer{
		Name: gosecName,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
	}
	return goanalysis.NewLinter(
		gosecName,
		"Inspects source code for security problems",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithContextSetter(func(lintCtx *linter.Context) {
		analyzer.Run = func(pass *analysis.Pass) (interface{}, error) {
			gosecAnalyzer := gosec.NewAnalyzer(gasConfig, true, settings.ExcludeGenerated, logger)
			gosecAnalyzer.LoadRules(ruleDefinitions.Builders())

			pkg := &packages.Package{
				Fset:      pass.Fset,
				Syntax:    pass.Files,
				Types:     pass.Pkg,
				TypesInfo: pass.TypesInfo,
			}
			gosecAnalyzer.Check(pkg)
			issues, _, _ := gosecAnalyzer.Report()
			if len(issues) == 0 {
				return nil, nil
			}

			res := make([]goanalysis.Issue, 0, len(issues))
			for _, i := range issues {
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

				res = append(res, goanalysis.NewIssue(&result.Issue{
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

			mu.Lock()
			resIssues = append(resIssues, res...)
			mu.Unlock()

			return nil, nil
		}
	}).WithIssuesReporter(func(*linter.Context) []goanalysis.Issue {
		return resIssues
	}).WithLoadMode(goanalysis.LoadModeTypesInfo)
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
