package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_sanitizeVersion(t *testing.T) {
	testCases := []struct {
		desc     string
		input    string
		expected string
	}{
		{
			desc:     "ampersand",
			input:    " te&st",
			expected: "test",
		},
		{
			desc:     "pipe",
			input:    " te|st",
			expected: "test",
		},
		{
			desc:     "version",
			input:    "v1.2.3",
			expected: "v1.2.3",
		},
		{
			desc:     "branch",
			input:    "feat/test",
			expected: "feat/test",
		},
		{
			desc:     "branch",
			input:    "value --key",
			expected: "valuekey",
		},
		{
			desc:     "hash",
			input:    "cd8b1177",
			expected: "cd8b1177",
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			v := sanitizeVersion(test.input)

			assert.Equal(t, test.expected, v)
		})
	}
}

func Test_filterGitEnviron(t *testing.T) {
	environ := []string{
		"GIT_ALLOW_PROTOCOL=https:git:ssh",
		"GIT_ASKPASS=/usr/bin/ssh-askpass",
		"GIT_AUTHOR_DATE=2026-04-20T17:52:14+02:00",
		"GIT_AUTHOR_EMAIL=jane.doe@example.com",
		"GIT_AUTHOR_NAME=Jane Doe",
		"GIT_COMMITTER_DATE=2026-04-20T17:52:14+02:00",
		"GIT_COMMITTER_EMAIL=john.smith@example.com",
		"GIT_COMMITTER_NAME=John Smith",
		"GIT_CONFIG_COUNT=2",
		"GIT_CONFIG_KEY_0=http.sslVerify",
		"GIT_CONFIG_KEY_1=user.email",
		"GIT_CONFIG_VALUE_0=false",
		"GIT_CONFIG_VALUE_1=bot@example.com",
		"GIT_DIR=/home/jane/project/.git",
		"GIT_EXEC_PATH=/usr/lib/git-core",
		"GIT_HTTP_PROXY_AUTHMETHOD=basic",
		"GIT_INDEX_FILE=/home/jane/project/.git/index",
		"GIT_SSH_COMMAND=ssh -i ~/.ssh/id_ed25519 -o StrictHostKeyChecking=no",
		"GIT_SSH=/usr/bin/ssh",
		"GIT_SSL_CAINFO=/etc/ssl/certs/ca-certificates.crt",
		"GIT_SSL_NO_VERIFY=true",
		"GIT_TERMINAL_PROMPT=0",
	}

	envs := filterGitEnviron(environ)

	expected := []string{
		"GIT_ALLOW_PROTOCOL=https:git:ssh",
		"GIT_ASKPASS=/usr/bin/ssh-askpass",
		"GIT_CONFIG_COUNT=2",
		"GIT_CONFIG_KEY_0=http.sslVerify",
		"GIT_CONFIG_KEY_1=user.email",
		"GIT_CONFIG_VALUE_0=false",
		"GIT_CONFIG_VALUE_1=bot@example.com",
		"GIT_EXEC_PATH=/usr/lib/git-core",
		"GIT_HTTP_PROXY_AUTHMETHOD=basic",
		"GIT_SSH_COMMAND=ssh -i ~/.ssh/id_ed25519 -o StrictHostKeyChecking=no",
		"GIT_SSH=/usr/bin/ssh",
		"GIT_SSL_CAINFO=/etc/ssl/certs/ca-certificates.crt",
		"GIT_SSL_NO_VERIFY=true",
	}

	assert.Equal(t, expected, envs)
}
