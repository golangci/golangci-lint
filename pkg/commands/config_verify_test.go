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
			desc: "commit unknown",
			info: BuildInfo{
				Commit: "unknown",
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
			expected: "parse version: Malformed version: example",
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
