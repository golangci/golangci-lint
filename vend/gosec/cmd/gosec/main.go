// (c) Copyright 2016 Hewlett Packard Enterprise Development LP
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"sort"
	"strings"

	"github.com/securego/gosec/v2"
	"github.com/securego/gosec/v2/analyzers"
	"github.com/securego/gosec/v2/autofix"
	"github.com/securego/gosec/v2/cmd/vflag"
	"github.com/securego/gosec/v2/issue"
	"github.com/securego/gosec/v2/report"
	"github.com/securego/gosec/v2/rules"
)

const (
	usageText = `
gosec - Golang security checker

gosec analyzes Go source code to look for common programming mistakes that
can lead to security problems.

VERSION: %s
GIT TAG: %s
BUILD DATE: %s

USAGE:

	# Check a single package
	$ gosec $GOPATH/src/github.com/example/project

	# Check all packages under the current directory and save results in
	# json format.
	$ gosec -fmt=json -out=results.json ./...

	# Run a specific set of rules (by default all rules will be run):
	$ gosec -include=G101,G203,G401  ./...

	# Run all rules except the provided
	$ gosec -exclude=G101 $GOPATH/src/github.com/example/project/...

	# Exclude specific rules from specific paths
	$ gosec --exclude-rules="cmd/.*:G204,G304" ./...

	# Exclude all rules from scripts directory
	$ gosec --exclude-rules="scripts/.*:*" ./...
`
	// Environment variable for AI API key.
	aiAPIKeyEnv = "GOSEC_AI_API_KEY" // #nosec G101

	// Exit codes
	exitSuccess = 0
	exitFailure = 1
)

type arrayFlags []string

func (a *arrayFlags) String() string {
	return strings.Join(*a, " ")
}

func (a *arrayFlags) Set(value string) error {
	*a = append(*a, value)
	return nil
}

var (
	// #nosec flag
	flagIgnoreNoSec = flag.Bool("nosec", false, "Ignores #nosec comments when set")

	// Path-based exclusions
	flagExcludeRules = flag.String("exclude-rules", "",
		`Path-based rule exclusions. Format: "path:rule1,rule2;path2:rule3"
Example: "cmd/.*:G204,G304;test/.*:G101"
Use "*" to exclude all rules for a path: "scripts/.*:*"`)

	// show ignored
	flagShowIgnored = flag.Bool("show-ignored", false, "If enabled, ignored issues are printed")

	// format output
	flagFormat = flag.String("fmt", "text", "Set output format. Valid options are: json, yaml, csv, junit-xml, html, sonarqube, golint, sarif or text")

	// #nosec alternative tag
	flagAlternativeNoSec = flag.String("nosec-tag", "", "Set an alternative string for #nosec. Some examples: #dontanalyze, #falsepositive")

	// flagEnableAudit enables audit mode
	flagEnableAudit = flag.Bool("enable-audit", false, "Enable audit mode")

	// output file
	flagOutput = flag.String("out", "", "Set output file for results")

	// config file
	flagConfig = flag.String("conf", "", "Path to optional config file")

	// quiet
	flagQuiet = flag.Bool("quiet", false, "Only show output when errors are found")

	// rules to explicitly include
	flagRulesInclude = flag.String("include", "", "Comma separated list of rules IDs to include. (see rule list)")

	// rules to explicitly exclude
	flagRulesExclude = vflag.ValidatedFlag{}

	// rules to explicitly exclude
	flagExcludeGenerated = flag.Bool("exclude-generated", false, "Exclude generated files")

	// log to file or stderr
	flagLogfile = flag.String("log", "", "Log messages to file rather than stderr")
	// sort the issues by severity
	flagSortIssues = flag.Bool("sort", true, "Sort issues by severity")

	// go build tags
	flagBuildTags = flag.String("tags", "", "Comma separated list of build tags")

	// fail by severity
	flagSeverity = flag.String("severity", "low", "Filter out the issues with a lower severity than the given value. Valid options are: low, medium, high")

	// fail by confidence
	flagConfidence = flag.String("confidence", "low", "Filter out the issues with a lower confidence than the given value. Valid options are: low, medium, high")

	// concurrency value
	flagConcurrency = flag.Int("concurrency", runtime.NumCPU(), "Concurrency value")

	// do not fail
	flagNoFail = flag.Bool("no-fail", false, "Do not fail the scanning, even if issues were found")

	// scan tests files
	flagScanTests = flag.Bool("tests", false, "Scan tests files")

	// print version and quit with exit code 0
	flagVersion = flag.Bool("version", false, "Print version and quit with exit code 0")

	// stdout the results as well as write it in the output file
	flagStdOut = flag.Bool("stdout", false, "Stdout the results as well as write it in the output file")

	// print the text report with color, this is enabled by default
	flagColor = flag.Bool("color", true, "Prints the text format report with colorization when it goes in the stdout")

	// append ./... to the target dir.
	flagRecursive = flag.Bool("r", false, "Appends \"./...\" to the target dir.")

	// overrides the output format when stdout the results while saving them in the output file
	flagVerbose = flag.String("verbose", "", "Overrides the output format when stdout the results while saving them in the output file.\nValid options are: json, yaml, csv, junit-xml, html, sonarqube, golint, sarif or text")

	// output suppression information for auditing purposes
	flagTrackSuppressions = flag.Bool("track-suppressions", false, "Output suppression information, including its kind and justification")

	// flagTerse shows only the summary of scan discarding all the logs
	flagTerse = flag.Bool("terse", false, "Shows only the results and summary")

	// AI platform provider to generate solutions to issues
	flagAiAPIProvider = flag.String("ai-api-provider", "", autofix.AIProviderFlagHelp)

	// key to implementing AI provider services
	flagAiAPIKey = flag.String("ai-api-key", "", "Key to access the AI API")

	// base URL for AI API (optional, for OpenAI-compatible APIs)
	flagAiBaseURL = flag.String("ai-base-url", "", "Base URL for AI API (e.g., for OpenAI-compatible services)")

	// skip SSL verification for AI API
	flagAiSkipSSL = flag.Bool("ai-skip-ssl", false, "Skip SSL certificate verification for AI API")

	// exclude the folders from scan
	flagDirsExclude arrayFlags

	logger *log.Logger
)

// #nosec
func usage() {
	usageText := fmt.Sprintf(usageText, Version, GitTag, BuildDate)
	fmt.Fprintln(os.Stderr, usageText)
	fmt.Fprint(os.Stderr, "OPTIONS:\n\n")
	flag.PrintDefaults()
	fmt.Fprint(os.Stderr, "\n\nRULES:\n\n")

	// sorted rule list for ease of reading
	rl := rules.Generate(*flagTrackSuppressions)
	al := analyzers.Generate(*flagTrackSuppressions)
	keys := make([]string, 0, len(rl.Rules)+len(al.Analyzers))
	for key := range rl.Rules {
		keys = append(keys, key)
	}
	for key := range al.Analyzers {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, k := range keys {
		var description string
		if rule, ok := rl.Rules[k]; ok {
			description = rule.Description
		} else if analyzer, ok := al.Analyzers[k]; ok {
			description = analyzer.Description
		}
		fmt.Fprintf(os.Stderr, "\t%s: %s\n", k, description)
	}
	fmt.Fprint(os.Stderr, "\n")
}

func loadConfig(configFile string) (gosec.Config, error) {
	config := gosec.NewConfig()
	if configFile != "" {
		// #nosec
		file, err := os.Open(configFile)
		if err != nil {
			return nil, err
		}
		defer file.Close() // #nosec G307
		if _, err := config.ReadFrom(file); err != nil {
			return nil, err
		}
	}
	if *flagIgnoreNoSec {
		config.SetGlobal(gosec.Nosec, "true")
	}
	if *flagShowIgnored {
		config.SetGlobal(gosec.ShowIgnored, "true")
	}
	if *flagAlternativeNoSec != "" {
		config.SetGlobal(gosec.NoSecAlternative, *flagAlternativeNoSec)
	}
	if *flagEnableAudit {
		config.SetGlobal(gosec.Audit, "true")
	}
	// set global option IncludeRules, when flag set or global option IncludeRules  is nil
	if v, _ := config.GetGlobal(gosec.IncludeRules); *flagRulesInclude != "" || v == "" {
		config.SetGlobal(gosec.IncludeRules, *flagRulesInclude)
	}
	// set global option ExcludeRules, when flag set or global option ExcludeRules  is nil
	if v, _ := config.GetGlobal(gosec.ExcludeRules); flagRulesExclude.String() != "" || v == "" {
		config.SetGlobal(gosec.ExcludeRules, flagRulesExclude.String())
	}
	return config, nil
}

func loadRules(include, exclude string) rules.RuleList {
	var filters []rules.RuleFilter
	if include != "" {
		logger.Printf("Including rules: %s", include)
		including := strings.Split(include, ",")
		filters = append(filters, rules.NewRuleFilter(false, including...))
	} else {
		logger.Println("Including rules: default")
	}

	if exclude != "" {
		logger.Printf("Excluding rules: %s", exclude)
		excluding := strings.Split(exclude, ",")
		filters = append(filters, rules.NewRuleFilter(true, excluding...))
	} else {
		logger.Println("Excluding rules: default")
	}
	return rules.Generate(*flagTrackSuppressions, filters...)
}

func loadAnalyzers(include, exclude string) *analyzers.AnalyzerList {
	var filters []analyzers.AnalyzerFilter
	if include != "" {
		logger.Printf("Including analyzers: %s", include)
		including := strings.Split(include, ",")
		filters = append(filters, analyzers.NewAnalyzerFilter(false, including...))
	} else {
		logger.Println("Including analyzers: default")
	}

	if exclude != "" {
		logger.Printf("Excluding analyzers: %s", exclude)
		excluding := strings.Split(exclude, ",")
		filters = append(filters, analyzers.NewAnalyzerFilter(true, excluding...))
	} else {
		logger.Println("Excluding analyzers: default")
	}
	return analyzers.Generate(*flagTrackSuppressions, filters...)
}

func getRootPaths(paths []string) ([]string, error) {
	rootPaths := make([]string, 0)
	for _, path := range paths {
		rootPath, err := gosec.RootPath(path)
		if err != nil {
			return nil, fmt.Errorf("failed to get the root path of the projects: %w", err)
		}
		rootPaths = append(rootPaths, rootPath)
	}
	return rootPaths, nil
}

// If verbose is defined it overwrites the defined format
// Otherwise the actual format is used
func getPrintedFormat(format string, verbose string) string {
	if verbose != "" {
		return verbose
	}
	return format
}

func printReport(format string, color bool, rootPaths []string, reportInfo *gosec.ReportInfo) error {
	return report.CreateReport(os.Stdout, format, color, rootPaths, reportInfo)
}

func saveReport(filename, format string, rootPaths []string, reportInfo *gosec.ReportInfo) error {
	outfile, err := os.Create(filename) // #nosec G304
	if err != nil {
		return err
	}
	defer outfile.Close() // #nosec G307
	return report.CreateReport(outfile, format, false, rootPaths, reportInfo)
}

func convertToScore(value string) (issue.Score, error) {
	value = strings.ToLower(value)
	switch value {
	case "low":
		return issue.Low, nil
	case "medium":
		return issue.Medium, nil
	case "high":
		return issue.High, nil
	default:
		return issue.Low, fmt.Errorf("provided value '%s' not valid. Valid options: low, medium, high", value)
	}
}

func filterIssues(issues []*issue.Issue, severity issue.Score, confidence issue.Score) ([]*issue.Issue, int) {
	result := make([]*issue.Issue, 0)
	trueIssues := 0
	for _, issue := range issues {
		if issue.Severity >= severity && issue.Confidence >= confidence {
			result = append(result, issue)
			if (!issue.NoSec || !*flagShowIgnored) && len(issue.Suppressions) == 0 {
				trueIssues++
			}
		}
	}
	return result, trueIssues
}

// computeExitCode determines the exit code based on issues found and noFail flag.
func computeExitCode(issues []*issue.Issue, errors map[string][]gosec.Error, noFail bool) int {
	nsi := 0
	for _, issue := range issues {
		if len(issue.Suppressions) == 0 {
			nsi++
		}
	}
	if (nsi > 0 || len(errors) > 0) && !noFail {
		return exitFailure
	}
	return exitSuccess
}

// buildPathExclusionFilter creates a PathExclusionFilter from config and CLI flags
func buildPathExclusionFilter(config gosec.Config, cliFlag string) (*gosec.PathExclusionFilter, error) {
	// Parse CLI exclude-rules
	cliRules, err := gosec.ParseCLIExcludeRules(cliFlag)
	if err != nil {
		return nil, fmt.Errorf("invalid --exclude-rules flag: %w", err)
	}

	// Get config file exclude-rules
	configRules, err := config.GetExcludeRules()
	if err != nil {
		return nil, fmt.Errorf("invalid exclude-rules in config: %w", err)
	}

	// Merge rules (CLI takes precedence)
	allRules := gosec.MergeExcludeRules(configRules, cliRules)

	// Create and return filter
	return gosec.NewPathExclusionFilter(allRules)
}

func main() {
	os.Exit(run())
}

func run() int {
	// Makes sure some version information is set
	prepareVersionInfo()

	// Setup usage description
	flag.Usage = usage

	// Setup the excluded folders from scan
	flag.Var(&flagDirsExclude, "exclude-dir", "Exclude folder from scan (can be specified multiple times)")
	if err := flag.Set("exclude-dir", "vendor"); err != nil {
		fmt.Fprintf(os.Stderr, "\nError: failed to exclude the %q directory from scan", "vendor")
	}
	if err := flag.Set("exclude-dir", "\\.git/"); err != nil {
		fmt.Fprintf(os.Stderr, "\nError: failed to exclude the %q directory from scan", "\\.git/")
	}

	// set for exclude
	flag.Var(&flagRulesExclude, "exclude", "Comma separated list of rules IDs to exclude. (see rule list)")

	// Parse command line arguments
	flag.Parse()

	if *flagVersion {
		fmt.Printf("Version: %s\nGit tag: %s\nBuild date: %s\n", Version, GitTag, BuildDate)
		return exitSuccess
	}

	// Ensure at least one file was specified or that the recursive -r flag was set.
	if flag.NArg() == 0 && !*flagRecursive {
		fmt.Fprintf(os.Stderr, "\nError: FILE [FILE...] or './...' or -r expected\n") // #nosec
		flag.Usage()
		return exitFailure
	}

	// Setup logging
	logWriter := os.Stderr
	if *flagLogfile != "" {
		var err error
		logWriter, err = os.Create(*flagLogfile)
		if err != nil {
			flag.Usage()
			log.Printf("failed to create log file: %v", err)
			return exitFailure
		}
		defer logWriter.Close() // #nosec
	}

	if *flagQuiet || *flagTerse {
		logger = log.New(io.Discard, "", 0)
	} else {
		logger = log.New(logWriter, "[gosec] ", log.LstdFlags)
	}

	// Initialize profiling after logger setup so it uses the same logger
	// (defers execute in LIFO order, so finishProfiling runs before logWriter.Close)
	profiler, err := initProfiling(logger)
	if err != nil {
		logger.Printf("failed to initialize profiling: %v", err)
		return exitFailure
	}
	defer finishProfiling(profiler)

	failSeverity, err := convertToScore(*flagSeverity)
	if err != nil {
		logger.Printf("Invalid severity value: %v", err)
		return exitFailure
	}

	failConfidence, err := convertToScore(*flagConfidence)
	if err != nil {
		logger.Printf("Invalid confidence value: %v", err)
		return exitFailure
	}

	// Load the analyzer configuration
	config, err := loadConfig(*flagConfig)
	if err != nil {
		logger.Printf("Failed to load config: %v", err)
		return exitFailure
	}

	// Load enabled rule definitions
	excludeRules, err := config.GetGlobal(gosec.ExcludeRules)
	if err != nil {
		logger.Printf("Failed to get exclude rules: %v", err)
		return exitFailure
	}
	includeRules, err := config.GetGlobal(gosec.IncludeRules)
	if err != nil {
		logger.Printf("Failed to get include rules: %v", err)
		return exitFailure
	}

	ruleList := loadRules(includeRules, excludeRules)

	analyzerList := loadAnalyzers(includeRules, excludeRules)

	if len(ruleList.Rules) == 0 && len(analyzerList.Analyzers) == 0 {
		logger.Print("No rules/analyzers are configured")
		return exitFailure
	}

	// Build path exclusion filter
	pathFilter, err := buildPathExclusionFilter(config, *flagExcludeRules)
	if err != nil {
		logger.Printf("Path exclusion filter error: %v", err)
		return exitFailure
	}

	// Create the analyzer
	analyzer := gosec.NewAnalyzer(config, *flagScanTests, *flagExcludeGenerated, *flagTrackSuppressions, *flagConcurrency, logger)
	analyzer.LoadRules(ruleList.RulesInfo())
	analyzer.LoadAnalyzers(analyzerList.AnalyzersInfo())

	excludedDirs := gosec.ExcludedDirsRegExp(flagDirsExclude)
	var packages []string

	paths := flag.Args()
	if len(paths) == 0 {
		paths = append(paths, "./...")
	}
	for _, path := range paths {
		pcks, err := gosec.PackagePaths(path, excludedDirs)
		if err != nil {
			logger.Printf("Failed to get package paths: %v", err)
			return exitFailure
		}
		packages = append(packages, pcks...)
	}

	if len(packages) == 0 {
		logger.Print("No packages found")
		return exitFailure
	}

	var buildTags []string
	if *flagBuildTags != "" {
		buildTags = strings.Split(*flagBuildTags, ",")
	}

	if err := analyzer.Process(buildTags, packages...); err != nil {
		logger.Printf("Analyzer error: %v", err)
		return exitFailure
	}

	// Collect the results
	issues, metrics, errors := analyzer.Report()

	// Apply path-based exclusions first
	var pathExcludedCount int
	issues, pathExcludedCount = pathFilter.FilterIssues(issues)
	if pathExcludedCount > 0 {
		logger.Printf("Excluded %d issues by path-based rules", pathExcludedCount)
	}

	// Sort the issue by severity
	if *flagSortIssues {
		sortIssues(issues)
	}

	// Filter the issues by severity and confidence
	var trueIssues int
	issues, trueIssues = filterIssues(issues, failSeverity, failConfidence)
	if metrics.NumFound != trueIssues {
		metrics.NumFound = trueIssues
	}

	// Exit quietly if nothing was found
	if len(issues) == 0 && *flagQuiet {
		return exitSuccess
	}

	// Create output report
	rootPaths, err := getRootPaths(flag.Args())
	if err != nil {
		logger.Printf("Failed to get root paths: %v", err)
		return exitFailure
	}

	reportInfo := gosec.NewReportInfo(issues, metrics, errors).WithVersion(Version)

	// Call AI request to solve the issues
	aiAPIKey := os.Getenv(aiAPIKeyEnv)
	if aiAPIKey == "" {
		aiAPIKey = *flagAiAPIKey
	}

	aiEnabled := *flagAiAPIProvider != ""

	if len(issues) > 0 && aiEnabled {
		err := autofix.GenerateSolution(*flagAiAPIProvider, aiAPIKey, *flagAiBaseURL, *flagAiSkipSSL, issues)
		if err != nil {
			logger.Print(err)
		}
	}

	if *flagOutput == "" || *flagStdOut {
		fileFormat := getPrintedFormat(*flagFormat, *flagVerbose)
		if err := printReport(fileFormat, *flagColor, rootPaths, reportInfo); err != nil {
			logger.Printf("Failed to print report: %v", err)
			return exitFailure
		}
	}
	if *flagOutput != "" {
		if err := saveReport(*flagOutput, *flagFormat, rootPaths, reportInfo); err != nil {
			logger.Printf("Failed to save report: %v", err)
			return exitFailure
		}
	}

	return computeExitCode(issues, errors, *flagNoFail)
}
