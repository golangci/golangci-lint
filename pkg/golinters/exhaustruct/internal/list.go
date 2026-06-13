package internal

import (
	"regexp"

	"dev.gaijin.team/go/golib/e"
	"dev.gaijin.team/go/golib/fields"
)

// ---
// Altered copy of https://github.com/GaijinEntertainment/go-exhaustruct/blob/v5.0.0/internal/pattern/list.go
// ---

func NewList(patterns ...string) ([]*regexp.Regexp, error) {
	if len(patterns) == 0 {
		return nil, nil
	}

	list := make([]*regexp.Regexp, 0, len(patterns))

	for _, pattern := range patterns {
		re, err := compilePattern(pattern)
		if err != nil {
			return nil, err
		}

		list = append(list, re)
	}

	return list, nil
}

func compilePattern(pattern string) (*regexp.Regexp, error) {
	if pattern == "" {
		return nil, e.New("empty regular expression is not allowed")
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, e.NewFrom("failed to compile regular expression", err, fields.F("pattern", pattern))
	}

	return re, nil
}
