package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/golangci/golangci-lint/pkg/config"
)

func (e *Executor) initVersionConfiguration(cmd *cobra.Command) {
	fs := cmd.Flags()
	fs.SortFlags = false // sort them as they are defined here
	initVersionFlagSet(fs, e.cfg)
}

func initVersionFlagSet(fs *pflag.FlagSet, cfg *config.Config) {
	// Version config
	vc := &cfg.Version
	fs.StringVar(&vc.Format, "format", "", wh("The version's format can be: 'short', 'json'"))
}

func (e *Executor) initVersion() {
	versionCmd := &cobra.Command{
		Use:               "version",
		Short:             "Version",
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE: func(cmd *cobra.Command, _ []string) error {
			switch strings.ToLower(e.cfg.Version.Format) {
			case "short":
				fmt.Println(e.buildInfo.Version)
				return nil

			case "json":
				return json.NewEncoder(os.Stdout).Encode(e.buildInfo)

			default:
				return printVersion(os.Stdout, e.buildInfo)
			}
		},
	}

	e.rootCmd.AddCommand(versionCmd)
	e.initVersionConfiguration(versionCmd)
}

func printVersion(w io.Writer, buildInfo BuildInfo) error {
	_, err := fmt.Fprintf(w, "golangci-lint has version %s built with %s from %s on %s\n",
		buildInfo.Version, buildInfo.GoVersion, buildInfo.Commit, buildInfo.Date)
	return err
}
