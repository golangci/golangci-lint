// Copyright 2019 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gitea

import (
	"bytes"
	"encoding/json"
	"fmt"

	"code.gitea.io/gitea/modules/structs"
)

// GitHook is equal to structs.GitHook
type GitHook = structs.GitHook

// GitHookList is equal to structs.GitHookList
type GitHookList = structs.GitHookList

// ListRepoGitHooks list all the Git hooks of one repository
func (c *Client) ListRepoGitHooks(user, repo string) (GitHookList, error) {
	hooks := make([]*GitHook, 0, 10)
	return hooks, c.getParsedResponse("GET", fmt.Sprintf("/repos/%s/%s/hooks/git", user, repo), nil, nil, &hooks)
}

// GetRepoGitHook get a Git hook of a repository
func (c *Client) GetRepoGitHook(user, repo, id string) (*GitHook, error) {
	h := new(GitHook)
	return h, c.getParsedResponse("GET", fmt.Sprintf("/repos/%s/%s/hooks/git/%s", user, repo, id), nil, nil, h)
}

// EditRepoGitHook modify one Git hook of a repository
func (c *Client) EditRepoGitHook(user, repo, id string, opt structs.EditGitHookOption) error {
	body, err := json.Marshal(&opt)
	if err != nil {
		return err
	}
	_, err = c.getResponse("PATCH", fmt.Sprintf("/repos/%s/%s/hooks/git/%s", user, repo, id), jsonHeader, bytes.NewReader(body))
	return err
}

// DeleteRepoGitHook delete one Git hook from a repository
func (c *Client) DeleteRepoGitHook(user, repo, id string) error {
	_, err := c.getResponse("DELETE", fmt.Sprintf("/repos/%s/%s/hooks/git/%s", user, repo, id), nil, nil)
	return err
}
