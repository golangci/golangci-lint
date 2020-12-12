package commands

import (
	"fmt"

	"github.com/anduril/golangci-lint/pkg/config"
)

var andurilDisableLinters = []string{
	"dogsled",
	"dupl",
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
	"goimports",
	"gofumpt",
	"gomnd",
	"interfacer",
	"lll",
	"maligned",
	"nakedret",
	"nestif",
	"nlreturn",
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
