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
func NewIssue(ruleID string, primaryLocation *Location, effortMinutes int) *Issue {
	return &Issue{
		RuleID:          ruleID,
		PrimaryLocation: primaryLocation,
		EffortMinutes:   effortMinutes,
	}
}

// NewImpact instantiate an Impact.
func NewImpact(softwareQuality string, severity string) *Impact {
	return &Impact{
		SoftwareQuality: softwareQuality,
		Severity:        severity,
	}
}

// NewRule instantiate a Rule.
func NewRule(id string, name string, description string, engineID string, cleanCodeAttribute string, impacts []*Impact) *Rule {
	return &Rule{
		ID:                 id,
		Name:               name,
		Description:        description,
		EngineID:           engineID,
		CleanCodeAttribute: cleanCodeAttribute,
		Impacts:            impacts,
	}
}
