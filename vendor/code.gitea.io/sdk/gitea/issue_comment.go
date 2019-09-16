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

// Comment is equal to structs.Comment
type Comment = structs.Comment

// ListIssueComments list comments on an issue.
func (c *Client) ListIssueComments(owner, repo string, index int64) ([]*Comment, error) {
	comments := make([]*Comment, 0, 10)
	return comments, c.getParsedResponse("GET", fmt.Sprintf("/repos/%s/%s/issues/%d/comments", owner, repo, index), nil, nil, &comments)
}

// ListRepoIssueComments list comments for a given repo.
func (c *Client) ListRepoIssueComments(owner, repo string) ([]*Comment, error) {
	comments := make([]*Comment, 0, 10)
	return comments, c.getParsedResponse("GET", fmt.Sprintf("/repos/%s/%s/issues/comments", owner, repo), nil, nil, &comments)
}

// CreateIssueComment create comment on an issue.
func (c *Client) CreateIssueComment(owner, repo string, index int64, opt structs.CreateIssueCommentOption) (*Comment, error) {
	body, err := json.Marshal(&opt)
	if err != nil {
		return nil, err
	}
	comment := new(Comment)
	return comment, c.getParsedResponse("POST", fmt.Sprintf("/repos/%s/%s/issues/%d/comments", owner, repo, index), jsonHeader, bytes.NewReader(body), comment)
}

// EditIssueComment edits an issue comment.
func (c *Client) EditIssueComment(owner, repo string, commentID int64, opt structs.EditIssueCommentOption) (*Comment, error) {
	body, err := json.Marshal(&opt)
	if err != nil {
		return nil, err
	}
	comment := new(Comment)
	return comment, c.getParsedResponse("PATCH", fmt.Sprintf("/repos/%s/%s/issues/comments/%d", owner, repo, commentID), jsonHeader, bytes.NewReader(body), comment)
}

// DeleteIssueComment deletes an issue comment.
func (c *Client) DeleteIssueComment(owner, repo string, commentID int64) error {
	_, err := c.getResponse("DELETE", fmt.Sprintf("/repos/%s/%s/issues/comments/%d", owner, repo, commentID), nil, nil)
	return err
}
