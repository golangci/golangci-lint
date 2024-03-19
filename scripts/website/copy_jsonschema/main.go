package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	hcversion "github.com/hashicorp/go-version"

	"github.com/golangci/golangci-lint/scripts/website/github"
)

func main() {
	err := copySchemas()
	if err != nil {
		log.Fatal(err)
	}
}

func copySchemas() error {
	dstDir := filepath.FromSlash("docs/static/jsonschema/")

	err := os.RemoveAll(dstDir)
	if err != nil {
		return fmt.Errorf("remove dir: %w", err)
	}

	err = os.MkdirAll(dstDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("make dir: %w", err)
	}

	// The key is the destination file.
	// The value is the source file.
	files := map[string]string{}

	entries, err := os.ReadDir("jsonschema")
	if err != nil {
		return fmt.Errorf("read dir: %w", err)
	}

	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".jsonschema.json") {
			files[entry.Name()] = entry.Name()
		}
	}

	latest, err := github.GetLatestVersion()
	if err != nil {
		return fmt.Errorf("get latest release version: %w", err)
	}

	version, err := hcversion.NewVersion(latest)
	if err != nil {
		return fmt.Errorf("parse version: %w", err)
	}

	versioned := fmt.Sprintf("golangci.v%d.%d.jsonschema.json", version.Segments()[0], version.Segments()[1])
	files[versioned] = "golangci.jsonschema.json"

	for dst, src := range files {
		err := copyFile(filepath.Join("jsonschema", src), filepath.Join(dstDir, dst))
		if err != nil {
			return fmt.Errorf("copy files: %w", err)
		}
	}

	return nil
}

func copyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("open file %s: %w", src, err)
	}

	defer func() { _ = source.Close() }()

	info, err := source.Stat()
	if err != nil {
		return fmt.Errorf("file %s not found: %w", src, err)
	}

	destination, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, info.Mode())
	if err != nil {
		return fmt.Errorf("create file %s: %w", dst, err)
	}

	defer func() { _ = destination.Close() }()

	_, err = io.Copy(destination, source)
	if err != nil {
		return fmt.Errorf("copy file %s to %s: %w", src, dst, err)
	}

	return nil
}
