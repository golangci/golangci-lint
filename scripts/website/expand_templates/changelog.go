package main

import (
	"bytes"
	"os"
	"path/filepath"
)

func copyChangelogs(dir string) error {
	data := map[string]string{
		"CHANGELOG.md":    filepath.Join(dir, "raw_changelog.tmp"),
		"CHANGELOG-v1.md": filepath.Join(dir, "raw_changelog_v1.tmp"),
	}

	for src, dst := range data {
		err := copyChangelog(src, dst)
		if err != nil {
			return err
		}
	}

	return nil
}

func copyChangelog(src, dst string) error {
	markerStart := "<!-- START --->"

	in, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}

	defer func() { _ = out.Close() }()

	var write bool

	for line := range bytes.Lines(bytes.ReplaceAll(in, []byte("### "), []byte("## "))) {
		if bytes.Equal(bytes.TrimSpace(line), []byte(markerStart)) {
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
