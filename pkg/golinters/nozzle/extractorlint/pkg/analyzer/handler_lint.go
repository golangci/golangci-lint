package analyzer

import (
	"fmt"
	"go/token"

	"golang.org/x/tools/go/analysis"
)

type handlerLint struct {
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

func (hl *handlerLint) suggestedFixMustUseRegisterHandler() analysis.SuggestedFix {
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

func (hl *handlerLint) validate() []*handlerLintIssue {
	var issues []*handlerLintIssue
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

type handlerLintIssue struct {
	Pos            token.Pos
	Message        string
	SuggestedFixes []analysis.SuggestedFix
}

func (issue *handlerLintIssue) addSuggestedFix(fix analysis.SuggestedFix) {
	issue.SuggestedFixes = append(issue.SuggestedFixes, fix)
}

func (hl *handlerLint) newIssuef(format string, args ...interface{}) *handlerLintIssue {
	msg := fmt.Sprintf(format, args...)
	return &handlerLintIssue{
		Pos:     hl.Pos,
		Message: fmt.Sprintf("handler '%s': %s", hl.Name, msg),
	}
}

func (hl *handlerLint) newIssue(msg string) *handlerLintIssue {
	return &handlerLintIssue{
		Pos:     hl.Pos,
		Message: fmt.Sprintf("handler '%s': %s", hl.Name, msg),
	}
}

func (hf *handlerField) newHandlerLintIssuef(format string, args ...interface{}) *handlerLintIssue {
	return &handlerLintIssue{
		Pos:     hf.Pos,
		Message: fmt.Sprintf(format, args...),
	}
}

// Diagnose returns an analysis.Diagnostic
func (issue *handlerLintIssue) Diagnose() analysis.Diagnostic {
	return analysis.Diagnostic{
		Pos:            issue.Pos,
		Message:        issue.Message,
		SuggestedFixes: issue.SuggestedFixes,
	}
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
