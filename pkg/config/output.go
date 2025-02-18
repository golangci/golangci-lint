package config

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

type Output struct {
	Formats     Formats  `mapstructure:"formats"`
	SortResults bool     `mapstructure:"sort-results"`
	SortOrder   []string `mapstructure:"sort-order"`
	PathPrefix  string   `mapstructure:"path-prefix"`
	ShowStats   bool     `mapstructure:"show-stats"`

	// Deprecated: use [Issues.UniqByLine] instead.
	UniqByLine *bool `mapstructure:"uniq-by-line"`
}

func (o *Output) Validate() error {
	if !o.SortResults && len(o.SortOrder) > 0 {
		return errors.New("sort-results should be 'true' to use sort-order")
	}

	validOrders := []string{"linter", "file", "severity"}

	all := strings.Join(o.SortOrder, " ")

	for _, order := range o.SortOrder {
		if strings.Count(all, order) > 1 {
			return fmt.Errorf("the sort-order name %q is repeated several times", order)
		}

		if !slices.Contains(validOrders, order) {
			return fmt.Errorf("unsupported sort-order name %q", order)
		}
	}

	return nil
}
