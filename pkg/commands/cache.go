package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/golangci/golangci-lint/internal/cache"
	"github.com/golangci/golangci-lint/pkg/fsutils"
	"github.com/golangci/golangci-lint/pkg/logutils"
)

func (e *Executor) initCache() {
	cacheCmd := &cobra.Command{
		Use:   "cache",
		Short: "Cache control and information",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}
	e.rootCmd.AddCommand(cacheCmd)

	cacheCmd.AddCommand(&cobra.Command{
		Use:               "clean",
		Short:             "Clean cache",
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE:              e.executeCleanCache,
	})
	cacheCmd.AddCommand(&cobra.Command{
		Use:               "status",
		Short:             "Show cache status",
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		Run:               e.executeCacheStatus,
	})

	// TODO: add trim command?
}

func (e *Executor) executeCleanCache(_ *cobra.Command, _ []string) error {
	cacheDir := cache.DefaultDir()
	if err := os.RemoveAll(cacheDir); err != nil {
		return fmt.Errorf("failed to remove dir %s: %w", cacheDir, err)
	}

	return nil
}

func (e *Executor) executeCacheStatus(_ *cobra.Command, _ []string) {
	cacheDir := cache.DefaultDir()
	fmt.Fprintf(logutils.StdOut, "Dir: %s\n", cacheDir)

	cacheSizeBytes, err := dirSizeBytes(cacheDir)
	if err == nil {
		fmt.Fprintf(logutils.StdOut, "Size: %s\n", fsutils.PrettifyBytesCount(cacheSizeBytes))
	}
}

func dirSizeBytes(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			size = info.Size()
		}
		return err
	})
	return size, err
}
