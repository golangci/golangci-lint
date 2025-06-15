package main

import (
	"bytes"
	"os"
	"path/filepath"
)

func copyChangelog(dir string) error {
	marker := "<!-- START --->"

	in, err := os.ReadFile("CHANGELOG.md")
	if err != nil {
		return err
	}

	out, err := os.Create(filepath.Join(dir, "raw_changelog.tmp"))
	if err != nil {
		return err
	}

	defer func() { _ = out.Close() }()

	var write bool

	for line := range bytes.Lines(bytes.ReplaceAll(in, []byte("### "), []byte("## "))) {
		if bytes.Equal(bytes.TrimSpace(line), []byte(marker)) {
			write = true
			continue
		}

		if !write {
			continue
		}

		_, err = out.Write(line)
		if err != nil {
			return err
		}
	}

	return nil
}
