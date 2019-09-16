// Copyright 2017 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gitea

import (
	"bytes"
	"encoding/json"
	"fmt"

	"code.gitea.io/gitea/modules/structs"
)

// TrackedTime is equal to structs.TrackedTime
type TrackedTime = structs.TrackedTime

// TrackedTimes is equal to structs.TrackedTimes
type TrackedTimes = structs.TrackedTimes

// GetUserTrackedTimes list tracked times of a user
func (c *Client) GetUserTrackedTimes(owner, repo, user string) (TrackedTimes, error) {
	times := make(TrackedTimes, 0, 10)
	return times, c.getParsedResponse("GET", fmt.Sprintf("/repos/%s/%s/times/%s", owner, repo, user), nil, nil, &times)
}

// GetRepoTrackedTimes list tracked times of a repository
func (c *Client) GetRepoTrackedTimes(owner, repo string) (TrackedTimes, error) {
	times := make(TrackedTimes, 0, 10)
	return times, c.getParsedResponse("GET", fmt.Sprintf("/repos/%s/%s/times", owner, repo), nil, nil, &times)
}

// GetMyTrackedTimes list tracked times of the current user
func (c *Client) GetMyTrackedTimes() (TrackedTimes, error) {
	times := make(TrackedTimes, 0, 10)
	return times, c.getParsedResponse("GET", "/user/times", nil, nil, &times)
}

// AddTime adds time to issue with the given index
func (c *Client) AddTime(owner, repo string, index int64, opt structs.AddTimeOption) (*TrackedTime, error) {
	body, err := json.Marshal(&opt)
	if err != nil {
		return nil, err
	}
	t := new(TrackedTime)
	return t, c.getParsedResponse("POST", fmt.Sprintf("/repos/%s/%s/issues/%d/times", owner, repo, index),
		jsonHeader, bytes.NewReader(body), t)
}

// ListTrackedTimes get tracked times of one issue via issue id
func (c *Client) ListTrackedTimes(owner, repo string, index int64) (TrackedTimes, error) {
	times := make(TrackedTimes, 0, 5)
	return times, c.getParsedResponse("GET", fmt.Sprintf("/repos/%s/%s/issues/%d/times", owner, repo, index), nil, nil, &times)
}
