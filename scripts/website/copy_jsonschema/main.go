package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

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

	err = os.CopyFS(dstDir, os.DirFS("jsonschema"))
	if err != nil {
		return fmt.Errorf("copy FS: %w", err)
	}

	err = copyLatestSchema()
	if err != nil {
		return fmt.Errorf("copy files: %w", err)
	}

	return nil
}

func copyLatestSchema() error {
	src := filepath.FromSlash("jsonschema/golangci.jsonschema.json")

	latest, err := github.GetLatestVersion()
	if err != nil {
		return fmt.Errorf("get latest release version: %w", err)
	}

	version, err := hcversion.NewVersion(latest)
	if err != nil {
		return fmt.Errorf("parse version: %w", err)
	}

	files := []string{
		fmt.Sprintf("golangci.v%d.jsonschema.json", version.Segments()[0]),
		fmt.Sprintf("golangci.v%d.%d.jsonschema.json", version.Segments()[0], version.Segments()[1]),
	}

	source, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("open file %s: %w", src, err)
	}

	defer func() { _ = source.Close() }()

	info, err := source.Stat()
	if err != nil {
		return fmt.Errorf("file %s not found: %w", src, err)
	}

	for _, dst := range files {
		err = copyFile(dst, source, info)
		if err != nil {
			return fmt.Errorf("copy file %s to %s: %w", src, dst, err)
		}
	}

	return nil
}

func copyFile(dst string, source io.Reader, info os.FileInfo) error {
	destination, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, info.Mode())
	if err != nil {
		return fmt.Errorf("create file %s: %w", dst, err)
	}

	defer func() { _ = destination.Close() }()

	_, err = io.Copy(destination, source)
	if err != nil {
		return err
	}

	return nil
}
