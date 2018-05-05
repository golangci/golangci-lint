package commands

import (
	"github.com/spf13/cobra"
)

func (e *Executor) initRoot() {
	rootCmd := &cobra.Command{
		Use:   "golangci-lint",
		Short: "golangci-lint is a smart linters runner.",
		Long:  `Smart, fast linters runner. Run it in cloud for every GitHub pull request on https://golangci.com`,
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}
	rootCmd.PersistentFlags().BoolVarP(&e.cfg.Common.IsVerbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringVar(&e.cfg.Common.CPUProfilePath, "cpu-profile-path", "", "Path to CPU profile output file")
	e.rootCmd = rootCmd
}
