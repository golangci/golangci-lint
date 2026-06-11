package internal

import (
	"fmt"
	"strings"

	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
)

func FormatCode(code string) string {
	if strings.Contains(code, "`") {
		return fmt.Sprintf("%q", code)
	}

	return fmt.Sprintf("`%s`", code)
}

func GetGoFileNames(pass *analysis.Pass) []string {
	var filenames []string

	for _, f := range pass.Files {
		position, b := goanalysis.GetGoFilePosition(pass, f)
		if !b {
			continue
		}

		filenames = append(filenames, position.Filename)
	}

	return filenames
}
