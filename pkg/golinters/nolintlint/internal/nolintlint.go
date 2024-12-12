// Package internal provides a linter to ensure that all //nolint directives are followed by explanations
package internal

import (
	"regexp"
	"strings"

	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/goanalysis"
	"github.com/golangci/golangci-lint/pkg/result"
)

const LinterName = "nolintlint"

const (
	NeedsMachineOnly Needs = 1 << iota
	NeedsSpecific
	NeedsExplanation
	NeedsUnused
	NeedsAll = NeedsMachineOnly | NeedsSpecific | NeedsExplanation
)

type Needs uint

var commentPattern = regexp.MustCompile(`^//\s*(nolint)(:\s*[\w-]+\s*(?:,\s*[\w-]+\s*)*)?\b`)

// matches a complete nolint directive
var fullDirectivePattern = regexp.MustCompile(`^//\s*nolint(?::(\s*[\w-]+\s*(?:,\s*[\w-]+\s*)*))?\s*(//.*)?\s*\n?$`)

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
		needs:           needs | NeedsMachineOnly,
		excludeByLinter: excludeByName,
	}, nil
}

var (
	leadingSpacePattern      = regexp.MustCompile(`^//(\s*)`)
	trailingBlankExplanation = regexp.MustCompile(`\s*(//\s*)?$`)
)

//nolint:funlen,gocyclo // the function is going to be refactored in the future
func (l Linter) Run(pass *analysis.Pass) ([]goanalysis.Issue, error) {
	var issues []goanalysis.Issue

	for _, file := range pass.Files {
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

				directiveWithOptionalLeadingSpace := "//"
				if leadingSpace != "" {
					directiveWithOptionalLeadingSpace += " "
				}

				split := strings.Split(strings.SplitN(comment.Text, ":", 2)[0], "//")
				directiveWithOptionalLeadingSpace += strings.TrimSpace(split[1])

				pos := pass.Fset.Position(comment.Pos())
				end := pass.Fset.Position(comment.End())

				// check for, report and eliminate leading spaces, so we can check for other issues
				if leadingSpace != "" {
					removeWhitespace := &result.Replacement{
						Inline: &result.InlineFix{
							StartCol: pos.Column + 1,
							Length:   len(leadingSpace),
						},
					}

					if (l.needs & NeedsMachineOnly) != 0 {
						issue := &result.Issue{
							FromLinter:  LinterName,
							Text:        formatNotMachine(comment.Text),
							Pos:         pos,
							Replacement: removeWhitespace,
						}

						issues = append(issues, goanalysis.NewIssue(issue, pass))
					} else if len(leadingSpace) > 1 {
						issue := &result.Issue{
							FromLinter:  LinterName,
							Text:        formatExtraLeadingSpace(comment.Text),
							Pos:         pos,
							Replacement: removeWhitespace,
						}

						issue.Replacement.Inline.NewString = " " // assume a single space was intended

						issues = append(issues, goanalysis.NewIssue(issue, pass))
					}
				}

				fullMatches := fullDirectivePattern.FindStringSubmatch(comment.Text)
				if len(fullMatches) == 0 {
					issue := &result.Issue{
						FromLinter: LinterName,
						Text:       formatParseError(comment.Text, directiveWithOptionalLeadingSpace),
						Pos:        pos,
					}

					issues = append(issues, goanalysis.NewIssue(issue, pass))

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
						issue := &result.Issue{
							FromLinter: LinterName,
							Text:       formatNotSpecific(comment.Text, directiveWithOptionalLeadingSpace),
							Pos:        pos,
						}

						issues = append(issues, goanalysis.NewIssue(issue, pass))
					}
				}

				// when detecting unused directives, we send all the directives through and filter them out in the nolint processor
				if (l.needs & NeedsUnused) != 0 {
					removeNolintCompletely := &result.Replacement{}

					startCol := pos.Column - 1

					if startCol == 0 {
						// if the directive starts from a new line, remove the line
						removeNolintCompletely.NeedOnlyDelete = true
					} else {
						removeNolintCompletely.Inline = &result.InlineFix{
							StartCol: startCol,
							Length:   end.Column - pos.Column,
						}
					}

					if len(linters) == 0 {
						issue := &result.Issue{
							FromLinter:   LinterName,
							Text:         formatUnusedCandidate(comment.Text, ""),
							Pos:          pos,
							ExpectNoLint: true,
							Replacement:  removeNolintCompletely,
						}

						issues = append(issues, goanalysis.NewIssue(issue, pass))
					} else {
						for _, linter := range linters {
							issue := &result.Issue{
								FromLinter:           LinterName,
								Text:                 formatUnusedCandidate(comment.Text, linter),
								Pos:                  pos,
								ExpectNoLint:         true,
								ExpectedNoLintLinter: linter,
							}

							// only offer replacement if there is a single linter
							// because of issues around commas and the possibility of all
							// linters being removed
							if len(linters) == 1 {
								issue.Replacement = removeNolintCompletely
							}

							issues = append(issues, goanalysis.NewIssue(issue, pass))
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

						issue := &result.Issue{
							FromLinter: LinterName,
							Text:       formatNoExplanation(comment.Text, fullDirectiveWithoutExplanation),
							Pos:        pos,
						}

						issues = append(issues, goanalysis.NewIssue(issue, pass))
					}
				}
			}
		}
	}

	return issues, nil
}
