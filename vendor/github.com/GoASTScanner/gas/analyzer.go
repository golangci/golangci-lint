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

// Modifications copyright (C) 2018 GolangCI

// Package gas holds the central scanning logic used by GAS
package gas

import (
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"go/types"
	"log"
	"os"
	"path"
	"reflect"
	"regexp"
	"strings"

	"path/filepath"

	"golang.org/x/tools/go/loader"
)

// The Context is populated with data parsed from the source code as it is scanned.
// It is passed through to all rule functions as they are called. Rules may use
// this data in conjunction withe the encoutered AST node.
type Context struct {
	FileSet  *token.FileSet
	Comments ast.CommentMap
	Info     *types.Info
	Pkg      *types.Package
	Root     *ast.File
	Config   map[string]interface{}
	Imports  *ImportTracker
	Ignores  []map[string]bool
}

// Metrics used when reporting information about a scanning run.
type Metrics struct {
	NumFiles int `json:"files"`
	NumLines int `json:"lines"`
	NumNosec int `json:"nosec"`
	NumFound int `json:"found"`
}

// Analyzer object is the main object of GAS. It has methods traverse an AST
// and invoke the correct checking rules as on each node as required.
type Analyzer struct {
	ignoreNosec bool
	ruleset     RuleSet
	context     *Context
	config      Config
	logger      *log.Logger
	issues      []*Issue
	stats       *Metrics
}

// NewAnalyzer builds a new anaylzer.
func NewAnalyzer(conf Config, logger *log.Logger) *Analyzer {
	ignoreNoSec := false
	if setting, err := conf.GetGlobal("nosec"); err == nil {
		ignoreNoSec = setting == "true" || setting == "enabled"
	}
	if logger == nil {
		logger = log.New(os.Stderr, "[gas]", log.LstdFlags)
	}
	return &Analyzer{
		ignoreNosec: ignoreNoSec,
		ruleset:     make(RuleSet),
		context:     &Context{},
		config:      conf,
		logger:      logger,
		issues:      make([]*Issue, 0, 16),
		stats:       &Metrics{},
	}
}

// LoadRules instantiates all the rules to be used when analyzing source
// packages
func (gas *Analyzer) LoadRules(ruleDefinitions map[string]RuleBuilder) {
	for id, def := range ruleDefinitions {
		r, nodes := def(id, gas.config)
		gas.ruleset.Register(r, nodes...)
	}
}

// Process kicks off the analysis process for a given package
func (gas *Analyzer) Process(buildTags []string, packagePaths ...string) error {
	ctx := build.Default
	ctx.BuildTags = append(ctx.BuildTags, buildTags...)
	packageConfig := loader.Config{
		Build:       &ctx,
		ParserMode:  parser.ParseComments,
		AllowErrors: true,
	}
	for _, packagePath := range packagePaths {
		abspath, err := filepath.Abs(packagePath)
		if err != nil {
			return err
		}
		if _, err := os.Stat(abspath); os.IsNotExist(err) {
			gas.logger.Printf("Skipping: %s. Path doesn't exist.", abspath)
			continue
		}
		gas.logger.Println("Searching directory:", abspath)

		basePackage, err := build.Default.ImportDir(packagePath, build.ImportComment)
		if err != nil {
			return err
		}

		var packageFiles []string
		for _, filename := range basePackage.GoFiles {
			packageFiles = append(packageFiles, path.Join(packagePath, filename))
		}

		packageConfig.CreateFromFilenames(basePackage.Name, packageFiles...)
	}

	builtPackage, err := packageConfig.Load()
	if err != nil {
		return err
	}

	gas.ProcessProgram(builtPackage)
	return nil
}

// ProcessProgram kicks off the analysis process for a given program
func (gas *Analyzer) ProcessProgram(builtPackage *loader.Program) {
	for _, pkg := range builtPackage.InitialPackages() {
		gas.logger.Println("Checking package:", pkg.String())
		for _, file := range pkg.Files {
			gas.logger.Println("Checking file:", builtPackage.Fset.File(file.Pos()).Name())
			gas.context.FileSet = builtPackage.Fset
			gas.context.Config = gas.config
			gas.context.Comments = ast.NewCommentMap(gas.context.FileSet, file, file.Comments)
			gas.context.Root = file
			gas.context.Info = &pkg.Info
			gas.context.Pkg = pkg.Pkg
			gas.context.Imports = NewImportTracker()
			gas.context.Imports.TrackPackages(gas.context.Pkg.Imports()...)
			ast.Walk(gas, file)
			gas.stats.NumFiles++
			gas.stats.NumLines += builtPackage.Fset.File(file.Pos()).LineCount()
		}
	}
}

// ignore a node (and sub-tree) if it is tagged with a "#nosec" comment
func (gas *Analyzer) ignore(n ast.Node) ([]string, bool) {
	if groups, ok := gas.context.Comments[n]; ok && !gas.ignoreNosec {
		for _, group := range groups {
			if strings.Contains(group.Text(), "#nosec") {
				gas.stats.NumNosec++

				// Pull out the specific rules that are listed to be ignored.
				re := regexp.MustCompile("(G\\d{3})")
				matches := re.FindAllStringSubmatch(group.Text(), -1)

				// If no specific rules were given, ignore everything.
				if matches == nil || len(matches) == 0 {
					return nil, true
				}

				// Find the rule IDs to ignore.
				var ignores []string
				for _, v := range matches {
					ignores = append(ignores, v[1])
				}
				return ignores, false
			}
		}
	}
	return nil, false
}

// Visit runs the GAS visitor logic over an AST created by parsing go code.
// Rule methods added with AddRule will be invoked as necessary.
func (gas *Analyzer) Visit(n ast.Node) ast.Visitor {
	// If we've reached the end of this branch, pop off the ignores stack.
	if n == nil {
		if len(gas.context.Ignores) > 0 {
			gas.context.Ignores = gas.context.Ignores[1:]
		}
		return gas
	}

	// Get any new rule exclusions.
	ignoredRules, ignoreAll := gas.ignore(n)
	if ignoreAll {
		return nil
	}

	// Now create the union of exclusions.
	ignores := make(map[string]bool, 0)
	if len(gas.context.Ignores) > 0 {
		for k, v := range gas.context.Ignores[0] {
			ignores[k] = v
		}
	}

	for _, v := range ignoredRules {
		ignores[v] = true
	}

	// Push the new set onto the stack.
	gas.context.Ignores = append([]map[string]bool{ignores}, gas.context.Ignores...)

	// Track aliased and initialization imports
	gas.context.Imports.TrackImport(n)

	for _, rule := range gas.ruleset.RegisteredFor(n) {
		if _, ok := ignores[rule.ID()]; ok {
			continue
		}
		issue, err := rule.Match(n, gas.context)
		if err != nil {
			file, line := GetLocation(n, gas.context)
			file = path.Base(file)
			gas.logger.Printf("Rule error: %v => %s (%s:%d)\n", reflect.TypeOf(rule), err, file, line)
		}
		if issue != nil {
			gas.issues = append(gas.issues, issue)
			gas.stats.NumFound++
		}
	}
	return gas
}

// Report returns the current issues discovered and the metrics about the scan
func (gas *Analyzer) Report() ([]*Issue, *Metrics) {
	return gas.issues, gas.stats
}

// Reset clears state such as context, issues and metrics from the configured analyzer
func (gas *Analyzer) Reset() {
	gas.context = &Context{}
	gas.issues = make([]*Issue, 0, 16)
	gas.stats = &Metrics{}
}
