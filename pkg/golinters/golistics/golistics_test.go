package golistics

import (
	"testing"

	"github.com/golangci/golangci-lint/v2/test/testshared/integration"
)

func TestGolistics(t *testing.T) {
	integration.RunTestdata(t)
}
