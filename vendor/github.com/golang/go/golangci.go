package govet

import (
	"fmt"
	"go/token"
	"os"
	"strings"
)

type Issue struct {
	Pos     token.Position
	Message string
}

var foundIssues []Issue

func Run(paths, buildTags []string, checkShadowing bool) ([]Issue, error) {
	foundIssues = nil

	if checkShadowing {
		experimental["shadow"] = false
	}
	for name, setting := range report {
		if *setting == unset && !experimental[name] {
			*setting = setTrue
		}
	}

	tagList = buildTags

	initPrintFlags()
	initUnusedFlags()

	for _, name := range paths {
		// Is it a directory?
		fi, err := os.Stat(name)
		if err != nil {
			warnf("error walking tree: %s", err)
			continue
		}
		if fi.IsDir() {
			dirsRun = true
		} else {
			filesRun = true
			if !strings.HasSuffix(name, "_test.go") {
				includesNonTest = true
			}
		}
	}
	if dirsRun && filesRun {
		return nil, fmt.Errorf("can't mix dirs and files")
	}
	if dirsRun {
		for _, name := range paths {
			walkDir(name)
		}
		return foundIssues, nil
	}
	if doPackage(paths, nil) == nil {
		return nil, nil
	}

	return foundIssues, nil
}
