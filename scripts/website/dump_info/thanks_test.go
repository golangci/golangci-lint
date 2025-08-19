package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/golangci/golangci-lint/v2/pkg/lint/linter"
	"github.com/golangci/golangci-lint/v2/pkg/result"
)

type FakeLinter struct {
	name string
}

func (*FakeLinter) Run(_ context.Context, _ *linter.Context) ([]*result.Issue, error) {
	return nil, nil
}

func (f *FakeLinter) Name() string {
	return f.name
}

func (*FakeLinter) Desc() string {
	return "fake linter"
}

func Test_extractInfo(t *testing.T) {
	testCases := []struct {
		desc     string
		lc       *linter.Config
		expected authorInfo
	}{
		{
			desc: "from GitHub",
			lc: &linter.Config{
				Linter:      &FakeLinter{name: "fake"},
				OriginalURL: "https://github.com/owner/linter",
			},
			expected: authorInfo{Author: "owner", Host: "github"},
		},
		{
			desc: "from GitLab",
			lc: &linter.Config{
				Linter:      &FakeLinter{name: "fake"},
				OriginalURL: "https://gitlab.com/owner/linter",
			},
			expected: authorInfo{Author: "owner", Host: "gitlab"},
		},
		{
			desc: "gostaticanalysis",
			lc: &linter.Config{
				Linter:      &FakeLinter{name: "fake"},
				OriginalURL: "https://github.com/gostaticanalysis/linter",
			},
			expected: authorInfo{Author: "tenntenn", Host: "github"},
		},
		{
			desc: "go-simpler",
			lc: &linter.Config{
				Linter:      &FakeLinter{name: "fake"},
				OriginalURL: "https://github.com/go-simpler/linter",
			},
			expected: authorInfo{Author: "tmzane", Host: "github"},
		},
		{
			desc: "curioswitch",
			lc: &linter.Config{
				Linter:      &FakeLinter{name: "fake"},
				OriginalURL: "https://github.com/curioswitch/linter",
			},
			expected: authorInfo{Author: "chokoswitch", Host: "github"},
		},
		{
			desc: "GaijinEntertainment",
			lc: &linter.Config{
				Linter:      &FakeLinter{name: "fake"},
				OriginalURL: "https://github.com/GaijinEntertainment/linter",
			},
			expected: authorInfo{Author: "xobotyi", Host: "github"},
		},
		{
			desc: "OpenPeeDeeP",
			lc: &linter.Config{
				Linter:      &FakeLinter{name: "fake"},
				OriginalURL: "https://github.com/OpenPeeDeeP/linter",
			},
			expected: authorInfo{Author: "dixonwille", Host: "github"},
		},
		{
			desc: "misspell",
			lc: &linter.Config{
				Linter:      &FakeLinter{name: "misspell"},
				OriginalURL: "https://github.com/myorg/linter",
			},
			expected: authorInfo{Author: "client9", Host: "github"},
		},
		{
			desc: "pkg.go.dev",
			lc: &linter.Config{
				Linter:      &FakeLinter{name: "fake"},
				OriginalURL: "https://pkg.go.dev/linter",
			},
			expected: authorInfo{Author: "golang", Host: "github"},
		},
		{
			desc: "golangci",
			lc: &linter.Config{
				Linter:      &FakeLinter{name: "fake"},
				OriginalURL: "https://github.com/golangci/linter",
			},
			expected: authorInfo{},
		},
		{
			desc: "invalid",
			lc: &linter.Config{
				Linter:      &FakeLinter{name: "fake"},
				OriginalURL: "https://example.com/linter",
			},
			expected: authorInfo{},
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			info := extractInfo(test.lc)

			assert.Equal(t, test.expected, info)
		})
	}
}
