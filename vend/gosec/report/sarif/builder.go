package sarif

// NewReport instantiate a SARIF Report
func NewReport(version string, schema string) *Report {
	return &Report{
		Version: version,
		Schema:  schema,
	}
}

// WithRuns defines runs for the current report
func (r *Report) WithRuns(runs ...*Run) *Report {
	r.Runs = runs
	return r
}

// NewMultiformatMessageString instantiate a MultiformatMessageString
func NewMultiformatMessageString(text string) *MultiformatMessageString {
	return &MultiformatMessageString{
		Text: text,
	}
}

// NewRun instantiate a Run
func NewRun(tool *Tool) *Run {
	return &Run{
		Tool: tool,
	}
}

// WithTaxonomies set the taxonomies for the current run
func (r *Run) WithTaxonomies(taxonomies ...*ToolComponent) *Run {
	r.Taxonomies = taxonomies
	return r
}

// WithResults set the results for the current run
func (r *Run) WithResults(results ...*Result) *Run {
	r.Results = results
	return r
}

// NewArtifactLocation instantiate an ArtifactLocation
func NewArtifactLocation(uri string) *ArtifactLocation {
	return &ArtifactLocation{
		URI: uri,
	}
}

// NewRegion instantiate a Region
func NewRegion(startLine int, endLine int, startColumn int, endColumn int, sourceLanguage string) *Region {
	return &Region{
		StartLine:      startLine,
		EndLine:        endLine,
		StartColumn:    startColumn,
		EndColumn:      endColumn,
		SourceLanguage: sourceLanguage,
	}
}

// WithSnippet defines the Snippet for the current Region
func (r *Region) WithSnippet(snippet *ArtifactContent) *Region {
	r.Snippet = snippet
	return r
}

// NewArtifactContent instantiate an ArtifactContent
func NewArtifactContent(text string) *ArtifactContent {
	return &ArtifactContent{
		Text: text,
	}
}

// NewTool instantiate a Tool
func NewTool(driver *ToolComponent) *Tool {
	return &Tool{
		Driver: driver,
	}
}

// NewResult instantiate a Result
func NewResult(ruleID string, ruleIndex int, level Level, message string, suppressions []*Suppression, autofix string) *Result {
	result := &Result{
		RuleID:       ruleID,
		RuleIndex:    ruleIndex,
		Level:        level,
		Message:      NewMessage(message),
		Suppressions: suppressions,
	}

	// Only create Fix when autofix content exists
	// Fixes with nil/null ArtifactChanges violate SARIF 2.1.0 schema
	if autofix != "" {
		result.Fixes = []*Fix{
			{
				Description: &Message{
					Text:     autofix,
					Markdown: autofix,
				},
				// ArtifactChanges MUST be a non-empty array per SARIF 2.1.0 schema
				ArtifactChanges: []*ArtifactChange{
					{
						ArtifactLocation: &ArtifactLocation{
							Description: NewMessage("File requiring changes"),
						},
						Replacements: []*Replacement{
							{
								DeletedRegion: NewRegion(1, 1, 1, 1, ""),
							},
						},
					},
				},
			},
		}
	}
	return result
}

// NewMessage instantiate a Message
func NewMessage(text string) *Message {
	return &Message{
		Text: text,
	}
}

// WithLocations define the current result's locations
func (r *Result) WithLocations(locations ...*Location) *Result {
	r.Locations = locations
	return r
}

// NewLocation instantiate a Location
func NewLocation(physicalLocation *PhysicalLocation) *Location {
	return &Location{
		PhysicalLocation: physicalLocation,
	}
}

// NewPhysicalLocation instantiate a PhysicalLocation
func NewPhysicalLocation(artifactLocation *ArtifactLocation, region *Region) *PhysicalLocation {
	return &PhysicalLocation{
		ArtifactLocation: artifactLocation,
		Region:           region,
	}
}

// NewToolComponent instantiate a ToolComponent
func NewToolComponent(name string, version string, informationURI string) *ToolComponent {
	return &ToolComponent{
		Name:           name,
		Version:        version,
		InformationURI: informationURI,
		GUID:           uuid3(name),
	}
}

// WithLanguage set Language for the current ToolComponent
func (t *ToolComponent) WithLanguage(language string) *ToolComponent {
	t.Language = language
	return t
}

// WithSemanticVersion set SemanticVersion for the current ToolComponent
func (t *ToolComponent) WithSemanticVersion(semanticVersion string) *ToolComponent {
	t.SemanticVersion = semanticVersion
	return t
}

// WithReleaseDateUtc set releaseDateUtc for the current ToolComponent
func (t *ToolComponent) WithReleaseDateUtc(releaseDateUtc string) *ToolComponent {
	t.ReleaseDateUtc = releaseDateUtc
	return t
}

// WithDownloadURI set downloadURI for the current ToolComponent
func (t *ToolComponent) WithDownloadURI(downloadURI string) *ToolComponent {
	t.DownloadURI = downloadURI
	return t
}

// WithOrganization set organization for the current ToolComponent
func (t *ToolComponent) WithOrganization(organization string) *ToolComponent {
	t.Organization = organization
	return t
}

// WithShortDescription set shortDescription for the current ToolComponent
func (t *ToolComponent) WithShortDescription(shortDescription *MultiformatMessageString) *ToolComponent {
	t.ShortDescription = shortDescription
	return t
}

// WithIsComprehensive set isComprehensive for the current ToolComponent
func (t *ToolComponent) WithIsComprehensive(isComprehensive bool) *ToolComponent {
	t.IsComprehensive = isComprehensive
	return t
}

// WithMinimumRequiredLocalizedDataSemanticVersion set MinimumRequiredLocalizedDataSemanticVersion for the current ToolComponent
func (t *ToolComponent) WithMinimumRequiredLocalizedDataSemanticVersion(minimumRequiredLocalizedDataSemanticVersion string) *ToolComponent {
	t.MinimumRequiredLocalizedDataSemanticVersion = minimumRequiredLocalizedDataSemanticVersion
	return t
}

// WithTaxa set taxa for the current ToolComponent
func (t *ToolComponent) WithTaxa(taxa ...*ReportingDescriptor) *ToolComponent {
	t.Taxa = taxa
	return t
}

// WithSupportedTaxonomies set the supported taxonomies for the current ToolComponent
func (t *ToolComponent) WithSupportedTaxonomies(supportedTaxonomies ...*ToolComponentReference) *ToolComponent {
	t.SupportedTaxonomies = supportedTaxonomies
	return t
}

// WithRules set the rules for the current ToolComponent
func (t *ToolComponent) WithRules(rules ...*ReportingDescriptor) *ToolComponent {
	t.Rules = rules
	return t
}

// NewToolComponentReference instantiate a ToolComponentReference
func NewToolComponentReference(name string) *ToolComponentReference {
	return &ToolComponentReference{
		Name: name,
		GUID: uuid3(name),
	}
}

// NewSuppression instantiate a Suppression
func NewSuppression(kind string, justification string) *Suppression {
	return &Suppression{
		Kind:          kind,
		Justification: justification,
	}
}
