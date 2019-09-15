package goanalysis

import (
	"context"
	"flag"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

type Linter struct {
	name, desc string
	analyzers  []*analysis.Analyzer
	cfg        map[string]map[string]interface{}
}

func NewLinter(name, desc string, analyzers []*analysis.Analyzer, cfg map[string]map[string]interface{}) *Linter {
	return &Linter{name: name, desc: desc, analyzers: analyzers, cfg: cfg}
}

func (lnt Linter) Name() string {
	return lnt.name
}

func (lnt Linter) Desc() string {
	return lnt.desc
}

func (lnt Linter) allAnalyzerNames() []string {
	var ret []string
	for _, a := range lnt.analyzers {
		ret = append(ret, a.Name)
	}
	return ret
}

func allFlagNames(fs *flag.FlagSet) []string {
	var ret []string
	fs.VisitAll(func(f *flag.Flag) {
		ret = append(ret, f.Name)
	})
	return ret
}

func valueToString(v interface{}) string {
	if ss, ok := v.([]string); ok {
		return strings.Join(ss, ",")
	}

	if is, ok := v.([]interface{}); ok {
		var ss []string
		for _, i := range is {
			ss = append(ss, fmt.Sprint(i))
		}
		return valueToString(ss)
	}

	return fmt.Sprint(v)
}

func (lnt Linter) configureAnalyzer(a *analysis.Analyzer, cfg map[string]interface{}) error {
	for k, v := range cfg {
		f := a.Flags.Lookup(k)
		if f == nil {
			validFlagNames := allFlagNames(&a.Flags)
			if len(validFlagNames) == 0 {
				return fmt.Errorf("analyzer doesn't have settings")
			}

			return fmt.Errorf("analyzer doesn't have setting %q, valid settings: %v",
				k, validFlagNames)
		}

		if err := f.Value.Set(valueToString(v)); err != nil {
			return errors.Wrapf(err, "failed to set analyzer setting %q with value %v", k, v)
		}
	}

	return nil
}

func (lnt Linter) configure() error {
	analyzersMap := map[string]*analysis.Analyzer{}
	for _, a := range lnt.analyzers {
		analyzersMap[a.Name] = a
	}

	for analyzerName, analyzerSettings := range lnt.cfg {
		a := analyzersMap[analyzerName]
		if a == nil {
			return fmt.Errorf("settings key %q must be valid analyzer name, valid analyzers: %v",
				analyzerName, lnt.allAnalyzerNames())
		}

		if err := lnt.configureAnalyzer(a, analyzerSettings); err != nil {
			return errors.Wrapf(err, "failed to configure analyzer %s", analyzerName)
		}
	}

	return nil
}

func (lnt Linter) Run(ctx context.Context, lintCtx *linter.Context) ([]result.Issue, error) {
	if err := analysis.Validate(lnt.analyzers); err != nil {
		return nil, errors.Wrap(err, "failed to validate analyzers")
	}

	if err := lnt.configure(); err != nil {
		return nil, errors.Wrap(err, "failed to configure analyzers")
	}

	runner := newRunner(lnt.name, lintCtx.Log.Child("goanalysis"), lintCtx.PkgCache, lintCtx.LoadGuard)

	diags, errs := runner.run(lnt.analyzers, lintCtx.Packages)
	for i := 1; i < len(errs); i++ {
		lintCtx.Log.Warnf("%s error: %s", lnt.Name(), errs[i])
	}
	if len(errs) != 0 {
		return nil, errs[0]
	}

	var issues []result.Issue
	for i := range diags {
		diag := &diags[i]
		issues = append(issues, result.Issue{
			FromLinter: lnt.Name(),
			Text:       fmt.Sprintf("%s: %s", diag.Analyzer.Name, diag.Message),
			Pos:        diag.Position,
		})
	}

	return issues, nil
}

func (lnt Linter) Analyzers() []*analysis.Analyzer {
	return lnt.analyzers
}

func (lnt Linter) Cfg() map[string]map[string]interface{} {
	return lnt.cfg
}

func (lnt Linter) AnalyzerToLinterNameMapping() map[*analysis.Analyzer]string {
	ret := map[*analysis.Analyzer]string{}
	for _, a := range lnt.analyzers {
		ret[a] = lnt.Name()
	}
	return ret
}
