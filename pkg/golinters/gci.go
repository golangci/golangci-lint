package golinters

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	gcicfg "github.com/daixiang0/gci/pkg/config"
	"github.com/daixiang0/gci/pkg/gci"
	"github.com/daixiang0/gci/pkg/io"
	"github.com/daixiang0/gci/pkg/log"
	"github.com/daixiang0/gci/pkg/section"
	"github.com/golangci/modinfo"
	"github.com/hexops/gotextdiff"
	"github.com/hexops/gotextdiff/myers"
	"github.com/hexops/gotextdiff/span"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
	"github.com/golangci/golangci-lint/pkg/golinters/internal"
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
		Requires: []*analysis.Analyzer{
			modinfo.Analyzer,
		},
	}

	var cfg *gcicfg.Config
	if settings != nil {
		rawCfg := gcicfg.YamlConfig{
			Cfg: gcicfg.BoolConfig{
				SkipGenerated: settings.SkipGenerated,
				CustomOrder:   settings.CustomOrder,
			},
			SectionStrings: settings.Sections,
		}

		if settings.LocalPrefixes != "" {
			prefix := []string{"standard", "default", fmt.Sprintf("prefix(%s)", settings.LocalPrefixes)}
			rawCfg.SectionStrings = prefix
		}

		var err error
		cfg, err = YamlConfig{origin: rawCfg}.Parse()
		if err != nil {
			internal.LinterLogger.Fatalf("gci: configuration parsing: %v", err)
		}
	}

	var lock sync.Mutex

	return goanalysis.NewLinter(
		gciName,
		"Gci controls Go package import order and makes it always deterministic.",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithContextSetter(func(lintCtx *linter.Context) {
		analyzer.Run = func(pass *analysis.Pass) (any, error) {
			var err error
			cfg.Sections, err = hackSectionList(pass, cfg)
			if err != nil {
				return nil, err
			}

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
	fileNames := internal.GetFileNames(pass)

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

		is, err := internal.ExtractIssuesFromPatch(diff, lintCtx, gciName, getIssuedTextGci)
		if err != nil {
			return nil, fmt.Errorf("can't extract issues from gci diff output %s: %w", diff, err)
		}

		for i := range is {
			issues = append(issues, goanalysis.NewIssue(&is[i], pass))
		}
	}

	return issues, nil
}

func getIssuedTextGci(settings *config.LintersSettings) string {
	text := "File is not `gci`-ed"

	hasOptions := settings.Gci.SkipGenerated || len(settings.Gci.Sections) > 0
	if !hasOptions {
		return text
	}

	text += " with"

	if settings.Gci.SkipGenerated {
		text += " --skip-generated"
	}

	if len(settings.Gci.Sections) > 0 {
		for _, sect := range settings.Gci.Sections {
			text += " -s " + sect
		}
	}

	if settings.Gci.CustomOrder {
		text += " --custom-order"
	}

	return text
}

func hackSectionList(pass *analysis.Pass, cfg *gcicfg.Config) (section.SectionList, error) {
	var sections section.SectionList

	for _, sect := range cfg.Sections {
		// local module hack
		if v, ok := sect.(*section.LocalModule); ok {
			info, err := modinfo.FindModuleFromPass(pass)
			if err != nil {
				return nil, err
			}

			if info.Path == "" {
				continue
			}

			v.Path = info.Path
		}

		sections = append(sections, sect)
	}

	return sections, nil
}

// diffFormattedFilesToArray is a copy of gci.DiffFormattedFilesToArray without io.StdInGenerator.
// gci.DiffFormattedFilesToArray uses gci.processStdInAndGoFilesInPaths that uses io.StdInGenerator but stdin is not active on CI.
// https://github.com/daixiang0/gci/blob/6f5cb16718ba07f0342a58de9b830ec5a6d58790/pkg/gci/gci.go#L63-L75
// https://github.com/daixiang0/gci/blob/6f5cb16718ba07f0342a58de9b830ec5a6d58790/pkg/gci/gci.go#L80
func diffFormattedFilesToArray(paths []string, cfg gcicfg.Config, diffs *[]string, lock *sync.Mutex) error {
	log.InitLogger()
	defer func() { _ = log.L().Sync() }()

	return gci.ProcessFiles(io.GoFilesInPathsGenerator(paths, true), cfg, func(filePath string, unmodifiedFile, formattedFile []byte) error {
		fileURI := span.URIFromPath(filePath)
		edits := myers.ComputeEdits(fileURI, string(unmodifiedFile), string(formattedFile))
		unifiedEdits := gotextdiff.ToUnified(filePath, filePath, string(unmodifiedFile), edits)
		lock.Lock()
		*diffs = append(*diffs, fmt.Sprint(unifiedEdits))
		lock.Unlock()
		return nil
	})
}

// Code bellow this comment is borrowed and modified from gci.
// https://github.com/daixiang0/gci/blob/4725b0c101801e7449530eee2ddb0c72592e3405/pkg/config/config.go

var defaultOrder = map[string]int{
	section.StandardType:    0,
	section.DefaultType:     1,
	section.CustomType:      2,
	section.BlankType:       3,
	section.DotType:         4,
	section.AliasType:       5,
	section.LocalModuleType: 6,
}

type YamlConfig struct {
	origin gcicfg.YamlConfig
}

//nolint:gocritic // code borrowed from gci and modified to fix LocalModule section behavior.
func (g YamlConfig) Parse() (*gcicfg.Config, error) {
	var err error

	sections, err := section.Parse(g.origin.SectionStrings)
	if err != nil {
		return nil, err
	}

	if sections == nil {
		sections = section.DefaultSections()
	}

	// if default order sorted sections
	if !g.origin.Cfg.CustomOrder {
		sort.Slice(sections, func(i, j int) bool {
			sectionI, sectionJ := sections[i].Type(), sections[j].Type()

			if strings.Compare(sectionI, sectionJ) == 0 {
				return strings.Compare(sections[i].String(), sections[j].String()) < 0
			}
			return defaultOrder[sectionI] < defaultOrder[sectionJ]
		})
	}

	sectionSeparators, err := section.Parse(g.origin.SectionSeparatorStrings)
	if err != nil {
		return nil, err
	}
	if sectionSeparators == nil {
		sectionSeparators = section.DefaultSectionSeparators()
	}

	return &gcicfg.Config{BoolConfig: g.origin.Cfg, Sections: sections, SectionSeparators: sectionSeparators}, nil
}
