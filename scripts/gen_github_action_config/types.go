package main

import "fmt"

type versionConfig struct {
	Error string `json:",omitempty"`

	TargetVersion string `json:",omitempty"`
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
