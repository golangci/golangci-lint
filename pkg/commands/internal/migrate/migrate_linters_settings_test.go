package migrate

import (
	"testing"

	"github.com/golangci/golangci-lint/v2/pkg/commands/internal/migrate/versionone"
	"github.com/golangci/golangci-lint/v2/pkg/commands/internal/migrate/versiontwo"
	"github.com/google/go-cmp/cmp"
)

func Test_toStaticCheckSettings(t *testing.T) {
	old := versionone.LintersSettings{
		Staticcheck: versionone.StaticCheckSettings{
			Checks: []string{"all", "-SA1000"},
		},
	}
	new := toStaticCheckSettings(old)

	expected := versiontwo.StaticCheckSettings{
		Checks: []string{"all", "-SA1000"},
	}

	if diff := cmp.Diff(new, expected); diff != "" {
		t.Errorf("toStaticCheckSettings() mismatch (-got +want):\n%s", diff)
	}
}
