package commands

import (
	"log"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/sirupsen/logrus"
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
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			runtime.GOMAXPROCS(e.cfg.Common.Concurrency)

			if e.cfg.Common.IsVerbose {
				logrus.SetLevel(logrus.InfoLevel)
			}

			if e.cfg.Common.CPUProfilePath != "" {
				f, err := os.Create(e.cfg.Common.CPUProfilePath)
				if err != nil {
					log.Fatal(err)
				}
				if err := pprof.StartCPUProfile(f); err != nil {
					log.Fatal(err)
				}
			}
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			if e.cfg.Common.CPUProfilePath != "" {
				pprof.StopCPUProfile()
			}
			os.Exit(e.exitCode)
		},
	}
	rootCmd.PersistentFlags().BoolVarP(&e.cfg.Common.IsVerbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringVar(&e.cfg.Common.CPUProfilePath, "cpu-profile-path", "", "Path to CPU profile output file")
	rootCmd.PersistentFlags().IntVarP(&e.cfg.Common.Concurrency, "concurrency", "j", runtime.NumCPU(), "Concurrency")
	e.rootCmd = rootCmd
}
