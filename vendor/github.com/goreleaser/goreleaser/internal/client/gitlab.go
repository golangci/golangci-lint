package client

import (
	"crypto/tls"
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/apex/log"
	"github.com/goreleaser/goreleaser/internal/artifact"
	"github.com/goreleaser/goreleaser/internal/tmpl"
	"github.com/goreleaser/goreleaser/pkg/config"
	"github.com/goreleaser/goreleaser/pkg/context"
	"github.com/xanzy/go-gitlab"
)

// ErrExtractHashFromFileUploadURL indicates the file upload hash could not ne extracted from the url
var ErrExtractHashFromFileUploadURL = errors.New("could not extract hash from gitlab file upload url")

type gitlabClient struct {
	client *gitlab.Client
}

// NewGitLab returns a gitlab client implementation
func NewGitLab(ctx *context.Context) (Client, error) {
	token := ctx.Token
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			// nolint: gosec
			InsecureSkipVerify: ctx.Config.GitLabURLs.SkipTLSVerify,
		},
	}
	httpClient := &http.Client{Transport: transport}
	client := gitlab.NewClient(httpClient, token)
	if ctx.Config.GitLabURLs.API != "" {
		err := client.SetBaseURL(ctx.Config.GitLabURLs.API)
		if err != nil {
			return &gitlabClient{}, err
		}
	}
	return &gitlabClient{client: client}, nil
}

// CreateFile gets a file in the repository at a given path
// and updates if it exists or creates it for later pipes in the pipeline
func (c *gitlabClient) CreateFile(
	ctx *context.Context,
	commitAuthor config.CommitAuthor,
	repo config.Repo,
	content []byte, // the content of the formula.rb
	path, // the path to the formula.rb
	message string, // the commit msg
) error {
	fileName := path
	// we assume having the formula in the master branch only
	ref := "master"
	branch := "master"
	opts := &gitlab.GetFileOptions{Ref: &ref}
	castedContent := string(content)
	projectID := repo.Owner + "/" + repo.Name

	log.WithFields(log.Fields{
		"owner": repo.Owner,
		"name":  repo.Name,
	}).Debug("projectID at brew")

	_, res, err := c.client.RepositoryFiles.GetFile(projectID, fileName, opts)
	if err != nil && (res == nil || res.StatusCode != 404) {
		log.WithFields(log.Fields{
			"fileName":   fileName,
			"ref":        ref,
			"projectID":  projectID,
			"statusCode": res.StatusCode,
			"err":        err.Error(),
		}).Error("error getting file for brew formula")
		return err
	}

	log.WithFields(log.Fields{
		"fileName":  fileName,
		"branch":    branch,
		"projectID": projectID,
	}).Debug("found already existing brew formula file")

	if res.StatusCode == 404 {
		log.WithFields(log.Fields{
			"fileName":  fileName,
			"ref":       ref,
			"projectID": projectID,
		}).Debug("creating brew formula")
		createOpts := &gitlab.CreateFileOptions{
			AuthorName:    &commitAuthor.Name,
			AuthorEmail:   &commitAuthor.Email,
			Content:       &castedContent,
			Branch:        &branch,
			CommitMessage: &message,
		}
		fileInfo, res, err := c.client.RepositoryFiles.CreateFile(projectID, fileName, createOpts)
		if err != nil {
			log.WithFields(log.Fields{
				"fileName":   fileName,
				"branch":     branch,
				"projectID":  projectID,
				"statusCode": res.StatusCode,
				"err":        err.Error(),
			}).Error("error creating brew formula file")
			return err
		}

		log.WithFields(log.Fields{
			"fileName":  fileName,
			"branch":    branch,
			"projectID": projectID,
			"filePath":  fileInfo.FilePath,
		}).Debug("created brew formula file")
		return nil
	}

	log.WithFields(log.Fields{
		"fileName":  fileName,
		"ref":       ref,
		"projectID": projectID,
	}).Debug("updating brew formula")
	updateOpts := &gitlab.UpdateFileOptions{
		AuthorName:    &commitAuthor.Name,
		AuthorEmail:   &commitAuthor.Email,
		Content:       &castedContent,
		Branch:        &branch,
		CommitMessage: &message,
	}

	updateFileInfo, res, err := c.client.RepositoryFiles.UpdateFile(projectID, fileName, updateOpts)
	if err != nil {
		log.WithFields(log.Fields{
			"fileName":   fileName,
			"branch":     branch,
			"projectID":  projectID,
			"statusCode": res.StatusCode,
			"err":        err.Error(),
		}).Error("error updating brew formula file")
		return err
	}

	log.WithFields(log.Fields{
		"fileName":   fileName,
		"branch":     branch,
		"projectID":  projectID,
		"filePath":   updateFileInfo.FilePath,
		"statusCode": res.StatusCode,
	}).Debug("updated brew formula file")
	return nil
}

// CreateRelease creates a new release or updates it by keeping
// the release notes if it exists
func (c *gitlabClient) CreateRelease(ctx *context.Context, body string) (releaseID string, err error) {
	title, err := tmpl.New(ctx).Apply(ctx.Config.Release.NameTemplate)
	if err != nil {
		return "", err
	}

	projectID := ctx.Config.Release.GitLab.Owner + "/" + ctx.Config.Release.GitLab.Name
	log.WithFields(log.Fields{
		"owner": ctx.Config.Release.GitLab.Owner,
		"name":  ctx.Config.Release.GitLab.Name,
	}).Debug("projectID")

	name := title
	tagName := ctx.Git.CurrentTag
	release, resp, err := c.client.Releases.GetRelease(projectID, tagName)
	if err != nil && (resp == nil || resp.StatusCode != 403) {
		return "", err
	}

	if resp.StatusCode == 403 {
		log.WithFields(log.Fields{
			"err": err.Error(),
		}).Debug("get release")

		description := body
		ref := ctx.Git.Commit
		gitURL := ctx.Git.URL

		log.WithFields(log.Fields{
			"name":        name,
			"description": description,
			"ref":         ref,
			"url":         gitURL,
		}).Debug("creating release")
		release, _, err = c.client.Releases.CreateRelease(projectID, &gitlab.CreateReleaseOptions{
			Name:        &name,
			Description: &description,
			Ref:         &ref,
			TagName:     &tagName,
		})

		if err != nil {
			log.WithFields(log.Fields{
				"err": err.Error(),
			}).Debug("error create release")
			return "", err
		}
		log.WithField("name", release.Name).Info("release created")
	} else {
		desc := body
		if release != nil && release.DescriptionHTML != "" {
			desc = release.DescriptionHTML
		}

		release, _, err = c.client.Releases.UpdateRelease(projectID, tagName, &gitlab.UpdateReleaseOptions{
			Name:        &name,
			Description: &desc,
		})
		if err != nil {
			log.WithFields(log.Fields{
				"err": err.Error(),
			}).Debug("error update release")
			return "", err
		}

		log.WithField("name", release.Name).Info("release updated")
	}

	return tagName, err // gitlab references a tag in a repo by its name
}

// Upload uploads a file into a release repository
func (c *gitlabClient) Upload(
	ctx *context.Context,
	releaseID string,
	artifact *artifact.Artifact,
	file *os.File,
) error {
	projectID := ctx.Config.Release.GitLab.Owner + "/" + ctx.Config.Release.GitLab.Name

	log.WithField("file", file.Name()).Debug("uploading file")
	projectFile, _, err := c.client.Projects.UploadFile(
		projectID,
		file.Name(),
		nil,
	)

	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"file": file.Name(),
		"url":  projectFile.URL,
	}).Debug("uploaded file")

	gitlabBaseURL := ctx.Config.GitLabURLs.Download
	// projectFile.URL from upload: /uploads/<hash>/filename.txt
	linkURL := gitlabBaseURL + "/" + projectID + projectFile.URL
	name := artifact.Name
	releaseLink, _, err := c.client.ReleaseLinks.CreateReleaseLink(
		projectID,
		releaseID,
		&gitlab.CreateReleaseLinkOptions{
			Name: &name,
			URL:  &linkURL,
		})

	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"id":  releaseLink.ID,
		"url": releaseLink.URL,
	}).Debug("created release link")

	fileUploadHash, err := extractProjectFileHashFrom(projectFile.URL)
	if err != nil {
		return err
	}

	// for checksums.txt the field is nil, so we initialize it
	if artifact.Extra == nil {
		artifact.Extra = make(map[string]interface{})
	}
	// we set this hash to be able to download the file
	// in following publish pipes like brew, scoop
	artifact.Extra["ArtifactUploadHash"] = fileUploadHash

	return err
}

// extractProjectFileHashFrom extracts the hash from the
// relative project file url of the format '/uploads/<hash>/filename.ext'
func extractProjectFileHashFrom(projectFileURL string) (string, error) {
	log.WithField("projectFileURL", projectFileURL).Debug("extract file hash from")
	splittedProjectFileURL := strings.Split(projectFileURL, "/")
	if len(splittedProjectFileURL) != 4 {
		log.WithField("projectFileURL", projectFileURL).Debug("could not extract file hash")
		return "", ErrExtractHashFromFileUploadURL
	}

	fileHash := splittedProjectFileURL[2]
	log.WithFields(log.Fields{
		"projectFileURL": projectFileURL,
		"fileHash":       fileHash,
	}).Debug("extracted file hash")
	return fileHash, nil
}
