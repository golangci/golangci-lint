package commands

import (
	"image/color"

	"charm.land/fang/v2"
	"charm.land/lipgloss/v2"
)

// Colors from https://golangci-lint.run
const (
	blue600 = "#2563eb" // primary
	blue700 = "#1d4ed8" // primary hover
	blue400 = "#60a5fa" // primary medium
	teal500 = "#0d9488" // accent (flags)
	gray900 = "#111827" // headings
	gray700 = "#374151" // body text
	gray50  = "#f9fafb" // light bg
	gray400 = "#9ca3af" // muted
)

func golangciColorScheme(ld lipgloss.LightDarkFunc) fang.ColorScheme {
	return fang.ColorScheme{
		Base:           ld(lipgloss.Color(gray900), lipgloss.Color(gray50)),
		Title:          lipgloss.Color(blue600),
		Description:    ld(lipgloss.Color(gray700), lipgloss.Color(gray400)),
		Codeblock:      ld(lipgloss.Color(gray50), lipgloss.Color("#1e1e2e")),
		Program:        lipgloss.Color(blue600),
		Command:        lipgloss.Color(teal500),
		DimmedArgument: ld(lipgloss.Color(gray400), lipgloss.Color("#6c7086")),
		Comment:        ld(lipgloss.Color(gray400), lipgloss.Color("#6c7086")),
		Flag:           lipgloss.Color(teal500),
		FlagDefault:    ld(lipgloss.Color(gray400), lipgloss.Color("#585b70")),
		Argument:       ld(lipgloss.Color(gray900), lipgloss.Color(gray50)),
		QuotedString:   lipgloss.Color(blue400),
		Help:           ld(lipgloss.Color(gray700), lipgloss.Color(gray400)),
		Dash:           ld(lipgloss.Color(gray400), lipgloss.Color("#6c7086")),
		ErrorHeader: [2]color.Color{
			lipgloss.Color("#ffffff"),
			lipgloss.Color(blue600),
		},
		ErrorDetails: lipgloss.Color(blue700),
	}
}
