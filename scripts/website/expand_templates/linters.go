package main

import (
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"

	"gopkg.in/yaml.v3"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/scripts/website/expand_templates/internal"
	"github.com/golangci/golangci-lint/v2/scripts/website/types"
)

const (
	keyLinters    = "linters"
	keyFormatters = "formatters"
	keySettings   = "settings"
)

func getLintersListMarkdown(enabled bool, src, section, latestVersion string) string {
	linters, err := readJSONFile[[]*types.LinterWrapper](src)
	if err != nil {
		panic(err)
	}

	var neededLcs []*types.LinterWrapper
	for _, lc := range linters {
		if lc.Internal {
			continue
		}

		if slices.Contains(slices.Collect(maps.Keys(lc.Groups)), config.GroupStandard) == enabled {
			neededLcs = append(neededLcs, lc)
		}
	}

	sort.Slice(neededLcs, func(i, j int) bool {
		return neededLcs[i].Name < neededLcs[j].Name
	})

	slices.SortFunc(neededLcs, func(a, b *types.LinterWrapper) int {
		if a.IsDeprecated() && b.IsDeprecated() {
			return strings.Compare(a.Name, b.Name)
		}

		if a.IsDeprecated() {
			return 1
		}

		if b.IsDeprecated() {
			return -1
		}

		return strings.Compare(a.Name, b.Name)
	})

	cards := internal.NewCards()

	for _, lc := range neededLcs {
		cards.Add(internal.NewCard().
			Link(fmt.Sprintf("/docs/usage/%s/configuration/#%s", section, lc.Name)).
			Title(lc.Name).
			Subtitle(getDesc(lc)).
			Tag(getTag(lc, latestVersion)),
		)
	}

	return cards.String()
}

func getTag(lc *types.LinterWrapper, latestVersion string) (content, style string) {
	if lc.Deprecation != nil {
		tagContent := "Deprecated"
		if lc.Deprecation.Replacement != "" {
			tagContent += fmt.Sprintf(" since %s", lc.Deprecation.Since)
		}

		return tagContent, "error"
	}

	if compareVersion(lc.Since, latestVersion) {
		return "New", "warning"
	}

	if lc.CanAutoFix {
		return "Autofix", "info"
	}

	return "", ""
}

func compareVersion(a, b string) bool {
	return a[:strings.LastIndex(a, ".")] == b[:strings.LastIndex(b, ".")]
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

	return strings.NewReplacer("\n", "<br/>", `"`, `'`).Replace(string(runes))
}

type SettingSnippets struct {
	ConfigurationFile  string
	LintersSettings    string
	FormattersSettings string
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

type ExampleSnippetsExtractor struct {
	referencePath string
	assetsPath    string
}

func NewExampleSnippetsExtractor() *ExampleSnippetsExtractor {
	return &ExampleSnippetsExtractor{
		referencePath: ".golangci.reference.yml",
		assetsPath:    "assets",
	}
}

func (e *ExampleSnippetsExtractor) GetExampleSnippets() (*SettingSnippets, error) {
	reference, err := os.ReadFile(e.referencePath)
	if err != nil {
		return nil, fmt.Errorf("can't read .golangci.reference.yml: %w", err)
	}

	snippets, err := e.extractExampleSnippets(reference)
	if err != nil {
		return nil, fmt.Errorf("can't extract example snippets from .golangci.reference.yml: %w", err)
	}

	return snippets, nil
}

//nolint:gocyclo // The complexity is expected because of raw YAML manipulations.
func (e *ExampleSnippetsExtractor) extractExampleSnippets(example []byte) (*SettingSnippets, error) {
	var data yaml.Node
	if err := yaml.Unmarshal(example, &data); err != nil {
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
		case "run", "output", keyLinters, keyFormatters, "issues", "severity", "version":
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

		if node.Value == "version" {
			n := &yaml.Node{
				HeadComment: fmt.Sprintf("See the dedicated %q documentation section.", node.Value),
				Kind:        node.Kind,
				Style:       node.Style,
				Tag:         node.Tag,
				Value:       node.Value,
				Content:     node.Content,
			}

			globalNode.Content = append(globalNode.Content, n, nextNode)
		} else {
			globalNode.Content = append(globalNode.Content, node, newNode)
		}

		if node.Value == keyLinters || node.Value == keyFormatters {
			for i := 0; i < len(nextNode.Content); i++ {
				if nextNode.Content[i].Value != keySettings {
					continue
				}

				settingSections, err := e.getSettingSections(node, nextNode.Content[i+1])
				if err != nil {
					return nil, err
				}

				switch node.Value {
				case keyLinters:
					snippets.LintersSettings = settingSections

				case keyFormatters:
					snippets.FormattersSettings = settingSections
				}

				nextNode.Content[i+1].Content = []*yaml.Node{
					{
						HeadComment: fmt.Sprintf(`See the dedicated "%s.%s" documentation section.`, node.Value, nextNode.Content[i].Value),
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
				}

				i++
			}
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

		_, _ = builder.WriteString(fmt.Sprintf("## `%s` configuration\n\n", node.Value))

		if node.Value == keyLinters || node.Value == keyFormatters {
			baseTitle := []rune(node.Value)
			r, _ := utf8.DecodeRuneInString(node.Value)
			baseTitle[0] = unicode.ToUpper(r)

			builder.WriteString(internal.NewCards().
				Cols(2).
				Add(internal.NewCard().
					Link(fmt.Sprintf("/docs/usage/%s", node.Value)).
					Title(string(baseTitle) + " Overview").
					Icon("collection")).
				Add(internal.NewCard().
					Link(fmt.Sprintf("/docs/usage/%s/configuration", node.Value)).
					Title(string(baseTitle) + "  Settings").
					Icon("adjustments")).
				String())
		}

		_, _ = builder.WriteString(fmt.Sprintf("\n\n%s", snippet))
	}

	overview, err := marshallSnippet(globalNode)
	if err != nil {
		return nil, err
	}

	snippets.ConfigurationFile = overview + builder.String()

	return &snippets, nil
}

func (e *ExampleSnippetsExtractor) getSettingSections(node, nextNode *yaml.Node) (string, error) {
	// Extract YAML settings
	allNodes := make(map[string]*yaml.Node)

	for i := 0; i < len(nextNode.Content); i += 2 {
		allNodes[nextNode.Content[i].Value] = &yaml.Node{
			Kind:  yaml.MappingNode,
			Tag:   nextNode.Tag,
			Value: node.Value,
			Content: []*yaml.Node{
				{
					Kind:  yaml.ScalarNode,
					Value: node.Value,
					Tag:   node.Tag,
				},
				{
					Kind: yaml.MappingNode,
					Content: []*yaml.Node{
						{
							Kind:  yaml.ScalarNode,
							Value: "settings",
							Tag:   node.Tag,
						},
						{
							Kind:    yaml.MappingNode,
							Tag:     nextNode.Tag,
							Content: []*yaml.Node{nextNode.Content[i], nextNode.Content[i+1]},
						},
					},
				},
			},
		}
	}

	// Using linter information

	linters, err := readJSONFile[[]*types.LinterWrapper](filepath.Join(e.assetsPath, fmt.Sprintf("%s-info.json", node.Value)))
	if err != nil {
		return "", err
	}

	builder := &strings.Builder{}
	for _, lc := range linters {
		if lc.Internal {
			continue
		}

		// it's important to use lc.Name() nor name because name can be alias
		_, _ = fmt.Fprintf(builder, "## %s\n\n", lc.Name)

		writeTags(builder, lc)

		if lc.Deprecation != nil {
			continue
		}

		settings, ok := allNodes[lc.Name]
		if !ok {
			if hasSettings(lc.Name) {
				return "", fmt.Errorf("can't find %s settings in .golangci.reference.yml", lc.Name)
			}

			_, _ = fmt.Fprintln(builder, "_No configuration_")

			continue
		}

		_, _ = fmt.Fprintln(builder, "```yaml")

		encoder := yaml.NewEncoder(builder)
		encoder.SetIndent(2)

		err := encoder.Encode(settings)
		if err != nil {
			return "", err
		}

		_, _ = fmt.Fprintln(builder, "```")
		_, _ = fmt.Fprintln(builder)
	}

	return builder.String(), nil
}

func writeTags(builder *strings.Builder, lc *types.LinterWrapper) {
	if lc == nil {
		return
	}

	_, _ = fmt.Fprintf(builder, "%s\n\n", getDesc(lc))

	_, _ = fmt.Fprintln(builder, internal.NewBadge().
		Content(fmt.Sprintf("Since golangci-lint %s", lc.Since)).
		Icon("calendar"))

	switch {
	case lc.IsDeprecated():
		content := "Deprecated"
		if lc.Deprecation.Replacement != "" {
			content += fmt.Sprintf(" since %s", lc.Deprecation.Since)
		}

		_, _ = fmt.Fprintln(builder, internal.NewBadge().
			Content(content).
			Type("error").
			Icon("sparkles"))

	case lc.CanAutoFix:
		_, _ = fmt.Fprintln(builder, internal.NewBadge().
			Content("Autofix").
			Type("info").
			Icon("sparkles"))
	}

	if lc.OriginalURL != "" {
		_, _ = fmt.Fprintln(builder, internal.NewBadge().
			Content("Repository").
			Link(lc.OriginalURL).
			Icon("github"))
	}

	_, _ = fmt.Fprintln(builder)
}

func hasSettings(name string) bool {
	tp := reflect.TypeOf(config.LintersSettings{})

	for i := range tp.NumField() {
		if strings.EqualFold(name, tp.Field(i).Name) {
			return true
		}
	}

	tp = reflect.TypeOf(config.FormatterSettings{})

	for i := range tp.NumField() {
		if strings.EqualFold(name, tp.Field(i).Name) {
			return true
		}
	}

	return false
}
