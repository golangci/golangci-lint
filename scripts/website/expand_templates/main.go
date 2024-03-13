package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/golangci/golangci-lint/internal/renameio"
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
		return fmt.Errorf("failed to walk dir: %w", err)
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
		return fmt.Errorf("failed to read %s: %w", path, err)
	}

	content := string(contentBytes)
	hasReplacements := false
	for key, replacement := range replacements {
		nextContent := content
		nextContent = strings.ReplaceAll(nextContent, fmt.Sprintf("{.%s}", key), replacement)

		// Yaml formatter in mdx code section makes extra spaces, need to match them too.
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
	if err = renameio.WriteFile(path, []byte(content), os.ModePerm); err != nil {
		return fmt.Errorf("failed to write changes to file %s: %w", path, err)
	}

	return nil
}

type latestRelease struct {
	TagName string `json:"tag_name"`
}

func getLatestVersion() (string, error) {
	req, err := http.NewRequest( //nolint:noctx
		http.MethodGet,
		"https://api.github.com/repos/golangci/golangci-lint/releases/latest",
		http.NoBody,
	)
	if err != nil {
		return "", fmt.Errorf("failed to prepare a http request: %w", err)
	}
	req.Header.Add("Accept", "application/vnd.github.v3+json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to get http response for the latest tag: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read a body for the latest tag: %w", err)
	}
	release := latestRelease{}
	err = json.Unmarshal(body, &release)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal the body for the latest tag: %w", err)
	}
	return release.TagName, nil
}

func buildTemplateContext() (map[string]string, error) {
	snippets, err := getExampleSnippets()
	if err != nil {
		return nil, err
	}

	pluginReference, err := getPluginReference()
	if err != nil {
		return nil, fmt.Errorf("failed to read plugin reference file: %w", err)
	}

	helps, err := readJSONFile[types.CLIHelp](filepath.Join("assets", "cli-help.json"))
	if err != nil {
		return nil, err
	}

	changeLog, err := os.ReadFile("CHANGELOG.md")
	if err != nil {
		return nil, err
	}

	latestVersion, err := getLatestVersion()
	if err != nil {
		return nil, fmt.Errorf("failed to get the latest version: %w", err)
	}

	exclusions, err := getDefaultExclusions()
	if err != nil {
		return nil, fmt.Errorf("default exclusions: %w", err)
	}

	return map[string]string{
		"CustomGCLReference":               pluginReference,
		"ConfigurationExample":             snippets.ConfigurationFile,
		"LintersCommandOutputEnabledOnly":  helps.Enable,
		"LintersCommandOutputDisabledOnly": helps.Disable,
		"EnabledByDefaultLinters":          getLintersListMarkdown(true),
		"DisabledByDefaultLinters":         getLintersListMarkdown(false),
		"DefaultExclusions":                exclusions,
		"ThanksList":                       getThanksList(),
		"RunHelpText":                      helps.Help,
		"ChangeLog":                        string(changeLog),
		"LatestVersion":                    latestVersion,
	}, nil
}

func readJSONFile[T any](src string) (T, error) {
	file, err := os.Open(src)
	if err != nil {
		var zero T
		return zero, err
	}

	var result T
	err = json.NewDecoder(file).Decode(&result)
	if err != nil {
		var zero T
		return zero, err
	}

	return result, nil
}
