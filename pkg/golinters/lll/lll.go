package lll

import (
	"bufio"
	"errors"
	"fmt"
	"go/ast"
	"os"
	"strings"
	"unicode/utf8"

	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
)

const linterName = "lll"

const goCommentDirectivePrefix = "//go:"

func New(settings *config.LllSettings) *goanalysis.Linter {
	analyzer := &analysis.Analyzer{
		Name: linterName,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
		Run: func(pass *analysis.Pass) (any, error) {
			err := runLll(pass, settings)
			if err != nil {
				return nil, err
			}

			return nil, nil
		},
	}

	return goanalysis.NewLinter(
		linterName,
		"Reports long lines",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}

func runLll(pass *analysis.Pass, settings *config.LllSettings) error {
	spaces := strings.Repeat(" ", settings.TabWidth)

	for _, file := range pass.Files {
		err := getLLLIssuesForFile(pass, file, settings.LineLength, spaces)
		if err != nil {
			return err
		}
	}

	return nil
}

func getLLLIssuesForFile(pass *analysis.Pass, file *ast.File, maxLineLen int, tabSpaces string) error {
	position := goanalysis.GetFilePosition(pass, file)

	if !strings.HasSuffix(position.Filename, ".go") {
		return nil
	}

	nonAdjPosition := pass.Fset.PositionFor(file.Pos(), false)

	f, err := os.Open(position.Filename)
	if err != nil {
		return fmt.Errorf("can't open file %s: %w", position.Filename, err)
	}

	defer f.Close()

	ft := pass.Fset.File(file.Pos())

	lineNumber := 0
	multiImportEnabled := false

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lineNumber++

		line := scanner.Text()
		line = strings.ReplaceAll(line, "\t", tabSpaces)

		if strings.HasPrefix(line, goCommentDirectivePrefix) {
			continue
		}

		if strings.HasPrefix(line, "import") {
			multiImportEnabled = strings.HasSuffix(line, "(")
			continue
		}

		if multiImportEnabled {
			if line == ")" {
				multiImportEnabled = false
			}

			continue
		}

		lineLen := utf8.RuneCountInString(line)
		if lineLen > maxLineLen {
			pass.Report(analysis.Diagnostic{
				Pos: ft.LineStart(goanalysis.AdjustPos(lineNumber, nonAdjPosition.Line, position.Line)),
				Message: fmt.Sprintf("The line is %d characters long, which exceeds the maximum of %d characters.",
					lineLen, maxLineLen),
			})
		}
	}

	if err := scanner.Err(); err != nil {
		// scanner.Scan() might fail if the line is longer than bufio.MaxScanTokenSize
		// In the case where the specified maxLineLen is smaller than bufio.MaxScanTokenSize
		// we can return this line as a long line instead of returning an error.
		// The reason for this change is that this case might happen with autogenerated files
		// The go-bindata tool for instance might generate a file with a very long line.
		// In this case, as it's an auto generated file, the warning returned by lll will
		// be ignored.
		// But if we return a linter error here, and this error happens for an autogenerated
		// file the error will be discarded (fine), but all the subsequent errors for lll will
		// be discarded for other files, and we'll miss legit error.
		if errors.Is(err, bufio.ErrTooLong) && maxLineLen < bufio.MaxScanTokenSize {
			pass.Report(analysis.Diagnostic{
				Pos:     ft.LineStart(goanalysis.AdjustPos(lineNumber, nonAdjPosition.Line, position.Line)),
				Message: fmt.Sprintf("line is more than %d characters", bufio.MaxScanTokenSize),
			})
		} else {
			return fmt.Errorf("can't scan file %s: %w", position.Filename, err)
		}
	}

	return nil
}
