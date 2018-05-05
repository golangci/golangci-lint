package executors

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/golangci/golangci-shared/pkg/analytics"
)

type TempDirShell struct {
	Shell
}

var _ Executor = &TempDirShell{}

var tmpRoot string

func init() {
	var err error
	tmpRoot, err = filepath.EvalSymlinks("/tmp")
	if err != nil {
		log.Fatalf("can't eval symlinks on /tmp: %s", err)
	}
}

func NewTempDirShell(tag string) (*TempDirShell, error) {
	wd, err := ioutil.TempDir(tmpRoot, fmt.Sprintf("golangci.%s", tag))
	if err != nil {
		return nil, fmt.Errorf("can't make temp dir: %s", err)
	}

	return &TempDirShell{
		Shell: *NewShell(wd),
	}, nil
}

func (s TempDirShell) Clean() {
	if err := os.RemoveAll(s.wd); err != nil {
		analytics.Log(context.TODO()).Warnf("Can't remove temp dir %s: %s", s.wd, err)
	}
}

func (s TempDirShell) WithEnv(k, v string) Executor {
	eCopy := s
	eCopy.SetEnv(k, v)
	return &eCopy
}

func (s TempDirShell) WithWorkDir(wd string) Executor {
	eCopy := s
	eCopy.wd = wd
	return &eCopy
}
