package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/rogpeppe/go-internal/lockedfile"

	"github.com/golangci/golangci-lint/scripts/website/github"
	"github.com/golangci/golangci-lint/scripts/website/types"
)

func main() {
	replacements, err := buildTemplateContext()
	if err != nil {
		log.Fatalf("Failed to build template context: %s", err)
	}

	if err := rewriteDocs(replacements); err != nil {
		log.Fatalf("Failed to rewrite docs: %s", err)
	}

	log.Print("Successfully expanded templates")
}

func rewriteDocs(replacements map[string]string) error {
	madeReplacements := map[string]bool{}

	err := filepath.Walk(filepath.Join("docs", "src", "docs"),
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			return processDoc(path, replacements, madeReplacements)
		})
	if err != nil {
		return fmt.Errorf("walk dir: %w", err)
	}

	if len(madeReplacements) != len(replacements) {
		for key := range replacements {
			if !madeReplacements[key] {
				log.Printf("Replacement %q wasn't performed", key)
			}
		}
		return fmt.Errorf("%d replacements weren't performed", len(replacements)-len(madeReplacements))
	}
	return nil
}

func processDoc(path string, replacements map[string]string, madeReplacements map[string]bool) error {
	contentBytes, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}

	content := string(contentBytes)
	hasReplacements := false
	for key, replacement := range replacements {
		nextContent := content
		nextContent = strings.ReplaceAll(nextContent, fmt.Sprintf("{.%s}", key), replacement)

		// YAML formatter in mdx code section makes extra spaces, need to match them too.
		nextContent = strings.ReplaceAll(nextContent, fmt.Sprintf("{ .%s }", key), replacement)

		if nextContent != content {
			hasReplacements = true
			madeReplacements[key] = true
			content = nextContent
		}
	}
	if !hasReplacements {
		return nil
	}

	log.Printf("Expanded template in %s, saving it", path)
	if err = lockedfile.Write(path, bytes.NewBufferString(content), os.ModePerm); err != nil {
		return fmt.Errorf("write changes to file %s: %w", path, err)
	}

	return nil
}

func buildTemplateContext() (map[string]string, error) {
	snippets, err := NewExampleSnippetsExtractor().GetExampleSnippets()
	if err != nil {
		return nil, err
	}

	pluginReference, err := getPluginReference()
	if err != nil {
		return nil, fmt.Errorf("read plugin reference file: %w", err)
	}

	helps, err := readJSONFile[types.CLIHelp](filepath.Join("assets", "cli-help.json"))
	if err != nil {
		return nil, err
	}

	changeLog, err := os.ReadFile("CHANGELOG.md")
	if err != nil {
		return nil, fmt.Errorf("read CHANGELOG.md: %w", err)
	}

	latestVersion, err := github.GetLatestVersion()
	if err != nil {
		return nil, fmt.Errorf("get the latest version: %w", err)
	}

	exclusions, err := getExclusionPresets()
	if err != nil {
		return nil, fmt.Errorf("default exclusions: %w", err)
	}

	return map[string]string{
		"CustomGCLReference":              pluginReference,
		"LintersExample":                  snippets.LintersSettings,
		"FormattersExample":               snippets.FormattersSettings,
		"ConfigurationExample":            snippets.ConfigurationFile,
		"LintersCommandOutputEnabledOnly": helps.Enable,
		"EnabledByDefaultLinters":         getLintersListMarkdown(true, filepath.Join("assets", "linters-info.json")),
		"DisabledByDefaultLinters":        getLintersListMarkdown(false, filepath.Join("assets", "linters-info.json")),
		"Formatters":                      getLintersListMarkdown(false, filepath.Join("assets", "formatters-info.json")),
		"ExclusionPresets":                exclusions,
		"ThanksList":                      getThanksList(),
		"RunHelpText":                     helps.Help,
		"ChangeLog":                       string(changeLog),
		"LatestVersion":                   latestVersion,
	}, nil
}

func readJSONFile[T any](src string) (T, error) {
	file, err := os.Open(src)
	if err != nil {
		var zero T
		return zero, fmt.Errorf("open file %s: %w", src, err)
	}

	defer func() { _ = file.Close() }()

	var result T
	err = json.NewDecoder(file).Decode(&result)
	if err != nil {
		var zero T
		return zero, fmt.Errorf("decode JSON file %s: %w", src, err)
	}

	return result, nil
}
