package config

import (
	"fmt"
	"slices"
	"strings"
)

type Output struct {
	Formats    Formats  `mapstructure:"formats"`
	SortOrder  []string `mapstructure:"sort-order"`
	PathPrefix string   `mapstructure:"path-prefix"`
	ShowStats  bool     `mapstructure:"show-stats"`
}

func (o *Output) Validate() error {
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
