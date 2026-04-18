package commands

import (
	"bytes"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/spf13/cobra"
	"go.yaml.in/yaml/v3"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goformatters"
	"github.com/golangci/golangci-lint/v2/pkg/lint/lintersdb"
	"github.com/golangci/golangci-lint/v2/pkg/logutils"
)

func (c *configCommand) executeUpdate(cmd *cobra.Command, _ []string) error {
	configFilePath := c.getUsedConfig()

	var cfgNode yaml.Node

	if configFilePath != "" {
		data, err := os.ReadFile(configFilePath)
		if err != nil {
			return fmt.Errorf("reading config file: %w", err)
		}
		if err = yaml.Unmarshal(data, &cfgNode); err != nil {
			return fmt.Errorf("parsing config file: %w", err)
		}
	} else {
		configFilePath = ".golangci.yml"
		cfgNode = newEmptyConfig()
		c.log.Infof("No config file found, creating %s", configFilePath)
	}

	// Build linter database.
	manager, err := lintersdb.NewManager(
		c.log.Child(logutils.DebugKeyLintersDB),
		config.NewDefault(),
		lintersdb.NewLinterBuilder(),
	)
	if err != nil {
		return fmt.Errorf("creating linter manager: %w", err)
	}

	update := configUpdate{}
	update.load(manager)

	// Get the root mapping node.
	if cfgNode.Kind != yaml.DocumentNode || len(cfgNode.Content) == 0 || cfgNode.Content[0].Kind != yaml.MappingNode {
		return fmt.Errorf("invalid YAML document structure")
	}
	rootMap := cfgNode.Content[0]

	// Ensure the linters section exists and update it with missing linters.
	update.addMissingLintersComments(rootMap)

	// Add missing formatter names as comments.
	update.addMissingFormatterComments(rootMap)

	// Marshal and write the updated config.
	var buf bytes.Buffer

	encoder := yaml.NewEncoder(&buf)
	encoder.SetIndent(2)

	if err = encoder.Encode(&cfgNode); err != nil {
		return fmt.Errorf("encoding YAML: %w", err)
	}

	if err = encoder.Close(); err != nil {
		return fmt.Errorf("closing encoder: %w", err)
	}

	if err = os.WriteFile(configFilePath, buf.Bytes(), 0o644); err != nil {
		return fmt.Errorf("writing config file: %w", err)
	}

	cmd.Printf("Configuration updated: %s\n", configFilePath)

	return nil
}

func newEmptyConfig() yaml.Node {
	return yaml.Node{
		Kind: yaml.DocumentNode,
		Content: []*yaml.Node{
			{
				Kind: yaml.MappingNode,
				Content: []*yaml.Node{
					{Kind: yaml.ScalarNode, Value: "version", Tag: "!!str"},
					{Kind: yaml.ScalarNode, Value: "2", Tag: "!!str", Style: yaml.DoubleQuotedStyle},
				},
			},
		},
	}
}

type configUpdate struct {
	supportedLinters    []string
	supportedFormatters []string
	descriptions        map[string]string
}

func (u *configUpdate) load(manager *lintersdb.Manager) {
	u.descriptions = map[string]string{}
	for _, lc := range manager.GetAllSupportedLinterConfigs() {
		if !lc.Internal && !lc.IsDeprecated() {
			name := lc.Name()
			u.descriptions[name], _, _ = strings.Cut(lc.Linter.Desc(), "\n")
			if goformatters.IsFormatter(name) {
				u.supportedFormatters = append(u.supportedFormatters, name)
			} else {
				u.supportedLinters = append(u.supportedLinters, name)
			}
		}
	}
}

func (u *configUpdate) addMissingLineComments(item *yaml.Node) {
	if item.LineComment == "" {
		item.LineComment = u.descriptions[item.Value]
	}
}

// addMissingLintersComments adds commented-out entries for linters
// that are not yet in the enable/disable lists.
// The allLinters slice must already be filtered (no internal, no deprecated, no formatters).
func (u *configUpdate) addMissingLintersComments(rootMap *yaml.Node) {
	existingLinters := make(map[string]bool)

	checkNode := rootMap
	editNode, header, indent := rootMap, "Linter configuration.\nlinters:\n  enable:", "    "

	lintersVal := findMappingValue(rootMap, "linters")
	if lintersVal != nil && lintersVal.Kind == yaml.MappingNode {
		checkNode = lintersVal
		editNode, header, indent = lintersVal, "Enable specific linter.\nenable:", "  "

		enableVal := findMappingValue(lintersVal, "enable")
		if enableVal != nil && enableVal.Kind == yaml.SequenceNode {
			editNode, header, indent = enableVal, "New linters", ""
			for _, item := range enableVal.Content {
				if item.Kind == yaml.ScalarNode {
					existingLinters[item.Value] = true
					u.addMissingLineComments(item)
				}
			}
		}

		disableVal := findMappingValue(lintersVal, "disable")
		if disableVal != nil && disableVal.Kind == yaml.SequenceNode {
			for _, item := range disableVal.Content {
				if item.Kind == yaml.ScalarNode {
					existingLinters[item.Value] = true
					u.addMissingLineComments(item)
				}
			}
		}
	}

	// Build list of missing linter comments.
	commentLines := []string{header}
	for _, name := range u.supportedLinters {
		if !existingLinters[name] && !textInNodeComments(checkNode, name) {
			commentLines = append(commentLines, fmt.Sprintf("%s- %s  # %s", indent, name, u.descriptions[name]))
		}
	}
	if len(commentLines) > 1 {
		appendToFootComment(editNode, strings.Join(commentLines, "\n"))
	}
}

// addMissingFormatterComments adds commented-out entries for formatters not yet enabled.
// The allFormatters slice must already be filtered to only contain formatter configs.
func (u *configUpdate) addMissingFormatterComments(rootMap *yaml.Node) {
	existingFormatters := make(map[string]bool)

	checkNode := rootMap
	editNode, header, indent := rootMap, "Formatters configuration.\nformatters:\n  enable:", "    "

	formattersVal := findMappingValue(rootMap, "formatters")
	if formattersVal != nil && formattersVal.Kind == yaml.MappingNode {
		checkNode = formattersVal
		editNode, header, indent = formattersVal, "Enable specific formatter.\nenable:", "  "

		enableVal := findMappingValue(formattersVal, "enable")
		if enableVal != nil && enableVal.Kind == yaml.SequenceNode {
			editNode, header, indent = enableVal, "New formatters", ""

			for _, item := range enableVal.Content {
				if item.Kind == yaml.ScalarNode {
					existingFormatters[item.Value] = true
					u.addMissingLineComments(item)
				}
			}
		}
	}

	commentLines := []string{header}
	for _, name := range u.supportedFormatters {
		if !existingFormatters[name] && !textInNodeComments(checkNode, name) {
			commentLines = append(commentLines, fmt.Sprintf("%s- %s  # %s", indent, name, u.descriptions[name]))
		}
	}
	if len(commentLines) > 1 {
		appendToFootComment(editNode, strings.Join(commentLines, "\n"))
	}
}

// --- YAML node helper functions ---

// findMappingValue returns the value node for a given key in a mapping node.
// Returns nil if the key is not found.
func findMappingValue(mapping *yaml.Node, key string) *yaml.Node {
	if mapping == nil || mapping.Kind != yaml.MappingNode {
		return nil
	}

	for i := 0; i < len(mapping.Content)-1; i += 2 {
		if mapping.Content[i].Kind == yaml.ScalarNode && mapping.Content[i].Value == key {
			return mapping.Content[i+1]
		}
	}

	return nil
}

// textInNodeComments checks recursively whether the given text appears in any comment of the node tree.s
func textInNodeComments(node *yaml.Node, text string) bool {
	if node == nil {
		return false
	}

	if strings.Contains(node.HeadComment, text) ||
		strings.Contains(node.LineComment, text) ||
		strings.Contains(node.FootComment, text) {
		return true
	}

	return slices.ContainsFunc(node.Content, func(child *yaml.Node) bool {
		return textInNodeComments(child, text)
	})
}

// appendToFootComment appends additional comment text to a node's FootComment.
func appendToFootComment(node *yaml.Node, comment string) {
	if node == nil {
		return
	}

	if (node.Kind == yaml.SequenceNode || node.Kind == yaml.MappingNode) && len(node.Content) > 0 {
		node = node.Content[len(node.Content)-1]
	}

	if node.FootComment == "" {
		node.FootComment = comment
	} else {
		node.FootComment += "\n\n" + comment
	}
}
