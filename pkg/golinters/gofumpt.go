package golinters

import (
	"bytes"
	"fmt"
	"go/token"
	"io"
	"io/ioutil"
	"os"
	"sync"

	"golang.org/x/tools/go/analysis"
	"mvdan.cc/gofumpt/format"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

const gofumptName = "gofumpt"

func NewGofumpt() *goanalysis.Linter {
	var mu sync.Mutex
	var resIssues []goanalysis.Issue

	analyzer := &analysis.Analyzer{
		Name: gofumptName,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
	}
	return goanalysis.NewLinter(
		gofumptName,
		"Gofumpt checks whether code was gofumpt-ed.",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithContextSetter(func(lintCtx *linter.Context) {
		analyzer.Run = func(pass *analysis.Pass) (interface{}, error) {
			var fileNames []string
			for _, f := range pass.Files {
				pos := pass.Fset.PositionFor(f.Pos(), false)
				fileNames = append(fileNames, pos.Filename)
			}

			var issues []goanalysis.Issue

			for _, f := range fileNames {
				input, err := ioutil.ReadFile(f)
				if err != nil {
					return nil, fmt.Errorf("unable to open file %s: %w", f, err)
				}
				output, err := format.Source(input, "")
				if err != nil {
					return nil, fmt.Errorf("error while running gofumpt: %w", err)
				}
				if !bytes.Equal(input, output) {
					reader, err := os.Open(f)
					if err != nil {
						return nil, fmt.Errorf("unable to open file %s: %w", f, err)
					}
					numberOfLines, err := lineCounter(reader)
					if err != nil {
						return nil, fmt.Errorf("unable to count the number of lines: %w", err)
					}

					issues = append(issues, goanalysis.NewIssue(&result.Issue{
						FromLinter: gofumptName,
						Text:       "File is not `gofumpt`-ed",
						Pos: token.Position{
							Filename: f,
							Line:     numberOfLines,
						},
					}, pass))
				}
			}

			if len(issues) == 0 {
				return nil, nil
			}

			mu.Lock()
			resIssues = append(resIssues, issues...)
			mu.Unlock()

			return nil, nil
		}
	}).WithIssuesReporter(func(*linter.Context) []goanalysis.Issue {
		return resIssues
	}).WithLoadMode(goanalysis.LoadModeSyntax)
}

func lineCounter(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}
