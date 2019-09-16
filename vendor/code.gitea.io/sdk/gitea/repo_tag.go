// Copyright 2019 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gitea

import (
	"fmt"

	"code.gitea.io/gitea/modules/structs"
)

// Tag is equal to structs.Tag
type Tag = structs.Tag

// ListRepoTags list all the branches of one repository
func (c *Client) ListRepoTags(user, repo string) ([]*Tag, error) {
	tags := make([]*Tag, 0, 10)
	return tags, c.getParsedResponse("GET", fmt.Sprintf("/repos/%s/%s/tags", user, repo), nil, nil, &tags)
}
