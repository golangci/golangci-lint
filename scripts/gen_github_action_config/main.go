package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

const noPatch = -1

func main() {
	if err := run(context.Background()); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	if len(os.Args) != 2 {
		return fmt.Errorf("usage: go run .../main.go out-path.json")
	}

	allReleases, err := fetchAllReleases(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch all releases: %w", err)
	}

	dest := os.Args[1]

	ext := filepath.Ext(dest)

	// https://github.com/golangci/golangci-lint-action/blob/5421a116d2bf2a1d53595d0dca7da6e18bd1cfd7/src/version.ts#L43-L47
	minAllowedVersionV1 := version{major: 1, minor: 28, patch: 3}

	// For compatibility with v1: it should always be related to v1 only.
	// TODO(ldez): it should be removed but I don't know when.
	err = generate(allReleases, minAllowedVersionV1, dest)
	if err != nil {
		return fmt.Errorf("failed to generate v1: %w", err)
	}

	destV1 := filepath.Join(filepath.Dir(dest), strings.TrimSuffix(filepath.Base(dest), ext)+"-v1"+ext)

	err = generate(allReleases, minAllowedVersionV1, destV1)
	if err != nil {
		return fmt.Errorf("failed to generate v1: %w", err)
	}

	return nil
}

func generate(allReleases []release, minAllowedVersion version, dest string) error {
	cfg, err := buildConfig(allReleases, minAllowedVersion)
	if err != nil {
		return fmt.Errorf("failed to build config: %w", err)
	}

	outFile, err := os.Create(dest)
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

func buildConfig(releases []release, minAllowedVersion version) (*actionConfig, error) {
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

	latestVersion := version{}
	latestVersionConfig := versionConfig{}

	for minorVersionedStr, maxPatchVersion := range maxPatchReleases {
		if minAllowedVersion.major < maxPatchVersion.major {
			minorVersionToConfig[minorVersionedStr] = versionConfig{
				Error: fmt.Sprintf("golangci-lint version '%s' isn't supported: only v%d versions are supported",
					minorVersionedStr, minAllowedVersion.major),
			}
			continue
		}

		if !maxPatchVersion.isAfterOrEq(&minAllowedVersion) {
			minorVersionToConfig[minorVersionedStr] = versionConfig{
				Error: fmt.Sprintf("golangci-lint version '%s' isn't supported: we support only %s and later versions",
					minorVersionedStr, minAllowedVersion),
			}
			continue
		}

		err := findLinuxAssetURL(&maxPatchVersion, versionToRelease[maxPatchVersion].ReleaseAssets.Nodes)
		if err != nil {
			return nil, fmt.Errorf("failed to find linux asset url for release %s: %w", maxPatchVersion, err)
		}

		minorVersionToConfig[minorVersionedStr] = versionConfig{
			TargetVersion: maxPatchVersion.String(),
		}

		if maxPatchVersion.isAfterOrEq(&latestVersion) {
			latestVersion = maxPatchVersion
			latestVersionConfig.TargetVersion = maxPatchVersion.String()
		}
	}

	minorVersionToConfig["latest"] = latestVersionConfig

	return &actionConfig{MinorVersionToConfig: minorVersionToConfig}, nil
}

func findLinuxAssetURL(ver *version, releaseAssets []releaseAsset) error {
	pattern := fmt.Sprintf("golangci-lint-%d.%d.%d-linux-amd64.tar.gz", ver.major, ver.minor, ver.patch)

	for _, relAsset := range releaseAssets {
		if strings.HasSuffix(relAsset.DownloadURL, pattern) {
			return nil
		}
	}

	return fmt.Errorf("no matched asset url for pattern %q", pattern)
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
