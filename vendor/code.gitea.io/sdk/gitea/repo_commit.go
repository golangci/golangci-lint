// Copyright 2018 The Gogs Authors. All rights reserved.
// Copyright 2019 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gitea

import (
	"fmt"

	"code.gitea.io/gitea/modules/structs"
)

// Commit is equal to structs.Commit
type Commit = structs.Commit

// GetSingleCommit returns a single commit
func (c *Client) GetSingleCommit(user, repo, commitID string) (*Commit, error) {
	commit := new(Commit)
	return commit, c.getParsedResponse("GET", fmt.Sprintf("/repos/%s/%s/git/commits/%s", user, repo, commitID), nil, nil, &commit)
}
