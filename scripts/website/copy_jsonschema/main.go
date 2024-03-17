package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	hcversion "github.com/hashicorp/go-version"

	"github.com/golangci/golangci-lint/scripts/website/gh"
)

func main() {
	dstDir := filepath.FromSlash("docs/static/jsonschema/")

	err := os.RemoveAll(dstDir)
	if err != nil {
		log.Fatalf("remove dir: %v", err)
	}

	err = os.MkdirAll(dstDir, os.ModePerm)
	if err != nil {
		log.Fatalf("make dir: %v", err)
	}

	// The key is the destination file.
	// The value is the source file.
	files := map[string]string{}

	entries, err := os.ReadDir("jsonschema")
	if err != nil {
		log.Fatalf("read dir: %v", err)
	}

	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".jsonschema.json") {
			files[entry.Name()] = entry.Name()
		}
	}

	latest, err := gh.GetLatestVersion()
	if err != nil {
		log.Fatalf("get latest release version: %v", err)
	}

	version, err := hcversion.NewVersion(latest)
	if err != nil {
		log.Fatalf("parse version: %v", err)
	}

	versioned := fmt.Sprintf("golangci.v%d.%d.jsonschema.json", version.Segments()[0], version.Segments()[1])
	files[versioned] = "golangci.jsonschema.json"

	for dst, src := range files {
		err := copyFile(filepath.Join("jsonschema", src), filepath.Join(dstDir, dst))
		if err != nil {
			log.Fatalf("Copy files: %v", err)
		}
	}
}

func copyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("open file %s: %w", src, err)
	}

	info, err := source.Stat()
	if err != nil {
		return fmt.Errorf("file %s not found: %w", src, err)
	}

	destination, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, info.Mode())
	if err != nil {
		return fmt.Errorf("create file %s: %w", dst, err)
	}

	_, err = io.Copy(destination, source)
	if err != nil {
		return fmt.Errorf("copy file %s to %s: %w", src, dst, err)
	}

	return nil
}
