package fsutils

import (
	"go/build"
	"path"
)

func GetProjectRoot() string {
	return path.Join(build.Default.GOPATH, "src", "github.com", "golangci", "golangci-worker")
}
