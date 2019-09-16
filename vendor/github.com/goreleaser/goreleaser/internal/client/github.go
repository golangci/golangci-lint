package client

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strconv"

	"github.com/apex/log"
	"github.com/google/go-github/v25/github"
	"github.com/goreleaser/goreleaser/internal/artifact"
	"github.com/goreleaser/goreleaser/internal/tmpl"
	"github.com/goreleaser/goreleaser/pkg/config"
	"github.com/goreleaser/goreleaser/pkg/context"
	"golang.org/x/oauth2"
)

type githubClient struct {
	client *github.Client
}

// NewGitHub returns a github client implementation
func NewGitHub(ctx *context.Context) (Client, error) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: ctx.Token},
	)
	httpClient := oauth2.NewClient(ctx, ts)
	base := httpClient.Transport.(*oauth2.Transport).Base
	// nolint: govet
	if base == nil || reflect.ValueOf(base).IsNil() {
		base = http.DefaultTransport
	}
	// nolint: gosec
	base.(*http.Transport).TLSClientConfig = &tls.Config{
		InsecureSkipVerify: ctx.Config.GitHubURLs.SkipTLSVerify,
	}
	httpClient.Transport.(*oauth2.Transport).Base = base
	client := github.NewClient(httpClient)
	if ctx.Config.GitHubURLs.API != "" {
		api, err := url.Parse(ctx.Config.GitHubURLs.API)
		if err != nil {
			return &githubClient{}, err
		}
		upload, err := url.Parse(ctx.Config.GitHubURLs.Upload)
		if err != nil {
			return &githubClient{}, err
		}
		client.BaseURL = api
		client.UploadURL = upload
	}

	return &githubClient{client: client}, nil
}

func (c *githubClient) CreateFile(
	ctx *context.Context,
	commitAuthor config.CommitAuthor,
	repo config.Repo,
	content []byte,
	path,
	message string,
) error {
	options := &github.RepositoryContentFileOptions{
		Committer: &github.CommitAuthor{
			Name:  github.String(commitAuthor.Name),
			Email: github.String(commitAuthor.Email),
		},
		Content: content,
		Message: github.String(message),
	}

	file, _, res, err := c.client.Repositories.GetContents(
		ctx,
		repo.Owner,
		repo.Name,
		path,
		&github.RepositoryContentGetOptions{},
	)
	if err != nil && res.StatusCode != 404 {
		return err
	}

	if res.StatusCode == 404 {
		_, _, err = c.client.Repositories.CreateFile(
			ctx,
			repo.Owner,
			repo.Name,
			path,
			options,
		)
		return err
	}
	options.SHA = file.SHA
	_, _, err = c.client.Repositories.UpdateFile(
		ctx,
		repo.Owner,
		repo.Name,
		path,
		options,
	)
	return err
}

func (c *githubClient) CreateRelease(ctx *context.Context, body string) (string, error) {
	var release *github.RepositoryRelease
	title, err := tmpl.New(ctx).Apply(ctx.Config.Release.NameTemplate)
	if err != nil {
		return "", err
	}

	var data = &github.RepositoryRelease{
		Name:       github.String(title),
		TagName:    github.String(ctx.Git.CurrentTag),
		Body:       github.String(body),
		Draft:      github.Bool(ctx.Config.Release.Draft),
		Prerelease: github.Bool(ctx.PreRelease),
	}
	release, _, err = c.client.Repositories.GetReleaseByTag(
		ctx,
		ctx.Config.Release.GitHub.Owner,
		ctx.Config.Release.GitHub.Name,
		ctx.Git.CurrentTag,
	)
	if err != nil {
		release, _, err = c.client.Repositories.CreateRelease(
			ctx,
			ctx.Config.Release.GitHub.Owner,
			ctx.Config.Release.GitHub.Name,
			data,
		)
	} else {
		// keep the pre-existing release notes
		if release.GetBody() != "" {
			data.Body = release.Body
		}
		release, _, err = c.client.Repositories.EditRelease(
			ctx,
			ctx.Config.Release.GitHub.Owner,
			ctx.Config.Release.GitHub.Name,
			release.GetID(),
			data,
		)
	}
	log.WithField("url", release.GetHTMLURL()).Info("release updated")
	githubReleaseID := strconv.FormatInt(release.GetID(), 10)
	return githubReleaseID, err
}

func (c *githubClient) Upload(
	ctx *context.Context,
	releaseID string,
	artifact *artifact.Artifact,
	file *os.File,
) error {
	githubReleaseID, err := strconv.ParseInt(releaseID, 10, 64)
	if err != nil {
		return err
	}
	_, _, err = c.client.Repositories.UploadReleaseAsset(
		ctx,
		ctx.Config.Release.GitHub.Owner,
		ctx.Config.Release.GitHub.Name,
		githubReleaseID,
		&github.UploadOptions{
			Name: artifact.Name,
		},
		file,
	)
	return err
}
