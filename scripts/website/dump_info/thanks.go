package main

import (
	"fmt"
	"maps"
	"regexp"
	"slices"
	"strings"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/lint/linter"
	"github.com/golangci/golangci-lint/v2/pkg/lint/lintersdb"
)

const (
	hostGitHub = "github"
	hostGitLab = "gitlab"
)

type authorDetails struct {
	Name    string   `json:"name"`
	Linters []string `json:"linters"`
	Profile string   `json:"profile"`
	Avatar  string   `json:"avatar"`
}

func saveThanksList(dst string) error {
	return saveToJSONFile(dst, getThanksList())
}

func getThanksList() []*authorDetails {
	addedAuthors := map[string]*authorDetails{}

	linters, _ := lintersdb.NewLinterBuilder().Build(config.NewDefault())

	for _, lc := range linters {
		if lc.Internal {
			continue
		}

		if lc.OriginalURL == "" {
			continue
		}

		if lc.IsDeprecated() && lc.Deprecation.Level > linter.DeprecationWarning {
			continue
		}

		info := extractInfo(lc)

		switch {
		case info.FromGitHub():
			if _, ok := addedAuthors[info.Author]; ok {
				addedAuthors[info.Author].Linters = append(addedAuthors[info.Author].Linters, lc.Name())
			} else {
				addedAuthors[info.Author] = &authorDetails{
					Name:    info.Author,
					Linters: []string{lc.Name()},
					Profile: fmt.Sprintf("https://github.com/sponsors/%s", info.Author),
					Avatar:  fmt.Sprintf("https://github.com/%[1]s.png", info.Author),
				}
			}

		case info.FromGitLab():
			if _, ok := addedAuthors[info.Author]; ok {
				addedAuthors[info.Author].Linters = append(addedAuthors[info.Author].Linters, lc.Name())
			} else {
				ghAuthor := info.Author
				if info.Author == "bosi" {
					ghAuthor = "bosix"
				}

				addedAuthors[info.Author] = &authorDetails{
					Name:    info.Author,
					Linters: []string{lc.Name()},
					Profile: fmt.Sprintf("https://gitlab.com/%[1]s", info.Author),
					Avatar:  fmt.Sprintf("https://github.com/%[1]s.png", ghAuthor),
				}
			}
		}
	}

	authors := slices.SortedFunc(maps.Keys(addedAuthors), func(a string, b string) int {
		return strings.Compare(strings.ToLower(a), strings.ToLower(b))
	})

	var details []*authorDetails

	for _, author := range authors {
		details = append(details, addedAuthors[author])
	}

	return details
}

type authorInfo struct {
	Author string
	Host   string
}

func extractInfo(lc *linter.Config) authorInfo {
	exp := regexp.MustCompile(`https://(github|gitlab)\.com/([^/]+)/.*`)

	switch lc.Name() {
	case "exhaustruct":
		return authorInfo{Author: "xobotyi", Host: hostGitHub}

	case "misspell":
		return authorInfo{Author: "client9", Host: hostGitHub}

	case "fatcontext":
		return authorInfo{Author: "Crocmagnon", Host: hostGitHub}

	case "godoclint":
		return authorInfo{Author: "babakks", Host: hostGitHub}

	case "goprintffuncname":
		return authorInfo{Author: "jirfag", Host: hostGitHub}

	default:
		if strings.HasPrefix(lc.OriginalURL, "https://pkg.go.dev/") {
			return authorInfo{Author: "golang", Host: hostGitHub}
		}

		if !exp.MatchString(lc.OriginalURL) {
			return authorInfo{}
		}

		submatch := exp.FindAllStringSubmatch(lc.OriginalURL, -1)

		info := authorInfo{
			Author: submatch[0][2],
			Host:   submatch[0][1],
		}

		switch info.Author {
		case "gostaticanalysis":
			info.Author = "tenntenn"

		case "go-simpler":
			info.Author = "tmzane"

		case "curioswitch":
			info.Author = "chokoswitch"

		case "GaijinEntertainment":
			info.Author = "xobotyi"

		case "OpenPeeDeeP":
			info.Author = "dixonwille"

		case "golangci":
			return authorInfo{}
		}

		return info
	}
}

func (i authorInfo) FromGitHub() bool {
	return i.Host == hostGitHub
}

func (i authorInfo) FromGitLab() bool {
	return i.Host == hostGitLab
}
