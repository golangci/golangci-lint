package migrate

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/golangci/golangci-lint/v2/pkg/commands/internal/migrate/ptr"
	"github.com/golangci/golangci-lint/v2/pkg/commands/internal/migrate/versionone"
)

func Test_disableAllFilter(t *testing.T) {
	testCases := []struct {
		desc     string
		old      versionone.Linters
		expected []string
	}{
		{
			desc: "no presets, fast",
			old: versionone.Linters{
				DisableAll: ptr.Pointer(true),
				Enable:     nil,
				Fast:       ptr.Pointer(false),
				Presets:    nil,
			},
			expected: nil,
		},
		{
			desc: "no presets, fast",
			old: versionone.Linters{
				DisableAll: ptr.Pointer(true),
				Enable:     nil,
				Fast:       ptr.Pointer(true),
				Presets:    nil,
			},
			expected: nil,
		},
		{
			desc: "no presets, enable",
			old: versionone.Linters{
				DisableAll: ptr.Pointer(true),
				Enable:     []string{"lll", "misspell", "govet"},
				Fast:       ptr.Pointer(false),
				Presets:    nil,
			},
			expected: []string{"govet", "lll", "misspell"},
		},
		{
			desc: "fast, no presets, enable",
			old: versionone.Linters{
				DisableAll: ptr.Pointer(true),
				Enable:     []string{"lll", "misspell", "govet"},
				Fast:       ptr.Pointer(true),
				Presets:    nil,
			},
			expected: []string{"govet", "lll", "misspell"},
		},
		{
			desc: "presets, enable",
			old: versionone.Linters{
				DisableAll: ptr.Pointer(true),
				Enable:     []string{"lll", "misspell", "govet"},
				Fast:       ptr.Pointer(false),
				Presets:    []string{"comment", "error", "format"},
			},
			expected: []string{
				"dupword",
				"err113",
				"errcheck",
				"errorlint",
				"gci",
				"godot",
				"godox",
				"gofmt",
				"gofumpt",
				"goimports",
				"govet",
				"lll",
				"misspell",
				"wrapcheck",
			},
		},
		{
			desc: "presets, enable, fast",
			old: versionone.Linters{
				DisableAll: ptr.Pointer(true),
				Enable:     []string{"lll", "misspell", "govet"},
				Fast:       ptr.Pointer(true),
				Presets:    []string{"comment", "error", "format"},
			},
			expected: []string{
				"dupword",
				"gci",
				"godot",
				"godox",
				"gofmt",
				"gofumpt",
				"goimports",
				"govet",
				"lll",
				"misspell",
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			results := disableAllFilter(test.old)

			assert.Equal(t, test.expected, results)
		})
	}
}

func Test_enableAllFilter(t *testing.T) {
	testCases := []struct {
		desc     string
		old      versionone.Linters
		expected []string
	}{
		{
			desc: "no options",
			old: versionone.Linters{
				EnableAll: ptr.Pointer(true),
				Disable:   nil,
				Fast:      ptr.Pointer(false),
				Presets:   nil,
			},
			expected: nil,
		},
		{
			desc: "presets (ignored)",
			old: versionone.Linters{
				EnableAll: ptr.Pointer(true),
				Disable:   nil,
				Fast:      ptr.Pointer(false),
				Presets:   []string{"comment", "error", "format"},
			},
			expected: nil,
		},
		{
			desc: "fast",
			old: versionone.Linters{
				EnableAll: ptr.Pointer(true),
				Disable:   nil,
				Fast:      ptr.Pointer(true),
				Presets:   nil,
			},
			expected: []string{"asasalint", "bodyclose", "canonicalheader", "containedctx", "contextcheck", "durationcheck", "err113", "errcheck", "errchkjson", "errname", "errorlint", "exhaustive", "exhaustruct", "exptostd", "fatcontext", "forbidigo", "forcetypeassert", "ginkgolinter", "gochecknoglobals", "gochecksumtype", "gocritic", "gosec", "gosimple", "gosmopolitan", "govet", "iface", "importas", "intrange", "ireturn", "loggercheck", "makezero", "mirror", "musttag", "nilerr", "nilnesserr", "nilnil", "noctx", "nonamedreturns", "paralleltest", "perfsprint", "protogetter", "reassign", "recvcheck", "revive", "rowserrcheck", "sloglint", "spancheck", "sqlclosecheck", "staticcheck", "stylecheck", "tagliatelle", "testifylint", "thelper", "tparallel", "unconvert", "unparam", "unused", "usetesting", "varnamelen", "wastedassign", "wrapcheck", "zerologlint"},
		},
		{
			desc: "disable",
			old: versionone.Linters{
				EnableAll: ptr.Pointer(true),
				Disable:   []string{"lll", "misspell", "govet"},
				Fast:      ptr.Pointer(false),
				Presets:   nil,
			},
			expected: []string{"govet", "lll", "misspell"},
		},
		{
			desc: "disable, fast",
			old: versionone.Linters{
				EnableAll: ptr.Pointer(true),
				Disable:   []string{"lll", "misspell", "govet"},
				Fast:      ptr.Pointer(true),
				Presets:   nil,
			},
			expected: []string{"asasalint", "bodyclose", "canonicalheader", "containedctx", "contextcheck", "durationcheck", "err113", "errcheck", "errchkjson", "errname", "errorlint", "exhaustive", "exhaustruct", "exptostd", "fatcontext", "forbidigo", "forcetypeassert", "ginkgolinter", "gochecknoglobals", "gochecksumtype", "gocritic", "gosec", "gosimple", "gosmopolitan", "govet", "iface", "importas", "intrange", "ireturn", "lll", "loggercheck", "makezero", "mirror", "misspell", "musttag", "nilerr", "nilnesserr", "nilnil", "noctx", "nonamedreturns", "paralleltest", "perfsprint", "protogetter", "reassign", "recvcheck", "revive", "rowserrcheck", "sloglint", "spancheck", "sqlclosecheck", "staticcheck", "stylecheck", "tagliatelle", "testifylint", "thelper", "tparallel", "unconvert", "unparam", "unused", "usetesting", "varnamelen", "wastedassign", "wrapcheck", "zerologlint"},
		},
		{
			desc: "disable, enable, fast",
			old: versionone.Linters{
				EnableAll: ptr.Pointer(true),
				Enable:    []string{"canonicalheader", "errname"},
				Disable:   []string{"lll", "misspell", "govet"},
				Fast:      ptr.Pointer(true),
				Presets:   nil,
			},
			expected: []string{"asasalint", "bodyclose", "containedctx", "contextcheck", "durationcheck", "err113", "errcheck", "errchkjson", "errorlint", "exhaustive", "exhaustruct", "exptostd", "fatcontext", "forbidigo", "forcetypeassert", "ginkgolinter", "gochecknoglobals", "gochecksumtype", "gocritic", "gosec", "gosimple", "gosmopolitan", "govet", "iface", "importas", "intrange", "ireturn", "lll", "loggercheck", "makezero", "mirror", "misspell", "musttag", "nilerr", "nilnesserr", "nilnil", "noctx", "nonamedreturns", "paralleltest", "perfsprint", "protogetter", "reassign", "recvcheck", "revive", "rowserrcheck", "sloglint", "spancheck", "sqlclosecheck", "staticcheck", "stylecheck", "tagliatelle", "testifylint", "thelper", "tparallel", "unconvert", "unparam", "unused", "usetesting", "varnamelen", "wastedassign", "wrapcheck", "zerologlint"},
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			results := enableAllFilter(test.old)

			assert.Equal(t, test.expected, results)
		})
	}
}

func Test_defaultLintersDisableFilter(t *testing.T) {
	testCases := []struct {
		desc     string
		old      versionone.Linters
		expected []string
	}{
		{
			desc: "no options",
			old: versionone.Linters{
				Enable:  nil,
				Disable: nil,
				Fast:    ptr.Pointer(false),
				Presets: nil,
			},
			expected: nil,
		},
		{
			desc: "presets (ignored)",
			old: versionone.Linters{
				Enable:  nil,
				Disable: nil,
				Fast:    ptr.Pointer(false),
				Presets: []string{"comment", "error", "format"},
			},
			expected: nil,
		},
		{
			desc: "fast",
			old: versionone.Linters{
				Enable:  nil,
				Disable: nil,
				Fast:    ptr.Pointer(true),
				Presets: nil,
			},
			expected: []string{"errcheck", "gosimple", "govet", "staticcheck", "unused"},
		},
		{
			desc: "enable",
			old: versionone.Linters{
				Enable:  []string{"lll", "misspell", "govet"},
				Disable: nil,
				Fast:    ptr.Pointer(false),
				Presets: nil,
			},
			expected: nil,
		},
		{
			desc: "disable",
			old: versionone.Linters{
				Enable:  nil,
				Disable: []string{"lll", "misspell", "govet"},
				Fast:    ptr.Pointer(false),
				Presets: nil,
			},
			expected: []string{"govet", "lll", "misspell"},
		},
		{
			desc: "disable, fast",
			old: versionone.Linters{
				Enable:  nil,
				Disable: []string{"lll", "misspell", "govet"},
				Fast:    ptr.Pointer(true),
				Presets: nil,
			},
			expected: []string{"errcheck", "gosimple", "govet", "lll", "misspell", "staticcheck", "unused"},
		},
		{
			desc: "enable, disable",
			old: versionone.Linters{
				Enable:  []string{"grouper", "importas", "errcheck"},
				Disable: []string{"lll", "misspell", "govet"},
				Fast:    ptr.Pointer(false),
				Presets: nil,
			},
			expected: []string{"govet", "lll", "misspell"},
		},
		{
			desc: "enable",
			old: versionone.Linters{
				Enable:  []string{"grouper", "importas", "errcheck"},
				Disable: nil,
				Fast:    ptr.Pointer(false),
				Presets: nil,
			},
			expected: nil,
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			results := toNames(defaultLintersDisableFilter(test.old))

			assert.Equal(t, test.expected, results)
		})
	}
}

func Test_defaultLintersEnableFilter(t *testing.T) {
	testCases := []struct {
		desc string

		old      versionone.Linters
		expected []string
	}{
		{
			desc: "no options",
			old: versionone.Linters{
				Enable:  nil,
				Disable: nil,
				Fast:    ptr.Pointer(false),
				Presets: nil,
			},
			expected: nil,
		},
		{
			desc: "enable",
			old: versionone.Linters{
				Enable:  []string{"grouper", "importas", "errcheck"},
				Disable: nil,
				Fast:    ptr.Pointer(false),
				Presets: nil,
			},
			expected: []string{"grouper", "importas"},
		},
		{
			desc: "enable, disable",
			old: versionone.Linters{
				Enable:  []string{"grouper", "importas", "errcheck"},
				Disable: []string{"lll", "misspell", "govet"},
				Fast:    ptr.Pointer(false),
				Presets: nil,
			},
			expected: []string{"grouper", "importas"},
		},
		{
			desc: "disable",
			old: versionone.Linters{
				Enable:  nil,
				Disable: []string{"lll", "misspell", "govet"},
				Fast:    ptr.Pointer(false),
				Presets: nil,
			},
			expected: nil,
		},
		{
			desc: "presets",
			old: versionone.Linters{
				Enable:  nil,
				Disable: nil,
				Fast:    ptr.Pointer(false),
				Presets: []string{"comment", "error", "format"},
			},
			expected: []string{"dupword", "err113", "errorlint", "gci", "godot", "godox", "gofmt", "gofumpt", "goimports", "misspell", "wrapcheck"},
		},
		{
			desc: "presets, fast",
			old: versionone.Linters{
				Enable:  nil,
				Disable: nil,
				Fast:    ptr.Pointer(true),
				Presets: []string{"comment", "error", "format"},
			},
			expected: []string{"dupword", "gci", "godot", "godox", "gofmt", "gofumpt", "goimports", "misspell"},
		},
	}

	// presets - slow + enable - default - [effective disable] => effective enable
	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			results := toNames(defaultLintersEnableFilter(test.old, defaultLintersDisableFilter(test.old)))

			assert.Equal(t, test.expected, results)
		})
	}
}

func Test_convertStaticcheckLinterNames(t *testing.T) {
	testCases := []struct {
		desc     string
		names    []string
		expected []string
	}{
		{
			desc:     "empty",
			names:    nil,
			expected: nil,
		},
		{
			desc:     "no staticcheck linters",
			names:    []string{"lll", "misspell", "govet"},
			expected: []string{"govet", "lll", "misspell"},
		},
		{
			desc:     "stylecheck",
			names:    []string{"lll", "misspell", "govet", "stylecheck"},
			expected: []string{"govet", "lll", "misspell", "staticcheck"},
		},
		{
			desc:     "gosimple",
			names:    []string{"lll", "misspell", "govet", "gosimple"},
			expected: []string{"govet", "lll", "misspell", "staticcheck"},
		},
		{
			desc:     "staticcheck",
			names:    []string{"lll", "misspell", "govet", "staticcheck"},
			expected: []string{"govet", "lll", "misspell", "staticcheck"},
		},
		{
			desc:     "staticcheck, stylecheck, gosimple",
			names:    []string{"lll", "misspell", "govet", "staticcheck", "stylecheck", "gosimple"},
			expected: []string{"govet", "lll", "misspell", "staticcheck"},
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			results := convertStaticcheckLinterNames(test.names)

			assert.Equal(t, test.expected, results)
		})
	}
}

func Test_unknownLinterNames(t *testing.T) {
	testCases := []struct {
		desc     string
		names    []string
		expected []string
	}{
		{
			desc:     "empty",
			names:    nil,
			expected: nil,
		},
		{
			desc:     "deprecated",
			names:    []string{"golint", "structcheck", "varcheck"},
			expected: nil,
		},
		{
			desc:     "deprecated and unknown",
			names:    []string{"golint", "structcheck", "varcheck", "a", "b"},
			expected: []string{"a", "b"},
		},
		{
			desc:     "deprecated and known",
			names:    []string{"golint", "structcheck", "varcheck", "gosec", "gofmt"},
			expected: nil,
		},
		{
			desc:     "only unknown",
			names:    []string{"a", "b", "c"},
			expected: []string{"a", "b", "c"},
		},
		{
			desc:     "unknown and known",
			names:    []string{"a", "gosec", "gofmt"},
			expected: []string{"a"},
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			names := unknownLinterNames(test.names, allLinters())

			assert.Equal(t, test.expected, names)
		})
	}
}
