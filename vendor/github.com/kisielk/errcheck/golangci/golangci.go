package golangci

import (
	"regexp"

	"github.com/kisielk/errcheck/internal/errcheck"
	"golang.org/x/tools/go/loader"
)

type Issue errcheck.UncheckedError

func Run(program *loader.Program, checkBlank, checkAsserts bool) ([]Issue, error) {
	checker := errcheck.NewChecker()
	checker.Blank = checkBlank
	checker.Asserts = checkAsserts

	checker.Ignore = map[string]*regexp.Regexp{
		"fmt": regexp.MustCompile(".*"),
	}

	if err := checker.CheckProgram(program); err != nil {
		if e, ok := err.(*errcheck.UncheckedErrors); ok {
			return makeIssues(e), nil
		}
		if err == errcheck.ErrNoGoFiles {
			return nil, nil
		}

		return nil, err
	}

	// no issues
	return nil, nil
}

func makeIssues(e *errcheck.UncheckedErrors) []Issue {
	var ret []Issue
	for _, uncheckedError := range e.Errors {
		ret = append(ret, Issue(uncheckedError))
	}

	return ret
}
