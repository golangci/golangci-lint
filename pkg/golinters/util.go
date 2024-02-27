package golinters

import (
	"fmt"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
)

func formatCode(code string, _ *config.Config) string {
	if strings.Contains(code, "`") {
		return code // TODO: properly escape or remove
	}

	return fmt.Sprintf("`%s`", code)
}

func getFileNames(pass *analysis.Pass) []string {
	var fileNames []string
	for _, f := range pass.Files {
		fileName := pass.Fset.PositionFor(f.Pos(), true).Filename
		ext := filepath.Ext(fileName)
		if ext != "" && ext != ".go" {
			// position has been adjusted to a non-go file, revert to original file
			fileName = pass.Fset.PositionFor(f.Pos(), false).Filename
		}
		fileNames = append(fileNames, fileName)
	}
	return fileNames
}
