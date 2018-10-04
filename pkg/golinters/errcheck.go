package golinters

import (
	"bufio"
	"context"
	"fmt"
	"os"

	errcheckAPI "github.com/golangci/errcheck/golangci"
	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
	"github.com/pkg/errors"
)

type Errcheck struct{}

func (Errcheck) Name() string {
	return "errcheck"
}

func (Errcheck) Desc() string {
	return "Errcheck is a program for checking for unchecked errors " +
		"in go programs. These unchecked errors can be critical bugs in some cases"
}

func (e Errcheck) Run(ctx context.Context, lintCtx *linter.Context) ([]result.Issue, error) {
	errCfg, err := genConfig(&lintCtx.Settings().Errcheck)
	if err != nil {
		return nil, err
	}
	issues, err := errcheckAPI.RunWithConfig(lintCtx.Program, errCfg)
	if err != nil {
		return nil, err
	}

	if len(issues) == 0 {
		return nil, nil
	}

	res := make([]result.Issue, 0, len(issues))
	for _, i := range issues {
		var text string
		if i.FuncName != "" {
			text = fmt.Sprintf("Error return value of %s is not checked", formatCode(i.FuncName, lintCtx.Cfg))
		} else {
			text = "Error return value is not checked"
		}
		res = append(res, result.Issue{
			FromLinter: e.Name(),
			Text:       text,
			Pos:        i.Pos,
		})
	}

	return res, nil
}

func genConfig(errCfg *config.ErrcheckSettings) (*errcheckAPI.Config, error) {
	c := &errcheckAPI.Config{
		Ignore:  errCfg.Ignore,
		Blank:   errCfg.CheckAssignToBlank,
		Asserts: errCfg.CheckTypeAssertions,
	}
	if errCfg.Exclude != "" {
		exclude, err := readExcludeFile(errCfg.Exclude)
		if err != nil {
			return nil, err
		}
		c.Exclude = exclude
	}
	return c, nil
}

func readExcludeFile(name string) (map[string]bool, error) {
	exclude := make(map[string]bool)
	fh, err := os.Open(name)
	if err != nil {
		return nil, errors.Wrapf(err, "failed reading exclude file: %s", name)
	}
	scanner := bufio.NewScanner(fh)
	for scanner.Scan() {
		exclude[scanner.Text()] = true
	}
	if err := scanner.Err(); err != nil {
		return nil, errors.Wrapf(err, "failed scanning file: %s", name)
	}
	return exclude, nil
}
