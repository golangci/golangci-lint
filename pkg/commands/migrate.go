package commands

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/pelletier/go-toml/v2"
	"github.com/santhosh-tekuri/jsonschema/v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"

	"github.com/golangci/golangci-lint/pkg/commands/internal/migrate"
	"github.com/golangci/golangci-lint/pkg/commands/internal/migrate/one"
	"github.com/golangci/golangci-lint/pkg/commands/internal/migrate/two"
	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/exitcodes"
	"github.com/golangci/golangci-lint/pkg/logutils"
)

type migrateOptions struct {
	config.LoaderOptions

	format         string // Flag only.
	skipValidation bool   // Flag only.
}
type migrateCommand struct {
	viper *viper.Viper
	cmd   *cobra.Command

	opts migrateOptions

	cfg *one.Config

	buildInfo BuildInfo

	log logutils.Log
}

func newMigrateCommand(log logutils.Log, info BuildInfo) *migrateCommand {
	c := &migrateCommand{
		viper:     viper.New(),
		cfg:       one.NewConfig(),
		buildInfo: info,
		log:       log,
	}

	migrateCmd := &cobra.Command{
		Use:               "migrate",
		Short:             "Migrate configuration file from v1 to v2",
		SilenceUsage:      true,
		SilenceErrors:     true,
		Args:              cobra.NoArgs,
		RunE:              c.execute,
		PreRunE:           c.preRunE,
		PersistentPreRunE: c.persistentPreRunE,
	}

	migrateCmd.SetOut(logutils.StdOut) // use custom output to properly color it in Windows terminals
	migrateCmd.SetErr(logutils.StdErr)

	fs := migrateCmd.Flags()
	fs.SortFlags = false // sort them as they are defined here

	setupConfigFileFlagSet(fs, &c.opts.LoaderOptions)

	fs.StringVar(&c.opts.format, "format", "",
		color.GreenString("By default, the file format is based on the configuration file extension.\n"+
			"Overrides file format detection.\nIt can be 'yml', 'yaml', 'toml', 'json'."))

	fs.BoolVar(&c.opts.skipValidation, "skip-validation", false,
		color.GreenString("Skip validating the configuration file on the JSONSchema of the v1."))

	c.cmd = migrateCmd

	return c
}

func (c *migrateCommand) execute(_ *cobra.Command, _ []string) error {
	if c.cfg.Version != "" {
		return fmt.Errorf("configuration version is already set: %s", c.cfg.Version)
	}

	srcPath := c.viper.ConfigFileUsed()
	if srcPath == "" {
		c.log.Warnf("No config file detected")
		os.Exit(exitcodes.NoConfigFileDetected)
	}

	err := c.backupConfigurationFile(srcPath)
	if err != nil {
		return err
	}

	c.cmd.Println("Migrating v1 configuration file:", srcPath)

	ext := filepath.Ext(srcPath)
	if strings.TrimSpace(c.opts.format) != "" {
		ext = "." + strings.TrimPrefix(c.opts.format, ".")
	}

	if !strings.EqualFold(filepath.Ext(srcPath), ext) {
		defer func() {
			_ = os.RemoveAll(srcPath)
		}()
	}

	newCfg := migrate.ToConfig(c.cfg)

	dstPath := strings.TrimSuffix(srcPath, filepath.Ext(srcPath)) + ext

	err = saveNewConfiguration(newCfg, dstPath)
	if err != nil {
		return fmt.Errorf("saving configuration file: %w", err)
	}

	c.cmd.Println("Migration done:", dstPath)

	return nil
}

func (c *migrateCommand) preRunE(cmd *cobra.Command, _ []string) error {
	if c.opts.skipValidation {
		return nil
	}

	usedConfigFile := c.viper.ConfigFileUsed()
	if usedConfigFile == "" {
		c.log.Warnf("No config file detected")
		os.Exit(exitcodes.NoConfigFileDetected)
	}

	c.cmd.Println("Validating v1 configuration file:", usedConfigFile)

	err := validateConfiguration("https://golangci-lint.run/jsonschema/golangci.v1.jsonschema.json", usedConfigFile)
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

func (c *migrateCommand) persistentPreRunE(_ *cobra.Command, args []string) error {
	c.log.Infof("%s", c.buildInfo.String())

	loader := config.NewBaseLoader(c.log.Child(logutils.DebugKeyConfigReader), c.viper, c.opts.LoaderOptions, c.cfg, args)

	err := loader.Load()
	if err != nil {
		return fmt.Errorf("can't load config: %w", err)
	}

	return nil
}

func (c *migrateCommand) backupConfigurationFile(srcPath string) error {
	filename := strings.TrimSuffix(filepath.Base(srcPath), filepath.Ext(srcPath)) + ".bck" + filepath.Ext(srcPath)
	dstPath := filepath.Join(filepath.Dir(srcPath), filename)

	c.cmd.Println("Saving the v1 configuration to:", dstPath)

	stat, err := os.Stat(srcPath)
	if err != nil {
		return err
	}

	data, err := os.ReadFile(srcPath)
	if err != nil {
		return err
	}

	err = os.WriteFile(dstPath, data, stat.Mode())
	if err != nil {
		return err
	}

	return nil
}

func saveNewConfiguration(newCfg *two.Config, dstPath string) error {
	dstFile, err := os.Create(dstPath)
	if err != nil {
		return err
	}

	defer func() { _ = dstFile.Close() }()

	ext := filepath.Ext(dstPath)

	switch strings.ToLower(ext) {
	case ".yml", ".yaml":
		encoder := yaml.NewEncoder(dstFile)
		encoder.SetIndent(2)

		return encoder.Encode(newCfg)

	case ".toml":
		encoder := toml.NewEncoder(dstFile)

		return encoder.Encode(newCfg)

	case ".json":
		// The JSON encoder converts empty struct to `{}` instead of nothing (even with omitempty JSON struct tags).
		// So we need to use the YAML encoder as bridge to create JSON file.

		var buf bytes.Buffer
		err := yaml.NewEncoder(&buf).Encode(newCfg)
		if err != nil {
			return err
		}

		raw := map[string]any{}
		err = yaml.NewDecoder(&buf).Decode(raw)
		if err != nil {
			return err
		}

		encoder := json.NewEncoder(dstFile)
		encoder.SetIndent("", "  ")

		return encoder.Encode(raw)

	default:
		return fmt.Errorf("unsupported file type: %s", ext)
	}
}
