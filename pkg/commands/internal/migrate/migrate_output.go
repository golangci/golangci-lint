package migrate

import (
	"slices"

	"github.com/golangci/golangci-lint/v2/pkg/commands/internal/migrate/ptr"
	"github.com/golangci/golangci-lint/v2/pkg/commands/internal/migrate/versionone"
	"github.com/golangci/golangci-lint/v2/pkg/commands/internal/migrate/versiontwo"
)

func toOutput(old *versionone.Config) versiontwo.Output {
	formats := versiontwo.Formats{}

	oldFormats := cleanIncompatibleFormats(old.Output.Formats, "colored-line-number", "line-number")
	oldFormats = cleanIncompatibleFormats(oldFormats, "colored-tab", "tab")
	oldFormats = cleanIncompatibleFormats(oldFormats, "junit-xml-extended", "junit-xml")

	for _, format := range oldFormats {
		switch ptr.Deref(format.Format) {
		case "colored-line-number":
			formats.Text.PrintLinterName = old.Output.PrintLinterName
			formats.Text.PrintIssuedLine = old.Output.PrintIssuedLine
			formats.Text.Colors = nil // color is true by default (flags).
			formats.Text.Path = new(defaultFormatPath(ptr.Deref(format.Path)))

		case "line-number":
			formats.Text.PrintLinterName = old.Output.PrintLinterName
			formats.Text.PrintIssuedLine = old.Output.PrintIssuedLine
			formats.Text.Colors = new(false)
			formats.Text.Path = new(defaultFormatPath(ptr.Deref(format.Path)))

		case "json":
			formats.JSON.Path = new(defaultFormatPath(ptr.Deref(format.Path)))

		case "colored-tab":
			formats.Tab.PrintLinterName = old.Output.PrintLinterName
			formats.Tab.Colors = nil // Colors is true by default (flags).
			formats.Tab.Path = new(defaultFormatPath(ptr.Deref(format.Path)))

		case "tab":
			formats.Tab.PrintLinterName = old.Output.PrintLinterName
			formats.Tab.Colors = new(false)
			formats.Tab.Path = new(defaultFormatPath(ptr.Deref(format.Path)))

		case "html":
			formats.HTML.Path = new(defaultFormatPath(ptr.Deref(format.Path)))

		case "checkstyle":
			formats.Checkstyle.Path = new(defaultFormatPath(ptr.Deref(format.Path)))

		case "code-climate":
			formats.CodeClimate.Path = new(defaultFormatPath(ptr.Deref(format.Path)))

		case "junit-xml":
			formats.JUnitXML.Extended = nil // Extended is false by default.
			formats.JUnitXML.Path = new(defaultFormatPath(ptr.Deref(format.Path)))

		case "junit-xml-extended":
			formats.JUnitXML.Extended = new(true)
			formats.JUnitXML.Path = new(defaultFormatPath(ptr.Deref(format.Path)))

		case "github-actions":
			// Ignored

		case "teamcity":
			formats.TeamCity.Path = new(defaultFormatPath(ptr.Deref(format.Path)))

		case "sarif":
			formats.Sarif.Path = new(defaultFormatPath(ptr.Deref(format.Path)))
		}
	}

	return versiontwo.Output{
		Formats:    formats,
		SortOrder:  old.Output.SortOrder,
		PathPrefix: old.Output.PathPrefix,
		ShowStats:  nil, // Enforce the new default. (nil -> omitempty -> true)
	}
}

func defaultFormatPath(p string) string {
	if p == "" {
		return "stdout"
	}

	return p
}

func cleanIncompatibleFormats(old versionone.OutputFormats, f1, f2 string) versionone.OutputFormats {
	index1 := slices.IndexFunc(old, func(format versionone.OutputFormat) bool {
		return ptr.Deref(format.Format) == f1
	})

	index2 := slices.IndexFunc(old, func(format versionone.OutputFormat) bool {
		return ptr.Deref(format.Format) == f2
	})

	if index1 >= 0 && index2 >= 0 {
		return slices.Delete(old, index2, index2+1)
	}

	return old
}
