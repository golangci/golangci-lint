package golinters

import (
	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/fsutils"
	"golang.org/x/tools/go/loader"
)

type Context struct {
	Paths        *fsutils.ProjectPaths
	Cfg          *config.Config
	Program      *loader.Program
	LoaderConfig *loader.Config
}

func (c *Context) RunCfg() *config.Run {
	return &c.Cfg.Run
}
