package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/golangci/golangci-lint/v2/scripts/website/github"
)

func main() {
	err := saveTmp(filepath.Join("docs", ".tmp"))
	if err != nil {
		log.Fatalf("Save tmp: %s", err)
	}

	err = saveData(filepath.Join("docs", "data"))
	if err != nil {
		log.Fatalf("Save data: %s", err)
	}

	log.Print("Successfully expanded templates")
}

func saveTmp(tmpDir string) error {
	err := os.RemoveAll(tmpDir)
	if err != nil {
		return fmt.Errorf("remove tmp dir: %w", err)
	}

	err = os.MkdirAll(tmpDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("create tmp dir: %w", err)
	}

	err = copyPluginReference(tmpDir)
	if err != nil {
		return fmt.Errorf("copy plugin reference: %w", err)
	}

	err = copyChangelogs(tmpDir)
	if err != nil {
		return fmt.Errorf("copy changelog: %w", err)
	}

	return nil
}

func saveData(dir string) error {
	latestVersion, err := github.GetLatestVersion()
	if err != nil {
		return fmt.Errorf("get the latest version: %w", err)
	}

	err = saveToJSONFile(filepath.Join(dir, "version.json"), map[string]string{"version": latestVersion})
	if err != nil {
		return fmt.Errorf("save latest version: %w", err)
	}

	snippets, err := NewExampleSnippetsExtractor().GetExampleSnippets()
	if err != nil {
		return fmt.Errorf("get example snippets: %w", err)
	}

	err = saveToJSONFile(filepath.Join(dir, "linter_settings.json"), snippets.LintersSettings)
	if err != nil {
		return fmt.Errorf("save linter snippets: %w", err)
	}

	err = saveToJSONFile(filepath.Join(dir, "formatter_settings.json"), snippets.FormattersSettings)
	if err != nil {
		return fmt.Errorf("save formatter snippets: %w", err)
	}

	err = saveToJSONFile(filepath.Join(dir, "configuration_file.json"), snippets.ConfigurationFile)
	if err != nil {
		return fmt.Errorf("save configuration file snippets: %w", err)
	}

	return nil
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
