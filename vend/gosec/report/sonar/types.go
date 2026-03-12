package sonar

// TextRange defines the text range of an issue's location
type TextRange struct {
	StartLine   int `json:"startLine"`
	EndLine     int `json:"endLine"`
	StartColumn int `json:"startColumn,omitempty"`
	EtartColumn int `json:"endColumn,omitempty"`
}

// Location defines a sonar issue's location
type Location struct {
	Message   string     `json:"message"`
	FilePath  string     `json:"filePath"`
	TextRange *TextRange `json:"textRange,omitempty"`
}

// Issue defines a sonar issue
type Issue struct {
	RuleID             string      `json:"ruleId"`
	PrimaryLocation    *Location   `json:"primaryLocation"`
	EffortMinutes      int         `json:"effortMinutes"`
	SecondaryLocations []*Location `json:"secondaryLocations,omitempty"`
}

// Impact defines the impact for a rule.
type Impact struct {
	SoftwareQuality string `json:"softwareQuality"`
	Severity        string `json:"severity"`
}

// Rule defines a sonar rule.
type Rule struct {
	ID                 string    `json:"id"`
	Name               string    `json:"name"`
	Description        string    `json:"description"`
	EngineID           string    `json:"engineId"`
	CleanCodeAttribute string    `json:"cleanCodeAttribute,omitempty"`
	Impacts            []*Impact `json:"impacts"`
}

// Report defines a sonar report
type Report struct {
	Rules  []*Rule  `json:"rules"`
	Issues []*Issue `json:"issues"`
}
