package packages

import (
	"fmt"

	"github.com/golangci/golangci-lint/pkg/lint/astcache"

	"golang.org/x/tools/go/packages"
)

//nolint:gocyclo
func ExtractErrors(pkg *packages.Package, astCache *astcache.Cache) []packages.Error {
	errors := extractErrorsImpl(pkg)
	if len(errors) == 0 {
		return errors
	}

	seenErrors := map[string]bool{}
	var uniqErrors []packages.Error
	for _, err := range errors {
		if seenErrors[err.Msg] {
			continue
		}
		seenErrors[err.Msg] = true
		uniqErrors = append(uniqErrors, err)
	}

	if len(pkg.GoFiles) != 0 {
		// errors were extracted from deps and have at leat one file in package
		for i := range uniqErrors {
			errPos, parseErr := ParseErrorPosition(uniqErrors[i].Pos)
			if parseErr != nil || astCache.Get(errPos.Filename) == nil {
				// change pos to local file to properly process it by processors (properly read line etc)
				uniqErrors[i].Msg = fmt.Sprintf("%s: %s", uniqErrors[i].Pos, uniqErrors[i].Msg)
				uniqErrors[i].Pos = fmt.Sprintf("%s:1", pkg.GoFiles[0])
			}
		}

		// some errors like "code in directory  expects import" don't have Pos, set it here
		for i := range uniqErrors {
			err := &uniqErrors[i]
			if err.Pos == "" {
				err.Pos = fmt.Sprintf("%s:1", pkg.GoFiles[0])
			}
		}
	}

	return uniqErrors
}

func extractErrorsImpl(pkg *packages.Package) []packages.Error {
	var errors []packages.Error
	packages.Visit([]*packages.Package{pkg},
		// only keep if pkg or its parents have errors,
		// otherwise save the effort, there are no errors to be found.
		func(pkg *packages.Package) bool { return pkg.IllTyped },
		func(pkg *packages.Package) { errors = append(errors, pkg.Errors...) })
	return errors
}
