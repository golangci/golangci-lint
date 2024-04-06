package printers

import (
	"fmt"
	"io"
	"strings"
	"unicode/utf8"

	"github.com/golangci/golangci-lint/pkg/result"
)

// Field limits.
const (
	smallLimit = 255
	largeLimit = 4000
)

// TeamCity printer for TeamCity format.
type TeamCity struct {
	w       io.Writer
	escaper *strings.Replacer
}

// NewTeamCity output format outputs issues according to TeamCity service message format.
func NewTeamCity(w io.Writer) *TeamCity {
	return &TeamCity{
		w: w,
		// https://www.jetbrains.com/help/teamcity/service-messages.html#Escaped+Values
		escaper: strings.NewReplacer(
			"'", "|'",
			"\n", "|n",
			"\r", "|r",
			"|", "||",
			"[", "|[",
			"]", "|]",
		),
	}
}

func (p *TeamCity) Print(issues []result.Issue) error {
	uniqLinters := map[string]struct{}{}

	for i := range issues {
		issue := issues[i]

		_, ok := uniqLinters[issue.FromLinter]
		if !ok {
			inspectionType := InspectionType{
				id:          issue.FromLinter,
				name:        issue.FromLinter,
				description: issue.FromLinter,
				category:    "Golangci-lint reports",
			}

			_, err := inspectionType.Print(p.w, p.escaper)
			if err != nil {
				return err
			}

			uniqLinters[issue.FromLinter] = struct{}{}
		}

		instance := InspectionInstance{
			typeID:   issue.FromLinter,
			message:  issue.Text,
			file:     issue.FilePath(),
			line:     issue.Line(),
			severity: issue.Severity,
		}

		_, err := instance.Print(p.w, p.escaper)
		if err != nil {
			return err
		}
	}

	return nil
}

// InspectionType is the unique description of the conducted inspection. Each specific warning or
// an error in code (inspection instance) has an inspection type.
// https://www.jetbrains.com/help/teamcity/service-messages.html#Inspection+Type
type InspectionType struct {
	id          string // (mandatory) limited by 255 characters.
	name        string // (mandatory) limited by 255 characters.
	description string // (mandatory) limited by 255 characters.
	category    string // (mandatory) limited by 4000 characters.
}

func (i InspectionType) Print(w io.Writer, escaper *strings.Replacer) (int, error) {
	return fmt.Fprintf(w, "##teamcity[inspectionType id='%s' name='%s' description='%s' category='%s']\n",
		cutVal(i.id, smallLimit), cutVal(i.name, smallLimit), cutVal(escaper.Replace(i.description), largeLimit), cutVal(i.category, smallLimit))
}

// InspectionInstance reports a specific defect, warning, error message.
// Includes location, description, and various optional and custom attributes.
// https://www.jetbrains.com/help/teamcity/service-messages.html#Inspection+Instance
type InspectionInstance struct {
	typeID   string // (mandatory) limited by 255 characters.
	message  string // (optional)  limited by 4000 characters.
	file     string // (mandatory) file path limited by 4000 characters.
	line     int    // (optional)  line of the file.
	severity string // (optional) any linter severity.
}

func (i InspectionInstance) Print(w io.Writer, replacer *strings.Replacer) (int, error) {
	return fmt.Fprintf(w, "##teamcity[inspection typeId='%s' message='%s' file='%s' line='%d' SEVERITY='%s']\n",
		cutVal(i.typeID, smallLimit),
		cutVal(replacer.Replace(i.message), largeLimit),
		cutVal(i.file, largeLimit),
		i.line, strings.ToUpper(i.severity))
}

func cutVal(s string, limit int) string {
	var size, count int
	for i := 0; i < limit && count < len(s); i++ {
		_, size = utf8.DecodeRuneInString(s[count:])
		count += size
	}

	return s[:count]
}
