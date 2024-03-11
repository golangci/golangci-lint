package main

import (
	"fmt"
	"sort"
	"strings"

	"golang.org/x/exp/maps"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/lint/lintersdb"
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

		linterURL := lc.OriginalURL
		if lc.Name() == "staticcheck" {
			linterURL = "https://github.com/dominikh/go-tools"
		}

		if author := extractAuthor(linterURL, "https://github.com/"); author != "" && author != "golangci" {
			if _, ok := addedAuthors[author]; ok {
				addedAuthors[author].Linters = append(addedAuthors[author].Linters, lc.Name())
			} else {
				addedAuthors[author] = &authorDetails{
					Linters: []string{lc.Name()},
					Profile: fmt.Sprintf("[%[1]s](https://github.com/sponsors/%[1]s)", author),
					Avatar:  fmt.Sprintf(`<img src="https://github.com/%[1]s.png" alt="%[1]s" style="max-width: 100%%;" width="20px;" />`, author),
				}
			}
		} else if author := extractAuthor(linterURL, "https://gitlab.com/"); author != "" {
			if _, ok := addedAuthors[author]; ok {
				addedAuthors[author].Linters = append(addedAuthors[author].Linters, lc.Name())
			} else {
				addedAuthors[author] = &authorDetails{
					Linters: []string{lc.Name()},
					Profile: fmt.Sprintf("[%[1]s](https://gitlab.com/%[1]s)", author),
				}
			}
		} else {
			continue
		}
	}

	authors := maps.Keys(addedAuthors)
	sort.Slice(authors, func(i, j int) bool {
		return strings.ToLower(authors[i]) < strings.ToLower(authors[j])
	})

	lines := []string{
		"|Author|Linter(s)|",
		"|---|---|",
	}

	for _, author := range authors {
		lines = append(lines, fmt.Sprintf("|%s %s|%s|",
			addedAuthors[author].Avatar, addedAuthors[author].Profile, strings.Join(addedAuthors[author].Linters, ", ")))
	}

	return strings.Join(lines, "\n")
}

func extractAuthor(originalURL, prefix string) string {
	if !strings.HasPrefix(originalURL, prefix) {
		return ""
	}

	return strings.SplitN(strings.TrimPrefix(originalURL, prefix), "/", 2)[0]
}
