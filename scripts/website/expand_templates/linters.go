package main

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"

	"gopkg.in/yaml.v3"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/scripts/website/types"
)

const listItemPrefix = "list-item-"

func getExampleSnippets() (*SettingSnippets, error) {
	reference, err := os.ReadFile(".golangci.reference.yml")
	if err != nil {
		return nil, fmt.Errorf("can't read .golangci.reference.yml: %w", err)
	}

	snippets, err := extractExampleSnippets(reference)
	if err != nil {
		return nil, fmt.Errorf("can't extract example snippets from .golangci.reference.yml: %w", err)
	}

	return snippets, nil
}

func getLintersListMarkdown(enabled bool) string {
	linters, err := readJSONFile[[]*types.LinterWrapper](filepath.Join("assets", "linters-info.json"))
	if err != nil {
		panic(err)
	}

	var neededLcs []*types.LinterWrapper
	for _, lc := range linters {
		if lc.Internal {
			continue
		}

		if lc.EnabledByDefault == enabled {
			neededLcs = append(neededLcs, lc)
		}
	}

	sort.Slice(neededLcs, func(i, j int) bool {
		return neededLcs[i].Name < neededLcs[j].Name
	})

	lines := []string{
		"|Name|Description|Presets|AutoFix|Since|",
		"|---|---|---|---|---|---|",
	}

	for _, lc := range neededLcs {
		line := fmt.Sprintf("|%s|%s|%s|%v|%s|",
			getName(lc),
			getDesc(lc),
			strings.Join(lc.InPresets, ", "),
			check(lc.CanAutoFix, "Auto fix supported"),
			lc.Since,
		)
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

func getName(lc *types.LinterWrapper) string {
	name := lc.Name

	if lc.OriginalURL != "" {
		name = fmt.Sprintf("[%s](%s)", name, lc.OriginalURL)
	}

	if hasSettings(lc.Name) {
		name = fmt.Sprintf("%s&nbsp;[%s](#%s)", name, spanWithID(listItemPrefix+lc.Name, "Configuration", "⚙️"), lc.Name)
	}

	if lc.Deprecation == nil {
		return name
	}

	title := "deprecated"
	if lc.Deprecation.Replacement != "" {
		title += fmt.Sprintf(" since %s", lc.Deprecation.Since)
	}

	return name + "&nbsp;" + span(title, "⚠")
}

func check(b bool, title string) string {
	if b {
		return span(title, "✔")
	}
	return ""
}

func getDesc(lc *types.LinterWrapper) string {
	desc := lc.Desc
	if lc.Deprecation != nil {
		desc = lc.Deprecation.Message
		if lc.Deprecation.Replacement != "" {
			desc += fmt.Sprintf(" Replaced by %s.", lc.Deprecation.Replacement)
		}
	}

	return formatDesc(desc)
}

func formatDesc(desc string) string {
	runes := []rune(desc)

	r, _ := utf8.DecodeRuneInString(desc)
	runes[0] = unicode.ToUpper(r)

	if runes[len(runes)-1] != '.' {
		runes = append(runes, '.')
	}

	return strings.ReplaceAll(string(runes), "\n", "<br/>")
}

func hasSettings(name string) bool {
	tp := reflect.TypeOf(config.LintersSettings{})

	for i := 0; i < tp.NumField(); i++ {
		if strings.EqualFold(name, tp.Field(i).Name) {
			return true
		}
	}

	return false
}

func span(title, icon string) string {
	return fmt.Sprintf(`<span title=%q>%s</span>`, title, icon)
}

func spanWithID(id, title, icon string) string {
	return fmt.Sprintf(`<span id=%q title=%q>%s</span>`, id, title, icon)
}

type SettingSnippets struct {
	ConfigurationFile string
	LintersSettings   string
}

func extractExampleSnippets(example []byte) (*SettingSnippets, error) {
	var data yaml.Node
	err := yaml.Unmarshal(example, &data)
	if err != nil {
		return nil, err
	}

	root := data.Content[0]

	globalNode := &yaml.Node{
		Kind:        root.Kind,
		Style:       root.Style,
		Tag:         root.Tag,
		Value:       root.Value,
		Anchor:      root.Anchor,
		Alias:       root.Alias,
		HeadComment: root.HeadComment,
		LineComment: root.LineComment,
		FootComment: root.FootComment,
		Line:        root.Line,
		Column:      root.Column,
	}

	snippets := SettingSnippets{}

	builder := strings.Builder{}

	for j, node := range root.Content {
		switch node.Value {
		case "run", "output", "linters", "linters-settings", "issues", "severity":
		default:
			continue
		}

		nextNode := root.Content[j+1]

		newNode := &yaml.Node{
			Kind: nextNode.Kind,
			Content: []*yaml.Node{
				{
					HeadComment: fmt.Sprintf("See the dedicated %q documentation section.", node.Value),
					Kind:        node.Kind,
					Style:       node.Style,
					Tag:         node.Tag,
					Value:       "option",
				},
				{
					Kind:  node.Kind,
					Style: node.Style,
					Tag:   node.Tag,
					Value: "value",
				},
			},
		}

		globalNode.Content = append(globalNode.Content, node, newNode)

		if node.Value == "linters-settings" {
			snippets.LintersSettings, err = getLintersSettingSections(node, nextNode)
			if err != nil {
				return nil, err
			}

			_, _ = builder.WriteString(
				fmt.Sprintf(
					"### `%s` configuration\n\nSee the dedicated [linters-settings](/usage/linters) documentation section.\n\n",
					node.Value,
				),
			)
			continue
		}

		nodeSection := &yaml.Node{
			Kind:    root.Kind,
			Style:   root.Style,
			Tag:     root.Tag,
			Value:   root.Value,
			Content: []*yaml.Node{node, nextNode},
		}

		snippet, errSnip := marshallSnippet(nodeSection)
		if errSnip != nil {
			return nil, errSnip
		}

		_, _ = builder.WriteString(fmt.Sprintf("### `%s` configuration\n\n%s", node.Value, snippet))
	}

	overview, err := marshallSnippet(globalNode)
	if err != nil {
		return nil, err
	}

	snippets.ConfigurationFile = overview + builder.String()

	return &snippets, nil
}

func getLintersSettingSections(node, nextNode *yaml.Node) (string, error) {
	linters, err := readJSONFile[[]*types.LinterWrapper](filepath.Join("assets", "linters-info.json"))
	if err != nil {
		return "", err
	}

	lintersDesc := make(map[string]string)
	for _, lc := range linters {
		if lc.Internal {
			continue
		}

		// it's important to use lc.Name() nor name because name can be alias
		lintersDesc[lc.Name] = getDesc(lc)
	}

	builder := &strings.Builder{}

	for i := 0; i < len(nextNode.Content); i += 2 {
		r := &yaml.Node{
			Kind:  nextNode.Kind,
			Style: nextNode.Style,
			Tag:   nextNode.Tag,
			Value: node.Value,
			Content: []*yaml.Node{
				{
					Kind:  node.Kind,
					Value: node.Value,
				},
				{
					Kind:    nextNode.Kind,
					Content: []*yaml.Node{nextNode.Content[i], nextNode.Content[i+1]},
				},
			},
		}

		_, _ = fmt.Fprintf(builder, "### %s\n\n", nextNode.Content[i].Value)
		_, _ = fmt.Fprintf(builder, "%s\n\n", lintersDesc[nextNode.Content[i].Value])
		_, _ = fmt.Fprintln(builder, "```yaml")

		encoder := yaml.NewEncoder(builder)
		encoder.SetIndent(2)

		err := encoder.Encode(r)
		if err != nil {
			return "", err
		}

		_, _ = fmt.Fprintln(builder, "```")
		_, _ = fmt.Fprintln(builder)
		_, _ = fmt.Fprintf(builder, "[%s](#%s)\n\n", span("Back to the top", "🔼"), listItemPrefix+nextNode.Content[i].Value)
		_, _ = fmt.Fprintln(builder)
	}

	return builder.String(), nil
}

func marshallSnippet(node *yaml.Node) (string, error) {
	builder := &strings.Builder{}

	if node.Value != "" {
		_, _ = fmt.Fprintf(builder, "### %s\n\n", node.Value)
	}
	_, _ = fmt.Fprintln(builder, "```yaml")

	encoder := yaml.NewEncoder(builder)
	encoder.SetIndent(2)

	err := encoder.Encode(node)
	if err != nil {
		return "", err
	}

	_, _ = fmt.Fprintln(builder, "```")
	_, _ = fmt.Fprintln(builder)

	return builder.String(), nil
}
