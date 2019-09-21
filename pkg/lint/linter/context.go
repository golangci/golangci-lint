package linter

import (
	"golang.org/x/tools/go/loader"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis/load"

	"github.com/golangci/golangci-lint/internal/pkgcache"

	"github.com/golangci/golangci-lint/pkg/fsutils"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/lint/astcache"
	"github.com/golangci/golangci-lint/pkg/logutils"
)

type Context struct {
	// Packages are deduplicated (test and normal packages) packages
	Packages []*packages.Package

	// OriginalPackages aren't deduplicated: they contain both normal and test
	// version for each of packages
	OriginalPackages []*packages.Package

	NotCompilingPackages []*packages.Package

	LoaderConfig *loader.Config  // deprecated, don't use for new linters
	Program      *loader.Program // deprecated, use Packages for new linters

	SSAProgram *ssa.Program // for unparam and interfacer but not for megacheck (it change it)

	Cfg       *config.Config
	ASTCache  *astcache.Cache
	FileCache *fsutils.FileCache
	LineCache *fsutils.LineCache
	Log       logutils.Log

	PkgCache         *pkgcache.Cache
	LoadGuard        *load.Guard
	NeedWholeProgram bool
}

func (c *Context) Settings() *config.LintersSettings {
	return &c.Cfg.LintersSettings
}
