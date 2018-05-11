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
			if err := cmd.Help(); err != nil {
				log.Fatal(err)
			}
		},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			runtime.GOMAXPROCS(e.cfg.Run.Concurrency)

			if e.cfg.Run.IsVerbose {
				logrus.SetLevel(logrus.InfoLevel)
			}

			if e.cfg.Run.CPUProfilePath != "" {
				f, err := os.Create(e.cfg.Run.CPUProfilePath)
				if err != nil {
					log.Fatal(err)
				}
				if err := pprof.StartCPUProfile(f); err != nil {
					log.Fatal(err)
				}
			}
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			if e.cfg.Run.CPUProfilePath != "" {
				pprof.StopCPUProfile()
			}
			os.Exit(e.exitCode)
		},
	}
	rootCmd.PersistentFlags().BoolVarP(&e.cfg.Run.IsVerbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringVar(&e.cfg.Run.CPUProfilePath, "cpu-profile-path", "", "Path to CPU profile output file")
	rootCmd.PersistentFlags().IntVarP(&e.cfg.Run.Concurrency, "concurrency", "j", runtime.NumCPU(), "Concurrency")

	e.rootCmd = rootCmd
}
