package packages

import (
	"fmt"

	"golang.org/x/tools/go/packages"
)

func ExtractErrors(pkg *packages.Package) []packages.Error {
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

	if len(pkg.Errors) == 0 && len(pkg.GoFiles) != 0 {
		// erorrs were extracted from deps and have at leat one file in package
		for i := range uniqErrors {
			// change pos to local file to properly process it by processors (properly read line etc)
			uniqErrors[i].Msg = fmt.Sprintf("%s: %s", uniqErrors[i].Pos, uniqErrors[i].Msg)
			uniqErrors[i].Pos = fmt.Sprintf("%s:1", pkg.GoFiles[0])
		}
	}

	return uniqErrors
}

func extractErrorsImpl(pkg *packages.Package) []packages.Error {
	if len(pkg.Errors) != 0 {
		return pkg.Errors
	}

	var errors []packages.Error
	for _, iPkg := range pkg.Imports {
		iPkgErrors := extractErrorsImpl(iPkg)
		if iPkgErrors != nil {
			errors = append(errors, iPkgErrors...)
		}
	}

	return errors
}
