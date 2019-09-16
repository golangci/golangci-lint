package client

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"code.gitea.io/gitea/modules/structs"
	"code.gitea.io/sdk/gitea"
	"github.com/apex/log"
	"github.com/goreleaser/goreleaser/internal/artifact"
	"github.com/goreleaser/goreleaser/internal/tmpl"
	"github.com/goreleaser/goreleaser/pkg/config"
	"github.com/goreleaser/goreleaser/pkg/context"
)

type giteaClient struct {
	client *gitea.Client
}

func getInstanceURL(apiURL string) (string, error) {
	u, err := url.Parse(apiURL)
	if err != nil {
		return "", err
	}
	u.Path = ""
	rawurl := u.String()
	if rawurl == "" {
		return "", fmt.Errorf("invalid URL: %v", apiURL)
	}
	return rawurl, nil
}

// NewGitea returns a gitea client implementation
func NewGitea(ctx *context.Context) (Client, error) {
	instanceURL, err := getInstanceURL(ctx.Config.GiteaURLs.API)
	if err != nil {
		return nil, err
	}
	client := gitea.NewClient(instanceURL, ctx.Token)
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			// nolint: gosec
			InsecureSkipVerify: ctx.Config.GiteaURLs.SkipTLSVerify,
		},
	}
	httpClient := &http.Client{Transport: transport}
	client.SetHTTPClient(httpClient)
	return &giteaClient{client: client}, nil
}

// CreateFile creates a file in the repository at a given path
// or updates the file if it exists
func (c *giteaClient) CreateFile(
	ctx *context.Context,
	commitAuthor config.CommitAuthor,
	repo config.Repo,
	content []byte,
	path,
	message string,
) error {
	//TODO: implement for brew and scoop support for Gitea-hosted repos
	return nil
}

func (c *giteaClient) createRelease(ctx *context.Context, title, body string) (*gitea.Release, error) {
	releaseConfig := ctx.Config.Release
	owner := releaseConfig.Gitea.Owner
	repoName := releaseConfig.Gitea.Name
	tag := ctx.Git.CurrentTag

	opts := structs.CreateReleaseOption{
		TagName:      tag,
		Target:       ctx.Git.Commit,
		Title:        title,
		Note:         body,
		IsDraft:      releaseConfig.Draft,
		IsPrerelease: ctx.PreRelease,
	}
	release, err := c.client.CreateRelease(owner, repoName, opts)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
		}).Debug("error creating Gitea release")
		return nil, err
	}
	log.WithField("id", release.ID).Info("Gitea release created")
	return release, nil
}

func (c *giteaClient) getExistingRelease(owner, repoName, tagName string) (*gitea.Release, error) {
	releases, err := c.client.ListReleases(owner, repoName)
	if err != nil {
		return nil, err
	}

	for _, release := range releases {
		if release.TagName == tagName {
			return release, nil
		}
	}

	return nil, nil
}

func (c *giteaClient) updateRelease(ctx *context.Context, title, body string, id int64) (*gitea.Release, error) {
	releaseConfig := ctx.Config.Release
	owner := releaseConfig.Gitea.Owner
	repoName := releaseConfig.Gitea.Name
	tag := ctx.Git.CurrentTag

	opts := structs.EditReleaseOption{
		TagName:      tag,
		Target:       ctx.Git.Commit,
		Title:        title,
		Note:         body,
		IsDraft:      &releaseConfig.Draft,
		IsPrerelease: &ctx.PreRelease,
	}

	release, err := c.client.EditRelease(owner, repoName, id, opts)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
		}).Debug("error updating Gitea release")
		return nil, err
	}
	log.WithField("id", release.ID).Info("Gitea release updated")
	return release, nil
}

// CreateRelease creates a new release or updates it by keeping
// the release notes if it exists
func (c *giteaClient) CreateRelease(ctx *context.Context, body string) (string, error) {
	var release *gitea.Release
	var err error

	releaseConfig := ctx.Config.Release

	title, err := tmpl.New(ctx).Apply(releaseConfig.NameTemplate)
	if err != nil {
		return "", err
	}

	release, err = c.getExistingRelease(
		releaseConfig.Gitea.Owner,
		releaseConfig.Gitea.Name,
		ctx.Git.CurrentTag,
	)
	if err != nil {
		return "", err
	}

	if release != nil {
		release, err = c.updateRelease(ctx, title, body, release.ID)
		if err != nil {
			return "", err
		}
	} else {
		release, err = c.createRelease(ctx, title, body)
		if err != nil {
			return "", err
		}
	}

	return strconv.FormatInt(release.ID, 10), nil
}

// Upload uploads a file into a release repository
func (c *giteaClient) Upload(
	ctx *context.Context,
	releaseID string,
	artifact *artifact.Artifact,
	file *os.File,
) error {
	giteaReleaseID, err := strconv.ParseInt(releaseID, 10, 64)
	if err != nil {
		return err
	}
	releaseConfig := ctx.Config.Release
	owner := releaseConfig.Gitea.Owner
	repoName := releaseConfig.Gitea.Name

	_, err = c.client.CreateReleaseAttachment(owner, repoName, giteaReleaseID, file, artifact.Name)
	return err
}
