package golinters

import (
	"fmt"
	"strings"
	"sync"

	"github.com/OpenPeeDeeP/depguard"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/loader" //nolint:staticcheck // require changes in github.com/OpenPeeDeeP/depguard

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

func parseDepguardSettings(dgSettings *config.DepGuardSettings) (map[*depguard.Depguard]map[string]string, error) {
	parsedDgSettings := make(map[*depguard.Depguard]map[string]string)
	dg := &depguard.Depguard{
		Packages:        dgSettings.Packages,
		IncludeGoRoot:   dgSettings.IncludeGoRoot,
		IgnoreFileRules: dgSettings.IgnoreFileRules,
	}

	if err := setDepguardListType(dg, dgSettings.ListType); err != nil {
		return nil, err
	}
	setupDepguardPackages(dg, dgSettings.PackagesWithErrorMessage)
	if dgSettings.PackagesWithErrorMessage != nil {
		parsedDgSettings[dg] = dgSettings.PackagesWithErrorMessage
	} else {
		parsedDgSettings[dg] = make(map[string]string)
	}

	for _, additionalGuard := range dgSettings.AdditionalGuards {
		additionalDg := &depguard.Depguard{
			Packages:        additionalGuard.Packages,
			IncludeGoRoot:   additionalGuard.IncludeGoRoot,
			IgnoreFileRules: additionalGuard.IgnoreFileRules,
		}

		if err := setDepguardListType(additionalDg, additionalGuard.ListType); err != nil {
			return nil, err
		}
		setupDepguardPackages(additionalDg, additionalGuard.PackagesWithErrorMessage)
		if additionalGuard.PackagesWithErrorMessage != nil {
			parsedDgSettings[additionalDg] = additionalGuard.PackagesWithErrorMessage
		} else {
			parsedDgSettings[additionalDg] = make(map[string]string)
		}
	}

	return parsedDgSettings, nil
}

func postProcessIssue(
	issue *depguard.Issue,
	dg *depguard.Depguard,
	packagesWithErrorMessage map[string]string,
	lintCtx *linter.Context,
) *result.Issue {
	msgSuffix := "is in the blacklist"
	if dg.ListType == depguard.LTWhitelist {
		msgSuffix = "is not in the whitelist"
	}

	userSuppliedMsgSuffix := packagesWithErrorMessage[issue.PackageName]
	if userSuppliedMsgSuffix != "" {
		userSuppliedMsgSuffix = ": " + userSuppliedMsgSuffix
	}

	return &result.Issue{
		Pos:        issue.Position,
		Text:       fmt.Sprintf("%s %s%s", formatCode(issue.PackageName, lintCtx.Cfg), msgSuffix, userSuppliedMsgSuffix),
		FromLinter: linterName,
	}
}

func setDepguardListType(dg *depguard.Depguard, listType string) error {
	var found bool
	dg.ListType, found = depguard.StringToListType[strings.ToLower(listType)]
	if !found {
		if listType != "" {
			return fmt.Errorf("unsure what list type %s is", listType)
		}
		dg.ListType = depguard.LTBlacklist
	}

	return nil
}

func setupDepguardPackages(
	dg *depguard.Depguard,
	packagesWithErrorMessage map[string]string,
) {
	if dg.ListType == depguard.LTBlacklist {
		// if the list type was a blacklist the packages with error messages should
		// be included in the blacklist package list

		noMessagePackages := make(map[string]bool)
		for _, pkg := range dg.Packages {
			noMessagePackages[pkg] = true
		}

		for pkg := range packagesWithErrorMessage {
			if _, ok := noMessagePackages[pkg]; !ok {
				dg.Packages = append(dg.Packages, pkg)
			}
		}
	}
}

const linterName = "depguard"

func NewDepguard() *goanalysis.Linter {
	var mu sync.Mutex
	var resIssues []goanalysis.Issue

	analyzer := &analysis.Analyzer{
		Name: linterName,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
	}
	return goanalysis.NewLinter(
		linterName,
		"Go linter that checks if package imports are in a list of acceptable packages",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithContextSetter(func(lintCtx *linter.Context) {
		dgSettings := &lintCtx.Settings().Depguard
		analyzer.Run = func(pass *analysis.Pass) (interface{}, error) {
			parsedDgSettings, err := parseDepguardSettings(dgSettings)
			if err != nil {
				return nil, err
			}

			loadConfig := &loader.Config{
				Cwd:   "",  // fallbacked to os.Getcwd
				Build: nil, // fallbacked to build.Default
			}
			prog := goanalysis.MakeFakeLoaderProgram(pass)

			for dg, packagesWithErrorMessage := range parsedDgSettings {
				issues, err := dg.Run(loadConfig, prog)
				if err != nil {
					return nil, err
				}
				res := make([]goanalysis.Issue, 0, len(issues))
				for _, i := range issues {
					lintIssue := postProcessIssue(i, dg, packagesWithErrorMessage, lintCtx)
					res = append(res, goanalysis.NewIssue(lintIssue, pass))
				}
				mu.Lock()
				resIssues = append(resIssues, res...)
				mu.Unlock()
			}
			return nil, nil
		}
	}).WithIssuesReporter(func(*linter.Context) []goanalysis.Issue {
		return resIssues
	}).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
