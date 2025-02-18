package config

type Formats struct {
	Text        Text        `mapstructure:"text"`
	JSON        SimpleStyle `mapstructure:"json"`
	Tab         Tab         `mapstructure:"tab"`
	HTML        SimpleStyle `mapstructure:"html"`
	Checkstyle  SimpleStyle `mapstructure:"checkstyle"`
	CodeClimate SimpleStyle `mapstructure:"code-climate"`
	JUnitXML    JUnitXML    `mapstructure:"junit-xml"`
	TeamCity    SimpleStyle `mapstructure:"team-city"`
	Sarif       SimpleStyle `mapstructure:"sarif"`
}

func (f *Formats) IsEmpty() bool {
	styles := []SimpleStyle{
		f.Text.SimpleStyle,
		f.JSON,
		f.Tab.SimpleStyle,
		f.HTML,
		f.Checkstyle,
		f.CodeClimate,
		f.JUnitXML.SimpleStyle,
		f.TeamCity,
		f.Sarif,
	}

	for _, style := range styles {
		if style.Path != "" {
			return false
		}
	}

	return true
}

type SimpleStyle struct {
	Path string `mapstructure:"path"`
}

type Text struct {
	SimpleStyle     `mapstructure:",squash"`
	PrintLinterName bool `mapstructure:"print-linter-name"`
	PrintIssuedLine bool `mapstructure:"print-issued-lines"`
	Colors          bool `mapstructure:"colors"`
}

type Tab struct {
	SimpleStyle     `mapstructure:",squash"`
	PrintLinterName bool `mapstructure:"print-linter-name"`
	UseColors       bool `mapstructure:"use-colors"`
}

type JUnitXML struct {
	SimpleStyle `mapstructure:",squash"`
	Extended    bool `mapstructure:"extended"`
}
