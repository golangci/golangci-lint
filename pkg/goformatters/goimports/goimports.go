package goimports

import (
	"golang.org/x/tools/imports"
)

const Name = "goimports"

type Formatter struct{}

func New() *Formatter {
	return &Formatter{}
}

func (*Formatter) Name() string {
	return Name
}

func (*Formatter) Format(filename string, src []byte) ([]byte, error) {
	// The `imports.LocalPrefix` (`settings.LocalPrefixes`) is a global var.
	return imports.Process(filename, src, nil)
}
