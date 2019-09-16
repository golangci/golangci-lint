// Copyright 2014 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gitea

import (
	"bytes"
	"encoding/json"
	"fmt"

	"code.gitea.io/gitea/modules/structs"
)

// Repository is equal to structs.Repository
type Repository = structs.Repository

// CreateRepoOption is equal to structs.CreateRepoOption
type CreateRepoOption = structs.CreateRepoOption

// ListMyRepos lists all repositories for the authenticated user that has access to.
func (c *Client) ListMyRepos() ([]*Repository, error) {
	repos := make([]*Repository, 0, 10)
	return repos, c.getParsedResponse("GET", "/user/repos", nil, nil, &repos)
}

// ListUserRepos list all repositories of one user by user's name
func (c *Client) ListUserRepos(user string) ([]*Repository, error) {
	repos := make([]*Repository, 0, 10)
	return repos, c.getParsedResponse("GET", fmt.Sprintf("/users/%s/repos", user), nil, nil, &repos)
}

// ListOrgRepos list all repositories of one organization by organization's name
func (c *Client) ListOrgRepos(org string) ([]*Repository, error) {
	repos := make([]*Repository, 0, 10)
	return repos, c.getParsedResponse("GET", fmt.Sprintf("/orgs/%s/repos", org), nil, nil, &repos)
}

// CreateRepo creates a repository for authenticated user.
func (c *Client) CreateRepo(opt CreateRepoOption) (*Repository, error) {
	body, err := json.Marshal(&opt)
	if err != nil {
		return nil, err
	}
	repo := new(Repository)
	return repo, c.getParsedResponse("POST", "/user/repos", jsonHeader, bytes.NewReader(body), repo)
}

// CreateOrgRepo creates an organization repository for authenticated user.
func (c *Client) CreateOrgRepo(org string, opt CreateRepoOption) (*Repository, error) {
	body, err := json.Marshal(&opt)
	if err != nil {
		return nil, err
	}
	repo := new(Repository)
	return repo, c.getParsedResponse("POST", fmt.Sprintf("/org/%s/repos", org), jsonHeader, bytes.NewReader(body), repo)
}

// GetRepo returns information of a repository of given owner.
func (c *Client) GetRepo(owner, reponame string) (*Repository, error) {
	repo := new(Repository)
	return repo, c.getParsedResponse("GET", fmt.Sprintf("/repos/%s/%s", owner, reponame), nil, nil, repo)
}

// EditRepo edit the properties of a repository
func (c *Client) EditRepo(owner, reponame string, opt structs.EditRepoOption) (*Repository, error) {
	body, err := json.Marshal(&opt)
	if err != nil {
		return nil, err
	}
	repo := new(Repository)
	return repo, c.getParsedResponse("PATCH", fmt.Sprintf("/repos/%s/%s", owner, reponame), jsonHeader, bytes.NewReader(body), repo)
}

// DeleteRepo deletes a repository of user or organization.
func (c *Client) DeleteRepo(owner, repo string) error {
	_, err := c.getResponse("DELETE", fmt.Sprintf("/repos/%s/%s", owner, repo), nil, nil)
	return err
}

// MigrateRepo migrates a repository from other Git hosting sources for the
// authenticated user.
//
// To migrate a repository for a organization, the authenticated user must be a
// owner of the specified organization.
func (c *Client) MigrateRepo(opt structs.MigrateRepoOption) (*Repository, error) {
	body, err := json.Marshal(&opt)
	if err != nil {
		return nil, err
	}
	repo := new(Repository)
	return repo, c.getParsedResponse("POST", "/repos/migrate", jsonHeader, bytes.NewReader(body), repo)
}

// MirrorSync adds a mirrored repository to the mirror sync queue.
func (c *Client) MirrorSync(owner, repo string) error {
	_, err := c.getResponse("POST", fmt.Sprintf("/repos/%s/%s/mirror-sync", owner, repo), nil, nil)
	return err
}
