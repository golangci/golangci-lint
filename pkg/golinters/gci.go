package golinters

import (
	"fmt"
	"strings"
	"sync"

	gcicfg "github.com/daixiang0/gci/pkg/config"
	"github.com/daixiang0/gci/pkg/gci"
	"github.com/daixiang0/gci/pkg/io"
	"github.com/daixiang0/gci/pkg/log"
	"github.com/hexops/gotextdiff"
	"github.com/hexops/gotextdiff/myers"
	"github.com/hexops/gotextdiff/span"
	"github.com/pkg/errors"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
)

const gciName = "gci"

func NewGci(settings *config.GciSettings) *goanalysis.Linter {
	var mu sync.Mutex
	var resIssues []goanalysis.Issue

	analyzer := &analysis.Analyzer{
		Name: gciName,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
		Run:  goanalysis.DummyRun,
	}

	var cfg *gcicfg.Config
	if settings != nil {
		rawCfg := gcicfg.YamlConfig{
			Cfg: gcicfg.BoolConfig{
				SkipGenerated: settings.SkipGenerated,
			},
			SectionStrings: settings.Sections,
		}

		if settings.LocalPrefixes != "" {
			prefix := []string{"standard", "default", fmt.Sprintf("prefix(%s)", settings.LocalPrefixes)}
			rawCfg.SectionStrings = prefix
		}

		cfg, _ = rawCfg.Parse()
	}

	var lock sync.Mutex

	return goanalysis.NewLinter(
		gciName,
		"Gci controls golang package import order and makes it always deterministic.",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithContextSetter(func(lintCtx *linter.Context) {
		analyzer.Run = func(pass *analysis.Pass) (interface{}, error) {
			issues, err := runGci(pass, lintCtx, cfg, &lock)
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

func runGci(pass *analysis.Pass, lintCtx *linter.Context, cfg *gcicfg.Config, lock *sync.Mutex) ([]goanalysis.Issue, error) {
	var fileNames []string
	for _, f := range pass.Files {
		pos := pass.Fset.PositionFor(f.Pos(), false)
		fileNames = append(fileNames, pos.Filename)
	}

	var diffs []string
	err := diffFormattedFilesToArray(fileNames, *cfg, &diffs, lock)
	if err != nil {
		return nil, err
	}

	var issues []goanalysis.Issue

	for _, diff := range diffs {
		if diff == "" {
			continue
		}

		is, err := extractIssuesFromPatch(diff, lintCtx, gciName)
		if err != nil {
			return nil, errors.Wrapf(err, "can't extract issues from gci diff output %s", diff)
		}

		for i := range is {
			issues = append(issues, goanalysis.NewIssue(&is[i], pass))
		}
	}

	return issues, nil
}

// diffFormattedFilesToArray is a copy of gci.DiffFormattedFilesToArray without io.StdInGenerator.
// gci.DiffFormattedFilesToArray uses gci.processStdInAndGoFilesInPaths that uses io.StdInGenerator but stdin is not active on CI.
// https://github.com/daixiang0/gci/blob/6f5cb16718ba07f0342a58de9b830ec5a6d58790/pkg/gci/gci.go#L63-L75
// https://github.com/daixiang0/gci/blob/6f5cb16718ba07f0342a58de9b830ec5a6d58790/pkg/gci/gci.go#L80
func diffFormattedFilesToArray(paths []string, cfg gcicfg.Config, diffs *[]string, lock *sync.Mutex) error {
	log.InitLogger()
	defer func() { _ = log.L().Sync() }()

	return gci.ProcessFiles(io.GoFilesInPathsGenerator(paths), cfg, func(filePath string, unmodifiedFile, formattedFile []byte) error {
		fileURI := span.URIFromPath(filePath)
		edits := myers.ComputeEdits(fileURI, string(unmodifiedFile), string(formattedFile))
		unifiedEdits := gotextdiff.ToUnified(filePath, filePath, string(unmodifiedFile), edits)
		lock.Lock()
		*diffs = append(*diffs, fmt.Sprint(unifiedEdits))
		lock.Unlock()
		return nil
	})
}

func getErrorTextForGci(settings config.GciSettings) string {
	text := "File is not `gci`-ed"

	hasOptions := settings.SkipGenerated || len(settings.Sections) > 0
	if !hasOptions {
		return text
	}

	text += " with"

	if settings.SkipGenerated {
		text += " -skip-generated"
	}

	if len(settings.Sections) > 0 {
		text += " -s " + strings.Join(settings.Sections, ",")
	}

	return text
}
