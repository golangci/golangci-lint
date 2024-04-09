package test

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/santhosh-tekuri/jsonschema/v5"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func Test_validateTestConfigurationFiles(t *testing.T) {
	err := validateTestConfigurationFiles("../jsonschema/golangci.next.jsonschema.json", ".")
	require.NoError(t, err)
}

func Test_validateTestConfigurationFilesLinters(t *testing.T) {
	err := validateTestConfigurationFiles("../jsonschema/golangci.next.jsonschema.json", "../pkg/golinters")
	require.NoError(t, err)
}

func validateTestConfigurationFiles(schemaPath, targetDir string) error {
	schema, err := loadSchema(filepath.FromSlash(schemaPath))
	if err != nil {
		return fmt.Errorf("load schema: %w", err)
	}

	yamlFiles, err := findConfigurationFiles(filepath.FromSlash(targetDir))
	if err != nil {
		return fmt.Errorf("find configuration files: %w", err)
	}

	var errAll error
	for _, filename := range yamlFiles {
		// internal tests
		if filename == filepath.FromSlash("testdata/withconfig/.golangci.yml") {
			continue
		}

		m, err := decodeYamlFile(filename)
		if err != nil {
			return err
		}

		err = schema.Validate(m)
		if err != nil {
			abs, _ := filepath.Abs(filename)
			errAll = errors.Join(errAll, fmt.Errorf("%s: %w", abs, err))
		}
	}

	return errAll
}

func loadSchema(schemaPath string) (*jsonschema.Schema, error) {
	compiler := jsonschema.NewCompiler()
	compiler.Draft = jsonschema.Draft7

	schemaFile, err := os.Open(schemaPath)
	if err != nil {
		return nil, fmt.Errorf("open schema file: %w", err)
	}

	defer func() { _ = schemaFile.Close() }()

	err = compiler.AddResource(filepath.Base(schemaPath), schemaFile)
	if err != nil {
		return nil, fmt.Errorf("add schema resource: %w", err)
	}

	schema, err := compiler.Compile(filepath.Base(schemaPath))
	if err != nil {
		return nil, fmt.Errorf("compile schema: %w", err)
	}

	return schema, nil
}

func findConfigurationFiles(targetDir string) ([]string, error) {
	var yamlFiles []string

	err := filepath.WalkDir(targetDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && (strings.EqualFold(filepath.Ext(d.Name()), ".yml") || strings.EqualFold(filepath.Ext(d.Name()), ".yaml")) {
			yamlFiles = append(yamlFiles, path)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walk dir: %w", err)
	}

	return yamlFiles, nil
}

func decodeYamlFile(filename string) (any, error) {
	yamlFile, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("[%s] file open: %w", filename, err)
	}

	defer func() { _ = yamlFile.Close() }()

	var m any
	err = yaml.NewDecoder(yamlFile).Decode(&m)
	if err != nil {
		return nil, fmt.Errorf("[%s] yaml decode: %w", filename, err)
	}

	return m, nil
}
