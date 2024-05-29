package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type BuildInfo struct {
	GoVersion string `json:"goVersion"`
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	Date      string `json:"date"`
}

func (b BuildInfo) String() string {
	return fmt.Sprintf("golangci-lint has version %s built with %s from %s on %s",
		b.Version, b.GoVersion, b.Commit, b.Date)
}

type versionInfo struct {
	Info      BuildInfo
	BuildInfo *debug.BuildInfo
}

type versionOptions struct {
	Format string
	Debug  bool
}

type versionCommand struct {
	cmd  *cobra.Command
	opts versionOptions

	info BuildInfo
}

func newVersionCommand(info BuildInfo) *versionCommand {
	c := &versionCommand{info: info}

	versionCmd := &cobra.Command{
		Use:               "version",
		Short:             "Version",
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE:              c.execute,
	}

	fs := versionCmd.Flags()
	fs.SortFlags = false // sort them as they are defined here

	fs.StringVar(&c.opts.Format, "format", "", color.GreenString("The version's format can be: 'short', 'json'"))
	fs.BoolVar(&c.opts.Debug, "debug", false, color.GreenString("Add build information"))

	c.cmd = versionCmd

	return c
}

func (c *versionCommand) execute(_ *cobra.Command, _ []string) error {
	if c.opts.Debug {
		info, ok := debug.ReadBuildInfo()
		if !ok {
			return nil
		}

		switch strings.ToLower(c.opts.Format) {
		case "json":
			return json.NewEncoder(os.Stdout).Encode(versionInfo{
				Info:      c.info,
				BuildInfo: info,
			})

		default:
			fmt.Println(info.String())
			return printVersion(os.Stdout, c.info)
		}
	}

	switch strings.ToLower(c.opts.Format) {
	case "short":
		fmt.Println(c.info.Version)
		return nil

	case "json":
		return json.NewEncoder(os.Stdout).Encode(c.info)

	default:
		return printVersion(os.Stdout, c.info)
	}
}

func printVersion(w io.Writer, info BuildInfo) error {
	_, err := fmt.Fprintln(w, info.String())
	return err
}
