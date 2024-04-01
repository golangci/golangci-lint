package commands

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	hcversion "github.com/hashicorp/go-version"
	"github.com/pelletier/go-toml/v2"
	"github.com/santhosh-tekuri/jsonschema/v5"
	"github.com/santhosh-tekuri/jsonschema/v5/httploader"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"

	"github.com/golangci/golangci-lint/pkg/exitcodes"
)

type verifyOptions struct {
	schemaURL string // For debugging purpose only (Flag only).
}

func (c *configCommand) executeVerify(cmd *cobra.Command, _ []string) error {
	usedConfigFile := c.getUsedConfig()
	if usedConfigFile == "" {
		c.log.Warnf("No config file detected")
		os.Exit(exitcodes.NoConfigFileDetected)
	}

	schemaURL, err := createSchemaURL(cmd.Flags(), c.buildInfo)
	if err != nil {
		return fmt.Errorf("get JSON schema: %w", err)
	}

	err = validateConfiguration(schemaURL, usedConfigFile)
	if err != nil {
		var v *jsonschema.ValidationError
		if !errors.As(err, &v) {
			return fmt.Errorf("[%s] validate: %w", usedConfigFile, err)
		}

		detail := v.DetailedOutput()

		printValidationDetail(cmd, &detail)

		return fmt.Errorf("the configuration contains invalid elements")
	}

	return nil
}

func createSchemaURL(flags *pflag.FlagSet, buildInfo BuildInfo) (string, error) {
	schemaURL, err := flags.GetString("schema")
	if err != nil {
		return "", fmt.Errorf("get schema flag: %w", err)
	}

	if schemaURL != "" {
		return schemaURL, nil
	}

	switch {
	case buildInfo.Version != "" && buildInfo.Version != "(devel)":
		version, err := hcversion.NewVersion(buildInfo.Version)
		if err != nil {
			return "", fmt.Errorf("parse version: %w", err)
		}

		schemaURL = fmt.Sprintf("https://golangci-lint.run/jsonschema/golangci.v%d.%d.jsonschema.json",
			version.Segments()[0], version.Segments()[1])

	case buildInfo.Commit != "" && buildInfo.Commit != "?":
		if buildInfo.Commit == "unknown" {
			return "", errors.New("unknown commit information")
		}

		commit := buildInfo.Commit

		if strings.HasPrefix(commit, "(") {
			c, _, ok := strings.Cut(strings.TrimPrefix(commit, "("), ",")
			if !ok {
				return "", errors.New("commit information not found")
			}

			commit = c
		}

		schemaURL = fmt.Sprintf("https://raw.githubusercontent.com/golangci/golangci-lint/%s/jsonschema/golangci.next.jsonschema.json",
			commit)

	default:
		return "", errors.New("version not found")
	}

	return schemaURL, nil
}

func validateConfiguration(schemaPath, targetFile string) error {
	httploader.Client = &http.Client{Timeout: 2 * time.Second}

	compiler := jsonschema.NewCompiler()
	compiler.Draft = jsonschema.Draft7

	schema, err := compiler.Compile(schemaPath)
	if err != nil {
		return fmt.Errorf("compile schema: %w", err)
	}

	var m any

	switch strings.ToLower(filepath.Ext(targetFile)) {
	case ".yaml", ".yml", ".json":
		m, err = decodeYamlFile(targetFile)
		if err != nil {
			return err
		}

	case ".toml":
		m, err = decodeTomlFile(targetFile)
		if err != nil {
			return err
		}

	default:
		// unsupported
		return errors.New("unsupported configuration format")
	}

	return schema.Validate(m)
}

func printValidationDetail(cmd *cobra.Command, detail *jsonschema.Detailed) {
	if detail.Error != "" {
		cmd.PrintErrf("jsonschema: %q does not validate with %q: %s\n",
			strings.ReplaceAll(strings.TrimPrefix(detail.InstanceLocation, "/"), "/", "."), detail.KeywordLocation, detail.Error)
	}

	for _, d := range detail.Errors {
		d := d
		printValidationDetail(cmd, &d)
	}
}

func decodeYamlFile(filename string) (any, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("[%s] file open: %w", filename, err)
	}

	defer func() { _ = file.Close() }()

	var m any
	err = yaml.NewDecoder(file).Decode(&m)
	if err != nil {
		return nil, fmt.Errorf("[%s] YAML decode: %w", filename, err)
	}

	return m, nil
}

func decodeTomlFile(filename string) (any, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("[%s] file open: %w", filename, err)
	}

	defer func() { _ = file.Close() }()

	var m any
	err = toml.NewDecoder(file).Decode(&m)
	if err != nil {
		return nil, fmt.Errorf("[%s] TOML decode: %w", filename, err)
	}

	return m, nil
}
