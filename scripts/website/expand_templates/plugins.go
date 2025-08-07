package main

import (
	"io"
	"os"
	"path/filepath"
)

func copyPluginReference(dir string) error {
	in, err := os.Open(".custom-gcl.reference.yml")
	if err != nil {
		return err
	}

	defer func() { _ = in.Close() }()

	out, err := os.Create(filepath.Join(dir, ".custom-gcl.reference.yml"))
	if err != nil {
		return err
	}

	defer func() { _ = out.Close() }()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	return nil
}
