package internal

import (
	"go/token"
	"strings"

	"github.com/daixiang0/gci/pkg/analyzer"
	"github.com/daixiang0/gci/pkg/log"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"

	"github.com/daixiang0/gci/pkg/config"
	"github.com/daixiang0/gci/pkg/gci"
	"github.com/daixiang0/gci/pkg/io"
)

const (
	NoInlineCommentsFlag  = "noInlineComments"
	NoPrefixCommentsFlag  = "noPrefixComments"
	SkipGeneratedFlag     = "skipGenerated"
	SectionsFlag          = "Sections"
	SectionSeparatorsFlag = "SectionSeparators"
	NoLexOrderFlag        = "NoLexOrder"
	CustomOrderFlag       = "CustomOrder"
	PrefixDelimiterFlag   = "PrefixDelimiter"
)

const SectionDelimiter = ","

var (
	noInlineComments     bool
	noPrefixComments     bool
	skipGenerated        bool
	sectionsStr          string
	sectionSeparatorsStr string
	noLexOrder           bool
	customOrder          bool
	prefixDelimiter      string
)

func NewAnalyzer() *analysis.Analyzer {
	log.InitLogger()
	_ = log.L().Sync()

	a := &analysis.Analyzer{
		Name:     "gci",
		Doc:      "Checks if code and import statements are formatted, it makes import statements always deterministic.",
		Run:      runAnalysis,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}

	a.Flags.BoolVar(&noInlineComments, NoInlineCommentsFlag, false,
		"If comments in the same line as the input should be present")
	a.Flags.BoolVar(&noPrefixComments, NoPrefixCommentsFlag, false,
		"If comments above an input should be present")
	a.Flags.BoolVar(&skipGenerated, SkipGeneratedFlag, false,
		"Skip generated files")
	a.Flags.StringVar(&sectionsStr, SectionsFlag, "",
		"Specify the Sections format that should be used to check the file formatting")
	a.Flags.StringVar(&sectionSeparatorsStr, SectionSeparatorsFlag, "",
		"Specify the Sections that are inserted as Separators between Sections")
	a.Flags.BoolVar(&noLexOrder, NoLexOrderFlag, false,
		"Drops lexical ordering for custom sections")
	a.Flags.BoolVar(&customOrder, CustomOrderFlag, false,
		"Enable custom order of sections")

	a.Flags.StringVar(&prefixDelimiter, PrefixDelimiterFlag, SectionDelimiter, "")

	return a
}

func runAnalysis(pass *analysis.Pass) (any, error) {
	var fileReferences []*token.File
	// extract file references for all files in the analyzer pass
	for _, pkgFile := range pass.Files {
		fileForPos := pass.Fset.File(pkgFile.Package)
		if fileForPos != nil {
			fileReferences = append(fileReferences, fileForPos)
		}
	}

	expectedNumFiles := len(pass.Files)
	foundNumFiles := len(fileReferences)
	if expectedNumFiles != foundNumFiles {
		return nil, InvalidNumberOfFilesInAnalysis{expectedNumFiles, foundNumFiles}
	}

	gciCfg, err := generateGciConfiguration(pass.Module.Path).Parse()
	if err != nil {
		return nil, err
	}

	for _, file := range fileReferences {
		unmodifiedFile, formattedFile, err := gci.LoadFormatGoFile(io.File{FilePath: file.Name()}, *gciCfg)
		if err != nil {
			return nil, err
		}

		fix, err := analyzer.GetSuggestedFix(file, unmodifiedFile, formattedFile)
		if err != nil {
			return nil, err
		}

		if fix == nil {
			// no difference
			continue
		}

		pass.Report(analysis.Diagnostic{
			Pos:            fix.TextEdits[0].Pos,
			Message:        "File is not properly formatted",
			SuggestedFixes: []analysis.SuggestedFix{*fix},
		})
	}

	return nil, nil
}

func generateGciConfiguration(modPath string) *config.YamlConfig {
	fmtCfg := config.BoolConfig{
		NoInlineComments: noInlineComments,
		NoPrefixComments: noPrefixComments,
		Debug:            false,
		SkipGenerated:    skipGenerated,
		CustomOrder:      customOrder,
		NoLexOrder:       noLexOrder,
	}

	var sectionStrings []string
	if sectionsStr != "" {
		s := strings.Split(sectionsStr, SectionDelimiter)
		for _, a := range s {
			sectionStrings = append(sectionStrings, strings.ReplaceAll(a, prefixDelimiter, SectionDelimiter))
		}
	}

	var sectionSeparatorStrings []string
	if sectionSeparatorsStr != "" {
		sectionSeparatorStrings = strings.Split(sectionSeparatorsStr, SectionDelimiter)
	}

	return &config.YamlConfig{Cfg: fmtCfg, SectionStrings: sectionStrings, SectionSeparatorStrings: sectionSeparatorStrings, ModPath: modPath}
}
