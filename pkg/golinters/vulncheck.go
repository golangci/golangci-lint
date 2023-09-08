package golinters

import (
	"bytes"
	"path/filepath"
	"sync"

	"golang.org/x/net/context"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/vuln/scan"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

const (
	vulncheckName = "vulncheck"
	vulncheckDoc  = "vulncheck detects uses of known vulnerabilities in Go programs."
)

func NewVulncheck(settings *config.VulncheckSettings) *goanalysis.Linter {
	var mu sync.Mutex
	var resIssues []goanalysis.Issue

	analyzer := &analysis.Analyzer{
		Name: vulncheckName,
		Doc:  vulncheckDoc,
		Run:  goanalysis.DummyRun,
	}

	return goanalysis.NewLinter(
		vulncheckName,
		vulncheckDoc,
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithContextSetter(func(lintCtx *linter.Context) {
		analyzer.Run = func(pass *analysis.Pass) (any, error) {
			issues, err := vulncheckRun(lintCtx, pass, settings)
			if err != nil {
				return nil, err
			}

			mu.Lock()
			resIssues = append(resIssues, issues...)
			mu.Unlock()

			return nil, nil
		}
	}).WithIssuesReporter(func(*linter.Context) []goanalysis.Issue {
		return resIssues
	})
}

func vulncheckRun(lintCtx *linter.Context, pass *analysis.Pass, _ *config.VulncheckSettings) ([]goanalysis.Issue, error) {
	files := getFileNames(pass)

	ctx := context.Background()
	lintCtx.Log.Errorf("%v\n", files)

	issues := []goanalysis.Issue{}
	for _, file := range files {
		lintCtx.Log.Errorf("%s %s %s %s\n", "-json", "-C", filepath.Dir(file), ".")
		cmd := scan.Command(ctx, "-json", "-C", filepath.Dir(file), ".")
		buf := &bytes.Buffer{}
		cmd.Stderr = buf
		cmd.Stdout = buf
		err := cmd.Start()
		if err != nil {
			return issues, err
		}
		err = cmd.Wait()
		if err != nil {
			return issues, err
		}
		issues = append(issues, goanalysis.NewIssue(&result.Issue{
			Text:       buf.String(),
			FromLinter: vulncheckName},
			pass))
	}

	lintCtx.Log.Errorf("%v\n", issues)
	return issues, nil
}
