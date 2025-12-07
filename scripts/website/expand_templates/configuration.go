package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"go.yaml.in/yaml/v3"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/scripts/website/types"
)

const (
	keyLinters    = "linters"
	keyFormatters = "formatters"
	keySettings   = "settings"
)

type SettingSnippets struct {
	ConfigurationFile  map[string]string
	LintersSettings    map[string]string
	FormattersSettings map[string]string
}

type ExampleSnippetsExtractor struct {
	referencePath string
	assetsPath    string
}

func NewExampleSnippetsExtractor() *ExampleSnippetsExtractor {
	return &ExampleSnippetsExtractor{
		referencePath: ".golangci.reference.yml",
		assetsPath:    filepath.Join("docs", "data"),
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

	snippets := SettingSnippets{
		ConfigurationFile: make(map[string]string),
	}

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

		buf := bytes.NewBuffer(nil)
		encoder := yaml.NewEncoder(buf)
		encoder.SetIndent(2)

		err := encoder.Encode(nodeSection)
		if err != nil {
			return nil, err
		}

		snippets.ConfigurationFile[node.Value] = buf.String()
	}

	buf := bytes.NewBuffer(nil)
	encoder := yaml.NewEncoder(buf)
	encoder.SetIndent(2)

	err := encoder.Encode(globalNode)
	if err != nil {
		return nil, err
	}

	snippets.ConfigurationFile["root"] = buf.String()

	return &snippets, nil
}

func (e *ExampleSnippetsExtractor) getSettingSections(node, nextNode *yaml.Node) (map[string]string, error) {
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

	linterSettings := make(map[string]string)

	// Using linter information
	linters, err := readJSONFile[[]*types.LinterWrapper](filepath.Join(e.assetsPath, fmt.Sprintf("%s_info.json", node.Value)))
	if err != nil {
		return nil, err
	}

	for _, lc := range linters {
		if lc.Internal {
			continue
		}

		if lc.Deprecation != nil {
			continue
		}

		settings, ok := allNodes[lc.Name]
		if !ok {
			if hasSettings(lc.Name) {
				log.Printf("can't find %s settings in .golangci.reference.yml", lc.Name)
			}

			continue
		}

		buf := bytes.NewBuffer(nil)
		encoder := yaml.NewEncoder(buf)
		encoder.SetIndent(2)

		err := encoder.Encode(settings)
		if err != nil {
			return nil, err
		}

		linterSettings[lc.Name] = buf.String()
	}

	return linterSettings, nil
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
