package golinters

import (
	"fmt"
	"strings"

	"github.com/golangci/golangci-lint/pkg/config"
)

func formatCode(code string, cfg *config.Run) string {
	if strings.Contains(code, "`") {
		return code // TODO: properly escape or remove
	}

	return fmt.Sprintf("`%s`", code)
}
