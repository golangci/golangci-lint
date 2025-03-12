package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goformatters"
	"github.com/golangci/golangci-lint/v2/pkg/lint/linter"
	"github.com/golangci/golangci-lint/v2/pkg/lint/lintersdb"
	"github.com/golangci/golangci-lint/v2/pkg/result/processors"
	"github.com/golangci/golangci-lint/v2/scripts/website/types"
)

func main() {
	err := saveLinters()
	if err != nil {
		log.Fatalf("Save linters: %v", err)
	}

	err = saveFormatters()
	if err != nil {
		log.Fatalf("Save formatters: %v", err)
	}

	err = saveDefaultExclusions()
	if err != nil {
		log.Fatalf("Save default exclusions: %v", err)
	}

	err = saveCLIHelp(filepath.Join("assets", "cli-help.json"))
	if err != nil {
		log.Fatalf("Save CLI help: %v", err)
	}
}

func saveFormatters() error {
	linters, _ := lintersdb.NewLinterBuilder().Build(config.NewDefault())

	var wraps []types.LinterWrapper
	for _, l := range linters {
		if l.IsDeprecated() && l.Deprecation.Level > linter.DeprecationWarning {
			continue
		}

		if !goformatters.IsFormatter(l.Name()) {
			continue
		}

		wrapper := types.LinterWrapper{
			Name:             l.Linter.Name(),
			Desc:             l.Linter.Desc(),
			Groups:           l.Groups,
			LoadMode:         l.LoadMode,
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

	return saveToJSONFile(filepath.Join("assets", "formatters-info.json"), wraps)
}

func saveLinters() error {
	linters, _ := lintersdb.NewLinterBuilder().Build(config.NewDefault())

	var wraps []types.LinterWrapper
	for _, l := range linters {
		if l.IsDeprecated() && l.Deprecation.Level > linter.DeprecationWarning {
			continue
		}

		if goformatters.IsFormatter(l.Name()) {
			continue
		}

		wrapper := types.LinterWrapper{
			Name:             l.Linter.Name(),
			Desc:             l.Linter.Desc(),
			Groups:           l.Groups,
			LoadMode:         l.LoadMode,
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

	return saveToJSONFile(filepath.Join("assets", "linters-info.json"), wraps)
}

func saveDefaultExclusions() error {
	data := make(map[string][]types.ExcludeRule)

	for name, rules := range processors.LinterExclusionPresets {
		for _, rule := range rules {
			data[name] = append(data[name], types.ExcludeRule{
				Linters:    rule.Linters,
				Path:       rule.Path,
				PathExcept: rule.PathExcept,
				Text:       rule.Text,
				Source:     rule.Source,
			})
		}
	}

	return saveToJSONFile(filepath.Join("assets", "exclusion-presets.json"), data)
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

	rumCmdHelp, err := getCmdHelp("run")
	if err != nil {
		return err
	}

	fmtCmdHelp, err := getCmdHelp("fmt")
	if err != nil {
		return err
	}

	data := types.CLIHelp{
		Enable:     string(lintersOutParts[0]),
		RunCmdHelp: rumCmdHelp,
		FmtCmdHelp: fmtCmdHelp,
	}

	return saveToJSONFile(dst, data)
}

func getCmdHelp(name string) (string, error) {
	helpCmd := exec.Command("./golangci-lint", name, "-h")
	helpCmd.Env = append(helpCmd.Env, os.Environ()...)

	help, err := helpCmd.Output()
	if err != nil {
		return "", fmt.Errorf("can't run help cmd: %w", err)
	}

	helpLines := bytes.Split(help, []byte("\n"))
	shortHelp := bytes.Join(helpLines[2:], []byte("\n"))

	return string(shortHelp), nil
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
