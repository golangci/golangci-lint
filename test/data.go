package test

import "path/filepath"

const testdataDir = "testdata"
const binName = "../golangci-lint"

var minimalPkg = getTestDataDir("minimalpkg")

func getProjectRoot() string {
	return filepath.Join("..", "...")
}

func getTestDataDir(names ...string) string {
	parts := append([]string{testdataDir}, names...)
	return filepath.Join(parts...)
}
