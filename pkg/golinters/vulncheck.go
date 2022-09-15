package golinters

import (
	"fmt"
	"strings"
	"sync"

	"golang.org/x/net/context"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/vuln/client"
	"golang.org/x/vuln/vulncheck"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

const (
	vulncheckName = "vulncheck"
	vulncheckDoc  = "Package vulncheck detects uses of known vulnerabilities in Go programs."
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
		"vulncheck",
		"Package vulncheck detects uses of known vulnerabilities in Go programs.",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithContextSetter(func(lintCtx *linter.Context) {
		analyzer.Run = func(pass *analysis.Pass) (interface{}, error) {
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

func vulncheckRun(lintCtx *linter.Context, pass *analysis.Pass, settings *config.VulncheckSettings) ([]goanalysis.Issue, error) {
	dbs := []string{"https://vuln.go.dev"}
	if len(settings.VulnDatabase) > 0 {
		dbs = settings.VulnDatabase
	}
	dbClient, err := client.NewClient(dbs, client.Options{})
	if err != nil {
		return nil, err
	}

	vcfg := &vulncheck.Config{Client: dbClient, SourceGoVersion: lintCtx.Cfg.Run.Go}
	vpkgs := vulncheck.Convert(lintCtx.Packages)
	ctx := context.Background()

	r, err := vulncheck.Source(ctx, vpkgs, vcfg)
	if err != nil {
		return nil, err
	}

	imports := vulncheck.ImportChains(r)
	issues := make([]goanalysis.Issue, 0, len(r.Vulns))

	for idx, vuln := range r.Vulns {
		issues = append(issues, goanalysis.NewIssue(&result.Issue{
			Text: writeVulnerability(idx, vuln.OSV.ID, vuln.OSV.Details, writeImorts(imports[vuln])),
		}, pass))
	}

	return issues, nil
}

func writeImorts(imports []vulncheck.ImportChain) string {
	var s strings.Builder
	for _, i := range imports {
		indent := 0
		for _, pkg := range i {
			s.WriteString(fmt.Sprintf("%s|_ %s", strings.Repeat(" ", indent), pkg.Name))
		}
	}

	return s.String()
}

func writeVulnerability(idx int, id, details, imports string) string {
	return fmt.Sprintf(`Vulnerability #%d: %s
%s
%s
  More info: https://pkg.go.dev/vuln/%s
`, idx, id, details, imports, id)
}
