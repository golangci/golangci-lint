package config

import (
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

// UpdateConfigFileWithNewLinters adds new linters to the "linters" config in the file at the provided path
func UpdateConfigFileWithNewLinters(configFilePath string, newLinters []string) error {
	configData, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return errors.Wrap(err, "could not read config file")
	}

	var docNode yaml.Node
	if err := yaml.Unmarshal(configData, &docNode); err != nil {
		return errors.Wrapf(err, "failed to unmarshal config file %q", configFilePath)
	}

	var configNode *yaml.Node
	if len(docNode.Content) > 0 {
		configNode = docNode.Content[0]
	} else {
		configNode = &yaml.Node{Kind: yaml.MappingNode}
		docNode.Content = append(docNode.Content, configNode)
	}

	// guess the indent level by looking at the column of second level nodes
	indentSpaces := 2
GuessSpaces:
	for _, n := range configNode.Content {
		for _, nn := range n.Content {
			indentSpaces = nn.Column - 1
			break GuessSpaces
		}
	}

	lintersNode := findOrInsertKeyedValue(configNode, "linters", &yaml.Node{Kind: yaml.MappingNode})

	// find the "linters" -> "enable" node (or create it)
	enableNode := findOrInsertKeyedValue(lintersNode, "enable", &yaml.Node{Kind: yaml.SequenceNode})

	for _, l := range newLinters {
		node := &yaml.Node{}
		node.SetString(l)
		enableNode.Content = append(enableNode.Content, node)
	}

	configFile, err := os.OpenFile(configFilePath, os.O_WRONLY|os.O_TRUNC, 0)
	if err != nil {
		return errors.Wrapf(err, "failed to open file %q for writing", configFilePath)
	}

	encoder := yaml.NewEncoder(configFile)
	encoder.SetIndent(indentSpaces)
	err = encoder.Encode(docNode.Content[0])
	if err == nil {
		err = encoder.Close()
	}
	if err != nil {
		err = configFile.Close()
	}
	return errors.Wrapf(err, "failed to update config file %q", configFilePath)
}

func findOrInsertKeyedValue(node *yaml.Node, key string, value *yaml.Node) *yaml.Node {
	for i, n := range node.Content {
		var childKey string
		err := n.Decode(&childKey)
		if err == nil && key == childKey {
			return node.Content[i+1]
		}
	}
	keyNode := &yaml.Node{}
	keyNode.SetString(key)
	node.Content = append(node.Content, keyNode, value)
	return value
}
