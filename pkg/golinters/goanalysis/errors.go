package goanalysis

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/tools/go/packages"

	"github.com/golangci/golangci-lint/pkg/lint/linter"
	libpackages "github.com/golangci/golangci-lint/pkg/packages"
	"github.com/golangci/golangci-lint/pkg/result"
)

type IllTypedError struct {
	Pkg *packages.Package
}

func (e *IllTypedError) Error() string {
	return fmt.Sprintf("errors in package: %v", e.Pkg.Errors)
}

type FailedPrerequisitesError struct {
	errors map[string][]string
}

func (f FailedPrerequisitesError) NotEmpty() bool {
	return len(f.errors) > 0
}

func (f *FailedPrerequisitesError) Consume(name string, err error) {
	if f.errors == nil {
		f.errors = map[string][]string{}
	}
	k := fmt.Sprintf("%v", err)
	f.errors[k] = append(f.errors[k], name)
}

func (f FailedPrerequisitesError) Error() string {
	var errs []string
	for err := range f.errors {
		errs = append(errs, err)
	}
	var groups []groupedPrerequisiteErr
	for _, err := range errs {
		groups = append(groups, groupedPrerequisiteErr{
			err:   err,
			names: f.errors[err],
		})
	}
	return fmt.Sprintf("failed prerequisites: %s", groups)
}

type groupedPrerequisiteErr struct {
	names []string
	err   string
}

func (g groupedPrerequisiteErr) String() string {
	if len(g.names) == 1 {
		return fmt.Sprintf("%s: %s", g.names[0], g.err)
	}
	return fmt.Sprintf("(%s): %s", strings.Join(g.names, ", "), g.err)
}

func buildIssuesFromErrorsForTypecheckMode(errs []error, lintCtx *linter.Context) ([]result.Issue, error) {
	var issues []result.Issue
	uniqReportedIssues := map[string]bool{}
	for _, err := range errs {
		itErr, ok := errors.Cause(err).(*IllTypedError)
		if !ok {
			return nil, err
		}
		for _, err := range libpackages.ExtractErrors(itErr.Pkg) {
			i, perr := parseError(err)
			if perr != nil { // failed to parse
				if uniqReportedIssues[err.Msg] {
					continue
				}
				uniqReportedIssues[err.Msg] = true
				lintCtx.Log.Errorf("typechecking error: %s", err.Msg)
			} else {
				i.Pkg = itErr.Pkg // to save to cache later
				issues = append(issues, *i)
			}
		}
	}
	return issues, nil
}

func parseError(srcErr packages.Error) (*result.Issue, error) {
	pos, err := libpackages.ParseErrorPosition(srcErr.Pos)
	if err != nil {
		return nil, err
	}

	return &result.Issue{
		Pos:        *pos,
		Text:       srcErr.Msg,
		FromLinter: "typecheck",
	}, nil
}
