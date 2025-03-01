package two

type FormatterSettings struct {
	Gci       GciSettings       `yaml:"gci,omitempty" toml:"gci,omitempty"`
	GoFmt     GoFmtSettings     `yaml:"gofmt,omitempty" toml:"gofmt,omitempty"`
	GoFumpt   GoFumptSettings   `yaml:"gofumpt,omitempty" toml:"gofumpt,omitempty"`
	GoImports GoImportsSettings `yaml:"goimports,omitempty" toml:"goimports,omitempty"`
	GoLines   GoLinesSettings   `yaml:"golines,omitempty" toml:"golines,omitempty"`
}

type GciSettings struct {
	Sections         []string `yaml:"sections,omitempty" toml:"sections,omitempty"`
	NoInlineComments *bool    `yaml:"no-inline-comments,omitempty" toml:"no-inline-comments,omitempty"`
	NoPrefixComments *bool    `yaml:"no-prefix-comments,omitempty" toml:"no-prefix-comments,omitempty"`
	CustomOrder      *bool    `yaml:"custom-order,omitempty" toml:"custom-order,omitempty"`
	NoLexOrder       *bool    `yaml:"no-lex-order,omitempty" toml:"no-lex-order,omitempty"`
}

type GoFmtSettings struct {
	Simplify     *bool              `yaml:"simplify,omitempty" toml:"simplify,omitempty"`
	RewriteRules []GoFmtRewriteRule `yaml:"rewrite-rules,omitempty" toml:"rewrite-rules,omitempty"`
}

type GoFmtRewriteRule struct {
	Pattern     *string `yaml:"pattern,omitempty" toml:"pattern,omitempty"`
	Replacement *string `yaml:"replacement,omitempty" toml:"replacement,omitempty"`
}

type GoFumptSettings struct {
	ModulePath *string `yaml:"module-path,omitempty" toml:"module-path,omitempty"`
	ExtraRules *bool   `yaml:"extra-rules,omitempty" toml:"extra-rules,omitempty"`

	LangVersion *string `yaml:"-,omitempty" toml:"-,omitempty"`
}

type GoImportsSettings struct {
	LocalPrefixes []string `yaml:"local-prefixes,omitempty" toml:"local-prefixes,omitempty"`
}

type GoLinesSettings struct {
	MaxLen          *int  `yaml:"max-len,omitempty" toml:"max-len,omitempty"`
	TabLen          *int  `yaml:"tab-len,omitempty" toml:"tab-len,omitempty"`
	ShortenComments *bool `yaml:"shorten-comments,omitempty" toml:"shorten-comments,omitempty"`
	ReformatTags    *bool `yaml:"reformat-tags,omitempty" toml:"reformat-tags,omitempty"`
	ChainSplitDots  *bool `yaml:"chain-split-dots,omitempty" toml:"chain-split-dots,omitempty"`
}
