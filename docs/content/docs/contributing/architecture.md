---
title: Architecture
weight: 3
---

There are the following golangci-lint execution steps:

```mermaid
graph LR
    init[Init]
    loadPackages[Load packages]
    runLinters[Run linters]
    postprocess[Postprocess issues]
    print[Print issues]

    init --> loadPackages --> runLinters --> postprocess --> print

```

## Init

The configuration is loaded from file and flags by `config.Loader` inside `PersistentPreRun` (or `PreRun`) of the commands that require configuration.

The linter database (`linterdb.Manager`) is fill based on the configuration:
- The linters ("internals" and plugins) are built by `linterdb.LinterBuilder` and `linterdb.PluginBuilder` builders.
- The configuration is validated by `linterdb.Validator`.

## Load Packages

Loading packages is listing all packages and their recursive dependencies for analysis.  
Also, depending on the enabled linters set some parsing of the source code can be performed at this step.

Packages loading starts here:

```go {base_url="https://github.com/golangci/golangci-lint/blob/main/", filename="pkg/lint/package.go"}
func (l *PackageLoader) Load(ctx context.Context, linters []*linter.Config) (pkgs, deduplicatedPkgs []*packages.Package, err error) {
	loadMode := findLoadMode(linters)

	pkgs, err = l.loadPackages(ctx, loadMode)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load packages: %w", err)
	}
// ...
```

First, we find a load mode as union of load modes for all enabled linters.
We use [go/packages](https://pkg.go.dev/golang.org/x/tools/go/packages) for packages loading and use it's enum `packages.Need*` for load modes.
Load mode sets which data does a linter needs for execution.

A linter that works only with AST need minimum of information: only filenames and AST.

There is no need for packages dependencies or type information.
AST is built during `go/analysis` execution to reduce memory usage.
Such AST-based linters are configured with the following code:

```go {base_url="https://github.com/golangci/golangci-lint/blob/main/", filename="pkg/lint/linter/config.go"}
func (lc *Config) WithLoadFiles() *Config {
	lc.LoadMode |= packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles
	return lc
}
```

If a linter uses `go/analysis` and needs type information, we need to extract more data by `go/packages`:

```go {base_url="https://github.com/golangci/golangci-lint/blob/main/", filename="pkg/lint/linter/config.go"}
func (lc *Config) WithLoadForGoAnalysis() *Config {
	lc = lc.WithLoadFiles()
	lc.LoadMode |= packages.NeedImports | packages.NeedDeps | packages.NeedExportFile | packages.NeedTypesSizes
	lc.IsSlow = true
	return lc
}
```

After finding a load mode, we run `go/packages`: 
he library get list of dirs (or `./...` as the default value) as input and outputs list of packages and requested information about them:
filenames, type information, AST, etc.

## Run Linters

First, we need to find all enabled linters. All linters are registered here:

```go {base_url="https://github.com/golangci/golangci-lint/blob/main/", filename="pkg/lint/lintersdb/builder_linter.go"}
func (b LinterBuilder) Build(cfg *config.Config) []*linter.Config {
	// ...
	return []*linter.Config{
		// ...
		linter.NewConfig(golinters.NewBodyclose()).
			WithSince("v1.18.0").
			WithLoadForGoAnalysis().
			WithURL("https://github.com/timakin/bodyclose"),
		// ...
		linter.NewConfig(golinters.NewGovet(govetCfg)).
			WithGroups(config.GroupStandard).
			WithSince("v1.0.0").
			WithLoadForGoAnalysis().
			WithURL("https://pkg.go.dev/cmd/vet"),
		// ...
	}
}
```

We filter requested in config and command-line linters in `EnabledSet`:

```go {base_url="https://github.com/golangci/golangci-lint/blob/main/", filename="pkg/lint/lintersdb/manager.go"}
func (m *Manager) GetEnabledLintersMap() (map[string]*linter.Config, error)
```

We merge enabled linters into one `MetaLinter` to improve execution time if we can:

```go {base_url="https://github.com/golangci/golangci-lint/blob/main/", filename="pkg/lint/lintersdb/manager.go"}
// GetOptimizedLinters returns enabled linters after optimization (merging) of multiple linters into a fewer number of linters.
// E.g. some go/analysis linters can be optimized into one metalinter for data reuse and speed up.
func (m *Manager) GetOptimizedLinters() ([]*linter.Config, error) {
	// ...
	m.combineGoAnalysisLinters(resultLintersSet)
	// ...
}
```

The `MetaLinter` just stores all merged linters inside to run them at once:

```go {base_url="https://github.com/golangci/golangci-lint/blob/main/", filename="pkg/goanalysis/metalinter.go"}
type MetaLinter struct {
	linters              []*Linter
	analyzerToLinterName map[*analysis.Analyzer]string
}
```

Currently, all linters except `unused` can be merged into this meta linter.  
The `unused` isn't merged because it has high memory usage.

Linters execution starts in `runAnalyzers`.
It's the most complex part of the golangci-lint.
We use custom [go/analysis](https://pkg.go.dev/golang.org/x/tools/go/analysis) runner there.  
It runs as much as it can in parallel.
It lazy-loads as much as it can to reduce memory usage.
Also, it sets all heavyweight data to `nil` as becomes unneeded to save memory.

We don't use existing [multichecker](https://pkg.go.dev/golang.org/x/tools/go/analysis/multichecker) because
it doesn't use caching and doesn't have some important performance optimizations.

All found by linters issues are represented with `result.Issue` struct:

```go {base_url="https://github.com/golangci/golangci-lint/blob/main/", filename="pkg/result/issue.go"}
type Issue struct {
	FromLinter string
	Text       string

	Severity string

	// Source lines of a code with the issue to show
	SourceLines []string

	// Pkg is needed for proper caching of linting results
	Pkg *packages.Package `json:"-"`

	Pos token.Position

	LineRange *Range `json:",omitempty"`

	// HunkPos is used only when golangci-lint is run over a diff
	HunkPos int `json:",omitempty"`

	// If we know how to fix the issue we can provide replacement lines
	SuggestedFixes []analysis.SuggestedFix `json:",omitempty"`

	// If we are expecting a nolint (because this is from nolintlint), record the expected linter
	ExpectNoLint         bool
	ExpectedNoLintLinter string

	// ...
}
```

## Postprocess Issues

We have an abstraction of `result.Processor` to postprocess found issues:

<!--
$ tree -L 1 ./pkg/result/processors/ | grep -v test
-->

{{< filetree/container >}}
  {{< filetree/folder name="./pkg/result/processors/" >}}
	{{< filetree/file name="cgo.go" >}}
	{{< filetree/file name="diff.go" >}}
	{{< filetree/file name="exclusion_generated_file_filter.go" >}}
	{{< filetree/file name="exclusion_generated_file_matcher.go" >}}
	{{< filetree/file name="exclusion_paths.go" >}}
	{{< filetree/file name="exclusion_presets.go" >}}
	{{< filetree/file name="exclusion_rules.go" >}}
	{{< filetree/file name="filename_unadjuster.go" >}}
	{{< filetree/file name="fixer.go" >}}
	{{< filetree/file name="identifier_marker.go" >}}
	{{< filetree/file name="invalid_issue.go" >}}
	{{< filetree/file name="issues.go" >}}
	{{< filetree/file name="max_from_linter.go" >}}
	{{< filetree/file name="max_per_file_from_linter.go" >}}
	{{< filetree/file name="max_same_issues.go" >}}
	{{< filetree/file name="nolint_filter.go" >}}
	{{< filetree/file name="path_absoluter.go" >}}
	{{< filetree/file name="path_prettifier.go" >}}
	{{< filetree/file name="path_relativity.go" >}}
	{{< filetree/file name="path_shortener.go" >}}
	{{< filetree/file name="processor.go" >}}
	{{< filetree/file name="severity.go" >}}
	{{< filetree/file name="sort_results.go" >}}
	{{< filetree/file name="source_code.go" >}}
	{{< filetree/file name="uniq_by_line.go" >}}
  {{< /filetree/folder >}}
{{< /filetree/container >}}

The abstraction is simple:

```go {base_url="https://github.com/golangci/golangci-lint/blob/main/", filename="pkg/result/processors/processor.go"}
type Processor interface {
	Process(issues []*result.Issue) ([]*result.Issue, error)
	Name() string
	Finish()
}
```

A processor can hide issues (`nolint`, `exclude`) or change issues (`path_prettifier`).

## Print Issues

We have an abstraction for printing found issues.

<!--
$ tree -L 1 ./pkg/printers/ | grep -v test
-->

{{< filetree/container >}}
  {{< filetree/folder name="./pkg/printers/" >}}
    {{< filetree/file name="checkstyle.go" >}}
    {{< filetree/file name="codeclimate.go" >}}
    {{< filetree/file name="html.go" >}}
    {{< filetree/file name="json.go" >}}
    {{< filetree/file name="junitxml.go" >}}
    {{< filetree/file name="printer.go" >}}
    {{< filetree/file name="sarif.go" >}}
    {{< filetree/file name="tab.go" >}}
    {{< filetree/file name="teamcity.go" >}}
    {{< filetree/file name="text.go" >}}
  {{< /filetree/folder >}}
{{< /filetree/container >}}


