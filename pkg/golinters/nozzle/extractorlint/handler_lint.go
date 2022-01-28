package extractorlint

import (
	"fmt"
	"go/token"

	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/result"
)

// HandlerLint contains information for a parsed ElementHandler
type HandlerLint struct {
	Pos  token.Pos
	Name string

	Reg              *handlerFieldReg
	ID               *handlerField
	IDFunc           *handlerField
	IDMatchedFunc    *handlerField
	AfterID          *handlerField
	AfterIDFunc      *handlerField
	OnText           *handlerField
	OnTextFunc       *handlerField
	Pop              *handlerField
	PopFunc          *handlerField
	ProbableHandlers *handlerField
}

type handlerField struct {
	FieldName string
	Pos       token.Pos
	Errs      []error
}

func (hf *handlerField) error(err error) {
	hf.Errs = append(hf.Errs, err)
}

func newHandlerField(fieldName string, pos token.Pos) *handlerField {
	return &handlerField{
		FieldName: fieldName,
		Pos:       pos,
	}
}

type handlerFieldReg struct {
	*handlerField
	Name    string
	NamePos token.Pos
}

func (hl *HandlerLint) suggestedFixMustUseRegisterHandler() analysis.SuggestedFix {
	var beginPos token.Pos
	switch {
	case hl.ID != nil:
		beginPos = hl.ID.Pos
	case hl.IDFunc != nil:
		beginPos = hl.ID.Pos
	default:
		beginPos = hl.Pos - token.Pos(1)
	}

	regHandlerStr := fmt.Sprintf(`Reg: parser.RegisterHandler("%q"),`, hl.Name)
	return newSuggestedFixInsert(regHandlerStr, beginPos)
}

func (hfReg *handlerFieldReg) suggestedFixReplaceName(handlerName string) analysis.SuggestedFix {
	return newSuggestedFixReplace(hfReg.Name, handlerName, hfReg.NamePos)
}

// Validate returns any issues for handlerLint
func (hl *HandlerLint) Validate() []*HandlerLintIssue {
	var issues []*HandlerLintIssue
	if hl.Reg == nil {
		issue := hl.newIssue("must use RegisterHandler")
		issue.addSuggestedFix(hl.suggestedFixMustUseRegisterHandler())
		issues = append(issues, issue)
	}

	if hl.Reg != nil && hl.Reg.Name != hl.Name {
		issue := hl.newIssuef("registered name '%s' must match var name", hl.Reg.Name)
		issue.addSuggestedFix(hl.Reg.suggestedFixReplaceName(hl.Name))
		issues = append(issues, issue)
	}

	if hl.IDFunc != nil {
		issues = append(issues, hl.IDFunc.newHandlerLintIssuef("handler '%s': IDFunc is deprecated", hl.Name))
	}

	return issues
}

// HandlerLintIssue contains reportable information for a handler lint issue
type HandlerLintIssue struct {
	Pos            token.Pos
	Message        string
	SuggestedFixes []analysis.SuggestedFix
}

func (issue *HandlerLintIssue) addSuggestedFix(fix analysis.SuggestedFix) {
	issue.SuggestedFixes = append(issue.SuggestedFixes, fix)
}

func (hl *HandlerLint) newIssuef(format string, args ...interface{}) *HandlerLintIssue {
	msg := fmt.Sprintf(format, args...)
	return &HandlerLintIssue{
		Pos:     hl.Pos,
		Message: fmt.Sprintf("handler '%s': %s", hl.Name, msg),
	}
}

func (hl *HandlerLint) newIssue(msg string) *HandlerLintIssue {
	return &HandlerLintIssue{
		Pos:     hl.Pos,
		Message: fmt.Sprintf("handler '%s': %s", hl.Name, msg),
	}
}

func (hf *handlerField) newHandlerLintIssuef(format string, args ...interface{}) *HandlerLintIssue {
	return &HandlerLintIssue{
		Pos:     hf.Pos,
		Message: fmt.Sprintf(format, args...),
	}
}

// Diagnose returns an analysis.Diagnostic
func (issue *HandlerLintIssue) Diagnose() analysis.Diagnostic {
	return analysis.Diagnostic{
		Pos:            issue.Pos,
		Message:        issue.Message,
		SuggestedFixes: issue.SuggestedFixes,
	}
}

// ResultIssue returns an analysis.Diagnostic
func (issue *HandlerLintIssue) ResultIssue() *result.Issue {
	rsltIssue := result.Issue{
		FromLinter: "extractorlint",
		Text:       issue.Message,
		Severity:   "warning",
		Pos:        token.Position{},
		// ExpectNoLint:         true,
		// ExpectedNoLintLinter: "extractorlint",
	}

	if len(issue.SuggestedFixes) > 0 {
		suggFix := issue.SuggestedFixes[0]
		var textEdit analysis.TextEdit
		var newLines []string
		if len(suggFix.TextEdits) > 0 {
			textEdit = suggFix.TextEdits[0]
			newLines = []string{string(textEdit.NewText)}
		}
		rsltIssue.Replacement = &result.Replacement{
			NeedOnlyDelete: false,
			NewLines:       newLines,
			Inline: &result.InlineFix{
				StartCol:  int(textEdit.Pos),
				Length:    int(textEdit.End - textEdit.Pos),
				NewString: string(textEdit.NewText),
			},
		}
	}

	return &rsltIssue
}

func newSuggestedFixInsert(insertExpr string, beginPos token.Pos) analysis.SuggestedFix {
	return analysis.SuggestedFix{
		Message: fmt.Sprintf("insert '%s'", insertExpr),
		TextEdits: []analysis.TextEdit{
			{
				Pos:     beginPos,
				End:     beginPos,
				NewText: []byte(insertExpr + "\n\t"),
			},
		},
	}
}

func newSuggestedFixReplace(oldExpr, newExpr string, beginPos token.Pos) analysis.SuggestedFix {
	endPos := token.Pos(int(beginPos) + len(oldExpr))
	return analysis.SuggestedFix{
		Message: fmt.Sprintf("replace '%s' with '%s'", oldExpr, newExpr),
		TextEdits: []analysis.TextEdit{
			{
				Pos:     beginPos,
				End:     endPos,
				NewText: []byte(newExpr),
			},
		},
	}
}
