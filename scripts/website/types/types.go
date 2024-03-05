package types

import "github.com/golangci/golangci-lint/pkg/lint/linter"

type CLIHelp struct {
	Enable  string `json:"enable"`
	Disable string `json:"disable"`
	Help    string `json:"help"`
}

type LinterWrapper struct {
	*linter.Config

	Name string `json:"name"`
	Desc string `json:"desc"`
}
