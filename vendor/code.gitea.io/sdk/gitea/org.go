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

// Organization is equal to structs.Organization
type Organization = structs.Organization

// ListMyOrgs list all of current user's organizations
func (c *Client) ListMyOrgs() ([]*Organization, error) {
	orgs := make([]*Organization, 0, 5)
	return orgs, c.getParsedResponse("GET", "/user/orgs", nil, nil, &orgs)
}

// ListUserOrgs list all of some user's organizations
func (c *Client) ListUserOrgs(user string) ([]*Organization, error) {
	orgs := make([]*Organization, 0, 5)
	return orgs, c.getParsedResponse("GET", fmt.Sprintf("/users/%s/orgs", user), nil, nil, &orgs)
}

// GetOrg get one organization by name
func (c *Client) GetOrg(orgname string) (*Organization, error) {
	org := new(Organization)
	return org, c.getParsedResponse("GET", fmt.Sprintf("/orgs/%s", orgname), nil, nil, org)
}

// CreateOrg creates an organization
func (c *Client) CreateOrg(opt structs.CreateOrgOption) (*Organization, error) {
	body, err := json.Marshal(&opt)
	if err != nil {
		return nil, err
	}
	org := new(Organization)
	return org, c.getParsedResponse("POST", "/orgs", jsonHeader, bytes.NewReader(body), org)
}

// EditOrg modify one organization via options
func (c *Client) EditOrg(orgname string, opt structs.EditOrgOption) error {
	body, err := json.Marshal(&opt)
	if err != nil {
		return err
	}
	_, err = c.getResponse("PATCH", fmt.Sprintf("/orgs/%s", orgname), jsonHeader, bytes.NewReader(body))
	return err
}

// DeleteOrg deletes an organization
func (c *Client) DeleteOrg(orgname string) error {
	_, err := c.getResponse("DELETE", fmt.Sprintf("/orgs/%s", orgname), nil, nil)
	return err
}
