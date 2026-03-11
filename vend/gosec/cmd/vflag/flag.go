package vflag

import (
	"errors"
	"strings"
)

// ValidatedFlag cli string type
type ValidatedFlag struct {
	Value string
}

func (f *ValidatedFlag) String() string {
	return f.Value
}

// Set will be called for flag that is of validateFlag type
func (f *ValidatedFlag) Set(value string) error {
	if strings.Contains(value, "-") {
		return errors.New("flag value cannot start with -")
	}

	f.Value = value
	return nil
}
