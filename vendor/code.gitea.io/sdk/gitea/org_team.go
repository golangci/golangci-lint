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

// Team is equal to structs.Team
type Team = structs.Team

// ListOrgTeams lists all teams of an organization
func (c *Client) ListOrgTeams(org string) ([]*Team, error) {
	teams := make([]*Team, 0, 10)
	return teams, c.getParsedResponse("GET", fmt.Sprintf("/orgs/%s/teams", org), nil, nil, &teams)
}

// ListMyTeams lists all the teams of the current user
func (c *Client) ListMyTeams() ([]*Team, error) {
	teams := make([]*Team, 0, 10)
	return teams, c.getParsedResponse("GET", "/user/teams", nil, nil, &teams)
}

// GetTeam gets a team by ID
func (c *Client) GetTeam(id int64) (*Team, error) {
	t := new(Team)
	return t, c.getParsedResponse("GET", fmt.Sprintf("/teams/%d", id), nil, nil, t)
}

// CreateTeam creates a team for an organization
func (c *Client) CreateTeam(org string, opt structs.CreateTeamOption) (*Team, error) {
	body, err := json.Marshal(&opt)
	if err != nil {
		return nil, err
	}
	t := new(Team)
	return t, c.getParsedResponse("POST", fmt.Sprintf("/orgs/%s/teams", org), jsonHeader, bytes.NewReader(body), t)
}

// EditTeam edits a team of an organization
func (c *Client) EditTeam(id int64, opt structs.EditTeamOption) error {
	body, err := json.Marshal(&opt)
	if err != nil {
		return err
	}
	_, err = c.getResponse("PATCH", fmt.Sprintf("/teams/%d", id), jsonHeader, bytes.NewReader(body))
	return err
}

// DeleteTeam deletes a team of an organization
func (c *Client) DeleteTeam(id int64) error {
	_, err := c.getResponse("DELETE", fmt.Sprintf("/teams/%d", id), nil, nil)
	return err
}

// ListTeamMembers lists all members of a team
func (c *Client) ListTeamMembers(id int64) ([]*User, error) {
	members := make([]*User, 0, 10)
	return members, c.getParsedResponse("GET", fmt.Sprintf("/teams/%d/members", id), nil, nil, &members)
}

// GetTeamMember gets a member of a team
func (c *Client) GetTeamMember(id int64, user string) (*User, error) {
	m := new(User)
	return m, c.getParsedResponse("GET", fmt.Sprintf("/teams/%d/members/%s", id, user), nil, nil, m)
}

// AddTeamMember adds a member to a team
func (c *Client) AddTeamMember(id int64, user string) error {
	_, err := c.getResponse("PUT", fmt.Sprintf("/teams/%d/members/%s", id, user), nil, nil)
	return err
}

// RemoveTeamMember removes a member from a team
func (c *Client) RemoveTeamMember(id int64, user string) error {
	_, err := c.getResponse("DELETE", fmt.Sprintf("/teams/%d/members/%s", id, user), nil, nil)
	return err
}

// ListTeamRepositories lists all repositories of a team
func (c *Client) ListTeamRepositories(id int64) ([]*Repository, error) {
	repos := make([]*Repository, 0, 10)
	return repos, c.getParsedResponse("GET", fmt.Sprintf("/teams/%d/repos", id), nil, nil, &repos)
}

// AddTeamRepository adds a repository to a team
func (c *Client) AddTeamRepository(id int64, org, repo string) error {
	_, err := c.getResponse("PUT", fmt.Sprintf("/teams/%d/repos/%s/%s", id, org, repo), nil, nil)
	return err
}

// RemoveTeamRepository removes a repository from a team
func (c *Client) RemoveTeamRepository(id int64, org, repo string) error {
	_, err := c.getResponse("DELETE", fmt.Sprintf("/teams/%d/repos/%s/%s", id, org, repo), nil, nil)
	return err
}
