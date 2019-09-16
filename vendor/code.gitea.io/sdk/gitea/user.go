// Copyright 2014 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gitea

import (
	"fmt"

	"code.gitea.io/gitea/modules/structs"
)

// User is equal to structs.User
type User = structs.User

// GetUserInfo get user info by user's name
func (c *Client) GetUserInfo(user string) (*User, error) {
	u := new(User)
	err := c.getParsedResponse("GET", fmt.Sprintf("/users/%s", user), nil, nil, u)
	return u, err
}

// GetMyUserInfo get user info of current user
func (c *Client) GetMyUserInfo() (*User, error) {
	u := new(User)
	err := c.getParsedResponse("GET", "/user", nil, nil, u)
	return u, err
}
