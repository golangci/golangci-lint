package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/lint/lintersdb"
	"github.com/golangci/golangci-lint/scripts/website/types"
)

func main() {
	linters := lintersdb.NewLinterBuilder().Build(config.NewDefault())

	var wraps []types.LinterWrapper
	for _, l := range linters {
		wrapper := types.LinterWrapper{
			Name:             l.Linter.Name(),
			Desc:             l.Linter.Desc(),
			EnabledByDefault: l.EnabledByDefault,
			LoadMode:         l.LoadMode,
			InPresets:        l.InPresets,
			AlternativeNames: l.AlternativeNames,
			OriginalURL:      l.OriginalURL,
			Internal:         l.Internal,
			CanAutoFix:       l.CanAutoFix,
			IsSlow:           l.IsSlow,
			DoesChangeTypes:  l.DoesChangeTypes,
			Since:            l.Since,
		}

		if l.Deprecation != nil {
			wrapper.Deprecation = &types.Deprecation{
				Since:       l.Deprecation.Since,
				Message:     l.Deprecation.Message,
				Replacement: l.Deprecation.Replacement,
			}
		}

		wraps = append(wraps, wrapper)
	}

	err := saveToJSONFile(filepath.Join("assets", "linters-info.json"), wraps)
	if err != nil {
		log.Fatalf("Save linters: %v", err)
	}

	err = saveToJSONFile(filepath.Join("assets", "default-exclusions.json"), config.DefaultExcludePatterns)
	if err != nil {
		log.Fatalf("Save default exclusions: %v", err)
	}

	err = saveCLIHelp(filepath.Join("assets", "cli-help.json"))
	if err != nil {
		log.Fatalf("Save CLI help: %v", err)
	}
}

func saveCLIHelp(dst string) error {
	err := exec.Command("make", "build").Run()
	if err != nil {
		return fmt.Errorf("can't run make build: %w", err)
	}

	lintersOut, err := exec.Command("./golangci-lint", "help", "linters").Output()
	if err != nil {
		return fmt.Errorf("can't run linters cmd: %w", err)
	}

	lintersOutParts := bytes.Split(lintersOut, []byte("\n\n"))

	helpCmd := exec.Command("./golangci-lint", "run", "-h")
	helpCmd.Env = append(helpCmd.Env, os.Environ()...)
	helpCmd.Env = append(helpCmd.Env, "HELP_RUN=1") // make default concurrency stable: don't depend on machine CPU number
	help, err := helpCmd.Output()
	if err != nil {
		return fmt.Errorf("can't run help cmd: %w", err)
	}

	helpLines := bytes.Split(help, []byte("\n"))
	shortHelp := bytes.Join(helpLines[2:], []byte("\n"))

	data := types.CLIHelp{
		Enable:  string(lintersOutParts[0]),
		Disable: string(lintersOutParts[1]),
		Help:    string(shortHelp),
	}

	return saveToJSONFile(dst, data)
}

func saveToJSONFile(dst string, data any) error {
	file, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("open file (%s): %w", dst, err)
	}

	defer func() { _ = file.Close() }()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	err = encoder.Encode(data)
	if err != nil {
		return fmt.Errorf("encode JSON (%s): %w", dst, err)
	}

	return nil
}
