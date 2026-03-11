// (c) Copyright gosec's authors
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

// Package goanalysis provides a standard golang.org/x/tools/go/analysis.Analyzer for gosec.
package goanalysis

import (
	"fmt"
	"go/token"
	"io"
	"log"
	"strconv"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/packages"

	"github.com/securego/gosec/v2"
	"github.com/securego/gosec/v2/analyzers"
	"github.com/securego/gosec/v2/issue"
	"github.com/securego/gosec/v2/rules"
)

const Doc = `gosec is a static analysis tool that scans Go code for security problems.`

// Analyzer is the standard go/analysis Analyzer for gosec.
var Analyzer = &analysis.Analyzer{
	Name:     "gosec",
	Doc:      Doc,
	Run:      run,
	Requires: []*analysis.Analyzer{buildssa.Analyzer},
}

var (
	flagIncludeRules     string
	flagExcludeRules     string
	flagExcludeGenerated bool
	flagMinSeverity      string
	flagMinConfidence    string
)

//nolint:gochecknoinits // Required for go/analysis Analyzer flag registration
func init() {
	Analyzer.Flags.StringVar(&flagIncludeRules, "include", "", "Comma-separated list of rule IDs to include (e.g., G101,G102)")
	Analyzer.Flags.StringVar(&flagExcludeRules, "exclude", "", "Comma-separated list of rule IDs to exclude (e.g., G104)")
	Analyzer.Flags.BoolVar(&flagExcludeGenerated, "exclude-generated", true, "Exclude generated code from analysis")
	Analyzer.Flags.StringVar(&flagMinSeverity, "severity", "low", "Minimum severity: low, medium, or high")
	Analyzer.Flags.StringVar(&flagMinConfidence, "confidence", "low", "Minimum confidence: low, medium, or high")
}

func run(pass *analysis.Pass) (any, error) {
	// Create gosec config and analyzer
	config := gosec.NewConfig()
	logger := log.New(io.Discard, "", 0) // Discard gosec's verbose logging
	gosecAnalyzer := gosec.NewAnalyzer(config, false, flagExcludeGenerated, false, 1, logger)

	// Build filters from include/exclude flags
	ruleFilters := buildFilters(flagIncludeRules, flagExcludeRules, rules.NewRuleFilter)
	analyzerFilters := buildFilters(flagIncludeRules, flagExcludeRules, analyzers.NewAnalyzerFilter)

	// Load rules and analyzers
	ruleList := rules.Generate(false, ruleFilters...)
	ruleBuilders, ruleSuppressed := ruleList.RulesInfo()
	gosecAnalyzer.LoadRules(ruleBuilders, ruleSuppressed)

	analyzerList := analyzers.Generate(false, analyzerFilters...)
	analyzerDefs, analyzerSuppressed := analyzerList.AnalyzersInfo()
	gosecAnalyzer.LoadAnalyzers(analyzerDefs, analyzerSuppressed)

	// Convert analysis.Pass to packages.Package
	pkg := convertPassToPackage(pass)

	// Run gosec AST-based rules on the package
	// This populates context.Ignores with nosec suppressions from comments
	gosecAnalyzer.CheckRules(pkg)

	// Run SSA-based analyzers using the cached SSA result provided by the analysis framework
	// This reuses the SSA already built, maintaining cache efficiency
	// Both AST and SSA issues will respect nosec comments via gosec's updateIssues()
	if ssaResult, ok := pass.ResultOf[buildssa.Analyzer].(*buildssa.SSA); ok && ssaResult != nil {
		gosecAnalyzer.CheckAnalyzersWithSSA(pkg, ssaResult)
	}

	// Get all results (both AST and SSA, with nosec filtering already applied)
	issues, _, _ := gosecAnalyzer.Report()

	// Report issues as diagnostics, filtering by severity and confidence
	minSev, err := parseScore(flagMinSeverity)
	if err != nil {
		return nil, fmt.Errorf("invalid severity %q: %w", flagMinSeverity, err)
	}
	minConf, err := parseScore(flagMinConfidence)
	if err != nil {
		return nil, fmt.Errorf("invalid confidence %q: %w", flagMinConfidence, err)
	}

	for _, iss := range issues {
		if iss.Severity < minSev || iss.Confidence < minConf {
			continue
		}

		pos := parsePosition(pass.Fset, iss)
		what := iss.What
		if iss.Cwe != nil && iss.Cwe.ID != "" {
			what = fmt.Sprintf("[%s] %s", iss.Cwe.SprintID(), iss.What)
		}
		msg := fmt.Sprintf("%s: %s (Severity: %s, Confidence: %s)", iss.RuleID, what, iss.Severity.String(), iss.Confidence.String())

		// If we can't locate the issue, report it anyway but note the location problem
		if pos == token.NoPos {
			msg = fmt.Sprintf("%s [unable to locate %s:%s]", msg, iss.File, iss.Line)
		}

		pass.Report(analysis.Diagnostic{
			Pos:      pos,
			Category: iss.RuleID,
			Message:  msg,
		})
	}

	return nil, nil
}

// convertPassToPackage converts an analysis.Pass to a packages.Package
// that gosec expects. This allows us to reuse gosec's existing analysis logic.
func convertPassToPackage(pass *analysis.Pass) *packages.Package {
	pkg := &packages.Package{
		Name:       pass.Pkg.Name(),
		Fset:       pass.Fset,
		Syntax:     pass.Files,
		Types:      pass.Pkg,
		TypesInfo:  pass.TypesInfo,
		TypesSizes: pass.TypesSizes,
	}

	// Populate file names for the package
	pkg.CompiledGoFiles = make([]string, len(pass.Files))
	for i, f := range pass.Files {
		pkg.CompiledGoFiles[i] = pass.Fset.File(f.Pos()).Name()
	}

	return pkg
}

// buildFilters creates include/exclude filters from comma-separated rule IDs
func buildFilters[T any](include, exclude string, newFilter func(bool, ...string) T) []T {
	var filters []T
	if include != "" {
		if ids := parseRuleIDs(include); len(ids) > 0 {
			filters = append(filters, newFilter(false, ids...))
		}
	}
	if exclude != "" {
		if ids := parseRuleIDs(exclude); len(ids) > 0 {
			filters = append(filters, newFilter(true, ids...))
		}
	}
	return filters
}

// parseRuleIDs parses a comma-separated list of rule IDs
func parseRuleIDs(s string) []string {
	parts := strings.Split(s, ",")
	ids := make([]string, 0, len(parts))
	for _, p := range parts {
		if id := strings.TrimSpace(p); id != "" {
			ids = append(ids, id)
		}
	}
	return ids
}

// parseScore converts a severity/confidence string to issue.Score
func parseScore(s string) (issue.Score, error) {
	switch strings.ToLower(s) {
	case "high":
		return issue.High, nil
	case "medium":
		return issue.Medium, nil
	case "low":
		return issue.Low, nil
	default:
		return issue.Low, fmt.Errorf("must be low, medium, or high")
	}
}

// parsePosition converts a gosec issue location to a token.Pos
func parsePosition(fset *token.FileSet, iss *issue.Issue) token.Pos {
	var file *token.File
	fset.Iterate(func(f *token.File) bool {
		if f.Name() == iss.File {
			file = f
			return false
		}
		return true
	})

	if file == nil {
		return token.NoPos
	}

	// Handle line ranges (e.g., "28-34") by using the start line
	lineStr := iss.Line
	if idx := strings.Index(lineStr, "-"); idx > 0 {
		lineStr = lineStr[:idx]
	}

	line, err := strconv.Atoi(lineStr)
	if err != nil || line < 1 || line > file.LineCount() {
		return token.NoPos
	}

	lineStart := file.LineStart(line)

	// Add column offset if available (column is 1-based)
	col, err := strconv.Atoi(iss.Col)
	if err != nil || col < 1 {
		return lineStart
	}

	// Calculate position: lineStart is the position of the first character,
	// so we add (col - 1) to get the correct column position
	pos := lineStart + token.Pos(col-1)

	// Ensure we don't exceed file bounds
	if int(pos) > file.Base()+file.Size() {
		return lineStart
	}

	return pos
}
