package two

type Formats struct {
	Text        Text         `yaml:"text,omitempty" toml:"text,omitempty"`
	JSON        SimpleFormat `yaml:"json,omitempty" toml:"json,omitempty"`
	Tab         Tab          `yaml:"tab,omitempty" toml:"tab,omitempty"`
	HTML        SimpleFormat `yaml:"html,omitempty" toml:"html,omitempty"`
	Checkstyle  SimpleFormat `yaml:"checkstyle,omitempty" toml:"checkstyle,omitempty"`
	CodeClimate SimpleFormat `yaml:"code-climate,omitempty" toml:"code-climate,omitempty"`
	JUnitXML    JUnitXML     `yaml:"junit-xml,omitempty" toml:"junit-xml,omitempty"`
	TeamCity    SimpleFormat `yaml:"teamcity,omitempty" toml:"teamcity,omitempty"`
	Sarif       SimpleFormat `yaml:"sarif,omitempty" toml:"sarif,omitempty"`
}

type SimpleFormat struct {
	Path *string `yaml:"path,omitempty" toml:"path,omitempty"`
}

type Text struct {
	SimpleFormat    `yaml:",inline"`
	PrintLinterName *bool `yaml:"print-linter-name,omitempty" toml:"print-linter-name,omitempty"`
	PrintIssuedLine *bool `yaml:"print-issued-lines,omitempty" toml:"print-issued-lines,omitempty"`
	Colors          *bool `yaml:"colors,omitempty" toml:"colors,omitempty"`
}

type Tab struct {
	SimpleFormat    `yaml:",inline"`
	PrintLinterName *bool `yaml:"print-linter-name,omitempty" toml:"print-linter-name,omitempty"`
	Colors          *bool `yaml:"colors,omitempty" toml:"colors,omitempty"`
}

type JUnitXML struct {
	SimpleFormat `yaml:",inline"`
	Extended     *bool `yaml:"extended,omitempty" toml:"extended,omitempty"`
}
