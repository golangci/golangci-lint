package main

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters"
	"github.com/golangci/golangci-shared/pkg/executors"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	var cfg config.Config
	config.ReadFromCommandLine(&cfg)

	linters := golinters.GetSupportedLinters()
	ctx := context.Background()

	ex, err := os.Executable()
	if err != nil {
		return err
	}
	exPath := filepath.Dir(ex)
	exec := executors.NewShell(exPath)

	for _, linter := range linters {
		res, err := linter.Run(ctx, exec)
		if err != nil {
			return err
		}
		log.Print(res)
	}

	return nil
}
