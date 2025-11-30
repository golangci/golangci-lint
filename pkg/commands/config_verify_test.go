package commands

import (
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_createSchemaURL(t *testing.T) {
	testCases := []struct {
		desc     string
		flag     string
		info     BuildInfo
		expected string
	}{
		{
			desc:     "schema flag only",
			flag:     "https://example.com",
			expected: "https://example.com",
		},
		{
			desc: "schema flag and build info",
			flag: "https://example.com",
			info: BuildInfo{
				Version: "v1.0.0",
				Commit:  "cd8b11773c6c1f595e8eb98c0d4310af20ae20df",
			},
			expected: "https://example.com",
		},
		{
			desc: "version and commit",
			info: BuildInfo{
				Version: "v1.0.0",
				Commit:  "cd8b11773c6c1f595e8eb98c0d4310af20ae20df",
			},
			expected: "https://golangci-lint.run/jsonschema/golangci.v1.0.jsonschema.json",
		},
		{
			desc: "commit only",
			info: BuildInfo{
				Commit: "cd8b11773c6c1f595e8eb98c0d4310af20ae20df",
			},
			expected: "https://raw.githubusercontent.com/golangci/golangci-lint/cd8b11773c6c1f595e8eb98c0d4310af20ae20df/jsonschema/golangci.next.jsonschema.json",
		},
		{
			desc: "version devel and commit",
			info: BuildInfo{
				Version: "(devel)",
				Commit:  "cd8b11773c6c1f595e8eb98c0d4310af20ae20df",
			},
			expected: "https://raw.githubusercontent.com/golangci/golangci-lint/cd8b11773c6c1f595e8eb98c0d4310af20ae20df/jsonschema/golangci.next.jsonschema.json",
		},
		{
			desc: "composite commit info",
			info: BuildInfo{
				Version: "",
				Commit:  `(cd8b11773c6c1f595e8eb98c0d4310af20ae20df, modified: "false", mod sum: "123")`,
			},
			expected: "https://raw.githubusercontent.com/golangci/golangci-lint/cd8b11773c6c1f595e8eb98c0d4310af20ae20df/jsonschema/golangci.next.jsonschema.json",
		},
		{
			desc: "v0 version",
			info: BuildInfo{
				Version: "v0.0.0-20250213211019-0a603e49e5e9",
				Commit:  `(0a603e49e5e9870f5f9f2035bcbe42cd9620a9d5, modified: "false", mod sum: "123")`,
			},
			expected: "https://raw.githubusercontent.com/golangci/golangci-lint/0a603e49e5e9870f5f9f2035bcbe42cd9620a9d5/jsonschema/golangci.next.jsonschema.json",
		},
		{
			desc: "dirty",
			info: BuildInfo{
				Version: "v1.64.6-0.20250225205237-3eecab1ebde9+dirty",
				Commit:  `(3eecab1ebde9, modified: "false", mod sum: "123")`,
			},
			expected: "https://golangci-lint.run/jsonschema/golangci.v1.64.jsonschema.json",
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
			flags.String("schema", "", "")
			if test.flag != "" {
				_ = flags.Set("schema", test.flag)
			}

			schemaURL, err := createSchemaURL(flags, test.info)
			require.NoError(t, err)

			assert.Equal(t, test.expected, schemaURL)
		})
	}
}

func Test_createSchemaURL_error(t *testing.T) {
	testCases := []struct {
		desc     string
		info     BuildInfo
		expected string
	}{
		{
			desc: "unknown commit",
			info: BuildInfo{
				Commit: "unknown",
			},
			expected: "unknown commit information",
		},
		{
			desc: "detailed unknown commit",
			info: BuildInfo{
				Version: "",
				Commit:  `(unknown, modified: ?, mod sum: "")`,
			},
			expected: "unknown commit information",
		},
		{
			desc: "commit ?",
			info: BuildInfo{
				Commit: "?",
			},
			expected: "version not found",
		},
		{
			desc: "version devel only",
			info: BuildInfo{
				Version: "(devel)",
			},
			expected: "version not found",
		},
		{
			desc: "invalid version",
			info: BuildInfo{
				Version: "example",
			},
			expected: "parse version: malformed version: example",
		},
		{
			desc: "invalid composite commit info",
			info: BuildInfo{
				Version: "",
				Commit:  `(cd8b11773c6c1f595e8eb98c0d4310af20ae20df)`,
			},
			expected: "commit information not found",
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
			flags.String("schema", "", "")

			_, err := createSchemaURL(flags, test.info)
			require.EqualError(t, err, test.expected)
		})
	}
}
