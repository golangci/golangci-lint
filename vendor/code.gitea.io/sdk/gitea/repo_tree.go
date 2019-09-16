// Copyright 2018 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gitea

import (
	"fmt"

	"code.gitea.io/gitea/modules/structs"
)

// GetTrees downloads a file of repository, ref can be branch/tag/commit.
// e.g.: ref -> master, tree -> macaron.go(no leading slash)
func (c *Client) GetTrees(user, repo, ref string, recursive bool) (*structs.GitTreeResponse, error) {
	var trees structs.GitTreeResponse
	var path = fmt.Sprintf("/repos/%s/%s/git/trees/%s", user, repo, ref)
	if recursive {
		path += "?recursive=1"
	}
	err := c.getParsedResponse("GET", path, nil, nil, &trees)
	return &trees, err
}
