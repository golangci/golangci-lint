package config

import (
	"testing"

	"github.com/golangci/golangci-lint/pkg/logutils"
)

func TestReader(t *testing.T) {
	var toCfg, commandLineCfg Config
	logger := logutils.NewStderrLog("")
	logger.SetLevel(logutils.LogLevelDebug)

	commandLineCfg.Run.Config = "../../.golangci.reference.yml"
	reader := NewFileReader(&toCfg, &commandLineCfg, logger)
	if err := reader.Read(); err != nil {
		t.Fatalf("unexpected error reading .golangci.reference.yml: %v", err)
	}
}
