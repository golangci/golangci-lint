package main

import (
	"fmt"
)

// ugh, needs to be turned into a config
func processSource(source, repo, path, file, exe, nametpl string) (out []byte, err error) {
	switch source {
	case "godownloader":
		// https://github.com/goreleaser/godownloader
		out, err = processGodownloader(repo, path, file)
	case "equinoxio":
		// https://equinox.io
		out, err = processEquinoxio(repo)
	case "raw":
		// raw mode is when people upload direct binaries
		// to GitHub releases that are not  not tar'ed or zip'ed.
		// For example:
		//   https://github.com/mvdan/sh/releases
		out, err = processRaw(repo, exe, nametpl)
	default:
		return nil, fmt.Errorf("unknown source %q", source)
	}
	return
}
