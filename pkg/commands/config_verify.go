package commands

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	hcversion "github.com/hashicorp/go-version"
	"github.com/pelletier/go-toml/v2"
	"github.com/santhosh-tekuri/jsonschema/v6"
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

		printValidationDetail(cmd, v.DetailedOutput())

		return errors.New("the configuration contains invalid elements")
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

		if version.Core().Equal(hcversion.Must(hcversion.NewVersion("v0.0.0"))) {
			commit, err := extractCommitHash(buildInfo)
			if err != nil {
				return "", err
			}

			return fmt.Sprintf("https://raw.githubusercontent.com/golangci/golangci-lint/%s/jsonschema/golangci.next.jsonschema.json",
				commit), nil
		}

		return fmt.Sprintf("https://golangci-lint.run/jsonschema/golangci.v%d.%d.jsonschema.json",
			version.Segments()[0], version.Segments()[1]), nil

	case buildInfo.Commit != "" && buildInfo.Commit != "?":
		commit, err := extractCommitHash(buildInfo)
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("https://raw.githubusercontent.com/golangci/golangci-lint/%s/jsonschema/golangci.next.jsonschema.json",
			commit), nil

	default:
		return "", errors.New("version not found")
	}
}

func extractCommitHash(buildInfo BuildInfo) (string, error) {
	if buildInfo.Commit == "" || buildInfo.Commit == "?" {
		return "", errors.New("empty commit information")
	}

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

	if commit == "unknown" {
		return "", errors.New("unknown commit information")
	}

	return commit, nil
}

func validateConfiguration(schemaPath, targetFile string) error {
	compiler := jsonschema.NewCompiler()
	compiler.UseLoader(jsonschema.SchemeURLLoader{
		"file":  jsonschema.FileLoader{},
		"https": newJSONSchemaHTTPLoader(),
	})
	compiler.DefaultDraft(jsonschema.Draft7)

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

func printValidationDetail(cmd *cobra.Command, detail *jsonschema.OutputUnit) {
	if detail.Error != nil {
		data, _ := json.Marshal(detail.Error)
		details, _ := strconv.Unquote(string(data))

		cmd.PrintErrf("jsonschema: %q does not validate with %q: %s\n",
			strings.ReplaceAll(strings.TrimPrefix(detail.InstanceLocation, "/"), "/", "."), detail.KeywordLocation, details)
	}

	for _, d := range detail.Errors {
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

type jsonschemaHTTPLoader struct {
	*http.Client
}

func newJSONSchemaHTTPLoader() *jsonschemaHTTPLoader {
	return &jsonschemaHTTPLoader{Client: &http.Client{
		Timeout: 2 * time.Second,
	}}
}

func (l jsonschemaHTTPLoader) Load(url string) (any, error) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, err
	}

	resp, err := l.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s returned status code %d", url, resp.StatusCode)
	}

	return jsonschema.UnmarshalJSON(resp.Body)
}
