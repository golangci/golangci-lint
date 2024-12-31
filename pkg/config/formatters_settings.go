package config

var defaultFormatterSettings = FormatterSettings{
	GoFmt: GoFmtSettings{
		Simplify: true,
	},
	Gci: GciSettings{
		Sections:      []string{"standard", "default"},
		SkipGenerated: true,
	},
	GoLines: GoLinesSettings{
		MaxLen:         100,
		TabLen:         4,
		ReformatTags:   true,
		ChainSplitDots: true,
	},
}

type FormatterSettings struct {
	Gci       GciSettings       `mapstructure:"gci"`
	GoFmt     GoFmtSettings     `mapstructure:"gofmt"`
	GoFumpt   GoFumptSettings   `mapstructure:"gofumpt"`
	GoImports GoImportsSettings `mapstructure:"goimports"`
	GoLines   GoLinesSettings   `mapstructure:"golines"`
}

type GciSettings struct {
	Sections         []string `mapstructure:"sections"`
	NoInlineComments bool     `mapstructure:"no-inline-comments"`
	NoPrefixComments bool     `mapstructure:"no-prefix-comments"`
	SkipGenerated    bool     `mapstructure:"skip-generated"`
	CustomOrder      bool     `mapstructure:"custom-order"`
	NoLexOrder       bool     `mapstructure:"no-lex-order"`

	// Deprecated: use Sections instead.
	LocalPrefixes string `mapstructure:"local-prefixes"`
}

type GoFmtSettings struct {
	Simplify     bool
	RewriteRules []GoFmtRewriteRule `mapstructure:"rewrite-rules"`
}

type GoFmtRewriteRule struct {
	Pattern     string
	Replacement string
}

type GoFumptSettings struct {
	ModulePath string `mapstructure:"module-path"`
	ExtraRules bool   `mapstructure:"extra-rules"`

	// Deprecated: use the global `run.go` instead.
	LangVersion string `mapstructure:"lang-version"`
}

type GoImportsSettings struct {
	LocalPrefixes string `mapstructure:"local-prefixes"`
}

type GoLinesSettings struct {
	MaxLen          int  `mapstructure:"max-len"`
	TabLen          int  `mapstructure:"tab-len"`
	ShortenComments bool `mapstructure:"shorten-comments"`
	ReformatTags    bool `mapstructure:"reformat-tags"`
	ChainSplitDots  bool `mapstructure:"chain-split-dots"`
}
