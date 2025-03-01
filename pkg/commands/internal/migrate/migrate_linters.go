package migrate

import (
	"github.com/golangci/golangci-lint/pkg/commands/internal/migrate/one"
	"github.com/golangci/golangci-lint/pkg/commands/internal/migrate/ptr"
	"github.com/golangci/golangci-lint/pkg/commands/internal/migrate/two"
)

func toLinters(old *one.Config) two.Linters {
	enable, disable := ProcessEffectiveLinters(old.Linters)

	return two.Linters{
		Default:    getDefaultName(old.Linters),
		Enable:     onlyLinterNames(convertStaticcheckLinterNames(enable)),
		Disable:    onlyLinterNames(convertStaticcheckLinterNames(disable)),
		FastOnly:   nil,
		Settings:   toLinterSettings(old.LintersSettings),
		Exclusions: toExclusions(old),
	}
}

func getDefaultName(old one.Linters) *string {
	switch {
	case ptr.Deref(old.DisableAll):
		return ptr.Pointer("none")
	case ptr.Deref(old.EnableAll):
		return ptr.Pointer("all")
	default:
		return nil // standard is the default
	}
}
