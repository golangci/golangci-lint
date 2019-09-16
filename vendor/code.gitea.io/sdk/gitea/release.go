// Copyright 2016 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gitea

import (
	"bytes"
	"encoding/json"
	"fmt"

	"code.gitea.io/gitea/modules/structs"
)

// Release is equal to structs.Release
type Release = structs.Release

// ListReleases list releases of a repository
func (c *Client) ListReleases(user, repo string) ([]*Release, error) {
	releases := make([]*Release, 0, 10)
	err := c.getParsedResponse("GET",
		fmt.Sprintf("/repos/%s/%s/releases", user, repo),
		nil, nil, &releases)
	return releases, err
}

// GetRelease get a release of a repository
func (c *Client) GetRelease(user, repo string, id int64) (*Release, error) {
	r := new(Release)
	err := c.getParsedResponse("GET",
		fmt.Sprintf("/repos/%s/%s/releases/%d", user, repo, id),
		nil, nil, &r)
	return r, err
}

// CreateRelease create a release
func (c *Client) CreateRelease(user, repo string, form structs.CreateReleaseOption) (*Release, error) {
	body, err := json.Marshal(form)
	if err != nil {
		return nil, err
	}
	r := new(Release)
	err = c.getParsedResponse("POST",
		fmt.Sprintf("/repos/%s/%s/releases", user, repo),
		jsonHeader, bytes.NewReader(body), r)
	return r, err
}

// EditRelease edit a release
func (c *Client) EditRelease(user, repo string, id int64, form structs.EditReleaseOption) (*Release, error) {
	body, err := json.Marshal(form)
	if err != nil {
		return nil, err
	}
	r := new(Release)
	err = c.getParsedResponse("PATCH",
		fmt.Sprintf("/repos/%s/%s/releases/%d", user, repo, id),
		jsonHeader, bytes.NewReader(body), r)
	return r, err
}

// DeleteRelease delete a release from a repository
func (c *Client) DeleteRelease(user, repo string, id int64) error {
	_, err := c.getResponse("DELETE",
		fmt.Sprintf("/repos/%s/%s/releases/%d", user, repo, id),
		nil, nil)
	return err
}
