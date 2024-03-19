package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const endpoint = "https://api.github.com/repos/golangci/golangci-lint/releases/latest"

type releaseInfo struct {
	TagName string `json:"tag_name"`
}

// GetLatestVersion gets latest release information.
func GetLatestVersion() (string, error) {
	//nolint:noctx // request timeout handled by the client
	req, err := http.NewRequest(http.MethodGet, endpoint, http.NoBody)
	if err != nil {
		return "", fmt.Errorf("prepare a HTTP request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{Timeout: 2 * time.Second}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("get HTTP response for the latest tag: %w", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read a body for the latest tag: %w", err)
	}

	release := releaseInfo{}

	err = json.Unmarshal(body, &release)
	if err != nil {
		return "", fmt.Errorf("unmarshal the body for the latest tag: %w", err)
	}

	return release.TagName, nil
}
