// Package nolintlint provides a linter to ensure that all //nolint directives are followed by explanations
package nolintlint

import (
	"fmt"
	"go/ast"
	"go/token"
	"regexp"
	"strings"

	"github.com/golangci/golangci-lint/pkg/result"
)

type baseIssue struct {
	fullDirective string
	position      token.Position
	replacement   *result.Replacement
}

func (b baseIssue) Position() token.Position {
	return b.position
}

func (b baseIssue) Replacement() *result.Replacement {
	return b.replacement
}

type ExtraLeadingSpace struct {
	baseIssue
}

func (i ExtraLeadingSpace) Details() string {
	return fmt.Sprintf("directive `%s` should not have more than one leading space", i.fullDirective)
}

func (i ExtraLeadingSpace) String() string { return toString(i) }

type NotSpecific struct {
	baseIssue
}

func (i NotSpecific) Details() string {
	return fmt.Sprintf("directive `%s` should mention specific linter such as `//nolint:my-linter`",
		i.fullDirective)
}

func (i NotSpecific) String() string { return toString(i) }

type ParseError struct {
	baseIssue
}

func (i ParseError) Details() string {
	return fmt.Sprintf("directive `%s` should match `//nolint[:<comma-separated-linters>] [// <explanation>]`",
		i.fullDirective)
}

func (i ParseError) String() string { return toString(i) }

type NoExplanation struct {
	baseIssue
	fullDirectiveWithoutExplanation string
}

//nolint:gocritic // TODO(ldez) must be change in the future.
func (i NoExplanation) Details() string {
	return fmt.Sprintf("directive `%s` should provide explanation such as `%s // this is why`",
		i.fullDirective, i.fullDirectiveWithoutExplanation)
}

func (i NoExplanation) String() string { return toString(i) }

type UnusedCandidate struct {
	baseIssue
	ExpectedLinter string
}

//nolint:gocritic // TODO(ldez) must be change in the future.
func (i UnusedCandidate) Details() string {
	details := fmt.Sprintf("directive `%s` is unused", i.fullDirective)
	if i.ExpectedLinter != "" {
		details += fmt.Sprintf(" for linter %q", i.ExpectedLinter)
	}
	return details
}

func (i UnusedCandidate) String() string { return toString(i) }

func toString(issue Issue) string {
	return fmt.Sprintf("%s at %s", issue.Details(), issue.Position())
}

type Issue interface {
	Details() string
	Position() token.Position
	String() string
	Replacement() *result.Replacement
}

type Needs uint

const (
	// Deprecated: NeedsMachineOnly is deprecated as leading spaces are no longer allowed,
	// making this condition always true. Consumers should adjust their code to assume
	// this as the default behavior and no longer rely on NeedsMachineOnly.
	NeedsMachineOnly Needs = 1 << iota
	NeedsSpecific
	NeedsExplanation
	NeedsUnused
	NeedsAll = NeedsSpecific | NeedsExplanation
)

var commentPattern = regexp.MustCompile(`^//\s*(nolint)(:\s*[\w-]+\s*(?:,\s*[\w-]+\s*)*)?\b`)

// matches a complete nolint directive
var fullDirectivePattern = regexp.MustCompile(`^//nolint(?::([\w-]+(?:,[\w-]+)*))?(?: (//.*))?\n?$`)

type Linter struct {
	needs           Needs // indicates which linter checks to perform
	excludeByLinter map[string]bool
}

// NewLinter creates a linter that enforces that the provided directives fulfill the provided requirements
func NewLinter(needs Needs, excludes []string) (*Linter, error) {
	excludeByName := make(map[string]bool)
	for _, e := range excludes {
		excludeByName[e] = true
	}

	return &Linter{
		needs:           needs,
		excludeByLinter: excludeByName,
	}, nil
}

var (
	leadingSpacePattern      = regexp.MustCompile(`^//(\s*)`)
	trailingBlankExplanation = regexp.MustCompile(`\s*(//\s*)?$`)
)

//nolint:funlen,gocyclo // the function is going to be refactored in the future
func (l Linter) Run(fset *token.FileSet, nodes ...ast.Node) ([]Issue, error) {
	var issues []Issue

	for _, node := range nodes {
		file, ok := node.(*ast.File)
		if !ok {
			continue
		}

		for _, c := range file.Comments {
			for _, comment := range c.List {
				if !commentPattern.MatchString(comment.Text) {
					continue
				}

				// check for a space between the "//" and the directive
				leadingSpaceMatches := leadingSpacePattern.FindStringSubmatch(comment.Text)

				var leadingSpace string
				if len(leadingSpaceMatches) > 0 {
					leadingSpace = leadingSpaceMatches[1]
				}

				pos := fset.Position(comment.Pos())
				end := fset.Position(comment.End())

				base := baseIssue{
					fullDirective: comment.Text,
					position:      pos,
				}

				fullMatches := fullDirectivePattern.FindStringSubmatch(comment.Text)
				if len(fullMatches) == 0 {
					issues = append(issues, ParseError{baseIssue: base})
					continue
				}

				lintersText, explanation := fullMatches[1], fullMatches[2]

				var linters []string
				if lintersText != "" && !strings.HasPrefix(lintersText, "all") {
					lls := strings.Split(lintersText, ",")
					linters = make([]string, 0, len(lls))
					rangeStart := (pos.Column - 1) + len("//") + len(leadingSpace) + len("nolint:")
					for i, ll := range lls {
						rangeEnd := rangeStart + len(ll)
						if i < len(lls)-1 {
							rangeEnd++ // include trailing comma
						}
						trimmedLinterName := strings.TrimSpace(ll)
						if trimmedLinterName != "" {
							linters = append(linters, trimmedLinterName)
						}
						rangeStart = rangeEnd
					}
				}

				if (l.needs & NeedsSpecific) != 0 {
					if len(linters) == 0 {
						issues = append(issues, NotSpecific{baseIssue: base})
					}
				}

				// when detecting unused directives, we send all the directives through and filter them out in the nolint processor
				if (l.needs & NeedsUnused) != 0 {
					removeNolintCompletely := &result.Replacement{
						Inline: &result.InlineFix{
							StartCol:  pos.Column - 1,
							Length:    end.Column - pos.Column,
							NewString: "",
						},
					}

					if len(linters) == 0 {
						issue := UnusedCandidate{baseIssue: base}
						issue.replacement = removeNolintCompletely
						issues = append(issues, issue)
					} else {
						for _, linter := range linters {
							issue := UnusedCandidate{baseIssue: base, ExpectedLinter: linter}
							// only offer replacement if there is a single linter
							// because of issues around commas and the possibility of all
							// linters being removed
							if len(linters) == 1 {
								issue.replacement = removeNolintCompletely
							}
							issues = append(issues, issue)
						}
					}
				}

				if (l.needs&NeedsExplanation) != 0 && (explanation == "" || strings.TrimSpace(explanation) == "//") {
					needsExplanation := len(linters) == 0 // if no linters are mentioned, we must have explanation
					// otherwise, check if we are excluding all the mentioned linters
					for _, ll := range linters {
						if !l.excludeByLinter[ll] { // if a linter does require explanation
							needsExplanation = true
							break
						}
					}

					if needsExplanation {
						fullDirectiveWithoutExplanation := trailingBlankExplanation.ReplaceAllString(comment.Text, "")
						issues = append(issues, NoExplanation{
							baseIssue:                       base,
							fullDirectiveWithoutExplanation: fullDirectiveWithoutExplanation,
						})
					}
				}
			}
		}
	}

	return issues, nil
}
