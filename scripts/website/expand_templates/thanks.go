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
	Linters []string
	Profile string
	Avatar  string
}

func getThanksList() string {
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
					Linters: []string{lc.Name()},
					Profile: fmt.Sprintf("[%[1]s](https://github.com/sponsors/%[1]s)", info.Author),
					Avatar:  getGitHubAvatar(info),
				}
			}

		case info.FromGitLab():
			if _, ok := addedAuthors[info.Author]; ok {
				addedAuthors[info.Author].Linters = append(addedAuthors[info.Author].Linters, lc.Name())
			} else {
				getGitHubAvatar(info)
				addedAuthors[info.Author] = &authorDetails{
					Linters: []string{lc.Name()},
					Profile: fmt.Sprintf("[%[1]s](https://gitlab.com/%[1]s)", info.Author),
					Avatar:  getGitHubAvatar(info),
				}
			}
		}
	}

	authors := slices.SortedFunc(maps.Keys(addedAuthors), func(a string, b string) int {
		return strings.Compare(strings.ToLower(a), strings.ToLower(b))
	})

	lines := []string{
		"|Author|Linter(s)|",
		"|---|---|",
	}

	for _, author := range authors {
		lines = append(lines, fmt.Sprintf("|<figure>%s<figcaption>%s</figcaption></figure>|%s|",
			addedAuthors[author].Avatar, addedAuthors[author].Profile, strings.Join(addedAuthors[author].Linters, ", ")))
	}

	return strings.Join(lines, "\n")
}

func getGitHubAvatar(info authorInfo) string {
	return fmt.Sprintf(
		`<img src="https://github.com/%[1]s.png" style="max-width: 100%%;" width="50px;" title="%[1]s" alt="%[1]s" loading="lazy">`,
		info.Author,
	)
}

type authorInfo struct {
	Author string
	Host   string
}

func extractInfo(lc *linter.Config) authorInfo {
	exp := regexp.MustCompile(`https://(github|gitlab)\.com/([^/]+)/.*`)

	switch lc.Name() {
	case "staticcheck":
		return authorInfo{Author: "dominikh", Host: hostGitHub}

	case "misspell":
		return authorInfo{Author: "client9", Host: hostGitHub}

	case "fatcontext":
		return authorInfo{Author: "Crocmagnon", Host: hostGitHub}

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

		case "bosi":
			info.Author = "bosix"

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
