package linter

import (
	"github.com/golangci/tools/go/ssa"
	"golang.org/x/tools/go/loader"
	"golang.org/x/tools/go/packages"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/lint/astcache"
	"github.com/golangci/golangci-lint/pkg/logutils"
)

type Context struct {
	Packages             []*packages.Package
	NotCompilingPackages []*packages.Package

	LoaderConfig *loader.Config  // deprecated, don't use for new linters
	Program      *loader.Program // deprecated, use Packages for new linters

	SSAProgram *ssa.Program

	Cfg      *config.Config
	ASTCache *astcache.Cache
	Log      logutils.Log
}

func (c *Context) Settings() *config.LintersSettings {
	return &c.Cfg.LintersSettings
}
