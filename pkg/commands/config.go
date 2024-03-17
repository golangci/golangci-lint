package commands

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	hcversion "github.com/hashicorp/go-version"
	"github.com/pelletier/go-toml/v2"
	"github.com/santhosh-tekuri/jsonschema/v5"
	_ "github.com/santhosh-tekuri/jsonschema/v5/httploader"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/exitcodes"
	"github.com/golangci/golangci-lint/pkg/fsutils"
	"github.com/golangci/golangci-lint/pkg/logutils"
)

type verifyOptions struct {
	schemaURL string // For debugging purpose only (Flag only).
}

type configCommand struct {
	viper *viper.Viper
	cmd   *cobra.Command

	opts       config.LoaderOptions
	verifyOpts verifyOptions

	buildInfo BuildInfo

	log logutils.Log
}

func newConfigCommand(log logutils.Log, info BuildInfo) *configCommand {
	c := &configCommand{
		viper:     viper.New(),
		log:       log,
		buildInfo: info,
	}

	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Config file information",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
		PersistentPreRunE: c.preRunE,
	}

	verifyCommand := &cobra.Command{
		Use:               "verify",
		Short:             "Verify configuration against JSON schema",
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE:              c.executeVerify,
	}

	configCmd.AddCommand(
		&cobra.Command{
			Use:               "path",
			Short:             "Print used config path",
			Args:              cobra.NoArgs,
			ValidArgsFunction: cobra.NoFileCompletions,
			Run:               c.execute,
		},
		verifyCommand,
	)

	flagSet := configCmd.PersistentFlags()
	flagSet.SortFlags = false // sort them as they are defined here

	setupConfigFileFlagSet(flagSet, &c.opts)

	// ex: --schema jsonschema/golangci.next.jsonschema.json
	verifyFlagSet := verifyCommand.Flags()
	verifyFlagSet.StringVar(&c.verifyOpts.schemaURL, "schema", "", color.GreenString("JSON schema URL"))
	_ = verifyFlagSet.MarkHidden("schema")

	c.cmd = configCmd

	return c
}

func (c *configCommand) preRunE(cmd *cobra.Command, _ []string) error {
	// The command doesn't depend on the real configuration.
	// It only needs to know the path of the configuration file.
	loader := config.NewLoader(c.log.Child(logutils.DebugKeyConfigReader), c.viper, cmd.Flags(), c.opts, config.NewDefault())

	if err := loader.Load(); err != nil {
		return fmt.Errorf("can't load config: %w", err)
	}

	return nil
}

func (c *configCommand) execute(cmd *cobra.Command, _ []string) {
	usedConfigFile := c.getUsedConfig()
	if usedConfigFile == "" {
		c.log.Warnf("No config file detected")
		os.Exit(exitcodes.NoConfigFileDetected)
	}

	cmd.Println(usedConfigFile)
}

// getUsedConfig returns the resolved path to the golangci config file,
// or the empty string if no configuration could be found.
func (c *configCommand) getUsedConfig() string {
	usedConfigFile := c.viper.ConfigFileUsed()
	if usedConfigFile == "" {
		return ""
	}

	prettyUsedConfigFile, err := fsutils.ShortestRelPath(usedConfigFile, "")
	if err != nil {
		c.log.Warnf("Can't pretty print config file path: %s", err)
		return usedConfigFile
	}

	return prettyUsedConfigFile
}

func (c *configCommand) executeVerify(cmd *cobra.Command, _ []string) error {
	usedConfigFile := c.getUsedConfig()
	if usedConfigFile == "" {
		c.log.Warnf("No config file detected")
		os.Exit(exitcodes.NoConfigFileDetected)
	}

	schemaURL, err := getSchemaURL(cmd.Flags(), c.buildInfo)
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
	}

	return nil
}

func printValidationDetail(cmd *cobra.Command, detail *jsonschema.Detailed) {
	if detail.Error != "" {
		cmd.PrintErrf("jsonschema: %s does not validate with %s: %s\n",
			strings.ReplaceAll(strings.TrimPrefix(detail.InstanceLocation, "/"), "/", "."), detail.KeywordLocation, detail.Error)
	}

	for _, d := range detail.Errors {
		d := d
		printValidationDetail(cmd, &d)
	}
}

func getSchemaURL(flags *pflag.FlagSet, buildInfo BuildInfo) (string, error) {
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
		if buildInfo.Commit != "unknown" {
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

	err = schema.Validate(m)
	if err != nil {
		return err
	}

	return nil
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
