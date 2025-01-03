package internal

import (
	"fmt"
	"strings"

	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
)

func FormatCode(code string, _ *config.Config) string {
	if strings.Contains(code, "`") {
		return code // TODO: properly escape or remove
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
