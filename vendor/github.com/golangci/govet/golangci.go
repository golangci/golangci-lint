package govet

import (
	"go/token"
	"strings"
)

type Issue struct {
	Pos     token.Position
	Message string
}

var foundIssues []Issue

func Run(files []string, checkShadowing bool) ([]Issue, error) {
	foundIssues = nil

	if checkShadowing {
		experimental["shadow"] = false
	}
	for name, setting := range report {
		if *setting == unset && !experimental[name] {
			*setting = setTrue
		}
	}

	initPrintFlags()
	initUnusedFlags()

	filesRun = true
	for _, name := range files {
		if !strings.HasSuffix(name, "_test.go") {
			includesNonTest = true
		}
	}
	if doPackage(files, nil) == nil {
		return nil, nil
	}

	return foundIssues, nil
}
