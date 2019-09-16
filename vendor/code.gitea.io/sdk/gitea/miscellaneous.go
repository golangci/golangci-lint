// Copyright 2015 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gitea

import "code.gitea.io/gitea/modules/structs"

// ServerVersion returns the version of the server
func (c *Client) ServerVersion() (string, error) {
	v := structs.ServerVersion{}
	return v.Version, c.getParsedResponse("GET", "/api/v1/version", nil, nil, &v)
}
