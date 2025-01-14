package processors

import (
	"fmt"
	"go/parser"
	"go/token"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/golangci/golangci-lint/pkg/result"
)

const (
	AutogeneratedModeLax     = "lax"
	AutogeneratedModeStrict  = "strict"
	AutogeneratedModeDisable = "disable"
)

// The values must be in lowercase.
const (
	genCodeGenerated = "code generated"
	genDoNotEdit     = "do not edit"

	// Related to easyjson.
	genAutoFile = "autogenerated file"

	//nolint:lll // Long URL
	// Related to Swagger Codegen.
	// https://github.com/swagger-api/swagger-codegen/blob/61cfeac3b9d855b4eb8bffa0d118bece117bcb7d/modules/swagger-codegen/src/main/resources/go/partial_header.mustache#L16
	// https://github.com/swagger-api/swagger-codegen/issues/12358
	genSwaggerCodegen = "* generated by: swagger codegen "
)

var _ Processor = (*AutogeneratedExclude)(nil)

type fileSummary struct {
	generated bool
}

// AutogeneratedExclude filters generated files.
// - mode "lax": see [isGeneratedFileLax] documentation.
// - mode "strict": see [isGeneratedFileStrict] documentation.
// - mode "disable": skips this processor.
type AutogeneratedExclude struct {
	debugf logutils.DebugFunc

	mode          string
	strictPattern *regexp.Regexp

	fileSummaryCache map[string]*fileSummary
}

func NewAutogeneratedExclude(mode string) *AutogeneratedExclude {
	return &AutogeneratedExclude{
		debugf:           logutils.Debug(logutils.DebugKeyAutogenExclude),
		mode:             mode,
		strictPattern:    regexp.MustCompile(`^// Code generated .* DO NOT EDIT\.$`),
		fileSummaryCache: map[string]*fileSummary{},
	}
}

func (*AutogeneratedExclude) Name() string {
	return "autogenerated_exclude"
}

func (p *AutogeneratedExclude) Process(issues []result.Issue) ([]result.Issue, error) {
	if p.mode == AutogeneratedModeDisable {
		return issues, nil
	}

	return filterIssuesErr(issues, p.shouldPassIssue)
}

func (*AutogeneratedExclude) Finish() {}

func (p *AutogeneratedExclude) shouldPassIssue(issue *result.Issue) (bool, error) {
	if filepath.Base(issue.FilePath()) == "go.mod" {
		return true, nil
	}

	// The file is already known.
	fs := p.fileSummaryCache[issue.FilePath()]
	if fs != nil {
		return !fs.generated, nil
	}

	fs = &fileSummary{}
	p.fileSummaryCache[issue.FilePath()] = fs

	if p.mode == AutogeneratedModeStrict {
		var err error
		fs.generated, err = p.isGeneratedFileStrict(issue.FilePath())
		if err != nil {
			return false, fmt.Errorf("failed to get doc (strict) of file %s: %w", issue.FilePath(), err)
		}
	} else {
		doc, err := getComments(issue.FilePath())
		if err != nil {
			return false, fmt.Errorf("failed to get doc (lax) of file %s: %w", issue.FilePath(), err)
		}

		fs.generated = p.isGeneratedFileLax(doc)
	}

	p.debugf("file %q is generated: %t", issue.FilePath(), fs.generated)

	// don't report issues for autogenerated files
	return !fs.generated, nil
}

// isGeneratedFileLax reports whether the source file is generated code.
// The function uses a bit laxer rules than isGeneratedFileStrict to match more generated code.
// See https://github.com/golangci/golangci-lint/issues/48 and https://github.com/golangci/golangci-lint/issues/72.
func (p *AutogeneratedExclude) isGeneratedFileLax(doc string) bool {
	markers := []string{genCodeGenerated, genDoNotEdit, genAutoFile, genSwaggerCodegen}

	doc = strings.ToLower(doc)

	for _, marker := range markers {
		if strings.Contains(doc, marker) {
			p.debugf("doc contains marker %q: file is generated", marker)

			return true
		}
	}

	p.debugf("doc of len %d doesn't contain any of markers: %s", len(doc), markers)

	return false
}

// isGeneratedFileStrict returns true if the source file has a line that matches the regular expression:
//
//	^// Code generated .* DO NOT EDIT\.$
//
// This line must appear before the first non-comment, non-blank text in the file.
// Based on https://go.dev/s/generatedcode.
func (p *AutogeneratedExclude) isGeneratedFileStrict(filePath string) (bool, error) {
	file, err := parser.ParseFile(token.NewFileSet(), filePath, nil, parser.PackageClauseOnly|parser.ParseComments)
	if err != nil {
		return false, fmt.Errorf("failed to parse file: %w", err)
	}

	if file == nil || len(file.Comments) == 0 {
		return false, nil
	}

	for _, comment := range file.Comments {
		if comment.Pos() > file.Package {
			return false, nil
		}

		for _, line := range comment.List {
			generated := p.strictPattern.MatchString(line.Text)
			if generated {
				p.debugf("doc contains ignore expression: file is generated")

				return true, nil
			}
		}
	}

	return false, nil
}

func getComments(filePath string) (string, error) {
	fset := token.NewFileSet()
	syntax, err := parser.ParseFile(fset, filePath, nil, parser.PackageClauseOnly|parser.ParseComments)
	if err != nil {
		return "", fmt.Errorf("failed to parse file: %w", err)
	}

	var docLines []string
	for _, c := range syntax.Comments {
		docLines = append(docLines, strings.TrimSpace(c.Text()))
	}

	return strings.Join(docLines, "\n"), nil
}
