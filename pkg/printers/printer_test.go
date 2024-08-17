package printers

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/golangci/golangci-lint/pkg/report"
	"github.com/golangci/golangci-lint/pkg/result"
)

func unmarshalFile(t *testing.T, filename string, v any) {
	t.Helper()

	file, err := os.ReadFile(filepath.Join("testdata", filename))
	require.NoError(t, err)

	err = json.Unmarshal(file, v)
	require.NoError(t, err)
}

func TestPrinter_Print_stdout(t *testing.T) {
	logger := logutils.NewStderrLog("skip")

	var issues []result.Issue
	unmarshalFile(t, "in-issues.json", &issues)

	data := &report.Data{}
	unmarshalFile(t, "in-report-data.json", data)

	testCases := []struct {
		desc     string
		cfg      *config.Output
		expected string
	}{
		{
			desc: "stdout (implicit)",
			cfg: &config.Output{
				Formats: []config.OutputFormat{
					{Format: "line-number"},
				},
			},
			expected: "golden-line-number.txt",
		},
		{
			desc: "stdout (explicit)",
			cfg: &config.Output{
				Formats: []config.OutputFormat{
					{
						Format: "line-number",
						Path:   "stdout",
					},
				},
			},
			expected: "golden-line-number.txt",
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			p, err := NewPrinter(logger, test.cfg, data)
			require.NoError(t, err)

			var stdOutBuffer bytes.Buffer
			p.stdOut = &stdOutBuffer

			var stdErrBuffer bytes.Buffer
			p.stdErr = &stdErrBuffer

			err = p.Print(issues)
			require.NoError(t, err)

			golden, err := os.ReadFile(filepath.Join("testdata", test.expected))
			require.NoError(t, err)

			assert.Equal(t, 0, stdErrBuffer.Len())
			assert.Equal(t, string(golden), stdOutBuffer.String())
		})
	}
}

func TestPrinter_Print_stderr(t *testing.T) {
	logger := logutils.NewStderrLog("skip")

	var issues []result.Issue
	unmarshalFile(t, "in-issues.json", &issues)

	data := &report.Data{}
	unmarshalFile(t, "in-report-data.json", data)

	cfg := &config.Output{
		Formats: []config.OutputFormat{
			{
				Format: "line-number",
				Path:   "stderr",
			},
		},
	}

	p, err := NewPrinter(logger, cfg, data)
	require.NoError(t, err)

	var stdOutBuffer bytes.Buffer
	p.stdOut = &stdOutBuffer

	var stdErrBuffer bytes.Buffer
	p.stdErr = &stdErrBuffer

	err = p.Print(issues)
	require.NoError(t, err)

	golden, err := os.ReadFile(filepath.Join("testdata", "golden-line-number.txt"))
	require.NoError(t, err)

	assert.Equal(t, 0, stdOutBuffer.Len())
	assert.Equal(t, string(golden), stdErrBuffer.String())
}

func TestPrinter_Print_file(t *testing.T) {
	logger := logutils.NewStderrLog("skip")

	var issues []result.Issue
	unmarshalFile(t, "in-issues.json", &issues)

	data := &report.Data{}
	unmarshalFile(t, "in-report-data.json", data)

	outputPath := filepath.Join(t.TempDir(), "report.txt")

	cfg := &config.Output{
		Formats: []config.OutputFormat{
			{
				Format: "line-number",
				Path:   outputPath,
			},
		},
	}

	p, err := NewPrinter(logger, cfg, data)
	require.NoError(t, err)

	var stdOutBuffer bytes.Buffer
	p.stdOut = &stdOutBuffer

	var stdErrBuffer bytes.Buffer
	p.stdErr = &stdErrBuffer

	err = p.Print(issues)
	require.NoError(t, err)

	golden, err := os.ReadFile(filepath.Join("testdata", "golden-line-number.txt"))
	require.NoError(t, err)

	assert.Equal(t, 0, stdOutBuffer.Len())
	assert.Equal(t, 0, stdErrBuffer.Len())

	actual, err := os.ReadFile(outputPath)
	require.NoError(t, err)

	assert.Equal(t, string(golden), string(actual))
}

func TestPrinter_Print_multiple(t *testing.T) {
	logger := logutils.NewStderrLog("skip")

	var issues []result.Issue
	unmarshalFile(t, "in-issues.json", &issues)

	data := &report.Data{}
	unmarshalFile(t, "in-report-data.json", data)

	outputPath := filepath.Join(t.TempDir(), "teamcity.txt")

	cfg := &config.Output{
		Formats: []config.OutputFormat{
			{
				Format: "teamcity",
				Path:   outputPath,
			},
			{
				Format: "json",
				Path:   "",
			},
			{
				Format: "line-number",
				Path:   "stderr",
			},
		},
	}

	p, err := NewPrinter(logger, cfg, data)
	require.NoError(t, err)

	var stdOutBuffer bytes.Buffer
	p.stdOut = &stdOutBuffer

	var stdErrBuffer bytes.Buffer
	p.stdErr = &stdErrBuffer

	err = p.Print(issues)
	require.NoError(t, err)

	goldenGitHub, err := os.ReadFile(filepath.Join("testdata", "golden-teamcity.txt"))
	require.NoError(t, err)

	actual, err := os.ReadFile(outputPath)
	require.NoError(t, err)

	assert.Equal(t, string(goldenGitHub), string(actual))

	goldenLineNumber, err := os.ReadFile(filepath.Join("testdata", "golden-line-number.txt"))
	require.NoError(t, err)

	assert.Equal(t, string(goldenLineNumber), stdErrBuffer.String())

	goldenJSON, err := os.ReadFile(filepath.Join("testdata", "golden-json.json"))
	require.NoError(t, err)

	assert.Equal(t, string(goldenJSON), stdOutBuffer.String())
}
