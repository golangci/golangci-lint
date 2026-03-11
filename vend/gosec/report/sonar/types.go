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
	EngineID           string      `json:"engineId"`
	RuleID             string      `json:"ruleId"`
	PrimaryLocation    *Location   `json:"primaryLocation"`
	Type               string      `json:"type"`
	Severity           string      `json:"severity"`
	EffortMinutes      int         `json:"effortMinutes"`
	SecondaryLocations []*Location `json:"secondaryLocations,omitempty"`
}

// Report defines a sonar report
type Report struct {
	Issues []*Issue `json:"issues"`
}
