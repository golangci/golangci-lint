package printers

import (
	"context"
	"fmt"
	"io"
	"strings"
	"unicode/utf8"

	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

// Field limits.
const (
	smallLimit = 255
	largeLimit = 4000
)

// Configer used for accessing linter.Config by its name for printing instanceType values.
type Configer interface {
	GetLinterConfigs(name string) []*linter.Config
}

// TeamCity printer for TeamCity format.
type TeamCity struct {
	w       io.Writer
	conf    Configer
	escaper *strings.Replacer
}

// NewTeamCity output format outputs issues according to TeamCity service message format.
func NewTeamCity(w io.Writer, conf Configer) *TeamCity {
	return &TeamCity{
		w:    w,
		conf: conf,
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

func (p *TeamCity) Print(_ context.Context, issues []result.Issue) error {
	uniqLinters := map[string]struct{}{}

	for i := range issues {
		issue := issues[i]

		_, ok := uniqLinters[issue.FromLinter]
		if !ok {
			linterConfigs := p.conf.GetLinterConfigs(issue.FromLinter)
			for _, config := range linterConfigs {
				inspectionType := inspectionType{
					id:          config.Linter.Name(),
					name:        config.Linter.Name(),
					description: config.Linter.Desc(),
					category:    "Golangci-lint reports",
				}

				_, err := inspectionType.Print(p.w, p.escaper)
				if err != nil {
					return err
				}
			}

			uniqLinters[issue.FromLinter] = struct{}{}
		}

		instance := inspectionInstance{
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

// inspectionType is the unique description of the conducted inspection. Each specific warning or
// an error in code (inspection instance) has an inspection type.
// https://www.jetbrains.com/help/teamcity/service-messages.html#Inspection+Type
type inspectionType struct {
	id          string // (mandatory) limited by 255 characters.
	name        string // (mandatory) limited by 255 characters.
	description string // (mandatory) limited by 255 characters.
	category    string // (mandatory) limited by 4000 characters.
}

func (i inspectionType) Print(w io.Writer, escaper *strings.Replacer) (int, error) {
	return fmt.Fprintf(w, "##teamcity[inspectionType id='%s' name='%s' description='%s' category='%s']\n",
		limit(i.id, smallLimit), limit(i.name, smallLimit), limit(escaper.Replace(i.description), largeLimit), limit(i.category, smallLimit))
}

// inspectionInstance reports a specific defect, warning, error message.
// Includes location, description, and various optional and custom attributes.
// https://www.jetbrains.com/help/teamcity/service-messages.html#Inspection+Instance
type inspectionInstance struct {
	typeID   string // (mandatory) limited by 255 characters.
	message  string // (optional)  limited by 4000 characters.
	file     string // (mandatory) file path limited by 4000 characters.
	line     int    // (optional)  line of the file.
	severity string // (optional)  severity attribute: INFO, ERROR, WARNING, WEAK WARNING.
}

func (i inspectionInstance) Print(w io.Writer, replacer *strings.Replacer) (int, error) {
	return fmt.Fprintf(w, "##teamcity[inspection typeId='%s' message='%s' file='%s' line='%d' SEVERITY='%s']\n",
		limit(i.typeID, smallLimit), limit(replacer.Replace(i.message), largeLimit), limit(i.file, largeLimit), i.line,
		strings.ToUpper(i.severity))
}

func limit(s string, max int) string {
	var size, count int
	for i := 0; i < max && count < len(s); i++ {
		_, size = utf8.DecodeRuneInString(s[count:])
		count += size
	}

	return s[:count]
}
