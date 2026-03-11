package sonar

// NewLocation instantiate a Location
func NewLocation(message string, filePath string, textRange *TextRange) *Location {
	return &Location{
		Message:   message,
		FilePath:  filePath,
		TextRange: textRange,
	}
}

// NewTextRange instantiate a TextRange
func NewTextRange(startLine int, endLine int) *TextRange {
	return &TextRange{
		StartLine: startLine,
		EndLine:   endLine,
	}
}

// NewIssue instantiate an Issue
func NewIssue(engineID string, ruleID string, primaryLocation *Location, issueType string, severity string, effortMinutes int) *Issue {
	return &Issue{
		EngineID:        engineID,
		RuleID:          ruleID,
		PrimaryLocation: primaryLocation,
		Type:            issueType,
		Severity:        severity,
		EffortMinutes:   effortMinutes,
	}
}
