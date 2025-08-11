package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	dataDir := filepath.Join("docs", "data")

	err := saveLinters(filepath.Join(dataDir, "linters_info.json"))
	if err != nil {
		log.Fatalf("Save linters: %v", err)
	}

	err = saveFormatters(filepath.Join(dataDir, "formatters_info.json"))
	if err != nil {
		log.Fatalf("Save formatters: %v", err)
	}

	err = saveDefaultExclusions(filepath.Join(dataDir, "exclusion_presets.json"))
	if err != nil {
		log.Fatalf("Save default exclusions: %v", err)
	}

	err = saveCLIHelp(context.Background(), filepath.Join(dataDir, "cli_help.json"))
	if err != nil {
		log.Fatalf("Save CLI help: %v", err)
	}

	err = saveThanksList(filepath.Join(dataDir, "thanks.json"))
	if err != nil {
		log.Fatalf("Save thanks list: %v", err)
	}
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
