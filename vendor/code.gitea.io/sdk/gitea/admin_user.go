// Copyright 2015 The Gogs Authors. All rights reserved.
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

// AdminListUsers lists all users
func (c *Client) AdminListUsers() ([]*User, error) {
	users := make([]*User, 0, 10)
	return users, c.getParsedResponse("GET", "/admin/users", nil, nil, &users)
}

// AdminCreateUser create a user
func (c *Client) AdminCreateUser(opt structs.CreateUserOption) (*User, error) {
	body, err := json.Marshal(&opt)
	if err != nil {
		return nil, err
	}
	user := new(User)
	return user, c.getParsedResponse("POST", "/admin/users", jsonHeader, bytes.NewReader(body), user)
}

// AdminEditUser modify user informations
func (c *Client) AdminEditUser(user string, opt structs.EditUserOption) error {
	body, err := json.Marshal(&opt)
	if err != nil {
		return err
	}
	_, err = c.getResponse("PATCH", fmt.Sprintf("/admin/users/%s", user), jsonHeader, bytes.NewReader(body))
	return err
}

// AdminDeleteUser delete one user according name
func (c *Client) AdminDeleteUser(user string) error {
	_, err := c.getResponse("DELETE", fmt.Sprintf("/admin/users/%s", user), nil, nil)
	return err
}

// AdminCreateUserPublicKey adds a public key for the user
func (c *Client) AdminCreateUserPublicKey(user string, opt structs.CreateKeyOption) (*PublicKey, error) {
	body, err := json.Marshal(&opt)
	if err != nil {
		return nil, err
	}
	key := new(PublicKey)
	return key, c.getParsedResponse("POST", fmt.Sprintf("/admin/users/%s/keys", user), jsonHeader, bytes.NewReader(body), key)
}

// AdminDeleteUserPublicKey deletes a user's public key
func (c *Client) AdminDeleteUserPublicKey(user string, keyID int) error {
	_, err := c.getResponse("DELETE", fmt.Sprintf("/admin/users/%s/keys/%d", user, keyID), nil, nil)
	return err
}
