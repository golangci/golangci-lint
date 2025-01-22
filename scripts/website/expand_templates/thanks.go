package main

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
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

		if lc.IsDeprecated() && lc.Deprecation.Level > linter.DeprecationWarning {
			continue
		}

		linterURL := extractLinterURL(lc)

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

	authors := slices.SortedFunc(maps.Keys(addedAuthors), func(a string, b string) int {
		return strings.Compare(strings.ToLower(a), strings.ToLower(b))
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

func extractLinterURL(lc *linter.Config) string {
	switch lc.Name() {
	case "staticcheck":
		return "https://github.com/dominikh/go-tools"

	case "depguard":
		return "https://github.com/dixonwille/depguard"

	default:
		if strings.HasPrefix(lc.OriginalURL, "https://github.com/gostaticanalysis/") {
			return "https://github.com/tenntenn/gostaticanalysis"
		}

		if strings.HasPrefix(lc.OriginalURL, "https://github.com/go-simpler/") {
			return "https://github.com/tmzane/go-simpler"
		}

		return lc.OriginalURL
	}
}

func extractAuthor(originalURL, prefix string) string {
	if !strings.HasPrefix(originalURL, prefix) {
		return ""
	}

	return strings.SplitN(strings.TrimPrefix(originalURL, prefix), "/", 2)[0]
}
