// Package env implements the Pipe interface providing validation of
// missing environment variables needed by the release process.
package env

import (
	"bufio"
	"os"

	"github.com/goreleaser/goreleaser/internal/pipe"
	"github.com/goreleaser/goreleaser/pkg/context"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
)

// ErrMissingToken indicates an error when GITHUB_TOKEN, GITLAB_TOKEN and GITEA_TOKEN are all missing in the environment
var ErrMissingToken = errors.New("missing GITHUB_TOKEN, GITLAB_TOKEN and GITEA_TOKEN")

// ErrMultipleTokens indicates that multiple tokens are defined. ATM only one of them if allowed
// See https://github.com/goreleaser/goreleaser/pull/809
var ErrMultipleTokens = errors.New("multiple tokens defined. Only one is allowed")

// Pipe for env
type Pipe struct{}

func (Pipe) String() string {
	return "loading environment variables"
}

// Default sets the pipe defaults
func (Pipe) Default(ctx *context.Context) error {
	var env = &ctx.Config.EnvFiles
	if env.GitHubToken == "" {
		env.GitHubToken = "~/.config/goreleaser/github_token"
	}
	if env.GitLabToken == "" {
		env.GitLabToken = "~/.config/goreleaser/gitlab_token"
	}
	if env.GiteaToken == "" {
		env.GiteaToken = "~/.config/goreleaser/gitea_token"
	}
	return nil
}

// Run the pipe
func (Pipe) Run(ctx *context.Context) error {
	githubToken, githubTokenErr := loadEnv("GITHUB_TOKEN", ctx.Config.EnvFiles.GitHubToken)
	gitlabToken, gitlabTokenErr := loadEnv("GITLAB_TOKEN", ctx.Config.EnvFiles.GitLabToken)
	giteaToken, giteaTokenErr := loadEnv("GITEA_TOKEN", ctx.Config.EnvFiles.GiteaToken)

	if ctx.SkipPublish {
		return pipe.ErrSkipPublishEnabled
	}
	if ctx.Config.Release.Disable {
		return pipe.Skip("release pipe is disabled")
	}

	numOfTokens := 0
	if githubToken != "" {
		numOfTokens++
	}
	if gitlabToken != "" {
		numOfTokens++
	}
	if giteaToken != "" {
		numOfTokens++
	}
	if numOfTokens > 1 {
		return ErrMultipleTokens
	}

	noTokens := githubToken == "" && gitlabToken == "" && giteaToken == ""
	noTokenErrs := githubTokenErr == nil && gitlabTokenErr == nil && giteaTokenErr == nil

	if noTokens && noTokenErrs {
		return ErrMissingToken
	}

	if gitlabTokenErr != nil {
		return errors.Wrap(gitlabTokenErr, "failed to load gitlab token")
	}

	if githubTokenErr != nil {
		return errors.Wrap(githubTokenErr, "failed to load github token")
	}

	if giteaTokenErr != nil {
		return errors.Wrap(giteaTokenErr, "failed to load gitea token")
	}

	if githubToken != "" {
		ctx.TokenType = context.TokenTypeGitHub
		ctx.Token = githubToken
	}

	if gitlabToken != "" {
		ctx.TokenType = context.TokenTypeGitLab
		ctx.Token = gitlabToken
	}

	if giteaToken != "" {
		ctx.TokenType = context.TokenTypeGitea
		ctx.Token = giteaToken
	}

	return nil
}

func loadEnv(env, path string) (string, error) {
	val := os.Getenv(env)
	if val != "" {
		return val, nil
	}
	path, err := homedir.Expand(path)
	if err != nil {
		return "", err
	}
	f, err := os.Open(path) // #nosec
	if os.IsNotExist(err) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	bts, _, err := bufio.NewReader(f).ReadLine()
	return string(bts), err
}
