package sarif

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/google/uuid"

	"github.com/securego/gosec/v2"
	"github.com/securego/gosec/v2/cwe"
	"github.com/securego/gosec/v2/issue"
)

// GenerateReport converts a gosec report into a SARIF report
func GenerateReport(rootPaths []string, data *gosec.ReportInfo) (*Report, error) {
	rules := []*ReportingDescriptor{}

	results := []*Result{}
	cweTaxa := []*ReportingDescriptor{}
	weaknesses := map[string]*cwe.Weakness{}

	for _, issue := range data.Issues {
		if issue.Cwe != nil {
			_, ok := weaknesses[issue.Cwe.ID]
			if !ok {
				weakness := cwe.Get(issue.Cwe.ID)
				weaknesses[issue.Cwe.ID] = weakness
				cweTaxon := parseSarifTaxon(weakness)
				cweTaxa = append(cweTaxa, cweTaxon)
			}
		}

		rule := parseSarifRule(issue)
		var ruleIndex int
		rules, ruleIndex = addRuleInOrder(rules, rule)

		location, err := parseSarifLocation(issue, rootPaths)
		if err != nil {
			return nil, err
		}

		result := NewResult(
			issue.RuleID,
			ruleIndex,
			getSarifLevel(issue.Severity.String()),
			issue.What,
			buildSarifSuppressions(issue.Suppressions),
			issue.Autofix,
		).WithLocations(location)

		results = append(results, result)
	}

	sort.SliceStable(cweTaxa, func(i, j int) bool { return cweTaxa[i].ID < cweTaxa[j].ID })

	tool := NewTool(buildSarifDriver(rules, data.GosecVersion))

	cweTaxonomy := buildCWETaxonomy(cweTaxa)

	run := NewRun(tool).
		WithTaxonomies(cweTaxonomy).
		WithResults(results...)

	return NewReport(Version, Schema).
		WithRuns(run), nil
}

// addRuleInOrder inserts a rule into the rules slice keeping the rules IDs order, it returns the new rules
// slice and the position where the rule was inserted
func addRuleInOrder(rules []*ReportingDescriptor, rule *ReportingDescriptor) ([]*ReportingDescriptor, int) {
	position := 0
	for i, r := range rules {
		if r.ID < rule.ID {
			continue
		}
		if r.ID == rule.ID {
			return rules, i
		}
		position = i
		break
	}
	rules = append(rules, nil)
	copy(rules[position+1:], rules[position:])
	rules[position] = rule
	return rules, position
}

// parseSarifRule return SARIF rule field struct
func parseSarifRule(i *issue.Issue) *ReportingDescriptor {
	cwe := issue.GetCweByRule(i.RuleID)
	name := i.RuleID
	if cwe != nil {
		name = cwe.Name
	}
	relationship := buildSarifReportingDescriptorRelationship(i.Cwe)
	rule := &ReportingDescriptor{
		ID:               i.RuleID,
		Name:             name,
		ShortDescription: NewMultiformatMessageString(i.What),
		FullDescription:  NewMultiformatMessageString(i.What),
		Help: NewMultiformatMessageString(fmt.Sprintf("%s\nSeverity: %s\nConfidence: %s\n",
			i.What, i.Severity.String(), i.Confidence.String())),
		Properties: &PropertyBag{
			"tags":      []string{"security", i.Severity.String()},
			"precision": strings.ToLower(i.Confidence.String()),
		},
		DefaultConfiguration: &ReportingConfiguration{
			Level: getSarifLevel(i.Severity.String()),
		},
	}
	if relationship != nil {
		rule.Relationships = []*ReportingDescriptorRelationship{relationship}
	}
	return rule
}

func buildSarifReportingDescriptorRelationship(weakness *cwe.Weakness) *ReportingDescriptorRelationship {
	if weakness == nil {
		return nil
	}
	return &ReportingDescriptorRelationship{
		Target: &ReportingDescriptorReference{
			ID:            weakness.ID,
			GUID:          uuid3(weakness.SprintID()),
			ToolComponent: NewToolComponentReference(cwe.Acronym),
		},
		Kinds: []string{"superset"},
	}
}

func buildCWETaxonomy(taxa []*ReportingDescriptor) *ToolComponent {
	return NewToolComponent(cwe.Acronym, cwe.Version, cwe.InformationURI).
		WithReleaseDateUtc(cwe.ReleaseDateUtc).
		WithDownloadURI(cwe.DownloadURI).
		WithOrganization(cwe.Organization).
		WithShortDescription(NewMultiformatMessageString(cwe.Description)).
		WithIsComprehensive(true).
		WithLanguage("en").
		WithMinimumRequiredLocalizedDataSemanticVersion(cwe.Version).
		WithTaxa(taxa...)
}

func parseSarifTaxon(weakness *cwe.Weakness) *ReportingDescriptor {
	return &ReportingDescriptor{
		ID:               weakness.ID,
		GUID:             uuid3(weakness.SprintID()),
		HelpURI:          weakness.SprintURL(),
		FullDescription:  NewMultiformatMessageString(weakness.Description),
		ShortDescription: NewMultiformatMessageString(weakness.Name),
	}
}

func parseSemanticVersion(version string) string {
	if len(version) == 0 {
		return "devel"
	}
	if strings.HasPrefix(version, "v") {
		return version[1:]
	}
	return version
}

func buildSarifDriver(rules []*ReportingDescriptor, gosecVersion string) *ToolComponent {
	semanticVersion := parseSemanticVersion(gosecVersion)
	return NewToolComponent("gosec", gosecVersion, "https://github.com/securego/gosec/").
		WithSemanticVersion(semanticVersion).
		WithSupportedTaxonomies(NewToolComponentReference(cwe.Acronym)).
		WithRules(rules...)
}

func uuid3(value string) string {
	return uuid.NewMD5(uuid.Nil, []byte(value)).String()
}

// parseSarifLocation return SARIF location struct
func parseSarifLocation(i *issue.Issue, rootPaths []string) (*Location, error) {
	region, err := parseSarifRegion(i)
	if err != nil {
		return nil, err
	}
	artifactLocation := parseSarifArtifactLocation(i, rootPaths)
	return NewLocation(NewPhysicalLocation(artifactLocation, region)), nil
}

func parseSarifArtifactLocation(i *issue.Issue, rootPaths []string) *ArtifactLocation {
	var filePath string
	for _, rootPath := range rootPaths {
		if strings.HasPrefix(i.File, rootPath) {
			filePath = strings.Replace(i.File, rootPath+"/", "", 1)
		}
	}
	return NewArtifactLocation(filePath)
}

func parseSarifRegion(i *issue.Issue) (*Region, error) {
	lines := strings.Split(i.Line, "-")
	startLine, err := strconv.Atoi(lines[0])
	if err != nil {
		return nil, err
	}
	endLine := startLine
	if len(lines) > 1 {
		endLine, err = strconv.Atoi(lines[1])
		if err != nil {
			return nil, err
		}
	}
	col, err := strconv.Atoi(i.Col)
	if err != nil {
		return nil, err
	}
	var code string
	line := startLine
	codeLines := strings.Split(i.Code, "\n")
	for _, codeLine := range codeLines {
		lineStart := fmt.Sprintf("%d:", line)
		if strings.HasPrefix(codeLine, lineStart) {
			code += strings.TrimSpace(
				strings.TrimPrefix(codeLine, lineStart))
			if endLine > startLine {
				code += "\n"
			}
			line++
			if line > endLine {
				break
			}
		}
	}
	snippet := NewArtifactContent(code)
	return NewRegion(startLine, endLine, col, col, "go").WithSnippet(snippet), nil
}

func getSarifLevel(s string) Level {
	switch s {
	case "LOW":
		return Warning
	case "MEDIUM":
		return Error
	case "HIGH":
		return Error
	default:
		return Note
	}
}

func buildSarifSuppressions(suppressions []issue.SuppressionInfo) []*Suppression {
	var sarifSuppressionList []*Suppression
	for _, s := range suppressions {
		sarifSuppressionList = append(sarifSuppressionList, NewSuppression(s.Kind, s.Justification))
	}
	return sarifSuppressionList
}
