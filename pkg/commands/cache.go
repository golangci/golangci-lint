package commands

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/golangci/golangci-lint/internal/cache"
	"github.com/golangci/golangci-lint/pkg/config"
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

	cacheCmd.AddCommand(&cobra.Command{
		Use:               "clean",
		Short:             "Clean cache",
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE:              e.executeCacheClean,
	})
	cacheCmd.AddCommand(&cobra.Command{
		Use:               "status",
		Short:             "Show cache status",
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		Run:               e.executeCacheStatus,
	})

	// TODO: add trim command?

	e.rootCmd.AddCommand(cacheCmd)
}

func (e *Executor) executeCacheClean(_ *cobra.Command, _ []string) error {
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
			size += info.Size()
		}
		return err
	})
	return size, err
}

// --- Related to cache but not used directly by the cache command.

func initHashSalt(version string, cfg *config.Config) error {
	binSalt, err := computeBinarySalt(version)
	if err != nil {
		return fmt.Errorf("failed to calculate binary salt: %w", err)
	}

	configSalt, err := computeConfigSalt(cfg)
	if err != nil {
		return fmt.Errorf("failed to calculate config salt: %w", err)
	}

	b := bytes.NewBuffer(binSalt)
	b.Write(configSalt)
	cache.SetSalt(b.Bytes())
	return nil
}

func computeBinarySalt(version string) ([]byte, error) {
	if version != "" && version != "(devel)" {
		return []byte(version), nil
	}

	if logutils.HaveDebugTag(logutils.DebugKeyBinSalt) {
		return []byte("debug"), nil
	}

	p, err := os.Executable()
	if err != nil {
		return nil, err
	}
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

// computeConfigSalt computes configuration hash.
// We don't hash all config fields to reduce meaningless cache invalidations.
// At least, it has a huge impact on tests speed.
// Fields: `LintersSettings` and `Run.BuildTags`.
func computeConfigSalt(cfg *config.Config) ([]byte, error) {
	lintersSettingsBytes, err := yaml.Marshal(cfg.LintersSettings)
	if err != nil {
		return nil, fmt.Errorf("failed to json marshal config linter settings: %w", err)
	}

	configData := bytes.NewBufferString("linters-settings=")
	configData.Write(lintersSettingsBytes)
	configData.WriteString("\nbuild-tags=%s" + strings.Join(cfg.Run.BuildTags, ","))

	h := sha256.New()
	if _, err := h.Write(configData.Bytes()); err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}
