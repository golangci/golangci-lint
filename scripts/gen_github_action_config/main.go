package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

const noPatch = -1

type logInfo struct {
	Warning string `json:",omitempty"`
	Info    string `json:",omitempty"`
}

type versionConfig struct {
	Error string `json:",omitempty"`

	Log *logInfo `json:",omitempty"`

	TargetVersion string `json:",omitempty"`
	AssetURL      string `json:",omitempty"`
}

type actionConfig struct {
	MinorVersionToConfig map[string]versionConfig
}

type version struct {
	major, minor, patch int
}

func (v version) String() string {
	ret := fmt.Sprintf("v%d.%d", v.major, v.minor)

	if v.patch != noPatch {
		ret += fmt.Sprintf(".%d", v.patch)
	}

	return ret
}

func (v version) isAfterOrEq(vv *version) bool {
	if v.major != vv.major {
		return v.major >= vv.major
	}

	if v.minor != vv.minor {
		return v.minor >= vv.minor
	}

	return v.patch >= vv.patch
}

type release struct {
	TagName       string
	ReleaseAssets struct {
		Nodes []releaseAsset
	} `graphql:"releaseAssets(first: 50)"`
}

type releaseAsset struct {
	DownloadURL string
}

func main() {
	if err := generate(context.Background()); err != nil {
		log.Fatal(err)
	}
}

func generate(ctx context.Context) error {
	if len(os.Args) != 2 {
		return fmt.Errorf("usage: go run .../main.go out-path.json")
	}

	allReleases, err := fetchAllReleases(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch all releases: %w", err)
	}

	cfg, err := buildConfig(allReleases)
	if err != nil {
		return fmt.Errorf("failed to build config: %w", err)
	}

	outFile, err := os.Create(os.Args[1])
	if err != nil {
		return fmt.Errorf("failed to create output config file: %w", err)
	}

	defer outFile.Close()

	enc := json.NewEncoder(outFile)
	enc.SetIndent("", "  ")

	if err = enc.Encode(cfg); err != nil {
		return fmt.Errorf("failed to json encode config: %w", err)
	}

	return nil
}

func fetchAllReleases(ctx context.Context) ([]release, error) {
	githubToken := os.Getenv("GITHUB_TOKEN")
	if githubToken == "" {
		return nil, errors.New("no GITHUB_TOKEN environment variable")
	}

	client := githubv4.NewClient(oauth2.NewClient(ctx, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: githubToken})))

	var q struct {
		Repository struct {
			Releases struct {
				Nodes    []release
				PageInfo struct {
					EndCursor   githubv4.String
					HasNextPage bool
				}
			} `graphql:"releases(first: 100, orderBy: { field: CREATED_AT, direction: DESC }, after: $releasesCursor)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	vars := map[string]any{
		"owner":          githubv4.String("golangci"),
		"name":           githubv4.String("golangci-lint"),
		"releasesCursor": (*githubv4.String)(nil),
	}

	var allReleases []release
	for {
		err := client.Query(ctx, &q, vars)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch releases page from GitHub: %w", err)
		}

		releases := q.Repository.Releases
		allReleases = append(allReleases, releases.Nodes...)

		if !releases.PageInfo.HasNextPage {
			break
		}

		vars["releasesCursor"] = githubv4.NewString(releases.PageInfo.EndCursor)
	}

	return allReleases, nil
}

func buildConfig(releases []release) (*actionConfig, error) {
	versionToRelease := map[version]release{}

	for _, rel := range releases {
		ver, err := parseVersion(rel.TagName)
		if err != nil {
			return nil, fmt.Errorf("failed to parse release %s version: %w", rel.TagName, err)
		}

		if _, ok := versionToRelease[*ver]; ok {
			return nil, fmt.Errorf("duplicate release %s", rel.TagName)
		}

		versionToRelease[*ver] = rel
	}

	maxPatchReleases := map[string]version{}

	for ver := range versionToRelease {
		key := fmt.Sprintf("v%d.%d", ver.major, ver.minor)

		if mapVer, ok := maxPatchReleases[key]; !ok || ver.isAfterOrEq(&mapVer) {
			maxPatchReleases[key] = ver
		}
	}

	minorVersionToConfig := map[string]versionConfig{}
	minAllowedVersion := version{major: 1, minor: 14, patch: 0}

	latestVersion := version{}
	latestVersionConfig := versionConfig{}

	for minorVersionedStr, maxPatchVersion := range maxPatchReleases {
		if !maxPatchVersion.isAfterOrEq(&minAllowedVersion) {
			minorVersionToConfig[minorVersionedStr] = versionConfig{
				Error: fmt.Sprintf("golangci-lint version '%s' isn't supported: we support only %s and later versions",
					minorVersionedStr, minAllowedVersion),
			}
			continue
		}

		maxPatchVersion := maxPatchVersion

		assetURL, err := findLinuxAssetURL(&maxPatchVersion, versionToRelease[maxPatchVersion].ReleaseAssets.Nodes)
		if err != nil {
			return nil, fmt.Errorf("failed to find linux asset url for release %s: %w", maxPatchVersion, err)
		}

		minorVersionToConfig[minorVersionedStr] = versionConfig{
			TargetVersion: maxPatchVersion.String(),
			AssetURL:      assetURL,
		}

		if maxPatchVersion.isAfterOrEq(&latestVersion) {
			latestVersion = maxPatchVersion
			latestVersionConfig.TargetVersion = maxPatchVersion.String()
			latestVersionConfig.AssetURL = assetURL
		}
	}

	minorVersionToConfig["latest"] = latestVersionConfig

	return &actionConfig{MinorVersionToConfig: minorVersionToConfig}, nil
}

func findLinuxAssetURL(ver *version, releaseAssets []releaseAsset) (string, error) {
	pattern := fmt.Sprintf("golangci-lint-%d.%d.%d-linux-amd64.tar.gz", ver.major, ver.minor, ver.patch)

	for _, relAsset := range releaseAssets {
		if strings.HasSuffix(relAsset.DownloadURL, pattern) {
			return relAsset.DownloadURL, nil
		}
	}

	return "", fmt.Errorf("no matched asset url for pattern %q", pattern)
}

func parseVersion(s string) (*version, error) {
	const vPrefix = "v"
	if !strings.HasPrefix(s, vPrefix) {
		return nil, fmt.Errorf("version %q should start with %q", s, vPrefix)
	}

	parts := strings.Split(strings.TrimPrefix(s, vPrefix), ".")

	var nums []int
	for _, part := range parts {
		num, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("failed to parse version %q: %w", s, err)
		}

		nums = append(nums, num)
	}

	switch len(nums) {
	case 2:
		return &version{major: nums[0], minor: nums[1], patch: noPatch}, nil
	case 3:
		return &version{major: nums[0], minor: nums[1], patch: nums[2]}, nil
	default:
		return nil, fmt.Errorf("invalid version format: %s", s)
	}
}
