package gci

import (
	gcicfg "github.com/daixiang0/gci/pkg/config"
	"github.com/daixiang0/gci/pkg/gci"
	"github.com/ldez/grignotin/gomod"

	"github.com/golangci/golangci-lint/pkg/config"
)

const Name = "gci"

type Formatter struct {
	config *gcicfg.Config
}

func New(cfg config.GciSettings) (*Formatter, error) {
	modPath, err := gomod.GetModulePath()
	if err != nil {
		return nil, err
	}

	parsedCfg, err := gcicfg.YamlConfig{
		Cfg: gcicfg.BoolConfig{
			NoInlineComments: cfg.NoInlineComments,
			NoPrefixComments: cfg.NoPrefixComments,
			SkipGenerated:    cfg.SkipGenerated,
			CustomOrder:      cfg.CustomOrder,
			NoLexOrder:       cfg.NoLexOrder,
		},
		SectionStrings: cfg.Sections,
		ModPath:        modPath,
	}.Parse()
	if err != nil {
		return nil, err
	}

	return &Formatter{config: parsedCfg}, nil
}

func (*Formatter) Name() string {
	return Name
}

func (f *Formatter) Format(filename string, src []byte) ([]byte, error) {
	_, formatted, err := gci.LoadFormat(src, filename, *f.config)
	return formatted, err
}
