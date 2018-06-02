package linter

import (
	"github.com/golangci/go-tools/ssa"
	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/fsutils"
	"github.com/golangci/golangci-lint/pkg/lint/astcache"
	"golang.org/x/tools/go/loader"
)

type Context struct {
	Paths                *fsutils.ProjectPaths
	Cfg                  *config.Config
	Program              *loader.Program
	SSAProgram           *ssa.Program
	LoaderConfig         *loader.Config
	ASTCache             *astcache.Cache
	NotCompilingPackages []*loader.PackageInfo
}

func (c *Context) Settings() *config.LintersSettings {
	return &c.Cfg.LintersSettings
}
