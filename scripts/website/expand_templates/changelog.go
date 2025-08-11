package main

import (
	"bufio"
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

	// TODO(ldez): use bytes.Lines when min go1.24 (and remove the new line)
	scanner := bufio.NewScanner(bytes.NewBuffer(bytes.ReplaceAll(in, []byte("### "), []byte("## "))))
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Bytes()

		if bytes.Equal(bytes.TrimSpace(line), []byte(marker)) {
			write = true
			continue
		}

		if !write {
			continue
		}

		line = append(line, '\n')

		_, err = out.Write(line)
		if err != nil {
			return err
		}
	}

	return nil
}
