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

const depguardLinterName = "depguard"

func NewDepguard() *goanalysis.Linter {
	var mu sync.Mutex
	var resIssues []goanalysis.Issue

	analyzer := &analysis.Analyzer{
		Name: depguardLinterName,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
	}
	return goanalysis.NewLinter(
		depguardLinterName,
		"Go linter that checks if package imports are in a list of acceptable packages",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithContextSetter(func(lintCtx *linter.Context) {
		parsedSettings, err := parseDepGuardSettings(&lintCtx.Settings().Depguard)

		analyzer.Run = func(pass *analysis.Pass) (interface{}, error) {
			if err != nil {
				return nil, err
			}

			loadConfig := &loader.Config{
				Cwd:   "",  // fallbacked to os.Getcwd
				Build: nil, // fallbacked to build.Default
			}

			prog := goanalysis.MakeFakeLoaderProgram(pass)

			for dg, pkgsWithErrorMessage := range parsedSettings {
				issues, errRun := runDepGuard(dg, pkgsWithErrorMessage, loadConfig, prog, pass)
				if errRun != nil {
					return nil, errRun
				}

				mu.Lock()
				resIssues = append(resIssues, issues...)
				mu.Unlock()
			}

			return nil, nil
		}
	}).WithIssuesReporter(func(*linter.Context) []goanalysis.Issue {
		return resIssues
	}).WithLoadMode(goanalysis.LoadModeTypesInfo)
}

func parseDepGuardSettings(settings *config.DepGuardSettings) (map[*depguard.Depguard]map[string]string, error) {
	parsedSettings := make(map[*depguard.Depguard]map[string]string)

	err := parseDGSettings(settings, parsedSettings)
	if err != nil {
		return nil, err
	}

	for _, additional := range settings.AdditionalGuards {
		add := additional
		err := parseDGSettings(&add, parsedSettings)
		if err != nil {
			return nil, err
		}
	}

	return parsedSettings, nil
}

func parseDGSettings(settings *config.DepGuardSettings, parsedSettings map[*depguard.Depguard]map[string]string) error {
	dg := &depguard.Depguard{
		Packages:        settings.Packages,
		IncludeGoRoot:   settings.IncludeGoRoot,
		IgnoreFileRules: settings.IgnoreFileRules,
	}

	var err error
	dg.ListType, err = getDepGuardListType(settings.ListType)
	if err != nil {
		return err
	}

	// if the list type was a blacklist the packages with error messages should  be included in the blacklist package list
	if dg.ListType == depguard.LTBlacklist {
		noMessagePackages := make(map[string]bool)
		for _, pkg := range dg.Packages {
			noMessagePackages[pkg] = true
		}

		for pkg := range settings.PackagesWithErrorMessage {
			if _, ok := noMessagePackages[pkg]; !ok {
				dg.Packages = append(dg.Packages, pkg)
			}
		}
	}

	if settings.PackagesWithErrorMessage != nil {
		parsedSettings[dg] = settings.PackagesWithErrorMessage
	} else {
		parsedSettings[dg] = make(map[string]string)
	}

	return nil
}

func getDepGuardListType(listType string) (depguard.ListType, error) {
	if listType == "" {
		return depguard.LTBlacklist, nil
	}

	listT, found := depguard.StringToListType[strings.ToLower(listType)]
	if !found {
		return depguard.LTBlacklist, fmt.Errorf("unsure what list type %s is", listType)
	}

	return listT, nil
}

func runDepGuard(dg *depguard.Depguard, pkgsWithErrorMessage map[string]string,
	loadConfig *loader.Config, prog *loader.Program, pass *analysis.Pass) ([]goanalysis.Issue, error) {
	issues, err := dg.Run(loadConfig, prog)
	if err != nil {
		return nil, err
	}

	res := make([]goanalysis.Issue, 0, len(issues))

	for _, issue := range issues {
		msgSuffix := "is in the blacklist"
		if dg.ListType == depguard.LTWhitelist {
			msgSuffix = "is not in the whitelist"
		}

		userSuppliedMsgSuffix := pkgsWithErrorMessage[issue.PackageName]
		if userSuppliedMsgSuffix != "" {
			userSuppliedMsgSuffix = ": " + userSuppliedMsgSuffix
		}

		res = append(res,
			goanalysis.NewIssue(&result.Issue{
				Pos:        issue.Position,
				Text:       fmt.Sprintf("%s %s%s", formatCode(issue.PackageName, nil), msgSuffix, userSuppliedMsgSuffix),
				FromLinter: depguardLinterName,
			}, pass),
		)
	}

	return res, nil
}
