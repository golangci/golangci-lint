package gofmt

import (
	"github.com/golangci/gofmt/gofmt"

	"github.com/golangci/golangci-lint/pkg/config"
)

const Name = "gofmt"

type Formatter struct {
	options gofmt.Options
}

func New(cfg config.GoFmtSettings) *Formatter {
	var rewriteRules []gofmt.RewriteRule
	for _, rule := range cfg.RewriteRules {
		rewriteRules = append(rewriteRules, gofmt.RewriteRule(rule))
	}

	return &Formatter{
		options: gofmt.Options{
			NeedSimplify: cfg.Simplify,
			RewriteRules: rewriteRules,
		},
	}
}

func (*Formatter) Name() string {
	return Name
}

func (f *Formatter) Format(filename string, src []byte) ([]byte, error) {
	return gofmt.Source(filename, src, f.options)
}
