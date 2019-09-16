// Copyright 2016 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gitea

import (
	"bytes"
	"encoding/json"
	"fmt"

	"code.gitea.io/gitea/modules/structs"
)

// Label is equal to structs.Label
type Label = structs.Label

// ListRepoLabels list labels of one repository
func (c *Client) ListRepoLabels(owner, repo string) ([]*Label, error) {
	labels := make([]*Label, 0, 10)
	return labels, c.getParsedResponse("GET", fmt.Sprintf("/repos/%s/%s/labels", owner, repo), nil, nil, &labels)
}

// GetRepoLabel get one label of repository by repo it
// TODO: maybe we need get a label by name
func (c *Client) GetRepoLabel(owner, repo string, id int64) (*Label, error) {
	label := new(Label)
	return label, c.getParsedResponse("GET", fmt.Sprintf("/repos/%s/%s/labels/%d", owner, repo, id), nil, nil, label)
}

// CreateLabel create one label of repository
func (c *Client) CreateLabel(owner, repo string, opt structs.CreateLabelOption) (*Label, error) {
	body, err := json.Marshal(&opt)
	if err != nil {
		return nil, err
	}
	label := new(Label)
	return label, c.getParsedResponse("POST", fmt.Sprintf("/repos/%s/%s/labels", owner, repo),
		jsonHeader, bytes.NewReader(body), label)
}

// EditLabel modify one label with options
func (c *Client) EditLabel(owner, repo string, id int64, opt structs.EditLabelOption) (*Label, error) {
	body, err := json.Marshal(&opt)
	if err != nil {
		return nil, err
	}
	label := new(Label)
	return label, c.getParsedResponse("PATCH", fmt.Sprintf("/repos/%s/%s/labels/%d", owner, repo, id), jsonHeader, bytes.NewReader(body), label)
}

// DeleteLabel delete one label of repository by id
// TODO: maybe we need delete by name
func (c *Client) DeleteLabel(owner, repo string, id int64) error {
	_, err := c.getResponse("DELETE", fmt.Sprintf("/repos/%s/%s/labels/%d", owner, repo, id), nil, nil)
	return err
}

// GetIssueLabels get labels of one issue via issue id
func (c *Client) GetIssueLabels(owner, repo string, index int64) ([]*Label, error) {
	labels := make([]*Label, 0, 5)
	return labels, c.getParsedResponse("GET", fmt.Sprintf("/repos/%s/%s/issues/%d/labels", owner, repo, index), nil, nil, &labels)
}

// AddIssueLabels add one or more labels to one issue
func (c *Client) AddIssueLabels(owner, repo string, index int64, opt structs.IssueLabelsOption) ([]*Label, error) {
	body, err := json.Marshal(&opt)
	if err != nil {
		return nil, err
	}
	var labels []*Label
	return labels, c.getParsedResponse("POST", fmt.Sprintf("/repos/%s/%s/issues/%d/labels", owner, repo, index), jsonHeader, bytes.NewReader(body), &labels)
}

// ReplaceIssueLabels replace old labels of issue with new labels
func (c *Client) ReplaceIssueLabels(owner, repo string, index int64, opt structs.IssueLabelsOption) ([]*Label, error) {
	body, err := json.Marshal(&opt)
	if err != nil {
		return nil, err
	}
	var labels []*Label
	return labels, c.getParsedResponse("PUT", fmt.Sprintf("/repos/%s/%s/issues/%d/labels", owner, repo, index), jsonHeader, bytes.NewReader(body), &labels)
}

// DeleteIssueLabel delete one label of one issue by issue id and label id
// TODO: maybe we need delete by label name and issue id
func (c *Client) DeleteIssueLabel(owner, repo string, index, label int64) error {
	_, err := c.getResponse("DELETE", fmt.Sprintf("/repos/%s/%s/issues/%d/labels/%d", owner, repo, index, label), nil, nil)
	return err
}

// ClearIssueLabels delete all the labels of one issue.
func (c *Client) ClearIssueLabels(owner, repo string, index int64) error {
	_, err := c.getResponse("DELETE", fmt.Sprintf("/repos/%s/%s/issues/%d/labels", owner, repo, index), nil, nil)
	return err
}
