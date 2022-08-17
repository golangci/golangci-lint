package commands

import (
	"fmt"

	"github.com/anduril/golangci-lint/pkg/config"
)

var andurilDisableLinters = []string{
	"dogsled",
	"dupl",
	"exhaustruct",
	"exhaustivestruct",
	"funlen",
	"gci",
	"gochecknoglobals",
	"gocognit",
	"goconst",
	"gocyclo",
	"godot",
	"godox",
	"goerr113",
	"goerr113",
	"goimports",
	"golint",
	"gofumpt",
	"gomnd",
	"interfacer",
	"ireturn",
	"lll",
	"maligned",
	"nakedret",
	"nestif",
	"nlreturn",
	"nosnakecase",
	"paralleltest",
	"prealloc",
	"scopelint",
	"stylecheck",
	"testpackage",
	"wrapcheck",
	"wsl",
}

func modifyConfigAnduril(c *config.Config) error {
	if !c.Linters.Anduril {
		return nil
	}

	if c.Linters.DisableAll {
		return fmt.Errorf("cannot set linters.disable-all when 'linters.anduril' == true")
	}
	c.Linters.EnableAll = true
	c.Linters.Disable = append(c.Linters.Disable, andurilDisableLinters...)
	return nil
}
