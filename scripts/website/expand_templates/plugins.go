package main

import (
	"fmt"
	"os"
)

func getPluginReference() (string, error) {
	reference, err := os.ReadFile(".custom-gcl.reference.yml")
	if err != nil {
		return "", fmt.Errorf("can't read .custom-gcl.reference.yml: %w", err)
	}

	return string(reference), nil
}
