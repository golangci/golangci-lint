// Copyright 2017 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gitea

import (
	"bytes"
	"encoding/json"
	"fmt"

	"code.gitea.io/gitea/modules/structs"
)

// Status is equal to structs.Status
type Status = structs.Status

// CreateStatus creates a new Status for a given Commit
//
// POST /repos/:owner/:repo/statuses/:sha
func (c *Client) CreateStatus(owner, repo, sha string, opts structs.CreateStatusOption) (*Status, error) {
	body, err := json.Marshal(&opts)
	if err != nil {
		return nil, err
	}
	status := &Status{}
	return status, c.getParsedResponse("POST", fmt.Sprintf("/repos/%s/%s/statuses/%s", owner, repo, sha),
		jsonHeader, bytes.NewReader(body), status)
}

// ListStatuses returns all statuses for a given Commit
//
// GET /repos/:owner/:repo/commits/:ref/statuses
func (c *Client) ListStatuses(owner, repo, sha string, opts structs.ListStatusesOption) ([]*Status, error) {
	statuses := make([]*structs.Status, 0, 10)
	return statuses, c.getParsedResponse("GET", fmt.Sprintf("/repos/%s/%s/commits/%s/statuses?page=%d", owner, repo, sha, opts.Page), nil, nil, &statuses)
}

// GetCombinedStatus returns the CombinedStatus for a given Commit
//
// GET /repos/:owner/:repo/commits/:ref/status
func (c *Client) GetCombinedStatus(owner, repo, sha string) (*structs.CombinedStatus, error) {
	status := &structs.CombinedStatus{}
	return status, c.getParsedResponse("GET", fmt.Sprintf("/repos/%s/%s/commits/%s/status", owner, repo, sha), nil, nil, status)
}
