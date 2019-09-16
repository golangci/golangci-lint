// Copyright 2017 Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gitea

import (
	"bytes"
	"encoding/json"
	"fmt"

	"code.gitea.io/gitea/modules/structs"
)

// GPGKey is equal to structs.GPGKey
type GPGKey = structs.GPGKey

// ListGPGKeys list all the GPG keys of the user
func (c *Client) ListGPGKeys(user string) ([]*GPGKey, error) {
	keys := make([]*GPGKey, 0, 10)
	return keys, c.getParsedResponse("GET", fmt.Sprintf("/users/%s/gpg_keys", user), nil, nil, &keys)
}

// ListMyGPGKeys list all the GPG keys of current user
func (c *Client) ListMyGPGKeys() ([]*GPGKey, error) {
	keys := make([]*GPGKey, 0, 10)
	return keys, c.getParsedResponse("GET", "/user/gpg_keys", nil, nil, &keys)
}

// GetGPGKey get current user's GPG key by key id
func (c *Client) GetGPGKey(keyID int64) (*GPGKey, error) {
	key := new(GPGKey)
	return key, c.getParsedResponse("GET", fmt.Sprintf("/user/gpg_keys/%d", keyID), nil, nil, &key)
}

// CreateGPGKey create GPG key with options
func (c *Client) CreateGPGKey(opt structs.CreateGPGKeyOption) (*GPGKey, error) {
	body, err := json.Marshal(&opt)
	if err != nil {
		return nil, err
	}
	key := new(GPGKey)
	return key, c.getParsedResponse("POST", "/user/gpg_keys", jsonHeader, bytes.NewReader(body), key)
}

// DeleteGPGKey delete GPG key with key id
func (c *Client) DeleteGPGKey(keyID int64) error {
	_, err := c.getResponse("DELETE", fmt.Sprintf("/user/gpg_keys/%d", keyID), nil, nil)
	return err
}
